// Package game provides core game logic for the bullet hell game
package game

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ============================================================
// 装备强化系统 (Equipment Enhancement System)
// ============================================================

// EnhancementType 强化类型
type EnhancementType int

const (
	EnhancementTypeAttack EnhancementType = iota // 攻击强化
	EnhancementTypeDefense                       // 防御强化
	EnhancementTypeHP                            // 生命强化
	EnhancementTypeCrit                          // 暴击强化
	EnhancementTypeDodge                         // 闪避强化
	EnhancementTypeSpeed                         // 速度强化
)

// Enhancement 强化定义
type Enhancement struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Type        EnhancementType `json:"type"`
	Level       int             `json:"level"`
	BaseValue   float64         `json:"base_value"`
	IncValue    float64         `json:"inc_value"`
	SuccessRate float64         `json:"success_rate"`
	Cost        int             `json:"cost"`
	MaxLevel    int             `json:"max_level"`
	Description string          `json:"description"`
}

// EnhancementResult 强化结果
type EnhancementResult struct {
	Success bool              `json:"success"`
	Level   int               `json:"level"`
	Item    *Item            `json:"item,omitempty"`
	Stats   map[string]float64 `json:"stats"`
	Message string            `json:"message"`
}

// EquipmentEnhancementManager 装备强化管理器
type EquipmentEnhancementManager struct {
	enhancements map[string]*Enhancement
	enchantments map[string]*Enchantment
}

// NewEquipmentEnhancementManager 创建强化管理器
func NewEquipmentEnhancementManager() *EquipmentEnhancementManager {
	m := &EquipmentEnhancementManager{
		enhancements: make(map[string]*Enhancement),
		enchantments: make(map[string]*Enchantment),
	}
	m.initDefaultEnhancements()
	m.initDefaultEnchantments()
	return m
}

func (m *EquipmentEnhancementManager) initDefaultEnhancements() {
	enhancements := []*Enhancement{
		{ID: "atk_1", Name: "锋利", Type: EnhancementTypeAttack, Level: 1, BaseValue: 5, IncValue: 3, SuccessRate: 0.9, Cost: 100, MaxLevel: 10, Description: "增加武器攻击力"},
		{ID: "atk_2", Name: "锐利", Type: EnhancementTypeAttack, Level: 2, BaseValue: 15, IncValue: 5, SuccessRate: 0.8, Cost: 300, MaxLevel: 10, Description: "大幅增加武器攻击力"},
		{ID: "def_1", Name: "坚固", Type: EnhancementTypeDefense, Level: 1, BaseValue: 5, IncValue: 3, SuccessRate: 0.9, Cost: 100, MaxLevel: 10, Description: "增加防御力"},
		{ID: "def_2", Name: "铁壁", Type: EnhancementTypeDefense, Level: 2, BaseValue: 15, IncValue: 5, SuccessRate: 0.8, Cost: 300, MaxLevel: 10, Description: "大幅增加防御力"},
		{ID: "hp_1", Name: "生命", Type: EnhancementTypeHP, Level: 1, BaseValue: 50, IncValue: 30, SuccessRate: 0.9, Cost: 100, MaxLevel: 10, Description: "增加生命上限"},
		{ID: "hp_2", Name: "活力", Type: EnhancementTypeHP, Level: 2, BaseValue: 150, IncValue: 50, SuccessRate: 0.8, Cost: 300, MaxLevel: 10, Description: "大幅增加生命上限"},
		{ID: "crit_1", Name: "精准", Type: EnhancementTypeCrit, Level: 1, BaseValue: 3, IncValue: 2, SuccessRate: 0.85, Cost: 150, MaxLevel: 10, Description: "增加暴击率"},
		{ID: "dodge_1", Name: "灵敏", Type: EnhancementTypeDodge, Level: 1, BaseValue: 3, IncValue: 2, SuccessRate: 0.85, Cost: 150, MaxLevel: 10, Description: "增加闪避率"},
		{ID: "speed_1", Name: "神速", Type: EnhancementTypeSpeed, Level: 1, BaseValue: 2, IncValue: 1, SuccessRate: 0.8, Cost: 200, MaxLevel: 10, Description: "增加移动速度"},
	}
	for _, e := range enhancements {
		m.enhancements[e.ID] = e
	}
}

