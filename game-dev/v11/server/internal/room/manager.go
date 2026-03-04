package room

import "sync"

type Manager struct {
    rooms map[string]*Room
    mutex sync.RWMutex
}

func NewManager() *Manager {
    return &Manager{
        rooms: make(map[string]*Room),
    }
}

func (m *Manager) CreateRoom() *Room {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    r := NewRoom()
    m.rooms[r.ID] = r
    return r
}

func (m *Manager) GetRoom(id string) *Room {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    return m.rooms[id]
}

func (m *Manager) DeleteRoom(id string) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    delete(m.rooms, id)
}
