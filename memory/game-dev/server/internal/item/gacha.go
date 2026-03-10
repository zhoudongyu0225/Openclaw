package item

import (
	"math/rand"
	"sync"
	"time"
)

// GachaPool 抽卡池
type GachaPool struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Items       []*GachaItem `json:"items"`
	RateUpItems []string   `json:"rate_up_items"` // 保底提升的道具ID
	CreatedAt   time.Time `json:"created_at"`
}

// GachaItem 抽卡道具
type GachaItem struct {
	ItemID     string  `json:"item_id"`
	Weight     int     `json:"weight"`      // 权重
	IsRateUp   bool    `json:"is_rate_up"`  // 是否为保底提升
	MinCount   int     `json:"min_count"`   // 最小数量
	MaxCount   int     `json:"max_count"`   // 最大数量
}

// GachaResult 抽卡结果
type GachaResult struct {
	Items      map[string]int `json:"items"`       // 获得的道具
	TotalCoins int            `json:"total_coins"` // 消耗的金币
	IsFirst    bool           `json:"is_first"`    // 是否为首次抽卡
	Timestamp  time.Time      `json:"timestamp"`
}

// GachaType 抽卡类型
type GachaType int

const (
	GachaTypeSingle GachaType = iota + 1
	GachaTypeMulti
	GachaTypeGuaranteed
)

// GachaManager 抽卡管理器
type GachaManager struct {
	pools      map[string]*GachaPool
	playerGacha map[string]*PlayerGachaRecord
	mu          sync.RWMutex
}

// PlayerGachaRecord 玩家抽卡记录
type PlayerGachaRecord struct {
	PlayerID      string            `json:"player_id"`
	PoolID        string            `json:"pool_id"`
	TotalCount    int               `json:"total_count"`     // 总抽卡次数
	NoSSRCount    int               `json:"no_ssr_count"`    // 连续未抽到SSR次数
	LastGachaTime time.Time         `json:"last_gacha_time"` // 上次抽卡时间
	History       []map[string]int  `json:"history"`          // 抽卡历史
}

// NewGachaManager 创建抽卡管理器
func NewGachaManager() *GachaManager {
	m := &GachaManager{
		pools:      make(map[string]*GachaPool),
		playerGacha: make(map[string]*PlayerGachaRecord),
	}
	m.initDefaultPools()
	return m
}

