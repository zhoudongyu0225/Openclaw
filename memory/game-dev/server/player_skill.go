package main

import (
	"math"
	"sync"
	"time"
)

// ============================================
// 玩家技能系统 (Player Skill System)
// 弹幕游戏核心玩法
// ============================================

// SkillCategory 技能分类
type SkillCategory int

const (
	SkillCategoryAttack SkillCategory = iota // 攻击技能
	SkillCategoryDefense                     // 防御技能
	SkillCategorySupport                     // 辅助技能
	SkillCategoryMovement                    // 移动技能
	SkillCategorySpecial                     // 特殊技能
)

// SkillTargetType 技能目标类型
type SkillTargetType int

const (
	SkillTargetNone SkillTargetType = iota
	SkillTargetSelf
	SkillTargetEnemy
	SkillTargetArea
	SkillTargetAlly
)

// SkillEffect 技能效果
type SkillEffect struct {
	Type           string  `json:"type"`            // 效果类型: damage, heal, shield, buff, debuff
	Value          float64 `json:"value"`           // 效果值
	Duration       int64   `json:"duration"`        // 持续时间(ms)
	Radius         float64 `json:"radius"`           // 范围半径
	Chance         float64 `json:"chance"`           // 触发概率
	Stacks         int     `json:"stacks"`           // 叠加层数
	Attribute      string  `json:"attribute"`       // 属性名
	PerStackBonus  float64 `json:"perStackBonus"`   // 每层加成
}

// Skill 技能定义
type Skill struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Category     SkillCategory `json:"category"`
	TargetType   SkillTargetType `json:"targetType"`
	Level        int           `json:"level"`           // 技能等级
	MaxLevel     int           `json:"maxLevel"`        // 最大等级
	Cooldown     int64         `json:"cooldown"`        // 冷却时间(ms)
	ManaCost     int           `json:"manaCost"`        // 魔法消耗
	CastTime     int64         `json:"castTime"`        // 施法时间(ms)
	Range        float64       `json:"range"`           // 施放范围
	Radius       float64       `json:"radius"`          // 效果范围
	Duration     int64         `json:"duration"`        // 持续时间
	Icon         string        `json:"icon"`            // 图标
	Effects      []SkillEffect `json:"effects"`        // 技能效果
	UpgradeCost  []int         `json:"upgradeCost"`     // 升级消耗
	UnlockLevel  int           `json:"unlockLevel"`     // 解锁等级
	RequiredSkills []string    `json:"requiredSkills"` // 前置技能
}

// PlayerSkillInstance 玩家技能实例
type PlayerSkillInstance struct {
	SkillID     string         `json:"skillId"`
	Level       int            `json:"level"`
	Exp         int            `json:"exp"`
	LastUsed    int64          `json:"lastUsed"`      // 上次使用时间
	CooldownEnd int64          `json:"cooldownEnd"`   // 冷却结束时间
	IsUnlocked  bool           `json:"isUnlocked"`    // 是否解锁
	IsActive    bool           `json:"isActive"`      // 是否激活(持续技能)
	ActiveEnd   int64          `json:"activeEnd"`     // 激活结束时间
}

// PlayerSkillTree 玩家技能树
type PlayerSkillTree struct {
	PlayerID    string                  `json:"playerId"`
	Skills      map[string]*PlayerSkillInstance `json:"skills"`
	AvailablePoints int                  `json:"availablePoints"` // 可用技能点
	TotalPoints int                     `json:"totalPoints"`     // 总技能点
	mu          sync.RWMutex
}

// SkillManager 技能管理器
type SkillManager struct {
	Skills     map[string]*Skill         // skillID -> Skill
	Trees      map[string]*PlayerSkillTree // playerID -> SkillTree
	Cooldowns  map[string]map[string]int64 // playerID -> {skillID -> cooldownEnd}
	mu         sync.RWMutex
}

