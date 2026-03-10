package game

import (
	"fmt"
	"sync"
	"time"
)

// ==================== 任务系统 ====================

// QuestCategory 任务分类
type QuestCategory int

const (
	QuestCategoryMain     QuestCategory = iota // 主线任务
	QuestCategoryDaily                         // 每日任务
	QuestCategoryWeekly                        // 每周任务
	QuestCategoryAchievement                   // 成就任务
	QuestCategoryEvent                         // 活动任务
	QuestCategorySide                          // 支线任务
)

// QuestType 任务类型
type QuestType int

const (
	QuestTypeCollect  QuestType = iota // 收集
	QuestTypeDefeat                   // 击败
	QuestTypeSurvive                  // 生存
	QuestTypeScore                    // 分数
	QuestTypeCombo                    // 连击
	QuestTypeTime                     // 时间
	QuestTypeWin                      // 胜利
	QuestTypeDraw                     // 抽卡
	QuestTypeUpgrade                  // 升级
	QuestTypeFriend                   // 好友
	QuestTypeGuild                    // 公会
	QuestTypeCustom                   // 自定义
)

// QuestTarget 任务目标
type QuestTarget struct {
	Type      QuestType `json:"type"`       // 目标类型
	TargetID  string    `json:"target_id"`  // 目标ID
	Count     int       `json:"count"`      // 需要数量
	Progress  int       `json:"progress"`   // 当前进度
	Completed bool      `json:"completed"`  // 是否完成
}

// Quest 任务
type Quest struct {
	ID          string        `json:"id"`           // 任务ID
	Name        string        `json:"name"`        // 任务名称
	Description string        `json:"description"` // 任务描述
	Category    QuestCategory `json:"category"`    // 任务分类
	Type        QuestType     `json:"type"`        // 任务类型
	Targets     []QuestTarget `json:"targets"`     // 任务目标
	Rewards     QuestRewards  `json:"rewards"`     // 任务奖励
	TimeLimit   int           `json:"time_limit"`  // 时间限制(秒), 0表示无限制
	LevelReq    int           `json:"level_req"`   // 等级要求
	Priority    int           `json:"priority"`     // 优先级
	AutoClaim   bool          `json:"auto_claim"`   // 是否自动领取
	Visible     bool          `json:"visible"`      // 是否可见
	CreatedAt   int64         `json:"created_at"`   // 创建时间
	UpdatedAt   int64         `json:"updated_at"`   // 更新时间
}

// QuestRewards 任务奖励
type QuestRewards struct {
	Coins  int            `json:"coins"`  // 金币
	Gems   int            `json:"gems"`    // 钻石
	Exp    int            `json:"exp"`     // 经验
	Items  map[string]int `json:"items"`   // 道具 {item_id: count}
	Title  string         `json:"title"`   // 称号
}

// PlayerQuest 玩家任务进度
type PlayerQuest struct {
	PlayerID     string        `json:"player_id"`      // 玩家ID
	QuestID      string        `json:"quest_id"`       // 任务ID
	Status       QuestStatus   `json:"status"`         // 任务状态
	Progress     []QuestTarget `json:"progress"`       // 进度
	StartTime    int64         `json:"start_time"`     // 开始时间
	CompleteTime int64         `json:"complete_time"`  // 完成时间
	Claimed      bool          `json:"claimed"`         // 是否已领取奖励
}

// QuestStatus 任务状态
type QuestStatus int

const (
	QuestStatusLocked  QuestStatus = iota // 锁定
	QuestStatusAvailable                   // 可接取
	QuestStatusInProgress                  // 进行中
	QuestStatusCompleted                   // 已完成(未领取)
	QuestStatusClaimed                      // 已领取
	QuestStatusExpired                      // 已过期
)

// QuestManager 任务管理器
type QuestManager struct {
	quests map[string]*Quest // 所有任务
	mu     sync.RWMutex
}

// NewQuestManager 创建任务管理器
func NewQuestManager() *QuestManager {
	return &QuestManager{
		quests: make(map[string]*Quest),
	}
}

