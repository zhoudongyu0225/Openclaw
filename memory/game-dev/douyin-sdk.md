# 抖音小游戏接入 - 笔记

## 开发者平台
https://developer.open-douyin.com/

## 资质要求
- 企业认证
- 营业执照
- 相关许可证

## 接口能力
- 用户登录
- 支付
- 分享
- 关系链

## 开发流程
1. 注册开发者账号
2. 创建应用
3. 申请接口权限
4. 开发调试
5. 提交审核
6. 上线发布

---

## 2026-03-07 完善：直播间 SDK 接入

### 核心能力矩阵

| 能力 | 用途 | 权限要求 |
|------|------|----------|
| 获取用户信息 | 登录/头像展示 | scope/user_info |
| 充值金币 | 礼物打赏 | scope/payment |
| 发送弹幕 | 互动消息 | 无 |
| 排行榜 | 榜单展示 | scope/leaderboard |
| 实时语音 | 连麦互动 | scope/voice |

### API 端点 (服务端)

```go
// 抖音开放平台 API
const (
    DouyinAPIHost      = "https://open.douyin.com"
    DouyinAuthURL      = "/oauth/access_token/"
    DouyinUserInfoURL  = "/oauth/userinfo/"
    DouyinPaymentURL  = "/order/create/"
)

// 鉴权
type DouyinAuth struct {
    ClientKey    string `json:"client_key"`
    ClientSecret string `json:"client_secret"`
    Code         string `json:"code"`
    GrantType    string `json:"grant_type"` // authorization_code
}

// 用户信息
type DouyinUserInfo struct {
    OpenID   string `json:"open_id"`
    Nickname string `json:"nickname"`
    Avatar   string `json:"avatar_url"`
    Gender   int    `json:"gender"` // 0未知 1男 2女
}
```

### 直播间推送 Webhook

```go
// 抖音直播间事件推送
type DouyinLiveEvent struct {
    EventType string    `json:"event_type"` // gift/danmaku/follow
    RoomID    string    `json:"room_id"`
    User      UserInfo  `json:"user"`
    Content   string    `json:"content,omitempty"`
    Gift      *GiftInfo `json:"gift,omitempty"`
    Timestamp int64     `json:"timestamp"`
}

type GiftInfo struct {
    GiftID   string `json:"gift_id"`
    GiftName string `json:"gift_name"`
    Count    int    `json:"count"`
    Value    int    `json:"value"`
}
```

### 对接流程

```
1. 前端获取授权码
   抖音APP → 扫码授权 → 获取 code

2. 服务端换 token
   POST /oauth/access_token/ 
   → access_token + open_id

3. 获取用户信息
   GET /oauth/userinfo/?access_token=xxx&open_id=xxx

4. 推送事件接收
   配置 Webhook URL
   → 接收礼物/弹幕/关注事件

5. 订单支付
   POST /order/create/
   → 唤起抖音支付
```

### 游戏内接入点

```go
// 直播间礼物事件处理
func (lr *LiveRoom) HandleDouyinGift(event *DouyinLiveEvent) {
    giftType := mapDouyinGift(event.Gift.GiftID)
    effect := lr.SendGift(event.User.OpenID, giftType)
    
    // 广播特效到前端
    lr.BroadcastGiftEffect(effect)
}

// 弹幕事件处理
func (lr *LiveRoom) HandleDouyinDanmaku(event *DouyinLiveEvent) {
    lr.SendDanmaku(event.User.OpenID, event.Content)
}
```

### 待接入功能

- [ ] 抖音登录授权流程
- [ ] 支付订单创建
- [ ] Webhook 事件接收服务
- [ ] 排行榜数据上报
- [ ] 实时语音连麦 (可选)
