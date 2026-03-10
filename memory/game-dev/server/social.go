package game

import (
	"fmt"
	"sync"
	"time"
)

// ==================== 聊天系统 ====================

// ChatChannelType 聊天频道类型
type ChatChannelType int

const (
	ChatChannelWorld  ChatChannelType = iota // 世界频道
	ChatChannelGuild                         // 公会频道
	ChatChannelPrivate                       // 私聊
	ChatChannelSystem                       // 系统频道
	ChatChannelBattle                       // 战斗频道
)

// ChatMessage 聊天消息
type ChatMessage struct {
	ID         string        `json:"id"`
	Channel    ChatChannelType `json:"channel"`
	SenderID   string        `json:"sender_id"`
	SenderName string        `json:"sender_name"`
	Content    string        `json:"content"`
	Timestamp  int64         `json:"timestamp"`
	TargetID   string        `json:"target_id,omitempty"` // 私聊目标
	GuildID    string        `json:"guild_id,omitempty"`  // 公会ID
	IsSystem   bool          `json:"is_system"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// ChatManager 聊天管理器
type ChatManager struct {
	mu          sync.RWMutex
	channels    map[ChatChannelType]*ChatChannel
	privateMsgs map[string][]*ChatMessage // playerID -> messages
	globalMsgs  []*ChatMessage
	maxHistory  int
}

type ChatChannel struct {
	mu      sync.RWMutex
	clients map[string]chan *ChatMessage
	messages []*ChatMessage
	maxSize int
}

// NewChatManager 创建聊天管理器
func NewChatManager() *ChatManager {
	m := &ChatManager{
		channels:    make(map[ChatChannelType]*ChatChannel),
		privateMsgs: make(map[string][]*ChatMessage),
		globalMsgs:  make([]*ChatMessage, 0),
		maxHistory:  1000,
	}

	// 初始化各频道
	for i := ChatChannelType(0); i <= ChatChannelBattle; i++ {
		m.channels[i] = &ChatChannel{
			clients:  make(map[string]chan *ChatMessage),
			messages: make([]*ChatMessage, 0),
			maxSize:  500,
		}
	}

	return m
}

// SendMessage 发送消息
func (m *ChatManager) SendMessage(msg *ChatMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	msg.ID = generateID()
	msg.Timestamp = time.Now().Unix()

	channel, ok := m.channels[msg.Channel]
	if !ok {
		return fmt.Errorf("invalid channel: %d", msg.Channel)
	}

	channel.mu.Lock()
	channel.messages = append(channel.messages, msg)
	if len(channel.messages) > channel.maxSize {
		channel.messages = channel.messages[1:]
	}
	channel.mu.Unlock()

	// 广播给频道内所有玩家
	channel.mu.RLock()
	for _, ch := range channel.clients {
		select {
		case ch <- msg:
		default:
		}
	}
	channel.mu.RUnlock()

	// 私聊单独发送
	if msg.Channel == ChatChannelPrivate && msg.TargetID != "" {
		m.deliverPrivateMsg(msg.TargetID, msg)
	}

	return nil
}

// deliverPrivateMsg 投递私聊消息
func (m *ChatManager) deliverPrivateMsg(playerID string, msg *ChatMessage) {
	m.privateMsgs[playerID] = append(m.privateMsgs[playerID], msg)
}

// Subscribe 订阅频道
func (m *ChatManager) Subscribe(playerID string, channel ChatChannelType) chan *ChatMessage {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch, ok := m.channels[channel]
	if !ok {
		return nil
	}

	msgChan := make(chan *ChatMessage, 50)
	ch.mu.Lock()
	ch.clients[playerID] = msgChan
	ch.mu.Unlock()

	return msgChan
}

// Unsubscribe 取消订阅
func (m *ChatManager) Unsubscribe(playerID string, channel ChatChannelType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch, ok := m.channels[channel]
	if !ok {
		return
	}

	ch.mu.Lock()
	delete(ch.clients, playerID)
	ch.mu.Unlock()
}

// GetHistory 获取历史消息
func (m *ChatManager) GetHistory(channel ChatChannelType, limit int) []*ChatMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ch, ok := m.channels[channel]
	if !ok {
		return nil
	}

	ch.mu.RLock()
	defer ch.mu.RUnlock()

	if limit > len(ch.messages) {
		limit = len(ch.messages)
	}

	start := len(ch.messages) - limit
	if start < 0 {
		start = 0
	}

	result := make([]*ChatMessage, limit)
	copy(result, ch.messages[start:])
	return result
}

// GetPrivateHistory 获取私聊历史
func (m *ChatManager) GetPrivateHistory(playerID string, otherID string, limit int) []*ChatMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()

	msgs := m.privateMsgs[playerID]
	result := make([]*ChatMessage, 0)

	for i := len(msgs) - 1; i >= 0 && len(result) < limit; i-- {
		msg := msgs[i]
		if msg.SenderID == otherID || msg.TargetID == otherID {
			result = append(result, msg)
		}
	}

	return result
}

// ==================== 邮件系统 ====================

// MailType 邮件类型
type MailType int

const (
	MailTypeSystem MailType = iota // 系统邮件
	MailTypePlayer                 // 玩家邮件
	MailTypeGift                   // 礼物邮件
	MailTypeAuction                // 拍卖邮件
)

// Mail 邮件
type Mail struct {
	ID         string        `json:"id"`
	Type       MailType      `json:"type"`
	Title      string        `json:"title"`
	Content    string        `json:"content"`
	SenderID   string        `json:"sender_id"`
	SenderName string        `json:"sender_name"`
	ReceiverID string        `json:"receiver_id"`
	Attachments []*MailAttachment `json:"attachments"`
	Read       bool          `json:"read"`
	Claimed    bool          `json:"claimed"`
	CreatedAt  int64         `json:"created_at"`
	ExpireAt   int64         `json:"expire_at"`
}

// MailAttachment 邮件附件
type MailAttachment struct {
	Type    string `json:"type"` // coin, gem, item
	ItemID  string `json:"item_id,omitempty"`
	Count   int    `json:"count"`
}

// MailManager 邮件管理器
type MailManager struct {
	mu      sync.RWMutex
	mails   map[string]map[string]*Mail // playerID -> mailID -> Mail
	expires time.Duration
}

// NewMailManager 创建邮件管理器
func NewMailManager() *MailManager {
	return &MailManager{
		mails:   make(map[string]map[string]*Mail),
		expires: 30 * 24 * time.Hour, // 30天过期
	}
}

// SendMail 发送邮件
func (m *MailManager) SendMail(mail *Mail) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	mail.ID = generateID()
	mail.CreatedAt = time.Now().Unix()
	mail.ExpireAt = time.Now().Add(m.expires).Unix()

	if m.mails[mail.ReceiverID] == nil {
		m.mails[mail.ReceiverID] = make(map[string]*Mail)
	}

	m.mails[mail.ReceiverID][mail.ID] = mail
	return nil
}

// SendSystemMail 发送系统邮件
func (m *MailManager) SendSystemMail(receiverID, title, content string, attachments []*MailAttachment) error {
	mail := &Mail{
		Type:       MailTypeSystem,
		Title:      title,
		Content:    content,
		SenderID:   "system",
		SenderName: "系统",
		ReceiverID: receiverID,
		Attachments: attachments,
	}
	return m.SendMail(mail)
}

// SendPlayerMail 发送玩家邮件
func (m *MailManager) SendPlayerMail(senderID, senderName, receiverID, title, content string, attachments []*MailAttachment) error {
	mail := &Mail{
		Type:       MailTypePlayer,
		Title:      title,
		Content:    content,
		SenderID:   senderID,
		SenderName: senderName,
		ReceiverID: receiverID,
		Attachments: attachments,
	}
	return m.SendMail(mail)
}

// GetMails 获取玩家邮件列表
func (m *MailManager) GetMails(playerID string) []*Mail {
	m.mu.RLock()
	defer m.mu.RUnlock()

	playerMails := m.mails[playerID]
	result := make([]*Mail, 0, len(playerMails))

	for _, mail := range playerMails {
		result = append(result, mail)
	}

	// 按时间倒序
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].CreatedAt > result[i].CreatedAt {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// GetMail 获取单封邮件
func (m *MailManager) GetMail(playerID, mailID string) *Mail {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.mails[playerID] == nil {
		return nil
	}
	return m.mails[playerID][mailID]
}

// ReadMail 读取邮件
func (m *MailManager) ReadMail(playerID, mailID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	mail := m.mails[playerID][mailID]
	if mail == nil {
		return fmt.Errorf("mail not found")
	}

	mail.Read = true
	return nil
}

// ClaimAttachments 领取附件
func (m *MailManager) ClaimAttachments(playerID, mailID string) ([]*MailAttachment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mail := m.mails[playerID][mailID]
	if mail == nil {
		return nil, fmt.Errorf("mail not found")
	}

	if mail.Claimed {
		return nil, fmt.Errorf("already claimed")
	}

	if len(mail.Attachments) == 0 {
		return nil, fmt.Errorf("no attachments")
	}

	mail.Claimed = true
	return mail.Attachments, nil
}

// DeleteMail 删除邮件
func (m *MailManager) DeleteMail(playerID, mailID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.mails[playerID] == nil {
		return fmt.Errorf("mail not found")
	}

	delete(m.mails[playerID], mailID)
	return nil
}

// CleanExpired 清理过期邮件
func (m *MailManager) CleanExpired() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().Unix()
	count := 0

	for playerID := range m.mails {
		for mailID, mail := range m.mails[playerID] {
			if mail.ExpireAt < now {
				delete(m.mails[playerID], mailID)
				count++
			}
		}
	}

	return count
}

// GetUnreadCount 获取未读邮件数
func (m *MailManager) GetUnreadCount(playerID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, mail := range m.mails[playerID] {
		if !mail.Read {
			count++
		}
	}
	return count
}

// ==================== 好友系统 ====================

// FriendStatus 好友状态
type FriendStatus int

const (
	FriendStatusNone    FriendStatus = iota // 非好友
	FriendStatusPending                     // 待确认
	FriendStatusFriends                     // 好友
	FriendStatusBlocked                     // 黑名单
)

// Friend 好友关系
type Friend struct {
	PlayerID   string       `json:"player_id"`
	FriendID   string       `json:"friend_id"`
	Status     FriendStatus `json:"status"`
	CreatedAt  int64        `json:"created_at"`
	Remark     string       `json:"remark,omitempty"` // 备注
	FriendName string       `json:"friend_name"`
	Level      int          `json:"level"`
	Avatar     string       `json:"avatar"`
}

// FriendManager 好友管理器
type FriendManager struct {
	mu         sync.RWMutex
	friends    map[string]map[string]*Friend // playerID -> friendID -> Friend
	blacklist  map[string]map[string]*Friend
	maxFriends int
	maxBlacklist int
}

// NewFriendManager 创建好友管理器
func NewFriendManager() *FriendManager {
	return &FriendManager{
		friends:    make(map[string]map[string]*Friend),
		blacklist:  make(map[string]map[string]*Friend),
		maxFriends: 500,
		maxBlacklist: 100,
	}
}

// AddFriend 添加好友
func (m *FriendManager) AddFriend(playerID, friendID, friendName string, level int, avatar string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if playerID == friendID {
		return fmt.Errorf("cannot add self")
	}

	// 检查是否被拉黑
	if m.isBlockedLocked(friendID, playerID) {
		return fmt.Errorf("player blocked you")
	}

	// 初始化玩家好友列表
	if m.friends[playerID] == nil {
		m.friends[playerID] = make(map[string]*Friend)
	}
	if m.friends[friendID] == nil {
		m.friends[friendID] = make(map[string]*Friend)
	}

	// 检查是否已达上限
	if len(m.friends[playerID]) >= m.maxFriends {
		return fmt.Errorf("friend list full")
	}

	// 双向好友
	friend := &Friend{
		PlayerID:   playerID,
		FriendID:   friendID,
		Status:     FriendStatusFriends,
		CreatedAt:  time.Now().Unix(),
		FriendName: friendName,
		Level:      level,
		Avatar:     avatar,
	}

	friendReverse := &Friend{
		PlayerID:   friendID,
		FriendID:   playerID,
		Status:     FriendStatusFriends,
		CreatedAt:  time.Now().Unix(),
		FriendName: "", // 对方信息
		Level:      0,
		Avatar:     "",
	}

	m.friends[playerID][friendID] = friend
	m.friends[friendID][playerID] = friendReverse

	return nil
}

// RemoveFriend 删除好友
func (m *FriendManager) RemoveFriend(playerID, friendID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.friends[playerID] == nil {
		return fmt.Errorf("not friends")
	}

	delete(m.friends[playerID], friendID)
	if m.friends[friendID] != nil {
		delete(m.friends[friendID], playerID)
	}

	return nil
}

// BlockPlayer 拉黑玩家
func (m *FriendManager) BlockPlayer(playerID, blockedID, remark string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if playerID == blockedID {
		return fmt.Errorf("cannot block self")
	}

	// 初始化黑名单
	if m.blacklist[playerID] == nil {
		m.blacklist[playerID] = make(map[string]*Friend)
	}

	if len(m.blacklist[playerID]) >= m.maxBlacklist {
		return fmt.Errorf("blacklist full")
	}

	// 先删除好友关系
	if m.friends[playerID] != nil {
		delete(m.friends[playerID], blockedID)
	}

	// 添加黑名单
	block := &Friend{
		PlayerID:  playerID,
		FriendID:  blockedID,
		Status:    FriendStatusBlocked,
		CreatedAt: time.Now().Unix(),
		Remark:    remark,
	}

	m.blacklist[playerID][blockedID] = block
	return nil
}

// UnblockPlayer 解除拉黑
func (m *FriendManager) UnblockPlayer(playerID, blockedID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.blacklist[playerID] == nil {
		return fmt.Errorf("not in blacklist")
	}

	delete(m.blacklist[playerID], blockedID)
	return nil
}

// isBlockedLocked 检查是否被拉黑（需持有锁）
func (m *FriendManager) isBlockedLocked(playerID, blockedID string) bool {
	if m.blacklist[playerID] == nil {
		return false
	}
	_, blocked := m.blacklist[playerID][blockedID]
	return blocked
}

// IsBlocked 检查是否被拉黑
func (m *FriendManager) IsBlocked(playerID, blockedID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isBlockedLocked(playerID, blockedID)
}

// GetFriends 获取好友列表
func (m *FriendManager) GetFriends(playerID string) []*Friend {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.friends[playerID] == nil {
		return nil
	}

	result := make([]*Friend, 0)
	for _, friend := range m.friends[playerID] {
		if friend.Status == FriendStatusFriends {
			result = append(result, friend)
		}
	}

	return result
}

// GetBlacklist 获取黑名单
func (m *FriendManager) GetBlacklist(playerID string) []*Friend {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.blacklist[playerID] == nil {
		return nil
	}

	result := make([]*Friend, 0)
	for _, friend := range m.blacklist[playerID] {
		result = append(result, friend)
	}

	return result
}

// GetFriendCount 获取好友数量
func (m *FriendManager) GetFriendCount(playerID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.friends[playerID] == nil {
		return 0
	}

	count := 0
	for _, friend := range m.friends[playerID] {
		if friend.Status == FriendStatusFriends {
			count++
		}
	}
	return count
}

// SetRemark 设置备注
func (m *FriendManager) SetRemark(playerID, friendID, remark string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.friends[playerID] == nil {
		return fmt.Errorf("not friends")
	}

	friend := m.friends[playerID][friendID]
	if friend == nil {
		return fmt.Errorf("not friends")
	}

	friend.Remark = remark
	return nil
}

// ==================== 商店系统 ====================

// ShopType 商店类型
type ShopType int

const (
	ShopTypeGold    ShopType = iota // 金币商店
	ShopTypeGem                      // 钻石商店
	ShopTypeItem                     // 道具商店
	ShopTypeSkill                    // 技能商店
	ShopTypeAvatar                   // 头像商店
	ShopTypeFrame                    // 头像框商店
	ShopTypeTitle                    // 称号商店
	ShopTypeGift                     // 礼包商店
	ShopTypeLimited                  // 限时商店
)

// ShopItem 商品
type ShopItem struct {
	ID          string            `json:"id"`
	Type        ShopType         `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       int               `json:"price"`
	PriceType   string            `json:"price_type"` // coin, gem
	Stock       int               `json:"stock"`       // -1 表示无限
	Sold        int               `json:"sold"`
	Discount    float64           `json:"discount"`    // 折扣 0.0-1.0
	LevelRequire int               `json:"level_require"`
	Items       []*ShopItemReward  `json:"items"`       // 奖励道具
	StartTime   int64             `json:"start_time"`
	EndTime     int64             `json:"end_time"`
	Tags        []string          `json:"tags"`
}

