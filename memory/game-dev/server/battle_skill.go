package main

import (
	"math"
	"sync"
	"time"
)

// ============================================
// Boss 技能系统
// ============================================

type BossSkillType int

const (
	BossSkillNone BossSkillType = iota
	BossSkillRage       // 愤怒: 短时间内提升攻速和伤害
	BossSkillShield     // 护盾: 免疫下一次伤害
	BossSkillSummon     // 召唤: 召唤小怪
	BossSkillAOE        // AOE: 全屏攻击
	BossSkillHeal       // 治疗: 回复自身HP
	BossSkillSlow       // 减速: 减速所有塔
	BossSkillExecute    // 处决: 对低血量敌人秒杀
)

// Boss技能
type BossSkill struct {
	Type       BossSkillType `json:"type"`
	Name       string        `json:"name"`
	Cooldown   int           `json:"cooldown"`   // 冷却时间(ms)
	Duration   int           `json:"duration"`   // 持续时间(ms)
	Damage     float64       `json:"damage"`    // 伤害值
	Value      float64       `json:"value"`      // 效果值 (如减速百分比)
	LastUsed   int64         `json:"lastUsed"`   // 上次使用时间
	IsActive   bool          `json:"isActive"`   // 是否激活
	ActiveTime int64         `json:"activeTime"` // 激活时间
}

// Boss配置
var BossConfigs = map[string]struct {
	BaseHP      float64
	Speed       float64
	Armor       float64
	Reward      int
	Skills      []BossSkillType
	SkillCooldown int // 技能冷却
}{
	"dragon": {
		BaseHP: 2000, Speed: 25, Armor: 30, Reward: 500,
		Skills: []BossSkillType{BossSkillRage, BossSkillAOE, BossSkillSummon},
		SkillCooldown: 15000,
	},
	"golem": {
		BaseHP: 3000, Speed: 15, Armor: 50, Reward: 800,
		Skills: []BossSkillType{BossSkillShield, BossSkillHeal, BossSkillAOE},
		SkillCooldown: 20000,
	},
	"demon": {
		BaseHP: 1500, Speed: 35, Armor: 20, Reward: 600,
		Skills: []BossSkillType{BossSkillRage, BossSkillExecute, BossSkillSlow},
		SkillCooldown: 10000,
	},
}

// Boss敌人
type BossEnemy struct {
	*Enemy
	SkillList   []BossSkill
	CurrentSkill *BossSkill
}

