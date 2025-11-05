# 沪金汇管理系统 - 前端

基于 Vue 2 + Element UI 的贵金属回收管理系统前端。

## 技术栈

- **Vue 2.6** - 渐进式 JavaScript 框架
- **Element UI 2.15** - 基于 Vue 2 的组件库
- **Vue Router 3.x** - 官方路由管理器
- **Vuex 3.x** - 状态管理
- **Axios** - HTTP 客户端
- **SCSS** - CSS 预处理器

## 功能特性

- ✅ 用户登录/登出
- ✅ 权限管理（基于角色的访问控制）
- ✅ 用户管理
- ✅ 角色管理
- ✅ 菜单管理
- ✅ 价格管理（基础价 + 差价模式）
- 🔄 响应式布局
- 🔄 动态路由
- 🔄 请求拦截
- 🔄 错误处理

## 目录结构

```
gold-admin-frontend/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API 接口
│   │   ├── auth.js        # 认证接口
│   │   ├── user.js        # 用户管理接口
│   │   ├── role.js        # 角色管理接口
│   │   ├── menu.js        # 菜单管理接口
│   │   └── price.js       # 价格管理接口
│   ├── assets/            # 静态资源
│   ├── components/        # 公共组件
│   ├── layout/            # 布局组件
│   │   ├── components/    # 布局子组件
│   │   │   ├── Navbar.vue
│   │   │   ├── Sidebar/
│   │   │   └── AppMain.vue
│   │   └── index.vue
│   ├── router/            # 路由配置
│   ├── store/             # Vuex 状态管理
│   │   ├── modules/
│   │   │   └── user.js
│   │   ├── getters.js
│   │   └── index.js
│   ├── styles/            # 全局样式
│   ├── utils/             # 工具函数
│   │   ├── request.js     # axios 封装
│   │   ├── auth.js        # 认证工具
│   │   └── validate.js    # 验证工具
│   ├── views/             # 页面组件
│   │   ├── login/         # 登录页
│   │   ├── dashboard/     # 首页
│   │   ├── system/        # 系统管理
│   │   │   ├── user/      # 用户管理
│   │   │   ├── role/      # 角色管理
│   │   │   └── menu/      # 菜单管理
│   │   ├── price/         # 价格管理
│   │   └── error-page/    # 错误页
│   ├── App.vue
│   ├── main.js
│   └── permission.js      # 权限控制
├── .env.development       # 开发环境变量
├── .env.production        # 生产环境变量
├── vue.config.js          # Vue CLI 配置
└── package.json
```

## 快速开始

### 1. 安装依赖

```bash
cd gold-admin-frontend
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
# 或
npm run serve
```

访问 `http://localhost:9527`

默认账号：
- 用户名: `admin`
- 密码: `admin123`

### 3. 编译生产环境

```bash
npm run build
```

编译后的文件在 `dist/` 目录

## 环境变量

### 开发环境 (`.env.development`)

```env
NODE_ENV=development
VUE_APP_BASE_API=/api
```

### 生产环境 (`.env.production`)

```env
NODE_ENV=production
VUE_APP_BASE_API=/api
```

## 开发规范

### 1. 命名规范

- **组件名**: 使用 PascalCase (如 `UserList.vue`)
- **文件名**: 使用 kebab-case (如 `user-list.vue`)
- **路由名**: 使用 PascalCase (如 `UserList`)
- **API 方法**: 使用 camelCase (如 `getUserList`)

### 2. 代码风格

项目使用 ESLint + Standard 规范

```bash
npm run lint
```

### 3. Git 提交规范

```
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式调整
refactor: 重构
test: 测试相关
chore: 构建/工具链相关
```

## 核心功能

### 1. 用户认证

登录后，Token 会存储在 Cookie 中，并在每次请求时自动携带。

```javascript
// 登录
this.$store.dispatch('user/login', { username, password })

// 登出
this.$store.dispatch('user/logout')

// 获取用户信息
this.$store.dispatch('user/getInfo')
```

### 2. 权限控制

使用路由守卫进行权限控制，未登录用户会被重定向到登录页。

```javascript
// src/permission.js
router.beforeEach(async(to, from, next) => {
  // 权限验证逻辑
})
```

### 3. 请求拦截

所有 API 请求都会经过拦截器处理：

- 请求拦截：自动携带 Token
- 响应拦截：统一错误处理

```javascript
// src/utils/request.js
service.interceptors.request.use(/* ... */)
service.interceptors.response.use(/* ... */)
```

### 4. 价格管理

价格管理采用 **基础价 + 差价** 模式：

- **基础价**: 从市场获取的标准价格
- **回购差价**: 通常为负数（低于基础价收购）
- **销售差价**: 通常为正数（高于基础价出售）

计算公式：
- 回购价 = 基础价 + 回购差价
- 销售价 = 基础价 + 销售差价

## 部署

### 1. 本地预览

```bash
npm run build
npm install -g serve
serve -s dist
```

### 2. Nginx 部署

```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /var/www/gold-admin-frontend/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. 嵌入 Go 后端

如果使用方案A（前后端打包在一起），将 `dist/` 目录内容复制到后端项目的 `web/` 目录。

## 常见问题

### 1. 开发环境跨域问题

已在 `vue.config.js` 中配置代理：

```javascript
devServer: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

### 2. 生产环境 API 地址

修改 `.env.production` 中的 `VUE_APP_BASE_API`

### 3. 路由 404 问题

使用 history 模式时，需要在服务器配置重写规则，将所有请求指向 `index.html`

## 浏览器支持

- Chrome (推荐)
- Firefox
- Safari
- Edge

不支持 IE

## 许可证

MIT

## 联系方式

- 作者: javapub
- 项目地址: https://github.com/javapub/mini-goldprice





