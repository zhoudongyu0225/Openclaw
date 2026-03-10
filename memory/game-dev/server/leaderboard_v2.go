package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

// ==================== 排行榜系统扩展 ====================

// LeaderboardType 排行榜类型
type LeaderboardType int

const (
	LeaderboardTypeScore   LeaderboardType = iota // 分数榜
	LeaderboardTypeLevel                          // 等级榜
	LeaderboardTypeCombat                         // 战力榜
	LeaderboardTypeAchievement                   // 成就榜
	LeaderboardTypeStreak                        // 连续登录榜
	LeaderboardTypeRichest                       // 富豪榜
	LeaderboardTypeKills                         // 击杀榜
	LeaderboardTypeWins                          // 胜率榜
)

// LeaderboardEntry 排行榜条目
type LeaderboardEntry struct {
	PlayerID    string    `json:"player_id"`     // 玩家ID
	PlayerName string    `json:"player_name"`  // 玩家名称
	Avatar      string    `json:"avatar"`       // 头像
	Value       float64   `json:"value"`         // 数值
	Rank        int       `json:"rank"`          // 排名
	Change      int       `json:"change"`        // 排名变化
	LastUpdate  time.Time `json:"last_update"` // 上次更新
}

// Leaderboard 排行榜
type Leaderboard struct {
	Type         LeaderboardType     `json:"type"`          // 排行榜类型
	Name         string             `json:"name"`          // 排行榜名称
	Entries      []*LeaderboardEntry `json:"entries"`       // 排行榜条目
	UpdatedAt    time.Time          `json:"updated_at"`    // 更新时间
	mu           sync.RWMutex
	maxSize      int                 // 最大条目数
	expireTime   time.Duration       // 缓存过期时间
}

// NewLeaderboard 创建排行榜
func NewLeaderboard(ltype LeaderboardType, name string, maxSize int) *Leaderboard {
	return &Leaderboard{
		Type:       ltype,
		Name:       name,
		Entries:    make([]*LeaderboardEntry, 0, maxSize),
		UpdatedAt:  time.Now(),
		maxSize:    maxSize,
		expireTime: 5 * time.Minute,
	}
}

// UpdateEntry 更新条目
func (l *Leaderboard) UpdateEntry(playerID, playerName, avatar string, value float64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 查找现有条目
	existingIdx := -1
	for i, e := range l.Entries {
		if e.PlayerID == playerID {
			existingIdx = i
			break
		}
	}

	oldRank := 0
	if existingIdx >= 0 {
		oldRank = l.Entries[existingIdx].Rank
		// 更新现有条目
		l.Entries[existingIdx].PlayerName = playerName
		l.Entries[existingIdx].Avatar = avatar
		l.Entries[existingIdx].Value = value
		l.Entries[existingIdx].LastUpdate = time.Now()
	} else {
		// 新增条目
		entry := &LeaderboardEntry{
			PlayerID:   playerID,
			PlayerName: playerName,
			Avatar:     avatar,
			Value:      value,
			LastUpdate: time.Now(),
		}
		l.Entries = append(l.Entries, entry)
	}

	// 重新排序
	l.sortEntries()

	// 更新排名
	for i, e := range l.Entries {
		e.Rank = i + 1
		e.Change = oldRank - e.Rank
	}

	// 截断超出大小的部分
	if len(l.Entries) > l.maxSize {
		l.Entries = l.Entries[:l.maxSize]
	}

	l.UpdatedAt = time.Now()
}

// sortEntries 根据排行榜类型排序
func (l *Leaderboard) sortEntries() {
	sort.Slice(l.Entries, func(i, j int) bool {
		switch l.Type {
		case LeaderboardTypeWins:
			// 胜率榜：降序
			return l.Entries[i].Value > l.Entries[j].Value
		default:
			// 其他榜：降序
			return l.Entries[i].Value > l.Entries[j].Value
		}
	})
}

// GetRank 获取玩家排名
func (l *Leaderboard) GetRank(playerID string) (int, float64) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for i, e := range l.Entries {
		if e.PlayerID == playerID {
			return i + 1, e.Value
		}
	}
	return -1, 0
}

