package main

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

// ============================================
// 弹幕模式系统 (Bullet Pattern System)
// 弹幕游戏核心模块
// ============================================

// BulletPatternType 弹幕类型
type BulletPatternType int

const (
	PatternNone BulletPatternType = iota
	PatternStraight   // 直线
	PatternSpiral     // 螺旋
	PatternRing       // 环形
	PatternWave       // 波浪
	PatternRain       // 雨幕
	PatternSpread     // 散射
	PatternAimed      // 瞄准
	PatternRandom     // 随机
	PatternComplex    // 复合模式
)

// Bullet 子弹实体
type Bullet struct {
	ID        string         `json:"id"`
	X         float64        `json:"x"`
	Y         float64        `json:"y"`
	VX        float64        `json:"vx"`
	VY        float64        `json:"vy"`
	Angle     float64        `json:"angle"`     // 飞行角度
	Speed     float64        `json:"speed"`     // 速度
	Damage    float64        `json:"damage"`    // 伤害
	Radius    float64        `json:"radius"`    // 大小
	Type      BulletPatternType `json:"type"`   // 弹幕类型
	Color     string         `json:"color"`     // 颜色
	OwnerID   string         `json:"ownerId"`   // 发射者ID
	IsEnemy   bool           `json:"isEnemy"`   // 是否敌弹
	LifeTime  int64          `json:"lifetime"`  // 存活时间(ms)
	MaxLife   int64          `json:"maxLife"`   // 最大存活时间
	CreatedAt int64          `json:"createdAt"` // 创建时间
	Extra     map[string]interface{} `json:"extra"` // 扩展数据
}

// BulletPattern 弹幕模式
type BulletPattern struct {
	Type           BulletPatternType `json:"type"`
	Name           string            `json:"name"`
	BulletCount    int               `json:"bulletCount"`    // 子弹数量
	Interval       int64             `json:"interval"`       // 发射间隔(ms)
	Speed          float64           `json:"speed"`          // 基础速度
	SpeedVariance  float64           `json:"speedVariance"`  // 速度随机范围
	AngleStart     float64           `json:"angleStart"`     // 起始角度
	AngleEnd       float64           `json:"angleEnd"`       // 结束角度
	AngleStep      float64           `json:"angleStep"`      // 角度步进
	Radius         float64           `json:"radius"`         // 子弹半径
	Damage         float64           `json:"damage"`         // 基础伤害
	DamageScale    float64           `json:"damageScale"`    // 伤害系数
	Color          string             `json:"color"`          // 颜色
	Duration       int64              `json:"duration"`       // 持续时间(ms)
	SpinSpeed      float64            `json:"spinSpeed"`      // 旋转速度
	SpreadAngle    float64            `json:"spreadAngle"`    // 散射角度
	AimAccuracy    float64            `json:"aimAccuracy"`    // 瞄准精度(0-1)
	Delay          int64              `json:"delay"`          // 延迟发射(ms)
	Loop           bool               `json:"loop"`           // 是否循环
	WaveAmplitude  float64            `json:"waveAmplitude"`  // 波浪振幅
	WaveFrequency  float64            `json:"waveFrequency"`  // 波浪频率
}

// PatternEmitter 弹幕发射器
type PatternEmitter struct {
	ID          string          `json:"id"`
	Pattern     *BulletPattern  `json:"pattern"`
	X           float64         `json:"x"`
	Y           float64         `json:"y"`
	TargetX     float64         `json:"targetX"`     // 目标X (用于瞄准)
	TargetY     float64         `json:"targetY"`     // 目标Y
	StartAngle  float64         `json:"startAngle"`  // 起始角度
	CurrentTick int64           `json:"currentTick"` // 当前tick
	LastEmit    int64           `json:"lastEmit"`    // 上次发射时间
	Active      bool            `json:"active"`      // 是否激活
	LoopCount   int             `json:"loopCount"`    // 循环次数
	mu          sync.RWMutex
}

// BulletManager 子弹管理器
type BulletManager struct {
	Bullets    map[string]*Bullet
	Emitters   map[string]*PatternEmitter
	Pool       *sync.Pool
	mu         sync.RWMutex
	NextBullet int64
	NextEmitter int64
}

// 创建子弹管理器
func NewBulletManager() *BulletManager {
	return &BulletManager{
		Bullets:    make(map[string]*Bullet),
		Emitters:   make(map[string]*PatternEmitter),
		Pool: &sync.Pool{
			New: func() interface{} {
				return &Bullet{
					Extra: make(map[string]interface{}),
				}
			},
		},
		NextBullet: 0,
		NextEmitter: 0,
	}
}

