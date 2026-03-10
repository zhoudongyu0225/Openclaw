// 签到系统 - danmaku_game/server/signin.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// SignInReward 签到奖励配置
type SignInReward struct {
	Day         int    `json:"day"`          // 天数 (1-30)
	ItemID      string `json:"item_id"`     // 道具ID
	ItemName    string `json:"item_name"`   // 道具名称
	ItemCount   int    `json:"item_count"`  // 道具数量
	ExtraItemID string `json:"extra_item_id"` // 额外奖励道具ID (连续签到)
	ExtraCount  int    `json:"extra_count"`  // 额外奖励数量
	VIPOnly     bool   `json:"vip_only"`     // 是否VIP专属
}

// SignInRecord 签到记录
type SignInRecord struct {
	PlayerID      int64     `json:"player_id"`       // 玩家ID
	TotalDays     int       `json:"total_days"`      // 总签到天数
	ContinuousDays int      `json:"continuous_days"` // 连续签到天数
	LastSignIn    time.Time `json:"last_sign_in"`    // 上次签到时间
	SignInDates   []string  `json:"sign_in_dates"`   // 签到日期列表
	ReceivedRewards []int   `json:"received_rewards"` // 已领取奖励天数
	Month         int       `json:"month"`           // 签到月份
	Year          int       `json:"year"`            // 签到年份
}

// SignInSystem 签到系统
type SignInSystem struct {
	db            *Database
	cache         *Cache
	inventory     *Inventory
	rewards       []*SignInReward
	monthRewards  map[int][]*SignInReward // 月份奖励
	maxDays       int                    // 最大签到天数
}

// NewSignInSystem 创建签到系统
func NewSignInSystem(db *Database, cache *Cache, inventory *Inventory) *SignInSystem {
	s := &SignInSystem{
		db:       db,
		cache:    cache,
		inventory: inventory,
		rewards:  make([]*SignInReward, 0),
		monthRewards: make(map[int][]*SignInReward),
		maxDays: 30,
	}

	// 初始化签到奖励
	s.initRewards()

	return s
}

// initRewards 初始化签到奖励配置
func (s *SignInSystem) initRewards() {
	// 普通签到奖励 (1-30天循环)
	s.rewards = []*SignInReward{
		{Day: 1, ItemID: "gold", ItemName: "金币", ItemCount: 100, ExtraCount: 50},
		{Day: 2, ItemID: "gem", ItemName: "钻石", ItemCount: 10, ExtraCount: 5},
		{Day: 3, ItemID: "gold", ItemName: "金币", ItemCount: 200, ExtraCount: 100},
		{Day: 4, ItemID: "exp_card", ItemName: "经验卡", ItemCount: 1, ExtraCount: 1},
		{Day: 5, ItemID: "gem", ItemName: "钻石", ItemCount: 20, ExtraCount: 10},
		{Day: 6, ItemID: "gold", ItemName: "金币", ItemCount: 300, ExtraCount: 150},
		{Day: 7, ItemID: "random_box", ItemName: "随机宝箱", ItemCount: 1, ExtraCount: 1, VIPOnly: true},
		{Day: 8, ItemID: "gem", ItemName: "钻石", ItemCount: 30, ExtraCount: 15},
		{Day: 9, ItemID: "gold", ItemName: "金币", ItemCount: 400, ExtraCount: 200},
		{Day: 10, ItemID: "skill_point", ItemName: "技能点", ItemCount: 50, ExtraCount: 25},
		{Day: 11, ItemID: "gem", ItemName: "钻石", ItemCount: 40, ExtraCount: 20},
		{Day: 12, ItemID: "gold", ItemName: "金币", ItemCount: 500, ExtraCount: 250},
		{Day: 13, ItemID: "exp_card_2", ItemName: "双倍经验卡", ItemCount: 1, ExtraCount: 1},
		{Day: 14, ItemID: "gem", ItemName: "钻石", ItemCount: 50, ExtraCount: 25, VIPOnly: true},
		{Day: 15, ItemID: "random_box", ItemName: "随机宝箱", ItemCount: 2, ExtraCount: 1},
		{Day: 16, ItemID: "gold", ItemName: "金币", ItemCount: 600, ExtraCount: 300},
		{Day: 17, ItemID: "gem", ItemName: "钻石", ItemCount: 60, ExtraCount: 30},
		{Day: 18, ItemID: "gold", ItemName: "金币", ItemCount: 700, ExtraCount: 350},
		{Day: 19, ItemID: "skill_point", ItemName: "技能点", ItemCount: 100, ExtraCount: 50},
		{Day: 20, ItemID: "gem", ItemName: "钻石", ItemCount: 70, ExtraCount: 35},
		{Day: 21, ItemID: "gold", ItemName: "金币", ItemCount: 800, ExtraCount: 400, VIPOnly: true},
		{Day: 22, ItemID: "gem", ItemName: "钻石", ItemCount: 80, ExtraCount: 40},
		{Day: 23, ItemID: "gold", ItemName: "金币", ItemCount: 900, ExtraCount: 450},
		{Day: 24, ItemID: "rare_box", ItemName: "稀有宝箱", ItemCount: 1, ExtraCount: 1},
		{Day: 25, ItemID: "gem", ItemName: "钻石", ItemCount: 100, ExtraCount: 50},
		{Day: 26, ItemID: "gold", ItemName: "金币", ItemCount: 1000, ExtraCount: 500},
		{Day: 27, ItemID: "gem", ItemName: "钻石", ItemCount: 100, ExtraCount: 50},
		{Day: 28, ItemID: "gold", ItemName: "金币", ItemCount: 1500, ExtraCount: 750},
		{Day: 29, ItemID: "epic_box", ItemName: "史诗宝箱", ItemCount: 1, ExtraCount: 1, VIPOnly: true},
		{Day: 30, ItemID: "legendary_box", ItemName: "传说宝箱", ItemCount: 1, ExtraCount: 1},
	}
}

