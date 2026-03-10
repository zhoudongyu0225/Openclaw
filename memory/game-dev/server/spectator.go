package game

import (
	"fmt"
	"sync"
	"time"
)

// SpectatorState 观战状态
type SpectatorState int

const (
	SpectatorWaiting SpectatorState = iota // 等待观战
	SpectatorWatching                      // 观看中
	SpectatorLeft                          // 已离开
)

// Spectator 观战者
type Spectator struct {
	UserID       string         `json:"user_id"`        // 用户ID
	UserName     string         `json:"user_name"`      // 用户名
	RoomID       string         `json:"room_id"`        // 房间ID
	State        SpectatorState `json:"state"`          // 观战状态
	JoinTime     time.Time      `json:"join_time"`     // 加入时间
	LastUpdateTime time.Time    `json:"last_update_time"` // 最后更新时间
	CameraAngle  int            `json:"camera_angle"`  // 视角 (0=跟随, 1=自由, 2=上帝)
	ShowDanmaku  bool           `json:"show_danmaku"`  // 是否显示弹幕
	ShowStats    bool           `json:"show_stats"`    // 是否显示统计
}

// SpectatorRoom 观战房间
type SpectatorRoom struct {
	RoomID       string       `json:"room_id"`        // 房间ID
	BattleRoomID string       `json:"battle_room_id"` // 对战房间ID
	HostUserID   string       `json:"host_user_id"`   // 房主用户ID
	Spectators   []*Spectator `json:"spectators"`     // 观战者列表
	MaxSpectators int         `json:"max_spectators"` // 最大观战人数
	IsRecording  bool         `json:"is_recording"`   // 是否录制
	ViewCount    int          `json:"view_count"`     // 总观看次数
	CreatedAt    time.Time    `json:"created_at"`     // 创建时间
}

// SpectatorManager 观战管理器
type SpectatorManager struct {
	mu sync.RWMutex

	// 观战房间: roomID -> SpectatorRoom
	spectatorRooms map[string]*SpectatorRoom

	// 观战者: userID -> Spectator
	spectators map[string]*Spectator

	// 用户对应观战房间: userID -> roomID
	userToRoom map[string]string

	// 战斗房间到观战房间: battleRoomID -> spectatorRoomID
	battleToSpectator map[string]string

	// 观战房间ID生成
	roomIDCounter int64
}

// NewSpectatorManager 创建观战管理器
func NewSpectatorManager() *SpectatorManager {
	return &SpectatorManager{
		spectatorRooms:   make(map[string]*SpectatorRoom),
		spectators:       make(map[string]*Spectator),
		userToRoom:       make(map[string]string),
		battleToSpectator: make(map[string]string),
		roomIDCounter:    time.Now().UnixNano(),
	}
}

// CreateSpectatorRoom 创建观战房间
func (sm *SpectatorManager) CreateSpectatorRoom(battleRoomID, hostUserID string) *SpectatorRoom {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.roomIDCounter++
	roomID := fmt.Sprintf("spec_%d", sm.roomIDCounter)

	spectatorRoom := &SpectatorRoom{
		RoomID:        roomID,
		BattleRoomID:  battleRoomID,
		HostUserID:    hostUserID,
		Spectators:    make([]*Spectator, 0),
		MaxSpectators: 100, // 默认最大观战人数
		IsRecording:  false,
		ViewCount:    0,
		CreatedAt:    time.Now(),
	}

	sm.spectatorRooms[roomID] = spectatorRoom
	sm.battleToSpectator[battleRoomID] = roomID

	return spectatorRoom
}

