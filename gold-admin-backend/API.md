# API 接口文档

## 基础信息

- 基础URL: `http://localhost:8080/api`
- 数据格式: JSON
- 字符编码: UTF-8

## 统一响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

- `code`: 状态码，200表示成功，其他表示失败
- `message`: 提示信息
- `data`: 返回数据

## 认证说明

除了登录接口外，其他接口都需要在请求头中携带 Token：

```
Authorization: Bearer <token>
```

---

## 1. 认证模块

### 1.1 登录

**接口**: `POST /api/login`

**请求参数**:
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_info": {
      "id": 1,
      "username": "admin",
      "real_name": "系统管理员",
      "avatar": "",
      "roles": ["admin"]
    },
    "menu_list": []
  }
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 1.2 登出

**接口**: `POST /api/logout`

**响应示例**:
```json
{
  "code": 200,
  "message": "登出成功",
  "data": null
}
```

### 1.3 获取当前用户信息

**接口**: `GET /api/user/info`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "real_name": "系统管理员",
    "avatar": "",
    "roles": ["admin"]
  }
}
```

---

## 2. 用户管理

### 2.1 获取用户列表

**接口**: `GET /api/users`

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）
- `username`: 用户名（模糊查询）
- `real_name`: 真实姓名（模糊查询）
- `status`: 状态（1:启用 0:禁用）

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "admin",
        "real_name": "系统管理员",
        "phone": "",
        "email": "",
        "avatar": "",
        "status": 1,
        "last_login_time": "2025-11-04T10:00:00Z",
        "last_login_ip": "127.0.0.1",
        "created_at": "2025-11-04T10:00:00Z",
        "updated_at": "2025-11-04T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

**curl 示例**:
```bash
curl http://localhost:8080/api/users?page=1&page_size=10 \
  -H "Authorization: Bearer <token>"
```

### 2.2 获取用户详情

**接口**: `GET /api/users/:id`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "real_name": "系统管理员",
      "status": 1
    },
    "role_ids": [1]
  }
}
```

### 2.3 创建用户

**接口**: `POST /api/users`

**请求参数**:
```json
{
  "username": "test",
  "password": "123456",
  "real_name": "测试用户",
  "phone": "13800138000",
  "email": "test@example.com",
  "avatar": "",
  "status": 1,
  "role_ids": [1]
}
```

**curl 示例**:
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "username": "test",
    "password": "123456",
    "real_name": "测试用户",
    "status": 1,
    "role_ids": [1]
  }'
```

### 2.4 更新用户

**接口**: `PUT /api/users/:id`

**请求参数**:
```json
{
  "real_name": "测试用户2",
  "phone": "13800138001",
  "email": "test2@example.com",
  "status": 1,
  "role_ids": [1, 2]
}
```

### 2.5 删除用户

**接口**: `DELETE /api/users/:id`

**curl 示例**:
```bash
curl -X DELETE http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer <token>"
```

### 2.6 修改用户密码

**接口**: `PUT /api/users/:id/password`

**请求参数**:
```json
{
  "password": "new_password"
}
```

---

## 3. 角色管理

### 3.1 获取角色列表

**接口**: `GET /api/roles`

**查询参数**:
- `page`: 页码
- `page_size`: 每页数量
- `name`: 角色名称
- `code`: 角色代码
- `status`: 状态

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "超级管理员",
        "code": "super_admin",
        "description": "拥有所有权限",
        "sort": 1,
        "status": 1,
        "created_at": "2025-11-04T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  }
}
```

### 3.2 获取所有角色（不分页）

**接口**: `GET /api/roles/all`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "超级管理员",
      "code": "super_admin",
      "status": 1
    }
  ]
}
```

### 3.3 获取角色详情

**接口**: `GET /api/roles/:id`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "role": {
      "id": 1,
      "name": "超级管理员",
      "code": "super_admin",
      "description": "拥有所有权限"
    },
    "menu_ids": [1, 2, 3]
  }
}
```

### 3.4 创建角色

**接口**: `POST /api/roles`

**请求参数**:
```json
{
  "name": "店长",
  "code": "shop_manager",
  "description": "店铺管理员",
  "sort": 2,
  "status": 1,
  "menu_ids": [1, 2, 5]
}
```

### 3.5 更新角色

**接口**: `PUT /api/roles/:id`

### 3.6 删除角色

**接口**: `DELETE /api/roles/:id`

---

## 4. 菜单管理

### 4.1 获取菜单列表（树形）