// GetTodayReward 获取今日奖励
func (s *SignInSystem) GetTodayReward(continuousDays int) *SignInReward {
	day := continuousDays%30 + 1
	if day == 0 {
		day = 1
	}

	for _, reward := range s.rewards {
		if reward.Day == day {
			return reward
		}
	}

	return nil
}

// SignIn 签到
func (s *SignInSystem) SignIn(playerID int64, vipLevel int) (*SignInRecord, *SignInReward, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// 获取签到记录
	record, err := s.GetRecord(playerID)
	if err != nil {
		return nil, nil, err
	}

	// 检查今天是否已签到
	lastSignInDate := record.LastSignIn.Format("2006-01-02")
	if lastSignInDate == today {
		return nil, nil, errors.New("今日已签到")
	}

	// 检查是否断签
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	isContinuous := lastSignInDate == yesterday

	if !isContinuous {
		// 断签，重置连续天数
		record.ContinuousDays = 0
	}

	// 增加连续签到天数
	record.ContinuousDays++
	record.TotalDays++

	// 更新签到日期
	record.SignInDates = append(record.SignInDates, today)
	record.LastSignIn = now

	// 更新年月
	record.Year = now.Year()
	record.Month = int(now.Month())

	// 获取奖励
	reward := s.GetTodayReward(record.ContinuousDays)
	if reward == nil {
		return nil, nil, errors.New("奖励配置错误")
	}

	// VIP专属检查
	if reward.VIPOnly && vipLevel < 1 {
		// 非VIP给予替代奖励
		reward = &SignInReward{
			Day:       reward.Day,
			ItemID:    "gold",
			ItemName:  "金币",
			ItemCount: reward.ItemCount,
			ExtraCount: reward.ExtraCount,
		}
	}

	// 发放基础奖励
	s.inventory.AddItem(reward.ItemID, reward.ItemCount)

	// 发放连续签到额外奖励 (第7天倍数)
	if record.ContinuousDays%7 == 0 && reward.ExtraItemID != "" {
		s.inventory.AddItem(reward.ExtraItemID, reward.ExtraCount)
	}

	// 补签检测 (可选功能)
	// if !isContinuous {
	//     // 可以设计补签卡道具
	// }

	// 保存记录
	s.SaveRecord(record)

	// 缓存清除
	s.cache.Del(fmt.Sprintf("signin:%d", playerID))

	return record, reward, nil
}

// GetRecord 获取签到记录
func (s *SignInSystem) GetRecord(playerID int64) (*SignInRecord, error) {
	cacheKey := fmt.Sprintf("signin:%d", playerID)

	// 尝试缓存
	if cached, err := s.cache.Get(cacheKey); err == nil {
		var record SignInRecord
		if json.Unmarshal([]byte(cached), &record) == nil {
			return &record, nil
		}
	}

	// 从数据库查询
	query := `SELECT player_id, total_days, continuous_days, last_sign_in, sign_in_dates, 
			  received_rewards, month, year FROM sign_in_records WHERE player_id = ?`

	record := &SignInRecord{}
	var lastSignIn []byte
	var signInDates []byte
	var receivedRewards []byte

	err := m.db.QueryRow(query, playerID).Scan(
		&record.PlayerID,
		&record.TotalDays,
		&record.ContinuousDays,
		&lastSignIn,
		&signInDates,
		&receivedRewards,
		&record.Month,
		&record.Year,
	)
	if err != nil {
		// 没有记录，创建新记录
		record = &SignInRecord{
			PlayerID: playerID,
			TotalDays: 0,
			ContinuousDays: 0,
			LastSignIn: time.Time{},
			SignInDates: make([]string, 0),
			ReceivedRewards: make([]int, 0),
		}
		return record, nil
	}

	json.Unmarshal(lastSignIn, &record.LastSignIn)
	json.Unmarshal(signInDates, &record.SignInDates)
	json.Unmarshal(receivedRewards, &record.ReceivedRewards)

	// 缓存
	if data, err := json.Marshal(record); err == nil {
		s.cache.SetEX(cacheKey, string(data), 3600)
	}

	return record, nil
}

