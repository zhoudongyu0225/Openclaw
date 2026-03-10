package game

import (
	"fmt"
	"time"
)

// ============================================================
// 成就徽章系统 Achievement Badge System
// ============================================================

// BadgeType 徽章类型
type BadgeType int

const (
	BadgeTypeCombat BadgeType = iota // 战斗徽章
	BadgeTypeCollection              // 收集徽章
	BadgeTypeSocial                  // 社交徽章
	BadgeTypeEvent                   // 活动徽章
	BadgeTypeSeason                  // 赛季徽章
	BadgeTypeSpecial                 // 特殊徽章
)

// BadgeRarity 徽章稀有度
type BadgeRarity int

const (
	BadgeRarityCommon BadgeRarity = iota // 普通
	BadgeRarityUncommon                   // 优秀
	BadgeRarityRare                       // 稀有
	BadgeRarityEpic                       // 史诗
	BadgeRarityLegendary                  // 传说
)

// Badge 徽章定义
type Badge struct {
	ID          string      `json:"id"`           // 徽章ID
	Name        string      `json:"name"`         // 徽章名称
	Description string      `json:"description"`  // 描述
	Type        BadgeType   `json:"type"`         // 类型
	Rarity      BadgeRarity `json:"rarity"`        // 稀有度
	Icon        string      `json:"icon"`          // 图标
	Condition   string      `json:"condition"`    // 获得条件
	StatID      string      `json:"stat_id"`       // 关联统计ID
	TargetValue int64       `json:"target_value"` // 目标值
	RewardCoins int64       `json:"reward_coins"`  // 金币奖励
	RewardGems  int64       `json:"reward_gems"`   // 钻石奖励
	RewardItem  string      `json:"reward_item"`   // 道具奖励
	Sort        int         `json:"sort"`          // 排序
	IsLimited   bool        `json:"is_limited"`    // 是否限时
	StartTime   int64       `json:"start_time"`    // 开始时间
	EndTime     int64       `json:"end_time"`      // 结束时间
}

// PlayerBadge 玩家徽章
type PlayerBadge struct {
	PlayerID   string `json:"player_id"`   // 玩家ID
	BadgeID    string `json:"badge_id"`     // 徽章ID
	ObtainedAt int64  `json:"obtained_at"`  // 获得时间
	IsNew      bool   `json:"is_new"`       // 是否新获得
	StatValue  int64  `json:"stat_value"`   // 当前统计值
}

// BadgeManager 徽章管理器
type BadgeManager struct {
	badges       map[string]*Badge           // badgeID -> Badge
	playerBadges map[string]map[string]*PlayerBadge // playerID -> badgeID -> PlayerBadge
	statBadges   map[string][]string         // statID -> []badgeID
}

// NewBadgeManager 创建徽章管理器
func NewBadgeManager() *BadgeManager {
	bm := &BadgeManager{
		badges:       make(map[string]*Badge),
		playerBadges: make(map[string]map[string]*PlayerBadge),
		statBadges:   make(map[string][]string),
	}
	bm.initBadges()
	return bm
}