// ShopItemReward 商店道具奖励
type ShopItemReward struct {
	ItemID string `json:"item_id"`
	Count  int    `json:"count"`
}

// ShopRecord 购买记录
type ShopRecord struct {
	ID        string `json:"id"`
	PlayerID  string `json:"player_id"`
	ItemID    string `json:"item_id"`
	Price     int    `json:"price"`
	PriceType string `json:"price_type"`
	Count     int    `json:"count"`
	Timestamp int64  `json:"timestamp"`
}

// ShopManager 商店管理器
type ShopManager struct {
	mu        sync.RWMutex
	shops     map[ShopType]*Shop
	records   map[string][]*ShopRecord // playerID -> records
	dailyLimit map[string]map[string]int // playerID -> itemID -> count
	resetTime time.Time
}

// Shop 商店
type Shop struct {
	Type     ShopType   `json:"type"`
	Name     string     `json:"name"`
	Items    []*ShopItem `json:"items"`
	OpenHours string    `json:"open_hours"` // e.g., "00:00-23:59"
	RefreshCD int       `json:"refresh_cd"` // 刷新冷却(秒)
}

// NewShopManager 创建商店管理器
func NewShopManager() *ShopManager {
	m := &ShopManager{
		shops:     make(map[ShopType]*Shop),
		records:   make(map[string][]*ShopRecord),
		dailyLimit: make(map[string]map[string]int),
		resetTime: time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour),
	}

	// 初始化各商店
	m.initShops()

	return m
}

