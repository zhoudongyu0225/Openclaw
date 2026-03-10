package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// AlertLevel 告警级别
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelError    AlertLevel = "error"
	AlertLevelCritical AlertLevel = "critical"
)

// Alert 告警结构
type Alert struct {
	ID          string                 `json:"id"`
	Level       AlertLevel             `json:"level"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Source      string                 `json:"source"`
	Timestamp   time.Time              `json:"timestamp"`
	Tags        map[string]string      `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	Acknowledged bool                  `json:"acknowledged"`
	Resolved    bool                   `json:"resolved"`
}

// AlertRule 告警规则
type AlertRule struct {
	Name        string                 `json:"name"`
	Condition   AlertCondition         `json:"condition"`
	Level       AlertLevel             `json:"level"`
	Message     string                 `json:"message"`
	Cooldown    time.Duration          `json:"cooldown"`
	Enabled     bool                   `json:"enabled"`
	Actions     []AlertAction          `json:"actions"`
}

// AlertCondition 告警条件
type AlertCondition struct {
	Metric      string   `json:"metric"`
	Operator    string   `json:"operator"` // gt, lt, eq, gte, lte
	Threshold   float64  `json:"threshold"`
	Duration    time.Duration `json:"duration"`
	Window      time.Duration `json:"window"`
}

// AlertAction 告警动作
type AlertAction struct {
	Type     string `json:"type"` // webhook, email, sms, slack
	Endpoint string `json:"endpoint"`
	Template string `json:"template"`
}

// AlertManager 告警管理器
type AlertManager struct {
	rules      map[string]*AlertRule
	alerts     map[string]*Alert
	history    []*Alert
	handlers   []AlertHandler
	mu         sync.RWMutex
	config     *AlertConfig
	lastAlert  map[string]time.Time
}

// AlertConfig 告警配置
type AlertConfig struct {
	Enabled        bool          `json:"enabled"`
	MaxAlerts      int           `json:"max_alerts"`       // 最大告警数
	RetentionDays  int           `json:"retention_days"`   // 保留天数
	CooldownPeriod time.Duration `json:"cooldown_period"` // 冷却期
}

// AlertHandler 告警处理器接口
type AlertHandler interface {
	Handle(alert *Alert) error
	Type() string
}

// WebhookHandler Webhook处理器
type WebhookHandler struct {
	endpoint string
	client   *HTTPClient
}

// EmailHandler 邮件处理器
type EmailHandler struct {
	smtpHost string
	smtpPort int
	from     string
	to       []string
}

// SlackHandler Slack处理器
type SlackHandler struct {
	webhookURL string
	channel    string
	username   string
}

// NewAlertManager 创建告警管理器
func NewAlertManager(config *AlertConfig) *AlertManager {
	if config == nil {
		config = &AlertConfig{
			Enabled:        true,
			MaxAlerts:      1000,
			RetentionDays:  30,
			CooldownPeriod: 5 * time.Minute,
		}
	}

	return &AlertManager{
		rules:     make(map[string]*AlertRule),
		alerts:    make(map[string]*Alert),
		history:   make([]*Alert, 0),
		handlers:  make([]AlertHandler, 0),
		config:    config,
		lastAlert: make(map[string]time.Time),
	}
}

// RegisterRule 注册告警规则
func (am *AlertManager) RegisterRule(rule *AlertRule) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if rule.Name == "" {
		return fmt.Errorf("rule name cannot be empty")
	}

	am.rules[rule.Name] = rule
	return nil
}

// UnregisterRule 注销告警规则
func (am *AlertManager) UnregisterRule(name string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	delete(am.rules, name)
}

// RegisterHandler 注册告警处理器
func (am *AlertManager) RegisterHandler(handler AlertHandler) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.handlers = append(am.handlers, handler)
}

// Evaluate 评估告警条件
func (am *AlertManager) Evaluate(metrics map[string]float64) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	for _, rule := range am.rules {
		if !rule.Enabled {
			continue
		}

		if am.shouldAlert(rule.Name) {
			if am.checkCondition(rule.Condition, metrics) {
				am.triggerAlert(rule, metrics)
			}
		}
	}
}