// initBadges 初始化徽章
func (bm *BadgeManager) initBadges() {
	badges := []*Badge{
		// 战斗徽章
		{
			ID: "badge_combat_001", Name: "初战告捷", Description: "完成第一场战斗",
			Type: BadgeTypeCombat, Rarity: BadgeRarityCommon,
			StatID: "total_battles", TargetValue: 1,
			RewardCoins: 100,
		},
		{
			ID: "badge_combat_002", Name: "百人斩", Description: "累计击杀100个敌人",
			Type: BadgeTypeCombat, Rarity: BadgeRarityRare,
			StatID: "total_kills", TargetValue: 100,
			RewardCoins: 1000, RewardGems: 10,
		},
		{
			ID: "badge_combat_003", Name: "千人斩", Description: "累计击杀1000个敌人",
			Type: BadgeTypeCombat, Rarity: BadgeRarityEpic,
			StatID: "total_kills", TargetValue: 1000,
			RewardCoins: 5000, RewardGems: 50,
		},
		{
			ID: "badge_combat_004", Name: "Boss杀手", Description: "击杀10个Boss",
			Type: BadgeTypeCombat, Rarity: BadgeRarityEpic,
			StatID: "boss_kills", TargetValue: 10,
			RewardCoins: 3000, RewardGems: 30,
		},
		{
			ID: "badge_combat_005", Name: "完美闪避", Description: "成功闪避100次",
			Type: BadgeTypeCombat, Rarity: BadgeRarityRare,
			StatID: "total_dodges", TargetValue: 100,
			RewardCoins: 1500, RewardGems: 15,
		},
		// 收集徽章
		{
			ID: "badge_collect_001", Name: "收藏家", Description: "收集50种不同道具",
			Type: BadgeTypeCollection, Rarity: BadgeRarityRare,
			StatID: "unique_items", TargetValue: 50,
			RewardCoins: 2000, RewardGems: 20,
		},
		{
			ID: "badge_collect_002", Name: "欧皇", Description: "获得10件传说装备",
			Type: BadgeTypeCollection, Rarity: BadgeRarityEpic,
			StatID: "legendary_items", TargetValue: 10,
			RewardCoins: 5000, RewardGems: 100,
		},
		// 社交徽章
		{
			ID: "badge_social_001", Name: "人脉广", Description: "添加10个好友",
			Type: BadgeTypeSocial, Rarity: BadgeRarityCommon,
			StatID: "friend_count", TargetValue: 10,
			RewardCoins: 500,
		},
		{
			ID: "badge_social_002", Name: "公会之光", Description: "公会等级达到10级",
			Type: BadgeTypeSocial, Rarity: BadgeRarityRare,
			StatID: "guild_level", TargetValue: 10,
			RewardCoins: 3000, RewardGems: 30,
		},
		// 赛季徽章
		{
			ID: "badge_season_001", Name: "赛季王者", Description: "达到赛季钻石段位",
			Type: BadgeTypeSeason, Rarity: BadgeRarityLegendary,
			StatID: "season_rank", TargetValue: 1,
			RewardGems: 500, RewardItem: "season_crown",
		},
		// 特殊徽章
		{
			ID: "badge_special_001", Name: "首充礼包", Description: "完成首次充值",
			Type: BadgeTypeSpecial, Rarity: BadgeRarityEpic,
			StatID: "first_recharge", TargetValue: 1,
			RewardItem: "first_recharge_box",
		},
		{
			ID: "badge_special_002", Name: "万元户", Description: "累计充值10000元",
			Type: BadgeTypeSpecial, Rarity: BadgeRarityLegendary,
			StatID: "total_recharge", TargetValue: 10000,
			RewardCoins: 10000, RewardGems: 1000,
		},
	}
	
	for _, b := range badges {
		bm.badges[b.ID] = b
		if b.StatID != "" {
			bm.statBadges[b.StatID] = append(bm.statBadges[b.StatID], b.ID)
		}
	}
}

// GetBadge 获取徽章
func (bm *BadgeManager) GetBadge(badgeID string) (*Badge, error) {
	b, ok := bm.badges[badgeID]
	if !ok {
		return nil, fmt.Errorf("徽章不存在")
	}
	return b, nil
}

// GetAllBadges 获取所有徽章
func (bm *BadgeManager) GetAllBadges() []*Badge {
	badges := make([]*Badge, 0, len(bm.badges))
	for _, b := range bm.badges {
		badges = append(badges, b)
	}
	return badges
}

// GetBadgesByType 按类型获取徽章
func (bm *BadgeManager) GetBadgesByType(btype BadgeType) []*Badge {
	badges := make([]*Badge, 0)
	for _, b := range bm.badges {
		if b.Type == btype {
			badges = append(badges, b)
		}
	}
	return badges
}