**接口**: `GET /api/menus`

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "parent_id": 0,
      "type": 1,
      "name": "system",
      "title": "系统管理",
      "icon": "el-icon-setting",
      "path": "/system",
      "component": "Layout",
      "sort": 1,
      "children": [
        {
          "id": 2,
          "parent_id": 1,
          "type": 2,
          "name": "user",
          "title": "用户管理",
          "icon": "el-icon-user",
          "path": "user",
          "component": "system/user/index",
          "sort": 1
        }
      ]
    }
  ]
}
```

### 4.2 获取菜单树（用于选择）

**接口**: `GET /api/menus/tree`

### 4.3 获取菜单详情

**接口**: `GET /api/menus/:id`

### 4.4 创建菜单

**接口**: `POST /api/menus`

**请求参数**:
```json
{
  "parent_id": 0,
  "type": 1,
  "name": "business",
  "title": "业务管理",
  "icon": "el-icon-s-management",
  "path": "/business",
  "component": "Layout",
  "sort": 10,
  "visible": 1,
  "status": 1
}
```

**类型说明**:
- `type`: 1-目录 2-菜单 3-按钮

### 4.5 更新菜单

**接口**: `PUT /api/menus/:id`

### 4.6 删除菜单

**接口**: `DELETE /api/menus/:id`

---

## 5. 价格管理

### 5.1 获取价格列表

**接口**: `GET /api/prices`

**查询参数**:
- `shop_id`: 店铺ID
- `code`: 品种代码
- `name`: 品种名称
- `status`: 状态

**响应示例**:
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "shop_id": null,
      "code": "gold_9999",
      "name": "黄金9999",
      "subtitle": "Au9999 · 千足金",
      "icon": "Au",
      "icon_color": "#FFD700",
      "base_price": 560.00,
      "buy_price_diff": -10.00,
      "sell_price_diff": 15.00,
      "buy_price": 550.00,
      "sell_price": 575.00,
      "sort": 1,
      "status": 1,
      "updated_at": "2025-11-04T14:30:00Z"
    }
  ]
}
```

### 5.2 获取价格详情

**接口**: `GET /api/prices/:id`

### 5.3 创建价格

**接口**: `POST /api/prices`

**请求参数**:
```json
{
  "shop_id": null,
  "code": "gold_9999",
  "name": "黄金9999",
  "subtitle": "Au9999 · 千足金",
  "icon": "Au",
  "icon_color": "#FFD700",
  "base_price": 560.00,
  "buy_price_diff": -10.00,
  "sell_price_diff": 15.00,
  "sort": 1,
  "status": 1
}
```

**字段说明**:
- `base_price`: 基础价格（从市场获取）
- `buy_price_diff`: 回购差价（通常为负数）
- `sell_price_diff`: 销售差价（通常为正数）
- 最终价格自动计算:
  - 回购价 = base_price + buy_price_diff
  - 销售价 = base_price + sell_price_diff

**curl 示例**:
```bash
curl -X POST http://localhost:8080/api/prices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "code": "gold_9999",
    "name": "黄金9999",
    "subtitle": "Au9999 · 千足金",
    "icon": "Au",
    "icon_color": "#FFD700",
    "base_price": 560.00,
    "buy_price_diff": -10.00,
    "sell_price_diff": 15.00,
    "sort": 1,
    "status": 1
  }'
```

### 5.4 更新价格

**接口**: `PUT /api/prices/:id`

**请求参数**:
```json
{
  "base_price": 565.00,
  "buy_price_diff": -12.00,
  "sell_price_diff": 18.00
}
```

### 5.5 删除价格

**接口**: `DELETE /api/prices/:id`

### 5.6 同步基础价格

**接口**: `POST /api/prices/sync`

**说明**: 从第三方API同步最新的贵金属基础价格（待实现）

---

## 错误码说明

| 错误码 | 说明 |
|-------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（未登录或Token失效） |
| 403 | 禁止访问（无权限） |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 开发建议

### Postman 测试

1. 创建环境变量:
   - `base_url`: `http://localhost:8080/api`
   - `token`: 从登录接口获取

2. 在请求中使用:
   - URL: `{{base_url}}/users`
   - Header: `Authorization: Bearer {{token}}`

### 前端集成

```javascript
// axios 配置示例
import axios from 'axios'

const service = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 5000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      // 错误处理
      return Promise.reject(new Error(res.message))
    }
    return res.data
  },
  error => {
    return Promise.reject(error)
  }
)

export default service
```