// checkCondition 检查条件
func (am *AlertManager) checkCondition(condition AlertCondition, metrics map[string]float64) bool {
	value, ok := metrics[condition.Metric]
	if !ok {
		return false
	}

	switch condition.Operator {
	case "gt":
		return value > condition.Threshold
	case "lt":
		return value < condition.Threshold
	case "eq":
		return value == condition.Threshold
	case "gte":
		return value >= condition.Threshold
	case "lte":
		return value <= condition.Threshold
	}

	return false
}

// shouldAlert 判断是否应该告警
func (am *AlertManager) shouldAlert(ruleName string) bool {
	lastTime, ok := am.lastAlert[ruleName]
	if !ok {
		return true
	}

	return time.Since(lastTime) > am.config.CooldownPeriod
}

// triggerAlert 触发告警
func (am *AlertManager) triggerAlert(rule *AlertRule, metrics map[string]float64) {
	alert := &Alert{
		ID:        generateAlertID(),
		Level:     rule.Level,
		Title:     rule.Name,
		Message:   rule.Message,
		Source:    "monitoring",
		Timestamp: time.Now(),
		Tags:      make(map[string]string),
		Metadata:  make(map[string]interface{}),
	}

	// 填充指标数据
	for k, v := range metrics {
		alert.Metadata[k] = v
	}

	am.mu.Lock()
	am.alerts[alert.ID] = alert
	am.history = append(am.history, alert)
	am.lastAlert[rule.Name] = time.Now()

	// 限制告警数量
	if len(am.alerts) > am.config.MaxAlerts {
		am.cleanupOldAlerts()
	}
	am.mu.Unlock()

	// 触发动作
	for _, action := range rule.Actions {
		am.executeAction(action, alert)
	}

	// 调用处理器
	for _, handler := range am.handlers {
		go handler.Handle(alert)
	}
}

// executeAction 执行告警动作
func (am *AlertManager) executeAction(action AlertAction, alert *Alert) {
	switch action.Type {
	case "webhook":
		am.sendWebhook(action.Endpoint, alert)
	case "email":
		am.sendEmail(action.Endpoint, alert)
	case "slack":
		am.sendSlack(action.Endpoint, alert)
	}
}

func (am *AlertManager) sendWebhook(endpoint string, alert *Alert) {
	// 实现 webhook 发送
	fmt.Printf("Sending webhook to %s: %s\n", endpoint, alert.Title)
}

func (am *AlertManager) sendEmail(to string, alert *Alert) {
	// 实现邮件发送
	fmt.Printf("Sending email to %s: %s\n", to, alert.Title)
}

func (am *AlertHandler) sendSlack(webhookURL string, alert *Alert) {
	// 实现 Slack 发送
}

// GetActiveAlerts 获取活跃告警
func (am *AlertManager) GetActiveAlerts() []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var active []*Alert
	for _, alert := range am.alerts {
		if !alert.Resolved {
			active = append(active, alert)
		}
	}

	return active
}

// GetAlertHistory 获取告警历史
func (am *AlertManager) GetAlertHistory(limit int) []*Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if limit > len(am.history) {
		limit = len(am.history)
	}

	result := make([]*Alert, limit)
	copy(result, am.history[len(am.history)-limit:])

	return result
}

