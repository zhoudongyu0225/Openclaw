package monitor

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// MetricsCollector 指标收集器
type MetricsCollector struct {
	// 连接指标
	activeConnections int64
	totalConnections  int64
	disconnected      int64

	// 房间指标
	activeRooms    int64
	totalRooms     int64
	roomsCreated   int64
	roomsDestroyed int64

	// 游戏指标
	activeGames      int64
	totalGames       int64
	gamesWon         int64
	gamesLost        int64
	totalDamageDealt int64
	totalEnemiesKilled int64

	// 战斗指标
	totalWavesSpawned   int64
	totalTowersPlaced   int64
	totalTowersUpgraded int64
	totalTowersSold     int64

	// 经济指标
	totalCoinsEarned   int64
	totalCoinsSpent    int64
	totalGiftsReceived int64
	totalGiftValue     int64

	// WebSocket 指标
	totalMessagesSent     int64
	totalMessagesReceived int64
	totalBytesSent        int64
	totalBytesReceived    int64

	// API 指标
	apiRequestsTotal   int64
	apiRequestsSuccess int64
	apiRequestsError   int64

	// 延迟指标 (毫秒)
	avgLatencyMs   int64
	maxLatencyMs   int64
	latencySamples int64

	// 支付指标
	totalPayments    int64
	paymentSuccess   int64
	paymentFailed    int64
	totalRevenue     int64 // 分

	// 弹幕指标
	totalDanmakuSent     int64
	filteredDanmaku     int64

	// 内存
	lastGCPause time.Duration

	mu         sync.RWMutex
	latencies  []int64
	latencyIdx int
}

var (
	defaultCollector *MetricsCollector
	collectorOnce    sync.Once
	startTime        = time.Now()
)

// GetCollector 获取默认指标收集器
func GetCollector() *MetricsCollector {
	collectorOnce.Do(func() {
		defaultCollector = NewCollector()
	})
	return defaultCollector
}

// NewCollector 创建新的指标收集器
func NewCollector() *MetricsCollector {
	return &MetricsCollector{
		latencies: make([]int64, 1000), // 保留最近1000个样本
	}
}

// --- 连接指标 ---

func (m *MetricsCollector) IncActiveConnections() {
	atomic.AddInt64(&m.activeConnections, 1)
	atomic.AddInt64(&m.totalConnections, 1)
}

func (m *MetricsCollector) DecActiveConnections() {
	atomic.AddInt64(&m.activeConnections, -1)
	atomic.AddInt64(&m.disconnected, 1)
}

func (m *MetricsCollector) GetActiveConnections() int64 {
	return atomic.LoadInt64(&m.activeConnections)
}

func (m *MetricsCollector) GetTotalConnections() int64 {
	return atomic.LoadInt64(&m.totalConnections)
}

// --- 房间指标 ---

func (m *MetricsCollector) IncActiveRooms() {
	atomic.AddInt64(&m.activeRooms, 1)
	atomic.AddInt64(&m.totalRooms, 1)
	atomic.AddInt64(&m.roomsCreated, 1)
}

func (m *MetricsCollector) DecActiveRooms() {
	atomic.AddInt64(&m.activeRooms, -1)
	atomic.AddInt64(&m.roomsDestroyed, 1)
}

func (m *MetricsCollector) GetActiveRooms() int64 {
	return atomic.LoadInt64(&m.activeRooms)
}

// --- 游戏指标 ---

func (m *MetricsCollector) IncActiveGames() {
	atomic.AddInt64(&m.activeGames, 1)
	atomic.AddInt64(&m.totalGames, 1)
}

func (m *MetricsCollector) DecActiveGames() {
	atomic.AddInt64(&m.activeGames, -1)
}

func (m *MetricsCollector) IncGamesWon() {
	atomic.AddInt64(&m.gamesWon, 1)
}

func (m *MetricsCollector) IncGamesLost() {
	atomic.AddInt64(&m.gamesLost, 1)
}

func (m *MetricsCollector) AddDamageDealt(damage int64) {
	atomic.AddInt64(&m.totalDamageDealt, damage)
}