// RegisterQuest 注册任务
func (qm *QuestManager) RegisterQuest(quest *Quest) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	quest.CreatedAt = time.Now().Unix()
	quest.UpdatedAt = time.Now().Unix()
	qm.quests[quest.ID] = quest
}

// GetQuest 获取任务
func (qm *QuestManager) GetQuest(id string) (*Quest, bool) {
	qm.mu.RLock()
	defer qm.mu.RUnlock()
	quest, ok := qm.quests[id]
	return quest, ok
}

// GetQuestsByCategory 按分类获取任务
func (qm *QuestManager) GetQuestsByCategory(category QuestCategory) []*Quest {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	var result []*Quest
	for _, quest := range qm.quests {
		if quest.Category == category && quest.Visible {
			result = append(result, quest)
		}
	}
	return result
}

// UpdateQuestProgress 更新任务进度
func (qm *QuestManager) UpdateQuestProgress(playerID string, questID string, questType QuestType, targetID string, count int, playerQuests map[string]*PlayerQuest) bool {
	qm.mu.RLock()
	quest, ok := qm.quests[questID]
	qm.mu.RUnlock()

	if !ok {
		return false
	}

	playerQuest, ok := playerQuests[questID]
	if !ok || playerQuest.Status != QuestStatusInProgress {
		return false
	}

	updated := false
	for i, target := range playerQuest.Progress {
		if target.Type == questType && (targetID == "" || target.TargetID == targetID) {
			playerQuest.Progress[i].Progress += count
			if playerQuest.Progress[i].Progress >= playerQuest.Progress[i].Count {
				playerQuest.Progress[i].Completed = true
			}
			updated = true
		}
	}

	// 检查是否全部完成
	allCompleted := true
	for _, target := range playerQuest.Progress {
		if !target.Completed {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		playerQuest.Status = QuestStatusCompleted
		playerQuest.CompleteTime = time.Now().Unix()
	}

	return updated
}

// ClaimQuestReward 领取任务奖励
func (qm *QuestManager) ClaimQuestReward(player *Player, questID string) (*QuestRewards, error) {
	qm.mu.RLock()
	quest, ok := qm.quests[questID]
	qm.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("quest not found: %s", questID)
	}

	playerQuest := player.GetQuest(questID)
	if playerQuest == nil {
		return nil, fmt.Errorf("player quest not found")
	}

	if playerQuest.Status != QuestStatusCompleted {
		return nil, fmt.Errorf("quest not completed")
	}

	if playerQuest.Claimed {
		return nil, fmt.Errorf("reward already claimed")
	}

	// 发放奖励
	player.AddCoins(quest.Rewards.Coins)
	player.AddGems(quest.Rewards.Gems)
	player.AddExp(quest.Rewards.Exp)

	for itemID, count := range quest.Rewards.Items {
		player.AddItem(itemID, count)
	}

	playerQuest.Claimed = true
	playerQuest.Status = QuestStatusClaimed

	return &quest.Rewards, nil
}

