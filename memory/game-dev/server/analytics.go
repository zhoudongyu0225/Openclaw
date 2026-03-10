package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ==================== 事件类型定义 ====================

// EventType 事件类型
type EventType string

const (
	EventTypeGameStart    EventType = "game_start"     // 游戏开始
	EventTypeGameEnd      EventType = "game_end"       // 游戏结束
	EventTypeLevelStart   EventType = "level_start"    // 关卡开始
	EventTypeLevelComplete EventType = "level_complete" // 关卡完成
	EventTypeLevelFail    EventType = "level_fail"      // 关卡失败
	EventTypeTowerBuild   EventType = "tower_build"    // 建造塔
	EventTypeTowerUpgrade EventType = "tower_upgrade"  // 升级塔
	EventTypeTowerSell    EventType = "tower_sell"     // 出售塔
	EventTypeEnemyKill    EventType = "enemy_kill"     // 击杀敌人
	EventTypeBossKill     EventType = "boss_kill"      // 击杀Boss
	EventTypeWaveStart    EventType = "wave_start"      // 波次开始
	EventTypeWaveComplete EventType = "wave_complete"  // 波次完成
	EventTypeGiftReceive  EventType = "gift_receive"   // 收到礼物
	EventTypeDanmaku      EventType = "danmaku"         // 发送弹幕
	EventTypeItemBuy      EventType = "item_buy"        // 购买道具
	EventTypeItemUse      EventType = "item_use"        // 使用道具
	EventTypeAchievement  EventType = "achievement"     // 成就解锁
	EventTypeQuestComplete EventType = "quest_complete" // 任务完成
	EventTypeCurrencyGain EventType = "currency_gain"  // 获得货币
	EventTypeCurrencySpend EventType = "currency_spend" // 消费货币
	EventTypePlayerDie    EventType = "player_die"      // 玩家死亡
	EventTypePlayerRevive EventType = "player_revive"  // 玩家复活
	EventTypeCustom       EventType = "custom"          // 自定义事件
)

// ==================== 事件结构 ====================

// GameEvent 游戏事件
type GameEvent struct {
	EventID    string                 `json:"event_id"`     // 事件唯一ID
	EventType  EventType               `json:"event_type"`  // 事件类型
	PlayerID   string                 `json:"player_id"`    // 玩家ID
	RoomID     string                 `json:"room_id"`      // 房间ID
	Timestamp  int64                  `json:"timestamp"`    // 时间戳(毫秒)
	Sequence   uint64                 `json:"sequence"`     // 序列号
	Properties map[string]interface{} `json:"properties"`   // 事件属性
	DeviceInfo *DeviceInfo             `json:"device_info"`  // 设备信息
	Location   *LocationInfo          `json:"location"`     // 位置信息
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	OS           string `json:"os"`            // 操作系统
	OSVersion    string `json:"os_version"`    // 系统版本
	DeviceModel  string `json:"device_model"`  // 设备型号
	AppVersion   string `json:"app_version"`   // APP版本
	NetworkType  string `json:"network_type"`  // 网络类型
	ScreenWidth  int    `json:"screen_width"`  // 屏幕宽度
	ScreenHeight int    `json:"screen_height"` // 屏幕高度
	Language     string `json:"language"`      // 语言
	Country      string `json:"country"`       // 国家
}

// LocationInfo 位置信息
type LocationInfo struct {
	Latitude  float64 `json:"latitude"`  // 纬度
	Longitude float64 `json:"longitude"` // 经度
	Altitude  float64 `json:"altitude"` // 海拔
}

// ==================== 用户画像 ====================

