package game

import (
	"errors"
	"time"
)

// SignInType 签到类型
type SignInType int

const (
	SignInTypeDaily    SignInType = 1 // 每日签到
	SignInTypeMonthly  SignInType = 2 // 月卡签到
	SignInTypeVIP      SignInType = 3 // VIP签到
)

// SignInState 签到状态
type SignInState int

const (
	SignInStateNone     SignInState = 0 // 未签到
	SignInStateSigned   SignInState = 1 // 已签到
	SignInStateRewarded SignInState = 2 // 已领取奖励
)

// SignInReward 签到奖励
type SignInReward struct {
	Day         int     // 第几天
	Coins       int64   // 金币
	Gems        int64   // 钻石
	Items       []int   // 物品ID
	ItemCounts  []int   // 物品数量
	IsSpecial   bool    // 是否为特殊奖励(第7天等)
}

// SignInRecord 玩家签到记录
type SignInRecord struct {
	PlayerID       int64       // 玩家ID
	SignInType     SignInType  // 签到类型
	LastSignInTime time.Time   // 上次签到时间
	TotalDays      int         // 累计签到天数
	ContinuousDays int         // 连续签到天数
	State          SignInState // 签到状态
	Rewards        []int       // 已领取的奖励天数
	CanBackfill    int         // 可补签次数
	LastBackfill   time.Time   // 上次补签时间
}

// SignInSystem 签到系统
type SignInSystem struct {
	dailyRewards   []SignInReward  // 每日签到奖励配置
	monthlyRewards []SignInReward  // 月卡签到奖励配置
	vipRewards     []SignInReward  // VIP签到奖励配置
	backfillCost   int64           // 补签消耗钻石
	maxBackfill    int             // 最多补签次数
}

// NewSignInSystem 创建签到系统
func NewSignInSystem() *SignInSystem {
	s := &SignInSystem{
		backfillCost:  50,  // 补签一次50钻
		maxBackfill:   3,   // 最多补签3次
	}

	// 初始化每日签到奖励(7天循环)
	s.dailyRewards = []SignInReward{
		{Day: 1, Coins: 100, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 2, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 3, Coins: 300, Gems: 5, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 4, Coins: 400, Gems: 5, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 5, Coins: 500, Gems: 10, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 6, Coins: 600, Gems: 10, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 7, Coins: 1000, Gems: 50, Items: []int{1}, ItemCounts: []int{1}, IsSpecial: true}, // 第7天大奖
	}

	// 月卡签到奖励
	s.monthlyRewards = []SignInReward{
		{Day: 1, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 2, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 3, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 4, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 5, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 6, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 7, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: true},
		{Day: 8, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 9, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 10, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 11, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 12, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 13, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 14, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: true},
		{Day: 15, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 16, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 17, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 18, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 19, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 20, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 21, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: true},
		{Day: 22, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 23, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 24, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 25, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 26, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 27, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 28, Coins: 800, Gems: 50, Items: []int{2}, ItemCounts: []int{1}, IsSpecial: true},
		{Day: 29, Coins: 200, Gems: 0, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 30, Coins: 1000, Gems: 100, Items: []int{3}, ItemCounts: []int{1}, IsSpecial: true},
	}

	// VIP签到奖励
	s.vipRewards = []SignInReward{
		{Day: 1, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 2, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 3, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 4, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 5, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 6, Coins: 500, Gems: 20, Items: nil, ItemCounts: nil, IsSpecial: false},
		{Day: 7, Coins: 2000, Gems: 100, Items: []int{4}, ItemCounts: []int{1}, IsSpecial: true},
	}

	return s
}

// GetRewards 获取签到奖励列表
func (s *SignInSystem) GetRewards(signType SignInType) []SignInReward {
	switch signType {
	case SignInTypeDaily:
		return s.dailyRewards
	case SignInTypeMonthly:
		return s.monthlyRewards
	case SignInTypeVIP:
		return s.vipRewards
	default:
		return s.dailyRewards
	}
}

// SignIn 签到
func (s *SignInSystem) SignIn(record *SignInRecord) (*SignInReward, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// 检查今天是否已签到
	if record.State == SignInStateSigned {
		lastSignDay := time.Date(record.LastSignInTime.Year(), record.LastSignInTime.Month(), record.LastSignInTime.Day(), 0, 0, 0, 0, time.UTC)
		if lastSignDay.Equal(today) {
			return nil, errors.New("今日已签到")
		}
	}

	// 更新连续签到
	if record.State == SignInStateSigned {
		yesterday := today.AddDate(0, 0, -1)
		lastSignDay := time.Date(record.LastSignInTime.Year(), record.LastSignInTime.Month(), record.LastSignInTime.Day(), 0, 0, 0, 0, time.UTC)
		if lastSignDay.Equal(yesterday) {
			record.ContinuousDays++
		} else {
			record.ContinuousDays = 1
		}
	} else {
		record.ContinuousDays = 1
	}

	// 更新总签到天数
	record.TotalDays++
	record.LastSignInTime = now
	record.State = SignInStateSigned

	// 获取奖励
	rewards := s.GetRewards(record.SignInType)
	dayIndex := (record.ContinuousDays - 1) % len(rewards)
	return &rewards[dayIndex], nil
}

