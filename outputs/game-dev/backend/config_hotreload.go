// Package config provides hot-reloadable configuration management
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server ServerConfig `json:"server"`
	
	// Game configuration  
	Game GameConfig `json:"game"`
	
	// Database configuration
	Database DatabaseConfig `json:"database"`
	
	// Redis configuration
	Redis RedisConfig `json:"redis"`
	
	// Security configuration
	Security SecurityConfig `json:"security"`
	
	// Monitoring configuration
	Monitor MonitorConfig `json:"monitor"`
	
	// Game server configuration
	GameServer GameServerConfig `json:"game_server"`
	
	// Version info
	Version string `json:"version"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	IdleTimeout  int    `json:"idle_timeout"`
	MaxConns     int    `json:"max_conns"`
	CertFile     string `json:"cert_file"`
	KeyFile      string `json:"key_file"`
}

// GameConfig holds game-specific configuration
type GameConfig struct {
	MaxPlayersPerRoom int           `json:"max_players_per_room"`
	RoomTimeout       int           `json:"room_timeout"`
	MatchTimeout      int           `json:"match_timeout"`
	FPS               int           `json:"fps"`
	FrameSkip         int           `json:"frame_skip"`
	BattleConfig      BattleConfig  `json:"battle_config"`
	DanmakuConfig     DanmakuConfig `json:"danmaku_config"`
}

// BattleConfig holds battle system configuration
type BattleConfig struct {
	MaxTowers       int     `json:"max_towers"`
	MaxEnemies      int     `json:"max_enemies"`
	MaxProjectiles  int     `json:"max_projectiles"`
	SpawnInterval   int     `json:"spawn_interval"`
	WaveInterval    int     `json:"wave_interval"`
	GoldPerKill     int     `json:"gold_per_kill"`
	InitialGold     int     `json:"initial_gold"`
	BaseHP          int     `json:"base_hp"`
}

// DanmakuConfig holds danmaku system configuration
type DanmakuConfig struct {
	MaxMessagesPerSec int `json:"max_messages_per_sec"`
	MaxLength          int `json:"max_length"`
	Cooldown           int `json:"cooldown"`
	RateLimit          int `json:"rate_limit"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URI        string `json:"uri"`
	Database   string `json:"database"`
	PoolSize   int    `json:"pool_size"`
	Timeout    int    `json:"timeout"`
	RetryCount int    `json:"retry_count"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Addr         string `json:"addr"`
	Password     string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
	DialTimeout  int    `json:"dial_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	RateLimitReqPerSec int  `json:"rate_limit_req_per_sec"`
	EnableIP whitelist bool `json:"enable_ip_whitelist"`
	IPWhitelist        []string `json:"ip_whitelist"`
	EnableCaptcha      bool `json:"enable_captcha"`
	SessionTimeout     int  `json:"session_timeout"`
	MaxLoginAttempts  int  `json:"max_login_attempts"`
}

// MonitorConfig holds monitoring configuration
type MonitorConfig struct {
	EnableMetrics    bool `json:"enable_metrics"`
	EnableTracing    bool `json:"enable_tracing"`
	EnableProfiling  bool `json:"enable_profiling"`
	MetricsPort      int  `json:"metrics_port"`
	TraceEndpoint    string `json:"trace_endpoint"`
	ProfilePort      int  `json:"profile_port"`
	ReportInterval   int  `json:"report_interval"`
}

// GameServerConfig holds game server specific configuration
type GameServerConfig struct {
	HeartbeatInterval int    `json:"heartbeat_interval"`
	MaxFrameDelay     int    `json:"max_frame_delay"`
	ReconnectTimeout  int    `json:"reconnect_timeout"`
	SaveReplay        bool   `json:"save_replay"`
	ReplayPath        string `json:"replay_path"`
	AIMode            string `json:"ai_mode"`
}

