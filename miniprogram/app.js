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
    // 尝试从后端获取最新数据
    this.updatePriceData()
      .then((priceData) => {
        console.log('价格数据更新成功');
      })
      .catch((err) => {
        console.error('价格数据初始化失败:', err);
        // 如果失败且没有缓存，使用默认数据作为兜底
        const cachedData = wx.getStorageSync('priceData');
        if (!cachedData) {
          console.warn('使用默认数据');
        }
      });
  },

  // 更新价格数据（从后端API获取）
  updatePriceData() {
    const api = require('./utils/api.js');
    
    return new Promise((resolve, reject) => {
      api.getPriceList({ status: 1 }) // 只获取启用的价格
        .then((priceList) => {
          // 转换数据格式，适配小程序需要的格式
          const formattedList = priceList.map(item => ({
            id: item.id,
            code: item.code,
            name: item.name,
            subtitle: item.subtitle || '',
            icon: item.icon || 'Au',
            iconColor: item.icon_color || '#FFD700',
            buyPrice: item.buy_price || 0,
            sellPrice: item.sell_price || 0,
            // 如果没有高低价，使用买卖价作为参考
            highPrice: item.sell_price || item.buy_price || 0,
            lowPrice: item.buy_price || 0
          }));

          const priceData = {
            updateTime: this.formatTime(new Date()),
            list: formattedList
          };

          // 缓存数据
          wx.setStorageSync('priceData', priceData);
          resolve(priceData);
        })
        .catch((err) => {
          console.error('获取价格数据失败:', err);
          // 失败时返回缓存的旧数据
          const cachedData = wx.getStorageSync('priceData');
          if (cachedData) {
            resolve(cachedData);
          } else {
            reject(err);
          }
        });
    });
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