// PlayerProfile 玩家画像
type PlayerProfile struct {
	PlayerID      string                 `json:"player_id"`       // 玩家ID
	RegisterTime  int64                  `json:"register_time"`   // 注册时间
	TotalPlayTime int64                  `json:"total_play_time"` // 总游戏时长(秒)
	TotalGames    int                    `json:"total_games"`     // 总游戏场次
	WinGames      int                    `json:"win_games"`       // 胜利场次
	Level         int                    `json:"level"`           // 玩家等级
	Experience    int64                  `json:"experience"`      // 经验值
	Currency      map[string int64        `json:"currency"`       // 货币数量
	Achievements  []string               `json:"achievements"`    // 已解锁成就
	TowersBuilt   int                     `json:"towers_built"`    // 建造塔总数
	EnemiesKilled int                     `json:"enemies_killed"`  // 击杀敌人总数
	BossKilled    int                     `json:"boss_killed"`     // 击杀Boss总数
	LastLoginTime int64                  `json:"last_login_time"` // 最后登录时间
	LoginDays     int                    `json:"login_days"`      // 登录天数
	VIPLevel      int                    `json:"vip_level"`       // VIP等级
	Tags          []string               `json:"tags"`            // 玩家标签
	Attributes    map[string]interface{} `json:"attributes"`      // 自定义属性
}

// ==================== 数据分析器 ====================

// Analytics 数据分析器
type Analytics struct {
	mu           sync.RWMutex
	eventBuffer  chan *GameEvent        // 事件缓冲通道
	playerCache  map[string]*PlayerProfile // 玩家画像缓存
	stats        *AnalyticsStats        // 统计数据
	flushInterval time.Duration         // 刷新间隔
	maxBufferSize int                    // 最大缓冲大小
	isRunning    bool                    // 是否运行中
}

// AnalyticsStats 分析统计数据
type AnalyticsStats struct {
	TotalEvents    int64            `json:"total_events"`     // 总事件数
	EventsByType   map[EventType]int64 `json:"events_by_type"` // 按类型统计
	ActivePlayers  int               `json:"active_players"` // 活跃玩家数
	NewPlayers     int               `json:"new_players"`     // 新增玩家数
	RetentionRate  float64          `json:"retention_rate"`  // 留存率
	AvgPlayTime    float64          `json:"avg_play_time"`   // 平均游戏时长
	PeakOnline     int               `json:"peak_online"`    // 峰值在线
	TodayGames     int               `json:"today_games"`     // 今日游戏场次
	TodayRevenue   float64          `json:"today_revenue"`   // 今日收入
}

// NewAnalytics 创建数据分析器
func NewAnalytics(bufferSize int, flushInterval time.Duration) *Analytics {
	a := &Analytics{
		eventBuffer:  make(chan *GameEvent, bufferSize),
		playerCache:  make(map[string]*PlayerProfile),
		stats:        &AnalyticsStats{
			EventsByType: make(map[EventType]int64),
		},
		flushInterval: flushInterval,
		maxBufferSize: bufferSize,
		isRunning:    false,
	}
	return a
}

// Start 启动数据分析器
func (a *Analytics) Start() {
	if a.isRunning {
		return
	}
	a.isRunning = true
	go a.eventLoop()
	go a.statsLoop()
}

// Stop 停止数据分析器
func (a *Analytics) Stop() {
	a.isRunning = false
	close(a.eventBuffer)
}

// eventLoop 事件处理循环
func (a *Analytics) eventLoop() {
	ticker := time.NewTicker(a.flushInterval)
	defer ticker.Stop()

	events := make([]*GameEvent, 0, a.maxBufferSize)

	for {
		select {
		case event, ok := <-a.eventBuffer:
			if !ok {
				// 处理剩余事件
				a.flushEvents(events)
				return
			}
			events = append(events, event)
			if len(events) >= a.maxBufferSize {
				a.flushEvents(events)
				events = make([]*GameEvent, 0, a.maxBufferSize)
			}
		case <-ticker.C:
			if len(events) > 0 {
				a.flushEvents(events)
				events = make([]*GameEvent, 0, a.maxBufferSize)
			}
		}
	}
}

// statsLoop 统计循环
func (a *Analytics) statsLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if !a.isRunning {
			return
		}
		a.updateStats()
	}
}

// flushEvents 刷新事件到存储
func (a *Analytics) flushEvents(events []*GameEvent) {
	if len(events) == 0 {
		return
	}
	// TODO: 实现实际存储逻辑 (写入数据库/文件/发送到分析服务)
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, event := range events {
		a.stats.TotalEvents++
		a.stats.EventsByType[event.EventType]++
	}

	// 模拟存储
	fmt.Printf("[Analytics] Flushed %d events\n", len(events))
}

// updateStats 更新统计数据
func (a *Analytics) updateStats() {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 计算活跃玩家
	a.stats.ActivePlayers = len(a.playerCache)
}

// ==================== 事件记录 ====================

