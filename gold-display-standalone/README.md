# Gold Display Standalone - 独立金价展示页面

## 项目简介

这是一个专门为 **32位低版本 Fully Kiosk Browser** 设计的独立金价展示页面。使用原生 HTML/CSS/JavaScript 实现，完全兼容旧版浏览器（Android 4.4+, Chrome 30+）。

## 特点

✅ **完全兼容旧浏览器**
- 使用 ES5 语法，不使用 async/await、箭头函数等新特性
- 使用 XMLHttpRequest 而不是 fetch
- 使用传统 CSS，避免使用太新的特性
- 兼容 Android 4.4+ 和 Chrome 30+

✅ **轻量级**
- 无依赖，纯原生实现
- 文件体积小，加载快速
- 适合在低性能设备上运行

✅ **功能完整**
- 视频循环播放
- 产品图片垂直滚动
- 实时金价表格（从 API 获取）
- 时间显示
- 跑马灯效果

## 文件结构

```
gold-display-standalone/
├── index.html          # 主页面
├── styles.css          # 样式文件
├── app.js              # 主逻辑（ES5语法）
├── config.js           # 配置文件（可选）
├── assets/             # 资源文件目录
│   ├── jewelry-video.mp4    # 视频文件
│   ├── product1.jpg         # 产品图片1
│   ├── product2.jpg         # 产品图片2
│   └── ...
└── README.md           # 说明文档
```

## 快速开始

### 1. 准备资源文件

在 `assets` 目录下添加：
- `jewelry-video.mp4` - 珠宝宣传视频（可选）
- `product1.jpg` ~ `product5.jpg` - 产品图片（可选）

如果没有这些文件，页面会显示占位符或使用默认图片。

### 2. 配置 API 地址

编辑 `app.js`，修改 `CONFIG.apiBaseUrl`：

```javascript
var CONFIG = {
    apiBaseUrl: '/api/v1',  // 如果前后端同域，使用相对路径
    // 或者
    // apiBaseUrl: 'http://your-backend-domain.com/api/v1',  // 跨域需要后端支持 CORS
};
```

### 3. 部署

#### 方式一：使用 Web 服务器

将整个 `gold-display-standalone` 目录放到 Web 服务器（如 Nginx、Apache）的网站根目录下。

#### 方式二：使用后端服务

如果后端支持静态文件服务，可以将文件放到后端的静态资源目录。

#### 方式三：直接打开（仅用于测试）

直接双击 `index.html` 打开（注意：这种方式可能无法正常调用 API）。

### 4. 访问

在浏览器或 Fully Kiosk Browser 中打开：
```
http://your-domain.com/gold-display-standalone/
```

## 配置说明

### 修改滚动速度

编辑 `app.js` 中的 `CONFIG`：

```javascript
var CONFIG = {
    scrollSpeed: 30,      // 产品图片滚动速度（像素/秒）
    marqueeSpeed: 100,    // 欢迎语滚动速度（像素/秒）
};
```

### 修改刷新频率

```javascript
var CONFIG = {
    refreshInterval: 5 * 60 * 1000,  // 5分钟，单位：毫秒
};
```

### 修改品牌信息

编辑 `index.html`：

```html
<span class="brand-name">你的品牌名称</span>
```

```html
<div id="marqueeText" class="marquee-text">你的欢迎语</div>
```

## API 接口要求

页面需要调用以下 API 获取金价数据：

**GET** `/api/v1/prices?page=1&page_size=100`

**响应格式**（支持多种格式）：

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
            }
        ]
    }
}
```

或者：

```json
{
    "code": 200,
    "data": [
        {
            "id": 1,
            "name": "足金(9999)",
            "price": 488,
            "fee": 10
        }
    ]
}
```

如果 API 调用失败，页面会自动使用模拟数据。

## 浏览器兼容性

| 浏览器/平台 | 最低版本 | 状态 |
|------------|---------|------|
| Chrome | 30+ | ✅ 完全支持 |
| Android WebView | 4.4+ | ✅ 完全支持 |
| Safari | 7+ | ✅ 完全支持 |
| Firefox | 25+ | ✅ 完全支持 |
| IE | 11+ | ⚠️ 部分支持（CSS 渐变可能有问题） |

## 性能优化

1. **图片优化**：建议使用压缩后的图片，文件大小控制在 200KB 以内
2. **视频优化**：建议使用 H.264 编码的 MP4 格式，分辨率不超过 1080p
3. **缓存策略**：建议在服务器端设置适当的缓存头

## 故障排查

### 1. 视频无法播放

- 检查视频文件路径是否正确
- 确认视频格式为 MP4（H.264 编码）
- 检查浏览器是否支持 HTML5 video

### 2. API 请求失败

- 检查 `apiBaseUrl` 配置是否正确
- 检查后端是否支持 CORS（如果跨域）
- 打开浏览器控制台查看错误信息

### 3. 图片无法显示

- 检查图片路径是否正确
- 确认图片文件存在
- 检查网络连接（如果使用在线图片）

### 4. 动画不流畅

- 检查设备性能
- 降低滚动速度
- 减少图片数量或大小

## 与 Vue 版本的对比

| 特性 | Vue 版本 | 原生版本 |
|------|---------|---------|
| 浏览器兼容性 | Chrome 80+ | Chrome 30+ |
| 文件大小 | ~500KB+ | ~50KB |
| 依赖 | Vue, Element UI | 无 |
| 构建工具 | 需要 npm build | 无需构建 |
| 性能 | 较好 | 优秀（低端设备） |
| 维护性 | 较好 | 一般 |

## 更新日志

### v1.0.0 (2025-01-XX)
- ✅ 初始版本
- ✅ 支持视频播放
- ✅ 支持图片滚动
- ✅ 支持实时金价显示
- ✅ 兼容旧版浏览器

## 许可证

与主项目保持一致。

## 技术支持

如有问题，请查看：
- 主项目 README
- 浏览器控制台错误信息
- 网络请求日志