// initShops 初始化商店
func (m *ShopManager) initShops() {
	// 金币商店
	m.shops[ShopTypeGold] = &Shop{
		Type:  ShopTypeGold,
		Name:  "金币商店",
		Items: []*ShopItem{
			{ID: "gold_100", Name: "100金币", Price: 10, PriceType: "gem", Items: []*ShopItemReward{{ItemID: "gold", Count: 100}}},
			{ID: "gold_500", Name: "500金币", Price: 45, PriceType: "gem", Items: []*ShopItemReward{{ItemID: "gold", Count: 500}}},
			{ID: "gold_1000", Name: "1000金币", Price: 80, PriceType: "gem", Items: []*ShopItemReward{{ItemID: "gold", Count: 1000}}},
			{ID: "gold_5000", Name: "5000金币", Price: 350, PriceType: "gem", Items: []*ShopItemReward{{ItemID: "gold", Count: 5000}}},
		},
		RefreshCD: 0,
	}

	// 钻石商店
	m.shops[ShopTypeGem] = &Shop{
		Type:  ShopTypeGem,
		Name:  "钻石商店",
		Items: []*ShopItem{
			{ID: "gem_60", Name: "60钻石", Price: 60, PriceType: "rmb", Items: []*ShopItemReward{{ItemID: "gem", Count: 60}}},
			{ID: "gem_300", Name: "300钻石", Price: 280, PriceType: "rmb", Items: []*ShopItemReward{{ItemID: "gem", Count: 300}}, Discount: 0.9},
			{ID: "gem_980", Name: "980钻石", Price: 980, PriceType: "rmb", Items: []*ShopItemReward{{ItemID: "gem", Count: 980}}, Discount: 0.85},
			{ID: "gem_1980", Name: "1980钻石", Price: 1980, PriceType: "rmb", Items: []*ShopItemReward{{ItemID: "gem", Count: 1980}}, Discount: 0.8},
			{ID: "gem_3280", Name: "3280钻石", Price: 3280, PriceType: "rmb", Items: []*ShopItemReward{{ItemID: "gem", Count: 3280}}, Discount: 0.75},
		},
		RefreshCD: 0,
	}

	// 道具商店
	m.shops[ShopTypeItem] = &Shop{
		Type:  ShopTypeItem,
		Name:  "道具商店",
		Items: []*ShopItem{
			{ID: "hp_potion", Name: "生命药水", Price: 50, PriceType: "coin", Items: []*ShopItemReward{{ItemID: "hp_potion", Count: 1}}},
			{ID: "mp_potion", Name: "魔法药水", Price: 50, PriceType: "coin", Items: []*ShopItemReward{{ItemID: "mp_potion", Count: 1}}},
			{ID: "revive", Name: "复活十字", Price: 100, PriceType: "gem", Items: []*ShopItemReward{{ItemID: "revive", Count: 1}}},
			{ID: "exp_book", Name: "经验书", Price: 200, PriceType: "gem", Items: []*ShopItemReward{{ItemID: "exp_book", Count: 1}}},
		},
		RefreshCD: 3600,
	}

	// 限时商店
	m.shops[ShopTypeLimited] = &Shop{
		Type:  ShopTypeLimited,
		Name:  "限时特惠",
		Items: []*ShopItem{},
		RefreshCD: 7200,
	}

	// 礼包商店
	m.shops[ShopTypeGift] = &Shop{
		Type:  ShopTypeGift,
		Name:  "礼包商店",
		Items: []*ShopItem{
			{ID: "starter_pack", Name: "新手礼包", Price: 60, PriceType: "rmb", Items: []*ShopItemReward{
				{ItemID: "gem", Count: 100},
				{ItemID: "gold", Count: 1000},
				{ItemID: "weapon_1", Count: 1},
			}},
			{ID: "daily_pack", Name: "每日礼包", Price: 30, PriceType: "rmb", Items: []*ShopItemReward{
				{ItemID: "gem", Count: 50},
				{ItemID: "gold", Count: 500},
			}},
			{ID: "month_card", Name: "月卡", Price: 300, PriceType: "rmb", Tags: []string{"subscription"}, Items: []*ShopItemReward{
				{ItemID: "gem", Count: 100},
			}},
		},
		RefreshCD: 86400,
	}
}