// CheckAndAward 检查并发放徽章
func (bm *BadgeManager) CheckAndAward(playerID, statID string, statValue int64) []*Badge {
	awarded := make([]*Badge, 0)
	badgeIDs, ok := bm.statBadges[statID]
	if !ok {
		return awarded
	}
	
	// 初始化玩家徽章数据
	if bm.playerBadges[playerID] == nil {
		bm.playerBadges[playerID] = make(map[string]*PlayerBadge)
	}
	
	for _, badgeID := range badgeIDs {
		badge, _ := bm.GetBadge(badgeID)
		if badge == nil {
			continue
		}
		
		// 检查是否已获得
		if _, exists := bm.playerBadges[playerID][badgeID]; exists {
			continue
		}
		
		// 检查是否达成条件
		if statValue >= badge.TargetValue {
			pb := &PlayerBadge{
				PlayerID:   playerID,
				BadgeID:    badgeID,
				ObtainedAt: time.Now().Unix(),
				IsNew:      true,
				StatValue:  statValue,
			}
			bm.playerBadges[playerID][badgeID] = pb
			awarded = append(awarded, badge)
		}
	}
	
	return awarded
}

// GetPlayerBadges 获取玩家徽章
func (bm *BadgeManager) GetPlayerBadges(playerID string) []*PlayerBadge {
	result := make([]*PlayerBadge, 0)
	for _, pb := range bm.playerBadges[playerID] {
		result = append(result, pb)
	}
	return result
}

// GetPlayerBadgeCount 获取玩家徽章数量
func (bm *BadgeManager) GetPlayerBadgeCount(playerID string) (int, int, int, int, int) {
	counts := make(map[BadgeRarity]int)
	for _, pb := range bm.playerBadges[playerID] {
		if badge, err := bm.GetBadge(pb.BadgeID); err == nil {
			counts[badge.Rarity]++
		}
	}
	return counts[BadgeRarityCommon], counts[BadgeRarityUncommon],
		counts[BadgeRarityRare], counts[BadgeRarityEpic], counts[BadgeRarityLegendary]
}

// MarkBadgeAsRead 标记徽章为已读
func (bm *BadgeManager) MarkBadgeAsRead(playerID, badgeID string) {
	if pb, ok := bm.playerBadges[playerID][badgeID]; ok {
		pb.IsNew = false
	}
}

// GetNewBadgeCount 获取新徽章数量
func (bm *BadgeManager) GetNewBadgeCount(playerID string) int {
	count := 0
	for _, pb := range bm.playerBadges[playerID] {
		if pb.IsNew {
			count++
		}
	}
	return count
}

// ============================================================
// 排行榜增强系统 Enhanced Leaderboard System
// ============================================================

// LeaderboardType 排行榜类型
type LeaderboardType int

const (
	LeaderboardTypeScore LeaderboardType = iota // 分数排行
	LeaderboardTypeKill                         // 击杀排行
	LeaderboardTypeWin                          // 胜率排行
	LeaderboardTypeCombo                        // 连击排行
	LeaderboardTypeGuild                        // 公会排行
	LeaderboardTypeRich                         // 财富排行
	LeaderboardTypeLevel                        // 等级排行
)

// LeaderboardEntry 排行榜条目
type LeaderboardEntry struct {
	Rank      int         `json:"rank"`       // 排名
	PlayerID  string      `json:"player_id"`  // 玩家ID
	PlayerName string    `json:"player_name"` // 玩家名
	Value     int64      `json:"value"`       // 数值
	GuildID   string     `json:"guild_id"`    // 公会ID
	GuildName string     `json:"guild_name"`  // 公会名
	Avatar    string     `json:"avatar"`      // 头像
	Title     string     `json:"title"`       // 称号
	UpdatedAt int64      `json:"updated_at"`  // 更新时间
}

// Leaderboard 排行榜
type Leaderboard struct {
	ID        string           `json:"id"`         // 排行榜ID
	Type      LeaderboardType `json:"type"`        // 类型
	Name      string          `json:"name"`        // 名称
	Season    int             `json:"season"`      // 赛季
	StartTime int64           `json:"start_time"`  // 开始时间
	EndTime   int64           `json:"end_time"`    // 结束时间
	Entries   []*LeaderboardEntry `json:"entries"` // 排行榜条目
	Top3      []*LeaderboardEntry `json:"top3"`    // 前3名缓存
}