// Enchantment 附魔定义
type Enchantment struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Effect      string  `json:"effect"`
	Value       float64 `json:"value"`
	Rarity      int     `json:"rarity"`
	SuccessRate float64 `json:"success_rate"`
	Cost        int     `json:"cost"`
	Description string  `json:"description"`
}

func (m *EquipmentEnhancementManager) initDefaultEnchantments() {
	enchantments := []*Enchantment{
		{ID: "ench_fire", Name: "火焰附魔", Effect: "fire_damage", Value: 10, Rarity: 2, SuccessRate: 0.7, Cost: 500, Description: "攻击时附加火焰伤害"},
		{ID: "ench_ice", Name: "冰霜附魔", Effect: "ice_damage", Value: 10, Rarity: 2, SuccessRate: 0.7, Cost: 500, Description: "攻击时附加冰霜伤害"},
		{ID: "ench_thunder", Name: "雷电附魔", Effect: "thunder_damage", Value: 12, Rarity: 3, SuccessRate: 0.5, Cost: 800, Description: "攻击时附加雷电伤害"},
		{ID: "ench_poison", Name: "毒液附魔", Effect: "poison_damage", Value: 5, Rarity: 2, SuccessRate: 0.7, Cost: 400, Description: "攻击时附加毒液伤害"},
		{ID: "ench_blessing", Name: "祝福附魔", Effect: "blessing", Value: 15, Rarity: 3, SuccessRate: 0.5, Cost: 1000, Description: "全属性提升"},
		{ID: "ench_light", Name: "光明附魔", Effect: "holy_damage", Value: 15, Rarity: 4, SuccessRate: 0.3, Cost: 2000, Description: "对暗影生物额外伤害"},
		{ID: "ench_shadow", Name: "暗影附魔", Effect: "dark_damage", Value: 15, Rarity: 4, SuccessRate: 0.3, Cost: 2000, Description: "对光明生物额外伤害"},
	}
	for _, e := range enchantments {
		m.enchantments[e.ID] = e
	}
}

// Enhance 装备强化
func (m *EquipmentEnhancementManager) Enhance(player *Player, item *Item, enhancementID string) (*EnhancementResult, error) {
	if player == nil || item == nil {
		return nil, errors.New("玩家或物品不存在")
	}

	enhancement, ok := m.enhancements[enhancementID]
	if !ok {
		return nil, errors.New("强化类型不存在")
	}

	if item.EnhanceLevel >= enhancement.MaxLevel {
		return &EnhancementResult{
			Success: false,
			Level:   item.EnhanceLevel,
			Message: "已达到强化等级上限",
		}, nil
	}

	if player.Gold < enhancement.Cost {
		return nil, errors.New("金币不足")
	}

	successRate := m.calculateSuccessRate(enhancement, item.EnhanceLevel)

	player.Gold -= enhancement.Cost
	success := rand.Float64() < successRate

	result := &EnhancementResult{
		Success: success,
		Level:   item.EnhanceLevel,
		Item:    item,
		Stats:   make(map[string]float64),
	}

	if success {
		item.EnhanceLevel++
		item.EnhancementType = enhancement.Type
		item.EnhancementValue = enhancement.BaseValue + float64(item.EnhanceLevel-1)*enhancement.IncValue
		result.Level = item.EnhanceLevel
		result.Message = fmt.Sprintf("强化成功！%s +%d", enhancement.Name, item.EnhanceLevel)
	} else {
		if rand.Float64() < 0.3 && item.EnhanceLevel > 0 {
			item.EnhanceLevel--
			result.Message = "强化失败，装备等级下降"
		} else {
			result.Message = "强化失败"
		}
	}

	result.Stats["attack"] = item.GetAttack()
	result.Stats["defense"] = item.GetDefense()
	result.Stats["hp"] = item.GetHP()

	return result, nil
}

func (m *EquipmentEnhancementManager) calculateSuccessRate(enh *Enhancement, currentLevel int) float64 {
	baseRate := enh.SuccessRate
	levelPenalty := float64(currentLevel) * 0.05
	rate := baseRate - levelPenalty
	if rate < 0.1 {
		rate = 0.1
	}
	return rate
}

