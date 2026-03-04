package room

import (
    "sync"
    "time"
    "github.com/gorilla/websocket"
)

type Room struct {
    ID        string
    Players   map[string]*Player
    Status    string // waiting, playing, finished
    CreatedAt time.Time
    mutex     sync.RWMutex
}

type Player struct {
    ID       string
    Name     string
    Conn     *websocket.Conn
    Send     chan []byte
    Team     string // blue/red
}

func NewRoom() *Room {
    return &Room{
        ID:        generateID(),
        Players:   make(map[string]*Player),
        Status:    "waiting",
        CreatedAt: time.Now(),
    }
}

func (r *Room) AddPlayer(p *Player) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    r.Players[p.ID] = p
}

func (r *Room) RemovePlayer(id string) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    delete(r.Players, id)
}

func (r *Room) Broadcast(msg []byte) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    for _, p := range r.Players {
        select {
        case p.Send <- msg:
        default:
        }
    }
}

func generateID() string {
    return "room_" + time.Now().Format("20060102150405")
}
