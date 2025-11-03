// pages/profile/profile.js
const app = getApp();

Page({
  data: {
    // 用户信息
    userInfo: {},
    // 统计数据
    stats: {
      appointments: 0,
      records: 0,
      calculations: 0
    },
    // 店铺信息
    storeInfo: {}
  },

  onLoad(options) {
    this.loadUserInfo();
    this.loadStats();
    this.loadStoreInfo();
  },

  onShow() {
    // 刷新统计数据
    this.loadStats();
  },

  // 加载用户信息
  loadUserInfo() {
    // 实际开发中应该从后端获取
    const userInfo = wx.getStorageSync('userInfo') || {
      nickName: '微信用户',
      avatarUrl: '',
      id: 'wx_guest'
    };

    this.setData({
      userInfo: userInfo
    });
  },

  // 加载统计数据
  loadStats() {
    // 从缓存获取计算历史数量
    const calcHistory = wx.getStorageSync('calcHistory') || [];
    
    // 实际开发中应该从后端获取真实数据
    this.setData({
      'stats.appointments': 2,
      'stats.records': 5,
      'stats.calculations': calcHistory.length
    });
  },

  // 加载店铺信息
  loadStoreInfo() {
    this.setData({
      storeInfo: app.globalData.storeInfo
    });
  },

  // 页面跳转
  onNavTo(e) {
    const url = e.currentTarget.dataset.url;
    
    if (url.startsWith('/pages/')) {
      // 如果是tabbar页面，使用switchTab
      const tabbarPages = ['/pages/index/index', '/pages/appointment/appointment', '/pages/calculator/calculator', '/pages/profile/profile'];
      
      if (tabbarPages.includes(url)) {
        wx.switchTab({ url });
      } else {
        wx.navigateTo({ url });
      }
    }
  },

  // 查看门店
  onViewStore() {
    const { storeInfo } = this.data;
    
    wx.openLocation({
      latitude: storeInfo.latitude,
      longitude: storeInfo.longitude,
      name: storeInfo.name,
      address: storeInfo.address,
      scale: 18
    });
  },

  // 拨打电话
  onCallPhone() {
    const phone = this.data.storeInfo.phone;
    wx.makePhoneCall({
      phoneNumber: phone.replace(/-/g, '')
    });
  },

  // 显示提示（功能开发中）
  onShowTip() {
    wx.showToast({
      title: '功能开发中',
      icon: 'none'
    });
  },

  // 分享
  onShareAppMessage() {
    return {
      title: '沪金汇 - 诚信经营，价格公道',
      path: '/pages/index/index'
    };
  }
});

