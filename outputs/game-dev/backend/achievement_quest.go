package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================
// 成就系统 (Achievement System)
// ============================================

// 成就类型
type AchievementType int

const (
	AchievementTypeKill AchievementType = iota // 击杀成就
	AchievementTypeScore      // 分数成就
	AchievementTypeWave      // 波次成就
	AchievementTypeTower     // 防御塔成就
	AchievementTypeGift      // 礼物成就
	AchievementTypeDanmaku   // 弹幕成就
	AchievementTypeTime      // 时间成就
	AchievementTypeSpecial   // 特殊成就
)

// 成就条件
type AchievementCondition struct {
	Type  AchievementType `json:"type"`
	Target int           `json:"target"` // 目标值
}

// 成就定义
type Achievement struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Icon        string               `json:"icon"`
	Type        AchievementType     `json:"type"`
	Condition   AchievementCondition `json:"condition"`
	Reward      *AchievementReward  `json:"reward"`
	Hidden      bool                `json:"hidden"` // 是否隐藏成就
}

// 成就奖励
type AchievementReward struct {
	Exp   int `json:"exp"`
	Coins int `json:"coins"`
	Gems  int `json:"gems"`
}

// 玩家成就进度
type PlayerAchievement struct {
	AchievementID string    `json:"achievementId"`
	Progress     int       `json:"progress"`
	Completed    bool      `json:"completed"`
	CompletedAt  time.Time `json:"completedAt"`
	RewardClaimed bool    `json:"rewardClaimed"`
}

// 成就管理器
type AchievementManager struct {
	// 成就定义
	Achievements map[string]*Achievement
	
	// 玩家成就进度
	PlayerProgress map[string]map[string]*PlayerAchievement // playerID -> achievementID -> progress
	
	mu sync.RWMutex
}

// 成就奖励
func NewAchievementManager() *AchievementManager {
	am := &AchievementManager{
		Achievements:    make(map[string]*Achievement),
		PlayerProgress:  make(map[string]map[string]*PlayerAchievement),
	}
	
	// 注册成就
	am.registerAchievements()
	
	return am
}

