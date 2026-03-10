package chat

import (
	"fmt"
	"sync"
	"time"
)

// ChannelType 频道类型
type ChannelType int

const (
	ChannelWorld   ChannelType = iota // 世界频道
	ChannelRoom                        // 房间频道
	ChannelGuild                       // 公会频道
	ChannelPrivate                     // 私聊频道
	ChannelSystem                      // 系统频道
	ChannelBroadcast                   // 广播频道
)

// MessageType 消息类型
type MessageType int

const (
	MsgTypeText    MessageType = iota // 文本
	MsgTypeImage                      // 图片
	MsgTypeVoice                      // 语音
	MsgTypeSystem                     // 系统消息
	MsgTypeGift                       // 礼物消息
	MsgTypeEmote                      // 表情消息
)

// ChatMessage 聊天消息
type ChatMessage struct {
	ID         string      `json:"id"`          // 消息ID
	ChannelID string      `json:"channel_id"`  // 频道ID
	SenderID   string      `json:"sender_id"`  // 发送者ID
	SenderName string      `json:"sender_name"`// 发送者昵称
	SenderAvatar string    `json:"sender_avatar"` // 发送者头像
	Content    string      `json:"content"`    // 消息内容
	MessageType MessageType `json:"message_type"` // 消息类型
	Timestamp  time.Time   `json:"timestamp"` // 时间戳
	ChannelType ChannelType `json:"channel_type"` // 频道类型
	Extra      string      `json:"extra"`      // 扩展数据 (礼物信息等)
	IsDeleted  bool        `json:"is_deleted"` // 是否被删除
}

// ChatChannel 聊天频道
type ChatChannel struct {
	ID          string      `json:"id"`           // 频道ID
	Name        string      `json:"name"`         // 频道名称
	ChannelType ChannelType `json:"channel_type"` // 频道类型
	Members     []string    `json:"members"`      // 成员列表
	MaxMembers  int         `json:"max_members"`  // 最大成员数
	CreatedAt   time.Time   `json:"created_at"`  // 创建时间
	OwnerID     string      `json:"owner_id"`    // 所有者ID
	IsMuted     bool        `json:"is_muted"`    // 是否禁言
	MuteEndTime time.Time   `json:"mute_end_time"`// 禁言结束时间
}

// ChatUser 用户聊天状态
type ChatUser struct {
	UserID       string    `json:"user_id"`        // 用户ID
	JoinChannels []string  `json:"join_channels"`  // 加入的频道
	MuteUntil    time.Time `json:"mute_until"`    // 禁言截止时间
	ChatLevel    int       `json:"chat_level"`     // 聊天等级
	TotalMsgs    int       `json:"total_msgs"`    // 总消息数
}

// ChatManager 聊天管理器
type ChatManager struct {
	mu sync.RWMutex

	// 频道: channelID -> ChatChannel
	channels map[string]*ChatChannel

	// 频道消息: channelID -> []ChatMessage
	messages map[string][]*ChatMessage

	// 用户聊天状态: userID -> ChatUser
	users map[string]*ChatUser

	// 私聊会话: userID1_userID2 -> channelID
	privateChats map[string]string

	// 消息ID生成
	msgIDCounter int64
}

// NewChatManager 创建聊天管理器
func NewChatManager() *ChatManager {
	return &ChatManager{
		channels:     make(map[string]*ChatChannel),
		messages:     make(map[string][]*ChatMessage),
		users:        make(map[string]*ChatUser),
		privateChats: make(map[string]string),
		msgIDCounter: time.Now().UnixNano(),
	}
}

// CreateChannel 创建频道
func (cm *ChatManager) CreateChannel(name string, channelType ChannelType, ownerID string) (*ChatChannel, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	channelID := fmt.Sprintf("ch_%d", time.Now().UnixNano())

	channel := &ChatChannel{
		ID:          channelID,
		Name:        name,
		ChannelType: channelType,
		Members:     []string{ownerID},
		MaxMembers:  100, // 默认最大成员数
		CreatedAt:   time.Now(),
		OwnerID:     ownerID,
		IsMuted:     false,
	}

	cm.channels[channelID] = channel
	cm.messages[channelID] = make([]*ChatMessage, 0)

	return channel, nil
}

