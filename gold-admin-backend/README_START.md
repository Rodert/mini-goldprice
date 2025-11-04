# 🎉 后台管理系统 - 完整启动指南

## 📋 项目概览

**沪金汇后台管理系统**已经开发完成，包含：

✅ **后端服务** (Go + Gin + GORM + SQLite)
- 端口：http://localhost:8080
- 完整的 RESTful API
- JWT 认证
- RBAC 权限管理

✅ **前端管理界面** (Vue 2 + Element UI)
- 端口：http://localhost:9527
- 响应式管理界面
- 动态路由
- 权限控制

---

## 🚀 一键启动（5分钟）

### 第一步：启动后端服务

```bash
# 进入后端目录
cd gold-admin-backend

# 启动服务（首次启动会自动初始化数据库）
go run main.go
```

**启动成功后会看到：**
```
✓ 创建默认角色完成
✓ 创建默认菜单完成
✓ 为超级管理员分配菜单完成
✓ 创建默认管理员完成
✓ 为管理员分配角色完成
✓ 创建默认店铺数据完成
✓ 创建默认价格数据完成
========================================
初始化数据完成！
默认管理员账号：admin
默认密码：admin123
========================================
服务器启动成功，监听地址: :8080
```

### 第二步：启动前端界面

```bash
# 新开一个终端窗口
cd gold-admin-frontend

# 安装依赖（仅首次需要）
npm install

# 启动开发服务器
npm run dev
```

**启动成功后会看到：**
```
App running at:
- Local:   http://localhost:9527/
```

### 第三步：访问系统

浏览器打开：**http://localhost:9527**

**登录信息：**
- 账号：`admin`
- 密码：`admin123`

---

## 🎯 已实现的功能模块

### ✅ 认证模块
- [x] 管理员登录
- [x] JWT Token 认证
- [x] 用户信息获取
- [x] 菜单权限加载
- [x] 登出

### ✅ 首页看板
- [x] 今日预约统计
- [x] 待处理预约数
- [x] 本月预约统计
- [x] 店铺/价格/用户统计
- [x] 最近动态列表
- [x] 预约趋势图表

### ✅ 系统管理

#### 用户管理
- [x] 用户列表（分页、搜索）
- [x] 创建用户
- [x] 编辑用户信息
- [x] 分配角色
- [x] 启用/禁用用户
- [x] 重置密码
- [x] 删除用户

#### 角色管理
- [x] 角色列表
- [x] 创建角色
- [x] 编辑角色
- [x] 分配菜单权限
- [x] 删除角色

#### 菜单管理
- [x] 菜单树形展示
- [x] 创建菜单（目录/菜单/按钮）
- [x] 编辑菜单
- [x] 删除菜单
- [x] 排序

### ✅ 业务管理

#### 价格管理
- [x] 价格列表（卡片展示）
- [x] 新增贵金属品种
- [x] 编辑价格（基础价 ± 差价模式）
- [x] 删除品种
- [x] 启用/禁用

#### 店铺管理
- [x] 店铺列表
- [x] 创建店铺
- [x] 编辑店铺信息（地址、电话、营业时间、坐标）
- [x] 启用/禁用店铺
- [x] 删除店铺

#### 预约管理
- [x] 预约列表（分页、筛选）
- [x] 预约详情
- [x] 更新预约状态（待确认/已确认/已完成/已取消）
- [x] 添加管理员备注
- [x] 删除预约
- [x] 预约统计

---

## 📊 核心数据说明

### 默认管理员账号
```
用户名：admin
密码：admin123
角色：超级管理员（拥有所有权限）
```

### 预设角色（5个）
1. **超级管理员** (super_admin) - 所有权限
2. **总部店长** (head_manager) - 管理所有店铺
3. **单店店长** (shop_manager) - 管理单个店铺
4. **店员** (shop_staff) - 只读权限
5. **财务** (finance) - 数据查看/导出

### 默认菜单结构
```
├─ 首页
├─ 业务管理
│  ├─ 价格管理
│  ├─ 预约管理
│  └─ 店铺管理
└─ 系统管理
   ├─ 用户管理
   ├─ 角色管理
   └─ 菜单管理
```

### 默认店铺数据（2个）
- **沪金汇总店** (shop1) - 上海市黄浦区
- **沪金汇浦东店** (shop2) - 上海市浦东新区