// initDefaultPools 初始化默认抽卡池
func (m *GachaManager) initDefaultPools() {
	// 武器池
	weaponPool := &GachaPool{
		ID:          "pool_weapon",
		Name:        "武器宝库",
		Description: "稀有能力武器等你来拿",
		RateUpItems: []string{"weapon_004", "weapon_005"},
		Items: []*GachaItem{
			{ItemID: "weapon_001", Weight: 5000, MinCount: 1, MaxCount: 1},
			{ItemID: "weapon_002", Weight: 3000, MinCount: 1, MaxCount: 1},
			{ItemID: "weapon_003", Weight: 1500, MinCount: 1, MaxCount: 1},
			{ItemID: "weapon_004", Weight: 400, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "weapon_005", Weight: 100, IsRateUp: true, MinCount: 1, MaxCount: 1},
		},
		CreatedAt: time.Now(),
	}
	m.pools[weaponPool.ID] = weaponPool

	// 护甲池
	armorPool := &GachaPool{
		ID:          "pool_armor",
		Name:        "防具宝库",
		Description: "强力护甲限时up",
		RateUpItems: []string{"armor_004", "armor_005"},
		Items: []*GachaItem{
			{ItemID: "armor_001", Weight: 5000, MinCount: 1, MaxCount: 1},
			{ItemID: "armor_002", Weight: 3000, MinCount: 1, MaxCount: 1},
			{ItemID: "armor_003", Weight: 1500, MinCount: 1, MaxCount: 1},
			{ItemID: "armor_004", Weight: 400, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "armor_005", Weight: 100, IsRateUp: true, MinCount: 1, MaxCount: 1},
		},
		CreatedAt: time.Now(),
	}
	m.pools[armorPool.ID] = armorPool

	// 综合池
	allPool := &GachaPool{
		ID:          "pool_all",
		Name:        "综合宝库",
		Description: "所有道具概率up",
		RateUpItems: []string{"weapon_004", "weapon_005", "armor_004", "armor_005", "acc_004", "acc_005"},
		Items: []*GachaItem{
			// 武器
			{ItemID: "weapon_001", Weight: 2000, MinCount: 1, MaxCount: 2},
			{ItemID: "weapon_002", Weight: 1500, MinCount: 1, MaxCount: 2},
			{ItemID: "weapon_003", Weight: 800, MinCount: 1, MaxCount: 2},
			{ItemID: "weapon_004", Weight: 200, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "weapon_005", Weight: 50, IsRateUp: true, MinCount: 1, MaxCount: 1},
			// 护甲
			{ItemID: "armor_001", Weight: 2000, MinCount: 1, MaxCount: 2},
			{ItemID: "armor_002", Weight: 1500, MinCount: 1, MaxCount: 2},
			{ItemID: "armor_003", Weight: 800, MinCount: 1, MaxCount: 2},
			{ItemID: "armor_004", Weight: 200, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "armor_005", Weight: 50, IsRateUp: true, MinCount: 1, MaxCount: 1},
			// 饰品
			{ItemID: "acc_001", Weight: 3000, MinCount: 1, MaxCount: 3},
			{ItemID: "acc_002", Weight: 2000, MinCount: 1, MaxCount: 2},
			{ItemID: "acc_003", Weight: 1000, MinCount: 1, MaxCount: 2},
			{ItemID: "acc_004", Weight: 150, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "acc_005", Weight: 50, IsRateUp: true, MinCount: 1, MaxCount: 1},
			// 材料
			{ItemID: "material_001", Weight: 5000, MinCount: 5, MaxCount: 10},
			{ItemID: "material_002", Weight: 2000, MinCount: 2, MaxCount: 5},
			{ItemID: "material_003", Weight: 500, MinCount: 1, MaxCount: 2},
		},
		CreatedAt: time.Now(),
	}
	m.pools[allPool.ID] = allPool

	// 限时概率up池
	rateUpPool := &GachaPool{
		ID:          "pool_rateup",
		Name:        "概率up宝库",
		Description: "限定SSR概率大幅提升",
		RateUpItems: []string{"weapon_005", "armor_005", "acc_005"},
		Items: []*GachaItem{
			{ItemID: "weapon_001", Weight: 1500, MinCount: 1, MaxCount: 2},
			{ItemID: "weapon_002", Weight: 1000, MinCount: 1, MaxCount: 2},
			{ItemID: "weapon_003", Weight: 500, MinCount: 1, MaxCount: 2},
			{ItemID: "weapon_004", Weight: 800, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "weapon_005", Weight: 300, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "armor_001", Weight: 1500, MinCount: 1, MaxCount: 2},
			{ItemID: "armor_002", Weight: 1000, MinCount: 1, MaxCount: 2},
			{ItemID: "armor_003", Weight: 500, MinCount: 1, MaxCount: 2},
			{ItemID: "armor_004", Weight: 800, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "armor_005", Weight: 300, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "acc_001", Weight: 2000, MinCount: 1, MaxCount: 3},
			{ItemID: "acc_002", Weight: 1500, MinCount: 1, MaxCount: 2},
			{ItemID: "acc_003", Weight: 800, MinCount: 1, MaxCount: 2},
			{ItemID: "acc_004", Weight: 600, IsRateUp: true, MinCount: 1, MaxCount: 1},
			{ItemID: "acc_005", Weight: 200, IsRateUp: true, MinCount: 1, MaxCount: 1},
		},
		CreatedAt: time.Now(),
	}
	m.pools[rateUpPool.ID] = rateUpPool
}

// GetPool 获取抽卡池
func (m *GachaManager) GetPool(poolID string) (*GachaPool, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pool, ok := m.pools[poolID]
	return pool, ok
}

// GetAllPools 获取所有抽卡池
func (m *GachaManager) GetAllPools() []*GachaPool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pools := make([]*GachaPool, 0, len(m.pools))
	for _, pool := range m.pools {
		pools = append(pools, pool)
	}

	return pools
}