// InitDefaultQuests 初始化默认任务
func (qm *QuestManager) InitDefaultQuests() {
	// 主线任务
	mainQuests := []*Quest{
		{
			ID:          "main_001",
			Name:        "初入战场",
			Description: "完成第一关",
			Category:    QuestCategoryMain,
			Type:        QuestTypeWin,
			Targets:     []QuestTarget{{Type: QuestTypeWin, TargetID: "level_1_1", Count: 1}},
			Rewards:     QuestRewards{Coins: 100, Gems: 10, Exp: 50},
			LevelReq:    1,
			Priority:    1,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "main_002",
			Name:        "小试牛刀",
			Description: "击败10个敌人",
			Category:    QuestCategoryMain,
			Type:        QuestTypeDefeat,
			Targets:     []QuestTarget{{Type: QuestTypeDefeat, TargetID: "", Count: 10}},
			Rewards:     QuestRewards{Coins: 200, Gems: 20, Exp: 100},
			LevelReq:    1,
			Priority:    2,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "main_003",
			Name:        "收集达人",
			Description: "收集50个金币",
			Category:    QuestCategoryMain,
			Type:        QuestTypeCollect,
			Targets:     []QuestTarget{{Type: QuestTypeCollect, TargetID: "coin", Count: 50}},
			Rewards:     QuestRewards{Coins: 300, Gems: 30, Exp: 150},
			LevelReq:    2,
			Priority:    3,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "main_004",
			Name:        "连胜勇士",
			Description: "达成100连击",
			Category:    QuestCategoryMain,
			Type:        QuestTypeCombo,
			Targets:     []QuestTarget{{Type: QuestTypeCombo, TargetID: "", Count: 100}},
			Rewards:     QuestRewards{Coins: 500, Gems: 50, Exp: 300},
			LevelReq:    3,
			Priority:    4,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "main_005",
			Name:        "Boss杀手",
			Description: "击败Boss",
			Category:    QuestCategoryMain,
			Type:        QuestTypeDefeat,
			Targets:     []QuestTarget{{Type: QuestTypeDefeat, TargetID: "boss", Count: 1}},
			Rewards:     QuestRewards{Coins: 1000, Gems: 100, Exp: 500, Items: map[string]int{"weapon_ssr": 1}},
			LevelReq:    5,
			Priority:    5,
			AutoClaim:   true,
			Visible:     true,
		},
	}

	// 每日任务
	dailyQuests := []*Quest{
		{
			ID:          "daily_001",
			Name:        "日常击杀",
			Description: "击败5个敌人",
			Category:    QuestCategoryDaily,
			Type:        QuestTypeDefeat,
			Targets:     []QuestTarget{{Type: QuestTypeDefeat, TargetID: "", Count: 5}},
			Rewards:     QuestRewards{Coins: 50, Gems: 5},
			TimeLimit:   86400,
			Priority:    1,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "daily_002",
			Name:        "日常胜利",
			Description: "完成1场战斗",
			Category:    QuestCategoryDaily,
			Type:        QuestTypeWin,
			Targets:     []QuestTarget{{Type: QuestTypeWin, TargetID: "", Count: 1}},
			Rewards:     QuestRewards{Coins: 100, Gems: 10},
			TimeLimit:   86400,
			Priority:    2,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "daily_003",
			Name:        "日常连击",
			Description: "达成20连击",
			Category:    QuestCategoryDaily,
			Type:        QuestTypeCombo,
			Targets:     []QuestTarget{{Type: QuestTypeCombo, TargetID: "", Count: 20}},
			Rewards:     QuestRewards{Coins: 80, Gems: 8},
			TimeLimit:   86400,
			Priority:    3,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "daily_004",
			Name:        "日常收集",
			Description: "收集20个金币",
			Category:    QuestCategoryDaily,
			Type:        QuestTypeCollect,
			Targets:     []QuestTarget{{Type: QuestTypeCollect, TargetID: "coin", Count: 20}},
			Rewards:     QuestRewards{Coins: 60, Gems: 6},
			TimeLimit:   86400,
			Priority:    4,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "daily_005",
			Name:        "日常生存",
			Description: "生存30秒",
			Category:    QuestCategoryDaily,
			Type:        QuestTypeSurvive,
			Targets:     []QuestTarget{{Type: QuestTypeSurvive, TargetID: "", Count: 30}},
			Rewards:     QuestRewards{Coins: 70, Gems: 7},
			TimeLimit:   86400,
			Priority:    5,
			AutoClaim:   true,
			Visible:     true,
		},
	}

	// 每周任务
	weeklyQuests := []*Quest{
		{
			ID:          "weekly_001",
			Name:        "周冠军",
			Description: "累计获得10000分",
			Category:    QuestCategoryWeekly,
			Type:        QuestTypeScore,
			Targets:     []QuestTarget{{Type: QuestTypeScore, TargetID: "", Count: 10000}},
			Rewards:     QuestRewards{Coins: 1000, Gems: 100, Exp: 500},
			TimeLimit:   604800,
			Priority:    1,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "weekly_002",
			Name:        "周连胜",
			Description: "累计50连击",
			Category:    QuestCategoryWeekly,
			Type:        QuestTypeCombo,
			Targets:     []QuestTarget{{Type: QuestTypeCombo, TargetID: "", Count: 50}},
			Rewards:     QuestRewards{Coins: 800, Gems: 80, Exp: 400},
			TimeLimit:   604800,
			Priority:    2,
			AutoClaim:   true,
			Visible:     true,
		},
		{
			ID:          "weekly_003",
			Name:        "周末BOSS",
			Description: "击败3个Boss",
			Category:    QuestCategoryWeekly,
			Type:        QuestTypeDefeat,
			Targets:     []QuestTarget{{Type: QuestTypeDefeat, TargetID: "boss", Count: 3}},
			Rewards:     QuestRewards{Coins: 1500, Gems: 150, Exp: 800, Items: map[string]int{"armor_rare": 1}},
			TimeLimit:   604800,
			Priority:    3,
			AutoClaim:   true,
			Visible:     true,
		},
	}

	for _, quest := range mainQuests {
		qm.RegisterQuest(quest)
	}
	for _, quest := range dailyQuests {
		qm.RegisterQuest(quest)
	}
	for _, quest := range weeklyQuests {
		qm.RegisterQuest(quest)
	}
}

