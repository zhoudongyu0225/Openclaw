# 🛒 氢镁胶囊 - Stripe/PayPal接入方案

## 支付接入文档

### 方案对比

| 方案 | 支持地区 | 手续费 | 特点 |
|------|----------|--------|------|
| **Stripe** | 135+国家 | 2.9%+$0.30 | 功能强，需技术对接 |
| **PayPal** | 200+国家 | 3.4%+$0.30 | 用户基数大，转化高 |
| **Shopify Payments** | 17国 | 2.9%+$0.30 | Shopify生态 |

### Stripe 接入步骤

#### 1. 注册账号
1. 访问 stripe.com
2. 使用邮箱注册
3. 完成企业认证

#### 2. 技术集成

**HTML示例代码：**

```html
<script src="https://js.stripe.com/v3/"></script>

<form id="payment-form">
  <div id="card-element"></div>
  <button type="submit">支付 ¥2988</button>
</form>

<script>
  var stripe = Stripe('your_publishable_key');
  var elements = stripe.elements();
  var card = elements.create('card');
  card.mount('#card-element');
  
  form.addEventListener('submit', function(event) {
    event.preventDefault();
    stripe.createToken(card).then(function(result) {
      // 处理支付
    });
  });
</script>
```

#### 3. webhook 配置
- 支付成功回调
- 支付失败回调
- 退款通知

### PayPal 接入步骤

#### 1. 注册商家账号
1. 访问 business.paypal.com
2. 完成企业认证

#### 2. 获取API凭证
- Client ID
- Secret

#### 3. HTML集成

```html
<script src="https://www.paypal.com/sdk/js?client-id=YOUR_CLIENT_ID&currency=USD"></script>

<div id="paypal-button-container"></div>

<script>
  paypal.Buttons({
    createOrder: function(data, actions) {
      return actions.order.create({
        purchase_units: [{
          amount: { value: '428' }
        }]
      });
    },
    onApprove: function(data, actions) {
      return actions.order.capture().then(function(details) {
        alert('支付成功！');
      });
    }
  }).render('#paypal-button-container');
</script>
```

### 推荐的支付方案

#### 方案A：Stripe（主推）

**优势：**
- 支持135+货币
- API 文档完善
- 支付成功率高
- 支持订阅（复购）

**接入难度：** 中等（需要技术人员）

#### 方案B：Stripe + PayPal 双通道

**配置：**
- 主支付：Stripe（信用卡）
- 备支付：PayPal（用户可选）

**代码示例：**

```html
<div class="payment-options">
  <button id="stripe-btn" class="active">信用卡支付</button>
  <button id="paypal-btn">PayPal支付</button>
</div>

<div id="stripe-form">
  <!-- Stripe 表单 -->
</div>

<div id="paypal-form" style="display:none;">
  <!-- PayPal 按钮 -->
</div>
```

### 合规要求

#### 1. 营养品销售声明

必须添加以下免责声明：

```text
声明：本产品为膳食补充剂，不用于诊断、治疗、治愈或预防任何疾病。
声明：本产品根据FDA规定生产，但未经FDA评估。
```

#### 2. 退款政策

建议：
- 30天无理由退款
- 运费由商家承担
- 客服响应时间：24小时内

#### 3. 运输政策

明确：
- 发货时间：1-3个工作日
- 物流方式：顺丰/DHL
- 送达时间：7-15天
- 关税：买家承担

### 下一步

1. [ ] 注册 Stripe 商家账户
2. [ ] 配置 Stripe API
3. [ ] 集成到 index-en.html
4. [ ] 测试支付流程

---

*产出时间：2026-03-02*
