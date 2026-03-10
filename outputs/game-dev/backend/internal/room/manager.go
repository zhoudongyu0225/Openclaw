package room

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"barrage-game/internal/model"
)

// Manager 房间管理器 - 线程安全
type Manager struct {
	rooms      map[string]*Room
	playerMap  map[string]string // playerID -> roomID
	roomIndex  map[string][]string // mode -> roomIDs
	mutex      sync.RWMutex
	heartbeat  time.Duration
	closeCh    chan struct{}
}

// NewManager 创建房间管理器
func NewManager(heartbeat time.Duration) *Manager {
	m := &Manager{
		rooms:     make(map[string]*Room),
		playerMap: make(map[string]string),
		roomIndex: make(map[string][]string),
		heartbeat: heartbeat,
		closeCh:   make(chan struct{}),
	}
	// 启动房间清理定时器
	go m.cleanupLoop()
	return m
}

// CreateRoom 创建房间
func (m *Manager) CreateRoom(req *model.CreateRoomReq) (*Room, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	settings := RoomSettings{
		MaxPlayers: 8,
		GameMode:    req.Mode,
		MapID:       req.MapID,
		TimeLimit:   300,
		WinScore:    1000,
	}
	if req.MaxPlayers > 0 {
		settings.MaxPlayers = req.MaxPlayers
	}

	r := newRoom(req.Name, settings)
	m.rooms[r.ID] = r
	m.roomIndex[req.Mode] = append(m.roomIndex[req.Mode], r.ID)

	return r, nil
}

// JoinRoom 玩家加入房间
func (m *Manager) JoinRoom(roomID, playerID, playerName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	r, ok := m.rooms[roomID]
	if !ok {
		return fmt.Errorf("房间不存在")
	}

	if r.Status != "waiting" {
		return fmt.Errorf("游戏已开始")
	}

	player := &Player{
		ID:       playerID,
		Name:     playerName,
		IsReady:  false,
		JoinAt:   time.Now(),
		Conn:     nil, // WebSocket连接由外部管理
	}

	if err := r.addPlayer(player); err != nil {
		return err
	}

	m.playerMap[playerID] = roomID
	return nil
}

// LeaveRoom 玩家离开房间
func (m *Manager) LeaveRoom(playerID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	roomID, ok := m.playerMap[playerID]
	if !ok {
		return fmt.Errorf("玩家不在房间中")
	}

	r, ok := m.rooms[roomID]
	if !ok {
		delete(m.playerMap, playerID)
		return nil
	}

	r.removePlayer(playerID)
	delete(m.playerMap, playerID)

	// 房间为空时删除
	if len(r.Players) == 0 {
		delete(m.rooms, roomID)
		// 从索引中移除
		for mode, ids := range m.roomIndex {
			for i, id := range ids {
				if id == roomID {
					m.roomIndex[mode] = append(ids[:i], ids[i+1:]...)
					break
				}
			}
		}
	}

	return nil
}

// GetRoom 获取房间
func (m *Manager) GetRoom(roomID string) *Room {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.rooms[roomID]
}

// ListRooms 列出所有房间
func (m *Manager) ListRooms() []*Room {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	list := make([]*Room, 0, len(m.rooms))
	for _, r := range m.rooms {
		list = append(list, r)
	}
	return list
}

// ListRoomsByMode 按模式列出房间
func (m *Manager) ListRoomsByMode(mode string) []*Room {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	ids, ok := m.roomIndex[mode]
	if !ok {
		return nil
	}

	list := make([]*Room, 0, len(ids))
	for _, id := range ids {
		if r := m.rooms[id]; r != nil {
			list = append(list, r)
		}
	}
	return list
}

// Broadcast 广播消息到房间
func (m *Manager) Broadcast(roomID string, msg *model.WSMessage) error {
	r := m.GetRoom(roomID)
	if r == nil {
		return fmt.Errorf("房间不存在")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, p := range r.Players {
		if p.Conn != nil {
			p.Conn.WriteMessage(1, data) // 1 = text
		}
	}
	return nil
}

// GetPlayerRoom 获取玩家所在房间
func (m *Manager) GetPlayerRoom(playerID string) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.playerMap[playerID]
}

// cleanupLoop 清理过期房间
func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(m.heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.closeCh:
			return
		}
	}
}

// cleanup 清理长时间无人的房间
func (m *Manager) cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	for id, r := range m.rooms {
		if len(r.Players) == 0 && now.Sub(r.CreatedAt) > 30*time.Minute {
			delete(m.rooms, id)
		}
	}
}

// Close 关闭管理器
func (m *Manager) Close() {
	close(m.closeCh)
}
