package game

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitState 熔断器状态
type CircuitState int

const (
	CircuitStateClosed CircuitState = iota // 关闭状态，正常
	CircuitStateOpen                       // 开启状态，熔断
	CircuitStateHalfOpen                   // 半开状态，尝试恢复
)

var stateNames = map[CircuitState]string{
	CircuitStateClosed:    "closed",
	CircuitStateOpen:      "open",
	CircuitStateHalfOpen: "half-open",
}

// String 转换为字符串
func (s CircuitState) String() string {
	return stateNames[s]
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	name             string
	failureThreshold int           // 失败阈值
	successThreshold int           // 成功阈值
	timeout          time.Duration // 熔断超时
	halfOpenMaxCalls int           // 半开状态最大并发数
	
	mu             sync.RWMutex
	state          CircuitState
	failures       int
	successes      int
	lastFailure    time.Time
	lastStateChange time.Time
	opensAt        time.Time
	
	// 回调函数
	onStateChange func(from, to CircuitState)
	onFailure     func(err error)
	onSuccess     func()
}

// CircuitBreakerOption 熔断器配置选项
type CircuitBreakerOption func(*CircuitBreaker)

// WithFailureThreshold 设置失败阈值
func WithFailureThreshold(threshold int) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.failureThreshold = threshold
	}
}

// WithSuccessThreshold 设置成功阈值
func WithSuccessThreshold(threshold int) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.successThreshold = threshold
	}
}

// WithTimeout 设置熔断超时
func WithTimeout(timeout time.Duration) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.timeout = timeout
	}
}

// WithHalfOpenMaxCalls 设置半开状态最大并发数
func WithHalfOpenMaxCalls(maxCalls int) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.halfOpenMaxCalls = maxCalls
	}
}

// WithOnStateChange 设置状态变更回调
func WithOnStateChange(fn func(from, to CircuitState)) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.onStateChange = fn
	}
}

// WithOnFailure 设置失败回调
func WithOnFailure(fn func(err error)) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.onFailure = fn
	}
}

// WithOnSuccess 设置成功回调
func WithOnSuccess(fn func()) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.onSuccess = fn
	}
}

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(name string, opts ...CircuitBreakerOption) *CircuitBreaker {
	cb := &CircuitBreaker{
		name:             name,
		failureThreshold: 5,
		successThreshold: 3,
		timeout:          60 * time.Second,
		halfOpenMaxCalls: 3,
		state:            CircuitStateClosed,
		lastStateChange:  time.Now(),
	}
	
	for _, opt := range opts {
		opt(cb)
	}
	
	return cb
}

// Execute 执行函数
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	// 检查状态
	state := cb.getState()
	
	switch state {
	case CircuitStateOpen:
		// 检查是否超时可以进入半开状态
		if time.Since(cb.opensAt) > cb.timeout {
			cb.setState(CircuitStateHalfOpen)
		} else {
			return errors.New("circuit breaker is open")
		}
	case CircuitStateHalfOpen:
		// 半开状态限制并发
		return errors.New("circuit breaker is half-open, max calls reached")
	}
	
	// 执行函数
	err := fn()
	
	if err != nil {
		cb.recordFailure(err)
	} else {
		cb.recordSuccess()
	}
	
	return err
}

// ExecuteWithRetry 执行函数，带重试
func (cb *CircuitBreaker) ExecuteWithRetry(ctx context.Context, fn func() error, maxRetries int) (err error) {
	for i := 0; i <= maxRetries; i++ {
		err = cb.Execute(ctx, fn)
		if err == nil {
			return nil
		}
		
		// 检查是否是熔断错误
		if errors.Is(err, errors.New("circuit breaker is open")) ||
			errors.Is(err, errors.New("circuit breaker is half-open, max calls reached")) {
			return err
		}
		
		// 等待后重试
		if i < maxRetries {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}
	
	return err
}

// recordFailure 记录失败
func (cb *CircuitBreaker) recordFailure(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.failures++
	cb.lastFailure = time.Now()
	
	if cb.onFailure != nil {
		cb.onFailure(err)
	}
	
	// 检查是否需要打开熔断器
	if cb.state == CircuitStateClosed && cb.failures >= cb.failureThreshold {
		cb.setStateLocked(CircuitStateOpen)
	} else if cb.state == CircuitStateHalfOpen {
		// 半开状态失败，重新打开
		cb.setStateLocked(CircuitStateOpen)
	}
}

// recordSuccess 记录成功
func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.successes++
	
	if cb.onSuccess != nil {
		cb.onSuccess()
	}
	
	// 检查是否需要关闭熔断器
	if cb.state == CircuitStateHalfOpen && cb.successes >= cb.successThreshold {
		cb.setStateLocked(CircuitStateClosed)
	}
}

