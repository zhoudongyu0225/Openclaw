package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// GrayReleaseConfig 灰度发布配置
type GrayReleaseConfig struct {
	Enabled        bool     `json:"enabled"`
	Percentage     int      `json:"percentage"`      // 灰度比例 0-100
	Whitelist      []string `json:"whitelist"`        // 白名单用户
	Blacklist      []string `json:"blacklist"`       // 黑名单用户
	Rules          []Rule   `json:"rules"`            // 自定义规则
	RolloutHistory []RolloutRecord `json:"history"`   //  rollout 历史
	mu             sync.RWMutex
}

// Rule 灰度规则
type Rule struct {
	Name       string            `json:"name"`
	Conditions []Condition       `json:"conditions"`
	Action     RuleAction        `json:"action"`
	Priority   int               `json:"priority"`
}

// Condition 规则条件
type Condition struct {
	Field    string   `json:"field"`    // 字段名
	Operator string   `json:"operator"` // 操作符: eq, neq, in, not_in, regex
	Value    string   `json:"value"`    // 值
}

// RuleAction 规则动作
type RuleAction struct {
	Type  string `json:"type"`  // include, exclude, redirect
	Value string `json:"value"` // 目标版本/地址
}

// RolloutRecord 灰度发布记录
type RolloutRecord struct {
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
	Percentage  int       `json:"percentage"`
	UserCount   int       `json:"user_count"`
	Status      string    `json:"status"` // pending, rolling, complete, rollback
	Description string    `json:"description"`
}

// GrayRelease 灰度发布管理器
type GrayRelease struct {
	config     *GrayReleaseConfig
	versions   map[string]*VersionInfo
	analytics *RolloutAnalytics
	mu         sync.RWMutex
}

// VersionInfo 版本信息
type VersionInfo struct {
	Name        string    `json:"name"`
	Features    []string  `json:"features"`
	BugFixes    []string  `json:"bug_fixes"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"` // active, deprecated, rolled_back
	Metrics     VersionMetrics
}

// VersionMetrics 版本指标
type VersionMetrics struct {
	TotalUsers     int            `json:"total_users"`
	ActiveUsers    int            `json:"active_users"`
	ErrorRate      float64        `json:"error_rate"`
	LatencyP50     time.Duration `json:"latency_p50"`
	LatencyP95     time.Duration `json:"latency_p95"`
	LatencyP99     time.Duration `json:"latency_p99"`
	CrashRate      float64        `json:"crash_rate"`
	UserFeedback   map[string]int `json:"user_feedback"` // feedback -> count
}

// RolloutAnalytics 灰度分析
type RolloutAnalytics struct {
	StartTime     time.Time
	EndTime       time.Time
	TargetUsers   int
	ActualUsers   int
	ConversionRate float64
	Issues        []Issue
	Recommendations []string
}

// Issue 问题记录
type Issue struct {
	Timestamp  time.Time
	Severity   string // critical, major, minor
	Type       string
	Description string
	AffectedUsers int
}

// NewGrayRelease 创建灰度发布管理器
func NewGrayRelease(config *GrayReleaseConfig) *GrayRelease {
	return &GrayRelease{
		config:     config,
		versions:   make(map[string]*VersionInfo),
		analytics: &RolloutAnalytics{},
	}
}

// IsUserInGray 判断用户是否在灰度范围
func (g *GrayRelease) IsUserInGray(userID string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// 检查白名单
	for _, id := range g.config.Whitelist {
		if id == userID {
			return true
		}
	}

	// 检查黑名单
	for _, id := range g.config.Blacklist {
		if id == userID {
			return false
		}
	}

	// 检查自定义规则
	for _, rule := range g.config.Rules {
		if g.matchRule(rule, userID) {
			return rule.Action.Type == "include"
		}
	}

	// 百分比灰度
	if g.config.Enabled && g.config.Percentage > 0 {
		hash := g.hashUserID(userID)
		return hash%100 < g.config.Percentage
	}

	return false
}

// hashUserID 用户ID哈希
func (g *GrayRelease) hashUserID(userID string) int {
	hash := 0
	for i, c := range userID {
		hash = hash*31 + int(c)*(i+1)
	}
	return abs(hash)
}

