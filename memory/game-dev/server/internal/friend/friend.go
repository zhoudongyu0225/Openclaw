package friend

import (
	"fmt"
	"sync"
	"time"
)

// FriendStatus 好友状态
type FriendStatus int

const (
	FriendStatusNone      FriendStatus = iota // 非好友
	FriendStatusPending                        // 待确认
	FriendStatusFriends                        // 好友
	FriendStatusBlocked                        // 已拉黑
)

// Friend 好友结构
type Friend struct {
	UserID      string    `json:"user_id"`       // 玩家ID
	Nickname    string    `json:"nickname"`       // 昵称
	Avatar      string    `json:"avatar"`         // 头像
	Level       int       `json:"level"`          // 等级
	Status      FriendStatus `json:"status"`     // 好友状态
	CreatedAt   time.Time `json:"created_at"`    // 添加时间
	LastOnline  time.Time `json:"last_online"`   // 最后在线时间
	IsOnline    bool      `json:"is_online"`     // 是否在线
	FriendPoints int      `json:"friend_points"` // 亲密度
}

// FriendRequest 好友请求
type FriendRequest struct {
	ID         string    `json:"id"`           // 请求ID
	FromUserID string    `json:"from_user_id"` // 申请人ID
	ToUserID   string    `json:"to_user_id"`   // 被申请人ID
	Message    string    `json:"message"`       // 附加消息
	Status     string    `json:"status"`        // pending/accepted/rejected
	CreatedAt  time.Time `json:"created_at"`    // 创建时间
}

// FriendManager 好友管理器
type FriendManager struct {
	mu sync.RWMutex

	// 好友关系: userID -> friendUserID -> Friend
	friends map[string]map[string]*Friend

	// 好友请求: userID -> []FriendRequest (收到的请求)
	requests map[string][]*FriendRequest

	// 黑名单: userID -> blockedUserID -> Friend
	blacklist map[string]map[string]*Friend

	// 亲密度: userID -> friendUserID -> points
	Intimacy map[string]map[string]int
}

// NewFriendManager 创建好友管理器
func NewFriendManager() *FriendManager {
	return &FriendManager{
		friends:    make(map[string]map[string]*Friend),
		requests:   make(map[string][]*FriendRequest),
		blacklist:  make(map[string]map[string]*Friend),
		Intimacy:   make(map[string]map[string]int),
	}
}

