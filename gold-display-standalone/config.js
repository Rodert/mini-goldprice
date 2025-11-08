/**
 * 配置文件
 * 可以根据实际部署环境修改这些配置
 */

var DISPLAY_CONFIG = {
    // API 基础地址
    // 如果前端和后端不在同一域名，需要修改这里
    apiBaseUrl: '/api/v1',

    // 如果使用 CORS，可以设置完整 URL
    // apiBaseUrl: 'http://your-backend-domain.com/api/v1',

    // 产品图片列表
    // 可以替换为本地图片路径，如: 'assets/product1.jpg'
    productImages: [
        'assets/product1.jpg',
        'assets/product2.jpg',
        'assets/product3.jpg',
        'assets/product4.jpg',
        'assets/product5.jpg'
    ],

    // 视频文件路径
    videoPath: 'assets/jewelry-video.mp4',

    // 滚动速度（像素/秒）
    scrollSpeed: 30,
    marqueeSpeed: 100,

    // 价格刷新间隔（毫秒）
    refreshInterval: 5 * 60 * 1000,  // 5分钟

    // 品牌信息
    brandName: 'SineGem 中国珠宝',
    welcomeText: '中国珠宝欢迎您！！！'
};

