package security

import (
	"net"
	"sync"
	"time"
)

// RateLimiter 滑动窗口限流器
type RateLimiter struct {
	mu       sync.RWMutex
	requests map[string]*window // key: client identifier
	config   RateLimitConfig
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	MaxRequests int           // 窗口内最大请求数
	WindowSize  time.Duration // 窗口大小
	BanDuration time.Duration // 超过限制后的封禁时长
	CleanupInterval time.Duration // 清理间隔
}

type window struct {
	requests []time.Time
	banned   bool
	banUntil time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*window),
		config:   config,
	}

	// 启动清理协程
	go rl.cleanup()

	return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) (allowed bool, retryAfter time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	w, exists := rl.requests[key]

	if !exists {
		w = &window{
			requests: make([]time.Time, 0, rl.config.MaxRequests),
		}
		rl.requests[key] = w
	}

	// 检查是否被封禁
	if w.banned && now.Before(w.banUntil) {
		return false, w.banUntil.Sub(now)
	}

	// 解封
	if w.banned && now.After(w.banUntil) {
		w.banned = false
		w.requests = w.requests[:0]
	}

	// 清理过期请求
	rl.cleanWindow(w, now)

	// 检查是否超过限制
	if len(w.requests) >= rl.config.MaxRequests {
		// 超过限制，封禁
		w.banned = true
		w.banUntil = now.Add(rl.config.BanDuration)
		return false, rl.config.BanDuration
	}

	// 记录请求
	w.requests = append(w.requests, now)
	return true, 0
}

// cleanWindow 清理窗口中的过期请求
func (rl *RateLimiter) cleanWindow(w *window, now time.Time) {
	cutoff := now.Add(-rl.config.WindowSize)
	i := 0
	for ; i < len(w.requests); i++ {
		if w.requests[i].After(cutoff) {
			break
		}
	}
	w.requests = w.requests[i:]
}

// cleanup 定期清理不活跃的客户端
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, w := range rl.requests {
			// 清理超过窗口大小的不活跃客户端
			if len(w.requests) == 0 && !w.banned {
				delete(rl.requests, key)
				continue
			}
			// 清理过期的请求
			rl.cleanWindow(w, now)
		}
		rl.mu.Unlock()
	}
}

// Reset 重置某个客户端的限流状态
func (rl *RateLimiter) Reset(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.requests, key)
}

// GetConfig 获取当前配置
func (rl *RateLimiter) GetConfig() RateLimitConfig {
	return rl.config
}

// --- 连接级别限流器 ---

// ConnectionRateLimiter 连接级别限流器
type ConnectionRateLimiter struct {
	mu           sync.RWMutex
	connections  map[string]*connInfo // key: IP
	config       RateLimitConfig
	maxConnPerIP int
}

type connInfo struct {
	count      int
	firstSeen  time.Time
	lastActive time.Time
}

// NewConnectionRateLimiter 创建连接限流器
func NewConnectionRateLimiter(config RateLimitConfig, maxConnPerIP int) *ConnectionRateLimiter {
	rl := &ConnectionRateLimiter{
		connections:  make(map[string]*connInfo),
		config:       config,
		maxConnPerIP: maxConnPerIP,
	}

	go rl.cleanup()
	return rl
}

// AllowConnection 允许新连接
func (rl *ConnectionRateLimiter) AllowConnection(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	info, exists := rl.connections[ip]

	if !exists {
		info = &connInfo{
			firstSeen:  now,
			lastActive: now,
		}
		rl.connections[ip] = info
	}

	// 更新活跃时间
	info.lastActive = now

	// 检查连接数
	if info.count >= rl.maxConnPerIP {
		return false
	}

	info.count++
	return true
}

// ReleaseConnection 释放连接
func (rl *ConnectionRateLimiter) ReleaseConnection(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if info, exists := rl.connections[ip]; exists {
		info.count--
		if info.count < 0 {
			info.count = 0
		}
	}
}