### 默认价格数据（3个）
- **黄金9999** (Au9999) - 回购价：550.00 元/克
- **黄金999** (Au999) - 回购价：543.00 元/克
- **白银999** (Ag999) - 回购价：6.00 元/克

---

## 🔧 API 接口地址

**Base URL**: `http://localhost:8080/api`

### 主要接口

#### 认证
- `POST /login` - 登录
- `POST /logout` - 登出
- `GET /user/info` - 获取用户信息

#### 首页看板
- `GET /dashboard/stats` - 统计数据
- `GET /dashboard/activities` - 最近动态
- `GET /dashboard/trend` - 预约趋势

#### 系统管理
- `GET /users` - 用户列表
- `GET /roles` - 角色列表
- `GET /menus` - 菜单树

#### 业务管理
- `GET /prices` - 价格列表
- `GET /shops` - 店铺列表
- `GET /appointments` - 预约列表

更多接口详见：[API.md](./API.md)

---

## 🎨 核心技术亮点

### 1. 价格管理（创新设计）⭐
```
基础价格（市场价） ± 差价 = 最终价格

示例：
- 基础价：560.00 元/克（随市场波动）
- 回购差价：-10.00 元/克（商家利润空间）
- 销售差价：+15.00 元/克（商家定价策略）

→ 回购价 = 560 - 10 = 550.00 元/克
→ 销售价 = 560 + 15 = 575.00 元/克

优势：
✓ 灵活调整商家策略
✓ 基础价和差价分开管理
✓ 支持不同客户差价
```

### 2. RBAC 权限系统
- 用户 ↔ 角色 ↔ 菜单（多对多关系）
- 页面级权限（已实现）
- 按钮级权限（已预留字段）
- 动态菜单加载

### 3. SQLite 零配置数据库
- 单文件存储：`./data/gold_admin.db`
- 零安装、零配置
- 便于备份和迁移
- 可无缝切换到 MySQL

### 4. JWT 无状态认证
- Token 过期时间：7天
- 刷新机制
- 前端自动携带

---

## 📂 项目结构

```
gold-admin-backend/           # 后端（Go）
├── api/v1/                  # API 处理器
│   ├── auth.go             # 认证（登录/登出）✅
│   ├── dashboard.go        # 首页看板 ✅
│   ├── user.go             # 用户管理 ✅
│   ├── role.go             # 角色管理 ✅
│   ├── menu.go             # 菜单管理 ✅
│   ├── price.go            # 价格管理 ✅
│   ├── shop.go             # 店铺管理 ✅
│   └── appointment.go      # 预约管理 ✅
├── models/                  # 数据模型
│   ├── init.go             # 数据库初始化 ✅
│   ├── admin_user.go       # 用户模型 ✅
│   ├── role.go             # 角色模型 ✅
│   ├── menu.go             # 菜单模型 ✅
│   ├── price.go            # 价格模型 ✅
│   ├── shop.go             # 店铺模型 ✅
│   └── appointment.go      # 预约模型 ✅
├── router/                  # 路由配置 ✅
├── middleware/              # 中间件
│   ├── cors.go             # 跨域 ✅
│   └── jwt.go              # JWT认证 ✅
├── utils/                   # 工具函数 ✅
├── config/                  # 配置文件 ✅
├── data/                    # SQLite 数据库文件
├── logs/                    # 日志文件
└── main.go                  # 入口文件 ✅

gold-admin-frontend/          # 前端（Vue 2）
├── src/
│   ├── api/                # API 调用
│   ├── layout/             # 布局组件
│   ├── router/             # 路由配置
│   ├── store/              # 状态管理
│   ├── utils/              # 工具函数
│   └── views/              # 页面组件
│       ├── dashboard/      # 首页看板
│       ├── login/          # 登录页
│       ├── price/          # 价格管理
│       └── system/         # 系统管理
│           ├── user/       # 用户管理
│           ├── role/       # 角色管理
│           └── menu/       # 菜单管理
└── vue.config.js           # Vue 配置
```

---

## 🔍 快速测试

### 1. 测试后端服务