// GetShop 获取商店
func (m *ShopManager) GetShop(shopType ShopType) *Shop {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.shops[shopType]
}

// GetAllShops 获取所有商店
func (m *ShopManager) GetAllShops() []*Shop {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Shop, 0, len(m.shops))
	for _, shop := range m.shops {
		result = append(result, shop)
	}
	return result
}

// Buy 购买商品
func (m *ShopManager) Buy(playerID, itemID string, count int) ([]*ShopItemReward, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 查找商品
	var shopItem *ShopItem
	var shop *Shop
	for _, s := range m.shops {
		for _, item := range s.Items {
			if item.ID == itemID {
				shopItem = item
				shop = s
				break
			}
		}
		if shopItem != nil {
			break
		}
	}

	if shopItem == nil {
		return nil, fmt.Errorf("item not found")
	}

	// 检查库存
	if shopItem.Stock >= 0 && shopItem.Sold+count > shopItem.Stock {
		return nil, fmt.Errorf("out of stock")
	}

	// 检查等级
	if shopItem.LevelRequire > 0 {
		// 需要检查玩家等级
		_ = playerID // TODO: check player level
	}

	// 检查每日限购
	if m.dailyLimit[playerID] == nil {
		m.dailyLimit[playerID] = make(map[string]int)
	}

	// 计算价格
	price := shopItem.Price
	if shopItem.Discount > 0 {
		price = int(float64(price) * shopItem.Discount)
	}
	totalPrice := price * count

	// 检查余额（需要调用外部接口）
	// 扣除货币（需要调用外部接口）

	// 更新销量
	shopItem.Sold += count

	// 记录购买
	record := &ShopRecord{
		ID:        generateID(),
		PlayerID:  playerID,
		ItemID:    itemID,
		Price:     totalPrice,
		PriceType: shopItem.PriceType,
		Count:     count,
		Timestamp: time.Now().Unix(),
	}
	m.records[playerID] = append(m.records[playerID], record)

	// 更新每日限购
	m.dailyLimit[playerID][itemID] += count

	// 返回奖励
	rewards := make([]*ShopItemReward, len(shopItem.Items))
	for i, item := range shopItem.Items {
		rewards[i] = &ShopItemReward{
			ItemID: item.ItemID,
			Count:  item.Count * count,
		}
	}

	return rewards, nil
}

