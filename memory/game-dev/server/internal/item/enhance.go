package item

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ==================== 装备强化系统 ====================

// EnhanceType 强化类型
type EnhanceType int

const (
	EnhanceTypeAttack EnhanceType = iota + 1 // 攻击强化
	EnhanceTypeDefense                         // 防御强化
	EnhanceTypeHP                              // 生命强化
	EnhanceTypeSpeed                           // 速度强化
	EnhanceTypeCrit                            // 暴击强化
	EnhanceTypeDodge                           // 闪避强化
)

// EnhanceResult 强化结果
type EnhanceResult int

const (
	EnhanceResultSuccess EnhanceResult = iota // 成功
	EnhanceResultFail                         // 失败
	EnhanceResultBreak                        // 装备破碎
)

// EquipmentEnhance 装备强化记录
type EquipmentEnhance struct {
	PlayerID      string        // 玩家ID
	ItemID        string        // 装备ID
	EnhanceLevel  int           // 强化等级
	EnhanceType   EnhanceType   // 强化类型
	EnhanceExp    int64         // 强化经验
	TotalCost     int64         // 总消耗金币
	SuccessCount  int           // 成功次数
	FailCount     int           // 失败次数
	LastEnhanceAt time.Time     // 上次强化时间
}

// EnhanceConfig 强化配置
type EnhanceConfig struct {
	BaseSuccessRate   float64   // 基础成功率
	LevelFailRate     float64   // 每级失败率加成
	BreakRate         float64   // 破碎率
	MaxLevel          int       // 最大强化等级
	EnhanceCostBase   int64     // 基础消耗
	EnhanceCostFactor float64   // 消耗系数
	ExpPerLevel       int64     // 每级所需经验
}

// DefaultEnhanceConfig 默认强化配置
var DefaultEnhanceConfig = EnhanceConfig{
	BaseSuccessRate:   0.95,
	LevelFailRate:     0.02,
	BreakRate:         0.01,
	MaxLevel:          15,
	EnhanceCostBase:   1000,
	EnhanceCostFactor: 1.5,
	ExpPerLevel:       1000,
}

// EquipmentEnhanceSystem 装备强化系统
type EquipmentEnhanceSystem struct {
	config      EnhanceConfig
	playerData  map[string]map[string]*EquipmentEnhance // playerID -> itemID -> enhance
}

// NewEquipmentEnhanceSystem 创建装备强化系统
func NewEquipmentEnhanceSystem() *EquipmentEnhanceSystem {
	return &EquipmentEnhanceSystem{
		config:     DefaultEnhanceConfig,
		playerData: make(map[string]map[string]*EquipmentEnhance),
	}
}

// GetEnhanceCost 获取强化消耗
func (e *EquipmentEnhanceSystem) GetEnhanceCost(level int) int64 {
	cost := float64(e.config.EnhanceCostBase) * math.Pow(e.config.EnhanceCostFactor, float64(level))
	return int64(cost)
}

// CalculateSuccessRate 计算成功率
func (e *EquipmentEnhanceSystem) CalculateSuccessRate(level int) float64 {
	failRate := e.config.LevelFailRate * float64(level)
	successRate := e.config.BaseSuccessRate - failRate
	if successRate < 0.1 {
		successRate = 0.1
	}
	return successRate
}

