package game

import (
	"encoding/json"
	"fmt"
	"time"
)

// ActivityType 活动类型
type ActivityType int

const (
	ActivityTypeDaily    ActivityType = iota // 每日活动
	ActivityTypeWeekly                        // 每周活动
	ActivityTypeLimited                       // 限时活动
	ActivityTypeEvent                         // 事件活动
	ActivityTypeBattle                        // 对战活动
	ActivityTypeGift                          // 礼物活动
	ActivityTypeRecharge                      // 充值活动
)

// ActivityState 活动状态
type ActivityState int

const (
	ActivityStatePending ActivityState = iota // 未开始
	ActivityStateActive                        // 进行中
	ActivityStateEnded                         // 已结束
	ActivityStateReward                        // 领奖中
)

// ActivityReward 活动奖励
type ActivityReward struct {
	Type      string `json:"type"`      // 奖励类型: coin/diamond/item
	Value     int64  `json:"value"`     // 奖励值
	ItemID    string `json:"item_id"`   // 物品ID
	ItemCount int    `json:"item_count"` // 物品数量
}

// ActivityCondition 活动条件
type ActivityCondition struct {
	Type        string `json:"type"`         // 条件类型
	Target:"target"`            string `json // 目标值
	Value       int64  `json:"value"`        // 数值
	Description string `json:" 描述
}

// Activity description"` //活动结构
type Activity struct {
	ID          string           `json:"id"`           // 活动ID
	Name        string           `json:"name"`         // 活动名称
	Type        ActivityType    `json:"type"`         // 活动类型
	State       ActivityState   `json:"state"`        // 活动状态
	StartTime   int64            `json:"start_time"`   // 开始时间
	EndTime     int64            `json:"end_time"`     // 结束时间
	Rewards     []ActivityReward `json:"rewards"`      // 奖励列表
	Conditions  []ActivityCondition `json:"conditions"` // 条件列表
	Content     string           `json:"content"`      // 活动内容
	Icon        string           `json:"icon"`         // 图标URL
	Sort        int              `json:"sort"`         // 排序
	Config      map[string]interface{} `json:"config"` // 扩展配置
	CreatedAt   int64            `json:"created_at"`   // 创建时间
	UpdatedAt   int64            `json:"updated_at"`   // 更新时间
}

// PlayerActivity 玩家活动进度
type PlayerActivity struct {
	PlayerID    string                 `json:"player_id"`    // 玩家ID
	ActivityID  string                 `json:"activity_id"` // 活动ID
	Progress    int64                  `json:"progress"`     // 进度
	Claimed     bool                   `json:"claimed"`      // 是否已领取
	ClaimedAt   int64                  `json:"claimed_at"`   // 领取时间
	ExtraData   map[string]interface{} `json:"extra_data"`   // 额外数据
	UpdatedAt   int64                  `json:"updated_at"`   // 更新时间
}

// ActivityManager 活动管理器
type ActivityManager struct {
	activities map[string]*Activity
	playerProgress map[string]map[string]*PlayerActivity // playerID -> activityID -> progress
	scheduledTasks map[string]*time.Timer
}

// NewActivityManager 创建活动管理器
func NewActivityManager() *ActivityManager {
	return &ActivityManager{
		activities:     make(map[string]*Activity),
		playerProgress: make(map[string]map[string]*PlayerActivity),
		scheduledTasks: make(map[string]*time.Timer),
	}
}

// Init 初始化活动管理器
func (am *ActivityManager) Init() error {
	// 加载活动配置
	am.loadActivities()
	// 启动活动状态检查
	am.startActivityChecker()
	return nil
}

// loadActivities 加载活动配置
func (am *ActivityManager) loadActivities() {
	// 从数据库加载活动
	// 这里简化处理，实际从DB加载
	now := time.Now().Unix()
	
	// 添加示例活动
	sampleActivities := []*Activity{
		{
			ID:        "daily_login_001",
			Name:      "每日登录",
			Type:      ActivityTypeDaily,
			State:     ActivityStateActive,
			StartTime: now - 86400,
			EndTime:   now + 86400*7,
			Rewards: []ActivityReward{
				{Type: "coin", Value: 1000},
				{Type: "diamond", Value: 10},
			},
			Conditions: []ActivityCondition{
				{Type: "login", Description: "登录1次"},
			},
			Content: "每日登录领取奖励",
			Icon:    "https://example.com/icon/daily_login.png",
			Sort:    1,
			Config:  map[string]interface{}{"repeat": true},
		},
		{
			ID:        "week_battle_001",
			Name:      "周末对战",
			Type:      ActivityTypeWeekly,
			State:     ActivityStateActive,
			StartTime: now - 86400*3,
			EndTime:   now + 86400*4,
			Rewards: []ActivityReward{
				{Type: "coin", Value: 5000},
				{Type: "item", ItemID: "battle_pass", ItemCount: 1},
			},
			Conditions: []ActivityCondition{
				{Type: "win_count", Target: "battle", Value: 10, Description: "胜利10次"},
			},
			Content: "周末对战活动，胜利指定次数领取奖励",
			Icon:    "https://example.com/icon/week_battle.png",
			Sort:    2,
			Config:  map[string]interface{}{"min_level": 5},
		},
		{
			ID:        "recharge_first_001",
			Name:      "首充大礼包",
			Type:      ActivityTypeRecharge,
			State:     ActivityStateActive,
			StartTime: now - 86400*30,
			EndTime:   now + 86400*365,
			Rewards: []ActivityReward{
				{Type: "diamond", Value: 600},
				{Type: "item", ItemID: "hero_card", ItemCount: 1},
			},
			Conditions: []ActivityCondition{
				{Type: "first_recharge", Value: 1, Description: "首次充值"},
			},
			Content: "首充任意金额即可领取大礼包",
			Icon:    "https://example.com/icon/first_recharge.png",
			Sort:    3,
			Config:  map[string]interface{}{"recharge_amount": 1},
		},
		{
			ID:        "limited_gift_001",
			Name:      "限时礼物狂欢",
			Type:      ActivityTypeLimited,
			State:     ActivityStateActive,
			StartTime: now,
			EndTime:   now + 86400,
			Rewards: []ActivityReward{
				{Type: "coin", Value: 8888},
			},
			Conditions: []ActivityCondition{
				{Type: "send_gift", Value: 1, Description: "送任意礼物"},
			},
			Content: "活动期间送任意礼物即可领取奖励",
			Icon:    "https://example.com/icon/limited_gift.png",
			Sort:    4,
			Config:  map[string]interface{}{"gift_multiplier": 2},
		},
	}
	
	for _, a := range sampleActivities {
		am.activities[a.ID] = a
		am.scheduleActivityTask(a)
	}
}

// scheduleActivityTask 调度活动任务
func (am *ActivityManager) scheduleActivityTask(activity *Activity) {
	now := time.Now().Unix()
	
	// 调度开始任务
	if activity.State == ActivityStatePending && activity.StartTime > now {
		duration := time.Duration(activity.StartTime-now) * time.Second
		if duration > 0 {
			timer := time.AfterFunc(duration, func() {
				am.StartActivity(activity.ID)
			})
			am.scheduledTasks[activity.ID+"_start"] = timer
		}
	}
	
	// 调度结束任务
	if activity.EndTime > now {
		duration := time.Duration(activity.EndTime-now) * time.Second
		if duration > 0 {
			timer := time.AfterFunc(duration, func() {
				am.EndActivity(activity.ID)
			})
			am.scheduledTasks[activity.ID+"_end"] = timer
		}
	}
}

// startActivityChecker 启动活动状态检查器
func (am *ActivityManager) startActivityChecker() {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			am.checkActivityState()
		}
	}()
}

