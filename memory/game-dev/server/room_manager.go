package main

import (
	"sync"
	"time"
)

// ============================================
// 房间管理器 (Room Manager)
// ============================================

type RoomManager struct {
	rooms    map[string]*LiveRoom
	playerRoom map[string]string // playerID -> roomID
	mu       sync.RWMutex
}

// 新建房间管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms:      make(map[string]*LiveRoom),
		playerRoom: make(map[string]string),
	}
}

// 创建房间
func (rm *RoomManager) CreateRoom(roomID, hostID, hostName string) *LiveRoom {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	room := NewLiveRoom(roomID, hostID, hostName)
	room.Battle.State.IsRunning = false
	
	rm.rooms[roomID] = room
	rm.playerRoom[hostID] = roomID
	
	return room
}

// 加入房间
func (rm *RoomManager) JoinRoom(roomID, playerID, playerName string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	room, ok := rm.rooms[roomID]
	if !ok {
		return nil // 房间不存在
	}
	
	room.JoinViewer(playerID, playerName)
	rm.playerRoom[playerID] = roomID
	
	return nil
}

// 离开房间
func (rm *RoomManager) LeaveRoom(playerID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	roomID, ok := rm.playerRoom[playerID]
	if !ok {
		return
	}
	
	if room, ok := rm.rooms[roomID]; ok {
		room.LeaveViewer(playerID)
	}
	
	delete(rm.playerRoom, playerID)
}

// 获取房间
func (rm *RoomManager) GetRoom(roomID string) *LiveRoom {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.rooms[roomID]
}

// 获取玩家所在房间
func (rm *RoomManager) GetPlayerRoom(playerID string) *LiveRoom {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	roomID, ok := rm.playerRoom[playerID]
	if !ok {
		return nil
	}
	
	return rm.rooms[roomID]
}

// 列出所有房间
func (rm *RoomManager) ListRooms() []*LiveRoom {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	rooms := make([]*LiveRoom, 0, len(rm.rooms))
	for _, room := range rm.rooms {
		rooms = append(rooms, room)
	}
	
	return rooms
}

// 删除房间
func (rm *RoomManager) DeleteRoom(roomID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	if room, ok := rm.rooms[roomID]; ok {
		// 清理玩家映射
		for playerID := range room.Viewers {
			delete(rm.playerRoom, playerID)
		}
	}
	
	delete(rm.rooms, roomID)
}

// ============================================
// 玩家管理器 (Player Manager)
// ============================================

type PlayerManager struct {
	players    map[string]*Player
	mu         sync.RWMutex
}

type Player struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Level        int            `json:"level"`
	Exp          int            `json:"exp"`
	Coins        int            `json:"coins"`
	Gems         int            `json:"gems"`
	WinCount     int            `json:"winCount"`
	LoseCount    int            `json:"loseCount"`
	TotalScore   int            `json:"totalScore"`
	HighScore    int            `json:"highScore"`
	Towers       map[string]int `json:"towers"` // towerID -> level
	Skins        []string       `json:"skins"`
	LastLogin    time.Time      `json:"lastLogin"`
	LoginDays    int            `json:"loginDays"`
}

// 新建玩家
func NewPlayer(id, name string) *Player {
	return &Player{
		ID:        id,
		Name:      name,
		Level:     1,
		Exp:       0,
		Coins:     100,
		Gems:      0,
		Towers:    make(map[string]int),
		Skins:     []string{"default"},
		LoginDays: 1,
		LastLogin: time.Now(),
	}
}

// 获取胜率
func (p *Player) GetWinRate() float64 {
	total := p.WinCount + p.LoseCount
	if total == 0 {
		return 0
	}
	return float64(p.WinCount) / float64(total) * 100
}

// 升级检查
func (p *Player) CheckLevelUp() bool {
	expNeeded := p.Level * 100 // 简单公式
	for p.Exp >= expNeeded {
		p.Exp -= expNeeded
		p.Level++
		expNeeded = p.Level * 100
		return true
	}
	return false
}

// 添加经验
func (p *Player) AddExp(exp int) {
	p.Exp += exp
	p.CheckLevelUp()
}

// 记录胜利
func (p *Player) RecordWin(score int) {
	p.WinCount++
	p.TotalScore += score
	if score > p.HighScore {
		p.HighScore = score
	}
}

// 记录失败
func (p *Player) RecordLose() {
	p.LoseCount++
}

// 新建玩家管理器
func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[string]*Player),
	}
}

// 创建或获取玩家
func (pm *PlayerManager) GetOrCreatePlayer(id, name string) *Player {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if player, ok := pm.players[id]; ok {
		return player
	}
	
	player := NewPlayer(id, name)
	pm.players[id] = player
	
	return player
}

// 获取玩家
func (pm *PlayerManager) GetPlayer(id string) *Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.players[id]
}

// 更新玩家
func (pm *PlayerManager) UpdatePlayer(player *Player) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.players[player.ID] = player
}

// 排行榜获取
func (pm *PlayerManager) GetLeaderboard(limit int) []*Player {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	players := make([]*Player, 0, len(pm.players))
	for _, p := range pm.players {
		players = append(players, p)
	}
	
	// 排序
	for i := 0; i < len(players)-1; i++ {
		for j := i + 1; j < len(players); j++ {
			if players[j].HighScore > players[i].HighScore {
				players[i], players[j] = players[j], players[i]
			}
		}
	}
	
	if limit > 0 && limit < len(players) {
		players = players[:limit]
	}
	
	return players
}

