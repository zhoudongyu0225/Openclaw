package game

import (
	"encoding/json"
	"fmt"
	"time"
)

// ============================================
// 每日任务系统 (Daily Task System)
// ============================================

// TaskCategory 任务分类
type TaskCategory int

const (
	TaskCategoryDaily    TaskCategory = iota // 每日任务
	TaskCategoryWeekly                       // 每周任务
	TaskCategoryAchievement                  // 成就任务
	TaskCategoryMain                         // 主线任务
	TaskCategorySide                         // 支线任务
)

// TaskType 任务类型
type TaskType int

const (
	TaskTypeWinBattle      TaskType = iota // 赢得战斗
	TaskTypeKillEnemy                       // 击杀敌人
	TaskTypeCollectItem                    // 收集道具
	TaskTypeUseSkill                       // 使用技能
	TaskTypeReachCombo                     // 达成连击
	TaskTypeReachScore                     // 达成分数
	TaskTypeUpgradeTower                   // 升级塔
	TaskTypeJoinRoom                       // 加入房间
	TaskTypeInviteFriend                   // 邀请好友
	TaskTypeSendGift                       // 送礼物
	TaskTypeSignIn                         // 签到
	TaskTypeDrawGacha                      // 抽卡
	TaskTypeEnhanceEquipment               // 强化装备
	TaskTypeCompleteLevel                  // 通关关卡
)

// TaskStatus 任务状态
type TaskStatus int

const (
	TaskStatusLocked   TaskStatus = iota // 未解锁
	TaskStatusAvailable                   // 可领取
	TaskStatusInProgress                  // 进行中
	TaskStatusCompleted                    // 已完成
	TaskStatusClaimed                     // 已领取
)

// DailyTask 每日任务
type DailyTask struct {
	TaskID       string            `json:"task_id"`       // 任务ID
	Category     TaskCategory      `json:"category"`     // 任务分类
	Type         TaskType          `json:"type"`         // 任务类型
	Title        string            `json:"title"`        // 任务标题
	Description  string            `json:"description"`   // 任务描述
	Target       int               `json:"target"`       // 目标数量
	Progress     int               `json:"progress"`     // 当前进度
	Rewards      TaskRewards       `json:"rewards"`      // 任务奖励
	Status       TaskStatus        `json:"status"`       // 任务状态
	StartTime    int64             `json:"start_time"`   // 开始时间
	EndTime      int64             `json:"end_time"`     // 结束时间
	ResetTime    int64             `json:"reset_time"`   // 重置时间
	Priority     int               `json:"priority"`     // 优先级
	Icon         string            `json:"icon"`         // 图标
	IsDaily      bool              `json:"is_daily"`     // 是否每日刷新
	Conditions   []TaskCondition   `json:"conditions"`   // 完成条件
}

// TaskCondition 任务条件
type TaskCondition struct {
	Type       string `json:"type"`       // 条件类型
	TargetID   string `json:"target_id"`  // 目标ID
	Operator   string `json:"operator"`   // 操作符 (>, <, =, >=, <=)
	Value      int    `json:"value"`       // 目标值
}

// TaskRewards 任务奖励
type TaskRewards struct {
	Coins   int `json:"coins"`   // 金币
	Gems    int `json:"gems"`    // 钻石
	Exp     int `json:"exp"`     // 经验
	Items   []RewardItem `json:"items"` // 道具
}

// RewardItem 奖励道具
type RewardItem struct {
	ItemID   string `json:"item_id"` // 道具ID
	Count    int    `json:"count"`   // 数量
}

// PlayerTaskProgress 玩家任务进度
type PlayerTaskProgress struct {
	PlayerID      string            `json:"player_id"`      // 玩家ID
	Tasks         map[string]*DailyTask `json:"tasks"`      // 任务列表
	CompletedCount int                `json:"completed_count"` // 今日完成任务数
	TotalPoints   int                `json:"total_points"`   // 总积分
	LastResetTime int64             `json:"last_reset_time"` // 上次重置时间
}

// DailyTaskManager 每日任务管理器
type DailyTaskManager struct {
	tasks      map[string]*DailyTask  // 任务模板
	presets    map[string][]string    // 预设任务组
	taskPools  map[TaskCategory][]string // 任务池
}

// NewDailyTaskManager 创建每日任务管理器
func NewDailyTaskManager() *DailyTaskManager {
	mgr := &DailyTaskManager{
		tasks:     make(map[string]*DailyTask),
		presets:   make(map[string][]string),
		taskPools: make(map[TaskCategory][]string),
	}
	mgr.initPresets()
	mgr.initTaskPools()
	return mgr
}