// checkActivityState 检查活动状态
func (am *ActivityManager) checkActivityState() {
	now := time.Now().Unix()
	
	for _, activity := range am.activities {
		oldState := activity.State
		
		switch {
		case now < activity.StartTime:
			activity.State = ActivityStatePending
		case now >= activity.StartTime && now < activity.EndTime:
			activity.State = ActivityStateActive
		case now >= activity.EndTime:
			activity.State = ActivityStateEnded
		}
		
		if oldState != activity.State {
			activity.UpdatedAt = now
			am.onActivityStateChange(activity, oldState)
		}
	}
}

// onActivityStateChange 活动状态变更回调
func (am *ActivityManager) onActivityStateChange(activity *Activity, oldState ActivityState) {
	// 通知玩家活动状态变更
	fmt.Printf("[Activity] Activity %s state changed: %d -> %d\n", activity.ID, oldState, activity.State)
	
	// 可以在这里发送通知给玩家
}

// StartActivity 开启活动
func (am *ActivityManager) StartActivity(activityID string) error {
	activity, ok := am.activities[activityID]
	if !ok {
		return fmt.Errorf("activity not found: %s", activityID)
	}
	
	activity.State = ActivityStateActive
	activity.UpdatedAt = time.Now().Unix()
	
	// 调度结束任务
	if activity.EndTime > time.Now().Unix() {
		duration := time.Duration(activity.EndTime-time.Now().Unix()) * time.Second
		timer := time.AfterFunc(duration, func() {
			am.EndActivity(activityID)
		})
		am.scheduledTasks[activityID+"_end"] = timer
	}
	
	return nil
}