// Enhance 强化装备
func (e *EquipmentEnhanceSystem) Enhance(playerID, itemID string, enhanceType EnhanceType, gold int64) *EnhanceResultInfo {
	// 获取或创建强化记录
	enhance := e.getOrCreateEnhance(playerID, itemID)

	// 检查是否已达最大等级
	if enhance.EnhanceLevel >= e.config.MaxLevel {
		return &EnhanceResultInfo{
			Result:    EnhanceResultFail,
			Message:   fmt.Sprintf("已达到最大强化等级 %d", e.config.MaxLevel),
			NewLevel:  enhance.EnhanceLevel,
			ExpGained: 0,
		}
	}

	// 计算消耗
	cost := e.GetEnhanceCost(enhance.EnhanceLevel)
	if gold < cost {
		return &EnhanceResultInfo{
			Result:    EnhanceResultFail,
			Message:   "金币不足",
			NewLevel:  enhance.EnhanceLevel,
			ExpGained: 0,
		}
	}

	// 计算成功率
	successRate := e.CalculateSuccessRate(enhance.EnhanceLevel)

	// 随机判定
	random := rand.Float64()
	var result EnhanceResult
	var expGained int64

	if random < successRate {
		// 成功
		result = EnhanceResultSuccess
		enhance.EnhanceLevel++
		enhance.SuccessCount++
		expGained = e.config.ExpPerLevel
	} else if random < successRate+e.config.BreakRate {
		// 破碎
		result = EnhanceResultBreak
		enhance.EnhanceLevel = 0
		enhance.FailCount++
		expGained = e.config.ExpPerLevel / 2
	} else {
		// 失败但不破碎
		result = EnhanceResultFail
		enhance.FailCount++
		expGained = e.config.ExpPerLevel
	}

	// 更新强化记录
	enhance.EnhanceType = enhanceType
	enhance.TotalCost += cost
	enhance.EnhanceExp += expGained
	enhance.LastEnhanceAt = time.Now()

	// 生成结果信息
	message := e.getResultMessage(result, enhance.EnhanceLevel)
	return &EnhanceResultInfo{
		Result:        result,
		Message:       message,
		NewLevel:     enhance.EnhanceLevel,
		ExpGained:    expGained,
		TotalCost:    enhance.TotalCost,
		SuccessRate:  successRate,
		SuccessCount: enhance.SuccessCount,
		FailCount:    enhance.FailCount,
	}
}

// EnhanceResultInfo 强化结果信息
type EnhanceResultInfo struct {
	Result        EnhanceResult
	Message       string
	NewLevel      int
	ExpGained     int64
	TotalCost     int64
	SuccessRate   float64
	SuccessCount  int
	FailCount     int
}

// getResultMessage 获取结果消息
func (e *EquipmentEnhanceSystem) getResultMessage(result EnhanceResult, level int) string {
	switch result {
	case EnhanceResultSuccess:
		return fmt.Sprintf("强化成功！装备强化等级提升至 %d", level)
	case EnhanceResultFail:
		return fmt.Sprintf("强化失败，装备等级不变（当前等级: %d）", level)
	case EnhanceResultBreak:
		return fmt.Sprintf("装备破碎！强化等级重置为 0")
	default:
		return "未知结果"
	}
}

// getOrCreateEnhance 获取或创建强化记录
func (e *EquipmentEnhanceSystem) getOrCreateEnhance(playerID, itemID string) *EquipmentEnhance {
	if e.playerData[playerID] == nil {
		e.playerData[playerID] = make(map[string]*EquipmentEnhance)
	}

	if e.playerData[playerID][itemID] == nil {
		e.playerData[playerID][itemID] = &EquipmentEnhance{
			PlayerID:      playerID,
			ItemID:        itemID,
			EnhanceLevel:  0,
			LastEnhanceAt: time.Now(),
		}
	}

	return e.playerData[playerID][itemID]
}

// GetEnhanceInfo 获取强化信息
func (e *EquipmentEnhanceSystem) GetEnhanceInfo(playerID, itemID string) *EquipmentEnhance {
	if e.playerData[playerID] == nil {
		return nil
	}
	return e.playerData[playerID][itemID]
}

// GetAllEnhances 获取玩家所有强化记录
func (e *EquipmentEnhanceSystem) GetAllEnhances(playerID string) []*EquipmentEnhance {
	if e.playerData[playerID] == nil {
		return nil
	}

	enhances := make([]*EquipmentEnhance, 0, len(e.playerData[playerID]))
	for _, enhance := range e.playerData[playerID] {
		enhances = append(enhances, enhance)
	}
	return enhances
}

// CalculateEnhanceBonus 计算强化加成
func (e *EquipmentEnhanceSystem) CalculateEnhanceBonus(level int, baseValue float64, enhanceType EnhanceType) float64 {
	// 每级增加 5% 加成
	bonusRate := 1.0 + float64(level)*0.05
	return baseValue * bonusRate
}

