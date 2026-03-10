package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// ============================================
// 游戏引擎核心 (Game Engine Core)
// ============================================

// 游戏状态
type GameState int

const (
	GameStateIdle GameState = iota // 空闲
	GameStateLobby     // 大厅
	GameStatePlaying   // 游戏中
	GameStatePaused    // 暂停
	GameStateEnded     // 结束
)

// 游戏引擎
type GameEngine struct {
	// 核心组件
	RoomManager    *RoomManager
	BattleManager  *BattleManager
	GiftManager    *GiftManager
	DanmakuManager *DanmakuManager
	
	// 状态
	State       GameState
	CurrentRoom *LiveRoom
	FrameCount  int64
	FPS         float64
	LastFrame   int64
	
	// 配置
	Config *GameConfig
	
	// 控制
	mu      sync.RWMutex
	running bool
}

type GameConfig struct {
	FPS           int     // 目标帧率
	Width         float64 // 地图宽度
	Height        float64 // 地图高度
	AutoStartWave bool    // 自动开始波次
	WaveDelay     int     // 波次间隔(秒)
	MaxLives      int     // 最大生命
	StartingMoney int     // 初始金币
}

// 全局游戏引擎
var GlobalEngine *GameEngine

// 初始化游戏引擎
func NewGameEngine() *GameEngine {
	return &GameEngine{
		RoomManager:    NewRoomManager(),
		BattleManager:  NewBattleManager(),
		GiftManager:    NewGiftManager(),
		DanmakuManager: NewDanmakuManager(),
		State:          GameStateIdle,
		Config: &GameConfig{
			FPS:           60,
			Width:         1200,
			Height:        800,
			AutoStartWave: true,
			WaveDelay:     10,
			MaxLives:      20,
			StartingMoney: 100,
		},
		running: false,
	}
}

// 创建游戏房间
func (ge *GameEngine) CreateRoom(roomID, hostID string) *LiveRoom {
	room := NewLiveRoom(roomID, hostID, "主播")
	room.Battle.State.Lives = ge.Config.MaxLives
	room.Battle.State.Money = ge.Config.StartingMoney
	room.Battle.State.IsRunning = false
	
	ge.mu.Lock()
	ge.CurrentRoom = room
	ge.State = GameStateLobby
	ge.mu.Unlock()
	
	return room
}

// 开始游戏
func (ge *GameEngine) StartGame() error {
	ge.mu.Lock()
	defer ge.mu.Unlock()
	
	if ge.CurrentRoom == nil {
		return fmt.Errorf("no room available")
	}
	
	ge.State = GameStatePlaying
	ge.CurrentRoom.Battle.State.IsRunning = true
	ge.CurrentRoom.Battle.State.Wave = 1
	ge.CurrentRoom.Battle.Spawner.StartWave(1)
	
	return nil
}

// 暂停游戏
func (ge *GameEngine) PauseGame() {
	ge.mu.Lock()
	defer ge.mu.Unlock()
	
	if ge.State == GameStatePlaying {
		ge.State = GameStatePaused
		ge.CurrentRoom.Battle.State.IsRunning = false
	}
}

// 恢复游戏
func (ge *GameEngine) ResumeGame() {
	ge.mu.Lock()
	defer ge.mu.Unlock()
	
	if ge.State == GameStatePaused {
		ge.State = GameStatePlaying
		ge.CurrentRoom.Battle.State.IsRunning = true
	}
}

// 结束游戏
func (ge *GameEngine) EndGame() *BattleResult {
	ge.mu.Lock()
	defer ge.mu.Unlock()
	
	ge.State = GameStateEnded
	ge.CurrentRoom.Battle.State.IsRunning = false
	
	// 计算战斗结果
	result := &BattleResult{
		Score:      ge.CurrentRoom.Battle.State.Score,
		Wave:       ge.CurrentRoom.Battle.State.Wave,
		LivesLeft:  ge.CurrentRoom.Battle.State.Lives,
		MoneyEarned: ge.CurrentRoom.Battle.State.Money,
		Duration:   time.Since(ge.CurrentRoom.StartTime),
		Win:        ge.CurrentRoom.Battle.State.Lives > 0,
	}
	
	return result
}

