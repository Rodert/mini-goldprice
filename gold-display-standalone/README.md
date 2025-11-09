# 金价展示页面 - 独立版本

这是一个专门为 **32位低版本 Fully Kiosk Browser** 设计的独立金价展示页面。使用原生 HTML/CSS/JavaScript 实现，完全兼容旧版浏览器（Android 4.4+, Chrome 30+）。

## 📋 特性

- ✅ **原生实现** - 不依赖任何框架（Vue、React等）
- ✅ **兼容性强** - 支持 Android 4.4+, Chrome 30+ 等旧版浏览器
- ✅ **独立运行** - 可直接在浏览器中打开，无需构建工具
- ✅ **自动刷新** - 金价数据每5分钟自动刷新
- ✅ **容错处理** - API失败时自动使用模拟数据
- ✅ **流畅动画** - 产品图片滚动、跑马灯效果
- ✅ **响应式设计** - 适配不同屏幕尺寸

## 🚀 快速开始

### 1. 准备资源文件（可选）

在 `assets/` 目录下添加以下文件：

- `jewelry-video.mp4` - 珠宝宣传视频（循环播放）
- `product1.jpg` ~ `product5.jpg` - 产品图片（滚动展示）

> 如果暂时没有资源文件，页面会显示占位符，不影响功能测试

### 2. 配置API地址

编辑 `app.js` 文件，修改API基础地址：

```javascript
var CONFIG = {
    apiBaseUrl: 'http://your-server:8080/api/v1',  // 修改为你的后端地址
    // ... 其他配置
};
```

### 3. 启动服务

#### 方式一：使用本地服务器（推荐）

```bash
# 使用Python简单HTTP服务器
cd gold-display-standalone
python -m SimpleHTTPServer 8000  # Python 2
# 或
python -m http.server 8000       # Python 3

# 使用Node.js http-server
npx http-server -p 8000

# 使用PHP内置服务器
php -S localhost:8000
```

#### 方式二：直接打开文件

直接在浏览器中打开 `index.html` 文件（注意：某些功能可能受限，如API请求可能被CORS阻止）

### 4. 访问页面

浏览器打开：`http://localhost:8000/index.html`

## 📁 文件结构

```
gold-display-standalone/
├── index.html          # 主页面
├── styles.css          # 样式文件
├── app.js              # JavaScript逻辑
├── assets/             # 资源文件目录
│   ├── jewelry-video.mp4    # 视频文件（需自行添加）
│   ├── product1.jpg          # 产品图片1（需自行添加）
│   ├── product2.jpg          # 产品图片2（需自行添加）
│   └── ...                    # 更多产品图片
└── README.md           # 本说明文件
```

## ⚙️ 配置说明

在 `app.js` 中可以修改以下配置：

```javascript
var CONFIG = {
    // API配置
    apiBaseUrl: 'http://localhost:8080/api/v1',  // API基础地址
    apiTimeout: 10000,                             // API超时时间（毫秒）
    
    // 动画配置
    scrollSpeed: 30,      // 产品图片滚动速度（像素/秒）
    marqueeSpeed: 100,    // 跑马灯滚动速度（像素/秒）
    animationFPS: 60,   // 动画帧率
    
    // 刷新配置
    refreshInterval: 5 * 60 * 1000,  // 数据刷新间隔（毫秒）
    
    // 资源路径
    videoPath: 'assets/jewelry-video.mp4',
    productImages: [
        // 产品图片URL数组
    ]
};
```

## 🎨 页面布局

```
┌──────────────────────────────────────────────┐
│  💎 SineGem 中国珠宝           今日金价       │
├───────────┬──────────┬──────────────────────┤
│           │          │                      │
│   视频    │  产品图  │   品名 | 价格 | 工费 │
│  (循环)   │  (滚动)  │   足金 | 488  | 10   │
│           │          │   Pt950| 388  | 10   │
│           │          │   ...  | ...  | ...  │
├───────────┴──────────┴──────────────────────┤
│ 2025-11-05 星期三 10:30:45 | 欢迎您>>>      │
└──────────────────────────────────────────────┘
```

## 🔧 兼容性说明

### 浏览器兼容性

- ✅ Android 4.4+ (WebView)
- ✅ Chrome 30+
- ✅ Fully Kiosk Browser (所有版本)
- ✅ 其他基于 Chromium 的旧版浏览器

### 技术特性

- 使用原生 JavaScript（ES5语法），不使用ES6+特性
- 使用 XMLHttpRequest 代替 fetch API
- 使用 CSS3 前缀（-webkit-, -moz-, -ms-, -o-）
- 使用 Flexbox 布局（带前缀）
- 避免使用现代CSS特性（如Grid、CSS变量等）

## 📡 API接口

### 获取价格列表

**请求：**
```
GET /api/v1/prices?page=1&page_size=100
```

**响应格式：**
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "name": "足金(9999)",
      "sell_price": 488,
      "fee": 10
    },
    ...
  ]
}
```

或

```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "name": "足金(9999)",
        "sell_price": 488,
        "fee": 10
      },
      ...
    ]
  }
}
```

## 🎯 使用场景

### 适用于：

- ✅ 32位Android设备
- ✅ 旧版平板电脑
- ✅ Fully Kiosk Browser
- ✅ 店铺电子显示屏
- ✅ 门店数字标牌

### 建议设置：

- 使用横屏显示
- 分辨率：1920x1080 或 1280x720
- 全屏模式（F11 或浏览器全屏）
- 保持浏览器标签页激活状态

## 🐛 常见问题

### Q: 视频不显示怎么办？

A: 
1. 检查 `assets/jewelry-video.mp4` 文件是否存在
2. 确保视频格式为 MP4 (H.264编码)
3. 如果视频不存在，页面会显示占位符

### Q: 产品图片不滚动？

A: 
1. 检查 `app.js` 中的 `productImages` 配置
2. 确保图片URL可访问
3. 检查浏览器控制台是否有错误

### Q: 金价数据不显示？

A: 
1. 检查 `app.js` 中的 `apiBaseUrl` 配置是否正确
2. 确保后端服务正常运行
3. 检查浏览器控制台是否有API错误
4. 如果API失败，会自动显示模拟数据

### Q: 页面在旧版浏览器中显示异常？

A: 
1. 确保使用支持HTML5的浏览器
2. 检查浏览器控制台是否有JavaScript错误
3. 尝试使用Chrome 30+或Fully Kiosk Browser

### Q: CORS跨域问题？

A: 
1. 如果API在不同域名，需要后端配置CORS
2. 或者使用同源策略（前后端同一域名）
3. 开发环境可以使用代理服务器

## 🔄 更新日志

### v1.0.0 (2025-01-XX)

- ✅ 初始版本
- ✅ 实现基础金价展示功能
- ✅ 实现产品图片滚动
- ✅ 实现跑马灯效果
- ✅ 实现视频播放
- ✅ 实现自动刷新
- ✅ 兼容旧版浏览器

## 📝 开发说明

### 技术栈

- **HTML5** - 页面结构
- **CSS3** - 样式（带前缀兼容）
- **JavaScript (ES5)** - 逻辑处理
- **XMLHttpRequest** - API请求

### 代码特点

- 使用IIFE（立即执行函数）避免全局污染
- 使用原生DOM API
- 避免使用现代JavaScript特性
- 添加详细的注释和日志

## 📄 许可证

本项目遵循项目主许可证。

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📞 支持

如有问题，请联系技术支持或查看项目文档。

---

**注意**：此版本专门为旧版浏览器设计，如需现代浏览器版本，请使用 `gold-admin-frontend` 中的 Vue 版本。

