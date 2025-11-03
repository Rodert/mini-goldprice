// pages/appointment/appointment.js
Page({
  data: {
    // 表单数据
    form: {
      metalIndex: 0,
      metalName: '',
      serviceType: 'store',
      date: '',
      time: '',
      timeIndex: 0,
      name: '',
      phone: '',
      address: '',
      note: ''
    },
    // 贵金属列表
    metalList: [],
    // 最小日期（今天）
    minDate: '',
    // 时间段
    timeSlots: [
      '09:00-10:00',
      '10:00-11:00',
      '11:00-12:00',
      '13:00-14:00',
      '14:00-15:00',
      '15:00-16:00',
      '16:00-17:00',
      '17:00-18:00'
    ]
  },

  onLoad(options) {
    this.initData();
    this.loadMetalList();
  },

  // 初始化数据
  initData() {
    const today = new Date();
    const dateStr = this.formatDate(today);
    
    this.setData({
      minDate: dateStr,
      'form.date': dateStr
    });
  },

  // 加载贵金属列表
  loadMetalList() {
    const priceData = wx.getStorageSync('priceData');
    if (priceData) {
      this.setData({
        metalList: priceData.list
      });
    }
  },

  // 选择品种
  onMetalChange(e) {
    const index = parseInt(e.detail.value);
    this.setData({
      'form.metalIndex': index,
      'form.metalName': this.data.metalList[index].name
    });
  },

  // 选择服务方式
  onServiceTypeChange(e) {
    const type = e.currentTarget.dataset.type;
    this.setData({
      'form.serviceType': type
    });
  },

  // 选择日期
  onDateChange(e) {
    this.setData({
      'form.date': e.detail.value
    });
  },

  // 选择时间
  onTimeChange(e) {
    const index = parseInt(e.detail.value);
    this.setData({
      'form.timeIndex': index,
      'form.time': this.data.timeSlots[index]
    });
  },

  // 输入姓名
  onNameInput(e) {
    this.setData({
      'form.name': e.detail.value
    });
  },

  // 输入电话
  onPhoneInput(e) {
    this.setData({
      'form.phone': e.detail.value
    });
  },

  // 输入地址
  onAddressInput(e) {
    this.setData({
      'form.address': e.detail.value
    });
  },

  // 输入备注
  onNoteInput(e) {
    this.setData({
      'form.note': e.detail.value
    });
  },

  // 提交预约
  onSubmit() {
    const { form } = this.data;

    // 验证
    if (!form.metalName) {
      wx.showToast({ title: '请选择回收品种', icon: 'none' });
      return;
    }

    if (!form.date) {
      wx.showToast({ title: '请选择预约日期', icon: 'none' });
      return;
    }

    if (!form.time) {
      wx.showToast({ title: '请选择预约时间', icon: 'none' });
      return;
    }

    if (!form.name) {
      wx.showToast({ title: '请输入姓名', icon: 'none' });
      return;
    }

    if (!form.phone || !/^1[3-9]\d{9}$/.test(form.phone)) {
      wx.showToast({ title: '请输入正确的手机号', icon: 'none' });
      return;
    }

    if (form.serviceType === 'home' && !form.address) {
      wx.showToast({ title: '请输入详细地址', icon: 'none' });
      return;
    }

    // 显示加载
    wx.showLoading({ title: '提交中...' });

    // 模拟提交（实际应该调用API）
    setTimeout(() => {
      wx.hideLoading();
      
      wx.showModal({
        title: '预约成功',
        content: '您的预约已提交，工作人员将在30分钟内与您联系确认。',
        showCancel: false,
        success: () => {
          // 清空表单
          this.setData({
            'form.metalIndex': 0,
            'form.metalName': '',
            'form.serviceType': 'store',
            'form.time': '',
            'form.timeIndex': 0,
            'form.name': '',
            'form.phone': '',
            'form.address': '',
            'form.note': ''
          });
        }
      });
    }, 1000);
  },

  // 格式化日期
  formatDate(date) {
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const day = date.getDate();
    
    return `${year}-${this.formatNumber(month)}-${this.formatNumber(day)}`;
  },

  formatNumber(n) {
    n = n.toString();
    return n[1] ? n : '0' + n;
  }
});