// 游戏主循环 (帧更新)
func (ge *GameEngine) Update() {
	ge.mu.Lock()
	
	if ge.State != GameStatePlaying {
		ge.mu.Unlock()
		return
	}
	
	room := ge.CurrentRoom
	ge.mu.Unlock()
	
	// 计算帧时间
	now := time.Now().UnixMilli()
	dt := float64(now-ge.LastFrame) / 1000.0
	if dt > 0 {
		ge.FPS = ge.FPS*0.9 + (1/dt)*0.1 // 平滑FPS
	}
	ge.LastFrame = now
	ge.FrameCount++
	
	// 限制最大帧时间，防止跳跃
	if dt > 0.1 {
		dt = 0.1
	}
	
	// 更新战斗
	room.Battle.Update(dt)
	
	// 更新弹幕
	room.DanmakuMgr.Update(dt)
	
	// 处理礼物特效
	room.GiftManager.ProcessGiftEffects()
	room.GiftManager.ApplyGiftEffect(room.Battle)
	
	// 波次管理
	ge.updateWave()
	
	// 检查游戏结束
	if room.Battle.State.Lives <= 0 {
		ge.EndGame()
	}
}

// 波次更新
func (ge *GameEngine) updateWave() {
	room := ge.CurrentRoom
	spawner := room.Battle.Spawner
	
	// 检查当前波次是否完成 (没有存活的敌人)
	hasAlive := false
	for _, e := range spawner.Enemies {
		if !e.IsDead() && e.Progress < 1.0 {
			hasAlive = true
			break
		}
	}
	
	if !hasAlive && ge.Config.AutoStartWave {
		// 波次完成，检查是否还有波次
		if spawner.Wave < len(WaveConfigs) {
			// 延时切换波次
			room.Battle.State.WaveTime--
			if room.Battle.State.WaveTime <= 0 {
				spawner.StartWave(spawner.Wave + 1)
				room.Battle.State.Wave = spawner.Wave
				room.Battle.State.WaveTime = ge.Config.WaveDelay
			}
		}
	}
}

// ============================================
// 塔操作API
// ============================================

// 放置防御塔
func (ge *GameEngine) PlaceTower(towerID string, towerType TowerType, x, y float64) (*Tower, error) {
	ge.mu.Lock()
	defer ge.mu.Unlock()
	
	if ge.CurrentRoom == nil {
		return nil, fmt.Errorf("no room")
	}
	
	// 检查位置是否有效
	if !ge.isValidPosition(x, y) {
		return nil, fmt.Errorf("invalid position")
	}
	
	// 检查金币
	cost := ge.getTowerCost(towerID)
	if ge.CurrentRoom.Battle.State.Money < cost {
		return nil, fmt.Errorf("insufficient money")
	}
	
	// 创建塔
	tower := NewTower(towerID, towerType, x, y, 1)
	ge.CurrentRoom.Battle.Towers.Add(tower)
	ge.CurrentRoom.Battle.State.Money -= cost
	
	return tower, nil
}

// 升级防御塔
func (ge *GameEngine) UpgradeTower(towerID string) error {
	ge.mu.Lock()
	defer ge.mu.Unlock()
	
	if ge.CurrentRoom == nil {
		return fmt.Errorf("no room")
	}
	
	tower := ge.CurrentRoom.Battle.Towers.Get(towerID)
	if tower == nil {
		return fmt.Errorf("tower not found")
	}
	
	// 检查金币
	cost := ge.getUpgradeCost(tower.Level)
	if ge.CurrentRoom.Battle.State.Money < cost {
		return fmt.Errorf("insufficient money")
	}
	
	tower.Upgrade()
	ge.CurrentRoom.Battle.State.Money -= cost
	
	return nil
}

// 出售防御塔
func (ge *GameEngine) SellTower(towerID string) (int, error) {
	ge.mu.Lock()
	defer ge.mu.Unlock()
	
	if ge.CurrentRoom == nil {
		return 0, fmt.Errorf("no room")
	}
	
	tower := ge.CurrentRoom.Battle.Towers.Get(towerID)
	if tower == nil {
		return 0, fmt.Errorf("tower not found")
	}
	
	// 返还50%金币
	refund := ge.getTowerCost(tower.ID) * tower.Level / 2
	ge.CurrentRoom.Battle.Towers.Remove(towerID)
	ge.CurrentRoom.Battle.State.Money += refund
	
	return refund, nil
}

// 检查位置是否有效
func (ge *GameEngine) isValidPosition(x, y float64) bool {
	// 边界检查
	if x < 50 || x > ge.Config.Width-50 || y < 50 || y > ge.Config.Height-50 {
		return false
	}
	
	// 检查与其他塔的碰撞
	for _, t := range ge.CurrentRoom.Battle.Towers.GetAll() {
		dist := math.Sqrt(math.Pow(t.X-x, 2) + math.Pow(t.Y-y, 2))
		if dist < 60 { // 塔之间最小距离
			return false
		}
	}
	
	return true
}

// 获取塔的价格
func (ge *GameEngine) getTowerCost(towerID string) int {
	costs := map[string]int{
		"arrow":     50,
		"cannon":    100,
		"ice":       80,
		"lightning": 120,
		"tower_heal": 60,
	}
	if cost, ok := costs[towerID]; ok {
		return cost
	}
	return 50
}