```bash
# 健康检查
curl http://localhost:8080/ping
# 返回：{"message":"pong"}

# 登录测试
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 2. 测试前端界面

访问：http://localhost:9527

**测试流程：**
1. 登录系统（admin/admin123）
2. 查看首页统计数据
3. 进入系统管理 → 用户管理
4. 创建一个新用户
5. 进入业务管理 → 价格管理
6. 修改价格数据
7. 进入业务管理 → 预约管理
8. 查看预约列表

---

## ⚙️ 配置说明

### 后端配置文件
位置：`config/config.yaml`

```yaml
server:
  port: 8080              # 后端端口
  mode: debug             # debug / release

database:
  type: sqlite            # 数据库类型
  path: ./data/gold_admin.db  # 数据库路径

jwt:
  secret: gold-admin-secret-key-change-me-in-production
  expire: 168             # Token 过期时间（小时）
```

### 前端配置文件
位置：`vue.config.js`

```javascript
devServer: {
  port: 9527,             // 前端端口
  proxy: {
    '/api': {
      target: 'http://localhost:8080',  // 后端地址
      changeOrigin: true
    }
  }
}
```

---

## 🐛 常见问题

### Q1: 后端启动失败？
**A**: 检查 8080 端口是否被占用
```bash
lsof -i :8080
# 如果被占用，杀掉进程或修改配置文件中的端口
```

### Q2: 前端启动失败？
**A**: 检查 9527 端口是否被占用，或删除 node_modules 重新安装
```bash
rm -rf node_modules package-lock.json
npm install
```

### Q3: 登录后没有菜单？
**A**: 检查用户是否分配了角色，角色是否有菜单权限

### Q4: 如何重置数据库？
**A**: 删除数据库文件后重启后端
```bash
rm -rf ./data/gold_admin.db
go run main.go
```

### Q5: 如何修改端口？
**A**: 
- 后端：修改 `config/config.yaml` 中的 `server.port`
- 前端：修改 `vue.config.js` 中的 `devServer.port`

---

## 📱 小程序对接

本系统提供的 API 可以直接给微信小程序调用：

### 小程序端需要调用的接口
- `GET /api/prices` - 获取价格列表
- `POST /api/appointments` - 创建预约
- `GET /api/shops/all` - 获取店铺列表

详见小程序文档：`../miniprogram/如何对接真实数据.md`

---

## 🚀 生产环境部署

### 方案一：前后端分离部署

**后端：**
```bash
go build -ldflags="-w -s" -o gold-admin main.go
./gold-admin
```

**前端：**
```bash
npm run build
# 将 dist 目录部署到 Nginx
```

### 方案二：嵌入式部署（推荐）

```bash
# 1. 打包前端
cd gold-admin-frontend && npm run build

# 2. 编译 Go（自动嵌入前端静态文件）
cd ../gold-admin-backend
go build -ldflags="-w -s" -o gold-admin main.go

# 3. 运行（前后端合一）
./gold-admin

# 访问：http://localhost:8080
```

详见：[DEPLOY.md](./DEPLOY.md)

---

## 📊 数据库查看

### 使用 SQLite 命令行
```bash
sqlite3 ./data/gold_admin.db

.tables                    # 查看所有表
SELECT * FROM admin_users; # 查看用户
SELECT * FROM roles;       # 查看角色
SELECT * FROM menus;       # 查看菜单
.quit                      # 退出
```

### 使用图形界面工具
推荐：[DB Browser for SQLite](https://sqlitebrowser.org/)

---

## 📚 相关文档

- [完整技术文档](./README.md)
- [API 接口文档](./API.md)
- [部署指南](./DEPLOY.md)
- [快速启动](./快速启动.md)

---

## 🎉 恭喜！

**沪金汇后台管理系统**已经完全开发完成并可以使用！

### 已完成的功能 ✅
- ✅ 用户、角色、菜单管理（RBAC 权限体系）
- ✅ 价格管理（基础价±差价模式）
- ✅ 店铺管理
- ✅ 预约管理
- ✅ 首页看板统计
- ✅ JWT 认证
- ✅ 动态路由
- ✅ 权限控制

### 系统架构 🏗️
- **后端**: Go + Gin + GORM + SQLite
- **前端**: Vue 2 + Element UI + Vuex + Vue Router
- **数据库**: SQLite（可切换 MySQL）
- **认证**: JWT
- **权限**: RBAC

### 快速访问 🚀
- **前端管理界面**: http://localhost:9527
- **后端 API**: http://localhost:8080
- **登录账号**: admin / admin123

---

**祝您使用愉快！如有问题，请参考相关文档。💰**