func (m *MetricsCollector) IncEnemiesKilled(n int64) {
	atomic.AddInt64(&m.totalEnemiesKilled, n)
}

// --- 战斗指标 ---

func (m *MetricsCollector) IncWavesSpawned(n int64) {
	atomic.AddInt64(&m.totalWavesSpawned, n)
}

func (m *MetricsCollector) IncTowersPlaced() {
	atomic.AddInt64(&m.totalTowersPlaced, 1)
}

func (m *MetricsCollector) IncTowersUpgraded() {
	atomic.AddInt64(&m.totalTowersUpgraded, 1)
}

func (m *MetricsCollector) IncTowersSold() {
	atomic.AddInt64(&m.totalTowersSold, 1)
}

// --- 经济指标 ---

func (m *MetricsCollector) AddCoinsEarned(coins int64) {
	atomic.AddInt64(&m.totalCoinsEarned, coins)
}

func (m *MetricsCollector) AddCoinsSpent(coins int64) {
	atomic.AddInt64(&m.totalCoinsSpent, coins)
}

func (m *MetricsCollector) IncGiftsReceived() {
	atomic.AddInt64(&m.totalGiftsReceived, 1)
}

func (m *MetricsCollector) AddGiftValue(value int64) {
	atomic.AddInt64(&m.totalGiftValue, value)
}

// --- WebSocket 指标 ---

func (m *MetricsCollector) IncMessagesSent(n int64) {
	atomic.AddInt64(&m.totalMessagesSent, n)
}

func (m *MetricsCollector) IncMessagesReceived(n int64) {
	atomic.AddInt64(&m.totalMessagesReceived, n)
}

func (m *MetricsCollector) AddBytesSent(n int64) {
	atomic.AddInt64(&m.totalBytesSent, n)
}

func (m *MetricsCollector) AddBytesReceived(n int64) {
	atomic.AddInt64(&m.totalBytesReceived, n)
}

// --- API 指标 ---

func (m *MetricsCollector) IncAPIRequest() {
	atomic.AddInt64(&m.apiRequestsTotal, 1)
}

func (m *MetricsCollector) IncAPISuccess() {
	atomic.AddInt64(&m.apiRequestsSuccess, 1)
}

func (m *MetricsCollector) IncAPIError() {
	atomic.AddInt64(&m.apiRequestsError, 1)
}

// --- 延迟指标 ---

func (m *MetricsCollector) RecordLatency(latencyMs int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 更新样本
	m.latencies[m.latencyIdx] = latencyMs
	m.latencyIdx = (m.latencyIdx + 1) % len(m.latencies)
	atomic.AddInt64(&m.latencySamples, 1)

	// 更新最大值
	for {
		current := atomic.LoadInt64(&m.maxLatencyMs)
		if latencyMs <= current {
			break
		}
		if atomic.CompareAndSwapInt64(&m.maxLatencyMs, current, latencyMs) {
			break
		}
	}
}

