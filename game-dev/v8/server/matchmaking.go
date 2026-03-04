package main

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
)

// MatchmakingService 匹配服务
type MatchmakingService struct {
    players    map[string]*Player
    waiting    []*Player
    matchQueue chan *Match
    mutex      sync.RWMutex
}

type Player struct {
    ID        string
    Rating    int
    JoinTime  time.Time
}

type Match struct {
    ID        string
    BlueTeam  []*Player
    RedTeam   []*Player
    CreateTime time.Time
}

func NewMatchmaking() *MatchmakingService {
    ms := &MatchmakingService{
        players:    make(map[string]*Player),
        waiting:    make([]*Player, 0),
        matchQueue: make(chan *Match, 100),
    }
    go ms.matchLoop()
    return ms
}

func (ms *MatchmakingService) AddPlayer(p *Player) {
    ms.mutex.Lock()
    defer ms.mutex.Unlock()
    ms.players[p.ID] = p
    ms.waiting = append(ms.waiting, p)
}

func (ms *MatchmakingService) matchLoop() {
    ticker := time.NewTicker(500 * time.Millisecond)
    for range ticker.C {
        ms.tryMatch()
    }
}

func (ms *MatchmakingService) tryMatch() {
    ms.mutex.Lock()
    defer ms.mutex.Unlock()
    if len(ms.waiting) < 2 { return }
    
    for i := 0; i < len(ms.waiting)-1; i++ {
        for j := i + 1; j < len(ms.waiting); j++ {
            if abs(ms.waiting[i].Rating - ms.waiting[j].Rating) <= 200 {
                match := &Match{
                    ID: "match_" + time.Now().Format("20060102150405"),
                    BlueTeam: []*Player{ms.waiting[i]},
                    RedTeam:  []*Player{ms.waiting[j]},
                    CreateTime: time.Now(),
                }
                ms.matchQueue <- match
                ms.waiting = append(ms.waiting[:j], ms.waiting[j+1:]...)
                ms.waiting = append(ms.waiting[:i], ms.waiting[i+1:]...)
                return
            }
        }
    }
}

func abs(x int) int { if x < 0 { return -x }; return x }

func (ms *MatchmakingService) HandleMatch(w http.ResponseWriter, r *http.Request) {
    var req struct {
        PlayerID string `json:"player_id"`
        Rating   int    `json:"rating"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    ms.AddPlayer(&Player{ID: req.PlayerID, Rating: req.Rating, JoinTime: time.Now()})
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "matching"})
}