// getState 获取当前状态
func (cb *CircuitBreaker) getState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	
	// 检查超时
	if cb.state == CircuitStateOpen && time.Since(cb.opensAt) > cb.timeout {
		return CircuitStateHalfOpen
	}
	
	return cb.state
}

// setState 设置状态
func (cb *CircuitBreaker) setState(state CircuitState) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.setStateLocked(state)
}

// setStateLocked 设置状态（需要持有锁）
func (cb *CircuitBreaker) setStateLocked(state CircuitState) {
	if cb.state == state {
		return
	}
	
	oldState := cb.state
	cb.state = state
	cb.lastStateChange = time.Now()
	
	// 重置计数器
	if state == CircuitStateClosed {
		cb.failures = 0
		cb.successes = 0
	} else if state == CircuitStateOpen {
		cb.opensAt = time.Now()
	} else if state == CircuitStateHalfOpen {
		cb.successes = 0
	}
	
	// 触发回调
	if cb.onStateChange != nil {
		cb.onStateChange(oldState, state)
	}
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetFailureCount 获取失败计数
func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// GetSuccessCount 获取成功计数
func (cb *CircuitBreaker) GetSuccessCount() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.successes
}

// GetLastFailureTime 获取最后失败时间
func (cb *CircuitBreaker) GetLastFailureTime() time.Time {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.lastFailure
}

// GetLastStateChangeTime 获取最后状态变更时间
func (cb *CircuitBreaker) GetLastStateChangeTime() time.Time {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.lastStateChange
}

// Reset 重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.failures = 0
	cb.successes = 0
	cb.state = CircuitStateClosed
	cb.lastStateChange = time.Now()
}

// CircuitBreakerInfo 熔断器信息
type CircuitBreakerInfo struct {
	Name              string        `json:"name"`
	State             string        `json:"state"`
	FailureThreshold int           `json:"failure_threshold"`
	SuccessThreshold int           `json:"success_threshold"`
	Timeout          string        `json:"timeout"`
	Failures          int           `json:"failures"`
	Successes         int           `json:"successes"`
	LastFailure       time.Time     `json:"last_failure"`
	LastStateChange   time.Time     `json:"last_state_change"`
}

// GetInfo 获取熔断器信息
func (cb *CircuitBreaker) GetInfo() CircuitBreakerInfo {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	
	return CircuitBreakerInfo{
		Name:              cb.name,
		State:             cb.state.String(),
		FailureThreshold:  cb.failureThreshold,
		SuccessThreshold: cb.successThreshold,
		Timeout:           cb.timeout.String(),
		Failures:          cb.failures,
		Successes:         cb.successes,
		LastFailure:       cb.lastFailure,
		LastStateChange:   cb.lastStateChange,
	}
}

// String 字符串表示
func (cb *CircuitBreaker) String() string {
	return fmt.Sprintf("CircuitBreaker(%s, state=%s, failures=%d, successes=%d)",
		cb.name, cb.GetState(), cb.GetFailureCount(), cb.GetSuccessCount())
}

// CircuitBreakerManager 熔断器管理器
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu        sync.RWMutex
}

// NewCircuitBreakerManager 创建熔断器管理器
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate 获取或创建熔断器
func (m *CircuitBreakerManager) GetOrCreate(name string, opts ...CircuitBreakerOption) *CircuitBreaker {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if cb, exists := m.breakers[name]; exists {
		return cb
	}
	
	cb := NewCircuitBreaker(name, opts...)
	m.breakers[name] = cb
	
	return cb
}

// Get 获取熔断器
func (m *CircuitBreakerManager) Get(name string) (*CircuitBreaker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if cb, exists := m.breakers[name]; exists {
		return cb, nil
	}
	
	return nil, fmt.Errorf("circuit breaker not found: %s", name)
}

// Remove 移除熔断器
func (m *CircuitBreakerManager) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.breakers, name)
}

// List 列出所有熔断器
func (m *CircuitBreakerManager) List() []*CircuitBreaker {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	breakers := make([]*CircuitBreaker, 0, len(m.breakers))
	for _, cb := range m.breakers {
		breakers = append(breakers, cb)
	}
	
	return breakers
}

// GetAllInfo 获取所有熔断器信息
func (m *CircuitBreakerManager) GetAllInfo() []CircuitBreakerInfo {
	breakers := m.List()
	infos := make([]CircuitBreakerInfo, len(breakers))
	
	for i, cb := range breakers {
		infos[i] = cb.GetInfo()
	}
	
	return infos
}

// ResetAll 重置所有熔断器
func (m *CircuitBreakerManager) ResetAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	for _, cb := range m.breakers {
		cb.Reset()
	}
}
