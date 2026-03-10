package game

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// BattleResult 战斗结果
type BattleResult int

const (
	BattleWin  BattleResult = iota // 胜利
	BattleDraw                     // 平局
	BattleLoss                     // 失败
	BattleQuit                     // 退出
)

// BattleMode 战斗模式
type BattleMode int

const (
	BattleModeRanked   BattleMode = iota // 排位赛
	BattleModeFriendly                    // 友谊赛
	BattleModePractice                    // 练习赛
	BattleModeEvent                       // 活动赛
)

// BattleHistory 战斗历史记录
type BattleHistory struct {
	ID            string       `json:"id"`             // 记录ID
	PlayerID      string       `json:"player_id"`     // 玩家ID
	RoomID        string       `json:"room_id"`       // 房间ID
	BattleMode    BattleMode   `json:"battle_mode"`   // 战斗模式
	Result        BattleResult `json:"result"`        // 战斗结果
	Score         int          `json:"score"`         // 获得分数
	RankScore     int          `json:"rank_score"`    // 排位分数变化
	KillCount     int          `json:"kill_count"`    // 击杀数
	DeathCount    int          `json:"death_count"`   // 死亡数
	AssistCount   int          `json:"assist_count"`  // 助攻数
	DamageDealt   int          `json:"damage_dealt"`  // 造成伤害
	DamageTaken   int          `json:"damage_taken"`  // 承受伤害
	HighestCombo  int          `json:"highest_combo"` // 最高连击
	TimeSurvived  int          `json:"time_survived"` // 存活时间(秒)
	TimePlayed    int          `json:"time_played"`  // 战斗时长(秒)
	TowersBuilt   int          `json:"towers_built"` // 建造塔数
	TowersUpgraded int          `json:"towers_upgraded"` // 升级塔数
	EnemiesKilled map[string]int `json:"enemies_killed"` // 各类敌人击杀数
	SkillsUsed    int          `json:"skills_used"`  // 技能使用次数
	GiftsReceived int          `json:"gifts_received"` // 收到礼物数
	DanmakuCount  int          `json:"danmaku_count"` // 弹幕数
	StartTime     time.Time    `json:"start_time"`   // 开始时间
	EndTime       time.Time    `json:"end_time"`     // 结束时间
	Timestamp     time.Time    `json:"timestamp"`    // 记录时间
}

// PlayerBattleStats 玩家战斗统计数据
type PlayerBattleStats struct {
	PlayerID         string          `json:"player_id"`          // 玩家ID
	TotalBattles     int             `json:"total_battles"`      // 总战斗场次
	WinCount         int             `json:"win_count"`          // 胜利场次
	DrawCount        int             `json:"draw_count"`         // 平局场次
	LossCount        int             `json:"loss_count"`         // 失败场次
	TotalKills       int             `json:"total_kills"`        // 总击杀数
	TotalDeaths      int             `json:"total_deaths"`       // 总死亡数
	TotalAssists     int             `json:"total_assists"`      // 总助攻数
	TotalDamageDealt int             `json:"total_damage_dealt"` // 总造成伤害
	TotalDamageTaken int             `json:"total_damage_taken"` // 总承受伤害
	HighestCombo     int             `json:"highest_combo"`      // 最高连击
	TotalPlayTime    int             `json:"total_play_time"`    // 总游戏时长(秒)
	CurrentStreak    int             `json:"current_streak"`     // 当前连胜/连败
	MaxWinStreak     int             `json:"max_win_streak"`     // 最大连胜
	MaxLossStreak    int             `json:"max_loss_streak"`    // 最大连败
	AverageScore     float64         `json:"average_score"`      // 平均得分
	TotalScore       int             `json:"total_score"`        // 总得分
	TotalRankScore   int             `json:"total_rank_score"`   // 总排位分数
	CurrentRankScore int             `json:"current_rank_score"` // 当前排位分数
	RankTier         string          `json:"rank_tier"`          // 段位
	MostUsedTowers   []string        `json:"most_used_towers"`   // 最常用塔
	MostUsedSkills  []string        `json:"most_used_skills"`   // 最常用技能
	ModeStats        map[string]BattleModeStats `json:"mode_stats"` // 各模式统计
	RecentBattles    []string        `json:"recent_battles"`    // 最近战斗ID列表
	LastBattleTime   time.Time       `json:"last_battle_time"`  // 最后战斗时间
	UpdatedAt        time.Time       `json:"updated_at"`        // 更新时间
}

