# 沪金汇 - 开发文档

## 📱 项目说明

**沪金汇（HJH）**是一个完整的微信小程序项目，包含贵金属回收店铺的所有核心功能。

## 🚀 快速开始

### 1. 安装微信开发者工具

下载并安装[微信开发者工具](https://developers.weixin.qq.com/miniprogram/dev/devtools/download.html)

### 2. 导入项目

1. 打开微信开发者工具
2. 点击"+"或"导入项目"
3. 选择 `miniprogram` 目录
4. 填写AppID（测试可使用测试号）
5. 点击"导入"

### 3. 运行项目

导入后会自动编译运行，可以在模拟器中预览效果。

## 📂 项目结构

```
miniprogram/
├── app.js                 # 小程序入口文件
├── app.json               # 小程序配置
├── app.wxss               # 全局样式
├── sitemap.json           # 搜索配置
├── project.config.json    # 项目配置
├── pages/                 # 页面目录
│   ├── index/            # 报价首页
│   │   ├── index.wxml    # 页面结构
│   │   ├── index.js      # 页面逻辑
│   │   ├── index.wxss    # 页面样式
│   │   └── index.json    # 页面配置
│   ├── appointment/       # 预约页面
│   ├── calculator/        # 计算器页面
│   └── profile/           # 个人中心
└── images/                # 图片资源（需要自行添加）
```

## 🎯 功能说明

### 1. 报价首页 (pages/index)

**功能**：
- 展示12种贵金属的回购价、销售价、高低价
- 分类筛选（全部、黄金、白银、铂金、钯金、国际、汇率）
- 下拉刷新价格
- 店铺信息展示
- 免责声明

**支持的品种**：
- 黄金9999
- 黄金T+D
- 黄金
- 白银
- 钯金
- 铂金
- 白银T+D
- 美黄金
- 美白银
- 伦敦金
- 伦敦银
- 美元兑换

**数据来源**：
目前使用模拟数据（`app.js` 中的 `updatePriceData` 方法），实际开发中需要对接后端API。

### 2. 预约回收 (pages/appointment)

**功能**：
- 选择回收品种
- 选择预约方式（到店/上门）
- 选择日期和时间
- 填写联系信息
- 表单验证

**数据提交**：
目前是模拟提交，实际开发中需要对接后端API。

### 3. 价格计算器 (pages/calculator)

**功能**：
- 选择贵金属品种
- 输入重量（克）
- 输入成色系数
- 实时计算回收价值
- 保存计算历史（本地存储）
- 查看和复用历史记录

**计算公式**：
```
回收价值 = 重量(克) × 回购价(元/克) × 成色系数(%)
```

### 4. 个人中心 (pages/profile)

**功能**：
- 用户信息展示
- 统计数据
- 我的预约（跳转）
- 回收记录
- 计算历史
- 优惠券
- 门店地址导航
- 联系客服

## 🔧 核心技术

### 数据存储

使用微信小程序本地存储（`wx.setStorageSync` / `wx.getStorageSync`）：

```javascript
// 存储价格数据
wx.setStorageSync('priceData', priceData);

// 读取价格数据
const priceData = wx.getStorageSync('priceData');

// 存储计算历史
wx.setStorageSync('calcHistory', history);
```

### 页面通信

使用 `app.globalData` 共享全局数据：

```javascript
// 设置全局数据
const app = getApp();
app.globalData.storeInfo = {...};

// 获取全局数据
const storeInfo = app.globalData.storeInfo;
```

### API调用（待实现）

实际开发中需要对接后端API：

```javascript
// 获取价格列表
wx.request({
  url: 'https://your-api.com/api/prices',
  method: 'GET',
  success: (res) => {
    console.log(res.data);
  }
});

// 提交预约
wx.request({
  url: 'https://your-api.com/api/appointments',
  method: 'POST',
  data: {...},
  success: (res) => {
    console.log(res.data);
  }
});
```

## 🎨 样式定制

### 主题颜色

在 `app.wxss` 中修改：

```css
/* 主色调 - 金色渐变 */
background: linear-gradient(135deg, #FFD700 0%, #FFA500 100%);

/* 成功色 - 绿色 */
color: #22C55E;

/* 危险色 - 红色 */
color: #FF4444;
```

### 顶部导航栏

在 `app.json` 中修改：

```json
"window": {
  "navigationBarTitleText": "贵金属回收",
  "navigationBarBackgroundColor": "#FFD700",
  "navigationBarTextStyle": "white"
}
```

## 📊 数据格式

### 价格数据格式

```javascript
{
  updateTime: "2025-11-03 14:32:15",
  list: [
    {
      id: 1,
      code: 'gold_9999',          // 品种代码
      name: '黄金9999',            // 品种名称
      subtitle: 'Au9999 · 千足金', // 副标题
      icon: 'Au',                  // 图标文字
      iconColor: '#FFD700',        // 图标颜色
      buyPrice: 558.50,            // 回购价
      sellPrice: 568.80,           // 销售价
      highPrice: 570.20,           // 最高价
      lowPrice: 556.30             // 最低价
    }
  ]
}
```

### 计算历史格式

```javascript
{
  id: 1699008000000,               // 时间戳
  metalName: '黄金9999',           // 品种名称
  metalCode: 'gold_9999',          // 品种代码
  buyPrice: 558.50,                // 回购价
  weight: 50,                      // 重量
  purity: 99.9,                    // 成色
  result: '27925.00',              // 计算结果
  formula: '50克 × ¥558.50/克 × 99.9%', // 公式
  time: '2025-11-03 14:32'        // 时间
}
```

## 🔌 需要对接的API

### 1. 获取价格列表
```
GET /api/prices
```

### 2. 创建预约
```
POST /api/appointments
Body: {
  metalType, serviceType, date, time,
  name, phone, address, note
}
```

### 3. 获取预约列表
```
GET /api/appointments?userId={userId}
```

### 4. 获取回收记录
```
GET /api/records?userId={userId}
```

### 5. 微信登录
```
POST /api/auth/login
Body: { code }
```

## 📝 待开发功能

- [ ] 后端API对接
- [ ] 微信登录
- [ ] 真实预约管理
- [ ] 回收记录管理
- [ ] 优惠券系统
- [ ] 消息推送
- [ ] 图片上传
- [ ] 在线客服

## ⚠️ 注意事项

### 1. 图标资源

✅ **已解决**：本项目使用 Emoji 和文字代替图片图标，无需准备任何图片资源！

- TabBar：💰 报价、📱 预约、🧮 计算器、👤 我的
- 用户头像：👤

### 2. 个人二维码

⚡ **重要功能**：首页右上角显示个人二维码
- **点击**：放大预览
- **长按**：识别二维码

**如何添加**：
1. 准备你的二维码图片（200x200px）
2. 放到 `miniprogram/images/qrcode.png`
3. 修改 `pages/index/index.js` 第29行：
   ```javascript
   qrcodeUrl: '/images/qrcode.png'
   ```

📖 详细说明：查看 [如何添加二维码.md](./如何添加二维码.md)

### 3. AppID配置

在 `project.config.json` 中将 `"appid"` 修改为你的小程序AppID：
```json
"appid": "你的AppID"
```

### 3. 服务器域名

实际上线前需要在微信公众平台配置服务器域名。

## 🚀 部署上线

1. 完善功能代码
2. 对接后端API
3. 准备所需图标
4. 在微信公众平台申请小程序
5. 配置服务器域名
6. 上传代码审核
7. 审核通过后发布

## 📞 技术支持

如有问题，欢迎反馈！

---

**开发环境**: 微信开发者工具 v1.06+  
**基础库版本**: 3.0.0+  
**开发语言**: JavaScript  
**当前版本**: v1.0.0

