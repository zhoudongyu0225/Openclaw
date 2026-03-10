package handler

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Player 玩家连接
type Player struct {
	ID         string
	Conn       *websocket.Conn
	Send       chan []byte
	LastActive time.Time
	JoinAt     time.Time
	mu         sync.Mutex
}

// PlayerManager 玩家管理器
type PlayerManager struct {
	players map[string]*Player
	mutex   sync.RWMutex
}

// NewPlayerManager 创建玩家管理器
func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[string]*Player),
	}
}

// Add 添加玩家
func (m *PlayerManager) Add(id string, conn *websocket.Conn) *Player {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	p := &Player{
		ID:         id,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		LastActive: time.Now(),
		JoinAt:     time.Now(),
	}
	m.players[id] = p
	return p
}

// Remove 移除玩家
func (m *PlayerManager) Remove(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.players, id)
}

// Get 获取玩家
func (m *PlayerManager) Get(id string) *Player {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.players[id]
}

// GetAll 获取所有玩家
func (m *PlayerManager) GetAll() []*Player {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	list := make([]*Player, 0, len(m.players))
	for _, p := range m.players {
		list = append(list, p)
	}
	return list
}

// Count 玩家数量
func (m *PlayerManager) Count() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.players)
}

// Cleanup 超时玩家
func (m *PlayerManager) Cleanup(timeout time.Duration) []*Player {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	timeoutPlayers := make([]*Player, 0)

	for id, p := range m.players {
		if now.Sub(p.LastActive) > timeout {
			timeoutPlayers = append(timeoutPlayers, p)
			delete(m.players, id)
		}
	}

	return timeoutPlayers
}
