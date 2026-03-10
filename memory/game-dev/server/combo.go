package main

import (
	"math"
	"sync"
	"time"
)

// ============================================
// 连击系统 (Combo System)
// 弹幕游戏核心评分系统
// ============================================

// ComboType 连击类型
type ComboType int

const (
	ComboTypeHit ComboType = iota // 击中
	ComboTypeKill                // 击杀
	ComboTypeDodge               // 闪避
	ComboTypePerfect             // 完美闪避
	ComboTypeSkill               // 技能命中
	ComboTypeBoss                // Boss击杀
	ComboTypeCollection          // 道具收集
)

// ComboEvent 连击事件
type ComboEvent struct {
	Type       ComboType `json:"type"`
	Points     int       `json:"points"`     // 基础分数
	Multiplier float64   `json:"multiplier"` // 乘数
	Timestamp  int64     `json:"timestamp"`
	SourceID   string    `json:"sourceId"`   // 来源ID
	TargetID   string    `json:"targetId"`   // 目标ID
}

// Combo 连击记录
type Combo struct {
	ID           string       `json:"id"`
	PlayerID     string       `json:"playerId"`
	Count        int          `json:"count"`         // 连击数
	MaxCount     int          `json:"maxCount"`      // 最大连击数
	Score        int64        `json:"score"`         // 总分数
	BaseScore    int64        `json:"baseScore"`    // 基础分数
	Multiplier   float64      `json:"multiplier"`   // 当前乘数
	MaxMultiplier float64     `json:"maxMultiplier"` // 最大乘数
	Events       []ComboEvent `json:"events"`       // 事件列表
	LastEvent    int64        `json:"lastEvent"`    // 上次事件时间
	Timeout      int64        `json:"timeout"`      // 连击超时时间(ms)
	IsActive     bool         `json:"isActive"`     // 是否活跃
	StartTime    int64        `json:"startTime"`    // 开始时间
	EndTime      int64        `json:"endTime"`      // 结束时间
	mu           sync.RWMutex
}

// ComboConfig 连击配置
type ComboConfig struct {
	BasePoints       map[ComboType]int   `json:"basePoints"`       // 基础分数
	MaxCombo         int                 `json:"maxCombo"`         // 最大连击数
	Timeout          int64               `json:"timeout"`          // 连击超时(ms)
	MaxMultiplier    float64             `json:"maxMultiplier"`   // 最大乘数
	MultiplierCurve  string              `json:"multiplierCurve"`  // 乘数曲线: linear, exponential,阶梯
	PerfectDodgeTime int64               `json:"perfectDodgeTime"` // 完美闪避判定时间(ms)
}

// 默认连击配置
var DefaultComboConfig = ComboConfig{
	BasePoints: map[ComboType]int{
		ComboTypeHit:       10,
		ComboTypeKill:      100,
		ComboTypeDodge:     50,
		ComboTypePerfect:   200,
		ComboTypeSkill:     150,
		ComboTypeBoss:      1000,
		ComboTypeCollection: 25,
	},
	MaxCombo:         999,
	Timeout:          2000,
	MaxMultiplier:    10.0,
	MultiplierCurve:  "exponential",
	PerfectDodgeTime: 200,
}

// ComboManager 连击管理器
type ComboManager struct {
	Combos   map[string]*Combo     // playerID -> Combo
	Config   *ComboConfig
	mu       sync.RWMutex
	Pool     *sync.Pool
}

// 创建连击管理器
func NewComboManager() *ComboManager {
	return &ComboManager{
		Combos: make(map[string]*Combo),
		Config: &DefaultComboConfig,
		Pool: &sync.Pool{
			New: func() interface{} {
				return &Combo{
					Events: make([]ComboEvent, 0, 100),
				}
			},
		},
	}
}

// 创建连击记录
func (cm *ComboManager) CreateCombo(playerID string) *Combo {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	combo := cm.Pool.Get().(*Combo)
	combo.ID = playerID + "_" + string(rune(time.Now().Unix()))
	combo.PlayerID = playerID
	combo.Count = 0
	combo.MaxCount = 0
	combo.Score = 0
	combo.BaseScore = 0
	combo.Multiplier = 1.0
	combo.MaxMultiplier = 1.0
	combo.Events = combo.Events[:0]
	combo.LastEvent = 0
	combo.Timeout = cm.Config.Timeout
	combo.IsActive = false
	combo.StartTime = time.Now().UnixMilli()
	combo.EndTime = 0
	
	cm.Combos[playerID] = combo
	return combo
}

// 获取或创建连击
func (cm *ComboManager) GetOrCreateCombo(playerID string) *Combo {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if combo, ok := cm.Combos[playerID]; ok {
		return combo
	}
	return cm.CreateCombo(playerID)
}

