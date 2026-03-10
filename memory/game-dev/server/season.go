package game

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

// SeasonSystem 赛季系统 - 管理游戏赛季、排行榜和奖励
type SeasonSystem struct {
	mu              sync.RWMutex
	seasons         map[string]*Season
	currentSeasonID string
	rewardPools     map[string]*RewardPool
	seasonRewards   map[string]map[string][]SeasonReward // seasonID -> rank -> rewards
	playerSeasons   map[string]*PlayerSeasonData          // playerID -> season data
	rankings        map[string]RankingCache               // seasonID -> ranking cache
}

// Season 赛季
type Season struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Status        string    `json:"status"` // upcoming/active/ended
	SeasonType    string    `json:"season_type"` // regular/champion/special
	Rules         SeasonRules `json:"rules"`
	ResetProgress bool      `json:"reset_progress"`
	CreatedAt     time.Time `json:"created_at"`
}

// SeasonRules 赛季规则
type SeasonRules struct {
	MaxRank           int     `json:"max_rank"`           // 最大段位
	PointsPerWin      int     `json:"points_per_win"`     // 胜利获得积分
	PointsPerLoss     int     `json:"points_per_loss"`    // 失败扣除积分
	PromotionPoints   int     `json:"promotion_points"`   // 晋升所需积分
	DemotionPoints    int     `json:"demotion_points"`    // 降级积分线
	PlacementMatches  int     `json:"placement_matches"`  // 定级赛场数
	MinPoints         int     `json:"min_points"`         // 最低积分
	BonusWinStreak    int     `json:"bonus_win_streak"`   // 连胜奖励场次
	BonusPointsFactor float64 `json:"bonus_points_factor"` // 连胜加成系数
}

// RewardPool 奖励池
type RewardPool struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	SeasonTypes []string      `json:"season_types"` // 适用的赛季类型
	Tiers       []RewardTier  `json:"tiers"`
	CreatedAt   time.Time     `json:"created_at"`
}

// RewardTier 奖励等级
type RewardTier struct {
	TierID     string  `json:"tier_id"`
	Name       string  `json:"name"`
	MinRank    int     `json:"min_rank"`
	MaxRank    int     `json:"max_rank"`
	Points     int     `json:"points"`      // 积分要求
	Rewards    []Item `json:"rewards"`      // 奖励列表
	IsExclusive bool  `json:"is_exclusive"` // 独家奖励
}

// SeasonReward 赛季奖励
type SeasonReward struct {
	RewardID  string `json:"reward_id"`
	ItemID    string `json:"item_id"`
	ItemType  string `json:"item_type"`
	Amount    int    `json:"amount"`
	Rarity    string `json:"rarity"` // common/rare/epic/legendary
	IsUnique  bool   `json:"is_unique"`
}

// PlayerSeasonData 玩家赛季数据
type PlayerSeasonData struct {
	PlayerID       string            `json:"player_id"`
	SeasonID       string            `json:"season_id"`
	CurrentRank    int               `json:"current_rank"`    // 当前段位 (1=青铜, 2=白银, 3=黄金, 4=钻石, 5=王者)
	Points         int               `json:"points"`         // 当前积分
	TotalWins      int               `json:"total_wins"`     // 总胜场
	TotalLosses    int               `json:"total_losses"`   // 总负场
	WinStreak      int               `json:"win_streak"`     // 当前连胜
	BestStreak     int               `json:"best_streak"`    // 最高连胜
	PlacementDone  int               `json:"placement_done"` // 定级赛完成场数
	PlacementWins  int               `json:"placement_wins"` // 定级赛胜场
	LastMatchTime  time.Time         `json:"last_match_time"`
	ReceivedRewards map[string]bool  `json:"received_rewards"` // 已领取奖励
	History        []SeasonHistory   `json:"history"`         // 历史记录
	UpdatedAt      time.Time         `json:"updated_at"`
}

// SeasonHistory 赛季历史
type SeasonHistory struct {
	SeasonID   string    `json:"season_id"`
	FinalRank  int       `json:"final_rank"`
	FinalPoints int      `json:"final_points"`
	TotalWins  int       `json:"total_wins"`
	BestStreak int       `json:"best_streak"`
	EndedAt    time.Time `json:"ended_at"`
}

// RankingCache 排行榜缓存
type RankingCache struct {
	SeasonID   string           `json:"season_id"`
	Rankings   []PlayerRank     `json:"rankings"`
	UpdatedAt  time.Time        `json:"updated_at"`
	TotalCount int              `json:"total_count"`
}

// PlayerRank 玩家排名
type PlayerRank struct {
	Rank       int    `json:"rank"`
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Avatar     string `json:"avatar"`
	RankLevel  int    `json:"rank_level"` // 段位等级
	Points     int    `json:"points"`
	Wins       int    `json:"wins"`
	WinRate    string `json:"win_rate"`
}