// 创建子弹
func (bm *BulletManager) CreateBullet() *Bullet {
	bullet := bm.Pool.Get().(*Bullet)
	bullet.ID = ""
	bullet.Extra = make(map[string]interface{})
	return bullet
}

// 回收子弹
func (bm *BulletManager) RecycleBullet(b *Bullet) {
	bm.Pool.Put(b)
}

// 生成唯一ID
func (bm *BulletManager) genBulletID() string {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.NextBullet++
	return "bullet_" + string(rune('a'+bm.NextBullet%26)) + string(rune('0'+bm.NextBullet%10))
}

// 创建子弹 (工厂方法)
func (bm *BulletManager) SpawnBullet(x, y, vx, vy, damage, radius float64, patternType BulletPatternType, ownerID string, isEnemy bool) *Bullet {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	
	bullet := &Bullet{
		ID:        bm.genBulletID(),
		X:         x,
		Y:         y,
		VX:        vx,
		VY:        vy,
		Angle:     math.Atan2(vy, vx),
		Speed:     math.Sqrt(vx*vx + vy*vy),
		Damage:    damage,
		Radius:    radius,
		Type:      patternType,
		Color:     "#FFFFFF",
		OwnerID:   ownerID,
		IsEnemy:   isEnemy,
		LifeTime:  0,
		MaxLife:   10000,
		CreatedAt: time.Now().UnixMilli(),
		Extra:     make(map[string]interface{}),
	}
	
	bm.Bullets[bullet.ID] = bullet
	return bullet
}

// 移除子弹
func (bm *BulletManager) RemoveBullet(id string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if bullet, ok := bm.Bullets[id]; ok {
		bm.RecycleBullet(bullet)
		delete(bm.Bullets, id)
	}
}

// 更新子弹位置
func (bm *BulletManager) UpdateBullets(dt int64) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	
	now := time.Now().UnixMilli()
	
	for id, bullet := range bm.Bullets {
		bullet.LifeTime = now - bullet.CreatedAt
		
		// 检查生命周期
		if bullet.LifeTime > bullet.MaxLife {
			delete(bm.Bullets, id)
			bm.RecycleBullet(bullet)
			continue
		}
		
		// 根据弹幕类型应用特殊运动
		switch bullet.Type {
		case PatternWave:
			// 波浪运动
			bullet.VY += math.Sin(float64(bullet.LifeTime)/100*math.Pi) * 0.5
		case PatternSpiral:
			// 螺旋 - 角度自然变化
			angle := math.Atan2(bullet.VY, bullet.VX)
			angle += 0.05 // 旋转
			speed := math.Sqrt(bullet.VX*bullet.VX + bullet.VY*bullet.VY)
			bullet.VX = math.Cos(angle) * speed
			bullet.VY = math.Sin(angle) * speed
		}
		
		// 更新位置
		bullet.X += bullet.VX * float64(dt) / 16.67
		bullet.Y += bullet.VY * float64(dt) / 16.67
		
		// 更新角度
		if bullet.VX != 0 || bullet.VY != 0 {
			bullet.Angle = math.Atan2(bullet.VY, bullet.VX)
		}
	}
}

// 发射直线弹幕
func (pe *PatternEmitter) EmitStraight(bm *BulletManager, count int, angleOffset float64) {
	for i := 0; i < count; i++ {
		angle := pe.Pattern.AngleStart + float64(i)*pe.Pattern.AngleStep + angleOffset
		rad := angle * math.Pi / 180
		
		speed := pe.Pattern.Speed
		if pe.Pattern.SpeedVariance > 0 {
			speed += (rand.Float64()*2 - 1) * pe.Pattern.SpeedVariance
		}
		
		vx := math.Cos(rad) * speed
		vy := math.Sin(rad) * speed
		
		bm.SpawnBullet(pe.X, pe.Y, vx, vy, pe.Pattern.Damage, pe.Pattern.Radius, pe.Pattern.Type, pe.ID, true)
	}
}

// 发射环形弹幕
func (pe *PatternEmitter) EmitRing(bm *BulletManager) {
	count := pe.Pattern.BulletCount
	for i := 0; i < count; i++ {
		angle := 360.0 / float64(count) * float64(i)
		rad := angle * math.Pi / 180
		
		vx := math.Cos(rad) * pe.Pattern.Speed
		vy := math.Sin(rad) * pe.Pattern.Speed
		
		bm.SpawnBullet(pe.X, pe.Y, vx, vy, pe.Pattern.Damage, pe.Pattern.Radius, PatternRing, pe.ID, true)
	}
}