// GetTopN 获取前N名
func (l *Leaderboard) GetTopN(n int) []*LeaderboardEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if n > len(l.Entries) {
		n = len(l.Entries)
	}
	result := make([]*LeaderboardEntry, n)
	copy(result, l.Entries[:n])
	return result
}

// GetAround 获取指定玩家附近的排名
func (l *Leaderboard) GetAround(playerID string, count int) ([]*LeaderboardEntry, int) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	playerIdx := -1
	for i, e := range l.Entries {
		if e.PlayerID == playerID {
			playerIdx = i
			break
		}
	}

	if playerIdx == -1 {
		return nil, -1
	}

	start := playerIdx - count/2
	if start < 0 {
		start = 0
	}
	end := playerIdx + count/2 + 1
	if end > len(l.Entries) {
		end = len(l.Entries)
	}

	result := make([]*LeaderboardEntry, end-start)
	copy(result, l.Entries[start:end])

	// 调整相对排名
	for i := range result {
		result[i].Rank = start + i + 1
	}

	return result, playerIdx + 1
}

// RemoveEntry 移除条目
func (l *Leaderboard) RemoveEntry(playerID string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, e := range l.Entries {
		if e.PlayerID == playerID {
			l.Entries = append(l.Entries[:i], l.Entries[i+1:]...)
			l.UpdatedAt = time.Now()
			break
		}
	}
}

// Clear 清空排行榜
func (l *Leaderboard) Clear()()
	defer l.mu {
	l.mu.Lock.Unlock()
	l.Entries = make([]*LeaderboardEntry, 0, l.maxSize)
	l.UpdatedAt = time.Now()
}

// GetTotalEntries 获取总条目数
func (l *Leaderboard) GetTotalEntries() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.Entries)
}

// LeaderboardManager 排行榜管理器
type LeaderboardManager struct {
	leaderboards map[LeaderboardType]*Leaderboard
	mu           sync.RWMutex
	friendRanks  map[string]map[LeaderboardType][]string // 好友排行榜缓存
}

// NewLeaderboardManager 创建排行榜管理器
func NewLeaderboardManager() *LeaderboardManager {
	m := &LeaderboardManager{
		leaderboards: make(map[LeaderboardType]*Leaderboard),
		friendRanks:  make(map[string]map[LeaderboardType][]string),
	}

	// 初始化各类排行榜
	m.leaderboards[LeaderboardTypeScore] = NewLeaderboard(LeaderboardTypeScore, "最强王者", 1000)
	m.leaderboards[LeaderboardTypeLevel] = NewLeaderboard(LeaderboardTypeLevel, "等级榜", 500)
	m.leaderboards[LeaderboardTypeCombat] = NewLeaderboard(LeaderboardTypeCombat, "战力榜", 500)
	m.leaderboards[LeaderboardTypeAchievement] = NewLeaderboard(LeaderboardTypeAchievement, "成就榜", 500)
	m.leaderboards[LeaderboardTypeStreak] = NewLeaderboard(LeaderboardTypeStreak, "连续登录榜", 200)
	m.leaderboards[LeaderboardTypeRichest] = NewLeaderboard(LeaderboardTypeRichest, "富豪榜", 200)
	m.leaderboards[LeaderboardTypeKills] = NewLeaderboard(LeaderboardTypeKills, "击杀榜", 500)
	m.leaderboards[LeaderboardTypeWins] = NewLeaderboard(LeaderboardTypeWins, "胜率榜", 500)

	return m
}

// GetLeaderboard 获取排行榜
func (m *LeaderboardManager) GetLeaderboard(ltype LeaderboardType) *Leaderboard {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.leaderboards[ltype]
}

// UpdateScore 更新分数
func (m *LeaderboardManager) UpdateScore(playerID, playerName, avatar string, score float64) {
	m.leaderboards[LeaderboardTypeScore].UpdateEntry(playerID, playerName, avatar, score)
}

// UpdateLevel 更新等级
func (m *LeaderboardManager) UpdateLevel(playerID, playerName, avatar string, level int) {
	m.leaderboards[LeaderboardTypeLevel].UpdateEntry(playerID, playerName, avatar, float64(level))
}

// UpdateCombatPower 更新战力
func (m *LeaderboardManager) UpdateCombatPower(playerID, playerName, avatar string, power float64) {
	m.leaderboards[LeaderboardTypeCombat].UpdateEntry(playerID, playerName, avatar, power)
}