// BattleModeStats 各战斗模式统计
type BattleModeStats struct {
	Battles   int     `json:"battles"`   // 场次
	Wins      int     `json:"wins"`      // 胜利
	WinRate   float64 `json:"win_rate"`  // 胜率
	AvgScore  float64 `json:"avg_score"` // 平均得分
	AvgKills  float64 `json:"avg_kills"` // 平均击杀
	AvgDeaths float64 `json:"avg_deaths"` // 平均死亡
}

// BattleHistoryManager 战绩管理器
type BattleHistoryManager struct {
	mu sync.RWMutex

	// 战斗历史: battleID -> BattleHistory
	battles map[string]*BattleHistory

	// 玩家战斗统计: playerID -> PlayerBattleStats
	stats map[string]*PlayerBattleStats

	// 玩家最近战斗: playerID -> []battleID
	recentBattles map[string][]string

	// 战斗ID生成
	battleIDCounter int64

	// 最大保留历史记录数
	maxHistoryPerPlayer int
}

// NewBattleHistoryManager 创建战绩管理器
func NewBattleHistoryManager() *BattleHistoryManager {
	return &BattleHistoryManager{
		battles:              make(map[string]*BattleHistory),
		stats:                make(map[string]*PlayerBattleStats),
		recentBattles:        make(map[string][]string),
		battleIDCounter:      time.Now().UnixNano(),
		maxHistoryPerPlayer: 1000,
	}
}

// CreateBattleRecord 创建战斗记录
func (bm *BattleHistoryManager) CreateBattleRecord(playerID, roomID string, mode BattleMode, startTime time.Time) string {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	bm.battleIDCounter++
	battleID := fmt.Sprintf("battle_%d", bm.battleIDCounter)

	battle := &BattleHistory{
		ID:           battleID,
		PlayerID:     playerID,
		RoomID:       roomID,
		BattleMode:   mode,
		StartTime:    startTime,
		Timestamp:    time.Now(),
		EnemiesKilled: make(map[string]int),
	}

	bm.battles[battleID] = battle

	// 添加到玩家最近战斗
	bm.recentBattles[playerID] = append([]string{battleID}, bm.recentBattles[playerID]...)
	if len(bm.recentBattles[playerID]) > bm.maxHistoryPerPlayer {
		bm.recentBattles[playerID] = bm.recentBattles[playerID][:bm.maxHistoryPerPlayer]
	}

	return battleID
}

// UpdateBattleRecord 更新战斗记录
func (bm *BattleHistoryManager) UpdateBattleRecord(battleID string, result BattleResult, score, rankScore int) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	battle, ok := bm.battles[battleID]
	if !ok {
		return fmt.Errorf("战斗记录不存在")
	}

	battle.Result = result
	battle.Score = score
	battle.RankScore = rankScore
	battle.EndTime = time.Now()
	battle.TimePlayed = int(battle.EndTime.Sub(battle.StartTime).Seconds())

	return nil
}

// UpdateBattleStats 更新战斗统计数据
func (bm *BattleHistoryManager) UpdateBattleStats(battleID string, killCount, deathCount, assistCount, damageDealt, damageTaken, highestCombo, timeSurvived int) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	battle, ok := bm.battles[battleID]
	if !ok {
		return fmt.Errorf("战斗记录不存在")
	}

	battle.KillCount = killCount
	battle.DeathCount = deathCount
	battle.AssistCount = assistCount
	battle.DamageDealt = damageDealt
	battle.DamageTaken = damageTaken
	battle.HighestCombo = highestCombo
	battle.TimeSurvived = timeSurvived

	return nil
}

// UpdateBattleActions 更新战斗行为数据
func (bm *BattleHistoryManager) UpdateBattleActions(battleID string, towersBuilt, towersUpgraded, skillsUsed, giftsReceived, danmakuCount int) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	battle, ok := bm.battles[battleID]
	if !ok {
		return fmt.Errorf("战斗记录不存在")
	}

	battle.TowersBuilt = towersBuilt
	battle.TowersUpgraded = towersUpgraded
	battle.SkillsUsed = skillsUsed
	battle.GiftsReceived = giftsReceived
	battle.DanmakuCount = danmakuCount

	return nil
}

