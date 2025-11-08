# 资源文件说明

## 必需文件

### 视频文件（可选）
- **文件名**: `jewelry-video.mp4`
- **格式**: MP4 (H.264 编码)
- **建议分辨率**: 1920x1080 或更低
- **建议大小**: < 50MB
- **说明**: 如果没有此文件，页面会显示占位符

## 产品图片（可选）

建议添加 5 张产品图片：
- `product1.jpg`
- `product2.jpg`
- `product3.jpg`
- `product4.jpg`
- `product5.jpg`

### 图片要求
- **格式**: JPG 或 PNG
- **建议尺寸**: 400x600 或类似比例
- **建议大小**: < 200KB/张
- **说明**: 如果没有这些文件，页面会使用默认的在线图片

## 使用本地图片

如果使用本地图片，需要修改 `app.js` 中的 `CONFIG.productImages`：

```javascript
var CONFIG = {
    productImages: [
        'assets/product1.jpg',
        'assets/product2.jpg',
        'assets/product3.jpg',
        'assets/product4.jpg',
        'assets/product5.jpg'
    ]
};
```

## 使用在线图片

也可以使用在线图片 URL，但需要注意：
- 确保图片服务器支持跨域访问
- 确保网络连接稳定
- 考虑加载速度