// TrackEvent 记录事件
func (a *Analytics) TrackEvent(event *GameEvent) {
	if !a.isRunning {
		return
	}
	event.Timestamp = time.Now().UnixMilli()
	if event.EventID == "" {
		event.EventID = fmt.Sprintf("%d_%s", event.Timestamp, generateUUID())
	}

	select {
	case a.eventBuffer <- event:
	default:
		fmt.Printf("[Analytics] Buffer full, dropping event: %s\n", event.EventType)
	}
}

// TrackGameStart 记录游戏开始
func (a *Analytics) TrackGameStart(playerID, roomID string, properties map[string]interface{}) {
	event := &GameEvent{
		EventType:  EventTypeGameStart,
		PlayerID:   playerID,
		RoomID:     roomID,
		Properties: properties,
	}
	a.TrackEvent(event)
}

// TrackGameEnd 记录游戏结束
func (a *Analytics) TrackGameEnd(playerID, roomID string, properties map[string]interface{}) {
	event := &GameEvent{
		EventType:  EventTypeGameEnd,
		PlayerID:   playerID,
		RoomID:     roomID,
		Properties: properties,
	}
	a.TrackEvent(event)
}

// TrackLevelEvent 记录关卡事件
func (a *Analytics) TrackLevelEvent(eventType EventType, playerID, roomID string, level int, properties map[string]interface{}) {
	if properties == nil {
		properties = make(map[string]interface{})
	}
	properties["level"] = level
	event := &GameEvent{
		EventType:  eventType,
		PlayerID:   playerID,
		RoomID:     roomID,
		Properties: properties,
	}
	a.TrackEvent(event)
}

// TrackTowerEvent 记录塔事件
func (a *Analytics) TrackTowerEvent(eventType EventType, playerID, roomID string, towerID string, towerType string, level int, properties map[string]interface{}) {
	if properties == nil {
		properties = make(map[string]interface{})
	}
	properties["tower_id"] = towerID
	properties["tower_type"] = towerType
	properties["tower_level"] = level
	event := &GameEvent{
		EventType:  eventType,
		PlayerID:   playerID,
		RoomID:     roomID,
		Properties: properties,
	}
	a.TrackEvent(event)
}

// TrackEnemyKill 记录击杀敌人
func (a *Analytics) TrackEnemyKill(playerID, roomID string, enemyType string, damage int64, properties map[string]interface{}) {
	if properties == nil {
		properties = make(map[string]interface{})
	}
	properties["enemy_type"] = enemyType
	properties["damage"] = damage
	event := &GameEvent{
		EventType:  EventTypeEnemyKill,
		PlayerID:   playerID,
		RoomID:     roomID,
		Properties: properties,
	}
	a.TrackEvent(event)
}

// TrackGift 记录礼物事件
func (a *Analytics) TrackGift(playerID, roomID string, giftType string, giftCount int, giftValue int64, senderID string) {
	event := &GameEvent{
		EventType: EventTypeGiftReceive,
		PlayerID:  playerID,
		RoomID:    roomID,
		Properties: map[string]interface{}{
			"gift_type":   giftType,
			"gift_count":  giftCount,
			"gift_value":  giftValue,
			"sender_id":   senderID,
		},
	}
	a.TrackEvent(event)
}

// TrackDanmaku 记录弹幕事件
func (a *Analytics) TrackDanmaku(playerID, roomID string, content string, danmakuType string) {
	event := &GameEvent{
		EventType: EventTypeDanmaku,
		PlayerID:  playerID,
		RoomID:    roomID,
		Properties: map[string]interface{}{
			"content":       content,
			"danmaku_type":  danmakuType,
			"content_length": len(content),
		},
	}
	a.TrackEvent(event)
}

// ==================== 玩家画像 ====================

// GetPlayerProfile 获取玩家画像
func (a *Analytics) GetPlayerProfile(playerID string) *PlayerProfile {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.playerCache[playerID]
}

// UpdatePlayerProfile 更新玩家画像
func (a *Analytics) UpdatePlayerProfile(profile *PlayerProfile) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.playerCache[profile.PlayerID] = profile
}