// JoinSpectator 加入观战
func (sm *SpectatorManager) JoinSpectator(userID, userName, battleRoomID string) (*Spectator, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 获取观战房间
	spectatorRoomID, ok := sm.battleToSpectator[battleRoomID]
	if !ok {
		// 自动创建观战房间
		sm.roomIDCounter++
		spectatorRoomID = fmt.Sprintf("spec_%d", sm.roomIDCounter)
		spectatorRoom := &SpectatorRoom{
			RoomID:        spectatorRoomID,
			BattleRoomID:  battleRoomID,
			HostUserID:    "",
			Spectators:    make([]*Spectator, 0),
			MaxSpectators: 100,
			ViewCount:    0,
			CreatedAt:    time.Now(),
		}
		sm.spectatorRooms[spectatorRoomID] = spectatorRoom
		sm.battleToSpectator[battleRoomID] = spectatorRoomID
	}

	spectatorRoom := sm.spectatorRooms[spectatorRoomID]

	// 检查观战人数是否已满
	if len(spectatorRoom.Spectators) >= spectatorRoom.MaxSpectators {
		return nil, fmt.Errorf("观战人数已满")
	}

	// 检查用户是否已在观战中
	if existing, ok := sm.spectators[userID]; ok {
		if existing.RoomID == spectatorRoomID {
			return existing, nil // 已在该房间观战
		}
		// 离开之前的观战房间
		sm.leaveSpectatorLocked(userID)
	}

	// 创建观战者
	spectator := &Spectator{
		UserID:          userID,
		UserName:        userName,
		RoomID:          spectatorRoomID,
		State:           SpectatorWaiting,
		JoinTime:        time.Now(),
		LastUpdateTime:  time.Now(),
		CameraAngle:     0,
		ShowDanmaku:     true,
		ShowStats:       true,
	}

	spectatorRoom.Spectators = append(spectatorRoom.Spectators, spectator)
	spectatorRoom.ViewCount++

	sm.spectators[userID] = spectator
	sm.userToRoom[userID] = spectatorRoomID

	return spectator, nil
}

// leaveSpectatorLocked 内部方法：离开观战（需持有锁）
func (sm *SpectatorManager) leaveSpectatorLocked(userID string) {
	spectator, ok := sm.spectators[userID]
	if !ok {
		return
	}

	// 从房间中移除
	spectatorRoom, ok := sm.spectatorRooms[spectator.RoomID]
	if ok {
		for i, s := range spectatorRoom.Spectators {
			if s.UserID == userID {
				spectatorRoom.Spectators = append(spectatorRoom.Spectators[:i], spectatorRoom.Spectators[i+1:]...)
				break
			}
		}
	}

	spectator.State = SpectatorLeft

	delete(sm.spectators, userID)
	delete(sm.userToRoom, userID)
}

// LeaveSpectator 离开观战
func (sm *SpectatorManager) LeaveSpectator(userID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, ok := sm.spectators[userID]
	if !ok {
		return fmt.Errorf("未在观战中")
	}

	sm.leaveSpectatorLocked(userID)
	return nil
}

// UpdateSpectatorState 更新观战者状态
func (sm *SpectatorManager) UpdateSpectatorState(userID string, state SpectatorState) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	spectator, ok := sm.spectators[userID]
	if !ok {
		return fmt.Errorf("观战者不存在")
	}

	spectator.State = state
	spectator.LastUpdateTime = time.Now()

	return nil
}

// SetCameraAngle 设置视角
func (sm *SpectatorManager) SetCameraAngle(userID string, angle int) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	spectator, ok := sm.spectators[userID]
	if !ok {
		return fmt.Errorf("观战者不存在")
	}

	if angle < 0 || angle > 2 {
		return fmt.Errorf("无效的视角")
	}

	spectator.CameraAngle = angle
	spectator.LastUpdateTime = time.Now()

	return nil
}

// ToggleDanmaku 切换弹幕显示
func (sm *SpectatorManager) ToggleDanmaku(userID string, show bool) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	spectator, ok := sm.spectators[userID]
	if !ok {
		return fmt.Errorf("观战者不存在")
	}

	spectator.ShowDanmaku = show
	spectator.LastUpdateTime = time.Now()

	return nil
}

// ToggleStats 切换统计显示
func (sm *SpectatorManager) ToggleStats(userID string, show bool) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	spectator, ok := sm.spectators[userID]
	if !ok {
		return fmt.Errorf("观战者不存在")
	}

	spectator.ShowStats = show
	spectator.LastUpdateTime = time.Now()

	return nil
}

// GetSpectator 获取观战者信息
func (sm *SpectatorManager) GetSpectator(userID string) (*Spectator, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	spectator, ok := sm.spectators[userID]
	if !ok {
		return nil, fmt.Errorf("观战者不存在")
	}

	return spectator, nil
}

// GetSpectatorRoom 获取观战房间
func (sm *SpectatorManager) GetSpectatorRoom(roomID string) (*SpectatorRoom, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	room, ok := sm.spectatorRooms[roomID]
	if !ok {
		return nil, fmt.Errorf("观战房间不存在")
	}

	return room, nil
}