// JoinChannel 加入频道
func (cm *ChatManager) JoinChannel(userID, channelID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	channel, ok := cm.channels[channelID]
	if !ok {
		return fmt.Errorf("频道不存在")
	}

	// 检查是否已加入
	for _, member := range channel.Members {
		if member == userID {
			return fmt.Errorf("已加入该频道")
		}
	}

	// 检查成员数限制
	if len(channel.Members) >= channel.MaxMembers {
		return fmt.Errorf("频道成员已满")
	}

	channel.Members = append(channel.Members, userID)

	// 更新用户状态
	if cm.users[userID] == nil {
		cm.users[userID] = &ChatUser{
			UserID:       userID,
			JoinChannels: []string{},
			ChatLevel:    1,
		}
	}
	cm.users[userID].JoinChannels = append(cm.users[userID].JoinChannels, channelID)

	return nil
}

// LeaveChannel 离开频道
func (cm *ChatManager) LeaveChannel(userID, channelID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	channel, ok := cm.channels[channelID]
	if !ok {
		return fmt.Errorf("频道不存在")
	}

	// 移除成员
	for i, member := range channel.Members {
		if member == userID {
			channel.Members = append(channel.Members[:i], channel.Members[i+1:]...)
			break
		}
	}

	// 更新用户状态
	if cm.users[userID] != nil {
		channels := cm.users[userID].JoinChannels
		for i, ch := range channels {
			if ch == channelID {
				cm.users[userID].JoinChannels = append(channels[:i], channels[i+1:]...)
				break
			}
		}
	}

	return nil
}

// SendMessage 发送消息
func (cm *ChatManager) SendMessage(senderID, senderName, senderAvatar, channelID, content string, msgType MessageType) (*ChatMessage, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	channel, ok := cm.channels[channelID]
	if !ok {
		return nil, fmt.Errorf("频道不存在")
	}

	// 检查用户是否被禁言
	if cm.users[senderID] != nil {
		if time.Now().Before(cm.users[senderID].MuteUntil) {
			return nil, fmt.Errorf("你已被禁言")
		}
	}

	// 检查频道是否被禁言
	if channel.IsMuted && time.Now().Before(channel.MuteEndTime) {
		return nil, fmt.Errorf("频道已被禁言")
	}

	// 生成消息ID
	cm.msgIDCounter++
	msgID := fmt.Sprintf("msg_%d", cm.msgIDCounter)

	msg := &ChatMessage{
		ID:           msgID,
		ChannelID:    channelID,
		SenderID:     senderID,
		SenderName:   senderName,
		SenderAvatar: senderAvatar,
		Content:      content,
		MessageType:  msgType,
		Timestamp:    time.Now(),
		ChannelType:  channel.ChannelType,
	}

	// 存储消息
	cm.messages[channelID] = append(cm.messages[channelID], msg)

	// 限制频道消息数量 (保留最近1000条)
	if len(cm.messages[channelID]) > 1000 {
		cm.messages[channelID] = cm.messages[channelID][-1000:]
	}

	// 更新用户消息数
	if cm.users[senderID] != nil {
		cm.users[senderID].TotalMsgs++
	}

	return msg, nil
}

// GetChannelMessages 获取频道消息
func (cm *ChatManager) GetChannelMessages(channelID string, limit int) []*ChatMessage {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	msgs := cm.messages[channelID]
	if len(msgs) == 0 {
		return []*ChatMessage{}
	}

	if limit > len(msgs) {
		limit = len(msgs)
	}

	// 返回最近的消息
	start := len(msgs) - limit
	result := make([]*ChatMessage, limit)
	copy(result, msgs[start:])

	return result
}