func (m *MetricsCollector) GetAvgLatencyMs() int64 {
	samples := atomic.LoadInt64(&m.latencySamples)
	if samples == 0 {
		return 0
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var sum int64
	count := len(m.latencies)
	if samples < int64(count) {
		count = int(samples)
	}
	for i := 0; i < count; i++ {
		sum += m.latencies[i]
	}
	return sum / int64(count)
}

func (m *MetricsCollector) GetMaxLatencyMs() int64 {
	return atomic.LoadInt64(&m.maxLatencyMs)
}

// --- 支付指标 ---

func (m *MetricsCollector) IncPayment() {
	atomic.AddInt64(&m.totalPayments, 1)
}

func (m *MetricsCollector) IncPaymentSuccess() {
	atomic.AddInt64(&m.paymentSuccess, 1)
}

func (m *MetricsCollector) IncPaymentFailed() {
	atomic.AddInt64(&m.paymentFailed, 1)
}

func (m *MetricsCollector) AddRevenue(cents int64) {
	atomic.AddInt64(&m.totalRevenue, cents)
}

// --- 弹幕指标 ---

func (m *MetricsCollector) IncDanmakuSent() {
	atomic.AddInt64(&m.totalDanmakuSent, 1)
}

func (m *MetricsCollector) IncFilteredDanmaku() {
	atomic.AddInt64(&m.filteredDanmaku, 1)
}

// --- 内存指标 ---

func (m *MetricsCollector) SetLastGCPause(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastGCPause = d
}

// --- 导出指标 ---

// Metrics 导出所有指标
type Metrics struct {
	// 连接
	ActiveConnections int64 `json:"active_connections"`
	TotalConnections  int64 `json:"total_connections"`
	Disconnected      int64 `json:"disconnected"`

	// 房间
	ActiveRooms     int64 `json:"active_rooms"`
	TotalRooms      int64 `json:"total_rooms"`
	RoomsCreated    int64 `json:"rooms_created"`
	RoomsDestroyed  int64 `json:"rooms_destroyed"`

	// 游戏
	ActiveGames      int64 `json:"active_games"`
	TotalGames       int64 `json:"total_games"`
	GamesWon         int64 `json:"games_won"`
	GamesLost        int64 `json:"games_lost"`
	TotalDamageDealt int64 `json:"total_damage_dealt"`
	TotalEnemiesKilled int64 `json:"total_enemies_killed"`

	// 战斗
	TotalWavesSpawned   int64 `json:"total_waves_spawned"`
	TotalTowersPlaced   int64 `json:"total_towers_placed"`
	TotalTowersUpgraded int64 `json:"total_towers_upgraded"`
	TotalTowersSold     int64 `json:"total_towers_sold"`

	// 经济
	TotalCoinsEarned   int64 `json:"total_coins_earned"`
	TotalCoinsSpent    int64 `json:"total_coins_spent"`
	TotalGiftsReceived int64 `json:"total_gifts_received"`
	TotalGiftValue     int64 `json:"total_gift_value"`

	// WebSocket
	TotalMessagesSent     int64 `json:"total_messages_sent"`
	TotalMessagesReceived int64 `json:"total_messages_received"`
	TotalBytesSent        int64 `json:"total_bytes_sent"`
	TotalBytesReceived    int64 `json:"total_bytes_received"`

	// API
	APIRequestsTotal   int64 `json:"api_requests_total"`
	APIRequestsSuccess int64 `json:"api_requests_success"`
	APIRequestsError   int64 `json:"api_requests_error"`

	// 延迟
	AvgLatencyMs int64 `json:"avg_latency_ms"`
	MaxLatencyMs int64 `json:"max_latency_ms"`

	// 支付
	TotalPayments  int64 `json:"total_payments"`
	PaymentSuccess int64 `json:"payment_success"`
	PaymentFailed  int64 `json:"payment_failed"`
	TotalRevenue   int64 `json:"total_revenue"`

	// 弹幕
	TotalDanmakuSent int64 `json:"total_danmaku_sent"`
	FilteredDanmaku int64 `json:"filtered_danmaku"`

	// 系统
	UptimeSeconds int64 `json:"uptime_seconds"`
}

func (m *MetricsCollector) Export() Metrics {
	return Metrics{
		ActiveConnections:   m.GetActiveConnections(),
		TotalConnections:    m.GetTotalConnections(),
		Disconnected:        atomic.LoadInt64(&m.disconnected),
		ActiveRooms:         atomic.LoadInt64(&m.activeRooms),
		TotalRooms:          atomic.LoadInt64(&m.totalRooms),
		RoomsCreated:        atomic.LoadInt64(&m.roomsCreated),
		RoomsDestroyed:      atomic.LoadInt64(&m.roomsDestroyed),
		ActiveGames:         atomic.LoadInt64(&m.activeGames),
		TotalGames:          atomic.LoadInt64(&m.totalGames),
		GamesWon:            atomic.LoadInt64(&m.gamesWon),
		GamesLost:           atomic.LoadInt64(&m.gamesLost),
		TotalDamageDealt:    atomic.LoadInt64(&m.totalDamageDealt),
		TotalEnemiesKilled:  atomic.LoadInt64(&m.totalEnemiesKilled),
		TotalWavesSpawned:    atomic.LoadInt64(&m.totalWavesSpawned),
		TotalTowersPlaced:    atomic.LoadInt64(&m.totalTowersPlaced),
		TotalTowersUpgraded:  atomic.LoadInt64(&m.totalTowersUpgraded),
		TotalTowersSold:      atomic.LoadInt64(&m.totalTowersSold),
		TotalCoinsEarned:    atomic.LoadInt64(&m.totalCoinsEarned),
		TotalCoinsSpent:     atomic.LoadInt64(&m.totalCoinsSpent),
		TotalGiftsReceived:  atomic.LoadInt64(&m.totalGiftsReceived),
		TotalGiftValue:      atomic.LoadInt64(&m.totalGiftValue),
		TotalMessagesSent:   atomic.LoadInt64(&m.totalMessagesSent),
		TotalMessagesReceived: atomic.LoadInt64(&m.totalMessagesReceived),
		TotalBytesSent:      atomic.LoadInt64(&m.totalBytesSent),
		TotalBytesReceived:  atomic.LoadInt64(&m.totalBytesReceived),
		APIRequestsTotal:    atomic.LoadInt64(&m.apiRequestsTotal),
		APIRequestsSuccess:  atomic.LoadInt64(&m.apiRequestsSuccess),
		APIRequestsError:    atomic.LoadInt64(&m.apiRequestsError),
		AvgLatencyMs:        m.GetAvgLatencyMs(),
		MaxLatencyMs:        m.GetMaxLatencyMs(),
		TotalPayments:       atomic.LoadInt64(&m.totalPayments),
		PaymentSuccess:      atomic.LoadInt64(&m.paymentSuccess),
		PaymentFailed:       atomic.LoadInt64(&m.paymentFailed),
		TotalRevenue:        atomic.LoadInt64(&m.totalRevenue),
		TotalDanmakuSent:    atomic.LoadInt64(&m.totalDanmakuSent),
		FilteredDanmaku:     atomic.LoadInt64(&m.filteredDanmaku),
		UptimeSeconds:       int64(time.Since(startTime).Seconds()),
	}
}

// GetMetrics 获取当前指标
func GetMetrics() Metrics {
	return GetCollector().Export()
}

// JSON 输出 JSON 格式
func (m Metrics) JSON() string {
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}

// PrometheusFormat 输出 Prometheus 格式
func (m Metrics) PrometheusFormat() string {
	return fmt.Sprintf(`# HELP active_connections Current active WebSocket connections
# TYPE active_connections gauge
active_connections %d
# HELP total_games Total games played
# TYPE total_games counter
total_games %d
# HELP games_won Total games won
# TYPE games_won counter
games_won %d
# HELP games_lost Total games lost
# TYPE games_lost counter
games_lost %d
# HELP total_damage_dealt Total damage dealt
# TYPE total_damage_dealt counter
total_damage_dealt %d
# HELP total_enemies_killed Total enemies killed
# TYPE total_enemies_killed counter
total_enemies_killed %d
# HELP total_revenue Total revenue in cents
# TYPE total_revenue counter
total_revenue %d
# HELP total_messages_sent Total WebSocket messages sent
# TYPE total_messages_sent counter
total_messages_sent %d
# HELP api_requests_success Total successful API requests
# TYPE api_requests_success counter
api_requests_success %d
# HELP avg_latency_ms Average latency in milliseconds
# TYPE avg_latency_ms gauge
avg_latency_ms %d
# HELP uptime_seconds Server uptime in seconds
# TYPE uptime_seconds gauge
uptime_seconds %d
`, m.ActiveConnections, m.TotalGames, m.GamesWon, m.GamesLost,
		m.TotalDamageDealt, m.TotalEnemiesKilled, m.TotalRevenue,
		m.TotalMessagesSent, m.APIRequestsSuccess, m.AvgLatencyMs, m.UptimeSeconds)
}
