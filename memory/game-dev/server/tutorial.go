package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ==================== 新手引导系统 ====================

// TutorialStep 新手引导步骤
type TutorialStep struct {
	ID          string   `json:"id"`           // 步骤ID
	Title       string   `json:"title"`        // 步骤标题
	Description string   `json:"description"`  // 步骤描述
	Trigger     string   `json:"trigger"`       // 触发条件类型
	TriggerData string   `json:"trigger_data"`  // 触发条件数据
	PositionX   float64  `json:"position_x"`   // 指引位置X
	PositionY   float64  `json:"position_y"`   // 指引位置Y
	Duration    int      `json:"duration"`      // 持续时间(毫秒)
	AutoNext    bool     `json:"auto_next"`    // 是否自动下一步
	Rewards     []Reward `json:"rewards"`      // 完成奖励
}

// TutorialCategory 引导分类
type TutorialCategory int

const (
	TutorialCategoryBasic    TutorialCategory = iota // 基础引导
	TutorialCategoryCombat                            // 战斗引导
	TutorialCategorySocial                            // 社交引导
	TutorialCategoryShop                              // 商城引导
	TutorialCategoryAdvanced                          // 高级引导
)

// Tutorial 新手引导
type Tutorial struct {
	ID         string           `json:"id"`          // 引导ID
	Category   TutorialCategory `json:"category"`    // 引导分类
	Name       string           `json:"name"`        // 引导名称
	Steps      []TutorialStep   `json:"steps"`       // 引导步骤
	MinLevel   int              `json:"min_level"`   // 最低等级
	MaxLevel   int              `json:"max_level"`   // 最高等级
	Required   bool             `json:"required"`     // 是否必须完成
	Skippable  bool             `json:"skippable"`   // 是否可跳过
	ResetOnFail bool            `json:"reset_on_fail"` // 失败后重置
}

// PlayerTutorial 玩家引导进度
type PlayerTutorial struct {
	PlayerID     string    `json:"player_id"`      // 玩家ID
	TutorialID   string    `json:"tutorial_id"`     // 引导ID
	CurrentStep  int       `json:"current_step"`    // 当前步骤
	Completed    bool      `json:"completed"`       // 是否完成
	Skipped      bool      `json:"skipped"`         // 是否跳过
	StartedAt    time.Time `json:"started_at"`      // 开始时间
	CompletedAt  time.Time `json:"completed_at"`    // 完成时间
	StepHistory  []int     `json:"step_history"`    // 步骤历史
}

// TutorialManager 新手引导管理器
type TutorialManager struct {
	tutorials    map[string]*Tutorial           // 引导列表
	playerProgress map[string]map[string]*PlayerTutorial // 玩家进度
	mu            sync.RWMutex
	once          sync.Once
}

// NewTutorialManager 创建引导管理器
func NewTutorialManager() *TutorialManager {
	m := &TutorialManager{
		tutorials:     make(map[string]*Tutorial),
		playerProgress: make(map[string]map[string]*PlayerTutorial),
	}
	m.initTutorials()
	return m
}

