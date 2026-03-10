package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// ============================================
// 排行榜系统 (Leaderboard System)
// ============================================

type LeaderboardType int

const (
	LeaderboardTypeScore LeaderboardType = iota // 分数榜
	LeaderboardTypeRich                        // 财富榜
	LeaderboardTypeWinRate                     // 胜率榜
	LeaderboardTypeKill                        // 击杀榜
)

// 排行榜条目
type LeaderboardEntry struct {
	Rank      int       `json:"rank"`
	UserID    string    `json:"userId"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Value     int64     `json:"value"`      // 分数/金币/击杀数
	WinCount  int       `json:"winCount"`   // 胜场
	LoseCount int       `json:"loseCount"`  // 负场
	WinRate   float64   `json:"winRate"`    // 胜率
	UpdatedAt time.Time `json:"updatedAt"`
}

// 排行榜
type Leaderboard struct {
	Type      LeaderboardType `json:"type"`
	Name      string         `json:"name"`
	Entries   []*LeaderboardEntry
	mu        sync.RWMutex
	updatedAt time.Time
}

// 排行榜管理器
type LeaderboardManager struct {
	Leaderboards map[LeaderboardType]*Leaderboard
	mu           sync.RWMutex
}

func NewLeaderboardManager() *LeaderboardManager {
	return &LeaderboardManager{
		Leaderboards: map[LeaderboardType]*Leaderboard{
			LeaderboardTypeScore:  {Type: LeaderboardTypeScore, Name: "积分榜", Entries: make([]*LeaderboardEntry, 0)},
			LeaderboardTypeRich:   {Type: LeaderboardTypeRich, Name: "财富榜", Entries: make([]*LeaderboardEntry, 0)},
			LeaderboardTypeWinRate:{Type: LeaderboardTypeWinRate, Name: "胜率榜", Entries: make([]*LeaderboardEntry, 0)},
			LeaderboardTypeKill:   {Type: LeaderboardTypeKill, Name: "击杀榜", Entries: make([]*LeaderboardEntry, 0)},
		},
	}
}

// 更新玩家数据
func (lm *LeaderboardManager) UpdatePlayer(userID, nickname, avatar string, score, money, kills int, win, lose bool) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	now := time.Now()

	// 更新各排行榜
	for _, lb := range lm.Leaderboards {
		entry := lb.getOrCreateEntry(userID)
		entry.Nickname = nickname
		entry.Avatar = avatar
		entry.UpdatedAt = now

		switch lb.Type {
		case LeaderboardTypeScore:
			entry.Value = int64(score)
		case LeaderboardTypeRich:
			entry.Value = int64(money)
		case LeaderboardTypeKill:
			entry.Value = int64(kills)
		}

		if win {
			entry.WinCount++
		} else if lose {
			entry.LoseCount++
		}

		// 计算胜率
		total := entry.WinCount + entry.LoseCount
		if total > 0 {
			entry.WinRate = float64(entry.WinCount) / float64(total) * 100
		}
	}

	// 重新排序各排行榜
	for _, lb := range lm.Leaderboards {
		lb.sort()
		lb.updateRanks()
		lb.updatedAt = now
	}
}

// 获取或创建排行榜条目
func (lb *Leaderboard) getOrCreateEntry(userID string) *LeaderboardEntry {
	for _, e := range lb.Entries {
		if e.UserID == userID {
			return e
		}
	}
	entry := &LeaderboardEntry{
		UserID:    userID,
		UpdatedAt: time.Now(),
	}
	lb.Entries = append(lb.Entries, entry)
	return entry
}

// 排序
func (lb *Leaderboard) sort() {
	switch lb.Type {
	case LeaderboardTypeScore, LeaderboardTypeRich, LeaderboardTypeKill:
		sort.Slice(lb.Entries, func(i, j int) bool {
			return lb.Entries[i].Value > lb.Entries[j].Value
		})
	case LeaderboardTypeWinRate:
		sort.Slice(lb.Entries, func(i, j int) bool {
			// 至少5场才上榜
			if lb.Entries[i].WinCount+lb.Entries[i].LoseCount < 5 {
				return false
			}
			if lb.Entries[j].WinCount+lb.Entries[j].LoseCount < 5 {
				return true
			}
			return lb.Entries[i].WinRate > lb.Entries[j].WinRate
		})
	}
}

// 更新排名
func (lb *Leaderboard) updateRanks() {
	for i, e := range lb.Entries {
		e.Rank = i + 1
	}
}

// 获取排行榜 (前N名)
func (lm *LeaderboardManager) GetLeaderboard(lbType LeaderboardType, limit int) []*LeaderboardEntry {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	lb, ok := lm.Leaderboards[lbType]
	if !ok {
		return nil
	}

	if limit <= 0 || limit > len(lb.Entries) {
		limit = len(lb.Entries)
	}
	return lb.Entries[:limit]
}

// 获取玩家排名
func (lm *LeaderboardManager) GetPlayerRank(userID string, lbType LeaderboardType) int {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	lb, ok := lm.Leaderboards[lbType]
	if !ok {
		return -1
	}

	for i, e := range lb.Entries {
		if e.UserID == userID {
			return i + 1
		}
	}
	return -1
}

// 获取玩家在各排行榜的排名
func (lm *LeaderboardManager) GetPlayerAllRanks(userID string) map[LeaderboardType]int {
	return map[LeaderboardType]int{
		LeaderboardTypeScore:   lm.GetPlayerRank(userID, LeaderboardTypeScore),
		LeaderboardTypeRich:   lm.GetPlayerRank(userID, LeaderboardTypeRich),
		LeaderboardTypeWinRate:lm.GetPlayerRank(userID, LeaderboardTypeWinRate),
		LeaderboardTypeKill:   lm.GetPlayerRank(userID, LeaderboardTypeKill),
	}
}

// ============================================
// 玩家数据管理器
// ============================================

type PlayerProfile struct {
	UserID       string    `json:"userId"`
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
	Level        int       `json:"level"`
	Exp          int       `json:"exp"`
	Score        int       `json:"score"`
	Money        int       `json:"money"`
	Gem          int       `json:"gem"`         // 钻石
	TotalScore   int       `json:"totalScore"`  // 历史总积分
	WinCount     int       `json:"winCount"`
	LoseCount    int       `json:"loseCount"`
	KillCount    int       `json:"killCount"`
	MaxWave      int       `json:"maxWave"`     // 最远波次
	TotalGames   int       `json:"totalGames"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	mu           sync.RWMutex
}

