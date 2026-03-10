package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ============================================================
// 公会战系统 Guild War System
// ============================================================

// GuildWarState 公会战状态
type GuildWarState int

const (
	GuildWarStatePrepare GuildWarState = iota // 准备阶段
	GuildWarStateSignup                      // 报名阶段
	GuildWarStateRunning                     // 进行中
	GuildWarStateEnded                       // 已结束
)

// GuildWarType 公会战类型
type GuildWarType int

const (
	GuildWarTypeNormal GuildWarType = iota // 普通公会战
	GuildWarTypeElite                      // 精英公会战
	GuildWarTypeBoss                        // Boss公会战
)

// GuildWarScore 公会战积分
type GuildWarScore struct {
	GuildID       string `json:"guild_id"`        // 公会ID
	TotalScore    int64  `json:"total_score"`     // 总积分
	WinCount      int    `json:"win_count"`        // 胜利次数
	LoseCount     int    `json:"lose_count"`      // 失败次数
	KillCount     int64  `json:"kill_count"`       // 击杀数
	DamageDealt    int64  `json:"damage_dealt"`    // 造成伤害
	DamageTaken    int64  `json:"damage_taken"`    // 承受伤害
	HealCount     int64  `json:"heal_count"`       // 治疗量
	LastUpdated   int64  `json:"last_updated"`    // 最后更新
}

// GuildWarMatch 公会战匹配
type GuildWarMatch struct {
	MatchID       string         `json:"match_id"`       // 匹配ID
	GuildWarID    string         `json:"guild_war_id"`   // 公会战ID
	GuildA        string         `json:"guild_a"`        // 公会A
	GuildB        string         `json:"guild_b"`        // 公会B
	ScoreA        int64          `json:"score_a"`        // 公会A得分
	ScoreB        int64          `json:"score_b"`        // 公会B得分
	State         string         `json:"state"`          // 状态
	StartTime     int64          `json:"start_time"`     // 开始时间
	EndTime       int64          `json:"end_time"`       // 结束时间
	Winner        string         `json:"winner"`         // 获胜公会
}

// GuildWarPlayer 公会战玩家数据
type GuildWarPlayer struct {
	PlayerID    string `json:"player_id"`     // 玩家ID
	GuildID     string `json:"guild_id"`      // 公会ID
	MatchID     string `json:"match_id"`      // 匹配ID
	Score       int64  `json:"score"`         // 得分
	KillCount   int    `json:"kill_count"`     // 击杀数
	DeathCount  int    `json:"death_count"`   // 死亡次数
	DamageDealt int64  `json:"damage_dealt"`   // 造成伤害
	HealDone    int64  `json:"heal_done"`     // 治疗量
}

// GuildWar 公会战
type GuildWar struct {
	ID          string         `json:"id"`           // 战次ID
	Name        string         `json:"name"`         // 战次名称
	Type        GuildWarType  `json:"type"`         // 战次类型
	State       GuildWarState `json:"state"`        // 状态
	Season      int            `json:"season"`       // 赛季
	StartTime   int64          `json:"start_time"`   // 开始时间
	EndTime     int64          `json:"end_time"`     // 结束时间
	SignupEnd   int64          `json:"signup_end"`   // 报名截止
	MinLevel    int            `json:"min_level"`    // 最低公会等级
	MaxTeams    int            `json:"max_teams"`    // 最大参赛队伍数
	Teams       []string       `json:"teams"`        // 参赛公会列表
	Matches     []*GuildWarMatch `json:"matches"`    // 匹配列表
	Config      map[string]interface{} `json:"config"` // 配置
}

// GuildWarManager 公会战管理器
type GuildWarManager struct {
	guildWars   map[string]*GuildWar        // guildWarID -> GuildWar
	playerData  map[string]*GuildWarPlayer  // playerID -> GuildWarPlayer
	guildScores map[string]map[string]*GuildWarScore // season -> guildID -> Score
	schedules   map[string]*time.Timer
	rand        *rand.Rand
}