// ClaimReward 领取奖励
func (s *SignInSystem) ClaimReward(record *SignInRecord, day int) (*SignInReward, error) {
	if record.State != SignInStateSigned {
		return nil, errors.New("今日尚未签到")
	}

	// 检查是否已领取
	for _, d := range record.Rewards {
		if d == day {
			return nil, errors.New("奖励已领取")
		}
	}

	rewards := s.GetRewards(record.SignInType)
	if day < 1 || day > len(rewards) {
		return nil, errors.New("无效的天数")
	}

	// 检查是否满足领取条件
	if day > record.ContinuousDays {
		return nil, errors.New("未满足领取条件")
	}

	// 记录已领取
	record.Rewards = append(record.Rewards, day)
	record.State = SignInStateRewarded

	reward := rewards[day-1]
	return &reward, nil
}

// Backfill 补签
func (s *SignInSystem) Backfill(record *SignInRecord, gems int64) (*SignInReward, error) {
	if record.CanBackfill <= 0 {
		return nil, errors.New("补签次数已用完")
	}

	if gems < s.backfillCost {
		return nil, errors.New("钻石不足")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// 检查今天是否已签到
	if record.State == SignInStateSigned {
		lastSignDay := time.Date(record.LastSignInTime.Year(), record.LastSignInTime.Month(), record.LastSignInTime.Day(), 0, 0, 0, 0, time.UTC)
		if lastSignDay.Equal(today) {
			return nil, errors.New("今日已签到,无法补签")
		}
	}

	// 更新签到记录
	record.CanBackfill--
	record.LastSignInTime = now
	record.ContinuousDays++
	record.TotalDays++
	record.State = SignInStateSigned

	// 返回补签奖励
	rewards := s.GetRewards(record.SignInType)
	dayIndex := (record.ContinuousDays - 1) % len(rewards)
	return &rewards[dayIndex], nil
}

// GetSignInInfo 获取签到信息
func (s *SignInSystem) GetSignInInfo(record *SignInRecord) map[string]interface{} {
	rewards := s.GetRewards(record.SignInType)
	continuousDay := record.ContinuousDays % len(rewards)
	if continuousDay == 0 {
		continuousDay = len(rewards)
	}

	info := map[string]interface{}{
		"player_id":        record.PlayerID,
		"sign_in_type":    record.SignInType,
		"total_days":      record.TotalDays,
		"continuous_days": record.ContinuousDays,
		"can_backfill":    record.CanBackfill,
		"rewards":         rewards,
		"claimed_rewards": record.Rewards,
		"current_day":     continuousDay,
		"today_signed":    record.State == SignInStateSigned,
		"backfill_cost":   s.backfillCost,
	}

	// 如果今天已签到,显示今天的奖励
	if record.State == SignInStateSigned {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		lastSignDay := time.Date(record.LastSignInTime.Year(), record.LastSignInTime.Month(), record.LastSignInTime.Day(), 0, 0, 0, 0, time.UTC)
		if lastSignDay.Equal(today) {
			info["today_reward"] = rewards[continuousDay-1]
		}
	}

	return info
}

// CreateSignInRecord 创建签到记录
func (s *SignInSystem) CreateSignInRecord(playerID int64, signType SignInType) *SignInRecord {
	return &SignInRecord{
		PlayerID:       playerID,
		SignInType:     signType,
		LastSignInTime: time.Time{},
		TotalDays:      0,
		ContinuousDays: 0,
		State:          SignInStateNone,
		Rewards:        []int{},
		CanBackfill:    s.maxBackfill,
		LastBackfill:   time.Time{},
	}
}

// GetContinuousBonus 获取连续签到加成
func (s *SignInSystem) GetContinuousBonus(continuousDays int) float64 {
	// 连续7天签到,金币加成50%
	if continuousDays >= 7 && continuousDays < 14 {
		return 1.5
	}
	// 连续14天签到,金币加成100%
	if continuousDays >= 14 && continuousDays < 30 {
		return 2.0
	}
	// 连续30天签到,金币加成150%
	if continuousDays >= 30 {
		return 2.5
	}
	// 基础加成
	return 1.0
}

// CheckAndResetDaily 检查并重置每日签到
func (s *SignInSystem) CheckAndResetDaily(record *SignInRecord) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	lastSignDay := time.Date(record.LastSignInTime.Year(), record.LastSignInTime.Month(), record.LastSignInTime.Day(), 0, 0, 0, 0, time.UTC)

	// 如果昨天没签到,重置连续天数
	yesterday := today.AddDate(0, 0, -1)
	if !lastSignDay.Equal(yesterday) && !lastSignDay.Before(yesterday) && record.TotalDays > 0 {
		record.ContinuousDays = 0
	}

	// 重置签到状态
	if !lastSignDay.Equal(today) {
		record.State = SignInStateNone
	}
}

// GetMonthlyCardExpireTime 获取月卡过期时间(示例)
func (s *SignInSystem) GetMonthlyCardExpireTime() time.Time {
	return time.Now().AddDate(0, 1, 0)
}

// HasMonthlyCard 检查是否有月卡
func (s *SignInSystem) HasMonthlyCard(record *SignInRecord) bool {
	return record.SignInType == SignInTypeMonthly
}

// UpgradeToMonthly 升级为月卡签到
func (s *SignInSystem) UpgradeToMonthly(record *SignInRecord) {
	record.SignInType = SignInTypeMonthly
}

// UpgradeToVIP 升级为VIP签到
func (s *SignInSystem) UpgradeToVIP(record *SignInRecord) {
	record.SignInType = SignInTypeVIP
}