// 玩家管理器
type PlayerManager struct {
	Players    map[string]*PlayerProfile
	Leaderboard *LeaderboardManager
	mu          sync.RWMutex
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		Players:    make(map[string]*PlayerProfile),
		Leaderboard: NewLeaderboardManager(),
	}
}

// 创建新玩家
func (pm *PlayerManager) CreatePlayer(userID, nickname, avatar string) *PlayerProfile {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, ok := pm.Players[userID]; ok {
		return pm.Players[userID]
	}

	player := &PlayerProfile{
		UserID:     userID,
		Nickname:   nickname,
		Avatar:     avatar,
		Level:      1,
		Exp:        0,
		Score:      1000,
		Money:      100,
		Gem:        0,
		TotalScore: 0,
		WinCount:   0,
		LoseCount:  0,
		KillCount:  0,
		MaxWave:    0,
		TotalGames: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	pm.Players[userID] = player
	return player
}

// 获取玩家
func (pm *PlayerManager) GetPlayer(userID string) *PlayerProfile {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.Players[userID]
}

// 更新玩家分数
func (pm *PlayerManager) UpdateScore(userID string, scoreDelta int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.Players[userID]
	if !ok {
		return
	}

	player.Score += scoreDelta
	if player.Score < 0 {
		player.Score = 0
	}
	player.TotalScore += scoreDelta
	player.UpdatedAt = time.Now()

	// 更新排行榜
	pm.Leaderboard.UpdatePlayer(
		player.UserID,
		player.Nickname,
		player.Avatar,
		player.Score,
		player.Money,
		player.KillCount,
		false, false,
	)
}

// 更新玩家金币
func (pm *PlayerManager) UpdateMoney(userID string, moneyDelta int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.Players[userID]
	if !ok {
		return
	}

	player.Money += moneyDelta
	if player.Money < 0 {
		player.Money = 0
	}
	player.UpdatedAt = time.Now()

	// 更新财富榜
	pm.Leaderboard.UpdatePlayer(
		player.UserID,
		player.Nickname,
		player.Avatar,
		player.Score,
		player.Money,
		player.KillCount,
		false, false,
	)
}

// 玩家胜利
func (pm *PlayerManager) PlayerWin(userID string, scoreGain, moneyGain, killCount int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.Players[userID]
	if !ok {
		return
	}

	player.WinCount++
	player.TotalGames++
	player.Score += scoreGain
	player.Money += moneyGain
	player.KillCount += killCount
	player.Exp += 10
	player.UpdatedAt = time.Now()

	// 升级检查
	pm.checkLevelUp(player)

	// 更新排行榜
	pm.Leaderboard.UpdatePlayer(
		player.UserID,
		player.Nickname,
		player.Avatar,
		player.Score,
		player.Money,
		player.KillCount,
		true, false,
	)
}

// 玩家失败
func (pm *PlayerManager) PlayerLose(userID string, moneyGain int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.Players[userID]
	if !ok {
		return
	}

	player.LoseCount++
	player.TotalGames++
	player.Money += moneyGain
	player.Exp += 2
	player.UpdatedAt = time.Now()

	// 升级检查
	pm.checkLevelUp(player)

	// 更新排行榜
	pm.Leaderboard.UpdatePlayer(
		player.UserID,
		player.Nickname,
		player.Avatar,
		player.Score,
		player.Money,
		player.KillCount,
		false, true,
	)
}

// 升级检查
func (pm *PlayerManager) checkLevelUp(player *PlayerProfile) {
	expNeeded := player.Level * 100 // 每级需要100*等级经验
	for player.Exp >= expNeeded {
		player.Exp -= expNeeded
		player.Level++
		expNeeded = player.Level * 100
	}
}

// 钻石购买 (模拟)
func (pm *PlayerManager) PurchaseWithGem(userID string, gemCost int, itemID string) bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	player, ok := pm.Players[userID]
	if !ok {
		return false
	}

	if player.Gem < gemCost {
		return false
	}

	player.Gem -= gemCost
	player.UpdatedAt = time.Now()

	// TODO: 发放物品
	fmt.Printf("Player %s purchased %s for %d gems\n", userID, itemID, gemCost)
	return true
}