// HotReloader manages configuration hot-reload
type HotReloader struct {
	configPath string
	config     *Config
	mu         sync.RWMutex
	watcher    *fileWatcher
	callbacks  []ConfigChangeCallback
	onChange   chan *Config
	stop       chan struct{}
}

// ConfigChangeCallback is called when configuration changes
type ConfigChangeCallback func(oldCfg, newCfg *Config)

// NewHotReloader creates a new configuration hot-reloader
func NewHotReloader(configPath string) (*HotReloader, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load initial config: %w", err)
	}
	
	return &HotReloader{
		configPath: configPath,
		config:     cfg,
		onChange:   make(chan *Config, 1),
		stop:       make(chan struct{}),
	}, nil
}

// LoadConfig loads configuration from file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	
	// Set defaults
	setDefaults(&cfg)
	
	return &cfg, nil
}

// setDefaults sets default values for missing configuration
func setDefaults(cfg *Config) {
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8888
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30
	}
	if cfg.Game.FPS == 0 {
		cfg.Game.FPS = 60
	}
	if cfg.Game.MaxPlayersPerRoom == 0 {
		cfg.Game.MaxPlayersPerRoom = 4
	}
	if cfg.Database.PoolSize == 0 {
		cfg.Database.PoolSize = 100
	}
	if cfg.Redis.PoolSize == 0 {
		cfg.Redis.PoolSize = 50
	}
	if cfg.Monitor.MetricsPort == 0 {
		cfg.Monitor.MetricsPort = 9090
	}
}

// Get returns the current configuration (thread-safe)
func (h *HotReloader) Get() *Config {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config
}

// GetServer returns server configuration
func (h *HotReloader) GetServer() ServerConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config.Server
}

// GetGame returns game configuration
func (h *HotReloader) GetGame() GameConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config.Game
}

// GetDatabase returns database configuration
func (h *HotReloader) GetDatabase() DatabaseConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config.Database
}

// GetRedis returns Redis configuration
func (h *HotReloader) GetRedis() RedisConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config.Redis
}

// GetSecurity returns security configuration
func (h *HotReloader) GetSecurity() SecurityConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config.Security
}

// GetMonitor returns monitor configuration
func (h *HotReloader) GetMonitor() MonitorConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config.Monitor
}

// GetGameServer returns game server configuration
func (h *HotReloader) GetGameServer() GameServerConfig {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.config.GameServer
}

// StartWatcher starts watching config file for changes
func (h *HotReloader) StartWatcher(interval time.Duration) error {
	watcher, err := newFileWatcher(h.configPath)
	if err != nil {
		return err
	}
	
	h.watcher = watcher
	
	go func() {
		defer watcher.Close()
		
		for {
			select {
			case <-h.stop:
				return
			case <-watcher.changed:
				h.reloadConfig()
			case <-time.After(interval):
				h.reloadConfig()
			}
		}
	}()
	
	return nil
}

// reloadConfig reloads configuration from file
func (h *HotReloader) reloadConfig() {
	newCfg, err := LoadConfig(h.configPath)
	if err != nil {
		fmt.Printf("Failed to reload config: %v\n", err)
		return
	}
	
	oldCfg := h.Get()
	
	h.mu.Lock()
	h.config = newCfg
	h.mu.Unlock()
	
	// Notify callbacks
	for _, cb := range h.callbacks {
		cb(oldCfg, newCfg)
	}
	
	// Send to channel
	select {
	case h.onChange <- newCfg:
	default:
	}
	
	fmt.Printf("Configuration reloaded: version=%s\n", newCfg.Version)
}

// OnChange registers a callback for configuration changes
func (h *HotReloader) OnChange(callback ConfigChangeCallback) {
	h.callbacks = append(h.callbacks, callback)
}

// OnChangeChannel returns a channel that receives config changes
func (h *HotReloader) OnChangeChannel() <-chan *Config {
	return h.onChange
}

// Stop stops the config watcher
func (h *HotReloader) Stop() {
	close(h.stop)
}