// 计算乘数
func (cm *ComboManager) calculateMultiplier(count int) float64 {
	maxMult := cm.Config.MaxMultiplier
	
	switch cm.Config.MultiplierCurve {
	case "linear":
		// 线性: 1 + count * 0.1
		return math.Min(1.0+float64(count)*0.1, maxMult)
	case "exponential":
		// 指数: 1.5^count
		return math.Min(math.Pow(1.5, float64(count)/10.0), maxMult)
	case "step":
		// 阶梯: 每50连击+1倍
		return math.Min(1.0+float64(count/50), maxMult)
	default:
		return math.Min(1.0+float64(count)*0.1, maxMult)
	}
}

// 添加连击事件
func (cm *ComboManager) AddEvent(playerID string, eventType ComboType, sourceID, targetID string) *Combo {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	now := time.Now().UnixMilli()
	
	// 获取或创建连击
	combo, ok := cm.Combos[playerID]
	if !ok {
		combo = cm.Pool.Get().(*Combo)
		combo.ID = playerID + "_" + string(rune(now))
		combo.PlayerID = playerID
		combo.Events = make([]ComboEvent, 0, 100)
		combo.StartTime = now
		cm.Combos[playerID] = combo
	}
	
	// 检查超时
	if now-combo.LastEvent > combo.Timeout {
		// 连击中断，保存最大连击
		if combo.Count > combo.MaxCount {
			combo.MaxCount = combo.Count
		}
		combo.Count = 0
		combo.Multiplier = 1.0
		combo.IsActive = false
	}
	
	// 计算分数
	basePoints := cm.Config.BasePoints[eventType]
	
	// 击杀和Boss有额外加分
	var extraMultiplier float64 = 1.0
	switch eventType {
	case ComboTypeKill:
		extraMultiplier = 1.5
	case ComboTypeBoss:
		extraMultiplier = 2.0
	case ComboTypePerfect:
		extraMultiplier = 2.0
	}
	
	// 更新连击
	combo.Count++
	if combo.Count > combo.MaxCount {
		combo.MaxCount = combo.Count
	}
	
	// 计算新乘数
	combo.Multiplier = cm.calculateMultiplier(combo.Count)
	if combo.Multiplier > combo.MaxMultiplier {
		combo.MaxMultiplier = combo.Multiplier
	}
	
	// 计算分数
	points := float64(basePoints) * combo.Multiplier * extraMultiplier
	combo.BaseScore += int64(basePoints)
	combo.Score += int64(points)
	
	// 记录事件
	event := ComboEvent{
		Type:       eventType,
		Points:     basePoints,
		Multiplier: combo.Multiplier,
		Timestamp:  now,
		SourceID:   sourceID,
		TargetID:   targetID,
	}
	combo.Events = append(combo.Events, event)
	
	// 保持事件列表合理大小
	if len(combo.Events) > 1000 {
		combo.Events = combo.Events[len(combo.Events)-500:]
	}
	
	combo.LastEvent = now
	combo.IsActive = true
	combo.EndTime = now
	
	return combo
}

// 简化接口
func (cm *ComboManager) AddHit(playerID string) *Combo {
	return cm.AddEvent(playerID, ComboTypeHit, "", "")
}

func (cm *ComboManager) AddKill(playerID string, targetID string) *Combo {
	return cm.AddEvent(playerID, ComboTypeKill, "", targetID)
}

func (cm *ComboManager) AddDodge(playerID string) *Combo {
	return cm.AddEvent(playerID, ComboTypeDodge, "", "")
}

func (cm *ComboManager) AddPerfectDodge(playerID string) *Combo {
	return cm.AddEvent(playerID, ComboTypePerfect, "", "")
}

func (cm *ComboManager) AddBossKill(playerID string, bossID string) *Combo {
	return cm.AddEvent(playerID, ComboTypeBoss, "", bossID)
}

func (cm *ComboManager) AddCollection(playerID string, itemID string) *Combo {
	return cm.AddEvent(playerID, ComboTypeCollection, itemID, "")
}

// 获取连击状态
func (cm *ComboManager) GetCombo(playerID string) *Combo {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.Combos[playerID]
}

// 更新连击状态 (检查超时)
func (cm *ComboManager) Update() map[string]*Combo {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	now := time.Now().UnixMilli()
	expired := make(map[string]*Combo)
	
	for playerID, combo := range cm.Combos {
		if combo.IsActive && now-combo.LastEvent > combo.Timeout {
			// 连击结束
			combo.IsActive = false
			combo.EndTime = now
			expired[playerID] = combo
		}
	}
	
	return expired
}