// 创建技能管理器
func NewSkillManager() *SkillManager {
	return &SkillManager{
		Skills:    make(map[string]*Skill),
		Trees:     make(map[string]*PlayerSkillTree),
		Cooldowns: make(map[string]map[string]int64),
	}
}

// 注册技能
func (sm *SkillManager) RegisterSkill(skill *Skill) {
	sm.Skills[skill.ID] = skill
}

// 获取技能
func (sm *SkillManager) GetSkill(skillID string) *Skill {
	return sm.Skills[skillID]
}

// 获取玩家技能树
func (sm *SkillManager) GetOrCreateTree(playerID string) *PlayerSkillTree {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	if tree, ok := sm.Trees[playerID]; ok {
		return tree
	}
	
	tree := &PlayerSkillTree{
		PlayerID:       playerID,
		Skills:         make(map[string]*PlayerSkillInstance),
		AvailablePoints: 0,
		TotalPoints:    0,
	}
	
	// 初始化所有技能为未解锁
	for skillID := range sm.Skills {
		tree.Skills[skillID] = &PlayerSkillInstance{
			SkillID:    skillID,
			Level:      0,
			Exp:        0,
			IsUnlocked: false,
		}
	}
	
	sm.Trees[playerID] = tree
	return tree
}

// 解锁技能
func (sm *SkillManager) UnlockSkill(playerID, skillID string) error {
	tree := sm.GetOrCreateTree(playerID)
	tree.mu.Lock()
	defer tree.mu.Unlock()
	
	skill := sm.Skills[skillID]
	if skill == nil {
		return ErrSkillNotFound
	}
	
	instance := tree.Skills[skillID]
	if instance.IsUnlocked {
		return ErrSkillAlreadyUnlocked
	}
	
	// 检查前置技能
	for _, reqSkillID := range skill.RequiredSkills {
		if !tree.Skills[reqSkillID].IsUnlocked {
			return ErrSkillPrerequisitesNotMet
		}
	}
	
	// 检查技能点
	if tree.AvailablePoints < 1 {
		return ErrNotEnoughSkillPoints
	}
	
	instance.IsUnlocked = true
	instance.Level = 1
	tree.AvailablePoints--
	tree.TotalPoints++
	
	return nil
}

// 升级技能
func (sm *SkillManager) UpgradeSkill(playerID, skillID string) error {
	tree := sm.GetOrCreateTree(playerID)
	tree.mu.Lock()
	defer tree.mu.Unlock()
	
	skill := sm.Skills[skillID]
	if skill == nil {
		return ErrSkillNotFound
	}
	
	instance := tree.Skills[skillID]
	if !instance.IsUnlocked {
		return ErrSkillNotUnlocked
	}
	
	if instance.Level >= skill.MaxLevel {
		return ErrSkillMaxLevel
	}
	
	// 检查技能点
	if tree.AvailablePoints < 1 {
		return ErrNotEnoughSkillPoints
	}
	
	// 检查升级消耗
	cost := skill.UpgradeCost[instance.Level-1]
	
	instance.Level++
	tree.AvailablePoints--
	tree.TotalPoints++
	
	return nil
}

// 使用技能
func (sm *SkillManager) UseSkill(playerID, skillID string, x, y float64, targetID string) (*SkillResult, error) {
	skill := sm.Skills[skillID]
	if skill == nil {
		return nil, ErrSkillNotFound
	}
	
	tree := sm.GetOrCreateTree(playerID)
	tree.mu.Lock()
	instance := tree.Skills[skillID]
	tree.mu.Unlock()
	
	if !instance.IsUnlocked {
		return nil, ErrSkillNotUnlocked
	}
	
	now := time.Now().UnixMilli()
	
	// 检查冷却
	if now < instance.CooldownEnd {
		return nil, ErrSkillOnCooldown
	}
	
	// 启动冷却
	instance.CooldownEnd = now + skill.Cooldown
	instance.LastUsed = now
	
	// 计算技能效果
	result := &SkillResult{
		SkillID:  skillID,
		CasterID: playerID,
		X:        x,
		Y:        y,
		TargetID: targetID,
		Effects:  make([]*SkillEffectResult, 0),
		Success:  true,
	}
	
	// 应用效果
	for _, effect := range skill.Effects {
		effectResult := sm.ApplyEffect(effect, skill, instance.Level, x, y, targetID)
		result.Effects = append(result.Effects, effectResult)
	}
	
	// 标记持续技能
	if skill.Duration > 0 {
		instance.IsActive = true
		instance.ActiveEnd = now + skill.Duration
	}
	
	return result, nil
}