// SaveRecord 保存签到记录
func (s *SignInSystem) SaveRecord(record *SignInRecord) error {
	signInDates, _ := json.Marshal(record.SignInDates)
	receivedRewards, _ := json.Marshal(record.ReceivedRewards)
	lastSignIn, _ := json.Marshal(record.LastSignIn)

	query := `INSERT INTO sign_in_records (player_id, total_days, continuous_days, last_sign_in, 
			  sign_in_dates, received_rewards, month, year) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?) 
			  ON DUPLICATE KEY UPDATE 
			  total_days = ?, continuous_days = ?, last_sign_in = ?, 
			  sign_in_dates = ?, received_rewards = ?, month = ?, year = ?`

	_, err := s.db.Exec(query,
		record.PlayerID, record.TotalDays, record.ContinuousDays, lastSignIn,
		signInDates, receivedRewards, record.Month, record.Year,
		record.TotalDays, record.ContinuousDays, lastSignIn,
		signInDates, receivedRewards, record.Month, record.Year,
	)

	return err
}

// CanSignIn 今日是否可以签到
func (s *SignInSystem) CanSignIn(playerID int64) (bool, error) {
	record, err := s.GetRecord(playerID)
	if err != nil {
		return false, err
	}

	today := time.Now().Format("2006-01-02")
	lastSignInDate := record.LastSignIn.Format("2006-01-02")

	return lastSignInDate != today, nil
}

// GetContinuousDays 获取连续签到天数
func (s *SignInSystem) GetContinuousDays(playerID int64) (int, error) {
	record, err := s.GetRecord(playerID)
	if err != nil {
		return 0, err
	}

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	lastSignInDate := record.LastSignIn.Format("2006-01-02")

	// 今天已签到，返回连续天数
	// 今天未签到但昨天签到了，返回连续天数
	// 否则返回0
	if lastSignInDate == today || lastSignInDate == yesterday {
		return record.ContinuousDays, nil
	}

	return 0, nil
}

// GetMonthSignInDays 获取本月签到天数
func (s *SignInSystem) GetMonthSignInDays(playerID int64) (int, error) {
	record, err := s.GetRecord(playerID)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	thisMonth := now.Month()
	thisYear := now.Year()

	// 如果不是本月，重置
	if record.Month != int(thisMonth) || record.Year != thisYear {
		return 0, nil
	}

	count := 0
	for _, date := range record.SignInDates {
		if len(date) >= 7 {
			month := date[5:7]
			if month == fmt.Sprintf("%02d", thisMonth) {
				count++
			}
		}
	}

	return count, nil
}

// GetSignInStatus 获取签到状态
type SignInStatus struct {
	CanSignIn      bool            `json:"can_sign_in"`       // 今日可签到
	ContinuousDays int             `json:"continuous_days"`  // 连续签到天数
	TotalDays      int             `json:"total_days"`       // 总签到天数
	MonthDays      int             `json:"month_days"`       // 本月签到天数
	TodayReward    *SignInReward   `json:"today_reward"`     // 今日奖励
	NextReward     *SignInReward   `json:"next_reward"`      // 明日奖励
	ReceivedDays   []int           `json:"received_days"`    // 已领取的奖励天数
}

