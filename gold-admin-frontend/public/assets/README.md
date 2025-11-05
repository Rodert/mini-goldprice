# 金价展示页面资源文件说明

## 目录结构

```
public/assets/
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
- **建议尺寸**: 1920x1080 或 1280x720
- **建议时长**: 20-60秒
- **格式**: MP4 (H.264编码)
- **用途**: 在展示页面左侧循环播放

### 产品图片（滚动展示）
- **文件名**: `product1.jpg` ~ `product5.jpg`
- **建议尺寸**: 400x400 或更大
- **格式**: JPG/PNG
- **数量**: 至少5张，可以更多
- **用途**: 在展示页面中间垂直滚动展示

### 如何添加更多图片

如果需要展示更多产品图片，请按以下步骤操作：

1. 将图片文件放入 `public/assets/` 目录，命名为 `product6.jpg`, `product7.jpg` 等
2. 编辑 `src/views/display/index.vue` 文件
3. 找到 `productImages` 数组，添加新图片路径：

```javascript
productImages: [
  '/assets/product1.jpg',
  '/assets/product2.jpg',
  '/assets/product3.jpg',
  '/assets/product4.jpg',
  '/assets/product5.jpg',
  '/assets/product6.jpg',  // 新增
  '/assets/product7.jpg'   // 新增
]
```

## 临时占位图片

如果暂时没有资源文件，页面会显示占位符：
- 视频区域：显示 "视频展示区域" 的提示
- 产品图片：使用现有的占位图片路径（可能显示错误图标）

建议尽快添加实际的视频和图片文件以获得最佳展示效果。

## 访问展示页面

启动前端服务后，访问：`http://localhost:8080/display`

该页面无需登录即可访问。