// initPresets 初始化预设任务组
func (m *DailyTaskManager) initPresets() {
	// 每日任务组
	m.presets["daily_beginner"] = []string{
		"task_daily_win_1",
		"task_daily_kill_10",
		"task_daily_combo_20",
		"task_daily_score_1000",
	}

	m.presets["daily_intermediate"] = []string{
		"task_daily_win_3",
		"task_daily_kill_30",
		"task_daily_combo_50",
		"task_daily_score_5000",
		"task_daily_use_skill_5",
	}

	m.presets["daily_advanced"] = []string{
		"task_daily_win_5",
		"task_daily_kill_50",
		"task_daily_combo_100",
		"task_daily_score_10000",
		"task_daily_upgrade_tower_3",
	}

	// 每周任务组
	m.presets["weekly_battle"] = []string{
		"task_weekly_win_20",
		"task_weekly_kill_200",
		"task_weekly_combo_500",
		"task_weekly_score_50000",
	}

	m.presets["weekly_social"] = []string{
		"task_weekly_invite_friend_5",
		"task_weekly_send_gift_10",
		"task_weekly_join_room_20",
	}
}

// initTaskPools 初始化任务池
func (m *DailyTaskManager) initTaskPools() {
	// 每日任务池
	m.taskPools[TaskCategoryDaily] = []string{
		"task_daily_win_1", "task_daily_win_2", "task_daily_win_3", "task_daily_win_5",
		"task_daily_kill_10", "task_daily_kill_20", "task_daily_kill_30", "task_daily_kill_50",
		"task_daily_combo_20", "task_daily_combo_50", "task_daily_combo_100",
		"task_daily_score_1000", "task_daily_score_3000", "task_daily_score_5000",
		"task_daily_use_skill_3", "task_daily_use_skill_5", "task_daily_use_skill_10",
		"task_daily_upgrade_tower_1", "task_daily_upgrade_tower_3", "task_daily_upgrade_tower_5",
		"task_daily_join_room_1", "task_daily_join_room_3", "task_daily_join_room_5",
		"task_daily_sign_in", "task_daily_draw_gacha", "task_daily_enhance_equipment",
	}

	// 每周任务池
	m.taskPools[TaskCategoryWeekly] = []string{
		"task_weekly_win_10", "task_weekly_win_20", "task_weekly_win_30",
		"task_weekly_kill_100", "task_weekly_kill_200", "task_weekly_kill_300",
		"task_weekly_combo_200", "task_weekly_combo_500", "task_weekly_combo_1000",
		"task_weekly_score_20000", "task_weekly_score_50000", "task_weekly_score_100000",
		"task_weekly_invite_friend_3", "task_weekly_invite_friend_5", "task_weekly_invite_friend_10",
		"task_weekly_send_gift_5", "task_weekly_send_gift_10", "task_weekly_send_gift_20",
		"task_weekly_join_room_10", "task_weekly_join_room_20", "task_weekly_join_room_30",
		"task_weekly_complete_level_5", "task_weekly_complete_level_10",
	}
}

// CreateTask 创建任务
func (m *DailyTaskManager) CreateTask(taskID, title, desc string, taskType TaskType, target int, rewards TaskRewards) *DailyTask {
	task := &DailyTask{
		TaskID:      taskID,
		Category:    TaskCategoryDaily,
		Type:        taskType,
		Title:       title,
		Description: desc,
		Target:      target,
		Progress:    0,
		Rewards:     rewards,
		Status:      TaskStatusAvailable,
		StartTime:   time.Now().Unix(),
		EndTime:     time.Now().Add(24 * time.Hour).Unix(),
		ResetTime:   time.Now().Add(24 * time.Hour).Unix(),
		IsDaily:     true,
	}
	m.tasks[taskID] = task
	return task
}

// GetTask 获取任务
func (m *DailyTaskManager) GetTask(taskID string) (*DailyTask, bool) {
	task, ok := m.tasks[taskID]
	return task, ok
}

// UpdateTaskProgress 更新任务进度
func (m *DailyTaskManager) UpdateTaskProgress(progress *PlayerTaskProgress, taskID string, value int) error {
	task, ok := progress.Tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if task.Status != TaskStatusInProgress && task.Status != TaskStatusAvailable {
		return fmt.Errorf("task not in progress: %s", taskID)
	}

	task.Progress += value
	if task.Progress >= task.Target {
		task.Status = TaskStatusCompleted
		progress.CompletedCount++
	}

	return nil
}

// ClaimTaskReward 领取任务奖励
func (m *DailyTaskManager) ClaimTaskReward(progress *PlayerTaskProgress, taskID string, player *Player) error {
	task, ok := progress.Tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if task.Status != TaskStatusCompleted {
		return fmt.Errorf("task not completed: %s", taskID)
	}

	// 发放奖励
	player.AddCoins(task.Rewards.Coins)
	player.AddGems(task.Rewards.Gems)
	player.AddExp(task.Rewards.Exp)

	for _, item := range task.Rewards.Items {
		player.Inventory.AddItem(item.ItemID, item.Count)
	}

	task.Status = TaskStatusClaimed
	progress.TotalPoints += task.Target

	return nil
}