// GetEnhanceStats 获取强化属性
func (e *EquipmentEnhanceSystem) GetEnhanceStats(level int, enhanceType EnhanceType) map[string]float64 {
	stats := make(map[string]float64)

	switch enhanceType {
	case EnhanceTypeAttack:
		stats["attack"] = float64(level * 10)
	case EnhanceTypeDefense:
		stats["defense"] = float64(level * 8)
	case EnhanceTypeHP:
		stats["hp"] = float64(level * 50)
	case EnhanceTypeSpeed:
		stats["speed"] = float64(level * 2)
	case EnhanceTypeCrit:
		stats["crit_rate"] = float64(level) * 0.5
	case EnhanceTypeDodge:
		stats["dodge_rate"] = float64(level) * 0.5
	}

	return stats
}

// ResetEnhance 重置强化
func (e *EquipmentEnhanceSystem) ResetEnhance(playerID, itemID string) bool {
	if e.playerData[playerID] == nil {
		return false
	}

	// 重置需要消耗钻石
	resetCost := int64(100)

	enhance := e.playerData[playerID][itemID]
	if enhance == nil {
		return false
	}

	// 重置等级和经验
	enhance.EnhanceLevel = 0
	enhance.EnhanceExp = 0
	enhance.TotalCost = 0

	return true
}

// TransferEnhance 转移强化（装备替换）
func (e *EquipmentEnhanceSystem) TransferEnhance(playerID, fromItemID, toItemID string) bool {
	fromEnhance := e.GetEnhanceInfo(playerID, fromItemID)
	if fromEnhance == nil || fromEnhance.EnhanceLevel == 0 {
		return false
	}

	// 创建新装备的强化记录
	toEnhance := e.getOrCreateEnhance(playerID, toItemID)
	toEnhance.EnhanceLevel = fromEnhance.EnhanceLevel
	toEnhance.EnhanceType = fromEnhance.EnhanceType
	toEnhance.EnhanceExp = fromEnhance.EnhanceExp

	// 清除原装备强化记录
	delete(e.playerData[playerID], fromItemID)

	return true
}

// GetEnhanceLeaderboard 获取强化排行榜
func (e *EquipmentEnhanceSystem) GetEnhanceLeaderboard(limit int) []*EnhanceLeaderboardEntry {
	type entry struct {
		playerID     string
		maxLevel     int
		totalCost    int64
		successCount int
	}

	playerStats := make(map[string]*entry)

	for playerID, items := range e.playerData {
		for _, enhance := range items {
			if _, ok := playerStats[playerID]; !ok {
				playerStats[playerID] = &entry{
					playerID:     playerID,
					maxLevel:     0,
					totalCost:    0,
					successCount: 0,
				}
			}

			if enhance.EnhanceLevel > playerStats[playerID].maxLevel {
				playerStats[playerID].maxLevel = enhance.EnhanceLevel
			}
			playerStats[playerID].totalCost += enhance.TotalCost
			playerStats[playerID].successCount += enhance.SuccessCount
		}
	}

	// 排序
	sorted := make([]*entry, 0, len(playerStats))
	for _, v := range playerStats {
		sorted = append(sorted, v)
	}

	// 按最大强化等级排序
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].maxLevel > sorted[i].maxLevel {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// 取前 N 名
	if limit > len(sorted) {
		limit = len(sorted)
	}

	result := make([]*EnhanceLeaderboardEntry, limit)
	for i := 0; i < limit; i++ {
		result[i] = &EnhanceLeaderboardEntry{
			Rank:         i + 1,
			PlayerID:     sorted[i].playerID,
			MaxLevel:     sorted[i].maxLevel,
			TotalCost:    sorted[i].totalCost,
			SuccessCount: sorted[i].successCount,
		}
	}

	return result
}

// EnhanceLeaderboardEntry 强化排行榜条目
type EnhanceLeaderboardEntry struct {
	Rank         int
	PlayerID     string
	MaxLevel     int
	TotalCost    int64
	SuccessCount int
}