// 发射螺旋弹幕
func (pe *PatternEmitter)EmitSpiral(bm *BulletManager, tick int64) {
	angle := pe.Pattern.AngleStart + float64(tick)*pe.Pattern.SpinSpeed
	rad := angle * math.Pi / 180
	
	vx := math.Cos(rad) * pe.Pattern.Speed
	vy := math.Sin(rad) * pe.Pattern.Speed
	
	bm.SpawnBullet(pe.X, pe.Y, vx, vy, pe.Pattern.Damage, pe.Pattern.Radius, PatternSpiral, pe.ID, true)
}

// 发射散射弹幕
func (pe *PatternEmitter) EmitSpread(bm *BulletManager, aimAngle float64) {
	count := pe.Pattern.BulletCount
	halfSpread := pe.Pattern.SpreadAngle / 2
	
	for i := 0; i < count; i++ {
		offset := -halfSpread + (pe.Pattern.SpreadAngle / float64(count-1)) * float64(i)
		if count == 1 {
			offset = 0
		}
		angle := aimAngle + offset
		rad := angle * math.Pi / 180
		
		vx := math.Cos(rad) * pe.Pattern.Speed
		vy := math.Sin(rad) * pe.Pattern.Speed
		
		bm.SpawnBullet(pe.X, pe.Y, vx, vy, pe.Pattern.Damage, pe.Pattern.Radius, PatternSpread, pe.ID, true)
	}
}

// 发射瞄准弹幕
func (pe *PatternEmitter) EmitAimed(bm *BulletManager) {
	// 计算瞄准角度
	dx := pe.TargetX - pe.X
	dy := pe.TargetY - pe.Y
	aimAngle := math.Atan2(dy, dx) * 180 / math.Pi
	
	// 添加精度偏移
	if pe.Pattern.AimAccuracy < 1.0 {
		inaccuracy := (1.0 - pe.Pattern.AimAccuracy) * 30 // 最大30度偏差
		aimAngle += (rand.Float64()*2 - 1) * inaccuracy
	}
	
	pe.EmitSpread(bm, aimAngle)
}

// 发射随机弹幕
func (pe *PatternEmitter) EmitRandom(bm *BulletManager) {
	for i := 0; i < pe.Pattern.BulletCount; i++ {
		angle := rand.Float64() * 360
		rad := angle * math.Pi / 180
		
		speed := pe.Pattern.Speed * (0.5 + rand.Float64())
		
		vx := math.Cos(rad) * speed
		vy := math.Sin(rad) * speed
		
		bm.SpawnBullet(pe.X, pe.Y, vx, vy, pe.Pattern.Damage, pe.Pattern.Radius, PatternRandom, pe.ID, true)
	}
}

// 更新发射器
func (pe *PatternEmitter) Update(bm *BulletManager, now int64) bool {
	if !pe.Active {
		return false
	}
	
	// 检查延迟
	if pe.Pattern.Delay > 0 && now - pe.CurrentTick*16 < pe.Pattern.Delay {
		return true
	}
	
	// 检查间隔
	if now-pe.LastEmit < pe.Pattern.Interval {
		return true
	}
	
	// 旋转起始角度 (用于螺旋)
	pe.StartAngle += pe.Pattern.SpinSpeed
	
	// 根据模式发射
	switch pe.Pattern.Type {
	case PatternStraight:
		pe.EmitStraight(bm, pe.Pattern.BulletCount, 0)
	case PatternRing:
		pe.EmitRing(bm)
	case PatternSpiral:
		pe.EmitSpiral(bm, pe.CurrentTick)
	case PatternSpread:
		pe.EmitAimed(bm)
	case PatternAimed:
		pe.EmitAimed(bm)
	case PatternRandom:
		pe.EmitRandom(bm)
	case PatternWave:
		pe.EmitStraight(bm, pe.Pattern.BulletCount, pe.StartAngle)
	case PatternRain:
		pe.EmitStraight(bm, pe.Pattern.BulletCount, -90) // 垂直向下
	}
	
	pe.LastEmit = now
	pe.CurrentTick++
	
	// 检查循环终止
	if !pe.Pattern.Loop && pe.CurrentTick >= pe.Pattern.Duration/pe.Pattern.Interval {
		pe.Active = false
		return false
	}
	
	return true
}

