# 逆动力氢镁胶囊 - 部署指南

## 快速部署（推荐）

### 方案1: Vercel 一键部署（免费）

1. 访问: https://vercel.com/new
2. 使用 GitHub 账号登录
3. 点击 "Import Project"
4. 上传 `projects/nidongli` 文件夹到你的 GitHub 仓库
5. 在 Vercel 导入该仓库
6. 点击 Deploy！

部署完成后，你会得到一个 URL，比如：`https://nidongli-xxx.vercel.app`

### 方案2: Netlify 一键部署（免费）

1. 访问: https://app.netlify.com/drop
2. 直接拖拽 `projects/nidongli` 文件夹到页面
3. 自动部署，得到 URL

---

## 支付功能配置

### PayPal 支付链接
1. 访问 https://www.paypal.com/paypalme/你的用户名
2. 创建支付链接，替换页面中的支付按钮

### Stripe 支付链接
1. 访问 https://dashboard.stripe.com/payment-links
2. 创建支付链接
3. 将链接嵌入页面

### 微信/支付宝
需要商家账号，建议先使用 PayPal 或 Stripe 测试

---

## 下一步
1. 部署上线
2. 配置支付
3. 测试购买流程
4. 推广销售