// 重置连击
func (cm *ComboManager) ResetCombo(playerID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if combo, ok := cm.Combos[playerID]; ok {
		combo.Count = 0
		combo.Multiplier = 1.0
		combo.IsActive = false
		combo.Events = combo.Events[:0]
	}
}

// 获取排行榜
func (cm *ComboManager) GetLeaderboard(limit int) []*Combo {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	combos := make([]*Combo, 0, len(cm.Combos))
	for _, combo := range cm.Combos {
		combos = append(combos, combo)
	}
	
	// 排序
	for i := 0; i < len(combos)-1; i++ {
		for j := i + 1; j < len(combos); j++ {
			if combos[j].Score > combos[i].Score {
				combos[i], combos[j] = combos[j], combos[i]
			}
		}
	}
	
	if limit > 0 && limit < len(combos) {
		combos = combos[:limit]
	}
	
	return combos
}

// 获取统计信息
func (cm *ComboManager) GetStats() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	totalScore := int64(0)
	maxCombo := 0
	maxMultiplier := 1.0
	activeCount := 0
	
	for _, combo := range cm.Combos {
		totalScore += combo.Score
		if combo.MaxCount > maxCombo {
			maxCombo = combo.MaxCount
		}
		if combo.MaxMultiplier > maxMultiplier {
			maxMultiplier = combo.MaxMultiplier
		}
		if combo.IsActive {
			activeCount++
		}
	}
	
	return map[string]interface{}{
		"totalPlayers":   len(cm.Combos),
		"activeCombos":   activeCount,
		"totalScore":     totalScore,
		"maxCombo":       maxCombo,
		"maxMultiplier": maxMultiplier,
	}
}

// ============================================
// 分数系统 (Score System)
// ============================================

// ScoreType 分数类型
type ScoreType int

const (
	ScoreKill ScoreType = iota
	ScoreWave
	ScoreTime
	ScoreCombo
	ScoreCollection
	ScoreAchievement
	ScoreMVP
)

// ScoreEntry 分数记录
type ScoreEntry struct {
	Type      ScoreType   `json:"type"`
	Points    int64       `json:"points"`
	Timestamp int64       `json:"timestamp"`
	Details   string      `json:"details"`
}

// ScoreManager 分数管理器
type ScoreManager struct {
	Entries  map[string][]ScoreEntry // playerID -> entries
	TotalScore map[string]int64
	mu       sync.RWMutex
}

func NewScoreManager() *ScoreManager {
	return &ScoreManager{
		Entries:    make(map[string][]ScoreEntry),
		TotalScore: make(map[string]int64),
	}
}

// 添加分数
func (sm *ScoreManager) AddScore(playerID string, scoreType ScoreType, points int64, details string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	entry := ScoreEntry{
		Type:      scoreType,
		Points:    points,
		Timestamp: time.Now().UnixMilli(),
		Details:   details,
	}
	
	sm.Entries[playerID] = append(sm.Entries[playerID], entry)
	sm.TotalScore[playerID] += points
	
	// 保持记录合理大小
	if len(sm.Entries[playerID]) > 10000 {
		sm.Entries[playerID] = sm.Entries[playerID][len(sm.Entries[playerID])-5000:]
	}
}

// 获取玩家总分
func (sm *ScoreManager) GetTotalScore(playerID string) int64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.TotalScore[playerID]
}

// 获取排行榜
func (sm *ScoreManager) GetLeaderboard(limit int) []struct {
	PlayerID string
	Score    int64
} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	type entry struct {
		playerID string
		score    int64
	}
	
	leaderboard := make([]entry, 0, len(sm.TotalScore))
	for pid, score := range sm.TotalScore {
		leaderboard = append(leaderboard, entry{pid, score})
	}
	
	// 排序
	for i := 0; i < len(leaderboard)-1; i++ {
		for j := i + 1; j < len(leaderboard); j++ {
			if leaderboard[j].score > leaderboard[i].score {
				leaderboard[i], leaderboard[j] = leaderboard[j], leaderboard[i]
			}
		}
	}
	
	if limit > 0 && limit < len(leaderboard) {
		leaderboard = leaderboard[:limit]
	}
	
	result := make([]struct {
		PlayerID string
		Score    int64
	}, len(leaderboard))
	
	for i, e := range leaderboard {
		result[i] = struct {
			PlayerID string
			Score    int64
		}{e.playerID, e.score}
	}
	
	return result
}

// 重置分数
func (sm *ScoreManager) Reset(playerID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	delete(sm.Entries, playerID)
	delete(sm.TotalScore, playerID)
}