// GetSignInStatus 获取签到状态
func (s *SignInSystem) GetSignInStatus(playerID int64, vipLevel int) (*SignInStatus, error) {
	record, err := s.GetRecord(playerID)
	if err != nil {
		return nil, err
	}

	canSignIn, _ := s.CanSignIn(playerID)
	continuousDays, _ := s.GetContinuousDays(playerID)
	monthDays, _ := s.GetMonthSignInDays(playerID)

	// 计算今日和明日奖励
	todayReward := s.GetTodayReward(continuousDays)
	if todayReward != nil && todayReward.VIPOnly && vipLevel < 1 {
		todayReward = &SignInReward{
			Day:       todayReward.Day,
			ItemID:    "gold",
			ItemName:  "金币",
			ItemCount: todayReward.ItemCount,
		}
	}

	nextReward := s.GetTodayReward(continuousDays + 1)
	if nextReward != nil && nextReward.VIPOnly && vipLevel < 1 {
		nextReward = &SignInReward{
			Day:       nextReward.Day,
			ItemID:    "gold",
			ItemName:  "金币",
			ItemCount: nextReward.ItemCount,
		}
	}

	return &SignInStatus{
		CanSignIn:      canSignIn,
		ContinuousDays: continuousDays,
		TotalDays:      record.TotalDays,
		MonthDays:      monthDays,
		TodayReward:    todayReward,
		NextReward:     nextReward,
		ReceivedDays:   record.ReceivedRewards,
	}, nil
}

// GetRewardByDay 获取指定天数的奖励
func (s *SignInSystem) GetRewardByDay(day int) *SignInReward {
	for _, reward := range s.rewards {
		if reward.Day == day {
			return reward
		}
	}
	return nil
}

// GetMonthlyRewards 获取月度签到奖励配置
func (s *SignInSystem) GetMonthlyRewards(month int) []*SignInReward {
	// 月度奖励通常是累计签到奖励
	// 这里返回30天的完整奖励配置
	return s.rewards
}

// ResetMonthlySignIn 重置月度签到 (月初)
func (s *SignInSystem) ResetMonthlySignIn(playerID int64) error {
	record, err := s.GetRecord(playerID)
	if err != nil {
		return err
	}

	now := time.Now()
	record.Month = int(now.Month())
	record.Year = now.Year()

	// 清空已领取奖励
	record.ReceivedRewards = make([]int, 0)

	return s.SaveRecord(record)
}

// GetSignInCalendar 获取签到日历
type SignInCalendar struct {
	Year      int      `json:"year"`       // 年
	Month     int      `json:"month"`      // 月
	SignInDays []int   `json:"sign_in_days"` // 签到日期列表
	Rewards   []*SignInReward `json:"rewards"` // 奖励配置
}

// GetSignInCalendar 获取签到日历
func (s *SignInSystem) GetSignInCalendar(playerID int64, year, month int) (*SignInCalendar, error) {
	record, err := s.GetRecord(playerID)
	if err != nil {
		return nil, err
	}

	calendar := &SignInCalendar{
		Year:     year,
		Month:    month,
		SignInDays: make([]int, 0),
		Rewards:  s.rewards,
	}

	// 筛选指定月份的签到日期
	for _, date := range record.SignInDates {
		if len(date) >= 7 {
			y, _ := fmt.Sscanf(date[:4], "%d", &year)
			m, _ := fmt.Sscanf(date[5:7], "%d", &month)
			d, _ := fmt.Sscanf(date[8:10], "%d", &d)

			if y == year && m == month {
				calendar.SignInDays = append(calendar.SignInDays, d)
			}
		}
	}

	return calendar, nil
}

// InitSignInTable 初始化签到表
func (s *SignInSystem) InitSignInTable() error {
	query := `CREATE TABLE IF NOT EXISTS sign_in_records (
		player_id BIGINT PRIMARY KEY,
		total_days INT NOT NULL DEFAULT 0,
		continuous_days INT NOT NULL DEFAULT 0,
		last_sign_in DATETIME,
		sign_in_dates JSON,
		received_rewards JSON,
		month INT NOT NULL DEFAULT 0,
		year INT NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	_, err := s.db.Exec(query)
	return err
}

// GetTotalSignInRank 获取签到排行榜
func (s *SignInSystem) GetTotalSignInRank(limit int) ([]*SignInRecord, error) {
	query := `SELECT player_id, total_days, continuous_days, last_sign_in, sign_in_dates, 
			  received_rewards, month, year 
			  FROM sign_in_records 
			  ORDER BY total_days DESC, continuous_days DESC 
			  LIMIT ?`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*SignInRecord
	for rows.Next() {
		record := &SignInRecord{}
		var lastSignIn []byte
		var signInDates []byte
		var receivedRewards []byte

		err := rows.Scan(&record.PlayerID, &record.TotalDays, &record.ContinuousDays,
			&lastSignIn, &signInDates, &receivedRewards, &record.Month, &record.Year)
		if err != nil {
			continue
		}

		json.Unmarshal(lastSignIn, &record.LastSignIn)
		json.Unmarshal(signInDates, &record.SignInDates)
		json.Unmarshal(receivedRewards, &record.ReceivedRewards)

		records = append(records, record)
	}

	return records, nil
}
