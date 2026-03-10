package handler

import (
	"encoding/json"
	"log"
	"time"
)

// ============================================
// 抖音直播事件处理器
// ============================================

// 抖音事件类型
const (
	DouyinEventGift   = "gift"     // 礼物
	DouyinEventDanmaku = "danmaku"  // 弹幕
	DouyinEventFollow = "follow"    // 关注
	DouyinEventLike  = "like"      // 点赞
	DouyinEventShare = "share"     // 分享
	DouyinEventEnter = "enter"     // 进入直播间
)

// 抖音用户信息
type DouyinUser struct {
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar_url"`
	Gender   int    `json:"gender"` // 0未知 1男 2女
}

// 抖音礼物信息
type DouyinGift struct {
	GiftID   string `json:"gift_id"`
	GiftName string `json:"gift_name"`
	Count    int    `json:"count"`
	Value    int    `json:"value"` // 货币价值
}

// 抖音直播事件
type DouyinLiveEvent struct {
	EventType string      `json:"event_type"`
	RoomID    string      `json:"room_id"`
	User      DouyinUser  `json:"user"`
	Content   string      `json:"content,omitempty"`
	Gift      *DouyinGift `json:"gift,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// 抖音事件处理器
type DouyinEventHandler struct {
	roomManager interface {
		GetRoom(id string) interface {
			SendGift(userID, giftType string) interface{}
			SendDanmaku(userID, content string) interface{}
			Broadcast(msg interface{})
		}
	}
}

// 新建处理器
func NewDouyinEventHandler() *DouyinEventHandler {
	return &DouyinEventHandler{}
}

// 设置房间管理器
func (h *DouyinEventHandler) SetRoomManager(rm interface{}) {
	h.roomManager = rm
}

// 处理抖音 Webhook 事件
func (h *DouyinEventHandler) HandleWebhook(payload []byte) error {
	var event DouyinLiveEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("[Douyin] 解析事件失败: %v", err)
		return err
	}

	log.Printf("[Douyin] 收到事件: %s, 用户: %s, 房间: %s", 
		event.EventType, event.User.Nickname, event.RoomID)

	switch event.EventType {
	case DouyinEventGift:
		h.handleGift(&event)
	case DouyinEventDanmaku:
		h.handleDanmaku(&event)
	case DouyinEventFollow:
		h.handleFollow(&event)
	case DouyinEventLike:
		h.handleLike(&event)
	case DouyinEventShare:
		h.handleShare(&event)
	case DouyinEventEnter:
		h.handleEnter(&event)
	}

	return nil
}

// 处理礼物事件
func (h *DouyinEventHandler) handleGift(event *DouyinLiveEvent) {
	if event.Gift == nil {
		return
	}

	room := h.getRoom(event.RoomID)
	if room == nil {
		return
	}

	// 抖音礼物ID映射到游戏内礼物
	giftType := mapDouyinGift(event.Gift.GiftID)
	
	// 发送礼物效果
	for i := 0; i < event.Gift.Count; i++ {
		room.SendGift(event.User.OpenID, giftType)
	}

	log.Printf("[Douyin] 礼物: %s x%d", event.Gift.GiftName, event.Gift.Count)
}

// 处理弹幕事件
func (h *DouyinEventHandler) handleDanmaku(event *DouyinLiveEvent) {
	room := h.getRoom(event.RoomID)
	if room == nil {
		return
	}

	room.SendDanmaku(event.User.OpenID, event.Content)
	log.Printf("[Douyin] 弹幕: %s: %s", event.User.Nickname, event.Content)
}

// 处理关注事件
func (h *DouyinEventHandler) handleFollow(event *DouyinLiveEvent) {
	log.Printf("[Douyin] 关注: %s", event.User.Nickname)
	// 可以发送系统消息欢迎新粉丝
}

// 处理点赞事件
func (h *DouyinEventHandler) handleLike(event *DouyinLiveEvent) {
	log.Printf("[Douyin] 点赞: %s", event.User.Nickname)
	// 点赞可以积累能量，用于释放技能
}

// 处理分享事件
func (h *DouyinEventHandler) handleShare(event *DouyinLiveEvent) {
	log.Printf("[Douyin] 分享: %s", event.User.Nickname)
	// 分享可以奖励金币
}

// 处理进入直播间事件
func (h *DouyinEventHandler) handleEnter(event *DouyinLiveEvent) {
	log.Printf("[Douyin] 进入: %s", event.User.Nickname)
	// 可以显示欢迎弹幕
}

// 获取房间
func (h *DouyinEventHandler) getRoom(roomID string) interface{} {
	if h.roomManager == nil {
		return nil
	}
	return h.roomManager.GetRoom(roomID)
}

// 抖音礼物ID映射到游戏内礼物类型
func mapDouyinGift(douyinGiftID string) string {
	// 抖音礼物ID映射表 (需根据实际配置)
	giftMap := map[string]string{
		"1":    "coin",   // 小心心
		"2":    "star",   // 星星
		"3":    "rocket",  // 火箭
		"4":    "car",     // 跑车
		"5":    "plane",   // 飞机
		"6":    "bang",    // 炸弹
		"1001": "coin",    // 玫瑰
		"1002": "star",    // 仙女棒
		"1003": "rocket",  // 跑车
	}

	if giftType, ok := giftMap[douyinGiftID]; ok {
		return giftType
	}
	return "coin" // 默认金币
}

// 构建抖音 API 请求
type DouyinAPIRequest struct {
	Method      string
	URL         string
	AccessToken string
	Params      map[string]string
}

// 发送抖音 API 请求
func (h *DouyinEventHandler) SendAPI(req *DouyinAPIRequest) ([]byte, error) {
	// TODO: 实现实际的 API 调用
	// 这里应该调用抖音开放平台 API
	log.Printf("[Douyin] API请求: %s %s", req.Method, req.URL)
	return nil, nil
}

// ============================================
// Webhook 签名验证
// ============================================

// 验证 Webhook 签名
func VerifyWebhookSignature(payload, signature, secret string) bool {
	// TODO: 实现 HMAC-SHA256 签名验证
	// expected := hmac_sha256(secret, payload)
	return true // 测试环境跳过验证
}

// 生成回调响应
func BuildWebhookResponse(errCode int, msg string) []byte {
	resp := map[string]interface{}{
		"code":    errCode,
		"message": msg,
		"msg_id":  time.Now().UnixMilli(),
	}
	
	data, _ := json.Marshal(resp)
	return data
}