// SkillResult 技能使用结果
type SkillResult struct {
	SkillID     string              `json:"skillId"`
	CasterID    string              `json:"casterId"`
	X           float64             `json:"x"`
	Y           float64              `json:"y"`
	TargetID    string              `json:"targetId"`
	Effects     []*SkillEffectResult `json:"effects"`
	Success     bool                `json:"success"`
	FailReason  string              `json:"failReason,omitempty"`
}

// SkillEffectResult 技能效果结果
type SkillEffectResult struct {
	Type       string  `json:"type"`
	Value      float64 `json:"value"`
	Affected   []string `json:"affected"` // 受影响的目标ID
	Duration   int64   `json:"duration"`
	TotalDamage float64 `json:"totalDamage"`
	TotalHeal  float64 `json:"totalHeal"`
}

// 应用技能效果
func (sm *SkillManager) ApplyEffect(effect *SkillEffect, skill *Skill, level int, x, y float64, targetID string) *SkillEffectResult {
	result := &SkillEffectResult{
		Type:     effect.Type,
		Duration: effect.Duration,
	}
	
	// 计算等级加成
	levelBonus := 1.0 + float64(level-1)*0.1
	value := effect.Value * levelBonus
	
	switch effect.Type {
	case "damage":
		result.TotalDamage = value
		result.Affected = []string{targetID}
		
	case "heal":
		result.TotalHeal = value
		result.Affected = []string{""}
		
	case "shield":
		result.Value = value
		result.Affected = []string{""}
		
	case "buff", "debuff":
		result.Value = value
		result.Affected = sm.FindTargetsInArea(x, y, effect.Radius, effect.Type == "debuff")
		
	case "area_damage":
		result.TotalDamage = value
		result.Affected = sm.FindTargetsInArea(x, y, effect.Radius, true)
	}
	
	return result
}

// 在区域内查找目标
func (sm *SkillManager) FindTargetsInArea(x, y, radius float64, enemyOnly bool) []string {
	// 简化实现 - 实际应该查询游戏实体
	return []string{}
}

// 技能错误
var (
	ErrSkillNotFound             = &SkillError{"skill_not_found", "技能不存在"}
	ErrSkillAlreadyUnlocked      = &SkillError{"skill_already_unlocked", "技能已解锁"}
	ErrSkillNotUnlocked          = &SkillError{"skill_not_unlocked", "技能未解锁"}
	ErrSkillOnCooldown           = &SkillError{"skill_on_cooldown", "技能冷却中"}
	ErrSkillMaxLevel             = &SkillError{"skill_max_level", "技能已达最大等级"}
	ErrNotEnoughSkillPoints      = &SkillError{"not_enough_skill_points", "技能点不足"}
	ErrSkillPrerequisitesNotMet  = &SkillError{"skill_prerequisites_not_met", "前置技能未满足"}
)

type SkillError struct {
	Code string
	Msg  string
}

func (e *SkillError) Error() string return e.Msg

// ============================================
// 预设技能
// ============================================