// GetPurchaseHistory 获取购买历史
func (m *ShopManager) GetPurchaseHistory(playerID string) []*ShopRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.records[playerID]
}

// RefreshShop 刷新商店
func (m *ShopManager) RefreshShop(shopType ShopType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	shop, ok := m.shops[shopType]
	if !ok {
		return fmt.Errorf("shop not found")
	}

	// 重置销量
	for _, item := range shop.Items {
		item.Sold = 0
	}

	return nil
}

// AddLimitedItem 添加限时商品
func (m *ShopManager) AddLimitedItem(item *ShopItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	shop, ok := m.shops[ShopTypeLimited]
	if !ok {
		return fmt.Errorf("limited shop not found")
	}

	shop.Items = append(shop.Items, item)
	return nil
}

// CleanExpiredItems 清理过期商品
func (m *ShopManager) CleanExpiredItems() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().Unix()
	count := 0

	shop := m.shops[ShopTypeLimited]
	if shop == nil {
		return 0
	}

	newItems := make([]*ShopItem, 0)
	for _, item := range shop.Items {
		if item.EndTime > 0 && item.EndTime < now {
			count++
			continue
		}
		newItems = append(newItems, item)
	}

	shop.Items = newItems
	return count
}

// ==================== 辅助函数 ====================

func generateID() string {
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), time.Now().Nanosecond()%1000)
}