// UpdateAchievementProgress 更新成就进度
func (m *LeaderboardManager) UpdateAchievementProgress(playerID, playerName, avatar string, count int) {
	m.leaderboards[LeaderboardTypeAchievement].UpdateEntry(playerID, playerName, avatar, float64(count))
}

// UpdateLoginStreak 更新连续登录
func (m *LeaderboardManager) UpdateLoginStreak(playerID, playerName, avatar string, days int) {
	m.leaderboards[LeaderboardTypeStreak].UpdateEntry(playerID, playerName, avatar, float64(days))
}

// UpdateRiches 更新财富
func (m *LeaderboardManager) UpdateRiches(playerID, playerName, avatar string, coins, gems float64) {
	// 财富值 = 金币/100 + 钻石*10
	wealth := coins/100 + gems*10
	m.leaderboards[LeaderboardTypeRichest].UpdateEntry(playerID, playerName, avatar, wealth)
}

// UpdateKills 更新击杀数
func (m *LeaderboardManager) UpdateKills(playerID, playerName, avatar string, kills int) {
	m.leaderboards[LeaderboardTypeKills].UpdateEntry(playerID, playerName, avatar, float64(kills))
}

// UpdateWinRate 更新胜率
func (m *LeaderboardManager) UpdateWinRate(playerID, playerName, avatar string, wins, losses int) {
	if wins+losses == 0 {
		return
	}
	rate := float64(wins) / float64(wins+losses) * 100
	m.leaderboards[LeaderboardTypeWins].UpdateEntry(playerID, playerName, avatar, rate)
}

// GetPlayerRank 获取玩家排名
func (m *LeaderboardManager) GetPlayerRank(playerID string, ltype LeaderboardType) (int, float64) {
	return m.leaderboards[ltype].GetRank(playerID)
}

// GetTopPlayers 获取top N玩家
func (m *LeaderboardManager) GetTopPlayers(ltype LeaderboardType, n int) []*LeaderboardEntry {
	return m.leaderboards[ltype].GetTopN(n)
}

// GetPlayerStats 获取玩家排行榜数据
func (m *LeaderboardManager) GetPlayerStats(playerID string) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]interface{})

	for ltype, lb := range m.leaderboards {
		rank, value := lb.GetRank(playerID)
		if rank > 0 {
			stats[ltype.String()] = map[string]interface{}{
				"rank":  rank,
				"value": value,
			}
		}
	}

	return stats
}

// GetAllLeaderboards 获取所有排行榜概览
func (m *LeaderboardManager) GetAllLeaderboards() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{})

	for ltype, lb := range m.leaderboards {
		result[ltype.String()] = map[string]interface{}{
			"name":         lb.Name,
			"total":        len(lb.Entries),
			"top3":         lb.GetTopN(3),
			"last_updated": lb.UpdatedAt,
		}
	}

	return result
}

// GetFriendLeaderboard 获取好友排行榜
func (m *LeaderboardManager) GetFriendLeaderboard(playerID string, friendIDs []string, ltype LeaderboardType) []*LeaderboardEntry {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 将玩家和好友ID合并
	allIDs := append([]string{playerID}, friendIDs...)

	// 获取所有相关玩家的数据
	entries := make([]*LeaderboardEntry, 0, len(allIDs))

	lb := m.leaderboards[ltype]
	for _, id := range allIDs {
		rank, value := lb.GetRank(id)
		if rank > 0 {
			for _, e := range lb.Entries {
				if e.PlayerID == id {
					entry := *e
					entry.Rank = rank
					entries = append(entries, &entry)
					break
				}
			}
		}
	}

	// 按排名排序
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Rank < entries[j].Rank
	})

	// 重新设置相对排名
	for i := range entries {
		entries[i].Rank = i + 1
		entries[i].Change = 0
	}

	return entries
}

// RefreshLeaderboard 刷新排行榜（从数据库重新加载）
func (m *LeaderboardManager) RefreshLeaderboard(ltype LeaderboardType) {
	// 这里应该从数据库加载数据
	// 简化实现，实际应该查询数据库
	fmt.Printf("Refreshing leaderboard: %s\n", ltype.String())
}