// ==================== 成就系统 ====================

// AchievementCategory 成就分类
type AchievementCategory int

const (
	AchievementCategoryCombat     AchievementCategory = iota // 战斗成就
	AchievementCategoryCollection                           // 收集成就
	AchievementCategorySocial                               // 社交成就
	AchievementCategoryTime                                 // 时间成就
	AchievementCategorySpecial                              // 特殊成就
)

// Achievement 成就
type Achievement struct {
	ID          string             `json:"id"`            // 成就ID
	Name        string             `json:"name"`          // 成就名称
	Description string             `json:"description"`  // 成就描述
	Category    AchievementCategory `json:"category"`    // 成就分类
	Condition   AchievementCondition `json:"condition"`   // 解锁条件
	Rewards     AchievementRewards  `json:"rewards"`      // 成就奖励
	Icon        string             `json:"icon"`          // 图标
	Rarity      Rarity             `json:"rarity"`        // 稀有度
	SortOrder   int                `json:"sort_order"`   // 排序
	Visible     bool               `json:"visible"`       // 是否可见
}

// AchievementCondition 成就条件
type AchievementCondition struct {
	Type      string `json:"type"`       // 条件类型
	TargetID  string `json:"target_id"`  // 目标ID
	Threshold int    `json:"threshold"`  // 阈值
}

// AchievementRewards 成就奖励
type AchievementRewards struct {
	Coins   int            `json:"coins"`    // 金币
	Gems    int            `json:"gems"`     // 钻石
	Title   string         `json:"title"`   // 称号
	Badge   string         `json:"badge"`   // 徽章
	Frame   string         `json:"frame"`   // 头像框
	Items   map[string]int `json:"items"`    // 道具
}

// PlayerAchievement 玩家成就进度
type PlayerAchievement struct {
	PlayerID      string `json:"player_id"`       // 玩家ID
	AchievementID string `json:"achievement_id"`   // 成就ID
	Progress      int    `json:"progress"`         // 当前进度
	Unlocked      bool   `json:"unlocked"`         // 是否解锁
	UnlockTime    int64  `json:"unlock_time"`      // 解锁时间
	Claimed       bool   `json:"claimed"`          // 是否已领取奖励
}

// AchievementManager 成就管理器
type AchievementManager struct {
	achievements map[string]*Achievement // 所有成就
	mu           sync.RWMutex
}

// NewAchievementManager 创建成就管理器
func NewAchievementManager() *AchievementManager {
	return &AchievementManager{
		achievements: make(map[string]*Achievement),
	}
}

// RegisterAchievement 注册成就
func (am *AchievementManager) RegisterAchievement(ach *Achievement) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.achievements[ach.ID] = ach
}

