// utils/api.js
// API基础配置
// 根据环境自动切换API地址
const ENV = {
  // 开发环境
  development: 'http://localhost:8080',
  // 生产环境
  production: 'https://gold.javapub.net.cn'
};

// 判断当前环境（可以根据实际情况调整判断逻辑）
// 方式1: 根据编译模式判断（推荐）
// 方式2: 手动设置环境变量
const isProduction = true; // 生产环境改为 true

const API_BASE_URL = isProduction ? ENV.production : ENV.development;

// 统一请求方法
function request(options) {
  return new Promise((resolve, reject) => {
    wx.request({
      url: API_BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        ...options.header
      },
      success: (res) => {
        if (res.statusCode === 200) {
          if (res.data.code === 200) {
            resolve(res.data.data);
          } else {
            reject(new Error(res.data.message || '请求失败'));
          }
        } else {
          reject(new Error(`HTTP错误: ${res.statusCode}`));
        }
      },
      fail: (err) => {
        reject(err);
      },
      complete: () => {
        if (options.complete) {
          options.complete();
        }
      }
    });
  });
}

// API方法
const api = {
  // 获取价格列表
  // params: { shop_id, code, name, status }
  getPriceList(params = {}) {
    let url = '/api/openapi/prices';
    const query = [];
    if (params.shop_id) query.push(`shop_id=${params.shop_id}`);
    if (params.code) query.push(`code=${params.code}`);
    if (params.name) query.push(`name=${params.name}`);
    if (params.status !== undefined) query.push(`status=${params.status}`);
    if (query.length > 0) url += '?' + query.join('&');
    
    return request({ url, method: 'GET' });
  }
};

module.exports = api;