// RankInfo 段位信息
type RankInfo struct {
	Level        int    `json:"level"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	MinPoints    int    `json:"min_points"`
	MaxPoints    int    `json:"max_points"`
	Promotion    string `json:"promotion"`
}

// NewSeasonSystem 创建赛季系统
func NewSeasonSystem() *SeasonSystem {
	ss := &SeasonSystem{
		seasons:         make(map[string]*Season),
		rewardPools:     make(map[string]*RewardPool),
		seasonRewards:   make(map[string]map[string][]SeasonReward),
		playerSeasons:   make(map[string]*PlayerSeasonData),
		rankings:        make(map[string]RankingCache),
	}
	ss.initDefaultRanks()
	ss.initDefaultSeasons()
	return ss
}

// initDefaultRanks 初始化默认段位
func (ss *SeasonSystem) initDefaultRanks() {
	// 默认段位在内存中定义
}

// initDefaultSeasons 初始化默认赛季
func (ss *SeasonSystem) initDefaultSeasons() {
	now := time.Now()
	season := &Season{
		ID:          "season_1",
		Name:        "S1 恐龙时代",
		StartTime:   now,
		EndTime:     now.Add(30 * 24 * time.Hour),
		Status:      "active",
		SeasonType:  "regular",
		Rules: SeasonRules{
			MaxRank:            5,
			PointsPerWin:       25,
			PointsPerLoss:      -15,
			PromotionPoints:    100,
			DemotionPoints:     0,
			PlacementMatches:   10,
			MinPoints:          0,
			BonusWinStreak:     3,
			BonusPointsFactor:  1.5,
		},
		ResetProgress: false,
		CreatedAt:     now,
	}
	ss.seasons[season.ID] = season
	ss.currentSeasonID = season.ID
}

// CreateSeason 创建赛季
func (ss *SeasonSystem) CreateSeason(name, seasonType string, startTime, endTime time.Time, rules SeasonRules) (*Season, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if endTime.Before(startTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	seasonID := fmt.Sprintf("season_%d", len(ss.seasons)+1)
	season := &Season{
		ID:          seasonID,
		Name:        name,
		StartTime:   startTime,
		EndTime:     endTime,
		Status:      "upcoming",
		SeasonType:  seasonType,
		Rules:       rules,
		ResetProgress: false,
		CreatedAt:   time.Now(),
	}

	ss.seasons[seasonID] = season
	return season, nil
}

// GetCurrentSeason 获取当前赛季
func (ss *SeasonSystem) GetCurrentSeason() *Season {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.seasons[ss.currentSeasonID]
}

// GetSeason 获取赛季
func (ss *SeasonSystem) GetSeason(seasonID string) *Season {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.seasons[seasonID]
}

// UpdateSeasonStatus 更新赛季状态
func (ss *SeasonSystem) UpdateSeasonStatus(seasonID string, status string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	season, ok := ss.seasons[seasonID]
	if !ok {
		return fmt.Errorf("season not found: %s", seasonID)
	}

	if status != "upcoming" && status != "active" && status != "ended" {
		return fmt.Errorf("invalid status: %s", status)
	}

	season.Status = status
	return nil
}

// CheckAndUpdateSeasons 检查并更新赛季状态
func (ss *SeasonSystem) CheckAndUpdateSeasons() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	now := time.Now()
	for _, season := range ss.seasons {
		if season.Status == "upcoming" && now.After(season.StartTime) {
			season.Status = "active"
			ss.currentSeasonID = season.ID
		} else if season.Status == "active" && now.After(season.EndTime) {
			season.Status = "ended"
			// 结算赛季奖励
			ss.settleSeason(season.ID)
		}
	}
}

// settleSeason 结算赛季
func (ss *SeasonSystem) settleSeason(seasonID string) {
	for playerID, data := range ss.playerSeasons {
		if data.SeasonID == seasonID {
			// 记录历史
			history := SeasonHistory{
				SeasonID:   seasonID,
				FinalRank:  data.CurrentRank,
				FinalPoints: data.Points,
				TotalWins:  data.TotalWins,
				BestStreak: data.BestStreak,
				EndedAt:    time.Now(),
			}
			data.History = append(data.History, history)
		}
	}
}

// GetOrCreatePlayerSeason 获取或创建玩家赛季数据
func (ss *SeasonSystem) GetOrCreatePlayerSeason(playerID, seasonID string) *PlayerSeasonData {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	key := fmt.Sprintf("%s_%s", playerID, seasonID)
	if data, ok := ss.playerSeasons[key]; ok {
		return data
	}

	data := &PlayerSeasonData{
		PlayerID:        playerID,
		SeasonID:        seasonID,
		CurrentRank:     1,
		Points:          0,
		TotalWins:       0,
		TotalLosses:     0,
		WinStreak:       0,
		BestStreak:      0,
		PlacementDone:   0,
		PlacementWins:   0,
		LastMatchTime:   time.Now(),
		ReceivedRewards: make(map[string]bool),
		History:         []SeasonHistory{},
		UpdatedAt:       time.Now(),
	}

	ss.playerSeasons[key] = data
	return data
}

// UpdateMatchResult 更新比赛结果
func (ss *SeasonSystem) UpdateMatchResult(playerID string, isWin bool) (*PlayerSeasonData, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	currentSeason := ss.seasons[ss.currentSeasonID]
	if currentSeason == nil {
		return nil, fmt.Errorf("no active season")
	}

	key := fmt.Sprintf("%s_%s", playerID, currentSeason.ID)
	data, ok := ss.playerSeasons[key]
	if !ok {
		data = ss.GetOrCreatePlayerSeason(playerID, currentSeason.ID)
	}

	rules := currentSeason.Rules

	// 计算积分变化
	pointsChange := rules.PointsPerLoss
	if isWin {
		pointsChange = rules.PointsPerWin
		// 连胜加成
		if data.WinStreak >= rules.BonusWinStreak {
			bonus := float64(pointsChange) * (rules.BonusPointsFactor - 1)
			pointsChange = int(math.Round(float64(pointsChange) + bonus))
		}
		data.WinStreak++
		if data.WinStreak > data.BestStreak {
			data.BestStreak = data.WinStreak
		}
		data.TotalWins++
	} else {
		data.WinStreak = 0
		data.TotalLosses++
	}

	// 更新定级赛
	if data.PlacementDone < rules.PlacementMatches {
		data.PlacementDone++
		if isWin {
			data.PlacementWins++
		}
		// 定级赛期间不计算积分
		data.LastMatchTime = time.Now()
		data.UpdatedAt = time.Now()
		return data, nil
	}

	// 更新积分
	data.Points = ss.calculateNewPoints(data.Points, pointsChange, rules)

	// 更新段位
	data.CurrentRank = ss.calculateNewRank(data.Points, rules)

	data.LastMatchTime = time.Now()
	data.UpdatedAt = time.Now()

	// 更新排行榜缓存
	ss.updateRankingCache(currentSeason.ID)

	return data, nil
}

// calculateNewPoints 计算新积分
func (ss *SeasonSystem) calculateNewPoints(currentPoints, change int, rules SeasonRules) int {
	newPoints := currentPoints + change
	if newPoints < rules.MinPoints {
		return rules.MinPoints
	}
	// 积分上限可以设置一个合理的最大值
	maxPoints := rules.PromotionPoints * rules.MaxRank
	if newPoints > maxPoints {
		return maxPoints
	}
	return newPoints
}

// calculateNewRank 计算新段位
func (ss *SeasonSystem) calculateNewRank(points int, rules SeasonRules) int {
	// 根据积分计算段位
	rank := (points / rules.PromotionPoints) + 1
	if rank > rules.MaxRank {
		return rules.MaxRank
	}
	if rank < 1 {
		return 1
	}
	return rank
}

// GetPlayerRank 获取玩家段位信息
func (ss *SeasonSystem) GetPlayerRank(playerID string) (*PlayerSeasonData, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	currentSeason := ss.seasons[ss.currentSeasonID]
	if currentSeason == nil {
		return nil, fmt.Errorf("no active season")
	}

	key := fmt.Sprintf("%s_%s", playerID, currentSeason.ID)
	data, ok := ss.playerSeasons[key]
	if !ok {
		return nil, fmt.Errorf("player season data not found")
	}

	return data, nil
}

// GetLeaderboard 获取排行榜
func (ss *SeasonSystem) GetLeaderboard(seasonID string, limit int) ([]PlayerRank, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	season, ok := ss.seasons[seasonID]
	if !ok {
		return nil, fmt.Errorf("season not found: %s", seasonID)
	}

	// 从缓存获取
	if cache, ok := ss.rankings[seasonID]; ok {
		if limit > len(cache.Rankings) {
			limit = len(cache.Rankings)
		}
		return cache.Rankings[:limit], nil
	}

	// 实时计算
	var rankings []PlayerRank
	for key, data := range ss.playerSeasons {
		if data.SeasonID != seasonID {
			continue
		}
		playerID := strings.Split(key, "_")[0]
		winRate := "0%"
		if data.TotalWins+data.TotalLosses > 0 {
			rate := float64(data.TotalWins) / float64(data.TotalWins+data.TotalLosses) * 100
			winRate = fmt.Sprintf("%.1f%%", rate)
		}
		rankings = append(rankings, PlayerRank{
			PlayerID:   playerID,
			RankLevel:  data.CurrentRank,
			Points:     data.Points,
			Wins:       data.TotalWins,
			WinRate:    winRate,
		})
	}

	// 排序
	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].Points != rankings[j].Points {
			return rankings[i].Points > rankings[j].Points
		}
		return rankings[i].Wins > rankings[j].Wins
	})

	// 设置排名
	for i := range rankings {
		rankings[i].Rank = i + 1
	}

	if limit > len(rankings) {
		limit = len(rankings)
	}

	return rankings[:limit], nil
}

// updateRankingCache 更新排行榜缓存
func (ss *SeasonSystem) updateRankingCache(seasonID string) {
	ranks, _ := ss.GetLeaderboard(seasonID, 1000)
	ss.rankings[seasonID] = RankingCache{
		SeasonID:   seasonID,
		Rankings:   ranks,
		UpdatedAt:  time.Now(),
		TotalCount: len(ranks),
	}
}

// CreateRewardPool 创建奖励池
func (ss *SeasonSystem) CreateRewardPool(name string, seasonTypes []string, tiers []RewardTier) *RewardPool {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	pool := &RewardPool{
		ID:          fmt.Sprintf("reward_pool_%d", len(ss.rewardPools)+1),
		Name:        name,
		SeasonTypes: seasonTypes,
		Tiers:       tiers,
		CreatedAt:   time.Now(),
	}

	ss.rewardPools[pool.ID] = pool
	return pool
}

// GetSeasonRewards 获取赛季奖励
func (ss *SeasonSystem) GetSeasonRewards(seasonID, rank string) []SeasonReward {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	if rewards, ok := ss.seasonRewards[seasonID]; ok {
		if tierRewards, ok := rewards[rank]; ok {
			return tierRewards
		}
	}
	return nil
}

// ClaimReward 领取奖励
func (ss *SeasonSystem) ClaimReward(playerID, rewardID string) (bool, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	currentSeason := ss.seasons[ss.currentSeasonID]
	if currentSeason == nil {
		return false, fmt.Errorf("no active season")
	}

	key := fmt.Sprintf("%s_%s", playerID, currentSeason.ID)
	data, ok := ss.playerSeasons[key]
	if !ok {
		return false, fmt.Errorf("player season data not found")
	}

	// 检查是否已领取
	if data.ReceivedRewards[rewardID] {
		return false, fmt.Errorf("reward already claimed")
	}

	// 标记已领取
	data.ReceivedRewards[rewardID] = true
	return true, nil
}

// GetRankInfo 获取段位信息
func (ss *SeasonSystem) GetRankInfo(rank int) *RankInfo {
	ranks := []*RankInfo{
		{Level: 1, Name: "青铜", Icon: "bronze", MinPoints: 0, MaxPoints: 99, Promotion: "白银"},
		{Level: 2, Name: "白银", Icon: "silver", MinPoints: 100, MaxPoints: 199, Promotion: "黄金"},
		{Level: 3, Name: "黄金", Icon: "gold", MinPoints: 200, MaxPoints: 399, Promotion: "钻石"},
		{Level: 4, Name: "钻石", Icon: "diamond", MinPoints: 400, MaxPoints: 799, Promotion: "王者"},
		{Level: 5, Name: "王者", Icon: "king", MinPoints: 800, MaxPoints: 999999, Promotion: ""},
	}

	if rank < 1 || rank > len(ranks) {
		return ranks[0]
	}
	return ranks[rank-1]
}

// GetSeasonStats 获取赛季统计
func (ss *SeasonSystem) GetSeasonStats(seasonID string) map[string]interface{} {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	stats := make(map[string]interface{})
	var totalWins, totalLosses int
	var activePlayers int

	for key, data := range ss.playerSeasons {
		if data.SeasonID != seasonID {
			continue
		}
		if strings.HasPrefix(key, seasonID) {
			activePlayers++
			totalWins += data.TotalWins
			totalLosses += data.TotalLosses
		}
	}

	stats["active_players"] = activePlayers
	stats["total_wins"] = totalWins
	stats["total_losses"] = totalLosses
	stats["total_matches"] = totalWins + totalLosses

	if totalWins+totalLosses > 0 {
		stats["global_win_rate"] = float64(totalWins) / float64(totalWins+totalLosses)
	}

	return stats
}

// MarshalJSON 序列化
func (ss *SeasonSystem) MarshalJSON() ([]byte, error) {
	type Alias SeasonSystem
	return json.Marshal(&struct {
		Alias
		CurrentSeason *Season `json:"current_season"`
	}{
		Alias:         Alias(*ss),
		CurrentSeason: ss.seasons[ss.currentSeasonID],
	})
}
