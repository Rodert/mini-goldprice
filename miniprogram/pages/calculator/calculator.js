// pages/calculator/calculator.js
const app = getApp();

Page({
  data: {
    // 当前Tab
    currentTab: 0,
    // 贵金属列表
    metalList: [],
    // 选中的金属
    selectedMetal: {},
    selectedMetalIndex: 0,
    // 重量
    weight: '',
    // 成色
    purity: '99.9',
    // 计算结果
    result: 0,
    // 计算公式
    formula: '',
    // 更新时间
    updateTime: '',
    // 快捷重量
    quickWeights: [10, 50, 100],
    // 快捷成色
    quickPurities: [99.9, 99.5, 99],
    // 计算历史
    historyList: []
  },

  onLoad(options) {
    this.loadMetalList();
    this.loadHistory();
  },

  onShow() {
    // 刷新价格数据
    this.loadMetalList();
  },

  // 加载贵金属列表
  loadMetalList() {
    const priceData = wx.getStorageSync('priceData');
    if (priceData) {
      this.setData({
        metalList: priceData.list,
        updateTime: priceData.updateTime
      });

      // 如果还没有选择金属，默认选择第一个
      if (!this.data.selectedMetal.name && priceData.list.length > 0) {
        this.setData({
          selectedMetal: priceData.list[0],
          selectedMetalIndex: 0
        });
      }
    }
  },

  // 加载历史记录
  loadHistory() {
    const history = wx.getStorageSync('calcHistory') || [];
    this.setData({
      historyList: history
    });
  },

  // 保存历史记录
  saveHistory(item) {
    let history = wx.getStorageSync('calcHistory') || [];
    
    // 添加新记录到开头
    history.unshift({
      id: Date.now(),
      metalName: item.metalName,
      metalCode: item.metalCode,
      buyPrice: item.buyPrice,
      weight: item.weight,
      purity: item.purity,
      result: item.result,
      formula: item.formula,
      time: this.formatTime(new Date())
    });

    // 只保留最近20条记录
    if (history.length > 20) {
      history = history.slice(0, 20);
    }

    wx.setStorageSync('calcHistory', history);
    this.loadHistory();
  },

  // Tab切换
  onTabChange(e) {
    const index = parseInt(e.currentTarget.dataset.index);
    this.setData({
      currentTab: index
    });

    if (index === 1) {
      this.loadHistory();
    }
  },

  // 选择金属
  onMetalChange(e) {
    const index = parseInt(e.detail.value);
    this.setData({
      selectedMetal: this.data.metalList[index],
      selectedMetalIndex: index,
      result: 0
    });
  },

  // 输入重量
  onWeightInput(e) {
    this.setData({
      weight: e.detail.value,
      result: 0
    });
  },

  // 输入成色
  onPurityInput(e) {
    this.setData({
      purity: e.detail.value,
      result: 0
    });
  },

  // 快捷选择重量
  onQuickWeight(e) {
    const weight = e.currentTarget.dataset.weight;
    this.setData({
      weight: weight.toString(),
      result: 0
    });
  },

  // 快捷选择成色
  onQuickPurity(e) {
    const purity = e.currentTarget.dataset.purity;
    this.setData({
      purity: purity.toString(),
      result: 0
    });
  },

  // 计算
  onCalculate() {
    const { selectedMetal, weight, purity } = this.data;

    // 验证
    if (!selectedMetal.name) {
      wx.showToast({
        title: '请选择贵金属品种',
        icon: 'none'
      });
      return;
    }

    if (!weight || parseFloat(weight) <= 0) {
      wx.showToast({
        title: '请输入正确的重量',
        icon: 'none'
      });
      return;
    }

    const purityValue = purity ? parseFloat(purity) : 100;
    if (purityValue <= 0 || purityValue > 100) {
      wx.showToast({
        title: '成色系数范围：0-100',
        icon: 'none'
      });
      return;
    }

    // 计算
    const weightNum = parseFloat(weight);
    const purityNum = purityValue / 100;
    const buyPrice = selectedMetal.buyPrice;
    const result = (weightNum * buyPrice * purityNum).toFixed(2);
    const formula = `${weightNum}克 × ¥${buyPrice}/克 × ${purityValue}% = ¥${result}`;

    this.setData({
      result: result,
      formula: formula
    });

    // 保存历史
    this.saveHistory({
      metalName: selectedMetal.name,
      metalCode: selectedMetal.code,
      buyPrice: buyPrice,
      weight: weightNum,
      purity: purityValue,
      result: result,
      formula: formula
    });

    wx.showToast({
      title: '计算完成',
      icon: 'success'
    });
  },

  // 使用历史配置
  onUseHistory(e) {
    const item = e.currentTarget.dataset.item;
    
    // 找到对应的金属
    const metalIndex = this.data.metalList.findIndex(m => m.code === item.metalCode);
    
    if (metalIndex >= 0) {
      this.setData({
        currentTab: 0,
        selectedMetal: this.data.metalList[metalIndex],
        selectedMetalIndex: metalIndex,
        weight: item.weight.toString(),
        purity: item.purity.toString()
      });

      wx.showToast({
        title: '已应用历史配置',
        icon: 'success'
      });
    }
  },

  // 删除历史记录
  onDeleteHistory(e) {
    const id = e.currentTarget.dataset.id;
    
    wx.showModal({
      title: '确认删除',
      content: '确定要删除这条计算记录吗？',
      success: (res) => {
        if (res.confirm) {
          let history = wx.getStorageSync('calcHistory') || [];
          history = history.filter(item => item.id !== id);
          wx.setStorageSync('calcHistory', history);
          this.loadHistory();

          wx.showToast({
            title: '删除成功',
            icon: 'success'
          });
        }
      }
    });
  },

  // 格式化时间
  formatTime(date) {
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hour = date.getHours();
    const minute = date.getMinutes();

    return `${year}-${this.formatNumber(month)}-${this.formatNumber(day)} ${this.formatNumber(hour)}:${this.formatNumber(minute)}`;
  },

  formatNumber(n) {
    n = n.toString();
    return n[1] ? n : '0' + n;
  }
});

