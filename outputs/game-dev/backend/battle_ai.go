package main

import (
	"math"
	"math/rand"
	"time"
)

// ============================================
// AI敌人行为系统
// ============================================

type AIBehavior struct {
	Type         string  `json:"type"`          // 行为类型
	TriggerHP    float64 `json:"triggerHP"`      // 触发血量百分比 (0-1)
	TriggerDist  float64 `json:"triggerDist"`    // 触发距离
	ActionChance float64 `json:"actionChance"`   // 行为执行概率
	Cooldown     int64   `json:"cooldown"`       // 冷却时间 (ms)
	LastAction   int64   `json:"lastAction"`     // 上次执行时间
}

// AI敌人行为类型
const (
	AIBehaviorRetreat  = "retreat"  // 撤退: 低血量时撤退回血
	AIBehaviorCharge   = "charge"   // 冲锋: 快速冲向终点
	AIBehaviorSplit    = "split"    // 分裂: 分成多个小怪
	AIBehaviorShield   = "shield"   // 护盾: 获得临时护盾
	AIBehaviorSpawn    = "spawn"    // 产卵: 产卵小怪
	AIBehaviorTeleport = "teleport" // 瞬移: 随机位置
	AIBehaviorSteal    = "steal"    // 偷取: 偷取金币
	AIBehaviorBurrow   = "burrow"   // 钻地: 暂时无敌
)

// 高级敌人 (带AI)
type AdvancedEnemy struct {
	*Enemy
	Behaviors []AIBehavior
	ShieldHP float64 // 护盾血量
	IsBurrowed bool  // 是否钻地
}

// 创建高级敌人
func NewAdvancedEnemy(id, enemyType string) *AdvancedEnemy {
	enemy := NewEnemySpawner().createEnemy(enemyType)
	
	behaviors := []AIBehavior{
		{Type: AIBehaviorRetreat, TriggerHP: 0.3, ActionChance: 0.5, Cooldown: 10000},
		{Type: AIBehaviorCharge, TriggerHP: 0.8, ActionChance: 0.3, Cooldown: 15000},
	}

	if enemyType == "tank" {
		behaviors = append(behaviors, AIBehavior{Type: AIBehaviorShield, TriggerHP: 0.5, ActionChance: 0.4, Cooldown: 20000})
	}

	return &AdvancedEnemy{
		Enemy:     enemy,
		Behaviors: behaviors,
		ShieldHP:  0,
		IsBurrowed: false,
	}
}

// AI行为更新
func (ae *AdvancedEnemy) UpdateAI(dt float64, battle *BattleManager) {
	now := time.Now().UnixMilli()
	
	// 检查是否在冷却中
	for i := range ae.Behaviors {
		behavior := &ae.Behaviors[i]
		if now-behavior.LastAction < behavior.Cooldown {
			continue
		}

		// 随机概率检查
		if rand.Float64() > behavior.ActionChance {
			continue
		}

		// 血量触发
		hpPercent := ae.HP / ae.MaxHP
		if hpPercent > behavior.TriggerHP {
			continue
		}

		// 执行行为
		ae.executeBehavior(behavior, battle)
		behavior.LastAction = now
	}
}

// 执行AI行为
func (ae *AdvancedEnemy) executeBehavior(behavior *AIBehavior, battle *BattleManager) {
	switch behavior.Type {
	case AIBehaviorRetreat:
		// 撤退: 快速回血 (每秒回复10%)
		go func() {
			for i := 0; i < 5; i++ {
				ae.HP = math.Min(ae.HP+ae.MaxHP*0.1, ae.MaxHP)
				time.Sleep(time.Second)
			}
		}()

	case AIBehaviorCharge:
		// 冲锋: 速度翻倍
		ae.Speed *= 2
		time.AfterFunc(3*time.Second, func() {
			ae.Speed /= 2
		})

	case AIBehaviorShield:
		// 护盾: 获得护盾
		ae.ShieldHP = ae.MaxHP * 0.5

	case AIBehaviorSpawn:
		// 产卵: 生成2个小怪
		for i := 0; i < 2; i++ {
			child := battle.Spawner.createEnemy("grunt")
			child.X = ae.X + float64(i*20)
			child.Y = ae.Y
			battle.Spawner.Enemies = append(battle.Spawner.Enemies, child)
		}

	case AIBehaviorTeleport:
		// 瞬移: 随机位置
		ae.X = rand.Float64() * 800
		ae.Y = rand.Float64() * 600

	case AIBehaviorSteal:
		// 偷取: 减少玩家金币
		battle.State.Money = math.Max(0, battle.State.Money-50)

	case AIBehaviorBurrow:
		// 钻地: 暂时无敌
		ae.IsBurrowed = true
		time.AfterFunc(2*time.Second, func() {
			ae.IsBurrowed = false
		})
	}
}