func (rl *ConnectionRateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, info := range rl.connections {
			// 清理长时间不活跃的连接记录
			if now.Sub(info.lastActive) > rl.config.WindowSize*2 {
				delete(rl.connections, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// --- WebSocket 消息限流器 ---

// MessageRateLimiter WebSocket消息限流器
type MessageRateLimiter struct {
	mu           sync.RWMutex
	clientMsgs   map[string]*msgWindow // key: client ID
	config       RateLimitConfig
	banThreshold int // 连续超过限制多少次后封禁
}

type msgWindow struct {
	messages    []time.Time
	violations  int // 连续违规次数
	banned      bool
	banUntil    time.Time
}

// NewMessageRateLimiter 创建消息限流器
func NewMessageRateLimiter(config RateLimitConfig, banThreshold int) *MessageRateLimiter {
	return &MessageRateLimiter{
		clientMsgs:   make(map[string]*msgWindow),
		config:       config,
		banThreshold: banThreshold,
	}
}

// AllowMessage 检查是否允许消息
func (rl *MessageRateLimiter) AllowMessage(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	w, exists := rl.clientMsgs[clientID]

	if !exists {
		w = &msgWindow{
			messages: make([]time.Time, 0, rl.config.MaxRequests),
		}
		rl.clientMsgs[clientID] = w
	}

	// 检查封禁
	if w.banned && now.Before(w.banUntil) {
		return false
	}

	// 解封
	if w.banned && now.After(w.banUntil) {
		w.banned = false
		w.violations = 0
		w.messages = w.messages[:0]
	}

	// 清理过期消息
	rl.cleanMessages(w, now)

	// 检查限制
	if len(w.messages) >= rl.config.MaxRequests {
		w.violations++
		if w.violations >= rl.banThreshold {
			w.banned = true
			w.banUntil = now.Add(rl.config.BanDuration)
		}
		return false
	}

	w.violations = 0
	w.messages = append(w.messages, now)
	return true
}

func (rl *MessageRateLimiter) cleanMessages(w *msgWindow, now time.Time) {
	cutoff := now.Add(-rl.config.WindowSize)
	i := 0
	for ; i < len(w.messages); i++ {
		if w.messages[i].After(cutoff) {
			break
		}
	}
	w.messages = w.messages[i:]
}

// --- IP 黑名单 ---

// IPBlackList IP黑名单
type IPBlackList struct {
	mu       sync.RWMutex
	blocked  map[string]time.Time // IP -> 过期时间
	duration time.Duration
}

// NewIPBlackList 创建IP黑名单
func NewIPBlackList(duration time.Duration) *IPBlackList {
	return &IPBlackList{
		blocked:  make(map[string]time.Time),
		duration: duration,
	}
}

// Block 封禁IP
func (bl *IPBlackList) Block(ip string) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	bl.blocked[ip] = time.Now().Add(bl.duration)
}

// BlockPermanent 永久封禁IP
func (bl *IPBlackList) BlockPermanent(ip string) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	bl.blocked[ip] = time.Now().Add(365 * 24 * time.Hour) // 1年
}

// Unblock 解封IP
func (bl *IPBlackList) Unblock(ip string) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	delete(bl.blocked, ip)
}

// IsBlocked 检查IP是否被封禁
func (bl *IPBlackList) IsBlocked(ip string) bool {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	if expire, exists := bl.blocked[ip]; exists {
		if time.Now().Before(expire) {
			return true
		}
		// 自动清理过期
		delete(bl.blocked, ip)
	}
	return false
}

// GetBlockedIPs 获取所有被封禁的IP
func (bl *IPBlackList) GetBlockedIPs() []string {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	ips := make([]string, 0, len(bl.blocked))
	now := time.Now()
	for ip, expire := range bl.blocked {
		if now.Before(expire) {
			ips = append(ips, ip)
		}
	}
	return ips
}

// --- 验证工具 ---

// ValidateIP 验证IP地址
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// GetClientIP 获取客户端真实IP
func GetClientIP(realIP, remoteAddr string) string {
	if realIP != "" {
		// 检查是否包含端口
		if idx := indexAny(realIP, ":"); idx > 0 {
			return realIP[:idx]
		}
		return realIP
	}
	// 从 remoteAddr 提取
	if idx := indexAny(remoteAddr, ":"); idx > 0 {
		return remoteAddr[:idx]
	}
	return remoteAddr
}

func indexAny(s, chars string) int {
	for i, c := range s {
		for _, cc := range chars {
			if c == cc {
				return i
			}
		}
	}
	return -1
}
