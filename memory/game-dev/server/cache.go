package game

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// CacheType 缓存类型
type CacheType string

const (
	CacheTypeMemory CacheType = "memory"
	CacheTypeRedis  CacheType = "redis"
)

// CacheItem 缓存项
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
	CreatedAt  time.Time
}

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
	GetMulti(ctx context.Context, keys []string) (map[string]interface{}, error)
	SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	items    map[string]*CacheItem
	mu       sync.RWMutex
	maxSize  int
	evictFn  func(key string, value interface{})
	defaultTTL time.Duration
}

// NewMemoryCache 创建内存缓存
func NewMemoryCache(maxSize int, defaultTTL time.Duration) *MemoryCache {
	cache := &MemoryCache{
		items:       make(map[string]*CacheItem),
		maxSize:     maxSize,
		defaultTTL:  defaultTTL,
	}
	
	// 启动清理goroutine
	go cache.cleanup()
	
	return cache
}

// Get 获取缓存
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	
	// 检查过期
	if time.Now().After(item.Expiration) {
		return nil, fmt.Errorf("key expired: %s", key)
	}
	
	return item.Value, nil
}

// Set 设置缓存
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// 检查容量
	if len(c.items) >= c.maxSize {
		c.evictOldest()
	}
	
	expiration := time.Now().Add(ttl)
	c.items[key] = &CacheItem{
		Value:      value,
		Expiration: expiration,
		CreatedAt:  time.Now(),
	}
	
	return nil
}

// Delete 删除缓存
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.items, key)
	return nil
}

// Clear 清空缓存
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items = make(map[string]*CacheItem)
	return nil
}

// GetMulti 批量获取
func (c *MemoryCache) GetMulti(ctx context.Context, keys []string) (map[string]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	result := make(map[string]interface{})
	now := time.Now()
	
	for _, key := range keys {
		item, exists := c.items[key]
		if exists && now.Before(item.Expiration) {
			result[key] = item.Value
		}
	}
	
	return result, nil
}

// SetMulti 批量设置
func (c *MemoryCache) SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.defaultTTL
	}
	
	c.mu.Lock()
	defer c.mu.Unlock()
	
	expiration := time.Now().Add(ttl)
	
	for key, value := range items {
		// 检查容量
		if len(c.items) >= c.maxSize {
			c.evictOldest()
		}
		
		c.items[key] = &CacheItem{
			Value:      value,
			Expiration: expiration,
			CreatedAt:  time.Now(),
		}
	}
	
	return nil
}

// evictOldest 清除最老的项
func (c *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true
	
	for key, item := range c.items {
		if first || item.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.CreatedAt
			first = false
		}
	}
	
	if oldestKey != "" {
		if c.evictFn != nil {
			c.evictFn(oldestKey, c.items[oldestKey].Value)
		}
		delete(c.items, oldestKey)
	}
}

// cleanup 定期清理过期项
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		
		for key, item := range c.items {
			if now.After(item.Expiration) {
				if c.evictFn != nil {
					c.evictFn(key, item.Value)
				}
				delete(c.items, key)
			}
		}
		
		c.mu.Unlock()
	}
}

// SetEvictCallback 设置淘汰回调
func (c *MemoryCache) SetEvictCallback(fn func(key string, value interface{})) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictFn = fn
}

// CacheStats 缓存统计
type CacheStats struct {
	Items      int64
	Hits       int64
	Misses     int64
	Evictions  int64
	HitRate    float64
}

// GetStats 获取缓存统计
func (c *MemoryCache) GetStats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	stats := CacheStats{
		Items: int64(len(c.items)),
	}
	
	// 计算命中率需要额外的计数器，这里简化返回
	return stats
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client    interface{} // redis.Client
	prefix    string
	encoder   Encoder
	decoder   Decoder
}

// NewRedisCache 创建Redis缓存
func NewRedisCache(client interface{}, prefix string) *RedisCache {
	return &RedisCache{
		client:  client,
		prefix:  prefix,
		encoder: JSONEncoder{},
		decoder: JSONDecoder{},
	}
}

// Get 获取缓存
func (r *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	// Redis实现需要实际的redis客户端
	// 这里提供接口定义
	return nil, fmt.Errorf("redis client not implemented")
}