// 高级敌人受到伤害
func (ae *AdvancedEnemy) TakeDamage(damage, armorReduce float64) {
	// 钻地无敌
	if ae.IsBurrowed {
		return
	}

	// 先扣护盾
	if ae.ShieldHP > 0 {
		damageLeft := damage
		if ae.ShieldHP >= damage {
			ae.ShieldHP -= damage
			return
		}
		damageLeft -= ae.ShieldHP
		ae.ShieldHP = 0
		damage = damageLeft
	}

	// 计算护甲减伤
	armor := ae.Armor * (1 - armorReduce)
	actualDamage := damage * (1 - armor/(armor+100))
	ae.HP -= actualDamage
}

// ============================================
// 路径寻路系统
// ============================================

type PathNode struct {
	X, Y     float64
	Cost     float64
	Parent   *PathNode
	IsBlocked bool
}

// A*寻路
type PathFinder struct {
	GridWidth  int
	GridHeight int
	GridSize   float64
}

// 新建寻路器
func NewPathFinder(width, height int, gridSize float64) *PathFinder {
	return &PathFinder{
		GridWidth:  width,
		GridHeight: height,
		GridSize:   gridSize,
	}
}

// 寻路 (A*算法简化版)
func (pf *PathFinder) FindPath(startX, startY, endX, endY float64, obstacles [][]bool) []Point {
	// 简化实现: 返回直线+小幅度随机偏移
	path := []Point{
		{X: startX, Y: startY},
	}

	// 计算方向
	dx := endX - startX
	dy := endY - startY
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		return path
	}

	// 添加中间点
	steps := int(dist / 50)
	for i := 1; i <= steps; i++ {
		progress := float64(i) / float64(steps)
		x := startX + dx*progress + (rand.Float64()-0.5)*20
		y := startY + dy*progress + (rand.Float64()-0.5)*20
		path = append(path, Point{X: x, Y: y})
	}

	path = append(path, Point{X: endX, Y: endY})
	return path
}

// 动态路径更新
type DynamicPath struct {
	BasePath   []Point
	Detours    []Detour
	CurrentIdx int
}

type Detour struct {
	StartTime time.Time
	Duration  time.Duration
	From      Point
	To        Point
}

// 获取实际路径
func (dp *DynamicPath) GetActualPath() []Point {
	if len(dp.Detours) == 0 {
		return dp.BasePath
	}

	now := time.Now()
	path := make([]Point, 0, len(dp.BasePath))

	for i, point := range dp.BasePath {
		modified := point
		
		// 检查是否有绕行
		for _, detour := range dp.Detours {
			if now.After(detour.StartTime) && now.Before(detour.StartTime.Add(detour.Duration)) {
				// 简单处理: 靠近绕行起点时使用绕行路径
				dist := math.Sqrt(math.Pow(point.X-detour.From.X, 2) + math.Pow(point.Y-detour.From.Y, 2))
				if dist < 50 {
					modified = detour.To
				}
			}
		}
		
		path = append(path, modified)
	}

	return path
}

// ============================================
// 敌人类型工厂
// ============================================

type EnemyFactory struct {
	Templates map[string]EnemyTemplate
}

type EnemyTemplate struct {
	Type         EnemyType
	BaseHP       float64
	Speed        float64
	Armor        float64
	Reward       int
	Behaviors    []string
	Sprite       string
	Animations   []string
}