// 注册成就
func (am *AchievementManager) registerAchievements() {
	// 击杀成就
	am.Register(&Achievement{
		ID:          "killer_100",
		Name:        "初出茅庐",
		Description: "累计击杀100个敌人",
		Icon:        "⚔️",
		Type:        AchievementTypeKill,
		Condition:   AchievementCondition{Type: AchievementTypeKill, Target: 100},
		Reward:      &AchievementReward{Exp: 100, Coins: 50},
	})
	
	am.Register(&Achievement{
		ID:          "killer_500",
		Name:        "小有名气",
		Description: "累计击杀500个敌人",
		Icon:        "🗡️",
		Type:        AchievementTypeKill,
		Condition:   AchievementCondition{Type: AchievementTypeKill, Target: 500},
		Reward:      &AchievementReward{Exp: 300, Coins: 150},
	})
	
	am.Register(&Achievement{
		ID:          "killer_1000",
		Name:        "无人能挡",
		Description: "累计击杀1000个敌人",
		Icon:        "⚔️",
		Type:        AchievementTypeKill,
		Condition:   AchievementCondition{Type: AchievementTypeKill, Target: 1000},
		Reward:      &AchievementReward{Exp: 500, Coins: 300, Gems: 5},
	})
	
	am.Register(&Achievement{
		ID:          "killer_5000",
		Name:        "传奇杀手",
		Description: "累计击杀5000个敌人",
		Icon:        "💀",
		Type:        AchievementTypeKill,
		Condition:   AchievementCondition{Type: AchievementTypeKill, Target: 5000},
		Reward:      &AchievementReward{Exp: 1000, Coins: 500, Gems: 20},
	})
	
	// 分数成就
	am.Register(&Achievement{
		ID:          "score_10k",
		Name:        "初试身手",
		Description: "单局获得10000分",
		Icon:        "📊",
		Type:        AchievementTypeScore,
		Condition:   AchievementCondition{Type: AchievementTypeScore, Target: 10000},
		Reward:      &AchievementReward{Exp: 100, Coins: 100},
	})
	
	am.Register(&Achievement{
		ID:          "score_50k",
		Name:        "高分选手",
		Description: "单局获得50000分",
		Icon:        "🏆",
		Type:        AchievementTypeScore,
		Condition:   AchievementCondition{Type: AchievementTypeScore, Target: 50000},
		Reward:      &AchievementReward{Exp: 300, Coins: 300, Gems: 3},
	})
	
	am.Register(&Achievement{
		ID:          "score_100k",
		Name:        "得分王者",
		Description: "单局获得100000分",
		Icon:        "👑",
		Type:        AchievementTypeScore,
		Condition:   AchievementCondition{Type: AchievementTypeScore, Target: 100000},
		Reward:      &AchievementReward{Exp: 500, Coins: 500, Gems: 10},
	})
	
	// 波次成就
	am.Register(&Achievement{
		ID:          "wave_5",
		Name:        "首战告捷",
		Description: "通过第5波",
		Icon:        "🌊",
		Type:        AchievementTypeWave,
		Condition:   AchievementCondition{Type: AchievementTypeWave, Target: 5},
		Reward:      &AchievementReward{Exp: 100, Coins: 100},
	})
	
	am.Register(&Achievement{
		ID:          "wave_10",
		Name:        "坚守阵地",
		Description: "通过第10波",
		Icon:        "🛡️",
		Type:        AchievementTypeWave,
		Condition:   AchievementCondition{Type: AchievementTypeWave, Target: 10},
		Reward:      &AchievementReward{Exp: 300, Coins: 300, Gems: 5},
	})
	
	am.Register(&Achievement{
		ID:          "wave_20",
		Name:        "不屈战士",
		Description: "通过第20波",
		Icon:        "🏰",
		Type:        AchievementTypeWave,
		Condition:   AchievementCondition{Type: AchievementTypeWave, Target: 20},
		Reward:      &AchievementReward{Exp: 500, Coins: 500, Gems: 10},
	})
	
	// 防御塔成就
	am.Register(&Achievement{
		ID:          "tower_50",
		Name:        "建筑大师",
		Description: "累计放置50座防御塔",
		Icon:        "🏗️",
		Type:        AchievementTypeTower,
		Condition:   AchievementCondition{Type: AchievementTypeTower, Target: 50},
		Reward:      &AchievementReward{Exp: 200, Coins: 200},
	})
	
	am.Register(&Achievement{
		ID:          "tower_200",
		Name:        "塔王",
		Description: "累计放置200座防御塔",
		Icon:        "🏛️",
		Type:        AchievementTypeTower,
		Condition:   AchievementCondition{Type: AchievementTypeTower, Target: 200},
		Reward:      &AchievementReward{Exp: 500, Coins: 500, Gems: 10},
	})
	
	// 礼物成就
	am.Register(&Achievement{
		ID:          "gift_100",
		Name:        "人气主播",
		Description: "收到100个礼物",
		Icon:        "🎁",
		Type:        AchievementTypeGift,
		Condition:   AchievementCondition{Type: AchievementTypeGift, Target: 100},
		Reward:      &AchievementReward{Exp: 200, Coins: 200},
	})
	
	am.Register(&Achievement{
		ID:          "gift_1000",
		Name:        "礼物达人",
		Description: "收到1000个礼物",
		Icon:        "💝",
		Type:        AchievementTypeGift,
		Condition:   AchievementCondition{Type: AchievementTypeGift, Target: 1000},
		Reward:      &AchievementReward{Exp: 500, Coins: 500, Gems: 15},
	})
	
	// 弹幕成就
	am.Register(&Achievement{
		ID:          "danmaku_500",
		Name:        "弹幕互动",
		Description: "收到500条弹幕",
		Icon:        "💬",
		Type:        AchievementTypeDanmaku,
		Condition:   AchievementCondition{Type: AchievementTypeDanmaku, Target: 500},
		Reward:      &AchievementReward{Exp: 150, Coins: 150},
	})
	
	// 特殊成就
	am.Register(&Achievement{
		ID:          "first_blood",
		Name:        "初战告捷",
		Description: "完成第一场战斗",
		Icon:        "🎖️",
		Type:        AchievementTypeSpecial,
		Condition:   AchievementCondition{Type: AchievementTypeSpecial, Target: 1},
		Reward:      &AchievementReward{Exp: 50, Coins: 50},
		Hidden:      false,
	})
	
	am.Register(&Achievement{
		ID:          "perfect_defense",
		Name:        "完美防御",
		Description: "0漏怪通过一关",
		Icon:        "🛡️",
		Type:        AchievementTypeSpecial,
		Condition:   AchievementCondition{Type: AchievementTypeSpecial, Target: 1},
		Reward:      &AchievementReward{Exp: 200, Coins: 200, Gems: 5},
		Hidden:      true,
	})
	
	am.Register(&Achievement{
		ID:          "speed_runner",
		Name:        "速通玩家",
		Description: "5分钟内通过第10波",
		Icon:        "⏱️",
		Type:        AchievementTypeTime,
		Condition:   AchievementCondition{Type: AchievementTypeTime, Target: 300}, // 300秒
		Reward:      &AchievementReward{Exp: 300, Coins: 300, Gems: 5},
		Hidden:      true,
	})
}