// LeaderboardManager 排行榜管理器
type LeaderboardManager struct {
	leaderboards    map[string]*Leaderboard          // leaderboardID -> Leaderboard
	playerRanks     map[string]map[LeaderboardType]int // playerID -> type -> rank
	cache           map[string]int64                 // cache key -> expire time
}

// NewLeaderboardManager 创建排行榜管理器
func NewLeaderboardManager() *LeaderboardManager {
	return &LeaderboardManager{
		leaderboards: make(map[string]*Leaderboard),
		playerRanks:  make(map[string]map[LeaderboardType]int),
		cache:        make(map[string]int64),
	}
}

// CreateLeaderboard 创建排行榜
func (lm *LeaderboardManager) CreateLeaderboard(ltype LeaderboardType, name string, season int, duration time.Duration) *Leaderboard {
	lb := &Leaderboard{
		ID:        fmt.Sprintf("lb_%d_%d", ltype, season),
		Type:      ltype,
		Name:      name,
		Season:    season,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Add(duration).Unix(),
		Entries:   make([]*LeaderboardEntry, 0),
		Top3:      make([]*LeaderboardEntry, 0),
	}
	lm.leaderboards[lb.ID] = lb
	return lb
}

// UpdateEntry 更新排行榜条目
func (lm *LeaderboardManager) UpdateEntry(lbID, playerID, playerName, guildID, guildName, avatar, title string, value int64) {
	lb, ok := lm.leaderboards[lbID]
	if !ok {
		return
	}
	
	// 查找现有条目
	found := false
	for _, entry := range lb.Entries {
		if entry.PlayerID == playerID {
			entry.Value = value
			entry.PlayerName = playerName
			entry.GuildID = guildID
			entry.GuildName = guildName
			entry.Avatar = avatar
			entry.Title = title
			entry.UpdatedAt = time.Now().Unix()
			found = true
			break
		}
	}
	
	// 新条目
	if !found {
		entry := &LeaderboardEntry{
			PlayerID:   playerID,
			PlayerName: playerName,
			Value:      value,
			GuildID:    guildID,
			GuildName:  guildName,
			Avatar:     avatar,
			Title:      title,
			UpdatedAt:  time.Now().Unix(),
		}
		lb.Entries = append(lb.Entries, entry)
	}
	
	// 排序
	lm.sortLeaderboard(lb)
	
	// 更新玩家排名
	lm.updatePlayerRank(playerID, lb.Type, lb.getRank(playerID))
	
	// 更新Top3缓存
	lm.updateTop3(lb)
}

// sortLeaderboard 排序排行榜
func (lm *LeaderboardManager) sortLeaderboard(lb *Leaderboard) {
	for i := 0; i < len(lb.Entries)-1; i++ {
		for j := i + 1; j < len(lb.Entries); j++ {
			if lb.Entries[j].Value > lb.Entries[i].Value {
				lb.Entries[i], lb.Entries[j] = lb.Entries[j], lb.Entries[i]
			}
		}
	}
	
	// 更新排名
	for i, entry := range lb.Entries {
		entry.Rank = i + 1
	}
}

// updateTop3 更新Top3缓存
func (lm *LeaderboardManager) updateTop3(lb *Leaderboard) {
	count := 3
	if len(lb.Entries) < 3 {
		count = len(lb.Entries)
	}
	lb.Top3 = make([]*LeaderboardEntry, count)
	copy(lb.Top3, lb.Entries[:count])
}

// getRank 获取玩家排名
func (lb *Leaderboard) getRank(playerID string) int {
	for _, entry := range lb.Entries {
		if entry.PlayerID == playerID {
			return entry.Rank
		}
	}
	return -1
}

