package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// ============================================
// 防御塔系统
// ============================================

type TowerType int

const (
	TowerTypeAttack TowerType = iota // 攻击塔
	TowerTypeDefense                 // 防御塔
	TowerTypeSupport                  // 辅助塔
)

// 防御塔
type Tower struct {
	ID       string    `json:"id"`
	Type     TowerType `json:"type"`
	X        float64   `json:"x"`
	Y        float64   `json:"y"`
	Level    int       `json:"level"`
	Range    float64   `json:"range"`    // 攻击范围
	Damage   float64   `json:"damage"`   // 伤害值
	FireRate float64   `json:"fireRate"` // 攻击频率 (次/秒)
	LastFire int64      `json:"lastFire"` // 上次攻击时间 (ms)

	// 特有属性
	ArmorReduce float64 // 护甲削减 (防御塔)
	BuffTarget  string  // buff目标 (辅助塔)
	BuffValue   float64 // buff数值

	mu sync.RWMutex
}

// 防御塔配置
var TowerConfigs = map[string]struct {
	BaseDamage  float64
	BaseRange   float64
	BaseFireRate float64
	LevelBonus  float64
}{
	"arrow":    {Damage: 10, Range: 150, FireRate: 1.0, LevelBonus: 2},
	"cannon":   {Damage: 30, Range: 120, FireRate: 0.5, LevelBonus: 5},
	"ice":      {Damage: 8, Range: 130, FireRate: 1.2, LevelBonus: 1.5, ArmorReduce: 0.3},
	"lightning":{Damage: 25, Range: 140, FireRate: 0.8, LevelBonus: 4},
	"tower_heal":{Damage: -20, Range: 100, FireRate: 1.0, LevelBonus: -5}, // 治疗塔
}

// 新建防御塔
func NewTower(id string, towerType TowerType, x, y float64, level int) *Tower {
	config := TowerConfigs[id]
	tower := &Tower{
		ID:         id,
		Type:       towerType,
		X:          x,
		Y:          y,
		Level:      level,
		Range:      config.BaseRange,
		Damage:     config.BaseDamage + float64(level-1)*config.LevelBonus,
		FireRate:   config.BaseFireRate,
		LastFire:   0,
		ArmorReduce: config.ArmorReduce,
	}
	return tower
}

// 获取攻击冷却时间 (毫秒)
func (t *Tower) GetCooldown() int64 {
	return int64(1000.0 / t.FireRate)
}

// 是否可以攻击
func (t *Tower) CanFire() bool {
	now := time.Now().UnixMilli()
	return now-t.LastFire >= t.GetCooldown()
}

// 寻找目标 - 优先攻击最接近终点的敌人
func (t *Tower) FindTarget(enemies []*Enemy) *Enemy {
	var target *Enemy
	var maxProgress float64 = -1

	for _, e := range enemies {
		if e.IsDead() {
			continue
		}
		dist := t.DistanceTo(e)
		if dist <= t.Range {
			if e.Progress > maxProgress {
				maxProgress = e.Progress
				target = e
			}
		}
	}
	return target
}

// 计算到敌人的距离
func (t *Tower) DistanceTo(e *Enemy) float64 {
	dx := t.X - e.X
	dy := t.Y - e.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// 攻击敌人
func (t *Tower) Attack(enemies []*Enemy) *Projectile {
	if !t.CanFire() {
		return nil
	}

	target := t.FindTarget(enemies)
	if target == nil {
		return nil
	}

	t.LastFire = time.Now().UnixMilli()

	// 创建投射物
	proj := &Projectile{
		ID:        generateID(),
		TowerID:   t.ID,
		X:         t.X,
		Y:         t.Y,
		TargetID:  target.ID,
		Damage:    t.Damage,
		Speed:     500, // 像素/秒
		Type:      "normal",
		ArmorReduce: t.ArmorReduce,
	}

	return proj
}

// 升级塔
func (t *Tower) Upgrade() {
	t.Level++
	config := TowerConfigs[t.ID]
	t.Damage = config.BaseDamage + float64(t.Level-1)*config.LevelBonus
	t.Range = config.BaseRange * (1 + float64(t.Level)*0.05)
}

// 防御塔管理器
type TowerManager struct {
	Towers   map[string]*Tower
	mu       sync.RWMutex
}

func NewTowerManager() *TowerManager {
	return &TowerManager{
		Towers: make(map[string]*Tower),
	}
}

func (tm *TowerManager) Add(tower *Tower) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.Towers[tower.ID] = tower
}