// CreatePlayerProfile 创建玩家画像
func (a *Analytics) CreatePlayerProfile(playerID string) *PlayerProfile {
	profile := &PlayerProfile{
		PlayerID:     playerID,
		RegisterTime: time.Now().Unix(),
		Currency:     make(map[string]int64),
		Achievements: make([]string, 0),
		Tags:         make([]string, 0),
		Attributes:   make(map[string]interface{}),
	}
	a.UpdatePlayerProfile(profile)
	return profile
}

// UpdatePlayerGameStats 更新玩家游戏统计
func (a *Analytics) UpdatePlayerGameStats(playerID string, playTime int64, won bool, enemiesKilled int, towersBuilt int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	profile, exists := a.playerCache[playerID]
	if !exists {
		profile = &PlayerProfile{
			PlayerID:     playerID,
			RegisterTime: time.Now().Unix(),
			Currency:     make(map[string]int64),
			Achievements: make([]string, 0),
			Tags:         make([]string, 0),
			Attributes:   make(map[string]interface{}),
		}
	}

	profile.TotalPlayTime += playTime
	profile.TotalGames++
	if won {
		profile.WinGames++
	}
	profile.EnemiesKilled += enemiesKilled
	profile.TowersBuilt += towersBuilt
	profile.LastLoginTime = time.Now().Unix()

	a.playerCache[playerID] = profile
}

// ==================== 漏斗分析 ====================

// FunnelStep 漏斗步骤
type FunnelStep struct {
	Name       string        `json:"name"`        // 步骤名称
	EventType  EventType     `json:"event_type"`  // 事件类型
	Conditions map[string]interface{} `json:"conditions"` // 条件
	TimeWindow time.Duration `json:"time_window"` // 时间窗口
}

// FunnelAnalyzer 漏斗分析器
type FunnelAnalyzer struct {
	steps     []FunnelStep
	analytics *Analytics
}

// NewFunnelAnalyzer 创建漏斗分析器
func NewFunnelAnalyzer(analytics *Analytics, steps []FunnelStep) *FunnelAnalyzer {
	return &FunnelAnalyzer{
		steps:     steps,
		analytics: analytics,
	}
}

// AnalyzeFunnel 分析漏斗
func (f *FunnelAnalyzer) AnalyzeFunnel(playerID string, startTime, endTime int64) *FunnelResult {
	result := &FunnelResult{
		PlayerID: playerID,
		Steps:    make([]FunnelStepResult, len(f.steps)),
	}

	// 简化实现：实际应查询事件存储
	for i, step := range f.steps {
		result.Steps[i] = FunnelStepResult{
			StepName: step.Name,
			Count:    0,
			Rate:     0,
		}
	}

	return result
}

// FunnelResult 漏斗结果
type FunnelResult struct {
	PlayerID string              `json:"player_id"`
	Steps    []FunnelStepResult  `json:"steps"`
	TotalRate float64            `json:"total_rate"`
}

// FunnelStepResult 漏斗步骤结果
type FunnelStepResult struct {
	StepName string  `json:"step_name"`
	Count    int     `json:"count"`
	Rate     float64 `json:"rate"`
}

// ==================== 留存分析 ====================

// RetentionAnalyzer 留存分析器
type RetentionAnalyzer struct {
	analytics *Analytics
}

// NewRetentionAnalyzer 创建留存分析器
func NewRetentionAnalyzer(analytics *Analytics) *RetentionAnalyzer {
	return &RetentionAnalyzer{analytics: analytics}
}

// CalculateRetention 计算留存率
func (r *RetentionAnalyzer) CalculateRetention(registerTime int64, currentTime int64, retentionDays []int) map[int]float64 {
	result := make(map[int]float64)

	for _, day := range retentionDays {
		// 简化实现
		result[day] = 0.0
	}

	return result
}

// GetDailyRetention 获取每日留存
func (r *RetentionAnalyzer) GetDailyRetention(date time.Time) *RetentionData {
	return &RetentionData{
		Date:         date,
		NewUsers:     0,
		Day1Retention: 0.0,
		Day7Retention: 0.0,
		Day30Retention: 0.0,
	}
}

// RetentionData 留存数据
type RetentionData struct {
	Date          time.Time `json:"date"`
	NewUsers      int       `json:"new_users"`
	Day1Retention float64   `json:"day1_retention"`
	Day7Retention float64   `json:"day7_retention"`
	Day30Retention float64  `json:"day30_retention"`
}

// ==================== 收入分析 ====================

