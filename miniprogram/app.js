// app.js
App({
  onLaunch() {
    // 小程序启动时执行
    console.log('贵金属回收小程序启动');
    
    // 检查更新
    this.checkUpdate();
    
    // 初始化数据
    this.initData();
  },

  onShow() {
    // 小程序显示时执行
    console.log('小程序显示');
  },

  onHide() {
    // 小程序隐藏时执行
    console.log('小程序隐藏');
  },

  // 检查小程序更新
  checkUpdate() {
    const updateManager = wx.getUpdateManager();
    
    updateManager.onCheckForUpdate((res) => {
      if (res.hasUpdate) {
        console.log('发现新版本');
      }
    });

    updateManager.onUpdateReady(() => {
      wx.showModal({
        title: '更新提示',
        content: '新版本已经准备好，是否重启应用？',
        success(res) {
          if (res.confirm) {
            updateManager.applyUpdate();
          }
        }
      });
    });

    updateManager.onUpdateFailed(() => {
      console.log('新版本下载失败');
    });
  },

  // 初始化数据
  initData() {
    // 如果没有缓存的价格数据，设置默认值
    const priceData = wx.getStorageSync('priceData');
    if (!priceData) {
      this.updatePriceData();
    }
  },

  // 更新价格数据（模拟数据，实际应该从后端获取）
  updatePriceData() {
    const priceData = {
      updateTime: this.formatTime(new Date()),
      list: [
        {
          id: 1,
          code: 'gold_9999',
          name: '黄金9999',
          subtitle: 'Au9999 · 千足金',
          icon: 'Au',
          iconColor: '#FFD700',
          buyPrice: 558.50,
          sellPrice: 568.80,
          highPrice: 570.20,
          lowPrice: 556.30
        },
        {
          id: 2,
          code: 'gold_td',
          name: '黄金T+D',
          subtitle: '上海金交所',
          icon: 'Au',
          iconColor: '#FFD700',
          buyPrice: 559.20,
          sellPrice: 569.50,
          highPrice: 571.00,
          lowPrice: 557.00
        },
        {
          id: 3,
          code: 'gold',
          name: '黄金',
          subtitle: '标准金条',
          icon: 'Au',
          iconColor: '#FFD700',
          buyPrice: 557.00,
          sellPrice: 567.30,
          highPrice: 569.50,
          lowPrice: 555.20
        },
        {
          id: 4,
          code: 'silver',
          name: '白银',
          subtitle: 'Ag9999 · 千足银',
          icon: 'Ag',
          iconColor: '#C0C0C0',
          buyPrice: 7.15,
          sellPrice: 7.45,
          highPrice: 7.50,
          lowPrice: 7.10
        },
        {
          id: 5,
          code: 'palladium',
          name: '钯金',
          subtitle: 'Pd990 · 钯金饰品',
          icon: 'Pd',
          iconColor: '#CED0DD',
          buyPrice: 715.80,
          sellPrice: 735.20,
          highPrice: 738.00,
          lowPrice: 712.50
        },
        {
          id: 6,
          code: 'platinum',
          name: '铂金',
          subtitle: 'Pt950 · 白金',
          icon: 'Pt',
          iconColor: '#E5E4E2',
          buyPrice: 208.30,
          sellPrice: 218.50,
          highPrice: 220.00,
          lowPrice: 206.50
        },
        {
          id: 7,
          code: 'silver_td',
          name: '白银T+D',
          subtitle: '上海金交所',
          icon: 'Ag',
          iconColor: '#C0C0C0',
          buyPrice: 7.20,
          sellPrice: 7.50,
          highPrice: 7.55,
          lowPrice: 7.15
        },
        {
          id: 8,
          code: 'us_gold',
          name: '美黄金',
          subtitle: 'COMEX黄金',
          icon: 'Au',
          iconColor: '#DAA520',
          buyPrice: 2650.30,
          sellPrice: 2665.80,
          highPrice: 2670.00,
          lowPrice: 2645.00
        },
        {
          id: 9,
          code: 'us_silver',
          name: '美白银',
          subtitle: 'COMEX白银',
          icon: 'Ag',
          iconColor: '#B0B0B0',
          buyPrice: 31.25,
          sellPrice: 31.85,
          highPrice: 32.00,
          lowPrice: 31.10
        },
        {
          id: 10,
          code: 'london_gold',
          name: '伦敦金',
          subtitle: 'LBMA黄金',
          icon: 'Au',
          iconColor: '#FFD700',
          buyPrice: 2648.50,
          sellPrice: 2664.20,
          highPrice: 2668.00,
          lowPrice: 2643.30
        },
        {
          id: 11,
          code: 'london_silver',
          name: '伦敦银',
          subtitle: 'LBMA白银',
          icon: 'Ag',
          iconColor: '#C0C0C0',
          buyPrice: 31.18,
          sellPrice: 31.78,
          highPrice: 31.95,
          lowPrice: 31.05
        },
        {
          id: 12,
          code: 'usd_cny',
          name: '美元兑换',
          subtitle: 'USD/CNY',
          icon: '$',
          iconColor: '#4CAF50',
          buyPrice: 7.25,
          sellPrice: 7.30,
          highPrice: 7.32,
          lowPrice: 7.23
        }
      ]
    };

    wx.setStorageSync('priceData', priceData);
    return priceData;
  },

  // 格式化时间
  formatTime(date) {
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    const hour = date.getHours();
    const minute = date.getMinutes();
    const second = date.getSeconds();

    return `${year}-${this.formatNumber(month)}-${this.formatNumber(day)} ${this.formatNumber(hour)}:${this.formatNumber(minute)}:${this.formatNumber(second)}`;
  },

  formatNumber(n) {
    n = n.toString();
    return n[1] ? n : '0' + n;
  },

  // 全局数据
  globalData: {
    userInfo: null,
    storeInfo: {
      name: '天雅珠宝城',
      brandLogo: 'HJH',
      address: '北京市东城区法华寺街129(天坛东门地铁站A8西北口步行500米)',
      phone: '400-888-8888',
      mobile: '138-0000-0000',
      hours: '周一至周日 09:00 - 18:00',
      latitude: 39.89,
      longitude: 116.42
    }
  }
});