func (tm *TowerManager) Get(id string) *Tower {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.Towers[id]
}

func (tm *TowerManager) GetAll() []*Tower {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	towers := make([]*Tower, 0, len(tm.Towers))
	for _, t := range tm.Towers {
		towers = append(towers, t)
	}
	return towers
}

func (tm *TowerManager) Remove(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	delete(tm.Towers, id)
}

// ============================================
// 投射物系统
// ============================================

type Projectile struct {
	ID           string  `json:"id"`
	TowerID      string  `json:"towerId"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	TargetID     string  `json:"targetId"`
	Damage       float64 `json:"damage"`
	Speed        float64 `json:"speed"` // 像素/秒
	Type         string  `json:"type"`  // normal, ice, lightning
	ArmorReduce  float64 `json:"armorReduce"`
	HitTime      int64   `json:"hitTime"`
}

// 更新投射物位置
func (p *Projectile) Update(dt float64, target *Enemy) (hit bool) {
	if target == nil || target.IsDead() {
		return true // 目标死亡，投射物消失
	}

	// 朝目标移动
	dx := target.X - p.X
	dy := target.Y - p.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 10 {
		// 命中
		return true
	}

	// 移动
	moveDist := p.Speed * dt
	p.X += dx / dist * moveDist
	p.Y += dy / dist * moveDist

	return false
}

// ============================================
// 敌人生成系统
// ============================================

type EnemyType int

const (
	EnemyTypeGrunt EnemyType = iota // 步兵
	EnemyTypeRanger                  // 远程
	EnemyTypeTank                    // 坦克
	EnemyTypeBoss                    // Boss
)

// 敌人
type Enemy struct {
	ID          string    `json:"id"`
	Type        EnemyType `json:"type"`
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Progress    float64   `json:"progress"` // 进度 0-1
	HP          float64   `json:"hp"`
	MaxHP       float64   `json:"maxHp"`
	Speed       float64   `json:"speed"`       // 移动速度
	Armor       float64   `json:"armor"`       // 护甲
	Reward      int       `json:"reward"`      // 击杀奖励
	PathIndex   int       `json:"pathIndex"`   // 路径点索引
	PathOffset  float64   `json:"pathOffset"`  // 路径偏移
	FreezeTime  int64     `json:"freezeTime"`  // 冰冻结束时间
	SlowFactor  float64    `json:"slowFactor"` // 减速因子
	mu          sync.RWMutex
}

// 敌人生成配置
var EnemyConfigs = map[string]struct {
	BaseHP  float64
	Speed   float64
	Armor   float64
	Reward  int
}{
	"grunt":     {HP: 50, Speed: 50, Armor: 0, Reward: 10},
	"ranger":    {HP: 30, Speed: 40, Armor: 5, Reward: 15},
	"tank":      {HP: 200, Speed: 25, Armor: 20, Reward: 30},
	"boss":      {HP: 1000, Speed: 20, Armor: 50, Reward: 100},
}

// 敌人生成器
type EnemySpawner struct {
	Wave       int
	Enemies    []*Enemy
	Path       []Point // 路径点
	SpawnQueue []SpawnInfo
	mu         sync.RWMutex
}

type Point struct {
	X float64
	Y float64
}

type SpawnInfo struct {
	Type     string
	Count    int
	Interval int // 毫秒
	Delay    int // 首次延迟 ms
}

// 波次配置
var WaveConfigs = []struct {
	Duration int // 波次持续时间 (秒)
	Enemies  []SpawnInfo
}{
	{
		Duration: 30,
		Enemies: []SpawnInfo{
			{Type: "grunt", Count: 10, Interval: 2000, Delay: 1000},
		},
	},
	{
		Duration: 45,
		Enemies: []SpawnInfo{
			{Type: "grunt", Count: 15, Interval: 1500, Delay: 1000},
			{Type: "ranger", Count: 5, Interval: 3000, Delay: 5000},
		},
	},
	{
		Duration: 60,
		Enemies: []SpawnInfo{
			{Type: "grunt", Count: 20, Interval: 1000, Delay: 1000},
			{Type: "ranger", Count: 10, Interval: 2000, Delay: 3000},
			{Type: "tank", Count: 3, Interval: 5000, Delay: 10000},
		},
	},
	{
		Duration: 90,
		Enemies: []SpawnInfo{
			{Type: "grunt", Count: 30, Interval: 800, Delay: 1000},
			{Type: "ranger", Count: 15, Interval: 1500, Delay: 2000},
			{Type: "tank", Count: 5, Interval: 4000, Delay: 8000},
			{Type: "boss", Count: 1, Interval: 0, Delay: 30000},
		},
	},
}

func NewEnemySpawner() *EnemySpawner {
	return &EnemySpawner{
		Wave:    0,
		Enemies: make([]*Enemy, 0),
		Path:    []Point{{X: 0, Y: 300}, {X: 200, Y: 300}, {X: 200, Y: 100}, {X: 600, Y: 100}, {X: 600, Y: 400}, {X: 800, Y: 400}},
	}
}

// 开始新波次
func (es *EnemySpawner) StartWave(waveNum int) {
	es.Wave = waveNum
	es.SpawnQueue = make([]SpawnInfo, 0)

	if waveNum-1 >= len(WaveConfigs) {
		// 最后一波重复
		waveNum = len(WaveConfigs)
	}

	config := WaveConfigs[waveNum-1]
	es.SpawnQueue = config.Enemies
}

// 生成敌人
func (es *EnemySpawner) Spawn() []*Enemy {
	newEnemies := make([]*Enemy, 0)

	for i, spawn := range es.SpawnQueue {
		if spawn.Count <= 0 {
			continue
		}

		enemy := es.createEnemy(spawn.Type)
		newEnemies = append(newEnemies, enemy)
		es.SpawnQueue[i].Count--
	}

	es.Enemies = append(es.Enemies, newEnemies...)
	return newEnemies
}

func (es *EnemySpawner) createEnemy(enemyType string) *Enemy {
	config := EnemyConfigs[enemyType]
	enemy := &Enemy{
		ID:         generateID(),
		Type:       EnemyTypeGrunt,
		X:          es.Path[0].X,
		Y:          es.Path[0].Y,
		Progress:   0,
		HP:         config.BaseHP,
		MaxHP:      config.BaseHP,
		Speed:      config.Speed,
		Armor:      config.Armor,
		Reward:     config.Reward,
		PathIndex:  0,
		PathOffset: 0,
		FreezeTime: 0,
		SlowFactor: 1.0,
	}

	switch enemyType {
	case "ranger":
		enemy.Type = EnemyTypeRanger
	case "tank":
		enemy.Type = EnemyTypeTank
	case "boss":
		enemy.Type = EnemyTypeBoss
	}

	return enemy
}

// 更新敌人位置
func (es *EnemySpawner) Update(dt float64) {
	now := time.Now().UnixMilli()

	for _, e := range es.Enemies {
		if e.IsDead() {
			continue
		}

		// 计算移动速度 (考虑减速效果)
		speed := e.Speed * e.SlowFactor
		if now < e.FreezeTime {
			speed = 0 // 冰冻
		}

		// 移动
		if e.PathIndex < len(es.Path)-1 {
			current := es.Path[e.PathIndex]
			next := es.Path[e.PathIndex+1]

			dx := next.X - current.X
			dy := next.Y - current.Y
			segmentLen := math.Sqrt(dx*dx + dy*dy)

			moveDist := speed * dt
			e.PathOffset += moveDist

			// 到达下一个路径点
			if e.PathOffset >= segmentLen {
				e.PathIndex++
				e.PathOffset = 0
			}

			// 更新坐标
			if e.PathIndex < len(es.Path)-1 {
				progress := e.PathOffset / segmentLen
				p1 := es.Path[e.PathIndex]
				p2 := es.Path[e.PathIndex+1]
				e.X = p1.X + (p2.X-p1.X)*progress
				e.Y = p1.Y + (p2.Y-p1.Y)*progress
			}
		}

		// 更新进度
		totalDist := es.CalcTotalPathLength()
		if totalDist > 0 {
			distCovered := 0.0
			for i := 0; i < e.PathIndex; i++ {
				p1 := es.Path[i]
				p2 := es.Path[i+1]
				distCovered += math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2))
			}
			distCovered += e.PathOffset
			e.Progress = distCovered / totalDist
		}
	}

	// 清理死亡和到达终点的敌人
	es.cleanup()
}

// 计算路径总长度
func (es *EnemySpawner) CalcTotalPathLength() float64 {
	total := 0.0
	for i := 0; i < len(es.Path)-1; i++ {
		p1 := es.Path[i]
		p2 := es.Path[i+1]
		total += math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2))
	}
	return total
}

// 清理敌人
func (es *EnemySpawner) cleanup() {
	alive := make([]*Enemy, 0)
	for _, e := range es.Enemies {
		if !e.IsDead() && e.Progress < 1.0 {
			alive = append(alive, e)
		}
	}
	es.Enemies = alive
}

// 是否死亡
func (e *Enemy) IsDead() bool {
	return e.HP <= 0
}

// 受到伤害
func (e *Enemy) TakeDamage(damage, armorReduce float64) {
	// 计算护甲减伤
	armor := e.Armor * (1 - armorReduce)
	actualDamage := damage * (1 - armor/(armor+100)) // 经典护甲公式

	e.HP -= actualDamage
}

// ============================================
// 战斗管理器
// ============================================

type BattleManager struct {
	Towers     *TowerManager
	Spawner    *EnemySpawner
	Projectiles []*Projectile
	State      BattleState
	mu         sync.RWMutex
}

type BattleState struct {
	Wave       int     `json:"wave"`
	Score      int     `json:"score"`
	Lives      int     `json:"lives"`      // 剩余生命
	Money      int     `json:"money"`      // 金币
	WaveTime   int     `json:"waveTime"`   // 波次剩余时间
	IsRunning  bool    `json:"isRunning"`
}

func NewBattleManager() *BattleManager {
	return &BattleManager{
		Towers:      NewTowerManager(),
		Spawner:     NewEnemySpawner(),
		Projectiles: make([]*Projectile, 0),
		State: BattleState{
			Wave:      1,
			Score:     0,
			Lives:     20,
			Money:     100,
			WaveTime:  30,
			IsRunning: false,
		},
	}
}

// 战斗主循环 (每帧调用)
func (bm *BattleManager) Update(dt float64) {
	if !bm.State.IsRunning {
		return
	}

	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 更新敌人
	bm.Spawner.Update(dt)

	// 塔攻击
	for _, tower := range bm.Towers.GetAll() {
		proj := tower.Attack(bm.Spawner.Enemies)
		if proj != nil {
			bm.Projectiles = append(bm.Projectiles, proj)
		}
	}

	// 更新投射物
	for i := len(bm.Projectiles) - 1; i >= 0; i-- {
		p := bm.Projectiles[i]
		target := bm.Spawner.GetEnemy(p.TargetID)
		hit := p.Update(dt, target)

		if hit && target != nil {
			target.TakeDamage(p.Damage, p.ArmorReduce)
			// 冰冻效果
			if p.Type == "ice" {
				target.FreezeTime = time.Now().UnixMilli() + 2000
				target.SlowFactor = 0.5
			}
			// 移除投射物
			bm.Projectiles = append(bm.Projectiles[:i], bm.Projectiles[i+1:]...)
		}
	}

	// 检查敌人死亡
	for _, e := range bm.Spawner.Enemies {
		if e.IsDead() {
			bm.State.Score += e.Reward
			bm.State.Money += e.Reward
		}
	}

	// 检查敌人到达终点
	for _, e := range bm.Spawner.Enemies {
		if e.Progress >= 1.0 {
			bm.State.Lives--
		}
	}

	// 游戏结束检查
	if bm.State.Lives <= 0 {
		bm.State.IsRunning = false
	}
}

func (es *EnemySpawner) GetEnemy(id string) *Enemy {
	for _, e := range es.Enemies {
		if e.ID == id {
			return e
		}
	}
	return nil
}

// 生成唯一ID
func generateID() string {
	return fmt.Sprintf("%d_%d", time.Now().UnixMilli(), now())
}

func now() int64 {
	return time.Now().UnixNano() % 10000
}

// ============================================
// 礼物系统 (Gift System)
// ============================================

type GiftType int

const (
	GiftTypeCoin GiftType = iota // 金币
	GiftTypeStar                  // 星星
	GiftTypeRocket               // 火箭
	GiftTypeCar                  // 跑车
	GiftTypePlane                // 飞机
	GiftTypeBang                 // 炸弹
)

// 礼物
type Gift struct {
	ID       string   `json:"id"`
	Type     GiftType `json:"type"`
	Name     string   `json:"name"`
	Price    int      `json:"price"`    // 价格(金币)
	Value    int      `json:"value"`    // 效果值
	Effect   string   `json:"effect"`   // 特效类型
	Duration int      `json:"duration"` // 持续时间(ms)
}

// 礼物配置
var GiftConfigs = map[string]Gift{
	"coin":    {Type: GiftTypeCoin, Name: "金币", Price: 1, Value: 10, Effect: "", Duration: 0},
	"star":    {Type: GiftTypeStar, Name: "星星", Price: 10, Value: 100, Effect: "star", Duration: 5000},
	"rocket":  {Type: GiftTypeRocket, Name: "火箭", Price: 50, Value: 500, Effect: "rocket", Duration: 3000},
	"car":     {Type: GiftTypeCar, Name: "跑车", Price: 100, Value: 1000, Effect: "car", Duration: 5000},
	"plane":   {Type: GiftTypePlane, Name: "飞机", Price: 200, Value: 2000, Effect: "plane", Duration: 8000},
	"bang":    {Type: GiftTypeBang, Name: "炸弹", Price: 30, Value: 300, Effect: "bomb", Duration: 0},
}

// 礼物管理器
type GiftManager struct {
	PendingGifts []*PendingGift // 待发放礼物
	ActiveEffects map[string]*GiftEffect // 活跃特效
	mu            sync.RWMutex
}

type PendingGift struct {
	Gift     *Gift
	SenderID string   // 发送者ID
	Receiver string   // 接收者 (空=全屏)
	Time     time.Time
}

type GiftEffect struct {
	GiftID    string
	Type      GiftType
	SenderID  string
	X         float64
	Y         float64
	StartTime time.Time
	Duration  int
	EndTime   time.Time
}

func NewGiftManager() *GiftManager {
	return &GiftManager{
		PendingGifts: make([]*PendingGift, 0),
		ActiveEffects: make(map[string]*GiftEffect),
	}
}

// 接收礼物
func (gm *GiftManager) ReceiveGift(giftType, senderID, receiver string) *GiftEffect {
	config, ok := GiftConfigs[giftType]
	if !ok {
		return nil
	}

	gift := &Gift{
		ID:       generateID(),
		Type:     config.Type,
		Name:     config.Name,
		Price:    config.Price,
		Value:    config.Value,
		Effect:   config.Effect,
		Duration: config.Duration,
	}

	pending := &PendingGift{
		Gift:     gift,
		SenderID: senderID,
		Receiver: receiver,
		Time:     time.Now(),
	}

	gm.mu.Lock()
	gm.PendingGifts = append(gm.PendingGifts, pending)

	// 创建特效
	effect := &GiftEffect{
		GiftID:    gift.ID,
		Type:      gift.Type,
		SenderID:  senderID,
		X:         400, // 默认屏幕中央
		Y:         300,
		StartTime: time.Now(),
		Duration:  gift.Duration,
		EndTime:   time.Now().Add(time.Duration(gift.Duration) * time.Millisecond),
	}
	gm.ActiveEffects[gift.ID] = effect
	gm.mu.Unlock()

	return effect
}

// 处理礼物效果
func (gm *GiftManager) ProcessGiftEffects() {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	now := time.Now()
	for id, effect := range gm.ActiveEffects {
		if now.After(effect.EndTime) {
			delete(gm.ActiveEffects, id)
		}
	}
}

// 获取活跃特效
func (gm *GiftManager) GetActiveEffects() []*GiftEffect {
	gm.mu.RLock()
	defer gm.mu.RUnlock()

	effects := make([]*GiftEffect, 0, len(gm.ActiveEffects))
	for _, e := range gm.ActiveEffects {
		effects = append(effects, e)
	}
	return effects
}

// 礼物效果应用到战斗
func (gm *GiftManager) ApplyGiftEffect(battle *BattleManager) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	for _, effect := range gm.ActiveEffects {
		switch effect.Type {
		case GiftTypeCoin:
			// 金币效果: 增加金钱
			battle.State.Money += GiftConfigs["coin"].Value
		case GiftTypeStar:
			// 星星效果: 增加分数
			battle.State.Score += GiftConfigs["star"].Value
		case GiftTypeBang:
			// 炸弹效果: 对所有敌人造成伤害
			for _, enemy := range battle.Spawner.Enemies {
				if !enemy.IsDead() {
					enemy.TakeDamage(300, 0) // 300伤害，无护甲削减
				}
			}
		}
	}
}

// ============================================
// 弹幕系统 (Bullet/Danmaku System)
// ============================================

type DanmakuType int

const (
	DanmakuTypeText DanmakuType = iota // 文字弹幕
	DanmakuTypeImage                    // 图片弹幕
	DanmakuTypeVoice                    // 语音弹幕
)

// 弹幕
type Danmaku struct {
	ID       string      `json:"id"`
	Type     DanmakuType `json:"type"`
	Content  string      `json:"content"`  // 内容
	SenderID string      `json:"senderId"` // 发送者ID
	Sender   string      `json:"sender"`   // 发送者名称
	X        float64     `json:"x"`        // 位置X
	Y        float64     `json:"y"`        // 位置Y
	Speed    float64     `json:"speed"`    // 移动速度
	Color    string      `json:"color"`    // 颜色
	FontSize int         `json:"fontSize"` // 字号
	Opacity  float64     `json:"opacity"`  // 透明度
	LifeTime int64       `json:"lifeTime"` // 生存时间(ms)
}

// 弹幕配置
var DanmakuConfigs = map[string]struct {
	Speed    float64
	Color    string
	FontSize int
	Opacity  float64
}{
	"normal":   {Speed: 200, Color: "#FFFFFF", FontSize: 24, Opacity: 1.0},
	"premium":  {Speed: 150, Color: "#FFD700", FontSize: 32, Opacity: 1.0},
	"rainbow":  {Speed: 180, Color: "#FF69B4", FontSize: 28, Opacity: 0.9},
}

// 弹幕管理器
type DanmakuManager struct {
	DanmakuList []*Danmaku
	MaxCount    int
	mu          sync.RWMutex
}

func NewDanmakuManager() *DanmakuManager {
	return &DanmakuManager{
		DanmakuList: make([]*Danmaku, 0),
		MaxCount:   100, // 最多同时显示100条弹幕
	}
}

// 发送弹幕
func (dm *DanmakuManager) Send(text, senderID, sender, configType string) *Danmaku {
	config, ok := DanmakuConfigs[configType]
	if !ok {
		config = DanmakuConfigs["normal"]
	}

	danmaku := &Danmaku{
		ID:       generateID(),
		Type:     DanmakuTypeText,
		Content:  text,
		SenderID: senderID,
		Sender:   sender,
		X:        1200, // 从屏幕右侧开始
		Y:        float64(50 + len(dm.DanmakuList)%20*30), // 随机Y轴位置
		Speed:    config.Speed,
		Color:    config.Color,
		FontSize: config.FontSize,
		Opacity:  config.Opacity,
		LifeTime: 10000, // 10秒
	}

	dm.mu.Lock()
	defer dm.mu.Unlock()

	// 超过最大数量，移除最老的
	if len(dm.DanmakuList) >= dm.MaxCount {
		dm.DanmakuList = dm.DanmakuList[1:]
	}

	dm.DanmakuList = append(dm.DanmakuList, danmaku)
	return danmaku
}

// 更新弹幕位置
func (dm *DanmakuManager) Update(dt float64) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	now := time.Now().UnixMilli()
	alive := make([]*Danmaku, 0)

	for _, d := range dm.DanmakuList {
		// 向左移动
		d.X -= d.Speed * dt

		// 超出屏幕或超时则移除
		if d.X > -200 && now < d.LifeTime {
			alive = append(alive, d)
		}
	}

	dm.DanmakuList = alive
}

// 获取所有弹幕
func (dm *DanmakuManager) GetAll() []*Danmaku {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	result := make([]*Danmaku, len(dm.DanmakuList))
	copy(result, dm.DanmakuList)
	return result
}

// 过滤弹幕 (敏感词过滤示例)
func (dm *DanmakuManager) Filter(text string) string {
	// 简单敏感词过滤示例
	sensitiveWords := []string{"bad", "test"}
	for _, word := range sensitiveWords {
		text = replaceAll(text, word, "***")
	}
	return text
}

func replaceAll(s, old, new string) string {
	result := s
	for {
		i := find(result, old)
		if i < 0 {
			break
		}
		result = result[:i] + new + result[i+len(old):]
	}
	return result
}

func find(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// ============================================
// 直播间系统 (Live Room System)
// ============================================

type LiveRoom struct {
	RoomID       string         // 房间ID
	AnchorID     string         // 主播ID
	AnchorName   string         // 主播名称
	Battle       *BattleManager // 战斗实例
	GiftManager  *GiftManager   // 礼物管理器
	DanmakuMgr   *DanmakuManager // 弹幕管理器
	Viewers      map[string]*Viewer // 观众列表
	TotalRevenue int             // 总收入
	TotalGift    int             // 总礼物数
	StartTime    time.Time
	mu           sync.RWMutex
}

type Viewer struct {
	ID       string
	Name     string
	Coin     int    // 金币
	IsVIP    bool   // 是否VIP
	JoinTime time.Time
}

func NewLiveRoom(roomID, anchorID, anchorName string) *LiveRoom {
	return &LiveRoom{
		RoomID:      roomID,
		AnchorID:    anchorID,
		AnchorName:  anchorName,
		Battle:      NewBattleManager(),
		GiftManager: NewGiftManager(),
		DanmakuMgr:  NewDanmakuManager(),
		Viewers:     make(map[string]*Viewer),
		StartTime:   time.Now(),
	}
}

// 观众进入
func (lr *LiveRoom) JoinViewer(viewerID, name string) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	lr.Viewers[viewerID] = &Viewer{
		ID:       viewerID,
		Name:     name,
		Coin:     100, // 初始金币
		IsVIP:    false,
		JoinTime: time.Now(),
	}
}

// 观众离开
func (lr *LiveRoom) LeaveViewer(viewerID string) {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	delete(lr.Viewers, viewerID)
}

// 发送礼物
func (lr *LiveRoom) SendGift(viewerID, giftType string) *GiftEffect {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	viewer, ok := lr.Viewers[viewerID]
	if !ok {
		return nil
	}

	config, ok := GiftConfigs[giftType]
	if !ok || viewer.Coin < config.Price {
		return nil
	}

	// 扣除金币
	viewer.Coin -= config.Price
	lr.TotalRevenue += config.Price
	lr.TotalGift++

	return lr.GiftManager.ReceiveGift(giftType, viewerID, "")
}

// 发送弹幕
func (lr *LiveRoom) SendDanmaku(viewerID, content string) *Danmaku {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	viewer, ok := lr.Viewers[viewerID]
	if !ok {
		return nil
	}

	// 过滤敏感词
	content = lr.DanmakuMgr.Filter(content)

	configType := "normal"
	if viewer.IsVIP {
		configType = "premium"
	}

	return lr.DanmakuMgr.Send(content, viewerID, viewer.Name, configType)
}

// 主循环更新
func (lr *LiveRoom) Update(dt float64) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	// 更新战斗
	lr.Battle.Update(dt)

	// 更新弹幕
	lr.DanmakuMgr.Update(dt)

	// 处理礼物特效
	lr.GiftManager.ProcessGiftEffects()

	// 应用礼物效果到战斗
	lr.GiftManager.ApplyGiftEffect(lr.Battle)
}