// RevenueAnalyzer 收入分析器
type RevenueAnalyzer struct {
	analytics *Analytics
}

// NewRevenueAnalyzer 创建收入分析器
func NewRevenueAnalyzer(analytics *Analytics) *RevenueAnalyzer {
	return &RevenueAnalyzer{analytics: analytics}
}

// CalculateARPPU 计算ARPPU (每付费用户平均收入)
func (r *RevenueAnalyzer) CalculateARPPU(startTime, endTime int64) float64 {
	// 简化实现
	return 0.0
}

// CalculateLTV 计算LTV (生命周期价值)
func (r *RevenueAnalyzer) CalculateLTV(playerID string) float64 {
	// 简化实现
	return 0.0
}

// GetDailyRevenue 获取每日收入
func (r *RevenueAnalyzer) GetDailyRevenue(date time.Time) *RevenueData {
	return &RevenueData{
		Date:         date,
		Revenue:      0,
		Orders:       0,
		PaidUsers:    0,
		ARPPU:        0,
		ARPDAU:       0,
	}
}

// RevenueData 收入数据
type RevenueData struct {
	Date      time.Time `json:"date"`
	Revenue   float64   `json:"revenue"`
	Orders    int       `json:"orders"`
	PaidUsers int       `json:"paid_users"`
	ARPPU     float64   `json:"arpu"`
	ARPDAU    float64   `json:"arpdau"`
}

// ==================== 实时统计 ====================

// RealtimeStats 实时统计
type RealtimeStats struct {
	OnlinePlayers    int     `json:"online_players"`     // 在线玩家
	ActiveRooms      int     `json:"active_rooms"`       // 活跃房间
	GamesPerMinute   float64 `json:"games_per_minute"`   // 每分钟游戏数
	MessagesPerMinute float64 `json:"messages_per_minute"` // 每分钟消息数
	AvgMatchTime     float64 `json:"avg_match_time"`     // 平均匹配时间(秒)
	AvgGameDuration  float64 `json:"avg_game_duration"` // 平均游戏时长(秒)
}

// GetRealtimeStats 获取实时统计
func (a *Analytics) GetRealtimeStats() *RealtimeStats {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return &RealtimeStats{
		OnlinePlayers:   a.stats.ActivePlayers,
		ActiveRooms:     0,
		GamesPerMinute:  0,
		MessagesPerMinute: 0,
		AvgMatchTime:    0,
		AvgGameDuration: a.stats.AvgPlayTime,
	}
}

// ==================== 报表生成 ====================

// Report 报表
type Report struct {
	ReportType string                 `json:"report_type"` // 报表类型
	DateRange  [2]int64               `json:"date_range"`  // 日期范围
	Metrics    map[string]interface{} `json:"metrics"`     // 指标
	GeneratedAt int64                 `json:"generated_at"` // 生成时间
}

// GenerateReport 生成报表
func (a *Analytics) GenerateReport(reportType string, startTime, endTime int64) *Report {
	report := &Report{
		ReportType: reportType,
		DateRange:  [2]int64{startTime, endTime},
		Metrics:    make(map[string]interface{}),
		GeneratedAt: time.Now().Unix(),
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	switch reportType {
	case "daily":
		report.Metrics["total_events"] = a.stats.TotalEvents
		report.Metrics["active_players"] = a.stats.ActivePlayers
		report.Metrics["today_games"] = a.stats.TodayGames
		report.Metrics["today_revenue"] = a.stats.TodayRevenue
	case "realtime":
		report.Metrics["online_players"] = a.stats.ActivePlayers
		report.Metrics["peak_online"] = a.stats.PeakOnline
	}

	return report
}

// ExportReport 导出报表
func (r *Report) ExportReport(format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(r, "", "  ")
	case "csv":
		// 简化CSV导出
		return []byte("report_type,date_range,metrics\n"), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// ==================== 辅助函数 ====================

// generateUUID 生成UUID
func generateUUID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), generateRandomID())
}

func generateRandomID() int64 {
	return time.Now().UnixNano() % 1000000
}

// GetStats 获取统计信息
func (a *Analytics) GetStats() *AnalyticsStats {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.stats
}

// ResetStats 重置统计
func (a *Analytics) ResetStats() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.stats = &AnalyticsStats{
		EventsByType: make(map[EventType]int64),
	}
}