// 注册所有预设技能
func (sm *SkillManager) RegisterDefaultSkills() {
	// 攻击技能
	sm.RegisterSkill(&Skill{
		ID: "slash", Name: "利刃", Description: "快速挥砍敌人",
		Category: SkillCategoryAttack, TargetType: SkillTargetEnemy,
		Cooldown: 500, ManaCost: 10, CastTime: 200, Range: 100,
		MaxLevel: 5, UnlockLevel: 1,
		Effects: []SkillEffect{
			{Type: "damage", Value: 50, Radius: 50},
		},
		UpgradeCost: []int{0, 100, 200, 400, 800},
	})
	
	sm.RegisterSkill(&Skill{
		ID: "fireball", Name: "火球", Description: "发射火球造成范围伤害",
		Category: SkillCategoryAttack, TargetType: SkillTargetArea,
		Cooldown: 2000, ManaCost: 25, CastTime: 500, Range: 300, Radius: 80,
		MaxLevel: 5, UnlockLevel: 3,
		Effects: []SkillEffect{
			{Type: "area_damage", Value: 80, Radius: 80},
		},
		UpgradeCost: []int{0, 150, 300, 600, 1200},
		RequiredSkills: []string{"slash"},
	})
	
	sm.RegisterSkill(&Skill{
		ID: "lightning", Name: "雷电", Description: "连锁闪电攻击多个敌人",
		Category: SkillCategoryAttack, TargetType: SkillTargetEnemy,
		Cooldown: 3000, ManaCost: 40, CastTime: 300, Range: 250,
		MaxLevel: 3, UnlockLevel: 5,
		Effects: []SkillEffect{
			{Type: "damage", Value: 100},
			{Type: "debuff", Value: 0.3, Duration: 2000, Attribute: "speed"},
		},
		UpgradeCost: []int{0, 500, 1000},
		RequiredSkills: []string{"fireball"},
	})
	
	// 防御技能
	sm.RegisterSkill(&Skill{
		ID: "shield", Name: "护盾", Description: "为自己施加护盾",
		Category: SkillCategoryDefense, TargetType: SkillTargetSelf,
		Cooldown: 8000, ManaCost: 20, CastTime: 100,
		MaxLevel: 5, UnlockLevel: 1,
		Effects: []SkillEffect{
			{Type: "shield", Value: 100, Duration: 5000},
		},
		UpgradeCost: []int{0, 80, 160, 320, 640},
	})
	
	sm.RegisterSkill(&Skill{
		ID: "dodge", Name: "闪避", Description: "瞬间闪避伤害",
		Category: SkillCategoryMovement, TargetType: SkillTargetSelf,
		Cooldown: 5000, ManaCost: 15, CastTime: 50,
		MaxLevel: 5, UnlockLevel: 1,
		Effects: []SkillEffect{
			{Type: "buff", Value: 2.0, Duration: 500, Attribute: "dodge"},
		},
		UpgradeCost: []int{0, 60, 120, 240, 480},
	})
	
	// 辅助技能
	sm.RegisterSkill(&Skill{
		ID: "heal", Name: "治疗", Description: "恢复自身生命",
		Category: SkillCategorySupport, TargetType: SkillTargetSelf,
		Cooldown: 10000, ManaCost: 30, CastTime: 500,
		MaxLevel: 5, UnlockLevel: 2,
		Effects: []SkillEffect{
			{Type: "heal", Value: 80},
		},
		UpgradeCost: []int{0, 120, 240, 480, 960},
	})
	
	// 特殊技能
	sm.RegisterSkill(&Skill{
		ID: "time_stop", Name: "时停", Description: "暂停时间",
		Category: SkillCategorySpecial, TargetType: SkillTargetArea,
		Cooldown: 60000, ManaCost: 100, CastTime: 1000, Range: 200,
		MaxLevel: 3, UnlockLevel: 8,
		Effects: []SkillEffect{
			{Type: "debuff", Value: 0, Duration: 3000, Attribute: "frozen"},
		},
		UpgradeCost: []int{0, 2000, 4000},
		RequiredSkills: []string{"lightning", "heal"},
	})
	
	sm.RegisterSkill(&Skill{
		ID: "meteor", Name: "流星", Description: "召唤流星雨",
		Category: SkillCategorySpecial, TargetType: SkillTargetArea,
		Cooldown: 90000, ManaCost: 150, CastTime: 2000, Range: 400, Radius: 150,
		MaxLevel: 3, UnlockLevel: 10,
		Effects: []SkillEffect{
			{Type: "area_damage", Value: 200, Radius: 150},
		},
		UpgradeCost: []int{0, 3000, 6000},
		RequiredSkills: []string{"time_stop"},
	})
}