// RecordEnemyKill 记录敌人击杀
func (bm *BattleHistoryManager) RecordEnemyKill(battleID, enemyType string) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	battle, ok := bm.battles[battleID]
	if !ok {
		return fmt.Errorf("战斗记录不存在")
	}

	battle.EnemiesKilled[enemyType]++

	return nil
}

// FinishBattle 完成战斗并更新统计
func (bm *BattleHistoryManager) FinishBattle(battleID string) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	battle, ok := bm.battles[battleID]
	if !ok {
		return fmt.Errorf("战斗记录不存在")
	}

	playerID := battle.PlayerID

	// 获取或创建玩家统计
	stats := bm.stats[playerID]
	if stats == nil {
		stats = &PlayerBattleStats{
			PlayerID:   playerID,
			ModeStats:  make(map[string]BattleModeStats),
			RecentBattles: []string{},
		}
		bm.stats[playerID] = stats
	}

	// 更新总战斗数
	stats.TotalBattles++

	// 更新胜负
	switch battle.Result {
	case BattleWin:
		stats.WinCount++
		stats.CurrentStreak++
		if stats.CurrentStreak > stats.MaxWinStreak {
			stats.MaxWinStreak = stats.CurrentStreak
		}
	case BattleDraw:
		stats.DrawCount++
		stats.CurrentStreak = 0
	case BattleLoss:
		stats.LossCount++
		stats.CurrentStreak--
		if stats.CurrentStreak < stats.MaxLossStreak {
			stats.MaxLossStreak = stats.CurrentStreak
		}
	case BattleQuit:
		stats.LossCount++
		stats.CurrentStreak = 0
	}

	// 更新击杀/死亡/助攻
	stats.TotalKills += battle.KillCount
	stats.TotalDeaths += battle.DeathCount
	stats.TotalAssists += battle.AssistCount

	// 更新伤害
	stats.TotalDamageDealt += battle.DamageDealt
	stats.TotalDamageTaken += battle.DamageTaken

	// 更新最高连击
	if battle.HighestCombo > stats.HighestCombo {
		stats.HighestCombo = battle.HighestCombo
	}

	// 更新总游戏时长
	stats.TotalPlayTime += battle.TimePlayed

	// 更新分数
	stats.TotalScore += battle.Score
	stats.TotalRankScore += battle.RankScore
	stats.CurrentRankScore += battle.RankScore

	// 计算平均得分
	if stats.TotalBattles > 0 {
		stats.AverageScore = float64(stats.TotalScore) / float64(stats.TotalBattles)
	}

	// 更新段位
	stats.RankTier = bm.CalculateRankTier(stats.CurrentRankScore)

	// 更新最近战斗
	stats.RecentBattles = append([]string{battleID}, stats.RecentBattles...)
	if len(stats.RecentBattles) > 20 {
		stats.RecentBattles = stats.RecentBattles[:20]
	}

	// 更新模式统计
	modeKey := fmt.Sprintf("%d", battle.BattleMode)
	modeStats := stats.ModeStats[modeKey]
	modeStats.Battles++
	if battle.Result == BattleWin {
		modeStats.Wins++
	}
	modeStats.WinRate = float64(modeStats.Wins) / float64(modeStats.Battles)
	modeStats.AvgScore = (modeStats.AvgScore*float64(modeStats.Battles-1) + float64(battle.Score)) / float64(modeStats.Battles)
	modeStats.AvgKills = (modeStats.AvgKills*float64(modeStats.Battles-1) + float64(battle.KillCount)) / float64(modeStats.Battles)
	modeStats.AvgDeaths = (modeStats.AvgDeaths*float64(modeStats.Battles-1) + float64(battle.DeathCount)) / float64(modeStats.Battles)
	stats.ModeStats[modeKey] = modeStats

	// 更新时间
	stats.LastBattleTime = time.Now()
	stats.UpdatedAt = time.Now()

	return nil
}