// Enchant 装备附魔
func (m *EquipmentEnhancementManager) Enchant(player *Player, item *Item, enchantID string) (*EnhancementResult, error) {
	if player == nil || item == nil {
		return nil, errors.New("玩家或物品不存在")
	}

	enchant, ok := m.enchantments[enchantID]
	if !ok {
		return nil, errors.New("附魔类型不存在")
	}

	for _, e := range item.Enchantments {
		if e == enchantID {
			return nil, errors.New("该附魔已存在")
		}
	}

	if len(item.Enchantments) >= 3 {
		return nil, errors.New("附魔数量已达上限")
	}

	if player.Gem < enchant.Cost {
		return nil, errors.New("钻石不足")
	}

	player.Gem -= enchant.Cost
	success := rand.Float64() < enchant.SuccessRate

	result := &EnhancementResult{
		Success: success,
		Item:    item,
	}

	if success {
		item.Enchantments = append(item.Enchantments, enchantID)
		result.Message = fmt.Sprintf("附魔成功！获得 %s", enchant.Name)
	} else {
		result.Message = "附魔失败"
	}

	return result, nil
}

// GetEnhancement 获取强化信息
func (m *EquipmentEnhancementManager) GetEnhancement(id string) (*Enhancement, bool) {
	e, ok := m.enhancements[id]
	return e, ok
}

// GetEnchantment 获取附魔信息
func (m *EquipmentEnhancementManager) GetEnchantment(id string) (*Enchantment, bool) {
	e, ok := m.enchantments[id]
	return e, ok
}

// ListEnhancements 列出所有强化
func (m *EquipmentEnhancementManager) ListEnhancements() []*Enhancement {
	result := make([]*Enhancement, 0, len(m.enhancements))
	for _, e := range m.enhancements {
		result = append(result, e)
	}
	return result
}

// ListEnchantments 列出所有附魔
func (m *EquipmentEnhancementManager) ListEnchantments() []*Enchantment {
	result := make([]*Enchantment, 0, len(m.enchantments))
	for _, e := range m.enchantments {
		result = append(result, e)
	}
	return result
}

// ============================================================
// 装备合成系统 (Equipment Synthesis System)
// ============================================================

// SynthesisRecipe 合成配方
type SynthesisRecipe struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	ResultID      string     `json:"result_id"`
	ResultCount   int        `json:"result_count"`
	Materials     []Material `json:"materials"`
	GoldCost      int        `json:"gold_cost"`
	SuccessRate   float64    `json:"success_rate"`
	RequiredLevel int        `json:"required_level"`
	Description   string     `json:"description"`
}

// Material 材料
type Material struct {
	ItemID string `json:"item_id"`
	Name   string `json:"name"`
	Count  int    `json:"count"`
}

// SynthesisResult 合成结果
type SynthesisResult struct {
	Success     bool        `json:"success"`
	ResultItem  *Item      `json:"result_item,omitempty"`
	ResultCount int         `json:"result_count"`
	Materials   []Material `json:"materials"`
	Message     string      `json:"message"`
}

// EquipmentSynthesisManager 装备合成管理器
type EquipmentSynthesisManager struct {
	recipes map[string]*SynthesisRecipe
}

// NewEquipmentSynthesisManager 创建合成管理器
func NewEquipmentSynthesisManager() *EquipmentSynthesisManager {
	m := &EquipmentSynthesisManager{
		recipes: make(map[string]*SynthesisRecipe),
	}
	m.initDefaultRecipes()
	return m
}