// NewGuildWarManager 创建公会战管理器
func NewGuildWarManager() *GuildWarManager {
	return &GuildWarManager{
		guildWars:  make(map[string]*GuildWar),
		playerData: make(map[string]*GuildWarPlayer),
		guildScores: make(map[string]map[string]*GuildWarScore),
		schedules:  make(map[string]*time.Timer),
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateGuildWar 创建公会战
func (gwm *GuildWarManager) CreateGuildWar(name string, warType GuildWarType, duration time.Duration) *GuildWar {
	war := &GuildWar{
		ID:        fmt.Sprintf("gw_%d", time.Now().UnixNano()),
		Name:      name,
		Type:      warType,
		State:     GuildWarStatePrepare,
		Season:    1,
		StartTime: time.Now().Add(duration).Unix(),
		EndTime:   time.Now().Add(duration * 2).Unix(),
		SignupEnd: time.Now().Add(duration - time.Hour).Unix(),
		MinLevel:  3,
		MaxTeams:  16,
		Teams:     make([]string, 0),
		Matches:   make([]*GuildWarMatch, 0),
		Config:    make(map[string]interface{}),
	}
	gwm.guildWars[war.ID] = war
	return war
}

// SignupGuild 公会报名
func (gwm *GuildWarManager) SignupGuild(warID, guildID string) error {
	war, ok := gwm.guildWars[warID]
	if !ok {
		return fmt.Errorf("公会战不存在")
	}
	if war.State != GuildWarStateSignup {
		return fmt.Errorf("当前不是报名阶段")
	}
	if len(war.Teams) >= war.MaxTeams {
		return fmt.Errorf("报名已满")
	}
	war.Teams = append(war.Teams, guildID)
	return nil
}

// StartGuildWar 开始公会战
func (gwm *GuildWarManager) StartGuildWar(warID string) error {
	war, ok := gwm.guildWars[warID]
	if !ok {
		return fmt.Errorf("公会战不存在")
	}
	if len(war.Teams) < 2 {
		return fmt.Errorf("参赛公会不足")
	}
	
	// 生成匹配
	war.Matches = gwm.generateMatches(war)
	war.State = GuildWarStateRunning
	return nil
}

// generateMatches 生成匹配
func (gwm *GuildWarManager) generateMatches(war *GuildWar) []*GuildWarMatch {
	matches := make([]*GuildWarMatch, 0)
	teams := war.Teams
	
	// 随机打乱
	for i := len(teams) - 1; i > 0; i-- {
		j := gwm.rand.Intn(i + 1)
		teams[i], teams[j] = teams[j], teams[i]
	}
	
	// 配对
	for i := 0; i < len(teams); i += 2 {
		if i+1 >= len(teams) {
			break
		}
		match := &GuildWarMatch{
			MatchID:    fmt.Sprintf("match_%d_%d", warID, i/2),
			GuildWarID: war.ID,
			GuildA:     teams[i],
			GuildB:     teams[i+1],
			State:      "pending",
			StartTime:  war.StartTime,
			EndTime:    war.EndTime,
		}
		matches = append(matches, match)
	}
	
	return matches
}

// UpdateMatchScore 更新匹配得分
func (gwm *GuildWarManager) UpdateMatchScore(matchID string, guildID string, scoreDelta int64) error {
	for _, war := range gwm.guildWars {
		for _, match := range war.Matches {
			if match.MatchID == matchID {
				if match.GuildA == guildID {
					match.ScoreA += scoreDelta
				} else if match.GuildB == guildID {
					match.ScoreB += scoreDelta
				}
				return nil
			}
		}
	}
	return fmt.Errorf("匹配不存在")
}

// GetGuildWar 获取公会战
func (gwm *GuildWarManager) GetGuildWar(warID string) (*GuildWar, error) {
	war, ok := gwm.guildWars[warID]
	if !ok {
		return nil, fmt.Errorf("公会战不存在")
	}
	return war, nil
}

// GetPlayerData 获取玩家公会战数据
func (gwm *GuildWarManager) GetPlayerData(playerID string) *GuildWarPlayer {
	return gwm.playerData[playerID]
}

// JoinWar 加入公会战
func (gwm *GuildWarManager) JoinWar(playerID, guildID, matchID string) *GuildWarPlayer {
	player := &GuildWarPlayer{
		PlayerID: playerID,
		GuildID:  guildID,
		MatchID:  matchID,
		Score:    0,
	}
	gwm.playerData[playerID] = player
	return player
}

// RecordKill 记录击杀
func (gwm *GuildWarManager) RecordKill(killerID, victimID string) {
	if killer, ok := gwm.playerData[killerID]; ok {
		killer.KillCount++
		killer.Score += 100
	}
	if victim, ok := gwm.playerData[victimID]; ok {
		victim.DeathCount++
	}
}

// RecordDamage 记录伤害
func (gwm *GuildWarManager) RecordDamage(playerID string, damage int64, isHeal bool) {
	if player, ok := gwm.playerData[playerID]; ok {
		if isHeal {
			player.HealDone += damage
			player.Score += damage / 2
		} else {
			player.DamageDealt += damage
			player.Score += damage / 10
		}
	}
}

// CalculateReward 计算奖励
func (gwm *GuildWarManager) CalculateReward(guildID string, rank int) map[string]int64 {
	rewards := make(map[string]int64)
	
	// 基础奖励
	baseScore := int64(1000)
	baseCoin := int64(10000)
	
	// 根据排名计算倍数
	multiplier := math.Max(0, 1-float64(rank-1)*0.1)
	
	rewards["score"] = int64(float64(baseScore) * multiplier)
	rewards["coin"] = int64(float64(baseCoin) * multiplier)
	
	// 前三名额外奖励
	if rank == 1 {
		rewards["diamond"] = 500
		rewards["title"] = 1 // 冠军称号
	} else if rank == 2 {
		rewards["diamond"] = 300
	} else if rank == 3 {
		rewards["diamond"] = 100
	}
	
	return rewards
}

// EndGuildWar 结束公会战
func (gwm *GuildWarManager) EndGuildWar(warID string) error {
	war, ok := gwm.guildWars[warID]
	if !ok {
		return fmt.Errorf("公会战不存在")
	}
	
	war.State = GuildWarStateEnded
	
	// 计算排名并发放奖励
	rank := 1
	for _, match := range war.Matches {
		if match.ScoreA > match.ScoreB {
			match.Winner = match.GuildA
			rewardA := gwm.CalculateReward(match.GuildA, rank)
			rewardB := gwm.CalculateReward(match.GuildB, rank+1)
			fmt.Printf("Guild %s wins, rewards: %v\n", match.GuildA, rewardA)
			fmt.Printf("Guild %s loses, rewards: %v\n", match.GuildB, rewardB)
			rank += 2
		} else if match.ScoreB > match.ScoreA {
			match.Winner = match.GuildB
			rewardA := gwm.CalculateReward(match.GuildA, rank+1)
			rewardB := gwm.CalculateReward(match.GuildB, rank)
			fmt.Printf("Guild %s wins, rewards: %v\n", match.GuildB, rewardB)
			fmt.Printf("Guild %s loses, rewards: %v\n", match.GuildA, rewardA)
			rank += 2
		}
	}
	
	return nil
}

// ============================================================
// 排行榜挑战系统 Rank Challenge System
// ============================================================

// ChallengeType 挑战类型
type ChallengeType int

const (
	ChallengeTypeScore ChallengeType = iota // 分数挑战
	ChallengeTypeTime                       // 限时挑战
	ChallengeTypeKill                       // 击杀挑战
	ChallengeTypeSurvival                   // 生存挑战
	ChallengeTypeCombo                      // 连击挑战
)

// Challenge 挑战
type Challenge struct {
	ID          string          `json:"id"`           // 挑战ID
	Name        string          `json:"name"`         // 挑战名称
	Type        ChallengeType  `json:"type"`         // 挑战类型
	Description string          `json:"description"`  // 描述
	Target      int64           `json:"target"`       // 目标值
	StartTime   int64           `json:"start_time"`   // 开始时间
	EndTime     int64           `json:"end_time"`     // 结束时间
	RewardCoins int64           `json:"reward_coins"` // 金币奖励
	RewardGems  int64           `json:"reward_gems"`  // 钻石奖励
	RewardItem  string          `json:"reward_item"`  // 道具奖励
	MaxAttempts int             `json:"max_attempts"` // 最大尝试次数
}

// PlayerChallenge 玩家挑战记录
type PlayerChallenge struct {
	PlayerID     string `json:"player_id"`     // 玩家ID
	ChallengeID  string `json:"challenge_id"`  // 挑战ID
	Attempts     int    `json:"attempts"`      // 尝试次数
	BestScore    int64  `json:"best_score"`    // 最佳成绩
	Completed    bool   `json:"completed"`     // 是否完成
	Claimed      bool   `json:"claimed"`       // 是否领取奖励
	LastAttempt  int64  `json:"last_attempt"`  // 上次尝试时间
	CompletedAt  int64  `json:"completed_at"`  // 完成时间
}

// ChallengeManager 挑战管理器
type ChallengeManager struct {
	challenges      map[string]*Challenge           // challengeID -> Challenge
	playerChallenges map[string]map[string]*PlayerChallenge // playerID -> challengeID -> PlayerChallenge
}

// NewChallengeManager 创建挑战管理器
func NewChallengeManager() *ChallengeManager {
	return &ChallengeManager{
		challenges:       make(map[string]*Challenge),
		playerChallenges:  make(map[string]map[string]*PlayerChallenge),
	}
}

// CreateChallenge 创建挑战
func (cm *ChallengeManager) CreateChallenge(name string, ctype ChallengeType, target int64, duration time.Duration) *Challenge {
	challenge := &Challenge{
		ID:          fmt.Sprintf("challenge_%d", time.Now().UnixNano()),
		Name:        name,
		Type:        ctype,
		Description: fmt.Sprintf("目标: %d", target),
		Target:      target,
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Add(duration).Unix(),
		RewardCoins: target * 10,
		RewardGems:  target / 100,
		MaxAttempts: 3,
	}
	cm.challenges[challenge.ID] = challenge
	return challenge
}

// GetChallenge 获取挑战
func (cm *ChallengeManager) GetChallenge(challengeID string) (*Challenge, error) {
	challenge, ok := cm.challenges[challengeID]
	if !ok {
		return nil, fmt.Errorf("挑战不存在")
	}
	return challenge, nil
}

// StartChallenge 开始挑战
func (cm *ChallengeManager) StartChallenge(playerID, challengeID string) (*PlayerChallenge, error) {
	challenge, err := cm.GetChallenge(challengeID)
	if err != nil {
		return nil, err
	}
	
	// 检查时间
	now := time.Now().Unix()
	if now < challenge.StartTime || now > challenge.EndTime {
		return nil, fmt.Errorf("挑战不在开放时间内")
	}
	
	// 初始化玩家挑战
	if cm.playerChallenges[playerID] == nil {
		cm.playerChallenges[playerID] = make(map[string]*PlayerChallenge)
	}
	
	pc, exists := cm.playerChallenges[playerID][challengeID]
	if !exists {
		pc = &PlayerChallenge{
			PlayerID:    playerID,
			ChallengeID: challengeID,
			Attempts:    0,
			BestScore:   0,
		}
		cm.playerChallenges[playerID][challengeID] = pc
	}
	
	// 检查次数
	if pc.Attempts >= challenge.MaxAttempts {
		return nil, fmt.Errorf("已达到最大尝试次数")
	}
	
	pc.Attempts++
	pc.LastAttempt = time.Now().Unix()
	
	return pc, nil
}

// UpdateScore 更新成绩
func (cm *ChallengeManager) UpdateScore(playerID, challengeID string, score int64) error {
	pc, ok := cm.playerChallenges[playerID][challengeID]
	if !ok {
		return fmt.Errorf("玩家挑战记录不存在")
	}
	
	// 更新最佳成绩
	if score > pc.BestScore {
		pc.BestScore = score
	}
	
	// 检查是否完成
	challenge, err := cm.GetChallenge(challengeID)
	if err != nil {
		return err
	}
	
	if score >= challenge.Target && !pc.Completed {
		pc.Completed = true
		pc.CompletedAt = time.Now().Unix()
	}
	
	return nil
}

// ClaimReward 领取奖励
func (cm *ChallengeManager) ClaimReward(playerID, challengeID string) (map[string]int64, error) {
	pc, ok := cm.playerChallenges[playerID][challengeID]
	if !ok {
		return nil, fmt.Errorf("玩家挑战记录不存在")
	}
	
	if !pc.Completed {
		return nil, fmt.Errorf("挑战未完成")
	}
	
	if pc.Claimed {
		return nil, fmt.Errorf("奖励已领取")
	}
	
	pc.Claimed = true
	
	challenge, _ := cm.GetChallenge(challengeID)
	rewards := map[string]int64{
		"coins": challenge.RewardCoins,
		"gems":  challenge.RewardGems,
	}
	
	return rewards, nil
}

// GetPlayerChallenges 获取玩家挑战列表
func (cm *ChallengeManager) GetPlayerChallenges(playerID string) []*PlayerChallenge {
	result := make([]*PlayerChallenge, 0)
	for _, pc := range cm.playerChallenges[playerID] {
		result = append(result, pc)
	}
	return result
}

// GetActiveChallenges 获取进行中的挑战
func (cm *ChallengeManager) GetActiveChallenges() []*Challenge {
	result := make([]*Challenge, 0)
	now := time.Now().Unix()
	for _, c := range cm.challenges {
		if now >= c.StartTime && now <= c.EndTime {
			result = append(result, c)
		}
	}
	return result
}

// InitDefaultChallenges 初始化默认挑战
func (cm *ChallengeManager) InitDefaultChallenges() {
	cm.CreateChallenge("分数达人", ChallengeTypeScore, 10000, 24*time.Hour)
	cm.CreateChallenge("击杀狂魔", ChallengeTypeKill, 100, 24*time.Hour)
	cm.CreateChallenge("生存大师", ChallengeTypeSurvival, 300, 24*time.Hour)
	cm.CreateChallenge("连击王者", ChallengeTypeCombo, 50, 24*time.Hour)
	cm.CreateChallenge("限时挑战", ChallengeTypeTime, 60, 2*time.Hour)
}

// ============================================================
// 签到奖励预览系统 Sign-In Reward Preview System
// ============================================================

// SignInPreview 签到预览
type SignInPreview struct {
	Day         int            `json:"day"`          // 天数
	DayName     string         `json:"day_name"`     // 日期名称
	Coins       int64          `json:"coins"`        // 金币
	Gems        int64          `json:"gems"`         // 钻石
	ItemID      string         `json:"item_id"`      // 道具ID
	ItemCount   int            `json:"item_count"`   // 道具数量
	IsSpecial   bool           `json:"is_special"`   // 是否特别奖励
	IsDouble    bool           `json:"is_double"`    // 是否双倍
}

// MonthlyCardReward 月卡奖励预览
type MonthlyCardReward struct {
	Days        int    `json:"days"`         // 月卡天数
	DailyCoins  int64  `json:"daily_coins"` // 每日金币
	DailyGems   int64  `json:"daily_gems"`  // 每日钻石
	TotalCoins  int64  `json:"total_coins"` // 总金币
	TotalGems   int64  `json:"total_gems"`  // 总钻石
	Price       int64  `json:"price"`       // 价格(钻石)
}

// SignInRewardManager 签到奖励管理器
type SignInRewardManager struct {
	dailyRewards   []SignInPreview     // 每日签到奖励
	monthlyRewards []MonthlyCardReward // 月卡奖励
	vipRewards     []SignInPreview     // VIP签到奖励
}

// NewSignInRewardManager 创建签到奖励管理器
func NewSignInRewardManager() *SignInRewardManager {
	srm := &SignInRewardManager{
		dailyRewards:   make([]SignInPreview, 7),
		monthlyRewards: make([]MonthlyCardReward, 0),
		vipRewards:     make([]SignInPreview, 7),
	}
	srm.initRewards()
	return srm
}

// initRewards 初始化奖励配置
func (srm *SignInRewardManager) initRewards() {
	// 每日签到奖励 (7天循环)
	daily := []struct {
		coins    int64
		gems     int64
		itemID   string
		itemCnt  int
		special  bool
	}{
		{100, 0, "", 0, false},
		{150, 5, "", 0, false},
		{200, 10, "gem_pack_small", 1, false},
		{250, 15, "", 0, false},
		{300, 20, "weapon_box", 1, false},
		{350, 25, "", 0, false},
		{500, 50, "ultimate_box", 1, true}, // 第7天大奖
	}
	
	for i, d := range daily {
		srm.dailyRewards[i] = SignInPreview{
			Day:        i + 1,
			DayName:    fmt.Sprintf("第%d天", i+1),
			Coins:      d.coins,
			Gems:       d.gems,
			ItemID:     d.itemID,
			ItemCount:  d.itemCnt,
			IsSpecial:  d.special,
		}
	}
	
	// VIP签到奖励
	vipDaily := []struct {
		coins   int64
		gems    int64
		itemID  string
		itemCnt int
	}{
		{200, 10, "", 0},
		{300, 15, "", 0},
		{400, 20, "gem_pack_medium", 1},
		{500, 25, "", 0},
		{600, 30, "weapon_box", 1},
		{700, 35, "", 0},
		{1000, 100, "ultimate_box", 1},
	}
	
	for i, v := range vipDaily {
		srm.vipRewards[i] = SignInPreview{
			Day:       i + 1,
			DayName:   fmt.Sprintf("VIP第%d天", i+1),
			Coins:     v.coins,
			Gems:      v.gems,
			ItemID:    v.itemID,
			ItemCount: v.itemCnt,
			IsSpecial: i == 6,
		}
	}
	
	// 月卡奖励
	srm.monthlyRewards = []MonthlyCardReward{
		{
			Days:       30,
			DailyCoins: 100,
			DailyGems:  10,
			TotalCoins: 3000,
			TotalGems:  300,
			Price:      300,
		},
		{
			Days:       30,
			DailyCoins: 200,
			DailyGems:  20,
			TotalCoins: 6000,
			TotalGems:  600,
			Price:      500,
		},
	}
}

// GetDailyRewards 获取每日签到奖励预览
func (srm *SignInRewardManager) GetDailyRewards() []SignInPreview {
	return srm.dailyRewards
}

// GetMonthlyRewards 获取月卡奖励预览
func (srm *SignInRewardManager) GetMonthlyRewards() []MonthlyCardReward {
	return srm.monthlyRewards
}

// GetVIPRewards 获取VIP签到奖励预览
func (srm *SignInRewardManager) GetVIPRewards() []SignInPreview {
	return srm.vipRewards
}

// CalculateStreakBonus 计算连续签到加成
func (srm *SignInRewardManager) CalculateStreakBonus(baseCoins, baseGems int64, streakDays int) (int64, int64) {
	bonus := 1.0
	
	// 连续签到加成
	if streakDays >= 7 {
		bonus = 1.5 // 150%
	} else if streakDays >= 3 {
		bonus = 1.2 // 120%
	}
	
	return int64(float64(baseCoins) * bonus), int64(float64(baseGems) * bonus)
}

// GetTotalMonthlyValue 计算月卡总价值
func (srm *SignInRewardManager) GetTotalMonthlyValue(cardType int) (int64, int64) {
	if cardType >= len(srm.monthlyRewards) {
		return 0, 0
	}
	mc := srm.monthlyRewards[cardType]
	return mc.TotalCoins, mc.TotalGems
}