// Gacha 抽卡
func (m *GachaManager) Gacha(playerID, poolID string, gachaType GachaType) (*GachaResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, ok := m.pools[poolID]
	if !ok {
		return nil, ErrInvalidPool
	}

	// 获取玩家抽卡记录
	recordKey := playerID + "_" + poolID
	record, exists := m.playerGacha[recordKey]
	if !exists {
		record = &PlayerGachaRecord{
			PlayerID: playerID,
			PoolID:   poolID,
			History:  make([]map[string]int, 0),
		}
		m.playerGacha[recordKey] = record
	}

	// 计算抽卡次数
	count := 1
	totalCoins := 100 // 单抽价格
	if gachaType == GachaTypeMulti {
		count = 10
		totalCoins = 1000 // 十连价格
	}

	// 开始抽卡
	result := &GachaResult{
		Items:      make(map[string]int),
		TotalCoins: totalCoins,
		Timestamp:  time.Now(),
	}

	// 计算权重
	totalWeight := 0
	for _, item := range pool.Items {
		weight := item.Weight
		// 保底提升
		if item.IsRateUp && record.NoSSRCount >= 50 {
			weight *= 2 // 保底时概率翻倍
		}
		totalWeight += weight
	}

	// 执行抽卡
	for i := 0; i < count; i++ {
		randSeed := rand.Intn(totalWeight)
		currentWeight := 0

		var selectedItem *GachaItem
		for _, item := range pool.Items {
			weight := item.Weight
			if item.IsRateUp && record.NoSSRCount >= 50 {
				weight *= 2
			}
			currentWeight += weight

			if randSeed < currentWeight {
				selectedItem = item
				break
			}
		}

		if selectedItem != nil {
			// 随机数量
			itemCount := selectedItem.MinCount
			if selectedItem.MaxCount > selectedItem.MinCount {
				itemCount = selectedItem.MinCount + rand.Intn(selectedItem.MaxCount-selectedItem.MinCount+1)
			}

			result.Items[selectedItem.ItemID] += itemCount

			// 检查是否为SSR
			if selectedItem.ItemID == "weapon_005" || selectedItem.ItemID == "armor_005" || selectedItem.ItemID == "acc_005" {
				record.NoSSRCount = 0
			} else {
				record.NoSSRCount++
			}
		}
	}

	record.TotalCount += count
	record.LastGachaTime = time.Now()
	record.History = append(record.History, result.Items)

	return result, nil
}

// GetPlayerRecord 获取玩家抽卡记录
func (m *GachaManager) GetPlayerRecord(playerID, poolID string) *PlayerGachaRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()

	recordKey := playerID + "_" + poolID
	if record, ok := m.playerGacha[recordKey]; ok {
		return record
	}

	return &PlayerGachaRecord{
		PlayerID: playerID,
		PoolID:   poolID,
		History:  make([]map[string]int, 0),
	}
}

// CalculateSSRChance 计算SSR概率
func (m *GachaManager) CalculateSSRChance(poolID string, totalGacha int) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pool, ok := m.pools[poolID]
	if !ok {
		return 0
	}

	ssrWeight := 0
	totalWeight := 0
	for _, item := range pool.Items {
		if item.ItemID == "weapon_005" || item.ItemID == "armor_005" || item.ItemID == "acc_005" {
			ssrWeight += item.Weight
		}
		totalWeight += item.Weight
	}

	baseChance := float64(ssrWeight) / float64(totalWeight) * 100

	// 保底概率提升
	guaranteedChance := float64(totalGacha) * 0.5 // 每次抽卡增加0.5%的基础概率
	if guaranteedChance > 50 {
		guaranteedChance = 50
	}

	return baseChance + guaranteedChance
}

// AddPool 添加抽卡池
func (m *GachaManager) AddPool(pool *GachaPool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.pools[pool.ID] = pool
}

// GachaError 抽卡错误
type GachaError string

func (e GachaError) Error() string {
	return string(e)
}

const (
	ErrInvalidPool    GachaError = "invalid pool ID"
	ErrInsufficientBalance GachaError = "insufficient balance"
	ErrPoolNotFound   GachaError = "pool not found"
)