func (m *EquipmentSynthesisManager) initDefaultRecipes() {
	recipes := []*SynthesisRecipe{
		{
			ID: "recipe_weapon_1", Name: "打造铁剑",
			ResultID: "weapon_iron_sword", ResultCount: 1,
			Materials: []Material{
				{ItemID: "iron_ingot", Name: "铁锭", Count: 5},
				{ItemID: "wood", Name: "木材", Count: 3},
			},
			GoldCost: 100, SuccessRate: 0.9, RequiredLevel: 1, Description: "打造基础武器",
		},
		{
			ID: "recipe_weapon_2", Name: "打造银剑",
			ResultID: "weapon_silver_sword", ResultCount: 1,
			Materials: []Material{
				{ItemID: "silver_ingot", Name: "银锭", Count: 5},
				{ItemID: "weapon_iron_sword", Name: "铁剑", Count: 1},
				{ItemID: "gem_ruby", Name: "红宝石", Count: 1},
			},
			GoldCost: 500, SuccessRate: 0.7, RequiredLevel: 5, Description: "打造进阶武器",
		},
		{
			ID: "recipe_weapon_3", Name: "打造金剑",
			ResultID: "weapon_golden_sword", ResultCount: 1,
			Materials: []Material{
				{ItemID: "gold_ingot", Name: "金锭", Count: 5},
				{ItemID: "weapon_silver_sword", Name: "银剑", Count: 1},
				{ItemID: "gem_diamond", Name: "钻石", Count: 2},
			},
			GoldCost: 2000, SuccessRate: 0.5, RequiredLevel: 10, Description: "打造高级武器",
		},
		{
			ID: "recipe_armor_1", Name: "打造铁甲",
			ResultID: "armor_iron", ResultCount: 1,
			Materials: []Material{
				{ItemID: "iron_ingot", Name: "铁锭", Count: 8},
				{ItemID: "leather", Name: "皮革", Count: 3},
			},
			GoldCost: 150, SuccessRate: 0.9, RequiredLevel: 1, Description: "打造基础防具",
		},
		{
			ID: "recipe_potion_health", Name: "制作生命药水",
			ResultID: "potion_health", ResultCount: 3,
			Materials: []Material{
				{ItemID: "herb_red", Name: "红色药草", Count: 2},
				{ItemID: "bottle", Name: "空瓶", Count: 1},
			},
			GoldCost: 50, SuccessRate: 0.95, RequiredLevel: 1, Description: "制作恢复药水",
		},
		{
			ID: "recipe_potion_mana", Name: "制作魔法药水",
			ResultID: "potion_mana", ResultCount: 3,
			Materials: []Material{
				{ItemID: "herb_blue", Name: "蓝色药草", Count: 2},
				{ItemID: "bottle", Name: "空瓶", Count: 1},
			},
			GoldCost: 50, SuccessRate: 0.95, RequiredLevel: 1, Description: "制作魔法药水",
		},
		{
			ID: "recipe_gem_1", Name: "合成初级宝石",
			ResultID: "gem_ruby", ResultCount: 1,
			Materials: []Material{
				{ItemID: "gem_shard_red", Name: "红宝石碎片", Count: 5},
			},
			GoldCost: 100, SuccessRate: 0.8, RequiredLevel: 3, Description: "合成红宝石",
		},
		{
			ID: "recipe_gem_2", Name: "合成中级宝石",
			ResultID: "gem_diamond", ResultCount: 1,
			Materials: []Material{
				{ItemID: "gem_shard_white", Name: "宝石碎片", Count: 10},
				{ItemID: "gem_ruby", Name: "红宝石", Count: 2},
			},
			GoldCost: 500, SuccessRate: 0.6, RequiredLevel: 7, Description: "合成钻石",
		},
		{
			ID: "recipe_elixir_1", Name: "合成力量药剂",
			ResultID: "elixir_power", ResultCount: 1,
			Materials: []Material{
				{ItemID: "herb_red", Name: "红色药草", Count: 5},
				{ItemID: "herb_gold", Name: "金色药草", Count: 1},
				{ItemID: "gem_shard_red", Name: "红宝石碎片", Count: 3},
			},
			GoldCost: 300, SuccessRate: 0.7, RequiredLevel: 5, Description: "临时提升攻击力",
		},
		{
			ID: "recipe_accessory_1", Name: "制作力量戒指",
			ResultID: "ring_power", ResultCount: 1,
			Materials: []Material{
				{ItemID: "gold_ingot", Name: "金锭", Count: 3},
				{ItemID: "gem_ruby", Name: "红宝石", Count: 1},
			},
			GoldCost: 400, SuccessRate: 0.75, RequiredLevel: 4, Description: "增加攻击力的戒指",
		},
	}
	for _, r := range recipes {
		m.recipes[r.ID] = r
	}
}