// ============================================
// 匹配系统 (Matchmaking)
// ============================================

type Matchmaker struct {
	queue     []*MatchRequest
	mu        sync.RWMutex
}

type MatchRequest struct {
	PlayerID   string
	PlayerName string
	MMR        int // 匹配分
	Timestamp  time.Time
}

// 新建匹配器
func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		queue: make([]*MatchRequest, 0),
	}
}

// 加入匹配队列
func (mm *Matchmaker) JoinQueue(req *MatchRequest) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	req.Timestamp = time.Now()
	mm.queue = append(mm.queue, req)
}

// 从匹配队列移除
func (mm *Matchmaker) LeaveQueue(playerID string) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	for i, req := range mm.queue {
		if req.PlayerID == playerID {
			mm.queue = append(mm.queue[:i], mm.queue[i+1:]...)
			return
		}
	}
}

// 尝试匹配
func (mm *Matchmaker) TryMatch() (string, string, bool) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	if len(mm.queue) < 2 {
		return "", "", false
	}
	
	// 简单匹配: 取队列前两个玩家
	req1 := mm.queue[0]
	req2 := mm.queue[1]
	
	// MMR差距不超过500
	if abs(req1.MMR-req2.MMR) > 500 {
		return "", "", false
	}
	
	// 移除已匹配的
	mm.queue = mm.queue[2:]
	
	return req1.PlayerID, req2.PlayerID, true
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// ============================================
// 排行榜管理器
// ============================================

type LeaderboardManager struct {
	scoreLeaderboard   []*Player
	winRateLeaderboard []*Player
	wealthLeaderboard  []*Player
	mu                 sync.RWMutex
}

// 新建排行榜管理器
func NewLeaderboardManager() *LeaderboardManager {
	return &LeaderboardManager{}
}

// 更新排行榜
func (lm *LeaderboardManager) Update(playerMgr *PlayerManager) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	players := playerMgr.GetLeaderboard(0)
	
	// 分数榜
	lm.scoreLeaderboard = make([]*Player, len(players))
	copy(lm.scoreLeaderboard, players)
	// 排序
	for i := 0; i < len(lm.scoreLeaderboard)-1; i++ {
		for j := i + 1; j < len(lm.scoreLeaderboard); j++ {
			if lm.scoreLeaderboard[j].HighScore > lm.scoreLeaderboard[i].HighScore {
				lm.scoreLeaderboard[i], lm.scoreLeaderboard[j] = lm.scoreLeaderboard[j], lm.scoreLeaderboard[i]
			}
		}
	}
	
	// 胜率榜 (至少10场)
	lm.winRateLeaderboard = make([]*Player, 0)
	for _, p := range players {
		if p.WinCount+p.LoseCount >= 10 {
			lm.winRateLeaderboard = append(lm.winRateLeaderboard, p)
		}
	}
	for i := 0; i < len(lm.winRateLeaderboard)-1; i++ {
		for j := i + 1; j < len(lm.winRateLeaderboard); j++ {
			if lm.winRateLeaderboard[j].GetWinRate() > lm.winRateLeaderboard[i].GetWinRate() {
				lm.winRateLeaderboard[i], lm.winRateLeaderboard[j] = lm.winRateLeaderboard[j], lm.winRateLeaderboard[i]
			}
		}
	}
	
	// 财富榜
	lm.wealthLeaderboard = make([]*Player, len(players))
	copy(lm.wealthLeaderboard, players)
	for i := 0; i < len(lm.wealthLeaderboard)-1; i++ {
		for j := i + 1; j < len(lm.wealthLeaderboard); j++ {
			if lm.wealthLeaderboard[j].Coins > lm.wealthLeaderboard[i].Coins {
				lm.wealthLeaderboard[i], lm.wealthLeaderboard[j] = lm.wealthLeaderboard[j], lm.wealthLeaderboard[i]
			}
		}
	}
}

// 获取分数榜
func (lm *LeaderboardManager) GetScoreLeaderboard(limit int) []*Player {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	if limit > 0 && limit < len(lm.scoreLeaderboard) {
		return lm.scoreLeaderboard[:limit]
	}
	return lm.scoreLeaderboard
}

// 获取胜率榜
func (lm *LeaderboardManager) GetWinRateLeaderboard(limit int) []*Player {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	if limit > 0 && limit < len(lm.winRateLeaderboard) {
		return lm.winRateLeaderboard[:limit]
	}
	return lm.winRateLeaderboard
}

// 获取财富榜
func (lm *LeaderboardManager) GetWealthLeaderboard(limit int) []*Player {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	if limit > 0 && limit < len(lm.wealthLeaderboard) {
		return lm.wealthLeaderboard[:limit]
	}
	return lm.wealthLeaderboard
}

// 获取玩家排名
func (lm *LeaderboardManager) GetPlayerRank(playerID string, playerMgr *PlayerManager) (int, int) {
	player := playerMgr.GetPlayer(playerID)
	if player == nil {
		return -1, -1
	}
	
	scoreRank := -1
	for i, p := range lm.scoreLeaderboard {
		if p.ID == playerID {
			scoreRank = i + 1
			break
		}
	}
	
	winRateRank := -1
	for i, p := range lm.winRateLeaderboard {
		if p.ID == playerID {
			winRateRank = i + 1
			break
		}
	}
	
	return scoreRank, winRateRank
}