// 敌人工厂
func NewEnemyFactory() *EnemyFactory {
	return &EnemyFactory{
		Templates: map[string]EnemyTemplate{
			"grunt": {
				Type: EnemyTypeGrunt, BaseHP: 50, Speed: 50, Armor: 0, Reward: 10,
				Sprite: "grunt", Animations: []string{"walk", "die"},
			},
			"ranger": {
				Type: EnemyTypeRanger, BaseHP: 30, Speed: 40, Armor: 5, Reward: 15,
				Behaviors: []string{"attack"}, Sprite: "ranger", Animations: []string{"walk", "attack", "die"},
			},
			"tank": {
				Type: EnemyTypeTank, BaseHP: 200, Speed: 25, Armor: 20, Reward: 30,
				Behaviors: []string{"shield"}, Sprite: "tank", Animations: []string{"walk", "shield", "die"},
			},
			"boss_dragon": {
				Type: EnemyTypeBoss, BaseHP: 2000, Speed: 25, Armor: 30, Reward: 500,
				Behaviors: []string{"rage", "aoe", "summon"}, Sprite: "dragon", Animations: []string{"walk", "rage", "aoe", "die"},
			},
			"boss_golem": {
				Type: EnemyTypeBoss, BaseHP: 3000, Speed: 15, Armor: 50, Reward: 800,
				Behaviors: []string{"shield", "heal", "aoe"}, Sprite: "golem", Animations: []string{"walk", "shield", "heal", "die"},
			},
			"boss_demon": {
				Type: EnemyTypeBoss, BaseHP: 1500, Speed: 35, Armor: 20, Reward: 600,
				Behaviors: []string{"rage", "execute", "slow"}, Sprite: "demon", Animations: []string{"walk", "rage", "execute", "die"},
			},
		},
	}
}

// 根据模板创建敌人
func (ef *EnemyFactory) Create(templateName string) *Enemy {
	template, ok := ef.Templates[templateName]
	if !ok {
		return nil
	}

	enemy := &Enemy{
		ID:         generateID(),
		Type:       template.Type,
		X:          0,
		Y:          300,
		Progress:   0,
		HP:         template.BaseHP,
		MaxHP:      template.BaseHP,
		Speed:      template.Speed,
		Armor:      template.Armor,
		Reward:     template.Reward,
		PathIndex:  0,
		PathOffset: 0,
		FreezeTime: 0,
		SlowFactor: 1.0,
	}

	return enemy
}

// ============================================
// 塔类型工厂
// ============================================

type TowerFactory struct {
	Templates map[string]TowerTemplate
}

type TowerTemplate struct {
	Type       TowerType
	BaseDamage float64
	BaseRange  float64
	BaseFireRate float64
	LevelBonus float64
	Special    string // 特殊效果
	Sprite     string
}

// 塔工厂
func NewTowerFactory() *TowerFactory {
	return &TowerFactory{
		Templates: map[string]TowerTemplate{
			"arrow": {
				Type: TowerTypeAttack, BaseDamage: 10, BaseRange: 150, BaseFireRate: 1.0, LevelBonus: 2,
				Sprite: "arrow_tower",
			},
			"cannon": {
				Type: TowerTypeAttack, BaseDamage: 30, BaseRange: 120, BaseFireRate: 0.5, LevelBonus: 5,
				Sprite: "cannon_tower",
			},
			"ice": {
				Type: TowerTypeAttack, BaseDamage: 8, BaseRange: 130, BaseFireRate: 1.2, LevelBonus: 1.5,
				Special: "slow", Sprite: "ice_tower",
			},
			"lightning": {
				Type: TowerTypeAttack, BaseDamage: 25, BaseRange: 140, BaseFireRate: 0.8, LevelBonus: 4,
				Special: "chain", Sprite: "lightning_tower",
			},
			"tower_heal": {
				Type: TowerTypeSupport, BaseDamage: -20, BaseRange: 100, BaseFireRate: 1.0, LevelBonus: -5,
				Special: "heal", Sprite: "heal_tower",
			},
			"wall": {
				Type: TowerTypeDefense, BaseDamage: 0, BaseRange: 0, BaseFireRate: 0, LevelBonus: 0,
				Special: "block", Sprite: "wall_tower",
			},
			"tower_buff": {
				Type: TowerTypeSupport, BaseDamage: 0, BaseRange: 80, BaseFireRate: 0, LevelBonus: 0,
				Special: "buff", Sprite: "buff_tower",
			},
		},
	}
}