// 初始化预设引导
func (m *TutorialManager) initTutorials() {
	// 1. 基础操作引导
	m.tutorials["basic_welcome"] = &Tutorial{
		ID:        "basic_welcome",
		Category:  TutorialCategoryBasic,
		Name:      "欢迎来到游戏",
		MinLevel:  1,
		MaxLevel:  5,
		Required:  true,
		Skippable: true,
		Steps: []TutorialStep{
			{
				ID:          "welcome_1",
				Title:       "欢迎",
				Description: "欢迎来到弹幕游戏世界！点击继续开始你的冒险之旅。",
				Trigger:     "auto",
				PositionX:   0.5,
				PositionY:   0.5,
				Duration:    0,
				AutoNext:    false,
				Rewards:     []Reward{{Type: "coins", Amount: 100}},
			},
			{
				ID:          "welcome_2",
				Title:       "角色移动",
				Description: "使用屏幕左下角的虚拟摇杆控制角色移动。",
				Trigger:     "touch_area",
				TriggerData: "joystick",
				PositionX:   0.15,
				PositionY:   0.25,
				Duration:    5000,
				AutoNext:    false,
			},
			{
				ID:          "welcome_3",
				Title:       "战斗操作",
				Description: "点击屏幕右侧的攻击按钮释放技能击败敌人。",
				Trigger:     "touch_area",
				TriggerData: "attack_button",
				PositionX:   0.85,
				PositionY:   0.25,
				Duration:    5000,
				AutoNext:    false,
			},
		},
	}

	// 2. 战斗入门引导
	m.tutorials["combat_basics"] = &Tutorial{
		ID:         "combat_basics",
		Category:   TutorialCategoryCombat,
		Name:       "战斗入门",
		MinLevel:   1,
		MaxLevel:   10,
		Required:   true,
		Skippable:  false,
		Steps: []TutorialStep{
			{
				ID:          "combat_1",
				Title:       "了解敌人",
				Description: "敌人会从上方发射弹幕，躲避弹幕并击败它们！",
				Trigger:     "enter_room",
				TriggerData: "combat_1",
				PositionX:   0.5,
				PositionY:   0.3,
				Duration:    3000,
				AutoNext:    true,
			},
			{
				ID:          "combat_2",
				Title:       "释放技能",
				Description: "点击技能按钮释放技能，造成大量伤害！",
				Trigger:     "skill_available",
				TriggerData: "skill_1",
				PositionX:   0.75,
				PositionY:   0.3,
				Duration:    5000,
				AutoNext:    false,
				Rewards:     []Reward{{Type: "gems", Amount: 10}},
			},
			{
				ID:          "combat_3",
				Title:       "闪避技巧",
				Description: "危急时刻使用闪避技能躲避致命弹幕！",
				Trigger:     "hp_low",
				TriggerData: "30",
				PositionX:   0.5,
				PositionY:   0.5,
				Duration:    5000,
				AutoNext:    false,
			},
		},
	}

	// 3. 社交引导
	m.tutorials["social_intro"] = &Tutorial{
		ID:         "social_intro",
		Category:   TutorialCategorySocial,
		Name:       "社交系统",
		MinLevel:   3,
		MaxLevel:   15,
		Required:   false,
		Skippable:  true,
		Steps: []TutorialStep{
			{
				ID:          "social_1",
				Title:       "添加好友",
				Description: "点击好友按钮，你可以添加其他玩家为好友，一起组队游戏！",
				Trigger:     "menu_open",
				TriggerData: "friend_menu",
				PositionX:   0.9,
				PositionY:   0.15,
				Duration:    5000,
				AutoNext:    false,
			},
			{
				ID:          "social_2",
				Title:       "加入公会",
				Description: "加入公会可以与其他玩家一起交流，参与公会活动！",
				Trigger:     "menu_open",
				TriggerData: "guild_menu",
				PositionX:   0.9,
				PositionY:   0.2,
				Duration:    5000,
				AutoNext:    false,
				Rewards:     []Reward{{Type: "coins", Amount: 200}},
			},
		},
	}

	// 4. 商城引导
	m.tutorials["shop_tutorial"] = &Tutorial{
		ID:         "shop_tutorial",
		Category:   TutorialCategoryShop,
		Name:       "商城系统",
		MinLevel:   5,
		MaxLevel:   20,
		Required:   false,
		Skippable:  true,
		Steps: []TutorialStep{
			{
				ID:          "shop_1",
				Title:       "商城入口",
				Description: "在这里你可以购买各种道具、角色和皮肤！",
				Trigger:     "menu_open",
				TriggerData: "shop_menu",
				PositionX:   0.85,
				PositionY:   0.1,
				Duration:    3000,
				AutoNext:    true,
			},
			{
				ID:          "shop_2",
				Title:       "首充礼包",
				Description: "首充任意金额即可获得超值大礼包！",
				Trigger:     "first_recharge",
				TriggerData: "",
				PositionX:   0.5,
				PositionY:   0.4,
				Duration:    0,
				AutoNext:    false,
				Rewards:     []Reward{{Type: "item", Amount: 1, ItemID: "first_recharge_pack"}},
			},
		},
	}

	// 5. 高级战斗引导
	m.tutorials["advanced_combat"] = &Tutorial{
		ID:         "advanced_combat",
		Category:   TutorialCategoryAdvanced,
		Name:       "高级战斗",
		MinLevel:   10,
		MaxLevel:   30,
		Required:   false,
		Skippable:  true,
		Steps: []TutorialStep{
			{
				ID:          "adv_1",
				Title:       "连击系统",
				Description: "连续击中敌人可以触发连击，获得额外分数和奖励！",
				Trigger:     "combo_start",
				TriggerData: "5",
				PositionX:   0.5,
				PositionY:   0.1,
				Duration:    5000,
				AutoNext:    true,
			},
			{
				ID:          "adv_2",
				Title:       "使用必杀技",
				Description: "积累能量后释放必杀技，瞬间清屏！",
				Trigger:     "ult_ready",
				TriggerData: "",
				PositionX:   0.8,
				PositionY:   0.2,
				Duration:    5000,
				AutoNext:    false,
				Rewards:     []Reward{{Type: "gems", Amount: 50}},
			},
		},
	}
}