// GetAchievement 获取成就
func (am *AchievementManager) GetAchievement(id string) (*Achievement, bool) {
	am.mu.RLock()
	defer am.mu.RUnlock()
	ach, ok := am.achievements[id]
	return ach, ok
}

// GetAchievementsByCategory 按分类获取成就
func (am *AchievementManager) GetAchievementsByCategory(category AchievementCategory) []*Achievement {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var result []*Achievement
	for _, ach := range am.achievements {
		if ach.Category == category && ach.Visible {
			result = append(result, ach)
		}
	}
	return result
}

// UpdateAchievementProgress 更新成就进度
func (am *AchievementManager) UpdateAchievementProgress(player *Player, achievementID string, progress int) bool {
	ach, ok := am.GetAchievement(achievementID)
	if !ok {
		return false
	}

	playerAch := player.GetAchievement(achievementID)
	if playerAch == nil {
		playerAch = &PlayerAchievement{
			PlayerID:      player.ID,
			AchievementID: achievementID,
			Progress:      0,
			Unlocked:      false,
		}
		player.Achievements = append(player.Achievements, playerAch)
	}

	oldProgress := playerAch.Progress
	playerAch.Progress = progress

	// 检查是否解锁
	if !playerAch.Unlocked && progress >= ach.Condition.Threshold {
		playerAch.Unlocked = true
		playerAch.UnlockTime = time.Now().Unix()
	}

	return playerAch.Progress != oldProgress
}

// IncrementAchievementProgress 增加成就进度
func (am *AchievementManager) IncrementAchievementProgress(player *Player, achievementID string, increment int) bool {
	playerAch := player.GetAchievement(achievementID)
	if playerAch == nil {
		return am.UpdateAchievementProgress(player, achievementID, increment)
	}

	ach, ok := am.GetAchievement(achievementID)
	if !ok {
		return false
	}

	newProgress := playerAch.Progress + increment
	if newProgress > ach.Condition.Threshold {
		newProgress = ach.Condition.Threshold
	}

	return am.UpdateAchievementProgress(player, achievementID, newProgress)
}

// ClaimAchievementReward 领取成就奖励
func (am *AchievementManager) ClaimAchievementReward(player *Player, achievementID string) (*AchievementRewards, error) {
	ach, ok := am.GetAchievement(achievementID)
	if !ok {
		return nil, fmt.Errorf("achievement not found: %s", achievementID)
	}

	playerAch := player.GetAchievement(achievementID)
	if playerAch == nil {
		return nil, fmt.Errorf("player achievement not found")
	}

	if !playerAch.Unlocked {
		return nil, fmt.Errorf("achievement not unlocked")
	}

	if playerAch.Claimed {
		return nil, fmt.Errorf("reward already claimed")
	}

	// 发放奖励
	player.AddCoins(ach.Rewards.Coins)
	player.AddGems(ach.Rewards.Gems)

	if ach.Rewards.Title != "" {
		player.AddTitle(ach.Rewards.Title)
	}

	if ach.Rewards.Badge != "" {
		player.AddBadge(ach.Rewards.Badge)
	}

	if ach.Rewards.Frame != "" {
		player.AddAvatarFrame(ach.Rewards.Frame)
	}

	for itemID, count := range ach.Rewards.Items {
		player.AddItem(itemID, count)
	}

	playerAch.Claimed = true

	return &ach.Rewards, nil
}