// 获取升级价格
func (ge *GameEngine) getUpgradeCost(level int) int {
	return 50 * level
}

// ============================================
// 战斗结算
// ============================================

type BattleResult struct {
	Score      int           `json:"score"`
	Wave       int           `json:"wave"`
	LivesLeft  int           `json:"livesLeft"`
	MoneyEarned int          `json:"moneyEarned"`
	Duration   time.Duration `json:"duration"`
	Win        bool          `json:"win"`
	Rewards    *RewardData   `json:"rewards"`
}

type RewardData struct {
	Exp     int `json:"exp"`     // 经验
	Coins   int `json:"coins"`  // 金币
	Gems    int `json:"gems"`   // 钻石
}

// 计算战斗奖励
func (br *BattleResult) CalcRewards() {
	baseExp := br.Wave * 100
	baseCoins := br.Score / 10
	
	if br.Win {
		br.Rewards = &RewardData{
			Exp:   baseExp + 50,
			Coins: baseCoins + 100,
			Gems:  br.Wave,
		}
	} else {
		br.Rewards = &RewardData{
			Exp:   baseExp / 2,
			Coins: baseCoins / 2,
			Gems:  0,
		}
	}
}

// ============================================
// 游戏事件系统
// ============================================

type GameEventType int

const (
	EventTowerPlaced GameEventType = iota
	EventTowerUpgraded
	EventTowerSold
	EventEnemyKilled
	EventEnemyEscaped
	EventWaveStarted
	EventWaveCompleted
	EventGiftReceived
	EventDanmakuSent
	EventGameStart
	EventGameEnd
)

// 游戏事件
type GameEvent struct {
	Type    GameEventType `json:"type"`
	Time    time.Time     `json:"time"`
	Data    interface{}  `json:"data"`
	RoomID  string        `json:"roomId"`
}

// 事件监听器
type GameEventListener func(event *GameEvent)

// 事件管理器
type EventManager struct {
	listeners map[GameEventType][]GameEventListener
	mu        sync.RWMutex
}

func NewEventManager() *EventManager {
	return &EventManager{
		listeners: make(map[GameEventType][]GameEventListener),
	}
}

// 注册事件监听
func (em *EventManager) On(eventType GameEventType, listener GameEventListener) {
	em.mu.Lock()
	defer em.mu.Unlock()
	
	em.listeners[eventType] = append(em.listeners[eventType], listener)
}

// 触发事件
func (em *EventManager) Emit(event *GameEvent) {
	em.mu.RLock()
	defer em.mu.RUnlock()
	
	if listeners, ok := em.listeners[event.Type]; ok {
		for _, listener := range listeners {
			listener(event)
		}
	}
}

// ============================================
// 存档系统
// ============================================

type SaveData struct {
	PlayerID    string       `json:"playerId"`
	PlayerName  string       `json:"playerName"`
	Level       int          `json:"level"`
	Exp         int          `json:"exp"`
	Coins       int          `json:"coins"`
	Gems        int          `json:"gems"`
	HighScore   int          `json:"highScore"`
	Towers      []string     `json:"towers"` // 已解锁塔
	Skins       []string     `json:"skins"`  // 已解锁皮肤
	Settings    *PlayerSettings `json:"settings"`
	LastLogin   time.Time    `json:"lastLogin"`
}

type PlayerSettings struct {
	MusicVolume   float64 `json:"musicVolume"`
	SFXVolume    float64 `json:"sfxVolume"`
	Graphics     string  `json:"graphics"` // low/medium/high
	ShowDanmaku  bool    `json:"showDanmaku"`
}

// 保存游戏数据
func (ge *GameEngine) SaveGame(playerID string) *SaveData {
	ge.mu.RLock()
	defer ge.mu.RUnlock()
	
	return &SaveData{
		PlayerID:   playerID,
		Level:      1,
		Exp:        0,
		Coins:      ge.CurrentRoom.Battle.State.Money,
		Gems:       0,
		HighScore:  ge.CurrentRoom.Battle.State.Score,
		LastLogin:  time.Now(),
		Settings: &PlayerSettings{
			MusicVolume:  0.7,
			SFXVolume:    0.8,
			Graphics:     "medium",
			ShowDanmaku:  true,
		},
	}
}

// 加载游戏数据
func (ge *GameEngine) LoadGame(data *SaveData) {
	ge.mu.Lock()
	defer ge.mu.Unlock()
	
	if ge.CurrentRoom != nil {
		ge.CurrentRoom.Battle.State.Money = data.Coins
	}
}

// ============================================
// 初始化全局引擎
// ============================================

func init() {
	GlobalEngine = NewGameEngine()
}