// GetTutorial 获取引导
func (m *TutorialManager) GetTutorial(id string) *Tutorial {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tutorials[id]
}

// GetTutorialsByCategory 按分类获取引导
func (m *TutorialManager) GetTutorialsByCategory(category TutorialCategory) []*Tutorial {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Tutorial
	for _, t := range m.tutorials {
		if t.Category == category {
			result = append(result, t)
		}
	}
	return result
}

// GetAvailableTutorials 获取玩家可用引导
func (m *TutorialManager) GetAvailableTutorials(playerID string, level int) []*Tutorial {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Tutorial
	progress := m.playerProgress[playerID]

	for _, t := range m.tutorials {
		// 检查等级范围
		if level < t.MinLevel || level > t.MaxLevel {
			continue
		}

		// 检查是否已完成
		if progress != nil {
			if p, ok := progress[t.ID]; ok && p.Completed {
				continue
			}
		}

		result = append(result, t)
	}
	return result
}

// StartTutorial 开始引导
func (m *TutorialManager) StartTutorial(playerID, tutorialID string) (*PlayerTutorial, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tutorial := m.tutorials[tutorialID]
	if tutorial == nil {
		return nil, fmt.Errorf("tutorial not found: %s", tutorialID)
	}

	// 初始化玩家进度
	if m.playerProgress[playerID] == nil {
		m.playerProgress[playerID] = make(map[string]*PlayerTutorial)
	}

	// 检查是否已有进度
	if existing, ok := m.playerProgress[playerID][tutorialID]; ok {
		if existing.Completed {
			return existing, fmt.Errorf("tutorial already completed")
		}
		if existing.Skipped && !tutorial.Required {
			// 重新开始
			existing.CurrentStep = 0
			existing.Skipped = false
			existing.StartedAt = time.Now()
			return existing, nil
		}
		return existing, nil
	}

	// 创建新进度
	progress := &PlayerTutorial{
		PlayerID:    playerID,
		TutorialID:  tutorialID,
		CurrentStep: 0,
		Completed:   false,
		Skipped:     false,
		StartedAt:   time.Now(),
		StepHistory: []int{},
	}

	m.playerProgress[playerID][tutorialID] = progress
	return progress, nil
}

// CompleteStep 完成步骤
func (m *TutorialManager) CompleteStep(playerID, tutorialID string) (*PlayerTutorial, []Reward, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	progress := m.playerProgress[playerID][tutorialID]
	if progress == nil {
		return nil, nil, fmt.Errorf("tutorial not started")
	}

	tutorial := m.tutorials[tutorialID]
	if tutorial == nil {
		return nil, nil, fmt.Errorf("tutorial not found")
	}

	// 记录步骤历史
	progress.StepHistory = append(progress.StepHistory, progress.CurrentStep)

	// 检查是否完成所有步骤
	if progress.CurrentStep >= len(tutorial.Steps)-1 {
		progress.Completed = true
		progress.CompletedAt = time.Now()

		// 发放完成奖励
		var rewards []Reward
		for _, step := range tutorial.Steps {
			rewards = append(rewards, step.Rewards...)
		}
		return progress, rewards, nil
	}

	// 进入下一步
	progress.CurrentStep++

	// 检查是否有步骤奖励
	var stepRewards []Reward
	currentStep := tutorial.Steps[progress.CurrentStep]
	if len(currentStep.Rewards) > 0 {
		stepRewards = currentStep.Rewards
	}

	return progress, stepRewards, nil
}