// updatePlayerRank 更新玩家排名缓存
func (lm *LeaderboardManager) updatePlayerRank(playerID string, ltype LeaderboardType, rank int) {
	if lm.playerRanks[playerID] == nil {
		lm.playerRanks[playerID] = make(map[LeaderboardType]int)
	}
	lm.playerRanks[playerID][ltype] = rank
}

// GetTop 获取排行榜Top N
func (lm *LeaderboardManager) GetTop(lbID string, n int) []*LeaderboardEntry {
	lb, ok := lm.leaderboards[lbID]
	if !ok {
		return nil
	}
	
	if n > len(lb.Entries) {
		n = len(lb.Entries)
	}
	
	result := make([]*LeaderboardEntry, n)
	copy(result, lb.Entries[:n])
	return result
}

// GetRank 获取玩家排名
func (lm *LeaderboardManager) GetRank(playerID string, ltype LeaderboardType) int {
	if ranks, ok := lm.playerRanks[playerID]; ok {
		if rank, ok := ranks[ltype]; ok {
			return rank
		}
	}
	return -1
}

// GetPlayerEntry 获取玩家条目
func (lm *LeaderboardManager) GetPlayerEntry(lbID, playerID string) *LeaderboardEntry {
	lb, ok := lm.leaderboards[lbID]
	if !ok {
		return nil
	}
	
	for _, entry := range lb.Entries {
		if entry.PlayerID == playerID {
			return entry
		}
	}
	return nil
}

// GetSurroundingPlayers 获取玩家周围排名
func (lm *LeaderboardManager) GetSurroundingPlayers(lbID, playerID string, count int) []*LeaderboardEntry {
	lb, ok := lm.leaderboards[lbID]
	if !ok {
		return nil
	}
	
	playerRank := lb.getRank(playerID)
	if playerRank == -1 {
		return nil
	}
	
	start := playerRank - count
	if start < 1 {
		start = 1
	}
	end := playerRank + count
	if end > len(lb.Entries) {
		end = len(lb.Entries)
	}
	
	result := make([]*LeaderboardEntry, 0)
	for i := start - 1; i < end; i++ {
		result = append(result, lb.Entries[i])
	}
	
	return result
}

// GetLeaderboard 获取排行榜
func (lm *LeaderboardManager) GetLeaderboard(lbID string) (*Leaderboard, error) {
	lb, ok := lm.leaderboards[lbID]
	if !ok {
		return nil, fmt.Errorf("排行榜不存在")
	}
	return lb, nil
}

// GetLeaderboardsByType 获取指定类型的排行榜
func (lm *LeaderboardManager) GetLeaderboardsByType(ltype LeaderboardType) []*Leaderboard {
	result := make([]*Leaderboard, 0)
	for _, lb := range lm.leaderboards {
		if lb.Type == ltype {
			result = append(result, lb)
		}
	}
	return result
}

// ClearExpired 清理过期排行榜
func (lm *LeaderboardManager) ClearExpired() int64 {
	now := time.Now().Unix()
	cleared := int64(0)
	
	for id, lb := range lm.leaderboards {
		if lb.EndTime < now {
			delete(lm.leaderboards, id)
			cleared++
		}
	}
	
	return cleared
}

// ============================================================
// 玩家统计系统 Player Statistics System
// ============================================================

// StatType 统计类型
type StatType string

const (
	StatBattles    StatType = "total_battles"    // 总战斗场次
	StatWins       StatType = "total_wins"       // 总胜利次数
	StatKills      StatType = "total_kills"      // 总击杀数
	StatDeaths     StatType = "total_deaths"     // 总死亡次数
	StatDamage     StatType = "total_damage"     // 总伤害
	StatHeal       StatType = "total_heal"       // 总治疗
	StatCoins      StatType = "total_coins"      // 总获得金币
	StatGems       StatType = "total_gems"       // 总获得钻石
	StatPlayTime   StatType = "total_play_time"  // 总游戏时间
	StatDodges     StatType = "total_dodges"      // 总闪避次数
	StatCombos     StatType = "max_combo"        // 最高连击
	StatBossKills  StatType = "boss_kills"       // Boss击杀
	StatWinStreak  StatType = "win_streak"       // 连胜
)