// CreatePrivateChat 创建私聊频道
func (cm *ChatManager) CreatePrivateChat(userID1, userID2 string) (string, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 生成固定的私聊频道ID (保证双方看到同一个频道)
	chatID := userID1
	if userID1 > userID2 {
		chatID = userID2 + "_" + userID1
	} else {
		chatID = userID1 + "_" + userID2
	}
	privateChatID := "private_" + chatID

	// 如果已存在，直接返回
	if _, ok := cm.channels[privateChatID]; ok {
		return privateChatID, nil
	}

	// 创建私聊频道
	channel := &ChatChannel{
		ID:          privateChatID,
		Name:        "私聊",
		ChannelType: ChannelPrivate,
		Members:     []string{userID1, userID2},
		MaxMembers:  2,
		CreatedAt:   time.Now(),
		OwnerID:     userID1,
	}

	cm.channels[privateChatID] = channel
	cm.messages[privateChatID] = make([]*ChatMessage, 0)
	cm.privateChats[chatID] = privateChatID

	return privateChatID, nil
}

// SendPrivateMessage 发送私聊消息
func (cm *ChatManager) SendPrivateMessage(senderID, receiverID, content string) (*ChatMessage, error) {
	// 创建或获取私聊频道
	channelID, err := cm.CreatePrivateChat(senderID, receiverID)
	if err != nil {
		return nil, err
	}

	return cm.SendMessage(senderID, senderID, "", channelID, content, MsgTypeText)
}

// MuteUser 禁言用户
func (cm *ChatManager) MuteUser(userID string, duration time.Duration) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.users[userID] == nil {
		cm.users[userID] = &ChatUser{
			UserID:    userID,
			ChatLevel: 1,
		}
	}

	cm.users[userID].MuteUntil = time.Now().Add(duration)
	return nil
}

// UnmuteUser 解除禁言
func (cm *ChatManager) UnmuteUser(userID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.users[userID] == nil {
		return fmt.Errorf("用户不存在")
	}

	cm.users[userID].MuteUntil = time.Time{}
	return nil
}

// MuteChannel 禁言频道
func (cm *ChatManager) MuteChannel(channelID string, duration time.Duration) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	channel, ok := cm.channels[channelID]
	if !ok {
		return fmt.Errorf("频道不存在")
	}

	channel.IsMuted = true
	channel.MuteEndTime = time.Now().Add(duration)

	return nil
}

// DeleteMessage 删除消息
func (cm *ChatManager) DeleteMessage(messageID, channelID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	msgs, ok := cm.messages[channelID]
	if !ok {
		return fmt.Errorf("频道不存在")
	}

	for i, msg := range msgs {
		if msg.ID == messageID {
			msgs[i].IsDeleted = true
			msgs[i].Content = "[消息已被删除]"
			return nil
		}
	}

	return fmt.Errorf("消息不存在")
}

// GetChannel 获取频道信息
func (cm *ChatManager) GetChannel(channelID string) (*ChatChannel, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	channel, ok := cm.channels[channelID]
	if !ok {
		return nil, fmt.Errorf("频道不存在")
	}

	return channel, nil
}

// GetUserChannels 获取用户加入的频道
func (cm *ChatManager) GetUserChannels(userID string) []*ChatChannel {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var result []*ChatChannel

	user := cm.users[userID]
	if user == nil {
		return result
	}

	for _, channelID := range user.JoinChannels {
		if channel, ok := cm.channels[channelID]; ok {
			result = append(result, channel)
		}
	}

	return result
}

// GetOnlineCount 获取在线人数 (世界频道)
func (cm *ChatManager) GetOnlineCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// 统计所有频道的去重用户数
	userCount := make(map[string]bool)
	for _, channel := range cm.channels {
		if channel.ChannelType == ChannelWorld {
			for _, member := range channel.Members {
				userCount[member] = true
			}
		}
	}

	return len(userCount)
}

// CleanupOldMessages 清理过期消息
func (cm *ChatManager) CleanupOldMessages(maxAge time.Duration) int {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	deletedCount := 0

	for channelID, msgs := range cm.messages {
		var remaining []*ChatMessage
		for _, msg := range msgs {
			if msg.Timestamp.After(cutoff) {
				remaining = append(remaining, msg)
			} else {
				deletedCount++
			}
		}
		cm.messages[channelID] = remaining
	}

	return deletedCount
}

// GetChannelMemberCount 获取频道成员数
func (cm *ChatManager) GetChannelMemberCount(channelID string) int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if channel, ok := cm.channels[channelID]; ok {
		return len(channel.Members)
	}

	return 0
}