// CalculateRankTier 根据排位分计算段位
func (bm *BattleHistoryManager) CalculateRankTier(rankScore int) string {
	switch {
	case rankScore >= 3000:
		return "王者"
	case rankScore >= 2500:
		return "宗师"
	case rankScore >= 2000:
		return "大师"
	case rankScore >= 1500:
		return "钻石"
	case rankScore >= 1000:
		return "铂金"
	case rankScore >= 500:
		return "黄金"
	case rankScore >= 200:
		return "白银"
	default:
		return "青铜"
	}
}

// GetBattleHistory 获取战斗记录
func (bm *BattleHistoryManager) GetBattleHistory(battleID string) (*BattleHistory, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	battle, ok := bm.battles[battleID]
	if !ok {
		return nil, fmt.Errorf("战斗记录不存在")
	}

	return battle, nil
}

// GetPlayerBattleHistory 获取玩家战斗历史
func (bm *BattleHistoryManager) GetPlayerBattleHistory(playerID string, limit int) []*BattleHistory {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	battleIDs := bm.recentBattles[playerID]
	if len(battleIDs) == 0 {
		return []*BattleHistory{}
	}

	if limit <= 0 || limit > len(battleIDs) {
		limit = len(battleIDs)
	}

	result := make([]*BattleHistory, 0, limit)
	for i := 0; i < limit; i++ {
		if battle, ok := bm.battles[battleIDs[i]]; ok {
			result = append(result, battle)
		}
	}

	return result
}

// GetPlayerStats 获取玩家战斗统计
func (bm *BattleHistoryManager) GetPlayerStats(playerID string) (*PlayerBattleStats, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	stats, ok := bm.stats[playerID]
	if !ok {
		return nil, fmt.Errorf("玩家统计不存在")
	}

	return stats, nil
}

// GetPlayerWinRate 获取玩家胜率
func (bm *BattleHistoryManager) GetPlayerWinRate(playerID string) float64 {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	stats, ok := bm.stats[playerID]
	if !ok || stats.TotalBattles == 0 {
		return 0.0
	}

	return float64(stats.WinCount) / float64(stats.TotalBattles) * 100
}

// GetPlayerKDA 获取玩家KDA
func (bm *BattleHistoryManager) GetPlayerKDA(playerID string) float64 {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	stats, ok := bm.stats[playerID]
	if !ok || stats.TotalDeaths == 0 {
		return 0.0
	}

	return float64(stats.TotalKills+stats.TotalAssists) / float64(stats.TotalDeaths)
}

// GetRecentBattles 获取最近战斗数
func (bm *BattleHistoryManager) GetRecentBattles(playerID string, count int) int {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if count <= 0 {
		count = 10
	}

	battleIDs := bm.recentBattles[playerID]
	if len(battleIDs) > count {
		return count
	}

	return len(battleIDs)
}

// GetLeaderboard 获取排行榜
func (bm *BattleHistoryManager) GetLeaderboard(sortBy string, limit int) []*PlayerBattleStats {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	allStats := make([]*PlayerBattleStats, 0, len(bm.stats))
	for _, stats := range bm.stats {
		allStats = append(allStats, stats)
	}

	// 排序
	switch sortBy {
	case "rank_score":
		sort.Slice(allStats, func(i, j int) bool {
			return allStats[i].CurrentRankScore > allStats[j].CurrentRankScore
		})
	case "win_rate":
		sort.Slice(allStats, func(i, j int) bool {
			return allStats[i].WinCount/allStats[i].TotalBattles > allStats[j].WinCount/allStats[j].TotalBattles
		})
	case "kills":
		sort.Slice(allStats, func(i, j int) bool {
			return allStats[i].TotalKills > allStats[j].TotalKills
		})
	case "wins":
		sort.Slice(allStats, func(i, j int) bool {
			return allStats[i].WinCount > allStats[j].WinCount
		})
	default:
		sort.Slice(allStats, func(i, j int) bool {
			return allStats[i].TotalBattles > allStats[j].TotalBattles
		})
	}

	if limit > len(allStats) {
		limit = len(allStats)
	}

	return allStats[:limit]
}

// CleanupOldRecords 清理旧记录
func (bm *BattleHistoryManager) CleanupOldRecords(maxAge time.Duration) int {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	deletedCount := 0

	for battleID, battle := range bm.battles {
		if battle.Timestamp.Before(cutoff) {
			delete(bm.battles, battleID)
			deletedCount++
		}
	}

	return deletedCount
}