// matchRule 匹配规则
func (g *GrayRelease) matchRule(rule Rule, userID string) bool {
	for _, cond := range rule.Conditions {
		switch cond.Field {
		case "user_id":
			switch cond.Operator {
			case "eq":
				return userID == cond.Value
			case "neq":
				return userID != cond.Value
			case "in":
				return contains(split(cond.Value, ","), userID)
			}
		}
	}
	return false
}

// StartRollout 开始灰度发布
func (g *GrayRelease) StartRollout(version string, percentage int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, ok := g.versions[version]; !ok {
		return fmt.Errorf("version %s not found", version)
	}

	record := RolloutRecord{
		Timestamp:   time.Now(),
		Version:     version,
		Percentage:  percentage,
		UserCount:   g.estimateUserCount(percentage),
		Status:      "rolling",
		Description: fmt.Sprintf("Starting rollout to %d%% users", percentage),
	}

	g.config.RolloutHistory = append(g.config.RolloutHistory, record)
	g.config.Percentage = percentage

	return nil
}

// CompleteRollout 完成灰度发布
func (g *GrayRelease) CompleteRollout(version string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.config.RolloutHistory) > 0 {
		record := &g.config.RolloutHistory[len(g.config.RolloutHistory)-1]
		if record.Version == version {
			record.Status = "complete"
		}
	}

	if v, ok := g.versions[version]; ok {
		v.Status = "active"
	}

	return nil
}

// Rollback 回滚
func (g *GrayRelease) Rollback(version string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.config.RolloutHistory) > 0 {
		record := &g.config.RolloutHistory[len(g.config.RolloutHistory)-1]
		record.Status = "rollback"
		record.Description = fmt.Sprintf("Rolled back from version %s", version)
	}

	if v, ok := g.versions[version]; ok {
		v.Status = "rolled_back"
	}

	return nil
}

// GetAnalytics 获取灰度分析数据
func (g *GrayRelease) GetAnalytics() *RolloutAnalytics {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.analytics
}

// UpdateMetrics 更新版本指标
func (g *GrayRelease) UpdateMetrics(version string, metrics VersionMetrics) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if v, ok := g.versions[version]; ok {
		v.Metrics = metrics
	}
}

// RegisterVersion 注册新版本
func (g *GrayRelease) RegisterVersion(info *VersionInfo) {
	g.mu.Lock()
	defer g.mu.Unlock()

	info.CreatedAt = time.Now()
	info.Status = "active"
	g.versions[info.Name] = info
}

// GetActiveVersion 获取用户应使用的版本
func (g *GrayRelease) GetActiveVersion(userID string, defaultVersion string) string {
	if g.IsUserInGray(userID) {
		// 获取最新活跃版本
		for name, v := range g.versions {
			if v.Status == "active" {
				return name
			}
		}
	}
	return defaultVersion
}

// estimateUserCount 估算用户数
func (g *GrayRelease) estimateUserCount(percentage int) int {
	baseUsers := 10000 // 假设基础用户数
	return baseUsers * percentage / 100
}

// GetConfig 获取配置
func (g *GrayRelease) GetConfig() *GrayReleaseConfig {
	g.mu.RLock()
	defer g.mu.RUnlock()

	configCopy := *g.config
	configCopy.Whitelist = make([]string, len(g.config.Whitelist))
	copy(configCopy.Whitelist, g.config.Whitelist)

	return &configCopy
}

// UpdateConfig 更新配置
func (g *GrayRelease) UpdateConfig(config *GrayReleaseConfig) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if config.Percentage < 0 || config.Percentage > 100 {
		return fmt.Errorf("percentage must be between 0 and 100")
	}

	g.config = config
	return nil
}

// MarshalJSON 自定义JSON序列化
func (g *GrayRelease) MarshalJSON() ([]byte, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return json.Marshal(struct {
		Config   *GrayReleaseConfig      `json:"config"`
		Versions map[string]*VersionInfo `json:"versions"`
	}{
		Config:   g.config,
		Versions: g.versions,
	})
}

// Helper functions
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func split(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GrayReleaseMiddleware 灰度发布中间件
func GrayReleaseMiddleware(gr *GrayRelease, defaultVersion string) func(userID string) string {
	return func(userID string) string {
		return gr.GetActiveVersion(userID, defaultVersion)
	}
}