// 根据模板创建塔
func (tf *TowerFactory) Create(templateName string, x, y float64, level int) *Tower {
	template, ok := tf.Templates[templateName]
	if !ok {
		return nil
	}

	tower := &Tower{
		ID:         generateID(),
		Type:       template.Type,
		X:          x,
		Y:          y,
		Level:      level,
		Range:      template.BaseRange * (1 + float64(level)*0.05),
		Damage:     template.BaseDamage + float64(level-1)*template.LevelBonus,
		FireRate:   template.BaseFireRate,
		LastFire:   0,
	}

	return tower
}

// ============================================
// 关卡编辑器数据结构
// ============================================

type Level struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Difficulty  int         `json:"difficulty"`  // 1-10
	Waves       []WaveConfig `json:"waves"`
	Path        []Point     `json:"path"`
	TowerSpots  []Point     `json:"towerSpots"` // 允许建塔的位置
	Background  string      `json:"background"`
	Music       string      `json:"music"`
}

type WaveConfig struct {
	WaveNum     int         `json:"waveNum"`
	Duration    int         `json:"duration"`   // 持续时间(秒)
	Enemies     []SpawnInfo `json:"enemies"`
	BossWave    bool        `json:"bossWave"`   // 是否Boss波次
	BossType    string      `json:"bossType"`    // Boss类型
}

// 预设关卡
var LevelConfigs = []Level{
	{
		ID: "level_1", Name: "新手村", Difficulty: 1,
		Waves: []WaveConfig{
			{WaveNum: 1, Duration: 30, Enemies: []SpawnInfo{{Type: "grunt", Count: 5, Interval: 3000}}},
			{WaveNum: 2, Duration: 30, Enemies: []SpawnInfo{{Type: "grunt", Count: 8, Interval: 2500}}},
		},
		Path: []Point{{X: 0, Y: 300}, {X: 200, Y: 300}, {X: 400, Y: 300}, {X: 600, Y: 300}, {X: 800, Y: 300}},
		TowerSpots: []Point{{X: 100, Y: 250}, {X: 300, Y: 350}, {X: 500, Y: 250}, {X: 700, Y: 350}},
		Background: "grass", Music: "level1",
	},
	{
		ID: "level_2", Name: "森林深处", Difficulty: 3,
		Waves: []WaveConfig{
			{WaveNum: 1, Duration: 30, Enemies: []SpawnInfo{{Type: "grunt", Count: 10, Interval: 2000}}},
			{WaveNum: 2, Duration: 40, Enemies: []SpawnInfo{{Type: "grunt", Count: 8, Interval: 1500}, {Type: "ranger", Count: 3, Interval: 4000}}},
			{WaveNum: 3, Duration: 45, Enemies: []SpawnInfo{{Type: "grunt", Count: 12, Interval: 1200}, {Type: "ranger", Count: 5, Interval: 3000}, {Type: "tank", Count: 2, Interval: 8000}}},
		},
		Path: []Point{{X: 0, Y: 100}, {X: 200, Y: 100}, {X: 200, Y: 400}, {X: 500, Y: 400}, {X: 500, Y: 200}, {X: 800, Y: 200}},
		TowerSpots: []Point{{X: 100, Y: 150}, {X: 250, Y: 350}, {X: 450, Y: 450}, {X: 550, Y: 150}, {X: 700, Y: 250}},
		Background: "forest", Music: "level2",
	},
	{
		ID: "level_3", Name: "Boss领地", Difficulty: 5,
		Waves: []WaveConfig{
			{WaveNum: 1, Duration: 30, Enemies: []SpawnInfo{{Type: "grunt", Count: 15, Interval: 1500}}},
			{WaveNum: 2, Duration: 45, Enemies: []SpawnInfo{{Type: "grunt", Count: 10, Interval: 1000}, {Type: "ranger", Count: 8, Interval: 2500}, {Type: "tank", Count: 3, Interval: 6000}}},
			{WaveNum: 3, Duration: 90, BossWave: true, BossType: "dragon", Enemies: []SpawnInfo{{Type: "grunt", Count: 20, Interval: 800}, {Type: "ranger", Count: 10, Interval: 1500}, {Type: "tank", Count: 5, Interval: 4000}}},
		},
		Path: []Point{{X: 0, Y: 300}, {X: 150, Y: 300}, {X: 150, Y: 150}, {X: 350, Y: 150}, {X: 350, Y: 450}, {X: 550, Y: 450}, {X: 550, Y: 200}, {X: 800, Y: 200}},
		TowerSpots: []Point{{X: 80, Y: 250}, {X: 200, Y: 200}, {X: 300, Y: 100}, {X: 400, Y: 400}, {X: 500, Y: 500}, {X: 600, Y: 300}, {X: 700, Y: 150}},
		Background: "boss_lair", Music: "boss",
	},
}