// Synthesize 执行合成
func (m *EquipmentSynthesisManager) Synthesize(player *Player, recipeID string) (*SynthesisResult, error) {
	if player == nil {
		return nil, errors.New("玩家不存在")
	}

	recipe, ok := m.recipes[recipeID]
	if !ok {
		return nil, errors.New("配方不存在")
	}

	result := &SynthesisResult{
		Materials: recipe.Materials,
	}

	materialsMap := make(map[string]int)
	for _, mat := range player.Inventory.Items {
		materialsMap[mat.ID] += mat.Count
	}

	for _, mat := range recipe.Materials {
		if materialsMap[mat.ItemID] < mat.Count {
			result.Message = fmt.Sprintf("材料不足：需要 %s x%d", mat.Name, mat.Count)
			return result, nil
		}
	}

	if player.Gold < recipe.GoldCost {
		return nil, errors.New("金币不足")
	}

	if player.Level < recipe.RequiredLevel {
		return nil, errors.New(fmt.Sprintf("需要角色等级 %d", recipe.RequiredLevel))
	}

	for _, mat := range recipe.Materials {
		player.RemoveItem(mat.ItemID, mat.Count)
	}

	player.Gold -= recipe.GoldCost

	success := rand.Float64() < recipe.SuccessRate
	result.Success = success

	if success {
		item := NewItem(recipe.ResultID)
		if item != nil {
			item.Count = recipe.ResultCount
			player.AddItem(item)
			result.ResultItem = item
			result.ResultCount = recipe.ResultCount
			result.Message = fmt.Sprintf("合成成功！获得 %s x%d", item.Name, recipe.ResultCount)
		} else {
			result.Message = "合成成功，但物品创建失败"
		}
	} else {
		result.Message = "合成失败，材料已消耗"
	}

	return result, nil
}

// GetRecipe 获取配方
func (m *EquipmentSynthesisManager) GetRecipe(id string) (*SynthesisRecipe, bool) {
	r, ok := m.recipes[id]
	return r, ok
}

// ListRecipes 列出所有配方
func (m *EquipmentSynthesisManager) ListRecipes() []*SynthesisRecipe {
	result := make([]*SynthesisRecipe, 0, len(m.recipes))
	for _, r := range m.recipes {
		result = append(result, r)
	}
	return result
}

// GetRecipesByResult 根据产物ID获取配方
func (m *EquipmentSynthesisManager) GetRecipesByResult(resultID string) []*SynthesisRecipe {
	result := make([]*SynthesisRecipe, 0)
	for _, r := range m.recipes {
		if r.ResultID == resultID {
			result = append(result, r)
		}
	}
	return result
}

// CanSynthesize 检查是否可以合成
func (m *EquipmentSynthesisManager) CanSynthesize(player *Player, recipeID string) (bool, string) {
	recipe, ok := m.recipes[recipeID]
	if !ok {
		return false, "配方不存在"
	}

	if player.Level < recipe.RequiredLevel {
		return false, fmt.Sprintf("需要角色等级 %d", recipe.RequiredLevel)
	}

	if player.Gold < recipe.GoldCost {
		return false, "金币不足"
	}

	materialsMap := make(map[string]int)
	for _, mat := range player.Inventory.Items {
		materialsMap[mat.ID] += mat.Count
	}

	for _, mat := range recipe.Materials {
		if materialsMap[mat.ItemID] < mat.Count {
			return false, fmt.Sprintf("材料不足：需要 %s x%d", mat.Name, mat.Count)
		}
	}

	return true, "可以合成"
}

// ============================================================
// 装备精炼系统 (Equipment Refinement System)
// ============================================================

// RefineResult 精炼结果
type RefineResult struct {
	Success    bool              `json:"success"`
	NewQuality int               `json:"new_quality"`
	Stats      map[string]float64 `json:"stats"`
	Message    string            `json:"message"`
}

// RefineEquipment 装备精炼
func (m *EquipmentEnhancementManager) RefineEquipment(player *Player, item *Item) (*RefineResult, error) {
	if player == nil || item == nil {
		return nil, errors.New("玩家或物品不存在")
	}

	if item.Type != ItemTypeWeapon && item.Type != ItemTypeArmor && item.Type != ItemTypeAccessory {
		return nil, errors.New("该类型装备无法精炼")
	}

	cost := item.Quality * 100
	if player.Gem < cost {
		return nil, errors.New("钻石不足")
	}

	baseRate := 0.8
	successRate := baseRate - float64(item.Quality)*0.05
	if successRate < 0.2 {
		successRate = 0.2
	}

	player.Gem -= cost
	success := rand.Float64() < successRate

	result := &RefineResult{
		Stats: make(map[string]float64),
	}

	if success {
		item.Quality++
		result.NewQuality = item.Quality
		result.Success = true
		result.Message = fmt.Sprintf("精炼成功！品质提升至 +%d", item.Quality)
	} else {
		result.NewQuality = item.Quality
		result.Success = false
		result.Message = "精炼失败"
	}

	qualityBonus := float64(item.Quality) * 0.1
	result.Stats["quality_bonus"] = qualityBonus
	result.Stats["attack"] = item.GetAttack()
	result.Stats["defense"] = item.GetDefense()
	result.Stats["hp"] = item.GetHP()

	return result, nil
}

