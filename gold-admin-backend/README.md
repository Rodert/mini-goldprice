# 沪金汇后台管理系统 - 技术设计文档

## 📋 项目概述

**项目名称**: 沪金汇后台管理系统  
**项目类型**: 贵金属回收店铺管理后台  
**技术架构**: 前后端分离 + 嵌入部署  
**开发日期**: 2025-11-04  
**版本**: v1.0.0

### 核心功能
- 多店铺管理
- 贵金属价格管理（支持基础价±差价模式）
- 预约管理
- 回收记录管理
- 权限管理（RBAC，页面级权限，预留按钮级）

---

## 🎯 技术栈选型

### 后端技术栈

```yaml
语言: Golang 1.20+
框架: Gin (Web框架)
ORM: GORM (数据库操作)
数据库: SQLite 3 (开发/演示) → MySQL 8.0+ (生产可选)
缓存: Redis (可选)
认证: JWT (github.com/golang-jwt/jwt/v4)
配置: Viper
日志: zap (可选)
```

**核心依赖包：**
```go
github.com/gin-gonic/gin              // Web框架
gorm.io/gorm                          // ORM
gorm.io/driver/sqlite                 // SQLite驱动 ⭐
github.com/golang-jwt/jwt/v4          // JWT认证
github.com/spf13/viper                // 配置管理
golang.org/x/crypto/bcrypt            // 密码加密
github.com/gin-contrib/cors           // CORS跨域
```

**数据库选型说明：**
- **开发阶段**: 使用 SQLite（无需安装数据库服务，单文件存储）
- **生产环境**: 可选择迁移到 MySQL（性能更好，支持并发）
- **迁移方式**: GORM 自动迁移或导出SQL，几乎无缝切换

### 前端技术栈

```yaml
框架: Vue 2.6.x
UI组件库: Element UI 2.15.x
路由: Vue Router 3.x
状态管理: Vuex 3.x
HTTP库: Axios
构建工具: Vue CLI / Webpack
模板: vue-admin-template (精简版)
```

**核心依赖包：**
```json
{
  "vue": "^2.6.14",
  "vue-router": "^3.5.1",
  "vuex": "^3.6.2",
  "element-ui": "^2.15.13",
  "axios": "^1.6.0",
  "js-cookie": "^3.0.1"
}
```

### 为什么选择这个技术栈？

| 技术 | 理由 |
|------|------|
| **Golang** | 性能强、部署简单（单一二进制）、并发能力强 |
| **Gin** | 轻量级、性能好、中文文档全 |
| **GORM** | 最流行的Go ORM、支持SQLite/MySQL、API友好 |
| **SQLite** | 零配置、单文件、开发便捷、可迁移MySQL |
| **Vue 2** | 成熟稳定、生态完善、学习曲线平缓 |
| **Element UI** | 组件丰富、企业级UI、开箱即用 |
| **JWT** | 无状态、适合前后端分离、扩展性好 |

---

## 🔐 权限设计方案（RBAC）

### 权限模型

```
用户（User） ←→ 角色（Role） ←→ 菜单/权限（Menu）

特点：
- 页面级权限（当前实现）✅
- 按钮级权限（预留字段）⭐
- 多店铺通过菜单区分
```

### 权限粒度

**当前阶段：页面级权限**
- 控制用户能访问哪些页面
- 通过菜单的 type=1（目录）和 type=2（菜单）实现
- 足够满足90%的场景

**扩展阶段：按钮级权限（预留）**
- 控制页面内的按钮显示/隐藏
- 通过菜单的 type=3（按钮）+ permission 字段实现
- 在 menus 表中已预留字段

### 角色设计（预设5个角色）

| 角色ID | 角色名称 | 角色编码 | 说明 |
|--------|----------|----------|------|
| 1 | 超级管理员 | super_admin | 所有权限 |
| 2 | 总部店长 | head_manager | 管理所有店铺 |
| 3 | 单店店长 | shop_manager | 管理单个店铺 |
| 4 | 店员 | shop_staff | 处理日常业务（只读） |
| 5 | 财务 | finance | 查看数据、导出报表 |

### 菜单结构（多店铺支持）

