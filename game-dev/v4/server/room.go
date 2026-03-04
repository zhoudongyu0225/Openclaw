package main

import (
    "sync"
    "time"
)

// Room 房间结构
type Room struct {
    ID        string
    Name      string
    Players   map[string]*Player
    Status    string // waiting, playing, paused
    Settings  RoomSettings
    CreatedAt time.Time
    StartedAt time.Time
    mutex     sync.RWMutex
}

// RoomSettings 房间设置
type RoomSettings struct {
    MaxPlayers    int     // 最大玩家数
    GameMode     string  // 游戏模式: rank, normal, practice
    MapID        string  // 地图ID
    TimeLimit    int     // 时间限制(秒)
    WinScore     int     // 获胜分数
}

// Player 玩家结构
type Player struct {
    ID       string
    Name     string
    Team     string // "blue" or "red"
    Score    int
    Kills    int
    Deaths   int
    IsReady  bool
    JoinAt  time.Time
}

// NewRoom 创建新房间
func NewRoom(name string, settings RoomSettings) *Room {
    return &Room{
        ID:        generateRoomID(),
        Name:      name,
        Players:   make(map[string]*Player),
        Status:    "waiting",
        Settings:  settings,
        CreatedAt: time.Now(),
    }
}

// AddPlayer 添加玩家
func (r *Room) AddPlayer(p *Player) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    if len(r.Players) >= r.Settings.MaxPlayers {
        return ErrRoomFull
    }
    
    // 自动分配队伍
    blue, red := 0, 0
    for _, player := range r.Players {
        if player.Team == "blue" {
            blue++
        } else {
            red++
        }
    }
    
    if blue <= red {
        p.Team = "blue"
    } else {
        p.Team = "red"
    }
    
    r.Players[p.ID] = p
    return nil
}

// RemovePlayer 移除玩家
func (r *Room) RemovePlayer(id string) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    delete(r.Players, id)
}

// StartGame 开始游戏
func (r *Room) StartGame() error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    if len(r.Players) < 2 {
        return ErrNotEnoughPlayers
    }
    
    r.Status = "playing"
    r.StartedAt = time.Now()
    return nil
}

// EndGame 结束游戏
func (r *Room) EndGame(winner string) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    r.Status = "finished"
}

func generateRoomID() string {
    return "room_" + time.Now().Format("20060102150405")
}

var (
    ErrRoomFull = fmt.Errorf("房间已满")
    ErrNotEnoughPlayers = fmt.Errorf("玩家数量不足")
)