// ============================================
// 技能点系统
// ============================================

// SkillPointRule 技能点规则
type SkillPointRule struct {
	LevelUpPoints    int // 升级获得技能点
	AchievementPoints int // 成就奖励技能点
	EventPoints      int // 活动奖励技能点
}

// 获取技能点
func (sm *SkillManager) AwardSkillPoints(playerID string, points int, reason string) {
	tree := sm.GetOrCreateTree(playerID)
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.AvailablePoints += points
}

// ============================================
// 技能状态效果 (Buff/Debuff System)
// ============================================

// Buff 状态效果
type Buff struct {
	ID          string    `json:"id"`
	SkillID     string    `json:"skillId"`
	CasterID    string    `json:"casterId"`
	TargetID    string    `json:"targetId"`
	Type        string    `json:"type"` // buff 或 debuff
	Attribute   string    `json:"attribute"`
	Value       float64   `json:"value"`
	Stacks      int       `json:"stacks"`
	MaxStacks   int       `json:"maxStacks"`
	Duration    int64     `json:"duration"`
	StartTime   int64     `json:"startTime"`
	EndTime     int64     `json:"endTime"`
}

// BuffManager Buff管理器
type BuffManager struct {
	Buffs    map[string]map[string]*Buff // targetID -> {buffID -> Buff}
	mu       sync.RWMutex
}

func NewBuffManager() *BuffManager {
	return &BuffManager{
		Buffs: make(map[string]map[string]*Buff),
	}
}

// 添加Buff
func (bm *BuffManager) AddBuff(buff *Buff) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	
	if bm.Buffs[buff.TargetID] == nil {
		bm.Buffs[buff.TargetID] = make(map[string]*Buff)
	}
	
	// 检查是否可叠加
	if existing, ok := bm.Buffs[buff.TargetID][buff.ID]; ok {
		if existing.Stacks < existing.MaxStacks {
			existing.Stacks++
			existing.Value += buff.Value
		}
		existing.EndTime = buff.EndTime
	} else {
		bm.Buffs[buff.TargetID][buff.ID] = buff
	}
}

// 移除Buff
func (bm *BuffManager) RemoveBuff(targetID, buffID string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	delete(bm.Buffs[targetID], buffID)
}

// 获取目标的所有Buff
func (bm *BuffManager) GetBuffs(targetID string) []*Buff {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	
	var result []*Buff
	for _, buff := range bm.Buffs[targetID] {
		result = append(result, buff)
	}
	return result
}

// 更新Buff状态
func (bm *BuffManager) Update() {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	
	now := time.Now().UnixMilli()
	
	for targetID, buffs := range bm.Buffs {
		for buffID, buff := range buffs {
			if now > buff.EndTime {
				delete(buffs, buffID)
			}
		}
		
		// 清理空目标
		if len(buffs) == 0 {
			delete(bm.Buffs, targetID)
		}
	}
}

// 获取属性加成
func (bm *BuffManager) GetAttributeBonus(targetID, attribute string) float64 {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	
	var total float64
	for _, buff := range bm.Buffs[targetID] {
		if buff.Attribute == attribute && buff.Type == "buff" {
			total += buff.Value * float64(buff.Stacks)
		}
	}
	return total
}