// SendFriendRequest 发送好友请求
func (fm *FriendManager) SendFriendRequest(fromUserID, toUserID, message string) (*FriendRequest, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// 检查是否已在黑名单
	if blacklist, ok := fm.blacklist[toUserID]; ok {
		if _, blocked := blacklist[fromUserID]; blocked {
			return nil, fmt.Errorf("对方已将你拉黑")
		}
	}

	// 检查是否已经是好友
	if friends, ok := fm.friends[fromUserID]; ok {
		if _, exists := friends[toUserID]; exists {
			return nil, fmt.Errorf("你们已经是好友了")
		}
	}

	// 检查是否已存在待处理的请求
	for _, req := range fm.requests[toUserID] {
		if req.FromUserID == fromUserID && req.Status == "pending" {
			return nil, fmt.Errorf("已发送过好友请求")
		}
	}

	req := &FriendRequest{
		ID:         fmt.Sprintf("req_%d", time.Now().UnixNano()),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	fm.requests[toUserID] = append(fm.requests[toUserID], req)
	return req, nil
}

// AcceptFriendRequest 接受好友请求
func (fm *FriendManager) AcceptFriendRequest(requestID, userID string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// 找到请求
	var req *FriendRequest
	var reqIndex int
	for i, r := range fm.requests[userID] {
		if r.ID == requestID {
			req = r
			reqIndex = i
			break
		}
	}

	if req == nil {
		return fmt.Errorf("请求不存在")
	}

	if req.Status != "pending" {
		return fmt.Errorf("请求已处理")
	}

	// 添加好友关系
	fm.addFriendPair(req.FromUserID, req.ToUserID)

	// 更新请求状态
	req.Status = "accepted"

	// 移除请求
	fm.requests[userID] = append(fm.requests[userID][:reqIndex], fm.requests[userID][reqIndex+1:]...)

	return nil
}

// RejectFriendRequest 拒绝好友请求
func (fm *FriendManager) RejectFriendRequest(requestID, userID string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for i, req := range fm.requests[userID] {
		if req.ID == requestID {
			req.Status = "rejected"
			fm.requests[userID] = append(fm.requests[userID][:i], fm.requests[userID][i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("请求不存在")
}

// RemoveFriend 删除好友
func (fm *FriendManager) RemoveFriend(userID, friendID string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if _, ok := fm.friends[userID]; !ok {
		return fmt.Errorf("你们不是好友")
	}

	delete(fm.friends[userID], friendID)
	delete(fm.friends[friendID], userID)

	// 清除亲密度
	delete(fm.Intimacy[userID], friendID)
	delete(fm.Intimacy[friendID], userID)

	return nil
}

// BlockUser 拉黑用户
func (fm *FriendManager) BlockUser(userID, blockedID string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// 先删除好友关系
	if friends, ok := fm.friends[userID]; ok {
		delete(friends, blockedID)
	}

	// 添加到黑名单
	if fm.blacklist[userID] == nil {
		fm.blacklist[userID] = make(map[string]*Friend)
	}

	fm.blacklist[userID][blockedID] = &Friend{
		UserID:     blockedID,
		CreatedAt:  time.Now(),
		Status:     FriendStatusBlocked,
	}

	// 清除亲密度
	if intimacy, ok := fm.Intimacy[userID]; ok {
		delete(intimacy, blockedID)
	}

	return nil
}

// UnblockUser 解除拉黑
func (fm *FriendManager) UnblockUser(userID, blockedID string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if _, ok := fm.blacklist[userID]; !ok {
		return fmt.Errorf("用户不在黑名单中")
	}

	delete(fm.blacklist[userID], blockedID)
	return nil
}

// GetFriends 获取好友列表
func (fm *FriendManager) GetFriends(userID string) []*Friend {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	var result []*Friend
	if friends, ok := fm.friends[userID]; ok {
		for _, f := range friends {
			result = append(result, f)
		}
	}

	return result
}

// GetFriendRequests 获取好友请求列表
func (fm *FriendManager) GetFriendRequests(userID string) []*FriendRequest {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	return fm.requests[userID]
}

// GetBlacklist 获取黑名单
func (fm *FriendManager) GetBlacklist(userID string) []*Friend {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	var result []*Friend
	if blacklist, ok := fm.blacklist[userID]; ok {
		for _, f := range blacklist {
			result = append(result, f)
		}
	}

	return result
}

// UpdateFriendPoints 更新亲密度
func (fm *FriendManager) UpdateFriendPoints(userID, friendID string, points int) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.Intimacy[userID] == nil {
		fm.Intimacy[userID] = make(map[string]int)
	}

	fm.Intimacy[userID][friendID] += points

	// 亲密度上限
	if fm.Intimacy[userID][friendID] > 10000 {
		fm.Intimacy[userID][friendID] = 10000
	}
}

// GetIntimacy 获取亲密度
func (fm *FriendManager) GetIntimacy(userID, friendID string) int {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if intimacy, ok := fm.Intimacy[userID]; ok {
		return intimacy[friendID]
	}

	return 0
}

// addFriendPair 添加双向好友关系
func (fm *FriendManager) addFriendPair(userID1, userID2 string) {
	// 用户1的好友列表
	if fm.friends[userID1] == nil {
		fm.friends[userID1] = make(map[string]*Friend)
	}
	fm.friends[userID1][userID2] = &Friend{
		UserID:      userID2,
		CreatedAt:  time.Now(),
		Status:     FriendStatusFriends,
		FriendPoints: 0,
	}

	// 用户2的好友列表
	if fm.friends[userID2] == nil {
		fm.friends[userID2] = make(map[string]*Friend)
	}
	fm.friends[userID2][userID1] = &Friend{
		UserID:      userID1,
		CreatedAt:  time.Now(),
		Status:     FriendStatusFriends,
		FriendPoints: 0,
	}

	// 初始化亲密度
	if fm.Intimacy[userID1] == nil {
		fm.Intimacy[userID1] = make(map[string]int)
	}
	if fm.Intimacy[userID2] == nil {
		fm.Intimacy[userID2] = make(map[string]int)
	}
	fm.Intimacy[userID1][userID2] = 0
	fm.Intimacy[userID2][userID1] = 0
}

// IsFriend 检查是否为好友
func (fm *FriendManager) IsFriend(userID1, userID2 string) bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if friends, ok := fm.friends[userID1]; ok {
		if _, exists := friends[userID2]; exists {
			return true
		}
	}

	return false
}

// GetFriendCount 获取好友数量
func (fm *FriendManager) GetFriendCount(userID string) int {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	if friends, ok := fm.friends[userID]; ok {
		return len(friends)
	}

	return 0
}
