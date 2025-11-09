// pages/index/index.js
const app = getApp();
const api = require('../../utils/api.js');

Page({
  data: {
    // 价格列表
    priceList: [],
    // 显示的列表（根据分类过滤）
    displayList: [],
    // 更新时间
    updateTime: '',
    // 分类列表
    categories: [
      { id: 'all', name: '全部' },
      { id: 'gold', name: '黄金' },
      { id: 'silver', name: '白银' },
      { id: 'platinum', name: '铂金' },
      { id: 'palladium', name: '钯金' },
      { id: 'foreign', name: '国际' },
      { id: 'exchange', name: '汇率' }
    ],
    // 当前选中的分类
    currentCategory: 'all',
    // 店铺信息
    storeInfo: {},
    // 加载状态
    loading: false,
    // 个人二维码地址（需要替换成实际的二维码图片地址）
    qrcodeUrl: 'https://i.52112.com/icon/jpg/256/20191118/66718/2887582.jpg'
  },

  onLoad(options) {
    console.log('首页加载');
    this.loadData();
    this.setStoreInfo();
  },

  onShow() {
    // 每次显示页面时刷新数据
    this.loadData();
  },

  onPullDownRefresh() {
    // 下拉刷新
    this.loadData();
  },

  // 加载数据
  loadData() {
    wx.showLoading({
      title: '加载中...',
      mask: true
    });

    // 从后端API获取真实数据
    api.getPriceList({ status: 1 })
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

        this.setData({
          priceList: formattedList,
          displayList: formattedList,
          updateTime: priceData.updateTime
        });

        // 缓存数据
        wx.setStorageSync('priceData', priceData);
      })
      .catch((err) => {
        console.error('加载数据失败:', err);
        wx.showToast({
          title: '加载失败，使用缓存数据',
          icon: 'none',
          duration: 2000
        });

        // 失败时使用缓存数据
        const cachedData = wx.getStorageSync('priceData');
        if (cachedData) {
          this.setData({
            priceList: cachedData.list,
            displayList: cachedData.list,
            updateTime: cachedData.updateTime
          });
        }
      })
      .finally(() => {
        wx.hideLoading();
        wx.stopPullDownRefresh();
      });
  },

  // 设置店铺信息
  setStoreInfo() {
    this.setData({
      storeInfo: app.globalData.storeInfo
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


  // 分类切换
  onCategoryChange(e) {
    const categoryId = e.currentTarget.dataset.id;
    this.setData({
      currentCategory: categoryId
    });

    this.filterByCategory(categoryId);
  },

  // 根据分类过滤数据
  filterByCategory(categoryId) {
    const { priceList } = this.data;
    
    if (categoryId === 'all') {
      this.setData({
        displayList: priceList
      });
      return;
    }

    let filteredList = [];

    switch (categoryId) {
      case 'gold':
        filteredList = priceList.filter(item => 
          item.code.includes('gold')
        );
        break;
      case 'silver':
        filteredList = priceList.filter(item => 
          item.code.includes('silver')
        );
        break;
      case 'platinum':
        filteredList = priceList.filter(item => 
          item.code.includes('platinum')
        );
        break;
      case 'palladium':
        filteredList = priceList.filter(item => 
          item.code.includes('palladium')
        );
        break;
      case 'foreign':
        filteredList = priceList.filter(item => 
          item.code.includes('us_') || item.code.includes('london_')
        );
        break;
      case 'exchange':
        filteredList = priceList.filter(item => 
          item.code.includes('usd')
        );
        break;
      default:
        filteredList = priceList;
    }

    this.setData({
      displayList: filteredList
    });
  },

  // 查看价格详情
  onPriceDetail(e) {
    const item = e.currentTarget.dataset.item;
    console.log('查看详情:', item);
    
    wx.showModal({
      title: item.name,
      content: `回购价: ¥${item.buyPrice}\n销售价: ¥${item.sellPrice}\n最高价: ¥${item.highPrice}\n最低价: ¥${item.lowPrice}`,
      confirmText: '使用计算器',
      cancelText: '关闭',
      success(res) {
        if (res.confirm) {
          // 跳转到计算器页面
          wx.switchTab({
            url: '/pages/calculator/calculator'
          });
        }
      }
    });
  },

  // 预览二维码（点击放大）
  onPreviewQRCode() {
    const { qrcodeUrl } = this.data;
    
    wx.previewImage({
      current: qrcodeUrl,
      urls: [qrcodeUrl]
    });
  },

  // 识别二维码（长按）
  onScanQRCode() {
    wx.showToast({
      title: '长按图片识别二维码',
      icon: 'none',
      duration: 2000
    });
  },

  // 拨打电话
  onCallPhone() {
    const phone = this.data.storeInfo.phone;
    wx.makePhoneCall({
      phoneNumber: phone.replace(/-/g, '')
    });
  },

  // 查看地址（导航）
  onViewAddress() {
    const { storeInfo } = this.data;
    
    wx.openLocation({
      latitude: storeInfo.latitude,
      longitude: storeInfo.longitude,
      name: storeInfo.name,
      address: storeInfo.address,
      scale: 18
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