// EndActivity 结束活动
func (am *ActivityManager) EndActivity(activityID string) error {
	activity, ok := am.activities[activityID]
	if !ok {
		return fmt.Errorf("activity not found: %s", activityID)
	}
	
	activity.State = ActivityStateEnded
	activity.UpdatedAt = time.Now().Unix()
	
	// 清理定时任务
	if timer, ok := am.scheduledTasks[activityID+"_end"]; ok {
		timer.Stop()
		delete(am.scheduledTasks, activityID+"_end")
	}
	
	return nil
}

// GetActivity 获取活动
func (am *ActivityManager) GetActivity(activityID string) (*Activity, error) {
	activity, ok := am.activities[activityID]
	if !ok {
		return nil, fmt.Errorf("activity not found: %s", activityID)
	}
	return activity, nil
}

// GetActiveActivities 获取进行中的活动
func (am *ActivityManager) GetActiveActivities() []*Activity {
	am.checkActivityState()
	
	var result []*Activity
	for _, activity := range am.activities {
		if activity.State == ActivityStateActive {
			result = append(result, activity)
		}
	}
	
	// 按sort排序
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Sort > result[j].Sort {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	return result
}

// GetActivitiesByType 按类型获取活动
func (am *ActivityManager) GetActivitiesByType(activityType ActivityType) []*Activity {
	var result []*Activity
	for _, activity := range am.activities {
		if activity.Type == activityType && activity.State == ActivityStateActive {
			result = append(result, activity)
		}
	}
	return result
}

// UpdatePlayerProgress 更新玩家活动进度
func (am *ActivityManager) UpdatePlayerProgress(playerID, activityID string, progress int64) error {
	activity, ok := am.activities[activityID]
	if !ok {
		return fmt.Errorf("activity not found: %s", activityID)
	}
	
	if activity.State != ActivityStateActive {
		return fmt.Errorf("activity not active: %s", activityID)
	}
	
	// 初始化玩家进度映射
	if _, ok := am.playerProgress[playerID]; !ok {
		am.playerProgress[playerID] = make(map[string]*PlayerActivity)
	}
	
	playerProgress, ok := am.playerProgress[playerID][activityID]
	if !ok {
		playerProgress = &PlayerActivity{
			PlayerID:   playerID,
			ActivityID: activityID,
			Progress:   0,
			Claimed:    false,
			ExtraData:  make(map[string]interface{}),
		}
		am.playerProgress[playerID][activityID] = playerProgress
	}
	
	// 更新进度
	playerProgress.Progress = progress
	playerProgress.UpdatedAt = time.Now().Unix()
	
	// 检查是否满足领取条件
	if am.CanClaimReward(playerID, activityID) {
		// 可以领取奖励
		fmt.Printf("[Activity] Player %s can claim reward for activity %s\n", playerID, activityID)
	}
	
	return nil
}

// CanClaimReward 检查是否可以领取奖励
func (am *ActivityManager) CanClaimReward(playerID, activityID string) bool {
	activity, ok := am.activities[activityID]
	if !ok {
		return false
	}
	
	playerProgress, ok := am.playerProgress[playerID][activityID]
	if !ok {
		return false
	}
	
	// 已领取
	if playerProgress.Claimed {
		return false
	}
	
	// 检查所有条件
	for _, condition := range activity.Conditions {
		if !am.checkCondition(playerID, activityID, condition) {
			return false
		}
	}
	
	return true
}

// checkCondition 检查条件是否满足
func (am *ActivityManager) checkCondition(playerID, activityID string, condition ActivityCondition) bool {
	playerProgress, ok := am.playerProgress[playerID][activityID]
	if !ok {
		return false
	}
	
	switch condition.Type {
	case "login", "win_count", "send_gift", "first_recharge":
		return playerProgress.Progress >= condition.Value
	default:
		return playerProgress.Progress >= condition.Value
	}
}

// ClaimReward 领取奖励
func (am *ActivityManager) ClaimReward(playerID, activityID string) ([]ActivityReward, error) {
	activity, ok := am.activities[activityID]
	if !ok {
		return nil, fmt.Errorf("activity not found: %s", activityID)
	}
	
	if !am.CanClaimReward(playerID, activityID) {
		return nil, fmt.Errorf("cannot claim reward: conditions not met")
	}
	
	playerProgress, ok := am.playerProgress[playerID][activityID]
	if !ok {
		return nil, fmt.Errorf("player progress not found")
	}
	
	// 标记已领取
	playerProgress.Claimed = true
	playerProgress.ClaimedAt = time.Now().Unix()
	playerProgress.UpdatedAt = time.Now().Unix()
	
	// 返回奖励列表
	return activity.Rewards, nil
}

// GetPlayerProgress 获取玩家活动进度
func (am *ActivityManager) GetPlayerProgress(playerID, activityID string) (*PlayerActivity, error) {
	if progressMap, ok := am.playerProgress[playerID]; ok {
		if progress, ok := progressMap[activityID]; ok {
			return progress, nil
		}
	}
	return nil, fmt.Errorf("player progress not found")
}

// GetPlayerActivities 获取玩家所有活动进度
func (am *ActivityManager) GetPlayerActivities(playerID string) []*PlayerActivity {
	var result []*PlayerActivity
	if progressMap, ok := am.playerProgress[playerID]; ok {
		for _, progress := range progressMap {
			result = append(result, progress)
		}
	}
	return result
}

// CreateActivity 创建活动
func (am *ActivityManager) CreateActivity(activity *Activity) error {
	if _, ok := am.activities[activity.ID]; ok {
		return fmt.Errorf("activity already exists: %s", activity.ID)
	}
	
	activity.CreatedAt = time.Now().Unix()
	activity.UpdatedAt = time.Now().Unix()
	
	am.activities[activity.ID] = activity
	am.scheduleActivityTask(activity)
	
	return nil
}

// UpdateActivity 更新活动
func (am *ActivityManager) UpdateActivity(activity *Activity) error {
	if _, ok := am.activities[activity.ID]; !ok {
		return fmt.Errorf("activity not found: %s", activity.ID)
	}
	
	activity.UpdatedAt = time.Now().Unix()
	am.activities[activity.ID] = activity
	
	return nil
}

// DeleteActivity 删除活动
func (am *ActivityManager) DeleteActivity(activityID string) error {
	if _, ok := am.activities[activityID]; !ok {
		return fmt.Errorf("activity not found: %s", activityID)
	}
	
	// 取消定时任务
	if timer, ok := am.scheduledTasks[activityID+"_start"]; ok {
		timer.Stop()
		delete(am.scheduledTasks, activityID+"_start")
	}
	if timer, ok := am.scheduledTasks[activityID+"_end"]; ok {
		timer.Stop()
		delete(am.scheduledTasks, activityID+"_end")
	}
	
	delete(am.activities, activityID)
	
	return nil
}

// GetActivityStats 获取活动统计
func (am *ActivityManager) GetActivityStats(activityID string) (map[string]interface{}, error) {
	activity, ok := am.activities[activityID]
	if !ok {
		return nil, fmt.Errorf("activity not found: %s", activityID)
	}
	
	stats := map[string]interface{}{
		"activity_id":   activity.ID,
		"activity_name": activity.Name,
		"state":         activity.State,
		"start_time":    activity.StartTime,
		"end_time":      activity.EndTime,
		"player_count":  0,
		"claim_count":   0,
	}
	
	// 统计参与玩家数
	for _, progressMap := range am.playerProgress {
		if progress, ok := progressMap[activityID]; ok {
			stats["player_count"] = stats["player_count"].(int) + 1
			if progress.Claimed {
				stats["claim_count"] = stats["claim_count"].(int) + 1
			}
		}
	}
	
	return stats, nil
}

// MarshalJSON 序列化活动
func (a *Activity) MarshalJSON() ([]byte, error) {
	type Alias Activity
	return json.Marshal(&struct {
		*Alias
		TypeName string `json:"type_name"`
		StateName string `json:"state_name"`
	}{
		Alias:     (*Alias)(a),
		TypeName:  a.GetTypeName(),
		StateName: a.GetStateName(),
	})
}

// GetTypeName 获取类型名称
func (a *Activity) GetTypeName() string {
	switch a.Type {
	case ActivityTypeDaily:
		return "daily"
	case ActivityTypeWeekly:
		return "weekly"
	case ActivityTypeLimited:
		return "limited"
	case ActivityTypeEvent:
		return "event"
	case ActivityTypeBattle:
		return "battle"
	case ActivityTypeGift:
		return "gift"
	case ActivityTypeRecharge:
		return "recharge"
	default:
		return "unknown"
	}
}

// GetStateName 获取状态名称
func (a *Activity) GetStateName() string {
	switch a.State {
	case ActivityStatePending:
		return "pending"
	case ActivityStateActive:
		return "active"
	case ActivityStateEnded:
		return "ended"
	case ActivityStateReward:
		return "reward"
	default:
		return "unknown"
	}
}