// PlayerStats 玩家统计
type PlayerStats struct {
	PlayerID      string           `json:"player_id"`       // 玩家ID
	Stats         map[StatType]int64 `json:"stats"`          // 统计数据
	TodayStats    map[StatType]int64 `json:"today_stats"`   // 今日统计
	WeekStats     map[StatType]int64 `json:"week_stats"`    // 本周统计
	SeasonStats   map[StatType]int64 `json:"season_stats"`  // 赛季统计
	LastReset     int64            `json:"last_reset"`      // 上次重置时间
}

// StatsManager 统计管理器
type StatsManager struct {
	playerStats map[string]*PlayerStats
}

// NewStatsManager 创建统计管理器
func NewStatsManager() *StatsManager {
	return &StatsManager{
		playerStats: make(map[string]*PlayerStats),
	}
}

// GetOrCreateStats 获取或创建玩家统计
func (sm *StatsManager) GetOrCreateStats(playerID string) *PlayerStats {
	ps, ok := sm.playerStats[playerID]
	if !ok {
		ps = &PlayerStats{
			PlayerID:    playerID,
			Stats:       make(map[StatType]int64),
			TodayStats:  make(map[StatType]int64),
			WeekStats:   make(map[StatType]int64),
			SeasonStats: make(map[StatType]int64),
			LastReset:   time.Now().Unix(),
		}
		sm.playerStats[playerID] = ps
	}
	return ps
}

// Increment 增加值
func (sm *StatsManager) Increment(playerID string, statType StatType, value int64) {
	ps := sm.GetOrCreateStats(playerID)
	ps.Stats[statType] += value
	ps.TodayStats[statType] += value
	ps.WeekStats[statType] += value
	ps.SeasonStats[statType] += value
}

// Set 设置值
func (sm *StatsManager) Set(playerID string, statType StatType, value int64) {
	ps := sm.GetOrCreateStats(playerID)
	ps.Stats[statType] = value
	
	// 检查是否是更高值
	if statType == StatCombos || statType == StatWinStreak {
		if value > ps.TodayStats[statType] {
			ps.TodayStats[statType] = value
		}
		if value > ps.WeekStats[statType] {
			ps.WeekStats[statType] = value
		}
		if value > ps.SeasonStats[statType] {
			ps.SeasonStats[statType] = value
		}
	} else {
		ps.TodayStats[statType] += value
		ps.WeekStats[statType] += value
		ps.SeasonStats[statType] += value
	}
}

// Get 获取统计值
func (sm *StatsManager) Get(playerID string, statType StatType) int64 {
	ps, ok := sm.playerStats[playerID]
	if !ok {
		return 0
	}
	return ps.Stats[statType]
}

// GetToday 获取今日统计
func (sm *StatsManager) GetToday(playerID string, statType StatType) int64 {
	ps, ok := sm.playerStats[playerID]
	if !ok {
		return 0
	}
	return ps.TodayStats[statType]
}

// ResetDaily 重置每日统计
func (sm *StatsManager) ResetDaily() int64 {
	now := time.Now().Unix()
	reset := int64(0)
	
	for _, ps := range sm.playerStats {
		if now - ps.LastReset >= 86400 { // 24小时
			ps.TodayStats = make(map[StatType]int64)
			reset++
		}
	}
	
	return reset
}

// GetWinRate 计算胜率
func (sm *StatsManager) GetWinRate(playerID string) float64 {
	battles := sm.Get(playerID, StatBattles)
	wins := sm.Get(playerID, StatWins)
	if battles == 0 {
		return 0
	}
	return float64(wins) / float64(battles) * 100
}

// GetKDA 计算KDA
func (sm *StatsManager) GetKDA(playerID string) float64 {
	kills := sm.Get(playerID, StatKills)
	deaths := sm.Get(playerID, StatDeaths)
	if deaths == 0 {
		return float64(kills)
	}
	return float64(kills) / float64(deaths)
}