// SkipTutorial 跳过引导
func (m *TutorialManager) SkipTutorial(playerID, tutorialID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tutorial := m.tutorials[tutorialID]
	if tutorial == nil {
		return fmt.Errorf("tutorial not found")
	}

	if !tutorial.Skippable {
		return fmt.Errorf("tutorial cannot be skipped")
	}

	progress := m.playerProgress[playerID][tutorialID]
	if progress == nil {
		return fmt.Errorf("tutorial not started")
	}

	progress.Skipped = true
	progress.CompletedAt = time.Now()

	return nil
}

// GetProgress 获取玩家引导进度
func (m *TutorialManager) GetProgress(playerID string) []*PlayerTutorial {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PlayerTutorial
	if progress, ok := m.playerProgress[playerID]; ok {
		for _, p := range progress {
			result = append(result, p)
		}
	}
	return result
}

// GetCurrentStep 获取当前步骤详情
func (m *TutorialManager) GetCurrentStep(playerID, tutorialID string) (*TutorialStep, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	progress := m.playerProgress[playerID][tutorialID]
	if progress == nil {
		return nil, fmt.Errorf("tutorial not started")
	}

	if progress.Completed || progress.Skipped {
		return nil, fmt.Errorf("tutorial already finished")
	}

	tutorial := m.tutorials[tutorialID]
	if tutorial == nil {
		return nil, fmt.Errorf("tutorial not found")
	}

	if progress.CurrentStep >= len(tutorial.Steps) {
		return nil, fmt.Errorf("invalid step index")
	}

	return &tutorial.Steps[progress.CurrentStep], nil
}

// ResetTutorial 重置引导
func (m *TutorialManager) ResetTutorial(playerID, tutorialID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tutorial := m.tutorials[tutorialID]
	if tutorial == nil {
		return fmt.Errorf("tutorial not found")
	}

	if progress, ok := m.playerProgress[playerID][tutorialID]; ok {
		progress.CurrentStep = 0
		progress.Completed = false
		progress.Skipped = false
		progress.StepHistory = []int{}
		progress.StartedAt = time.Now()
		progress.CompletedAt = time.Time{}
	}

	return nil
}

// Serialize 序列化引导数据
func (m *TutorialManager) Serialize() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.tutorials)
}

// Deserialize 反序列化引导数据
func (m *TutorialManager) Deserialize(data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return json.Unmarshal(data, &m.tutorials)
}

// GetTutorialStats 获取引导统计
func (m *TutorialManager) GetTutorialStats(playerID string) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]interface{})
	progress := m.playerProgress[playerID]

	total := len(m.tutorials)
	completed := 0
	inProgress := 0
	skipped := 0

	for _, t := range m.tutorials {
		if p, ok := progress[t.ID]; ok {
			if p.Completed {
				completed++
			} else if p.Skipped {
				skipped++
			} else {
				inProgress++
			}
		}
	}

	stats["total_tutorials"] = total
	stats["completed"] = completed
	stats["in_progress"] = inProgress
	stats["skipped"] = skipped
	stats["not_started"] = total - completed - inProgress - skipped

	return stats
}

// TutorialTriggerHandler 触发条件处理器
type TutorialTriggerHandler struct {
	manager *TutorialManager
}

func NewTutorialTriggerHandler(m *TutorialManager) *TutorialTriggerHandler {
	return &TutorialTriggerHandler{manager: m}
}

// HandleTrigger 处理触发事件
func (h *TutorialTriggerHandler) HandleTrigger(playerID, triggerType, triggerData string) {
	tutorials := h.manager.GetAvailableTutorials(playerID, 1) // 简化，实际应传入玩家等级

	for _, t := range tutorials {
		for i, step := range t.Steps {
			if step.Trigger == triggerType {
				// 检查触发数据是否匹配
				if step.TriggerData == "" || step.TriggerData == triggerData {
					progress, _ := h.manager.playerProgress[playerID][t.ID]
					if progress != nil && progress.CurrentStep == i && !progress.Completed {
						// 触发下一步
						fmt.Printf("Trigger tutorial %s step %d for player %s\n", t.ID, i, playerID)
					}
				}
			}
		}
	}
}