// Set 设置缓存
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return fmt.Errorf("redis client not implemented")
}

// Delete 删除缓存
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return fmt.Errorf("redis client not implemented")
}

// Clear 清空缓存
func (r *RedisCache) Clear(ctx context.Context) error {
	return fmt.Errorf("redis client not implemented")
}

// GetMulti 批量获取
func (r *RedisCache) GetMulti(ctx context.Context, keys []string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("redis client not implemented")
}

// SetMulti 批量设置
func (r *RedisCache) SetMulti(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	return fmt.Errorf("redis client not implemented")
}

// Encoder 编码器接口
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

// Decoder 解码器接口
type Decoder interface {
	Decode(data []byte, v interface{}) error
}

// JSONEncoder JSON编码器
type JSONEncoder struct{}

// Encode 编码
func (e *JSONEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// JSONDecoder JSON解码器
type JSONDecoder struct{}

// Decode 解码
func (d *JSONDecoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// CacheWrapper 缓存包装器，提供通用缓存功能
type CacheWrapper struct {
	cache Cache
	stats *CacheStats
	mu    sync.Mutex
}

// NewCacheWrapper 创建缓存包装器
func NewCacheWrapper(cache Cache) *CacheWrapper {
	return &CacheWrapper{
		cache: cache,
		stats: &CacheStats{},
	}
}

// Get 获取缓存，带统计
func (w *CacheWrapper) Get(ctx context.Context, key string, loader func() (interface{}, time.Duration, error)) (interface{}, error) {
	// 尝试从缓存获取
	value, err := w.cache.Get(ctx, key)
	if err == nil {
		w.mu.Lock()
		w.stats.Hits++
		w.mu.Unlock()
		return value, nil
	}
	
	// 缓存未命中
	w.mu.Lock()
	w.stats.Misses++
	w.mu.Unlock()
	
	// 调用loader加载数据
	if loader != nil {
		value, ttl, err := loader()
		if err != nil {
			return nil, err
		}
		
		// 存入缓存
		w.cache.Set(ctx, key, value, ttl)
		
		return value, nil
	}
	
	return nil, fmt.Errorf("key not found: %s", key)
}

// Invalidate 失效缓存
func (w *CacheWrapper) Invalidate(ctx context.Context, key string) error {
	return w.cache.Delete(ctx, key)
}

// GetStats 获取统计
func (w *CacheWrapper) GetStats() CacheStats {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	stats := *w.stats
	if stats.Hits+stats.Misses > 0 {
		stats.HitRate = float64(stats.Hits) / float64(stats.Hits+stats.Misses)
	}
	
	return stats
}

// CacheManager 缓存管理器
type CacheManager struct {
	caches    map[string]Cache
	defaultTTL time.Duration
	mu        sync.RWMutex
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(defaultTTL time.Duration) *CacheManager {
	return &CacheManager{
		caches:     make(map[string]Cache),
		defaultTTL: defaultTTL,
	}
}

// RegisterCache 注册缓存
func (m *CacheManager) RegisterCache(name string, cache Cache) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.caches[name] = cache
}

// GetCache 获取缓存
func (m *CacheManager) GetCache(name string) (Cache, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	cache, exists := m.caches[name]
	if !exists {
		return nil, fmt.Errorf("cache not found: %s", name)
	}
	
	return cache, nil
}

// GetOrCreateCache 获取或创建缓存
func (m *CacheManager) GetOrCreateCache(name string, cacheType CacheType, maxSize int) (Cache, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if cache, exists := m.caches[name]; exists {
		return cache, nil
	}
	
	var cache Cache
	switch cacheType {
	case CacheTypeMemory:
		cache = NewMemoryCache(maxSize, m.defaultTTL)
	case CacheTypeRedis:
		return nil, fmt.Errorf("redis cache requires client")
	default:
		cache = NewMemoryCache(maxSize, m.defaultTTL)
	}
	
	m.caches[name] = cache
	return cache, nil
}

// ClearAll 清空所有缓存
func (m *CacheManager) ClearAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for _, cache := range m.caches {
		if err := cache.Clear(ctx); err != nil {
			return err
		}
	}
	
	return nil
}
