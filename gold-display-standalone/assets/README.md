# 资源文件说明

## 目录结构

```
assets/
├── jewelry-video.mp4    # 珠宝宣传视频（循环播放）
├── product1.jpg         # 产品图片1（滚动展示）
├── product2.jpg         # 产品图片2（滚动展示）
├── product3.jpg         # 产品图片3（滚动展示）
├── product4.jpg         # 产品图片4（滚动展示）
├── product5.jpg         # 产品图片5（滚动展示）
└── README.md           # 本说明文件
```

## 资源要求

### 视频文件

- **文件名**: `jewelry-video.mp4`
- **格式**: MP4 (H.264编码)
- **建议尺寸**: 1920x1080 或 1280x720
- **建议时长**: 20-60秒
- **用途**: 在展示页面左侧循环播放

### 产品图片

- **文件名**: `product1.jpg` ~ `product5.jpg`
- **格式**: JPG/PNG
- **建议尺寸**: 400x400 或更大
- **数量**: 至少5张，可以更多
- **用途**: 在展示页面中间垂直滚动展示

## 如何添加资源

### 方式一：直接添加文件

1. 将视频文件命名为 `jewelry-video.mp4`，放入 `assets/` 目录
2. 将产品图片命名为 `product1.jpg`, `product2.jpg` 等，放入 `assets/` 目录

### 方式二：修改配置使用在线资源

编辑 `app.js` 文件，修改配置：

```javascript
var CONFIG = {
    // 视频路径
    videoPath: 'https://your-cdn.com/video.mp4',
    
    // 产品图片URL数组
    productImages: [
        'https://your-cdn.com/product1.jpg',
        'https://your-cdn.com/product2.jpg',
        'https://your-cdn.com/product3.jpg',
        'https://your-cdn.com/product4.jpg',
        'https://your-cdn.com/product5.jpg'
    ]
};
```

## 资源获取建议

### 免费视频资源

- Pexels: https://www.pexels.com/zh-cn/search/videos/jewelry/
- Pixabay: https://pixabay.com/zh/videos/search/jewelry/

### 免费图片资源

- Unsplash: https://unsplash.com/s/photos/jewelry
- Pexels: https://www.pexels.com/zh-cn/search/jewelry/

## 注意事项

1. **文件大小**: 建议视频文件不超过50MB，图片文件不超过2MB
2. **格式兼容**: 确保视频格式为MP4 (H.264)，图片格式为JPG或PNG
3. **路径正确**: 确保资源文件路径与配置中的路径一致
4. **CORS问题**: 如果使用在线资源，确保资源服务器允许跨域访问

## 测试占位资源

如果暂时没有资源文件，可以使用以下占位资源进行测试：

### 在线占位图片

```javascript
productImages: [
    'https://via.placeholder.com/400x400/FFD700/000000?text=Product+1',
    'https://via.placeholder.com/400x400/C0C0C0/000000?text=Product+2',
    'https://via.placeholder.com/400x400/CD7F32/000000?text=Product+3',
    'https://via.placeholder.com/400x400/FFD700/000000?text=Product+4',
    'https://via.placeholder.com/400x400/C0C0C0/000000?text=Product+5'
]
```

### 在线测试视频

可以使用以下测试视频URL：

```javascript
videoPath: 'https://media.w3.org/2010/05/sintel/trailer.mp4'
```

---

**提示**: 资源文件不是必需的，页面会在资源缺失时显示占位符，不影响核心功能测试。