// GetSeasonLeaderboard 获取赛季排行榜
func (m *LeaderboardManager) GetSeasonLeaderboard(seasonID string, ltype LeaderboardType) *Leaderboard {
	// 赛季排行榜应该有独立的存储
	// 简化实现，返回主排行榜
	return m.leaderboards[ltype]
}

// ExportLeaderboard 导出排行榜数据
func (m *LeaderboardManager) ExportLeaderboard(ltype LeaderboardType) ([]byte, error) {
	lb := m.leaderboards[ltype]
	return json.Marshal(lb.Entries)
}

// ImportLeaderboard 导入排行榜数据
func (m *LeaderboardManager) ImportLeaderboard(ltype LeaderboardType, data []byte) error {
	var entries []*LeaderboardEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	m.leaderboards[ltype].mu.Lock()
	defer m.leaderboards[ltype].mu.Unlock()

	m.leaderboards[ltype].Entries = entries
	m.leaderboards[ltype].UpdatedAt = time.Now()

	return nil
}

// String 转换为字符串
func (lt LeaderboardType) String() string {
	switch lt {
	case LeaderboardTypeScore:
		return "score"
	case LeaderboardTypeLevel:
		return "level"
	case LeaderboardTypeCombat:
		return "combat"
	case LeaderboardTypeAchievement:
		return "achievement"
	case LeaderboardTypeStreak:
		return "streak"
	case LeaderboardTypeRichest:
		return "richest"
	case LeaderboardTypeKills:
		return "kills"
	case LeaderboardTypeWins:
		return "wins"
	default:
		return "unknown"
	}
}

// LeaderboardEvent 排行榜事件
type LeaderboardEvent struct {
	Type      LeaderboardType `json:"type"`
	PlayerID  string          `json:"player_id"`
	PlayerName string         `json:"player_name"`
	Avatar    string          `json:"avatar"`
	Value     float64         `json:"value"`
	Timestamp time.Time       `json:"timestamp"`
}

// LeaderboardEventHandler 排行榜事件处理器
type LeaderboardEventHandler struct {
	manager *LeaderboardManager
	eventCh chan *LeaderboardEvent
	stopCh  chan struct{}
}

// NewLeaderboardEventHandler 创建事件处理器
func NewLeaderboardEventHandler(m *LeaderboardManager) *LeaderboardEventHandler {
	return &LeaderboardEventHandler{
		manager: m,
		eventCh: make(chan *LeaderboardEvent, 1000),
		stopCh:  make(chan struct{}),
	}
}

// Start 启动事件处理
func (h *LeaderboardEventHandler) Start() {
	go func() {
		for {
			select {
			case event := <-h.eventCh:
				h.handleEvent(event)
			case <-h.stopCh:
				return
			}
		}
	}()
}

// Stop 停止事件处理
func (h *LeaderboardEventHandler) Stop() {
	close(h.stopCh)
}

// PushEvent 推送事件
func (h *LeaderboardEventHandler) PushEvent(event *LeaderboardEvent) {
	select {
	case h.eventCh <- event:
	default:
		fmt.Println("Leaderboard event channel full")
	}
}

// handleEvent 处理事件
func (h *LeaderboardEventHandler) handleEvent(event *LeaderboardEvent) {
	switch event.Type {
	case LeaderboardTypeScore:
		h.manager.UpdateScore(event.PlayerID, event.PlayerName, event.Avatar, event.Value)
	case LeaderboardTypeLevel:
		h.manager.UpdateLevel(event.PlayerID, event.PlayerName, event.Avatar, int(event.Value))
	case LeaderboardTypeCombat:
		h.manager.UpdateCombatPower(event.PlayerID, event.PlayerName, event.Avatar, event.Value)
	case LeaderboardTypeAchievement:
		h.manager.UpdateAchievementProgress(event.PlayerID, event.PlayerName, event.Avatar, int(event.Value))
	case LeaderboardTypeStreak:
		h.manager.UpdateLoginStreak(event.PlayerID, event.PlayerName, event.Avatar, int(event.Value))
	case LeaderboardTypeKills:
		h.manager.UpdateKills(event.PlayerID, event.PlayerName, event.Avatar, int(event.Value))
	case LeaderboardTypeWins:
		// 胜率需要额外处理 wins 和 losses
		break
	}
}