// Acknowledge 确认告警
func (am *AlertManager) Acknowledge(alertID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, ok := am.alerts[alertID]
	if !ok {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	alert.Acknowledged = true
	return nil
}

// Resolve 解决告警
func (am *AlertManager) Resolve(alertID string, message string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	alert, ok := am.alerts[alertID]
	if !ok {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	alert.Resolved = true
	if alert.Metadata == nil {
		alert.Metadata = make(map[string]interface{})
	}
	alert.Metadata["resolved_at"] = time.Now()
	alert.Metadata["resolution_message"] = message

	return nil
}

// cleanupOldAlerts 清理旧告警
func (am *AlertManager) cleanupOldAlerts() {
	cutoff := time.Now().AddDate(0, 0, -am.config.RetentionDays)

	var active []*Alert
	for _, alert := range am.alerts {
		if !alert.Resolved || alert.Timestamp.After(cutoff) {
			active = append(active, alert)
		}
	}

	am.alerts = make(map[string]*Alert)
	for _, alert := range active {
		am.alerts[alert.ID] = alert
	}
}

// GetStats 获取告警统计
func (am *AlertManager) GetStats() map[string]interface{} {
	am.mu.RLock()
	defer am.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["total"] = len(am.history)
	stats["active"] = 0
	stats["acknowledged"] = 0
	stats["resolved"] = 0
	stats["by_level"] = make(map[string]int)

	for _, alert := range am.alerts {
		if !alert.Resolved {
			stats["active"] = (stats["active"].(int)) + 1
		}
		if alert.Acknowledged {
			stats["acknowledged"] = (stats["acknowledged"].(int)) + 1
		}
		if alert.Resolved {
			stats["resolved"] = (stats["resolved"].(int)) + 1
		}

		levelCount := stats["by_level"].(map[string]int)
		levelCount[string(alert.Level)]++
		stats["by_level"] = levelCount
	}

	return stats
}

// Handle 实现 AlertHandler 接口
func (h *WebhookHandler) Handle(alert *Alert) error {
	return h.send(alert)
}

func (h *WebhookHandler) Type() string {
	return "webhook"
}

func (h *WebhookHandler) send(alert *Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	fmt.Printf("Webhook sent to %s: %s\n", h.endpoint, string(data))
	return nil
}

// Handle 实现 AlertHandler 接口
func (h *EmailHandler) Handle(alert *Alert) error {
	return h.send(alert)
}

func (h *EmailHandler) Type() string {
	return "email"
}

func (h *EmailHandler) send(alert *Alert) error {
	fmt.Printf("Email sent to %v: %s\n", h.to, alert.Title)
	return nil
}

// Handle 实现 AlertHandler 接口
func (h *SlackHandler) Handle(alert *Alert) error {
	return h.send(alert)
}

func (h *SlackHandler) Type() string {
	return "slack"
}

func (h *SlackHandler) send(alert *Alert) error {
	message := fmt.Sprintf("[%s] %s: %s", alert.Level, alert.Title, alert.Message)
	fmt.Printf("Slack message sent to %s: %s\n", h.channel, message)
	return nil
}

// NewWebhookHandler 创建 Webhook 处理器
func NewWebhookHandler(endpoint string) *WebhookHandler {
	return &WebhookHandler{
		endpoint: endpoint,
		client:   NewHTTPClient(10 * time.Second),
	}
}

// NewEmailHandler 创建邮件处理器
func NewEmailHandler(smtpHost string, smtpPort int, from string, to []string) *EmailHandler {
	return &EmailHandler{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		from:     from,
		to:       to,
	}
}

// NewSlackHandler 创建 Slack 处理器
func NewSlackHandler(webhookURL string, channel string) *SlackHandler {
	return &SlackHandler{
		webhookURL: webhookURL,
		channel:    channel,
		username:   "Game Server Alert",
	}
}

// HTTPClient 简化的 HTTP 客户端
type HTTPClient struct {
	timeout time.Duration
}

// NewHTTPClient 创建 HTTP 客户端
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{timeout: timeout}
}

// Post 发送 POST 请求
func (c *HTTPClient) Post(url string, data []byte) error {
	fmt.Printf("POST %s: %s\n", url, string(data))
	return nil
}

// generateAlertID 生成告警ID
func generateAlertID() string {
	return fmt.Sprintf("alert_%d_%d", time.Now().Unix(), rand.Intn(10000))
}

// AlertMiddleware 创建告警中间件
func AlertMiddleware(am *AlertManager, metricsFunc func() map[string]float64) func() {
	return func() {
		metrics := metricsFunc()
		am.Evaluate(metrics)
	}
}