// 注册成就
func (am *AchievementManager) Register(achievement *Achievement) {
	am.Achievements[achievement.ID] = achievement
}

// 初始化玩家成就进度
func (am *AchievementManager) InitPlayer(playerID string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	if _, ok := am.PlayerProgress[playerID]; !ok {
		am.PlayerProgress[playerID] = make(map[string]*PlayerAchievement)
		
		// 为玩家初始化所有成就进度
		for _, achievement := range am.Achievements {
			am.PlayerProgress[playerID][achievement.ID] = &PlayerAchievement{
				AchievementID: achievement.ID,
				Progress:     0,
				Completed:    false,
			}
		}
	}
}

// 更新成就进度
func (am *AchievementManager) UpdateProgress(playerID string, achievementType AchievementType, value int) *PlayerAchievement {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	// 确保玩家已初始化
	am.initPlayerIfNeeded(playerID)
	
	// 查找匹配的成就
	var completed *PlayerAchievement
	for _, achievement := range am.Achievements {
		if achievement.Type != achievementType {
			continue
		}
		
		progress := am.PlayerProgress[playerID][achievement.ID]
		if progress.Completed || progress.RewardClaimed {
			continue
		}
		
		// 更新进度
		oldProgress := progress.Progress
		progress.Progress += value
		
		// 检查是否完成
		if progress.Progress >= achievement.Condition.Target && oldProgress < achievement.Condition.Target {
			progress.Completed = true
			progress.CompletedAt = time.Now()
			completed = progress
		}
	}
	
	return completed
}

// 初始化玩家（内部用）
func (am *AchievementManager) initPlayerIfNeeded(playerID string) {
	if _, ok := am.PlayerProgress[playerID]; !ok {
		am.PlayerProgress[playerID] = make(map[string]*PlayerAchievement)
		for _, achievement := range am.Achievements {
			am.PlayerProgress[playerID][achievement.ID] = &PlayerAchievement{
				AchievementID: achievement.ID,
				Progress:     0,
				Completed:    false,
			}
		}
	}
}

// 获取玩家成就列表
func (am *AchievementManager) GetPlayerAchievements(playerID string) []*PlayerAchievement {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	am.initPlayerIfNeeded(playerID)
	
	result := make([]*PlayerAchievement, 0, len(am.PlayerProgress[playerID]))
	for _, pa := range am.PlayerProgress[playerID] {
		result = append(result, pa)
	}
	
	return result
}

// 获取玩家已完成成就数
func (am *AchievementManager) GetCompletedCount(playerID string) int {
	am.mu.RLock()
	defer am.mu.RUnlock()
	
	am.initPlayerIfNeeded(playerID)
	
	count := 0
	for _, pa := range am.PlayerProgress[playerID] {
		if pa.Completed {
			count++
		}
	}
	
	return count
}

// 领取成就奖励
func (am *AchievementManager) ClaimReward(playerID, achievementID string) (*AchievementReward, error) {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	am.initPlayerIfNeeded(playerID)
	
	progress, ok := am.PlayerProgress[playerID][achievementID]
	if !ok {
		return nil, fmt.Errorf("achievement not found")
	}
	
	if !progress.Completed {
		return nil, fmt.Errorf("achievement not completed")
	}
	
	if progress.RewardClaimed {
		return nil, fmt.Errorf("reward already claimed")
	}
	
	achievement, ok := am.Achievements[achievementID]
	if !ok {
		return nil, fmt.Errorf("achievement definition not found")
	}
	
	progress.RewardClaimed = true
	
	return achievement.Reward, nil
}

