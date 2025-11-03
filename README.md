# 沪金汇 - 贵金属回收小程序

## 🏪 项目简介

**沪金汇（HJH）**是一个贵金属回收店铺的微信小程序，用于展示店铺的回购和销售价格，提供在线预约、价格计算等服务。适用于黄金回收店、典当行、贵金属交易商等实体店铺。

## 🎯 项目定位

- **业务类型**：贵金属回收/销售商店
- **目标用户**：需要回收或购买贵金属的客户
- **核心价值**：展示实时报价、提供便捷预约、快速价值估算
- **使用场景**：查看报价、预约回收、计算价值、查询记录

## ✨ 核心功能

### 1. 报价展示（首页）
- **回购价格**：店铺收购贵金属的价格
- **销售价格**：店铺销售贵金属的价格
- **当日高/低价**：今日价格波动范围
- **多品种支持**：黄金、白银、铂金、钯金等
- **店铺信息**：地址、电话、营业时间、导航
- **免责声明**：价格说明、风险提示

### 2. 预约回收
- **品种选择**：选择要回收的贵金属类型
- **预约方式**：到店回收或上门回收
- **时间选择**：灵活选择预约日期和时间段
- **信息填写**：姓名、电话、地址、备注
- **预约记录**：查看历史预约状态

### 3. 价格计算器
- **品种选择**：选择不同的贵金属
- **重量输入**：输入物品重量（克）
- **成色系数**：可选输入成色百分比
- **实时计算**：自动计算预估回收价值
- **计算历史**：保存计算记录

### 4. 个人中心
- **用户信息**：微信头像、昵称
- **预约管理**：查看和管理预约记录
- **回收记录**：历史回收交易记录
- **计算历史**：查看价格计算记录
- **优惠券**：查看可用优惠券
- **设置帮助**：通知设置、常见问题、门店地址

## 📂 项目结构

```
mini-goldprice/
├── prototype/                 # HTML原型图
│   ├── start.html            # 导航页（推荐从这里开始）
│   ├── index.html            # 报价首页
│   ├── appointment.html      # 预约回收页
│   ├── calculator.html       # 价格计算器页
│   ├── profile.html          # 个人中心页
│   └── README.md             # 原型说明文档
└── README.md                 # 项目说明文档
```

## 🎨 原型图预览

### 快速开始
```bash
# 在浏览器中打开导航页
open prototype/start.html
```

### 页面说明
1. **报价首页（index.html）**
   - 展示回购价、销售价、高低价
   - 多品种卡片展示
   - 店铺信息（地址、电话、营业时间）
   - 免责声明
   - 底部导航：报价、预约、计算器、我的

2. **预约回收（appointment.html）**
   - 选择回收品种（黄金、白银、铂金、钯金）
   - 选择预约方式（到店/上门）
   - 选择日期和时间段
   - 填写联系信息
   - 查看预约记录

3. **价格计算器（calculator.html）**
   - 选择贵金属品种
   - 输入重量（克）
   - 输入成色系数（可选）
   - 实时显示计算结果
   - 查看计算历史

4. **个人中心（profile.html）**
   - 用户信息和统计
   - 我的预约、回收记录、计算历史
   - 优惠券管理
   - 消息通知、门店地址、常见问题
   - 关于和反馈

## 🚀 技术栈建议

### 前端（微信小程序）

#### 推荐方案一：原生小程序
```
- 微信小程序原生开发
- Vant Weapp（UI组件库）
- ECharts（可选，如需图表）
```

#### 推荐方案二：跨平台框架
```
- uni-app（Vue语法，可跨平台）
- Taro（React语法）
```

### 后端服务

#### 方案一：Node.js
```
技术栈：
- Express / Koa（Web框架）
- MySQL / MongoDB（数据库）
- Redis（缓存）
- JWT（用户认证）

主要功能：
- 用户管理（微信登录）
- 价格管理（手动更新报价）
- 预约管理（创建、查询、更新）
- 计算记录
- 消息推送（预约提醒）
```

#### 方案二：Java
```
技术栈：
- Spring Boot（框架）
- MyBatis（ORM）
- MySQL（数据库）
- Redis（缓存）

主要功能：
- 同上
```

### 数据库设计

```sql
-- 用户表
users (id, openid, nickname, avatar, phone, created_at)

-- 价格表
prices (id, metal_type, buy_price, sell_price, high_price, low_price, updated_at)

-- 预约表
appointments (id, user_id, metal_type, service_type, appointment_time, name, phone, address, note, status, created_at)

-- 计算记录表
calculations (id, user_id, metal_type, weight, purity, result, created_at)

-- 回收记录表
recycling_records (id, user_id, metal_type, weight, price, total_amount, created_at)

-- 优惠券表
coupons (id, user_id, title, discount, status, expire_at, created_at)
```

## 🔄 业务流程

### 用户流程
```
1. 打开小程序 → 查看当日报价
2. 使用计算器 → 估算物品价值
3. 在线预约 → 选择时间和方式
4. 接到电话 → 店铺确认预约
5. 到店/上门 → 现场鉴定和交易
6. 完成交易 → 记录保存到系统
```