// ============================================
// 示例代码
// ============================================

/*
// 使用示例
func main() {
	// 创建玩家管理器
	pm := NewPlayerManager()

	// 创建玩家
	p1 := pm.CreatePlayer("user001", "主播A", "avatar1.png")
	p2 := pm.CreatePlayer("user002", "观众B", "avatar2.png")

	// 模拟游戏结果
	pm.PlayerWin("user001", 50, 100, 10)
	pm.PlayerWin("user002", 30, 80, 5)
	pm.PlayerLose("user002", 50)

	// 获取排行榜
	scoreBoard := pm.Leaderboard.GetLeaderboard(LeaderboardTypeScore, 10)
	fmt.Println("=== 积分榜 ===")
	for _, e := range scoreBoard {
		fmt.Printf("#%d %s: %d\n", e.Rank, e.Nickname, e.Value)
	}

	richBoard := pm.Leaderboard.GetLeaderboard(LeaderboardTypeRich, 10)
	fmt.Println("\n=== 财富榜 ===")
	for _, e := range richBoard {
		fmt.Printf("#%d %s: %d\n", e.Rank, e.Nickname, e.Value)
	}

	// 获取玩家排名
	rank := pm.Leaderboard.GetPlayerRank("user001", LeaderboardTypeScore)
	fmt.Printf("\n用户user001积分榜排名: #%d\n", rank)
}
*/