```
首页
业务管理/
  ├─ 店铺管理
  ├─ 沪金汇1店/
  │   ├─ 价格管理
  │   └─ 预约管理
  ├─ 沪金汇2店/
  │   ├─ 价格管理
  │   └─ 预约管理
  └─ 沪金汇3店/
      ├─ 价格管理
      └─ 预约管理
系统管理/
  ├─ 用户管理
  ├─ 角色管理
  └─ 菜单管理
```

**说明：**
- 每个店铺是独立的菜单目录
- 店员只能看到自己店铺的菜单
- 通过角色分配菜单权限实现多店铺隔离

---

## 💾 数据库设计（15张表）

### 核心表关系图

```
admin_users ←→ user_roles ←→ roles ←→ role_menus ←→ menus
                                 │
                                 └─→ shops → prices
                                           → appointments
                                           → recycling_records
```

### 表结构设计

#### 1. 权限相关（5张表）

**1.1 admin_users - 管理员用户表**
```sql
-- SQLite 语法
CREATE TABLE admin_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(50) UNIQUE NOT NULL,          -- 用户名
  password VARCHAR(255) NOT NULL,                -- 密码（bcrypt加密）
  real_name VARCHAR(50),                         -- 真实姓名
  phone VARCHAR(20),                             -- 手机号
  email VARCHAR(100),                            -- 邮箱
  avatar VARCHAR(255),                           -- 头像URL
  status TINYINT DEFAULT 1,                      -- 状态 1:启用 0:禁用
  last_login_time DATETIME,                      -- 最后登录时间
  last_login_ip VARCHAR(50),                     -- 最后登录IP
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);

CREATE INDEX idx_admin_users_username ON admin_users(username);
CREATE INDEX idx_admin_users_status ON admin_users(status);
```

**GORM 模型定义：**
```go
type AdminUser struct {
    ID            uint      `gorm:"primaryKey" json:"id"`
    Username      string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Password      string    `gorm:"size:255;not null" json:"-"`
    RealName      string    `gorm:"size:50" json:"real_name"`
    Phone         string    `gorm:"size:20" json:"phone"`
    Email         string    `gorm:"size:100" json:"email"`
    Avatar        string    `gorm:"size:255" json:"avatar"`
    Status        int8      `gorm:"default:1" json:"status"`
    LastLoginTime time.Time `json:"last_login_time"`
    LastLoginIP   string    `gorm:"size:50" json:"last_login_ip"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}