// 获取成就定义
func (am *AchievementManager) GetAchievement(id string) *Achievement {
	return am.Achievements[id]
}

// 获取所有成就
func (am *AchievementManager) GetAllAchievements() []*Achievement {
	result := make([]*Achievement, 0, len(am.Achievements))
	for _, a := range am.Achievements {
		result = append(result, a)
	}
	return result
}

// ============================================
// 任务系统 (Quest System)
// ============================================

// 任务类型
type QuestType int

const (
	QuestTypeDaily QuestType = iota // 每日任务
	QuestTypeWeekly                 // 每周任务
	QuestTypeAchievement            // 成就任务
	QuestTypeMainLine               // 主线任务
)

// 任务状态
type QuestStatus int

const (
	QuestStatusLocked QuestStatus = iota // 锁定
	QuestStatusAvailable                 // 可接
	QuestStatusInProgress                // 进行中
	QuestStatusCompleted                  // 已完成
	QuestStatusRewarded                   // 已领取奖励
)

// 任务目标
type QuestTarget struct {
	Type      string `json:"type"`      // kill_enemy/place_tower/get_score/send_gift
	Target    int    `json:"target"`   // 目标数量
	Progress  int    `json:"progress"` // 当前进度
}

// 任务定义
type Quest struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        QuestType   `json:"type"`
	Targets     []QuestTarget `json:"targets"`
	Reward      *AchievementReward `json:"reward"`
	ExpiresAt   time.Time  `json:"expiresAt"` // 过期时间
}