// GetSpectatorsByBattleRoom 根据战斗房间获取观战者列表
func (sm *SpectatorManager) GetSpectatorsByBattleRoom(battleRoomID string) []*Spectator {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	spectatorRoomID, ok := sm.battleToSpectator[battleRoomID]
	if !ok {
		return []*Spectator{}
	}

	spectatorRoom, ok := sm.spectatorRooms[spectatorRoomID]
	if !ok {
		return []*Spectator{}
	}

	return spectatorRoom.Spectators
}

// GetSpectatorCount 获取观战人数
func (sm *SpectatorManager) GetSpectatorCount(battleRoomID string) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	spectatorRoomID, ok := sm.battleToSpectator[battleRoomID]
	if !ok {
		return 0
	}

	spectatorRoom, ok := sm.spectatorRooms[spectatorRoomID]
	if !ok {
		return 0
	}

	count := 0
	for _, s := range spectatorRoom.Spectators {
		if s.State == SpectatorWatching {
			count++
		}
	}

	return count
}

// IsSpectating 检查用户是否在观战
func (sm *SpectatorManager) IsSpectating(userID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	_, ok := sm.spectators[userID]
	return ok
}

// GetSpectatorBattleRoom 获取观战者所在的战斗房间
func (sm *SpectatorManager) GetSpectatorBattleRoom(userID string) (string, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	spectator, ok := sm.spectators[userID]
	if !ok {
		return "", fmt.Errorf("未在观战中")
	}

	spectatorRoom, ok := sm.spectatorRooms[spectator.RoomID]
	if !ok {
		return "", fmt.Errorf("观战房间不存在")
	}

	return spectatorRoom.BattleRoomID, nil
}

// StartRecording 开始录制观战
func (sm *SpectatorManager) StartRecording(roomID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	spectatorRoom, ok := sm.spectatorRooms[roomID]
	if !ok {
		return fmt.Errorf("观战房间不存在")
	}

	spectatorRoom.IsRecording = true
	return nil
}

// StopRecording 停止录制观战
func (sm *SpectatorManager) StopRecording(roomID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	spectatorRoom, ok := sm.spectatorRooms[roomID]
	if !ok {
		return fmt.Errorf("观战房间不存在")
	}

	spectatorRoom.IsRecording = false
	return nil
}

// CloseSpectatorRoom 关闭观战房间
func (sm *SpectatorManager) CloseSpectatorRoom(battleRoomID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	spectatorRoomID, ok := sm.battleToSpectator[battleRoomID]
	if !ok {
		return nil // 没有观战房间
	}

	spectatorRoom, ok := sm.spectatorRooms[spectatorRoomID]
	if !ok {
		return nil
	}

	// 标记所有观战者为已离开
	for _, spectator := range spectatorRoom.Spectators {
		spectator.State = SpectatorLeft
		delete(sm.spectators, spectator.UserID)
		delete(sm.userToRoom, spectator.UserID)
	}

	// 删除观战房间
	delete(sm.spectatorRooms, spectatorRoomID)
	delete(sm.battleToSpectator, battleRoomID)

	return nil
}

// GetAllSpectatorRooms 获取所有观战房间
func (sm *SpectatorManager) GetAllSpectatorRooms() []*SpectatorRoom {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	rooms := make([]*SpectatorRoom, 0, len(sm.spectatorRooms))
	for _, room := range sm.spectatorRooms {
		rooms = append(rooms, room)
	}

	return rooms
}

// GetPopularBattleRooms 获取热门观战房间
func (sm *SpectatorManager) GetPopularBattleRooms(limit int) []*SpectatorRoom {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	rooms := make([]*SpectatorRoom, 0, len(sm.spectatorRooms))
	for _, room := range sm.spectatorRooms {
		rooms = append(rooms, room)
	}

	// 按观看次数排序
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].ViewCount > rooms[j].ViewCount
	})

	if limit > 0 && limit < len(rooms) {
		rooms = rooms[:limit]
	}

	return rooms
}

// CleanupInactiveSpectators 清理不活跃的观战者
func (sm *SpectatorManager) CleanupInactiveSpectators(timeout time.Duration) int {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	cutoff := time.Now().Add(-timeout)
	cleanedCount := 0

	for userID, spectator := range sm.spectators {
		if spectator.LastUpdateTime.Before(cutoff) {
			sm.leaveSpectatorLocked(userID)
			cleanedCount++
		}
	}

	return cleanedCount
}
