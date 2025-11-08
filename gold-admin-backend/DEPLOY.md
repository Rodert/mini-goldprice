# 部署指南

## 开发环境部署

### 1. 后端部署

#### 前置要求
- Go 1.20+
- Make（可选）

#### 步骤

```bash
# 1. 克隆项目
cd gold-admin-backend

# 2. 安装依赖
go mod tidy

# 3. 初始化数据库和管理员账号
go run tools/create_admin.go

# 4. 启动服务
go run main.go
# 或使用 Make
make run
```

服务将在 `http://localhost:8080` 启动

默认管理员账号:
- 用户名: `admin`
- 密码: `admin123`

### 2. 前端部署（待开发）

前端项目将使用 Vue 2 + Element UI 开发，部署步骤见前端项目文档。

---

## 生产环境部署

### 方案一：单体部署（推荐）

使用 Go 的 `embed` 功能将前端打包文件嵌入到后端二进制文件中，实现单文件部署。

#### 1. 编译后端

```bash
# 编译生成二进制文件
make build
# 或
go build -ldflags="-w -s" -o gold-admin main.go
```

#### 2. 准备配置文件

```bash
# 创建生产环境配置
cp config/config.yaml config/prod.yaml

# 修改生产环境配置
vim config/prod.yaml
```

生产环境配置示例:
```yaml
server:
  port: 8080
  mode: release  # 生产模式

database:
  type: sqlite
  path: /data/gold_admin.db  # 使用绝对路径
  max_idle_conns: 10
  max_open_conns: 100

jwt:
  secret: your-very-secure-secret-key-here  # 修改为安全的密钥
  expire: 168

log:
  level: warn
  path: /var/log/gold-admin
```

#### 3. 启动服务

```bash
# 创建必要目录
mkdir -p /data /var/log/gold-admin

# 启动服务
./gold-admin

# 后台运行
nohup ./gold-admin > /var/log/gold-admin/app.log 2>&1 &
```

#### 4. 使用 systemd 管理（Linux）

创建服务文件 `/etc/systemd/system/gold-admin.service`:

```ini
[Unit]
Description=Gold Admin Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/gold-admin
ExecStart=/opt/gold-admin/gold-admin
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

管理服务:
```bash
# 启动服务
sudo systemctl start gold-admin

# 设置开机自启
sudo systemctl enable gold-admin

# 查看状态
sudo systemctl status gold-admin

# 查看日志
sudo journalctl -u gold-admin -f
```

### 方案二：前后端分离部署

#### 1. 后端部署

同方案一，但需要修改 CORS 配置允许前端域名。

修改 `middleware/cors.go`:
```go
AllowOrigins: []string{
    "https://admin.yourdomain.com",  // 生产环境前端地址
},
```

#### 2. 前端部署

使用 Nginx 部署前端静态文件:

```nginx
server {
    listen 80;
    server_name admin.yourdomain.com;

    root /var/www/gold-admin-frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # 代理 API 请求
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

---

## Docker 部署

### 1. 创建 Dockerfile

```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags="-w -s" -o gold-admin main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /root/
COPY --from=builder /app/gold-admin .
COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./gold-admin"]
```

### 2. 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  gold-admin:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/root/data
      - ./logs:/root/logs
      - ./config/prod.yaml:/root/config/config.yaml
    restart: unless-stopped
    environment:
      - TZ=Asia/Shanghai
```

### 3. 启动服务

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

---

## 性能优化

### 1. 数据库优化

对于生产环境，建议使用 MySQL/PostgreSQL 替代 SQLite：

修改 `config/prod.yaml`:
```yaml
database:
  type: mysql
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: gold_admin
  charset: utf8mb4
```

修改 `models/init.go` 添加 MySQL 支持：
```go
import (
    "gorm.io/driver/mysql"
)

// 在 InitDB 函数中添加
case "mysql":
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
        cfg.Username, cfg.Password, cfg.Host, cfg.Port, 
        cfg.Database, cfg.Charset)
    dialector = mysql.Open(dsn)
```

### 2. 添加缓存（可选）

对于频繁查询的数据（如菜单、价格），可以添加 Redis 缓存。

### 3. 日志分级

生产环境将日志级别设置为 `warn` 或 `error`，减少 I/O 操作。

---

## 安全建议

### 1. JWT Secret

生产环境必须修改 JWT 密钥：
```yaml
jwt:
  secret: use-a-very-long-and-random-secret-key-here
```

生成安全密钥:
```bash
openssl rand -base64 32
```

### 2. 修改默认管理员密码

首次部署后立即修改默认管理员密码。

### 3. HTTPS

生产环境务必使用 HTTPS，可以使用 Nginx 反向代理并配置 SSL 证书：

```nginx
server {
    listen 443 ssl http2;
    server_name admin.yourdomain.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name admin.yourdomain.com;
    return 301 https://$server_name$request_uri;
}
```

### 4. 数据备份

定期备份数据库文件：

```bash
# SQLite 备份
cp /data/gold_admin.db /backup/gold_admin_$(date +%Y%m%d).db

# MySQL 备份
mysqldump -u root -p gold_admin > /backup/gold_admin_$(date +%Y%m%d).sql
```

设置定时任务：
```cron
# 每天凌晨 2 点备份
0 2 * * * /path/to/backup.sh
```

---

## 监控和维护

### 1. 健康检查

```bash
# 检查服务是否运行
curl http://localhost:8080/ping
```

### 2. 日志监控

```bash
# 实时查看日志
tail -f /var/log/gold-admin/app.log

# 查看错误日志
grep "ERROR" /var/log/gold-admin/app.log
```

### 3. 性能监控

可以集成 Prometheus + Grafana 进行性能监控。

---

## 故障排查

### 1. 服务无法启动

检查：
- 配置文件路径是否正确
- 数据库连接是否正常
- 端口是否被占用

```bash
# 检查端口占用
lsof -i :8080

# 查看详细日志
./gold-admin 2>&1 | tee error.log
```

### 2. 数据库连接失败

- 检查数据库配置
- 确保数据库服务已启动
- 检查防火墙规则

### 3. Token 过期

默认 Token 有效期为 7 天，可在配置文件中修改：
```yaml
jwt:
  expire: 168  # 小时
```

---

## 升级指南

### 1. 备份数据

```bash
# 备份数据库
cp /data/gold_admin.db /backup/

# 备份配置
cp config/config.yaml /backup/
```

### 2. 更新代码

```bash
git pull origin main
go mod tidy
```

### 3. 数据库迁移

如果有新的表结构变更，GORM 会自动迁移。但建议先在测试环境验证。

### 4. 重启服务

```bash
sudo systemctl restart gold-admin
```