// CalculateRefineCost 计算精炼费用
func (m *EquipmentEnhancementManager) CalculateRefineCost(item *Item) int {
	return item.Quality * 100
}

// GetRefineSuccessRate 获取精炼成功率
func (m *EquipmentEnhancementManager) GetRefineSuccessRate(item *Item) float64 {
	baseRate := 0.8
	successRate := baseRate - float64(item.Quality)*0.05
	if successRate < 0.2 {
		return 0.2
	}
	return math.Round(successRate*100) / 100
}

// ============================================================
// 装备分解系统 (Equipment Disassembly System)
// ============================================================

// DisassemblyResult 分解结果
type DisassemblyResult struct {
	Items   map[string]int `json:"items"`
	Gold    int            `json:"gold"`
	Message string         `json:"message"`
}

// DisassembleEquipment 装备分解
func (m *EquipmentSynthesisManager) DisassembleEquipment(player *Player, item *Item) (*DisassemblyResult, error) {
	if player == nil || item == nil {
		return nil, errors.New("玩家或物品不存在")
	}

	result := &DisassemblyResult{
		Items: make(map[string]int),
	}

	goldReward := item.Quality * 50
	result.Gold = goldReward
	player.Gold += goldReward

	switch item.Type {
	case ItemTypeWeapon:
		result.Items["iron_ingot"] = item.Quality
		if item.Quality >= 3 {
			result.Items["gem_shard_red"] = 1
		}
	case ItemTypeArmor:
		result.Items["iron_ingot"] = item.Quality
		if item.Quality >= 3 {
			result.Items["gem_shard_blue"] = 1
		}
	case ItemTypeAccessory:
		result.Items["gold_ingot"] = 1
		if item.Quality >= 2 {
			result.Items["gem_shard_white"] = 1
		}
	default:
		result.Items["metal_scrap"] = item.Quality
	}

	for itemID, count := range result.Items {
		material := NewItem(itemID)
		if material != nil {
			material.Count = count
			player.AddItem(material)
		}
	}

	player.RemoveItem(item.ID, 1)

	result.Message = fmt.Sprintf("分解成功！获得 %d 金币", goldReward)
	return result, nil
}