// 创建预设弹幕模式
var PresetPatterns = map[string]*BulletPattern{
	// 基础模式
	"basic_ring": {
		Type: PatternRing, Name: "基础环", BulletCount: 12,
		Speed: 3, AngleStep: 30, Radius: 8, Damage: 10,
		Color: "#FF6B6B", Interval: 1000, Loop: true,
	},
	"basic_spiral": {
		Type: PatternSpiral, Name: "基础螺旋", BulletCount: 1,
		Speed: 4, AngleStart: 0, SpinSpeed: 5, Radius: 6, Damage: 8,
		Color: "#4ECDC4", Interval: 50, Loop: true,
	},
	"basic_spread": {
		Type: PatternSpread, Name: "散射", BulletCount: 5,
		Speed: 5, SpreadAngle: 60, Radius: 7, Damage: 12,
		Color: "#FFE66D", Interval: 800, AimAccuracy: 0.9,
	},
	// 中级模式
	"mid_double_spiral": {
		Type: PatternSpiral, Name: "双螺旋", BulletCount: 2,
		Speed: 4.5, SpinSpeed: 4, Radius: 6, Damage: 10,
		Color: "#A855F7", Interval: 40, Loop: true,
	},
	"mid_wave": {
		Type: PatternWave, Name: "波浪", BulletCount: 8,
		Speed: 3, AngleStart: 0, AngleStep: 15, WaveAmplitude: 20,
		WaveFrequency: 0.1, Radius: 6, Damage: 8,
		Color: "#06B6D4", Interval: 100, Duration: 5000,
	},
	// 高级模式
	"hard_complex": {
		Type: PatternComplex, Name: "复合弹幕", BulletCount: 24,
		Speed: 5, AngleStart: 0, SpinSpeed: 3, Radius: 5, Damage: 15,
		Color: "#EF4444", Interval: 30, Duration: 8000, Loop: true,
	},
	"hard_rain": {
		Type: PatternRain, Name: "弹幕雨", BulletCount: 10,
		Speed: 6, AngleStep: 36, Radius: 4, Damage: 8,
		Color: "#3B82F6", Interval: 100, Loop: true,
	},
	// Boss模式
	"boss_helix": {
		Type: PatternSpiral, Name: "螺旋风暴", BulletCount: 3,
		Speed: 5, SpinSpeed: 6, Radius: 7, Damage: 20,
		Color: "#DC2626", Interval: 30, Loop: true,
	},
	"boss_flower": {
		Type: PatternRing, Name: "花朵绽放", BulletCount: 36,
		Speed: 4, Radius: 8, Damage: 25,
		Color: "#F472B6", Interval: 1500, Duration: 5000,
	},
}

// 获取预设模式
func GetPattern(name string) *BulletPattern {
	if p, ok := PresetPatterns[name]; ok {
		return p
	}
	return nil
}

// ============================================
// 弹幕效果器 (Danmaku Effect System)
// ============================================

// DanmakuEffect 弹幕特效
type DanmakuEffect struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	X         float64        `json:"x"`
	Y         float64        `json:"y"`
	Scale     float64        `json:"scale"`
	Alpha     float64        `json:"alpha"`
	Rotation  float64        `json:"rotation"`
	Duration  int64          `json:"duration"`
	StartTime int64          `json:"startTime"`
	EndTime   int64          `json:"endTime"`
}

// DanmakuEffectManager 特效管理器
type DanmakuEffectManager struct {
	Effects map[string]*DanmakuEffect
	mu       sync.RWMutex
}

func NewDanmakuEffectManager() *DanmakuEffectManager {
	return &DanmakuEffectManager{
		Effects: make(map[string]*DanmakuEffect),
	}
}

// 添加特效
func (dm *DanmakuEffectManager) AddEffect(e *DanmakuEffect) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.Effects[e.ID] = e
}

// 移除特效
func (dm *DanmakuEffectManager) RemoveEffect(id string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	delete(dm.Effects, id)
}

// 更新特效
func (dm *DanmakuEffectManager) Update(now int64) []*DanmakuEffect {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	
	var expired []string
	var result []*DanmakuEffect
	
	for id, e := range dm.Effects {
		if now > e.EndTime {
			expired = append(expired, id)
			continue
		}
		
		// 计算当前进度
		progress := float64(now-e.StartTime) / float64(e.Duration)
		e.Alpha = 1.0 - progress
		e.Scale = 1.0 + progress*0.5
		
		result = append(result, e)
	}
	
	for _, id := range expired {
		delete(dm.Effects, id)
	}
	
	return result
}