// ResetDailyTasks 重置每日任务
func (m *DailyTaskManager) ResetDailyTasks(progress *PlayerTaskProgress) {
	now := time.Now().Unix()
	progress.LastResetTime = now

	for _, task := range progress.Tasks {
		if task.IsDaily && task.Status != TaskStatusLocked {
			task.Progress = 0
			task.Status = TaskStatusAvailable
			task.EndTime = now + 24*3600
		}
	}

	progress.CompletedCount = 0
}

// GenerateDailyTasks 生成每日任务
func (m *DailyTaskManager) GenerateDailyTasks(playerID string) *PlayerTaskProgress {
	progress := &PlayerTaskProgress{
		PlayerID:      playerID,
		Tasks:         make(map[string]*DailyTask),
		CompletedCount: 0,
		TotalPoints:   0,
		LastResetTime: time.Now().Unix(),
	}

	// 从任务池中随机选择任务
	dailyPool := m.taskPools[TaskCategoryDaily]
	taskCount := 5 + len(dailyPool)%3 // 5-7个任务

	for i := 0; i < taskCount && i < len(dailyPool); i++ {
		taskID := dailyPool[i]
		if task, ok := m.tasks[taskID]; ok {
			newTask := *task
			newTask.Progress = 0
			newTask.Status = TaskStatusAvailable
			newTask.StartTime = time.Now().Unix()
			newTask.EndTime = time.Now().Add(24 * time.Hour).Unix()
			progress.Tasks[taskID] = &newTask
		}
	}

	return progress
}

// GetTaskList 获取任务列表
func (m *DailyTaskManager) GetTaskList(progress *PlayerTaskProgress) []*DailyTask {
	tasks := make([]*DailyTask, 0, len(progress.Tasks))
	for _, task := range progress.Tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetTaskByCategory 按分类获取任务
func (m *DailyTaskManager) GetTaskByCategory(progress *PlayerTaskProgress, category TaskCategory) []*DailyTask {
	tasks := make([]*DailyTask, 0)
	for _, task := range progress.Tasks {
		if task.Category == category {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// ============================================
// 任务事件处理器
// ============================================

// OnBattleWin 战斗胜利事件
func (m *DailyTaskManager) OnBattleWin(progress *PlayerTaskProgress, isBoss bool) {
	taskType := TaskTypeWinBattle
	target := 1
	if isBoss {
		target = 1
	}

	for _, task := range progress.Tasks {
		if task.Type == taskType {
			m.UpdateTaskProgress(progress, task.TaskID, target)
		}
	}
}

// OnEnemyKill 击杀敌人事件
func (m *DailyTaskManager) OnEnemyKill(progress *PlayerTaskProgress, enemyType string) {
	for _, task := range progress.Tasks {
		if task.Type == TaskTypeKillEnemy {
			m.UpdateTaskProgress(progress, task.TaskID, 1)
		}
	}
}

// OnComboReached 达成连击事件
func (m *DailyTaskManager) OnComboReached(progress *PlayerTaskProgress, comboCount int) {
	for _, task := range progress.Tasks {
		if task.Type == TaskTypeReachCombo && task.Target <= comboCount {
			m.UpdateTaskProgress(progress, task.TaskID, comboCount)
		}
	}
}

// OnScoreReached 达成分数事件
func (m *DailyTaskManager) OnScoreReached(progress *PlayerTaskProgress, score int) {
	for _, task := range progress.Tasks {
		if task.Type == TaskTypeReachScore && task.Target <= score {
			m.UpdateTaskProgress(progress, task.TaskID, score)
		}
	}
}

// OnSkillUsed 使用技能事件
func (m *DailyTaskManager) OnSkillUsed(progress *PlayerTaskProgress, skillID string) {
	for _, task := range progress.Tasks {
		if task.Type == TaskTypeUseSkill {
			m.UpdateTaskProgress(progress, task.TaskID, 1)
		}
	}
}

// OnTowerUpgraded 升级塔事件
func (m *DailyTaskManager) OnTowerUpgraded(progress *PlayerTaskProgress) {
	for _, task := range progress.Tasks {
		if task.Type == TaskTypeUpgradeTower {
			m.UpdateTaskProgress(progress, task.TaskID, 1)
		}
	}
}

// OnLevelCompleted 通关关卡事件
func (m *DailyTaskManager) OnLevelCompleted(progress *PlayerTaskProgress, levelID string, stars int) {
	for _, task := range progress.Tasks {
		if task.Type == TaskTypeCompleteLevel {
			m.UpdateTaskProgress(progress, task.TaskID, 1)
		}
	}
}

// ============================================
// 序列化方法
// ============================================

// MarshalJSON 序列化
func (t *DailyTask) MarshalJSON() ([]byte, error) {
	type Alias DailyTask
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}

// UnmarshalJSON 反序列化
func (t *DailyTask) UnmarshalJSON(data []byte) error {
	type Alias DailyTask
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}