// 创建Boss
func NewBossEnemy(id, bossType string) *BossEnemy {
	config := BossConfigs[bossType]
	enemy := &Enemy{
		ID:         id,
		Type:       EnemyTypeBoss,
		X:          0,
		Y:          300,
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

	// 初始化技能
	skills := make([]BossSkill, len(config.Skills))
	for i, skillType := range config.Skills {
		skills[i] = BossSkill{
			Type:     skillType,
			Name:     GetSkillName(skillType),
			Cooldown: config.SkillCooldown,
			LastUsed: 0,
			IsActive: false,
		}
	}

	return &BossEnemy{
		Enemy:       enemy,
		SkillList:   skills,
		CurrentSkill: nil,
	}
}

func GetSkillName(skillType BossSkillType) string {
	switch skillType {
	case BossSkillRage:
		return "愤怒"
	case BossSkillShield:
		return "护盾"
	case BossSkillSummon:
		return "召唤"
	case BossSkillAOE:
		return "AOE攻击"
	case BossSkillHeal:
		return "治疗"
	case BossSkillSlow:
		return "减速"
	case BossSkillExecute:
		return "处决"
	default:
		return "未知"
	}
}

// Boss释放技能
func (b *BossEnemy) UseSkill(battle *BattleManager) bool {
	now := time.Now().UnixMilli()

	// 检查是否有可用技能
	for i := range b.SkillList {
		skill := &b.SkillList[i]
		
		// 冷却中或已激活
		if now-skill.LastUsed < skill.Cooldown || skill.IsActive {
			continue
		}

		// 使用技能
		skill.LastUsed = now
		skill.IsActive = true
		skill.ActiveTime = now
		b.CurrentSkill = skill

		// 执行技能效果
		b.executeSkill(skill, battle)
		return true
	}

	// 检查技能是否结束
	for i := range b.SkillList {
		skill := &b.SkillList[i]
		if skill.IsActive && now-skill.ActiveTime > skill.Duration {
			skill.IsActive = false
			b.CurrentSkill = nil
		}
	}

	return false
}

// 执行Boss技能
func (b *BossEnemy) executeSkill(skill *BossSkill, battle *BattleManager) {
	switch skill.Type {
	case BossSkillRage:
		// 愤怒: 提升伤害和攻速
		b.Speed *= 1.5
		b.Armor += 10

	case BossSkillShield:
		// 护盾: 免疫伤害 (通过设置护盾标记)
		// 实际逻辑在TakeDamage中检查

	case BossSkillSummon:
		// 召唤: 生成小怪
		for i := 0; i < 3; i++ {
			enemy := battle.Spawner.createEnemy("grunt")
			enemy.X = b.X + float64(i*30)
			enemy.Y = b.Y
			battle.Spawner.Enemies = append(battle.Spawner.Enemies, enemy)
		}

	case BossSkillAOE:
		// AOE: 对所有塔造成伤害
		for _, tower := range battle.Towers.GetAll() {
			dist := math.Sqrt(math.Pow(tower.X-b.X, 2) + math.Pow(tower.Y-b.Y, 2))
			if dist < 200 {
				// 简单处理: 塔受到伤害会降级或损坏
				tower.Damage *= 0.8
			}
		}

	case BossSkillHeal:
		// 治疗: 回复HP
		healAmount := b.MaxHP * 0.3
		b.HP = math.Min(b.HP+healAmount, b.MaxHP)

	case BossSkillSlow:
		// 减速: 减速所有塔 (降低攻速)
		for _, tower := range battle.Towers.GetAll() {
			tower.FireRate *= 0.5
		}

	case BossSkillExecute:
		// 处决: 对低血量敌人秒杀 (<20% HP)
		executeThreshold := 0.2
		for _, enemy := range battle.Spawner.Enemies {
			if enemy.HP/enemy.MaxHP < executeThreshold {
				enemy.HP = 0
			}
		}
	}
}

// Boss受到伤害 (检查护盾)
func (b *BossEnemy) TakeDamage(damage, armorReduce float64) bool {
	for _, skill := range b.SkillList {
		if skill.IsActive && skill.Type == BossSkillShield {
			// 护盾生效，免疫伤害
			return false
		}
	}

	// 计算护甲减伤
	armor := b.Armor * (1 - armorReduce)
	actualDamage := damage * (1 - armor/(armor+100))
	b.HP -= actualDamage
	return true
}

// ============================================
// 特殊事件系统
// ============================================

type EventType int

const (
	EventMeteorShower EventType = iota // 流星雨: 全屏伤害
	EventTreasureBox                    // 宝箱: 随机奖励
	EventEarthquake                     // 地震: 塔暂时失效
	EventGoldRush                       // 黄金 rush: 金币加成
	EventBossRush                       // Boss rush: 快速召唤Boss
)

// 特殊事件
type SpecialEvent struct {
	Type      EventType `json:"type"`
	Name      string    `json:"name"`
	Duration  int       `json:"duration"`  // 持续时间(ms)
	StartTime int64     `json:"startTime"`
	EndTime   int64     `json:"endTime"`
	Value     float64   `json:"value"`     // 效果值
}

// 事件管理器
type EventManager struct {
	ActiveEvents map[EventType]*SpecialEvent
	EventQueue   []SpecialEvent // 待触发事件
	mu           sync.RWMutex
}

var EventConfigs = map[EventType]struct {
	Name      string
	Duration  int
	Weight    int // 权重 (越高越容易触发)
}{
	EventMeteorShower: {Name: "流星雨", Duration: 5000, Weight: 20},
	EventTreasureBox:  {Name: "宝箱", Duration: 0, Weight: 15},
	EventEarthquake:   {Name: "地震", Duration: 8000, Weight: 10},
	EventGoldRush:     {Name: "黄金 rush", Duration: 10000, Weight: 25},
	EventBossRush:     {Name: "Boss rush", Duration: 0, Weight: 5},
}

func NewEventManager() *EventManager {
	return &EventManager{
		ActiveEvents: make(map[EventType]*SpecialEvent),
		EventQueue:   make([]SpecialEvent, 0),
	}
}

// 触发随机事件
func (em *EventManager) TriggerRandomEvent() *SpecialEvent {
	// 简单随机选择
	totalWeight := 0
	for _, config := range EventConfigs {
		totalWeight += config.Weight
	}

	random := now() % int64(totalWeight)
	current := int64(0)

	var selectedType EventType
	for eventType, config := range EventConfigs {
		current += int64(config.Weight)
		if random < current {
			selectedType = eventType
			break
		}
	}

	return em.TriggerEvent(selectedType)
}

// 触发指定事件
func (em *EventManager) TriggerEvent(eventType EventType) *SpecialEvent {
	config, ok := EventConfigs[eventType]
	if !ok {
		return nil
	}

	now := time.Now().UnixMilli()
	event := &SpecialEvent{
		Type:      eventType,
		Name:      config.Name,
		Duration:  config.Duration,
		StartTime:  now,
		EndTime:   now + int64(config.Duration),
		Value:      0,
	}

	em.mu.Lock()
	em.ActiveEvents[eventType] = event
	em.mu.Unlock()

	return event
}

// 应用事件效果
func (em *EventManager) ApplyEventEffect(battle *BattleManager) {
	em.mu.Lock()
	defer em.mu.Unlock()

	now := time.Now().UnixMilli()

	for eventType, event := range em.ActiveEvents {
		if now > event.EndTime {
			// 事件结束，移除
			delete(em.ActiveEvents, eventType)
			continue
		}

		// 应用事件效果
		switch eventType {
		case EventMeteorShower:
			// 流星雨: 每秒对随机敌人造成伤害
			if now%1000 < 16 { // 简单每秒触发
				for _, enemy := range battle.Spawner.Enemies {
					if !enemy.IsDead() {
						enemy.TakeDamage(50, 0)
					}
				}
			}

		case EventGoldRush:
			// 黄金 rush: 金币加成
			battle.State.Money += 1 // 每帧+1金币

		case EventEarthquake:
			// 地震: 塔攻速减半 (在事件结束时恢复)
			// 实际效果在事件结束时恢复
		}
	}
}

// 结束事件后清理
func (em *EventManager) OnEventEnd(eventType EventType, battle *BattleManager) {
	switch eventType {
	case EventEarthquake:
		// 恢复塔的攻速
		for _, tower := range battle.Towers.GetAll() {
			config := TowerConfigs[tower.ID]
			tower.FireRate = config.BaseFireRate
		}
	}
}

// ============================================
// 战斗加成系统 (Buff/Debuff)
// ============================================

type BuffType int

const (
	BuffDamage BuffType = iota // 伤害加成
	BuffAttackSpeed           // 攻速加成
	BuffRange                 // 范围加成
	BuffGold                  // 金币加成
	BuffXP                    // 经验加成
	DebuffSlow                // 减速
	DebuffPoison              // 中毒
	DebuffStun                // 眩晕
)

// 战斗Buff
type Buff struct {
	Type      BuffType   `json:"type"`
	Name      string     `json:"name"`
	Value     float64    `json:"value"`     // 效果值
	Duration  int        `json:"duration"`  // 持续时间(ms)
	StartTime int64      `json:"startTime"`
	EndTime   int64      `json:"endTime"`
	Stackable bool       `json:"stackable"` // 是否可叠加
	MaxStack  int        `json:"maxStack"`  // 最大叠加层数
	Stack     int        `json:"stack"`     // 当前层数
}

// Buff管理器
type BuffManager struct {
	PlayerBuffs map[BuffType]*Buff  // 玩家Buff
	TowerBuffs  map[string]map[BuffType]*Buff // 塔Buff [towerID][buffType]
	EnemyBuffs  map[string]map[BuffType]*Buff // 敌人Debuff
	mu          sync.RWMutex
}

func NewBuffManager() *BuffManager {
	return &BuffManager{
		PlayerBuffs: make(map[BuffType]*Buff),
		TowerBuffs:  make(map[string]map[BuffType]*Buff),
		EnemyBuffs:  make(map[string]map[BuffType]*Buff),
	}
}

// 添加Buff
func (bm *BuffManager) AddBuff(targetType string, targetID string, buffType BuffType, value float64, duration int) *Buff {
	now := time.Now().UnixMilli()
	
	buff := &Buff{
		Type:      buffType,
		Name:      GetBuffName(buffType),
		Value:     value,
		Duration: duration,
		StartTime: now,
		EndTime:   now + int64(duration),
		Stackable: false,
		MaxStack:  1,
		Stack:     1,
	}

	bm.mu.Lock()
	defer bm.mu.Unlock()

	switch targetType {
	case "player":
		if existing, ok := bm.PlayerBuffs[buffType]; ok && existing.Stackable {
			if existing.Stack < existing.MaxStack {
				existing.Stack++
				existing.EndTime = now + int64(duration)
			}
			return existing
		}
		bm.PlayerBuffs[buffType] = buff

	case "tower":
		if _, ok := bm.TowerBuffs[targetID]; !ok {
			bm.TowerBuffs[targetID] = make(map[BuffType]*Buff)
		}
		bm.TowerBuffs[targetID][buffType] = buff

	case "enemy":
		if _, ok := bm.EnemyBuffs[targetID]; !ok {
			bm.EnemyBuffs[targetID] = make(map[BuffType]*Buff)
		}
		bm.EnemyBuffs[targetID][buffType] = buff
	}

	return buff
}

// 获取Buff效果值
func (bm *BuffManager) GetBuffValue(targetType string, targetID string, buffType BuffType) float64 {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	var buff *Buff
	switch targetType {
	case "player":
		buff = bm.PlayerBuffs[buffType]
	case "tower":
		if buffs, ok := bm.TowerBuffs[targetID]; ok {
			buff = buffs[buffType]
		}
	case "enemy":
		if buffs, ok := bm.EnemyBuffs[targetID]; ok {
			buff = buffs[buffType]
		}
	}

	if buff == nil || time.Now().UnixMilli() > buff.EndTime {
		return 0
	}

	return buff.Value * float64(buff.Stack)
}

// 更新Buff状态
func (bm *BuffManager) Update() {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	now := time.Now().UnixMilli()

	// 清理过期Buff
	for buffType, buff := range bm.PlayerBuffs {
		if now > buff.EndTime {
			delete(bm.PlayerBuffs, buffType)
		}
	}

	for targetID, buffs := range bm.TowerBuffs {
		for buffType, buff := range buffs {
			if now > buff.EndTime {
				delete(buffs, buffType)
			}
		}
		if len(buffs) == 0 {
			delete(bm.TowerBuffs, targetID)
		}
	}

	for targetID, buffs := range bm.EnemyBuffs {
		for buffType, buff := range buffs {
			if now > buff.EndTime {
				delete(buffs, buffType)
			}
		}
		if len(buffs) == 0 {
			delete(bm.EnemyBuffs, targetID)
		}
	}
}

func GetBuffName(buffType BuffType) string {
	switch buffType {
	case BuffDamage:
		return "伤害加成"
	case BuffAttackSpeed:
		return "攻速加成"
	case BuffRange:
		return "范围加成"
	case BuffGold:
		return "金币加成"
	case BuffXP:
		return "经验加成"
	case DebuffSlow:
		return "减速"
	case DebuffPoison:
		return "中毒"
	case DebuffStun:
		return "眩晕"
	default:
		return "未知"
	}
}

// ============================================
// 战斗结算系统
// ============================================

type BattleResult struct {
	Win            bool      `json:"win"`
	Wave           int       `json:"wave"`
	Score          int       `json:"score"`
	MoneyEarned    int       `json:"moneyEarned"`
	XPEarned       int       `json:"xpEarned"`
	TowersBuilt    int       `json:"towersBuilt"`
	EnemiesKilled  int       `json:"enemiesKilled"`
	BossKilled     int       `json:"bossKilled"`
	DamageDealt    float64   `json:"damageDealt"`
	Duration       int       `json:"duration"` // 秒
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
}

// 战斗统计
type BattleStats struct {
	TotalDamage    float64
	EnemiesKilled  int
	BossKilled     int
	TowersBuilt    int
	GoldEarned     int
	XPEarned       int
}

// 开始战斗
func (bm *BattleManager) StartBattle() {
	bm.Stats = &BattleStats{}
	bm.State.IsRunning = true
	bm.State.Wave = 1
	bm.State.Lives = 20
	bm.State.Score = 0
	bm.State.Money = 100
	bm.StartTime = time.Now()
	bm.Spawner.StartWave(1)
}

// 结束战斗
func (bm *BattleManager) EndBattle() *BattleResult {
	bm.State.IsRunning = false
	bm.EndTime = time.Now()

	result := &BattleResult{
		Win:           bm.State.Lives > 0,
		Wave:          bm.State.Wave,
		Score:         bm.State.Score,
		MoneyEarned:   bm.Stats.GoldEarned,
		XPEarned:      bm.Stats.XPEarned,
		TowersBuilt:   bm.Stats.TowersBuilt,
		EnemiesKilled: bm.Stats.EnemiesKilled,
		BossKilled:    bm.Stats.BossKilled,
		DamageDealt:   bm.Stats.TotalDamage,
		Duration:      int(bm.EndTime.Sub(bm.StartTime).Seconds()),
		StartTime:     bm.StartTime,
		EndTime:       bm.EndTime,
	}

	return result
}

// 扩展BattleManager
type BattleManagerEx struct {
	*TowerManager
	*EnemySpawner
	*GiftManager
	*DanmakuManager
	*BuffManager
	*EventManager
	*BattleStats
	State     BattleState
	StartTime time.Time
	EndTime   time.Time
	mu        sync.RWMutex
}

func NewBattleManagerEx() *BattleManagerEx {
	return &BattleManagerEx{
		TowerManager:  NewTowerManager(),
		EnemySpawner:  NewEnemySpawner(),
		GiftManager:  NewGiftManager(),
		DanmakuMgr:   NewDanmakuManager(),
		BuffManager:  NewBuffManager(),
		EventManager: NewEventManager(),
		Stats:        &BattleStats{},
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

// 主循环 (扩展版)
func (bm *BattleManagerEx) Update(dt float64) {
	if !bm.State.IsRunning {
		return
	}

	bm.mu.Lock()
	defer bm.mu.Unlock()

	// 更新敌人
	bm.EnemySpawner.Update(dt)

	// 塔攻击
	for _, tower := range bm.TowerManager.GetAll() {
		// 应用Buff效果
		damageBonus := bm.BuffManager.GetBuffValue("tower", tower.ID, BuffDamage)
		fireRateBonus := bm.BuffManager.GetBuffValue("tower", tower.ID, BuffAttackSpeed)
		
		tower.Damage *= (1 + damageBonus)
		tower.FireRate *= (1 + fireRateBonus)

		proj := tower.Attack(bm.EnemySpawner.Enemies)
		if proj != nil {
			// 投射物逻辑...
		}
	}

	// 更新Buff
	bm.BuffManager.Update()

	// 事件效果
	bm.EventManager.ApplyEventEffect((*BattleManager)(nil))

	// 更新礼物/弹幕
	bm.GiftManager.ProcessGiftEffects()
	bm.DanmakuMgr.Update(dt)
}