// Validate validates the configuration
func (h *HotReloader) Validate() error {
	cfg := h.Get()
	
	// Validate server config
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	
	// Validate game config
	if cfg.Game.MaxPlayersPerRoom <= 0 || cfg.Game.MaxPlayersPerRoom > 100 {
		return fmt.Errorf("invalid max players per room: %d", cfg.Game.MaxPlayersPerRoom)
	}
	if cfg.Game.FPS <= 0 || cfg.Game.FPS > 144 {
		return fmt.Errorf("invalid FPS: %d", cfg.Game.FPS)
	}
	
	// Validate database config
	if cfg.Database.URI == "" {
		return fmt.Errorf("database URI is required")
	}
	
	// Validate Redis config
	if cfg.Redis.Addr == "" {
		return fmt.Errorf("Redis address is required")
	}
	
	return nil
}

// Save saves configuration to file
func (h *HotReloader) Save() error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	data, err := json.MarshalIndent(h.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Ensure directory exists
	dir := filepath.Dir(h.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	if err := os.WriteFile(h.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// Update updates specific configuration values
func (h *HotReloader) Update(updateFunc func(*Config)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	updateFunc(h.config)
}

// fileWatcher watches a file for changes
type fileWatcher struct {
	path    string
	modTime time.Time
	ch      chan struct{}
	closed  bool
}

// newFileWatcher creates a new file watcher
func newFileWatcher(path string) (*fileWatcher, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	
	return &fileWatcher{
		path:    path,
		modTime: stat.ModTime(),
		ch:      make(chan struct{}, 1),
	}, nil
}

// Changed returns true if the file has been modified
func (w *fileWatcher) Changed() bool {
	stat, err := os.Stat(w.path)
	if err != nil {
		return false
	}
	
	if stat.ModTime().After(w.modTime) {
		w.modTime = stat.ModTime()
		return true
	}
	
	return false
}

// changed returns the change channel
func (w *fileWatcher) changed() <-chan struct{} {
	return w.ch
}

// Close closes the watcher
func (w *fileWatcher) Close() error {
	w.closed = true
	close(w.ch)
	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8888,
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  120,
			MaxConns:     10000,
		},
		Game: GameConfig{
			MaxPlayersPerRoom: 4,
			RoomTimeout:       300,
			MatchTimeout:      60,
			FPS:               60,
			FrameSkip:         3,
			BattleConfig: BattleConfig{
				MaxTowers:      20,
				MaxEnemies:     100,
				MaxProjectiles: 500,
				SpawnInterval:  2000,
				WaveInterval:   30000,
				GoldPerKill:    10,
				InitialGold:    100,
				BaseHP:         100,
			},
			DanmakuConfig: DanmakuConfig{
				MaxMessagesPerSec: 100,
				MaxLength:          100,
				Cooldown:           1000,
				RateLimit:          10,
			},
		},
		Database: DatabaseConfig{
			URI:        "mongodb://localhost:27017",
			Database:   "danmaku_game",
			PoolSize:   100,
			Timeout:    10,
			RetryCount: 3,
		},
		Redis: RedisConfig{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			PoolSize:     50,
			MinIdleConns: 10,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
		},
		Security: SecurityConfig{
			RateLimitReqPerSec: 100,
			EnableIP whitelist: false,
			IPWhitelist:        []string{},
			EnableCaptcha:      false,
			SessionTimeout:     3600,
			MaxLoginAttempts:   5,
		},
		Monitor: MonitorConfig{
			EnableMetrics:   true,
			EnableTracing:    false,
			EnableProfiling:  false,
			MetricsPort:     9090,
			TraceEndpoint:   "",
			ProfilePort:     6060,
			ReportInterval:  60,
		},
		GameServer: GameServerConfig{
			HeartbeatInterval: 5,
			MaxFrameDelay:     10,
			ReconnectTimeout:  30,
			SaveReplay:         true,
			ReplayPath:         "./replays",
			AIMode:             "normal",
		},
		Version: "1.0.0",
	}
}