// BatchDisassemble 批量分解
func (m *EquipmentSynthesisManager) BatchDisassemble(player *Player, itemIDs []string) ([]*DisassemblyResult, error) {
	results := make([]*DisassemblyResult, 0)

	for _, itemID := range itemIDs {
		item := player.GetItem(itemID)
		if item != nil {
			result, err := m.DisassembleEquipment(player, item)
			if err != nil {
				continue
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// ============================================================
// 装备回收系统 (Equipment Recycle System)
// ============================================================

// RecycleResult 回收结果
type RecycleResult struct {
	Gem     int    `json:"gem"`
	Exp     int    `json:"exp"`
	Message string `json:"message"`
}

// RecycleEquipment 装备回收
func (m *EquipmentEnhancementManager) RecycleEquipment(player *Player, item *Item) (*RecycleResult, error) {
	if player == nil || item == nil {
		return nil, errors.New("玩家或物品不存在")
	}

	result := &RecycleResult{}

	expReward := item.Quality * 10
	result.Exp = expReward

	gemReward := 0
	if item.Quality >= 2 {
		gemReward = item.Quality * 5
	}
	result.Gem = gemReward

	player.AddExp(expReward)
	if gemReward > 0 {
		player.Gem += gemReward
	}

	player.RemoveItem(item.ID, 1)

	result.Message = fmt.Sprintf("回收成功！获得 %d 经验和 %d 钻石", expReward, gemReward)
	return result, nil
}

// ============================================================
// 装备打孔系统 (Equipment Socket System)
// ============================================================

// Socket 打孔位
type Socket struct {
	Index   int    `json:"index"`
	GemID   string `json:"gem_id,omitempty"`
	IsEmpty bool   `json:"is_empty"`
}

// AddSocket 打孔
func (m *EquipmentEnhancementManager) AddSocket(player *Player, item *Item) (*Item, error) {
	if player == nil || item == nil {
		return nil, errors.New("玩家或物品不存在")
	}

	cost := 200
	if player.Gem < cost {
		return nil, errors.New("钻石不足")
	}

	if len(item.Sockets) >= 3 {
		return nil, errors.New("打孔数已达上限")
	}

	player.Gem -= cost

	socket := Socket{
		Index:   len(item.Sockets),
		IsEmpty: true,
	}
	item.Sockets = append(item.Sockets, socket)

	return item, nil
}

// InsertGem 镶嵌宝石
func (m *EquipmentEnhancementManager) InsertGem(player *Player, item *Item, socketIndex int, gemID string) error {
	if player == nil || item == nil {
		return errors.New("玩家或物品不存在")
	}

	if socketIndex < 0 || socketIndex >= len(item.Sockets) {
		return errors.New("打孔位不存在")
	}

	socket := &item.Sockets[socketIndex]
	if !socket.IsEmpty {
		return errors.New("打孔位已有宝石")
	}

	gem := player.GetItem(gemID)
	if gem == nil {
		return errors.New("宝石不存在")
	}

	player.RemoveItem(gemID, 1)

	socket.GemID = gemID
	socket.IsEmpty = false

	return nil
}

// RemoveGem 取下宝石
func (m *EquipmentEnhancementManager) RemoveGem(player *Player, item *Item, socketIndex int) (*Item, error) {
	if player == nil || item == nil {
		return nil, errors.New("玩家或物品不存在")
	}

	if socketIndex < 0 || socketIndex >= len(item.Sockets) {
		return nil, errors.New("打孔位不存在")
	}

	socket := &item.Sockets[socketIndex]
	if socket.IsEmpty {
		return nil, errors.New("打孔位为空")
	}

	gem := NewItem(socket.GemID)
	if gem != nil {
		gem.Count = 1
		player.AddItem(gem)
	}

	socket.GemID = ""
	socket.IsEmpty = true

	return item, nil
}

// ============================================================
// 装备耐久度系统 (Equipment Durability System)
// ============================================================

// RepairEquipment 装备修理
func (m *EquipmentEnhancementManager) RepairEquipment(player *Player, item *Item) (int, error) {
	if player == nil || item == nil {
		return 0, errors.New("玩家或物品不存在")
	}

	if item.Durability >= item.MaxDurability {
		return 0, errors.New("装备无需修理")
	}

	lostDurability := item.MaxDurability - item.Durability
	cost := lostDurability * 10

	if player.Gold < cost {
		return 0, errors.New("金币不足")
	}

	player.Gold -= cost
	item.Durability = item.MaxDurability

	return cost, nil
}

// DamageEquipment 装备损坏
func (m *EquipmentEnhancementManager) DamageEquipment(item *Item, damage int) {
	if item == nil {
		return
	}
	item.Durability -= damage
	if item.Durability < 0 {
		item.Durability = 0
	}
}

// CheckEquipmentBroken 检查装备是否损坏
func (m *EquipmentEnhancementManager) CheckEquipmentBroken(item *Item) bool {
	if item == nil {
		return false
	}
	return item.Durability <= 0
}

// GetEquipmentDurabilityPercent 获取装备耐久度百分比
func (m *EquipmentEnhancementManager) GetEquipmentDurabilityPercent(item *Item) float64 {
	if item == nil || item.MaxDurability == 0 {
		return 0
	}
	return float64(item.Durability) / float64(item.MaxDurability) * 100
}

// ============================================================
// 全局管理器实例
// ============================================================

var (
	enhancementMgr *EquipmentEnhancementManager
	synthesisMgr   *EquipmentSynthesisManager
)

func init() {
	enhancementMgr = NewEquipmentEnhancementManager()
	synthesisMgr = NewEquipmentSynthesisManager()
}

// GetEnhancementManager 获取强化管理器
func GetEnhancementManager() *EquipmentEnhancementManager {
	return enhancementMgr
}

// GetSynthesisManager 获取合成管理器
func GetSynthesisManager() *EquipmentSynthesisManager {
	return synthesisMgr
}