### 店铺管理流程
```
1. 登录管理后台
2. 更新当日报价（回购价、销售价）
3. 查看预约列表
4. 联系客户确认预约
5. 完成交易后录入系统
6. 查看营业数据统计
```

## 💡 扩展功能建议

### 基础版（MVP）
- [x] 价格展示
- [x] 在线预约
- [x] 价格计算器
- [x] 个人中心

### 增强版
- [ ] 优惠券系统
- [ ] 积分系统
- [ ] 会员等级
- [ ] 消息推送
- [ ] 在线客服
- [ ] 图片上传（拍照鉴定）

### 高级版
- [ ] 多店铺管理
- [ ] 店铺管理后台
- [ ] 数据统计分析
- [ ] 营销活动管理
- [ ] 财务对账系统
- [ ] 鉴定报告生成

## ⚠️ 重要提示

### 合规性要求
1. **营业资质**：需要有合法的贵金属回收营业执照
2. **实名认证**：需要实名认证客户身份
3. **交易记录**：需要保存完整的交易记录
4. **免责声明**：必须包含价格波动、仅供参考等说明
5. **用户协议**：需要完善的用户协议和隐私政策

### 价格管理
1. **更新频率**：建议每日更新报价，跟随市场行情
2. **价格来源**：可参考国际金价、上海黄金交易所等
3. **价差设置**：回购价和销售价需要保持合理价差
4. **成色鉴定**：需要专业设备鉴定贵金属成色
5. **重量复核**：需要精确电子秤称重

### 用户体验
1. **简洁明了**：价格展示要清晰，不要过于复杂
2. **即时响应**：预约后需要及时联系确认
3. **专业服务**：提供专业的鉴定和咨询服务
4. **诚信经营**：价格公道，不压价欺诈
5. **安全保障**：保护客户隐私和交易安全

## 📊 数据接口示例

### 获取价格列表
```javascript
GET /api/prices

Response:
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "metal_type": "gold_9999",
      "name": "黄金 Au9999",
      "buy_price": 558.50,
      "sell_price": 568.80,
      "high_price": 570.20,
      "low_price": 556.30,
      "updated_at": "2025-11-03 14:32:15"
    }
  ]
}
```

### 创建预约
```javascript
POST /api/appointments

Request:
{
  "metal_type": "gold",
  "service_type": "store",  // store: 到店, home: 上门
  "appointment_time": "2025-11-04 10:00",
  "name": "张三",
  "phone": "13800000000",
  "address": "",  // 到店回收可为空
  "note": "黄金项链约30克"
}

Response:
{
  "code": 200,
  "message": "预约成功，工作人员会尽快联系您",
  "data": {
    "id": 123,
    "status": "pending"
  }
}
```

### 计算价格
```javascript
POST /api/calculate

Request:
{
  "metal_type": "gold_9999",
  "weight": 50,
  "purity": 99.9
}

Response:
{
  "code": 200,
  "data": {
    "buy_price": 558.50,
    "weight": 50,
    "purity": 99.9,
    "result": 27925.00,
    "formula": "50克 × ¥558.50/克 × 99.9%"
  }
}
```

## 📱 小程序配置

### app.json 配置示例
```json
{
  "pages": [
    "pages/index/index",
    "pages/appointment/appointment",
    "pages/calculator/calculator",
    "pages/profile/profile"
  ],
  "window": {
    "navigationBarTitleText": "贵金属回收",
    "navigationBarBackgroundColor": "#FFD700",
    "navigationBarTextStyle": "white"
  },
  "tabBar": {
    "list": [
      {
        "pagePath": "pages/index/index",
        "text": "报价",
        "iconPath": "images/icon-home.png",
        "selectedIconPath": "images/icon-home-active.png"
      },
      {
        "pagePath": "pages/appointment/appointment",
        "text": "预约",
        "iconPath": "images/icon-appointment.png",
        "selectedIconPath": "images/icon-appointment-active.png"
      },
      {
        "pagePath": "pages/calculator/calculator",
        "text": "计算器",
        "iconPath": "images/icon-calculator.png",
        "selectedIconPath": "images/icon-calculator-active.png"
      },
      {
        "pagePath": "pages/profile/profile",
        "text": "我的",
        "iconPath": "images/icon-profile.png",
        "selectedIconPath": "images/icon-profile-active.png"
      }
    ]
  }
}
```

## 🎯 开发计划

### 第一阶段：原型设计 ✅
- [x] 完成页面原型图设计
- [x] 确定功能模块
- [x] 设计UI/UX

### 第二阶段：基础开发
- [ ] 搭建小程序项目
- [ ] 开发后端API服务
- [ ] 实现基础页面
- [ ] 微信登录集成

### 第三阶段：核心功能
- [ ] 价格管理系统
- [ ] 预约功能
- [ ] 计算器功能
- [ ] 个人中心

### 第四阶段：优化上线
- [ ] 功能测试
- [ ] 性能优化
- [ ] 小程序审核
- [ ] 正式上线

## 📞 联系方式

如有疑问或需要定制开发，欢迎联系！

---

**项目类型**: 贵金属回收小程序  
**当前版本**: v1.0.0 (原型阶段)  
**创建日期**: 2025-11-03  
**适用行业**: 黄金回收、典当行、贵金属交易