```

**1.2 roles - 角色表**
```sql
CREATE TABLE roles (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL COMMENT '角色名称',
  code VARCHAR(50) UNIQUE NOT NULL COMMENT '角色标识',
  description TEXT COMMENT '角色描述',
  sort INT DEFAULT 0 COMMENT '排序',
  status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_code (code),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';
```

**1.3 menus - 菜单表**
```sql
CREATE TABLE menus (
  id INT PRIMARY KEY AUTO_INCREMENT,
  parent_id INT DEFAULT 0 COMMENT '父菜单ID，0为顶级',
  type TINYINT DEFAULT 1 COMMENT '类型 1:目录 2:菜单 3:按钮',
  name VARCHAR(50) NOT NULL COMMENT '菜单名称（英文）',
  title VARCHAR(50) NOT NULL COMMENT '菜单标题（中文）',
  icon VARCHAR(50) COMMENT '图标',
  path VARCHAR(100) COMMENT '路由路径',
  component VARCHAR(100) COMMENT '组件路径',
  permission VARCHAR(100) COMMENT '权限标识（预留）',
  sort INT DEFAULT 0 COMMENT '排序',
  visible TINYINT DEFAULT 1 COMMENT '是否显示 1:显示 0:隐藏',
  status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_parent_id (parent_id),
  INDEX idx_type (type),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';
```

**1.4 user_roles - 用户角色关联表**
```sql
CREATE TABLE user_roles (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL COMMENT '用户ID',
  role_id INT NOT NULL COMMENT '角色ID',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_role (user_id, role_id),
  INDEX idx_user_id (user_id),
  INDEX idx_role_id (role_id),
  FOREIGN KEY (user_id) REFERENCES admin_users(id) ON DELETE CASCADE,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';
```

**1.5 role_menus - 角色菜单关联表**
```sql
CREATE TABLE role_menus (
  id INT PRIMARY KEY AUTO_INCREMENT,
  role_id INT NOT NULL COMMENT '角色ID',
  menu_id INT NOT NULL COMMENT '菜单ID',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_role_menu (role_id, menu_id),
  INDEX idx_role_id (role_id),
  INDEX idx_menu_id (menu_id),
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色菜单关联表';
```

#### 2. 业务相关（7张表）

**2.1 shops - 店铺表**
```sql
CREATE TABLE shops (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL COMMENT '店铺名称',
  code VARCHAR(50) UNIQUE NOT NULL COMMENT '店铺编码（shop1, shop2）',
  address VARCHAR(255) COMMENT '地址',
  phone VARCHAR(20) COMMENT '固定电话',
  mobile VARCHAR(20) COMMENT '手机号',
  hours VARCHAR(100) COMMENT '营业时间',
  latitude DECIMAL(10,7) COMMENT '纬度（用于小程序导航）',
  longitude DECIMAL(10,7) COMMENT '经度',
  description TEXT COMMENT '店铺介绍',
  status TINYINT DEFAULT 1 COMMENT '状态 1:营业 0:停业',
  sort INT DEFAULT 0 COMMENT '排序',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_code (code),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='店铺表';
```

**2.2 prices - 价格表（核心）⭐**
```sql
CREATE TABLE prices (
  id INT PRIMARY KEY AUTO_INCREMENT,
  shop_id INT COMMENT '店铺ID（可为空表示全局价格）',
  code VARCHAR(50) NOT NULL COMMENT '唯一标识（gold_9999, silver_999）',
  name VARCHAR(50) NOT NULL COMMENT '品种名称',
  subtitle VARCHAR(100) COMMENT '副标题',
  icon VARCHAR(10) COMMENT '图标（Au, Ag）',
  icon_color VARCHAR(20) COMMENT '图标颜色',
  
  base_price DECIMAL(10,2) NOT NULL COMMENT '基础价格（元/克）',
  buy_price_diff DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '回购差价（可为负）',
  sell_price_diff DECIMAL(10,2) NOT NULL DEFAULT 0 COMMENT '销售差价（可为正）',
  
  -- 自动计算字段（或应用层计算）
  buy_price DECIMAL(10,2) GENERATED ALWAYS AS (base_price + buy_price_diff) STORED COMMENT '回购价',
  sell_price DECIMAL(10,2) GENERATED ALWAYS AS (base_price + sell_price_diff) STORED COMMENT '销售价',
  
  sort INT DEFAULT 0 COMMENT '排序',
  status TINYINT DEFAULT 1 COMMENT '状态 1:启用 0:禁用',
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  
  UNIQUE KEY uk_shop_code (shop_id, code),
  INDEX idx_code (code),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='价格表';
```

**价格计算说明：**
```
回购价 = base_price + buy_price_diff
销售价 = base_price + sell_price_diff

示例：
- 基础价: 560.00
- 回购差价: -10.00 → 回购价 = 550.00
- 销售差价: +15.00 → 销售价 = 575.00
```

**2.3 price_histories - 价格历史表**
```sql
CREATE TABLE price_histories (
  id INT PRIMARY KEY AUTO_INCREMENT,
  price_id INT NOT NULL COMMENT '价格记录ID',
  base_price DECIMAL(10,2) COMMENT '基础价格',
  buy_price_diff DECIMAL(10,2) COMMENT '回购差价',
  sell_price_diff DECIMAL(10,2) COMMENT '销售差价',
  buy_price DECIMAL(10,2) COMMENT '回购价',
  sell_price DECIMAL(10,2) COMMENT '销售价',
  change_reason VARCHAR(255) COMMENT '变动原因',
  operator_id INT COMMENT '操作人ID',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_price_id (price_id),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='价格历史表';
```

**2.4 appointments - 预约表**
```sql
CREATE TABLE appointments (
  id INT PRIMARY KEY AUTO_INCREMENT,
  shop_id INT COMMENT '店铺ID',
  user_id INT COMMENT '小程序用户ID',
  openid VARCHAR(100) COMMENT '微信openid',
  metal_type VARCHAR(50) COMMENT '品种',
  service_type VARCHAR(20) COMMENT '服务类型（store:到店 home:上门）',
  appointment_time DATETIME COMMENT '预约时间',
  name VARCHAR(50) COMMENT '姓名',
  phone VARCHAR(20) COMMENT '电话',
  address VARCHAR(255) COMMENT '地址（上门回收）',
  note TEXT COMMENT '客户备注',
  admin_remark TEXT COMMENT '管理员备注',
  status VARCHAR(20) DEFAULT 'pending' COMMENT '状态（pending/confirmed/completed/cancelled）',
  confirmed_at DATETIME COMMENT '确认时间',
  completed_at DATETIME COMMENT '完成时间',
  cancelled_at DATETIME COMMENT '取消时间',
  handler_id INT COMMENT '处理人ID',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_shop_id (shop_id),
  INDEX idx_status (status),
  INDEX idx_appointment_time (appointment_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预约表';
```

**2.5 recycling_records - 回收记录表**
```sql
CREATE TABLE recycling_records (
  id INT PRIMARY KEY AUTO_INCREMENT,
  record_no VARCHAR(50) UNIQUE NOT NULL COMMENT '回收单号',
  shop_id INT COMMENT '店铺ID',
  user_id INT COMMENT '小程序用户ID',
  appointment_id INT COMMENT '关联预约ID',
  customer_name VARCHAR(50) COMMENT '客户姓名',
  customer_phone VARCHAR(20) COMMENT '客户电话',
  metal_type VARCHAR(50) COMMENT '品种',
  weight DECIMAL(10,3) COMMENT '重量（克）',
  purity DECIMAL(5,2) COMMENT '成色（%）',
  unit_price DECIMAL(10,2) COMMENT '单价（元/克）',
  total_amount DECIMAL(12,2) COMMENT '总金额（元）',
  payment_method VARCHAR(20) COMMENT '支付方式（cash:现金 transfer:转账）',
  note TEXT COMMENT '备注',
  operator_id INT COMMENT '操作员ID',
  operator_name VARCHAR(50) COMMENT '操作员姓名',
  status TINYINT DEFAULT 1 COMMENT '状态 1:正常 0:已作废',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_record_no (record_no),
  INDEX idx_shop_id (shop_id),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='回收记录表';
```

**2.6 miniapp_users - 小程序用户表**
```sql
CREATE TABLE miniapp_users (
  id INT PRIMARY KEY AUTO_INCREMENT,
  openid VARCHAR(100) UNIQUE NOT NULL COMMENT '微信openid',
  unionid VARCHAR(100) COMMENT '微信unionid',
  nickname VARCHAR(100) COMMENT '昵称',
  avatar VARCHAR(255) COMMENT '头像',
  phone VARCHAR(20) COMMENT '手机号',
  gender TINYINT COMMENT '性别 0:未知 1:男 2:女',
  tags VARCHAR(255) COMMENT '标签',
  last_visit_time DATETIME COMMENT '最后访问时间',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_openid (openid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='小程序用户表';
```

**2.7 calculations - 计算历史表**
```sql
CREATE TABLE calculations (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT COMMENT '用户ID',
  openid VARCHAR(100) COMMENT '微信openid',
  metal_type VARCHAR(50) COMMENT '品种',
  weight DECIMAL(10,3) COMMENT '重量（克）',
  purity DECIMAL(5,2) COMMENT '成色（%）',
  unit_price DECIMAL(10,2) COMMENT '单价（元/克）',
  result DECIMAL(12,2) COMMENT '计算结果（元）',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='计算历史表';
```

#### 3. 日志相关（3张表）

**3.1 login_logs - 登录日志表**
```sql
CREATE TABLE login_logs (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT COMMENT '用户ID',
  username VARCHAR(50) COMMENT '用户名',
  ip VARCHAR(50) COMMENT 'IP地址',
  location VARCHAR(100) COMMENT 'IP归属地',
  device VARCHAR(100) COMMENT '设备信息',
  status TINYINT COMMENT '状态 1:成功 0:失败',
  message VARCHAR(255) COMMENT '失败原因',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';
```

**3.2 operation_logs - 操作日志表**
```sql
CREATE TABLE operation_logs (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT COMMENT '操作人ID',
  username VARCHAR(50) COMMENT '操作人用户名',
  module VARCHAR(50) COMMENT '模块',
  action VARCHAR(50) COMMENT '操作（create/update/delete）',
  content TEXT COMMENT '操作内容',
  ip VARCHAR(50) COMMENT 'IP地址',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_user_id (user_id),
  INDEX idx_module (module),
  INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';
```

**3.3 appointment_status_logs - 预约状态流转表**
```sql
CREATE TABLE appointment_status_logs (
  id INT PRIMARY KEY AUTO_INCREMENT,
  appointment_id INT NOT NULL COMMENT '预约ID',
  from_status VARCHAR(20) COMMENT '原状态',
  to_status VARCHAR(20) COMMENT '新状态',
  remark TEXT COMMENT '备注',
  operator_id INT COMMENT '操作人ID',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_appointment_id (appointment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预约状态流转表';
```

---

## 📋 功能模块清单

### MVP版本（第一阶段）

#### 1. 登录认证模块 🔐
- [x] 管理员登录
- [x] JWT Token 认证
- [x] 获取用户信息
- [x] 获取用户菜单权限
- [x] 登出

#### 2. 首页看板 📊
- [x] 核心数据统计（今日预约、回收、用户数）
- [x] 快捷入口
- [x] 最新动态列表

#### 3. 系统管理 ⚙️
- [x] 用户管理（CRUD、分配角色、启禁用、重置密码）
- [x] 角色管理（CRUD、分配菜单权限）
- [x] 菜单管理（树形结构、CRUD）

#### 4. 业务管理 - 店铺管理 🏪
- [x] 店铺列表
- [x] 新增店铺
- [x] 编辑店铺信息
- [x] 启用/禁用店铺
- [ ] 删除店铺（需检查业务数据）

#### 5. 业务管理 - 价格管理 💰（核心）
- [x] 价格列表（卡片展示）
- [x] 新增贵金属品种（唯一标识）
- [x] 编辑价格（基础价 ± 差价模式）
- [x] 删除品种
- [x] 价格历史记录
- [ ] 同步国际金价API（可选）

#### 6. 业务管理 - 预约管理 📅
- [x] 预约列表（分页、筛选）
- [x] 预约详情
- [x] 更新预约状态
- [x] 添加管理员备注
- [ ] 转换为回收记录

### 增强版（第二阶段）

#### 7. 业务管理 - 回收记录 📝
- [ ] 回收记录列表
- [ ] 新增回收记录
- [ ] 查看详情
- [ ] 导出Excel
- [ ] 统计汇总

#### 8. 小程序用户管理 📱
- [ ] 小程序用户列表
- [ ] 用户详情（预约、回收、计算历史）
- [ ] 用户标签

#### 9. 操作日志 📋
- [ ] 登录日志列表
- [ ] 操作日志列表
- [ ] 日志详情
- [ ] 日志导出

### 完整版（第三阶段）

#### 10. 数据统计 📈
- [ ] 营业数据统计
- [ ] 趋势分析图表
- [ ] 排行榜（店铺、员工、品种）
- [ ] 数据导出

---

## 🚀 部署方案（方案A：嵌入部署）

### 最终产物

```
gold-admin              # 单一可执行文件（~20-25MB）
config.yaml             # 配置文件（可选）
```

### 部署架构

```
┌────────────────────────────────────┐
│  gold-admin (单一可执行文件)        │
│  ├─ Go 后端（Gin + GORM）          │
│  ├─ 前端静态文件（嵌入到二进制）    │
│  └─ API 路由 + 静态文件服务        │
└────────────────────────────────────┘
         │
         ├─ /api/*        → Go API处理
         ├─ /static/*     → 嵌入的静态文件
         └─ /*            → index.html (Vue SPA)
```

### 嵌入实现（Go 1.16+）

```go
package main

import (
    "embed"
    "io/fs"
    "net/http"
    "github.com/gin-gonic/gin"
)

//go:embed dist/*
var staticFiles embed.FS

func main() {
    r := gin.Default()
    
    // API 路由
    api := r.Group("/api")
    {
        api.POST("/login", Login)
        api.GET("/user/list", GetUserList)
        // ...
    }
    
    // 静态文件服务
    staticFS, _ := fs.Sub(staticFiles, "dist")
    r.StaticFS("/static", http.FS(staticFS))
    
    // SPA 路由（所有非API请求返回 index.html）
    r.NoRoute(func(c *gin.Context) {
        data, _ := staticFiles.ReadFile("dist/index.html")
        c.Data(200, "text/html; charset=utf-8", data)
    })
    
    r.Run(":8080")
}
```

### 构建命令

```bash
# 开发环境（前后端分离）
## 后端
go run main.go

## 前端
cd web && npm run dev

# 生产环境（嵌入打包）
## 1. 打包前端
cd web && npm run build

## 2. 编译Go（会自动嵌入 dist 目录）
go build -ldflags="-w -s" -o gold-admin main.go

## 3. 运行
./gold-admin
```

### 部署优势

- ✅ 单一文件部署，极简
- ✅ 无需 Nginx
- ✅ 跨平台（Linux/Mac/Windows）
- ✅ 性能好（Go 静态文件服务）
- ✅ 防篡改（静态文件在二进制中）

---

## 📂 项目目录结构

```
gold-admin-backend/
├── api/                        # API 处理器
│   └── v1/
│       ├── auth.go             # 登录认证
│       ├── user.go             # 用户管理
│       ├── role.go             # 角色管理
│       ├── menu.go             # 菜单管理
│       ├── shop.go             # 店铺管理
│       ├── price.go            # 价格管理
│       ├── appointment.go      # 预约管理
│       ├── record.go           # 回收记录
│       ├── dashboard.go        # 首页看板
│       └── log.go              # 日志管理
│
├── models/                     # 数据模型
│   ├── init.go                 # 数据库初始化
│   ├── admin_user.go
│   ├── role.go
│   ├── menu.go
│   ├── user_role.go
│   ├── role_menu.go
│   ├── shop.go
│   ├── price.go
│   ├── appointment.go
│   ├── record.go
│   └── log.go
│
├── router/                     # 路由
│   └── router.go
│
├── middleware/                 # 中间件
│   ├── cors.go                 # 跨域
│   ├── jwt.go                  # JWT认证
│   ├── permission.go           # 权限验证
│   └── logger.go               # 日志
│
├── utils/                      # 工具函数
│   ├── response.go             # 统一响应
│   ├── jwt.go                  # JWT工具
│   ├── password.go             # 密码加密
│   ├── validator.go            # 参数验证
│   └── paginate.go             # 分页
│
├── config/                     # 配置
│   ├── config.go
│   └── config.yaml
│
├── sql/                        # SQL 脚本
│   ├── schema.sql              # 数据库结构
│   └── init_data.sql           # 初始化数据
│
├── docs/                       # 文档
│   ├── ui-prototype/           # UI 原型（HTML）
│   ├── README.md
│   └── API.md
│
├── web/                        # Vue 前端项目
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── vue.config.js
│
├── dist/                       # 前端打包文件（自动生成）
│
├── main.go                     # 入口文件
├── go.mod
├── go.sum
├── Makefile                    # 构建脚本
├── .gitignore
└── README.md                   # 本文档
```

---

## 🔌 API 接口设计

### 统一响应格式

```go
type Response struct {
    Code    int         `json:"code"`    // 200:成功 其他:失败
    Data    interface{} `json:"data"`    // 返回数据
    Message string      `json:"message"` // 提示信息
}

// 成功
{
  "code": 200,
  "data": {...},
  "message": "success"
}

// 失败
{
  "code": 400,
  "data": null,
  "message": "参数错误"
}
```

### 核心接口列表

#### 认证相关
```
POST   /api/login                 登录
POST   /api/logout                登出
GET    /api/user/info             获取当前用户信息
GET    /api/user/menus            获取用户菜单权限
POST   /api/refresh-token         刷新Token
```

#### 用户管理
```
GET    /api/admin/user/list       用户列表
POST   /api/admin/user/add        新增用户
PUT    /api/admin/user/update/:id 更新用户
DELETE /api/admin/user/delete/:id 删除用户
PUT    /api/admin/user/status/:id 更新状态
POST   /api/admin/user/assign-role 分配角色
PUT    /api/admin/user/reset-pwd/:id 重置密码
```

#### 角色管理
```
GET    /api/admin/role/list       角色列表
POST   /api/admin/role/add        新增角色
PUT    /api/admin/role/update/:id 更新角色
DELETE /api/admin/role/delete/:id 删除角色
POST   /api/admin/role/assign-menu 分配菜单
GET    /api/admin/role/menus/:id  获取角色菜单
```

#### 菜单管理
```
GET    /api/admin/menu/tree       菜单树
GET    /api/admin/menu/list       菜单列表
POST   /api/admin/menu/add        新增菜单
PUT    /api/admin/menu/update/:id 更新菜单
DELETE /api/admin/menu/delete/:id 删除菜单
```

#### 店铺管理
```
GET    /api/shop/list             店铺列表
POST   /api/shop/add              新增店铺
PUT    /api/shop/update/:id       更新店铺
DELETE /api/shop/delete/:id       删除店铺
PUT    /api/shop/status/:id       更新状态
```

#### 价格管理
```
GET    /api/price/list            价格列表
GET    /api/price/detail/:id      价格详情
POST   /api/price/add             新增品种
PUT    /api/price/update/:id      更新价格
DELETE /api/price/delete/:id      删除品种
GET    /api/price/history/:id     价格历史
POST   /api/price/refresh         刷新基础价（同步API）
```

#### 预约管理
```
GET    /api/appointment/list      预约列表
GET    /api/appointment/detail/:id 预约详情
PUT    /api/appointment/status/:id 更新状态
PUT    /api/appointment/remark/:id 添加备注
POST   /api/appointment/to-record 转为回收记录
DELETE /api/appointment/delete/:id 删除预约
```

#### 回收记录
```
GET    /api/record/list           回收记录列表
POST   /api/record/add            新增记录
GET    /api/record/detail/:id     记录详情
PUT    /api/record/update/:id     更新记录
DELETE /api/record/delete/:id     删除记录
GET    /api/record/export         导出Excel
GET    /api/record/stats          统计数据
```

#### 首页看板
```
GET    /api/dashboard/stats       统计数据
GET    /api/dashboard/trends      趋势数据
GET    /api/dashboard/recent      最新动态
```

---

## 📝 开发优先级

### 第一阶段（核心功能）- 3-5天
```
✅ 1. 项目初始化
   - Go 项目结构
   - 数据库连接
   - 基础配置

✅ 2. 认证模块
   - JWT 中间件
   - 登录/登出
   - 密码加密

✅ 3. 用户/角色/菜单管理（RBAC核心）
   - CRUD 接口
   - 权限分配
   - 菜单树生成

✅ 4. 前端框架搭建
   - vue-admin-template 集成
   - 路由配置
   - 权限路由守卫
```

### 第二阶段（业务功能）- 5-7天
```
✅ 5. 店铺管理
✅ 6. 价格管理（基础价±差价模式）
✅ 7. 预约管理
✅ 8. 首页看板
```

### 第三阶段（完善优化）- 3-5天
```
✅ 9. 回收记录
✅ 10. 操作日志
✅ 11. 数据统计
✅ 12. 导出功能
```

### 第四阶段（测试部署）- 2-3天
```
✅ 13. 功能测试
✅ 14. 前端打包
✅ 15. Go 编译（嵌入静态文件）
✅ 16. 部署上线
```

---

## 🎯 关键技术点

### 1. 价格管理（核心业务逻辑）

**设计思路：**
```
基础价格（国际金价/市场价） ± 差价 = 最终价格

优势：
- 灵活：商家可快速调整差价策略
- 透明：基础价和差价分开管理
- 扩展：可为不同客户设置不同差价
```

**实现示例：**
```go
type Price struct {
    ID            uint    `gorm:"primaryKey"`
    Code          string  `gorm:"uniqueIndex;not null"` // gold_9999
    Name          string  `gorm:"not null"`
    BasePrice     float64 `gorm:"not null"`             // 基础价
    BuyPriceDiff  float64 `gorm:"not null;default:0"`  // 回购差价（可为负）
    SellPriceDiff float64 `gorm:"not null;default:0"`  // 销售差价（可为正）
}

// 计算方法
func (p *Price) GetBuyPrice() float64 {
    return p.BasePrice + p.BuyPriceDiff
}

func (p *Price) GetSellPrice() float64 {
    return p.BasePrice + p.SellPriceDiff
}
```

### 2. 权限验证流程

**登录流程：**
```
1. 用户输入账号密码
2. bcrypt 验证密码
3. 生成 JWT Token
4. 查询用户角色
5. 查询角色菜单（type=1,2）
6. 返回 Token + 菜单树
7. 前端存储 Token 和菜单
8. 动态生成路由
```

**请求验证：**
```
1. 前端请求携带 Token (Authorization: Bearer xxx)
2. JWT 中间件验证 Token
3. 解析出 user_id
4. 权限中间件验证（可选）
5. 业务逻辑处理
```

### 3. 菜单树生成算法

```go
// 递归生成菜单树
func BuildMenuTree(menus []Menu, parentID int) []MenuTree {
    var tree []MenuTree
    for _, menu := range menus {
        if menu.ParentID == parentID {
            node := MenuTree{
                ID:       menu.ID,
                Title:    menu.Title,
                Children: BuildMenuTree(menus, menu.ID),
            }
            tree = append(tree, node)
        }
    }
    return tree
}
```

---

## 📌 注意事项

### 安全相关
1. **密码加密**：使用 bcrypt，cost >= 10
2. **JWT 密钥**：使用随机字符串，不要提交到代码库
3. **SQL 注入**：使用 GORM 参数化查询
4. **XSS 防护**：前端输入过滤
5. **CORS 配置**：生产环境限制来源

### 性能优化
1. **数据库索引**：常用查询字段建索引
2. **分页查询**：避免一次性加载大量数据
3. **Redis 缓存**：菜单、配置等热数据
4. **静态资源**：前端打包压缩、CDN

### 数据一致性
1. **事务处理**：涉及多表操作使用事务
2. **外键约束**：ON DELETE CASCADE 谨慎使用
3. **软删除**：重要数据不物理删除

---

## 🔧 配置文件示例

**config.yaml**
```yaml
server:
  port: 8080
  mode: release # debug / release

database:
  type: sqlite                           # sqlite / mysql
  path: ./data/gold_admin.db             # SQLite 数据库文件路径 ⭐
  # MySQL 配置（如需切换）
  # host: 127.0.0.1
  # port: 3306
  # username: root
  # password: 123456
  # database: gold_admin
  # charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

jwt:
  secret: your-secret-key-change-me
  expire: 168 # 小时（7天）

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0
  enabled: false                         # 是否启用 Redis

log:
  level: info
  path: ./logs
```

**数据库初始化代码：**
```go
package models

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
    var err error
    DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        return err
    }

    // 自动迁移所有表
    err = DB.AutoMigrate(
        &AdminUser{},
        &Role{},
        &Menu{},
        &UserRole{},
        &RoleMenu{},
        &Shop{},
        &Price{},
        &PriceHistory{},
        &Appointment{},
        &RecyclingRecord{},
        &MiniappUser{},
        &Calculation{},
        &LoginLog{},
        &OperationLog{},
        &AppointmentStatusLog{},
    )
    
    return err
}
```

---

## 📖 相关文档

- [UI 原型设计](./docs/ui-prototype/README.md)
- [API 接口文档](./docs/API.md)
- [数据库设计](./sql/schema.sql)
- [部署文档](./docs/DEPLOY.md)

---

## 🎉 开始开发

**克隆项目后：**
```bash
# 1. 安装依赖
go mod download

# 2. 修改配置（可选）
vim config/config.yaml

# 3. 运行后端（首次运行会自动创建 SQLite 数据库）
go run main.go

# 4. 运行前端（另一个终端）
cd web
npm install
npm run dev
```

**SQLite 优势：**
- ✅ 零配置：无需安装数据库服务
- ✅ 单文件：`gold_admin.db` 一个文件搞定
- ✅ 便携性：复制文件即可备份/迁移
- ✅ 开发快：适合快速开发和演示
- ✅ 可迁移：后期可无缝切换到 MySQL

---

**设计文档版本**: v1.0.0  
**创建日期**: 2025-11-04  
**作者**: AI Assistant  
**适用项目**: 沪金汇后台管理系统

