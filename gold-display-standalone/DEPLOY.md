# 部署指南 - Gold Display Standalone

## 快速部署

### 方式一：使用 Nginx

1. 将 `gold-display-standalone` 目录复制到服务器
2. 配置 Nginx：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/gold-display-standalone;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理（如果需要）
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 方式二：使用后端静态文件服务

如果后端是 Go 语言，可以在 `main.go` 中添加：

```go
// 静态文件服务
router.Static("/gold-display", "./gold-display-standalone")
```

然后访问：`http://your-domain.com/gold-display/`

### 方式三：直接放在后端 public 目录

将 `gold-display-standalone` 目录内容复制到后端的静态文件目录。

## 配置 API 地址

### 同域部署

如果前后端在同一域名下，使用相对路径：

```javascript
var CONFIG = {
    apiBaseUrl: '/api/v1',
};
```

### 跨域部署

如果前后端在不同域名，需要：

1. 修改 `app.js` 中的 `apiBaseUrl`：

```javascript
var CONFIG = {
    apiBaseUrl: 'http://your-backend-domain.com/api/v1',
};
```

2. 确保后端支持 CORS（跨域资源共享）

## 测试

1. 在浏览器中打开页面
2. 按 F12 打开开发者工具
3. 检查控制台是否有错误
4. 检查网络请求是否正常

## 故障排查

### 1. 页面空白

- 检查文件路径是否正确
- 检查浏览器控制台错误
- 确认所有文件都已上传

### 2. API 请求失败

- 检查 `apiBaseUrl` 配置
- 检查后端服务是否运行
- 检查 CORS 配置（如果跨域）

### 3. 视频/图片无法显示

- 检查文件路径
- 检查文件是否存在
- 检查文件权限

## 性能优化建议

1. **启用 Gzip 压缩**（Nginx 配置）：
```nginx
gzip on;
gzip_types text/css application/javascript image/svg+xml;
```

2. **设置缓存头**：
```nginx
location ~* \.(jpg|jpeg|png|gif|mp4)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

3. **使用 CDN**：将静态资源放到 CDN 上

## 安全建议

1. 如果使用跨域 API，确保后端有适当的 CORS 配置
2. 不要在代码中硬编码敏感信息
3. 使用 HTTPS（生产环境）