// 加载关卡
func LoadLevel(levelID string) *Level {
	for _, level := range LevelConfigs {
		if level.ID == levelID {
			return &level
		}
	}
	return nil
}

// ============================================
// 战斗引擎接口
// ============================================

// 战斗引擎接口
type BattleEngine interface {
	Init()
	Start()
	Pause()
	Resume()
	End()
	Update(dt float64)
	GetState() *BattleState
	GetResult() *BattleResult
}

// 战斗引擎实现
type StandardBattleEngine struct {
	Manager     *BattleManagerEx
	Level       *Level
	CurrentWave int
	IsPaused    bool
	mu          sync.RWMutex
}

// 新建战斗引擎
func NewBattleEngine(levelID string) *StandardBattleEngine {
	level := LoadLevel(levelID)
	if level == nil {
		return nil
	}

	return &StandardBattleEngine{
		Manager:     NewBattleManagerEx(),
		Level:       level,
		CurrentWave: 0,
		IsPaused:    false,
	}
}

// 初始化
func (e *StandardBattleEngine) Init() {
	e.Manager.State = BattleState{
		Wave:      1,
		Score:     0,
		Lives:     20,
		Money:     100,
		WaveTime:  30,
		IsRunning: false,
	}
}

// 开始战斗
func (e *StandardBattleEngine) Start() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Manager.State.IsRunning = true
	e.Manager.StartTime = time.Now()
	e.CurrentWave = 1
	
	// 加载第一波
	e.loadWave(e.CurrentWave)
}

// 暂停
func (e *StandardBattleEngine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.IsPaused = true
	e.Manager.State.IsRunning = false
}

// 恢复
func (e *StandardBattleEngine) Resume() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.IsPaused = false
	e.Manager.State.IsRunning = true
}

// 结束
func (e *StandardBattleEngine) End() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Manager.State.IsRunning = false
}

// 更新
func (e *StandardBattleEngine) Update(dt float64) {
	if e.IsPaused || !e.Manager.State.IsRunning {
		return
	}

	e.Manager.Update(dt)

	// 检查是否需要开始下一波
	if len(e.Manager.EnemySpawner.Enemies) == 0 && e.Manager.State.Lives > 0 {
		e.NextWave()
	}
}

// 下一波
func (e *StandardBattleEngine) NextWave() {
	e.CurrentWave++
	if e.CurrentWave > len(e.Level.Waves) {
		// 战斗胜利
		e.End()
		return
	}
	e.loadWave(e.CurrentWave)
}

// 加载波次
func (e *StandardBattleEngine) loadWave(waveNum int) {
	waveIndex := waveNum - 1
	if waveIndex >= len(e.Level.Waves) {
		return
	}

	wave := e.Level.Waves[waveIndex]
	e.Manager.EnemySpawner.Wave = waveNum
	e.Manager.EnemySpawner.SpawnQueue = wave.Enemies
	e.Manager.State.Wave = waveNum
	e.Manager.State.WaveTime = wave.Duration
}

// 获取状态
func (e *StandardBattleEngine) GetState() *BattleState {
	return &e.Manager.State
}

// 获取结果
func (e *StandardBattleEngine) GetResult() *BattleResult {
	return e.Manager.EndBattle()
}