// 玩家任务进度
type PlayerQuest struct {
	QuestID   string     `json:"questId"`
	Status    QuestStatus `json:"status"`
	Progress  int        `json:"progress"`
	StartedAt time.Time  `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
}

// 任务管理器
type QuestManager struct {
	// 任务定义
	Quests map[string]*Quest
	
	// 玩家任务进度
	PlayerQuests map[string]map[string]*PlayerQuest // playerID -> questID -> progress
	
	mu sync.RWMutex
}

// 新建任务管理器
func NewQuestManager() *QuestManager {
	qm := &QuestManager{
		Quests:       make(map[string]*Quest),
		PlayerQuests: make(map[string]map[string]*PlayerQuest),
	}
	
	// 注册任务
	qm.registerQuests()
	
	return qm
}

// 注册任务
func (qm *QuestManager) registerQuests() {
	// 每日任务
	qm.Register(&Quest{
		ID:          "daily_kill_10",
		Name:        "日常击杀",
		Description: "击杀10个敌人",
		Type:        QuestTypeDaily,
		Targets:     []QuestTarget{{Type: "kill_enemy", Target: 10}},
		Reward:      &AchievementReward{Exp: 50, Coins: 50},
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	})
	
	qm.Register(&Quest{
		ID:          "daily_score_1000",
		Name:        "日常得分",
		Description: "获得1000分",
		Type:        QuestTypeDaily,
		Targets:     []QuestTarget{{Type: "get_score", Target: 1000}},
		Reward:      &AchievementReward{Exp: 50, Coins: 50},
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	})
	
	qm.Register(&Quest{
		ID:          "daily_place_5",
		Name:        "日常建造",
		Description: "放置5座防御塔",
		Type:        QuestTypeDaily,
		Targets:     []QuestTarget{{Type: "place_tower", Target: 5}},
		Reward:      &AchievementReward{Exp: 50, Coins: 50},
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	})
	
	qm.Register(&Quest{
		ID:          "daily_wave_3",
		Name:        "日常推图",
		Description: "通过第3波",
		Type:        QuestTypeDaily,
		Targets:     []QuestTarget{{Type: "pass_wave", Target: 3}},
		Reward:      &AchievementReward{Exp: 100, Coins: 100},
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	})
	
	// 每周任务
	qm.Register(&Quest{
		ID:          "weekly_kill_100",
		Name:        "周常击杀",
		Description: "击杀100个敌人",
		Type:        QuestTypeWeekly,
		Targets:     []QuestTarget{{Type: "kill_enemy", Target: 100}},
		Reward:      &AchievementReward{Exp: 300, Coins: 300, Gems: 3},
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
	})
	
	qm.Register(&Quest{
		ID:          "weekly_score_10k",
		Name:        "周常得分",
		Description: "获得10000分",
		Type:        QuestTypeWeekly,
		Targets:     []QuestTarget{{Type: "get_score", Target: 10000}},
		Reward:      &AchievementReward{Exp: 300, Coins: 300, Gems: 3},
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour),
	})
	
	// 主线任务
	qm.Register(&Quest{
		ID:          "main_wave_1",
		Name:        "初入战场",
		Description: "通过第1波",
		Type:        QuestTypeMainLine,
		Targets:     []QuestTarget{{Type: "pass_wave", Target: 1}},
		Reward:      &AchievementReward{Exp: 100, Coins: 100},
	})
	
	qm.Register(&Quest{
		ID:          "main_wave_5",
		Name:        "小试牛刀",
		Description: "通过第5波",
		Type:        QuestTypeMainLine,
		Targets:     []QuestTarget{{Type: "pass_wave", Target: 5}},
		Reward:      &AchievementReward{Exp: 200, Coins: 200},
	})
	
	qm.Register(&Quest{
		ID:          "main_wave_10",
		Name:        "锋芒毕露",
		Description: "通过第10波",
		Type:        QuestTypeMainLine,
		Targets:     []QuestTarget{{Type: "pass_wave", Target: 10}},
		Reward:      &AchievementReward{Exp: 500, Coins: 500, Gems: 5},
	})
	
	qm.Register(&Quest{
		ID:          "main_wave_20",
		Name:        "无人能挡",
		Description: "通过第20波",
		Type:        QuestTypeMainLine,
		Targets:     []QuestTarget{{Type: "pass_wave", Target: 20}},
		Reward:      &AchievementReward{Exp: 1000, Coins: 1000, Gems: 10},
	})
}

// 注册任务
func (qm *QuestManager) Register(quest *Quest) {
	qm.Quests[quest.ID] = quest
}

// 初始化玩家任务
func (qm *QuestManager) InitPlayer(playerID string) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	
	if _, ok := qm.PlayerQuests[playerID]; !ok {
		qm.PlayerQuests[playerID] = make(map[string]*PlayerQuest)
		
		// 初始化所有任务
		for _, quest := range qm.Quests {
			qm.PlayerQuests[playerID][quest.ID] = &PlayerQuest{
				QuestID:  quest.ID,
				Status:   QuestStatusAvailable,
				Progress: 0,
			}
		}
	}
}

// 接受任务
func (qm *QuestManager) AcceptQuest(playerID, questID string) error {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	
	qm.initPlayerIfNeeded(playerID)
	
	pq, ok := qm.PlayerQuests[playerID][questID]
	if !ok {
		return fmt.Errorf("quest not found")
	}
	
	if pq.Status != QuestStatusAvailable {
		return fmt.Errorf("quest not available")
	}
	
	pq.Status = QuestStatusInProgress
	pq.StartedAt = time.Now()
	
	return nil
}

// 更新任务进度
func (qm *QuestManager) UpdateProgress(playerID, targetType string, value int) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	
	qm.initPlayerIfNeeded(playerID)
	
	// 查找匹配的任务
	for questID, pq := range qm.PlayerQuests[playerID] {
		if pq.Status != QuestStatusInProgress {
			continue
		}
		
		quest, ok := qm.Quests[questID]
		if !ok {
			continue
		}
		
		// 检查任务目标类型
		for i, target := range quest.Targets {
			if target.Type == targetType {
				oldProgress := pq.Progress
				pq.Progress += value
				
				// 更新任务目标进度
				quest.Targets[i].Progress = pq.Progress
				
				// 检查是否完成
				if pq.Progress >= target.Target && oldProgress < target.Target {
					pq.Status = QuestStatusCompleted
					pq.CompletedAt = time.Now()
				}
				
				break
			}
		}
	}
}

// 完成任务并领取奖励
func (qm *QuestManager) CompleteQuest(playerID, questID string) (*AchievementReward, error) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	
	qm.initPlayerIfNeeded(playerID)
	
	pq, ok := qm.PlayerQuests[playerID][questID]
	if !ok {
		return nil, fmt.Errorf("quest not found")
	}
	
	if pq.Status != QuestStatusCompleted {
		return nil, fmt.Errorf("quest not completed")
	}
	
	quest, ok := qm.Quests[questID]
	if !ok {
		return nil, fmt.Errorf("quest definition not found")
	}
	
	pq.Status = QuestStatusRewarded
	
	return quest.Reward, nil
}

// 获取玩家任务列表
func (qm *QuestManager) GetPlayerQuests(playerID string) []*PlayerQuest {
	qm.mu.RLock()
	defer qm.mu.RUnlock()
	
	qm.initPlayerIfNeeded(playerID)
	
	result := make([]*PlayerQuest, 0, len(qm.PlayerQuests[playerID]))
	for _, pq := range qm.PlayerQuests[playerID] {
		result = append(result, pq)
	}
	
	return result
}

// 获取每日任务
func (qm *QuestManager) GetDailyQuests(playerID string) []*Quest {
	qm.mu.RLock()
	defer qm.mu.RUnlock()
	
	qm.initPlayerIfNeeded(playerID)
	
	result := make([]*Quest, 0)
	for _, quest := range qm.Quests {
		if quest.Type == QuestTypeDaily {
			result = append(result, quest)
		}
	}
	
	return result
}

// 初始化玩家（内部用）
func (qm *QuestManager) initPlayerIfNeeded(playerID string) {
	if _, ok := qm.PlayerQuests[playerID]; !ok {
		qm.PlayerQuests[playerID] = make(map[string]*PlayerQuest)
		for _, quest := range qm.Quests {
			qm.PlayerQuests[playerID][quest.ID] = &PlayerQuest{
				QuestID:  quest.ID,
				Status:   QuestStatusAvailable,
				Progress: 0,
			}
		}
	}
}

// ============================================
// 玩家统计系统 (Player Statistics)
// ============================================

// 玩家统计数据
type PlayerStats struct {
	PlayerID string `json:"playerId"`
	
	// 战斗统计
	TotalGames    int `json:"totalGames"`    // 总场次
	WinGames      int `json:"winGames"`      // 胜利场次
	LoseGames     int `json:"loseGames"`     // 失败场次
	WinRate       float64 `json:"winRate"`   // 胜率
	
	// 击杀统计
	TotalKills   int `json:"totalKills"`    // 总击杀
	MaxKills      int `json:"maxKills"`      // 最高击杀
	AvgKills      float64 `json:"avgKills"`  // 平均击杀
	
	// 波次统计
	MaxWave       int `json:"maxWave"`       // 最高波次
	TotalWaves    int `json:"totalWaves"`   // 累计波次
	
	// 分数统计
	TotalScore    int `json:"totalScore"`   // 总得分
	MaxScore      int `json:"maxScore"`     // 最高得分
	AvgScore      float64 `json:"avgScore"` // 平均得分
	
	// 经济统计
	TotalMoneyEarned int `json:"totalMoneyEarned"` // 累计获得金币
	TotalMoneySpent int `json:"totalMoneySpent"`  // 累计花费金币
	
	// 塔统计
	TowersPlaced   int `json:"towersPlaced"`   // 放置塔数
	TowersUpgraded int `json:"towersUpgraded"` // 升级塔数
	TowersSold     int `json:"towersSold"`     // 出售塔数
	
	// 时间统计
	TotalPlayTime int64 `json:"totalPlayTime"` // 总游玩时间(秒)
	FirstPlayTime time.Time `json:"firstPlayTime"` // 首次游玩时间
	LastPlayTime  time.Time `json:"lastPlayTime"`  // 最后游玩时间
	
	// 社交统计
	TotalGiftsReceived int `json:"totalGiftsReceived"` // 收到礼物
	TotalDanmakuReceived int `json:"totalDanmakuReceived"` // 收到弹幕
	
	mu sync.RWMutex
}

// 玩家统计管理器
type StatsManager struct {
	Stats map[string]*PlayerStats // playerID -> stats
	mu    sync.RWMutex
}

// 新建统计管理器
func NewStatsManager() *StatsManager {
	return &StatsManager{
		Stats: make(map[string]*PlayerStats),
	}
}

// 获取玩家统计
func (sm *StatsManager) GetStats(playerID string) *PlayerStats {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	if stats, ok := sm.Stats[playerID]; ok {
		return stats
	}
	
	// 返回默认统计
	return &PlayerStats{
		PlayerID:     playerID,
		FirstPlayTime: time.Now(),
	}
}

// 获取或创建玩家统计
func (sm *StatsManager) GetOrCreateStats(playerID string) *PlayerStats {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if stats, ok := sm.Stats[playerID]; ok {
		return stats
	}
	
	stats := &PlayerStats{
		PlayerID:     playerID,
		FirstPlayTime: time.Now(),
		LastPlayTime:  time.Now(),
	}
	sm.Stats[playerID] = stats
	
	return stats
}

// 记录游戏开始
func (sm *StatsManager) RecordGameStart(playerID string) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.TotalGames++
	stats.LastPlayTime = time.Now()
}

// 记录游戏胜利
func (sm *StatsManager) RecordGameWin(playerID string, kills, score, wave int) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.WinGames++
	if stats.TotalGames > 0 {
		stats.WinRate = float64(stats.WinGames) / float64(stats.TotalGames) * 100
	}
	
	// 更新击杀
	stats.TotalKills += kills
	if kills > stats.MaxKills {
		stats.MaxKills = kills
	}
	stats.AvgKills = float64(stats.TotalKills) / float64(stats.TotalGames)
	
	// 更新波次
	if wave > stats.MaxWave {
		stats.MaxWave = wave
	}
	stats.TotalWaves += wave
	
	// 更新分数
	stats.TotalScore += score
	if score > stats.MaxScore {
		stats.MaxScore = score
	}
	stats.AvgScore = float64(stats.TotalScore) / float64(stats.TotalGames)
}

// 记录游戏失败
func (sm *StatsManager) RecordGameLose(playerID string, kills, score, wave int) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.LoseGames++
	if stats.TotalGames > 0 {
		stats.WinRate = float64(stats.WinGames) / float64(stats.TotalGames) * 100
	}
	
	// 更新击杀
	stats.TotalKills += kills
	if kills > stats.MaxKills {
		stats.MaxKills = kills
	}
	stats.AvgKills = float64(stats.TotalKills) / float64(stats.TotalGames)
	
	// 更新波次
	if wave > stats.MaxWave {
		stats.MaxWave = wave
	}
	stats.TotalWaves += wave
	
	// 更新分数
	stats.TotalScore += score
	if score > stats.MaxScore {
		stats.MaxScore = score
	}
	stats.AvgScore = float64(stats.TotalScore) / float64(stats.TotalGames)
}

// 记录塔操作
func (sm *StatsManager) RecordTowerPlaced(playerID string) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.TowersPlaced++
}

// 记录塔升级
func (sm *StatsManager) RecordTowerUpgraded(playerID string) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.TowersUpgraded++
}

// 记录塔出售
func (sm *StatsManager) RecordTowerSold(playerID string) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.TowersSold++
}

// 记录金币变化
func (sm *StatsManager) RecordMoneyChange(playerID string, earned, spent int) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.TotalMoneyEarned += earned
	stats.TotalMoneySpent += spent
}

// 记录收到礼物
func (sm *StatsManager) RecordGiftReceived(playerID string, count int) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.TotalGiftsReceived += count
}

// 记录收到弹幕
func (sm *StatsManager) RecordDanmakuReceived(playerID string, count int) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.TotalDanmakuReceived += count
}

// 更新游戏时长
func (sm *StatsManager) UpdatePlayTime(playerID string, seconds int64) {
	stats := sm.GetOrCreateStats(playerID)
	stats.mu.Lock()
	defer stats.mu.Unlock()
	
	stats.TotalPlayTime += seconds
	stats.LastPlayTime = time.Now()
}

// ============================================
// 全局管理器
// ============================================

var (
	GlobalAchievementManager *AchievementManager
	GlobalQuestManager       *QuestManager
	GlobalStatsManager       *StatsManager
)

func init() {
	GlobalAchievementManager = NewAchievementManager()
	GlobalQuestManager = NewQuestManager()
	GlobalStatsManager = NewStatsManager()
}