// InitDefaultAchievements 初始化默认成就
func (am *AchievementManager) InitDefaultAchievements() {
	// 战斗成就
	combatAchievements := []*Achievement{
		{
			ID:          "combat_first_blood",
			Name:        "初战告捷",
			Description: "完成第一场战斗",
			Category:    AchievementCategoryCombat,
			Condition:   AchievementCondition{Type: "win_count", Threshold: 1},
			Rewards:     AchievementRewards{Coins: 100, Gems: 10},
			Icon:        "⚔️",
			Rarity:      RarityCommon,
			SortOrder:   1,
			Visible:     true,
		},
		{
			ID:          "combat_100_kills",
			Name:        "百人斩",
			Description: "累计击败100个敌人",
			Category:    AchievementCategoryCombat,
			Condition:   AchievementCondition{Type: "kill_count", Threshold: 100},
			Rewards:     AchievementRewards{Coins: 500, Gems: 50, Badge: "kill_100"},
			Icon:        "🗡️",
			Rarity:      RarityUncommon,
			SortOrder:   2,
			Visible:     true,
		},
		{
			ID:          "combat_1000_kills",
			Name:        "千人斩",
			Description: "累计击败1000个敌人",
			Category:    AchievementCategoryCombat,
			Condition:   AchievementCondition{Type: "kill_count", Threshold: 1000},
			Rewards:     AchievementRewards{Coins: 2000, Gems: 200, Badge: "kill_1000", Title: "千人斩"},
			Icon:        "⚔️",
			Rarity:      RarityEpic,
			SortOrder:   3,
			Visible:     true,
		},
		{
			ID:          "combat_boss_killer",
			Name:        "Boss杀手",
			Description: "累计击败10个Boss",
			Category:    AchievementCategoryCombat,
			Condition:   AchievementCondition{Type: "boss_kill_count", Threshold: 10},
			Rewards:     AchievementRewards{Coins: 3000, Gems: 300, Badge: "boss_killer", Title: "Boss杀手"},
			Icon:        "👹",
			Rarity:      RarityLegendary,
			SortOrder:   4,
			Visible:     true,
		},
		{
			ID:          "combat_no_damage",
			Name:        "完美闪避",
			Description: "无伤通过一关",
			Category:    AchievementCategoryCombat,
			Condition:   AchievementCondition{Type: "no_damage_win", Threshold: 1},
			Rewards:     AchievementRewards{Coins: 500, Gems: 50, Badge: "perfect_dodge"},
			Icon:        "💫",
			Rarity:      RarityRare,
			SortOrder:   5,
			Visible:     true,
		},
	}

	// 收集成就
	collectionAchievements := []*Achievement{
		{
			ID:          "collection_first_item",
			Name:        "获得第一件道具",
			Description: "获得第一件道具",
			Category:    AchievementCategoryCollection,
			Condition:   AchievementCondition{Type: "item_count", Threshold: 1},
			Rewards:     AchievementRewards{Coins: 50, Gems: 5},
			Icon:        "📦",
			Rarity:      RarityCommon,
			SortOrder:   1,
			Visible:     true,
		},
		{
			ID:          "collection_100_items",
			Name:        "收藏家",
			Description: "累计获得100件道具",
			Category:    AchievementCategoryCollection,
			Condition:   AchievementCondition{Type: "item_count", Threshold: 100},
			Rewards:     AchievementRewards{Coins: 800, Gems: 80, Badge: "collector"},
			Icon:        "🗃️",
			Rarity:      RarityRare,
			SortOrder:   2,
			Visible:     true,
		},
		{
			ID:          "collection_rare_item",
			Name:        "稀有收获",
			Description: "获得第一件稀有及以上装备",
			Category:    AchievementCategoryCollection,
			Condition:   AchievementCondition{Type: "rare_item_count", Threshold: 1},
			Rewards:     AchievementRewards{Coins: 300, Gems: 30, Badge: "rare_finder"},
			Icon:        "💎",
			Rarity:      RarityUncommon,
			SortOrder:   3,
			Visible:     true,
		},
	}

	// 社交成就
	socialAchievements := []*Achievement{
		{
			ID:          "social_first_friend",
			Name:        "认识新朋友",
			Description: "添加第一个好友",
			Category:    AchievementCategorySocial,
			Condition:   AchievementCondition{Type: "friend_count", Threshold: 1},
			Rewards:     AchievementRewards{Coins: 50, Gems: 5},
			Icon:        "🤝",
			Rarity:      RarityCommon,
			SortOrder:   1,
			Visible:     true,
		},
		{
			ID:          "social_10_friends",
			Name:        "人脉广泛",
			Description: "拥有10个好友",
			Category:    AchievementCategorySocial,
			Condition:   AchievementCondition{Type: "friend_count", Threshold: 10},
			Rewards:     AchievementRewards{Coins: 300, Gems: 30, Badge: "popular"},
			Icon:        "👥",
			Rarity:      RarityUncommon,
			SortOrder:   2,
			Visible:     true,
		},
		{
			ID:          "social_guild_master",
			Name:        "公会之长",
			Description: "创建或加入公会",
			Category:    AchievementCategorySocial,
			Condition:   AchievementCondition{Type: "guild_joined", Threshold: 1},
			Rewards:     AchievementRewards{Coins: 500, Gems: 50, Title: "公会之长"},
			Icon:        "🏰",
			Rarity:      RarityRare,
			SortOrder:   3,
			Visible:     true,
		},
	}

	// 时间成就
	timeAchievements := []*Achievement{
		{
			ID:          "time_7_days",
			Name:        "坚持不懈",
			Description: "连续签到7天",
			Category:    AchievementCategoryTime,
			Condition:   AchievementCondition{Type: "consecutive_days", Threshold: 7},
			Rewards:     AchievementRewards{Coins: 300, Gems: 30, Badge: "week_streak"},
			Icon:        "📅",
			Rarity:      RarityUncommon,
			SortOrder:   1,
			Visible:     true,
		},
		{
			ID:          "time_30_days",
			Name:        "持之以恒",
			Description: "连续签到30天",
			Category:    AchievementCategoryTime,
			Condition:   AchievementCondition{Type: "consecutive_days", Threshold: 30},
			Rewards:     AchievementRewards{Coins: 1000, Gems: 100, Badge: "month_streak", Title: "坚持不懈"},
			Icon:        "🗓️",
			Rarity:      RarityEpic,
			SortOrder:   2,
			Visible:     true,
		},
		{
			ID:          "time_100_hours",
			Name:        "资深玩家",
			Description: "累计在线100小时",
			Category:    AchievementCategoryTime,
			Condition:   AchievementCondition{Type: "play_time", Threshold: 360000},
			Rewards:     AchievementRewards{Coins: 2000, Gems: 200, Badge: "veteran", Frame: "veteran_frame"},
			Icon:        "⏰",
			Rarity:      RarityEpic,
			SortOrder:   3,
			Visible:     true,
		},
	}

	for _, ach := range combatAchievements {
		am.RegisterAchievement(ach)
	}
	for _, ach := range collectionAchievements {
		am.RegisterAchievement(ach)
	}
	for _, ach := range socialAchievements {
		am.RegisterAchievement(ach)
	}
	for _, ach := range timeAchievements {
		am.RegisterAchievement(ach)
	}
}

// ==================== 玩家任务/成就扩展 ====================

// GetQuest 获取玩家任务
func (p *Player) GetQuest(questID string) *PlayerQuest {
	for _, q := range p.Quests {
		if q.QuestID == questID {
			return q
		}
	}
	return nil
}

// GetAchievement 获取玩家成就
func (p *Player) GetAchievement(achievementID string) *PlayerAchievement {
	for _, a := range p.Achievements {
		if a.AchievementID == achievementID {
			return a
		}
	}
	return nil
}

// AddTitle 添加称号
func (p *Player) AddTitle(title string) {
	for _, t := range p.Titles {
		if t == title {
			return
		}
	}
	p.Titles = append(p.Titles, title)
}

// AddBadge 添加徽章
func (p *Player) AddBadge(badge string) {
	for _, b := range p.Badges {
		if b == badge {
			return
		}
	}
	p.Badges = append(p.Badges, badge)
}

// AddAvatarFrame 添加头像框
func (p *Player) AddAvatarFrame(frame string) {
	for _, f := range p.AvatarFrames {
		if f == frame {
			return
		}
	}
	p.AvatarFrames = append(p.AvatarFrames, frame)
}
