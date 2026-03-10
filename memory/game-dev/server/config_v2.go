package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// GameConfig 游戏配置主结构
type GameConfig struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Log      LogConfig      `json:"log"`
	Game     GameSettings   `json:"game"`
	Security SecurityConfig `json:"security"`
	Monitor  MonitorConfig  `json:"monitor"`

	mu          sync.RWMutex
	watcher     *fileWatcher
	lastLoadAt  time.Time
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        `json:"port"`
	Host         string        `json:"host"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	IdleTimeout  time.Duration `json:"idleTimeout"`
	MaxConns     int           `json:"maxConns"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	MongoURI        string        `json:"mongoUri"`
	MongoDB         string        `json:"mongoDB"`
	MongoPoolSize   int           `json:"mongoPoolSize"`
	MongoTimeout    time.Duration `json:"mongoTimeout"`
	SQLitePath      string        `json:"sqlitePath"`
	UseSQLite       bool          `json:"useSqlite"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string        `json:"addr"`
	Password     string        `json:"password"`
	DB           int           `json:"db"`
	PoolSize     int           `json:"poolSize"`
	DialTimeout  time.Duration `json:"dialTimeout"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string        `json:"level"`
	Path       string        `json:"path"`
	MaxSize    int64         `json:"maxSize"`
	MaxBackups int           `json:"maxBackups"`
	MaxAge     int           `json:"maxAge"`
	Compress   bool          `json:"compress"`
	Console    bool          `json:"console"`
}

// GameSettings 游戏设置
type GameSettings struct {
	Battle     BattleConfig     `json:"battle"`
	Room       RoomConfig       `json:"room"`
	AI         AIConfig         `json:"ai"`
	Economy    EconomyConfig    `json:"economy"`
	Tower      TowerConfigs     `json:"tower"`
	Enemy      EnemyConfigs     `json:"enemy"`
	Wave       WaveConfigs     `json:"wave"`
	Gift       GiftConfigs      `json:"gift"`
	Danmaku    DanmakuConfig    `json:"danmaku"`
}

// BattleConfig 战斗配置
type BattleConfig struct {
	FPS            int     `json:"fps"`
	TickRate       int     `json:"tickRate"`
	MaxPlayers     int     `json:"maxPlayers"`
	PreloadTime    float64 `json:"preloadTime"`
	AutoStartDelay int     `json:"autoStartDelay"`
	WinScore       int     `json:"winScore"`
	LoseScore      int     `json:"loseScore"`
}

// RoomConfig 房间配置
type RoomConfig struct {
	MaxRooms           int           `json:"maxRooms"`
	MaxPlayersPerRoom  int           `json:"maxPlayersPerRoom"`
	RoomTimeout        time.Duration `json:"roomTimeout"`
	HeartbeatInterval  time.Duration `json:"heartbeatInterval"`
	ReconnectTimeout   time.Duration `json:"reconnectTimeout"`
	AutoMatchTimeout   time.Duration `json:"autoMatchTimeout"`
}

// AIConfig AI配置
type AIConfig struct {
	Enable           bool    `json:"enable"`
	Difficulty       string  `json:"difficulty"`
	ThinkInterval    int     `json:"thinkInterval"`
	MaxThinkTime     int     `json:"maxThinkTime"`
	AggressionFactor float64 `json:"aggressionFactor"`
	DefendFactor     float64 `json:"defendFactor"`
}

// EconomyConfig 经济系统配置
type EconomyConfig struct {
	InitialCoins      int `json:"initialCoins"`
	InitialDiamonds   int `json:"initialDiamonds"`
	KillReward        int `json:"killReward"`
	WaveClearReward   int `json:"waveClearReward"`
	WinReward         int `json:"winReward"`
	LoseReward        int `json:"loseReward"`
	DailyBonus        int `json:"dailyBonus"`
	LoginBonus        int `json:"loginBonus"`
}

// TowerConfigs 防御塔配置组
type TowerConfigs struct {
	Arrow  TowerConfig `json:"arrow"`
	Cannon TowerConfig `json:"cannon"`
	Ice    TowerConfig `json:"ice"`
	Lightning TowerConfig `json:"lightning"`
	Heal   TowerConfig `json:"heal"`
}

// TowerConfig 单个防御塔配置
type TowerConfig struct {
	Name        string  `json:"name"`
	BaseDamage  float64 `json:"baseDamage"`
	BaseRange   float64 `json:"baseRange"`
	BaseAttackSpeed float64 `json:"baseAttackSpeed"`
	BaseCost    int     `json:"baseCost"`
	LevelFactor float64 `json:"levelFactor"`
	Special     string  `json:"special"`
}

// EnemyConfigs 敌人配置组
type EnemyConfigs struct {
	Grunt  EnemyConfig `json:"grunt"`
	Ranger EnemyConfig `json:"ranger"`
	Tank   EnemyConfig `json:"tank"`
	Boss   EnemyConfig `json:"boss"`
}

// EnemyConfig 单个敌人配置
type EnemyConfig struct {
	Name         string  `json:"name"`
	BaseHP       float64 `json:"baseHP"`
	BaseArmor    float64 `json:"baseArmor"`
	BaseSpeed    float64 `json:"baseSpeed"`
	Bounty       int     `json:"bounty"`
	SpawnCost    int     `json:"spawnCost"`
}

// WaveConfigs 波次配置组
type WaveConfigs struct {
	Wave1 WaveConfig `json:"wave1"`
	Wave2 WaveConfig `json:"wave2"`
	Wave3 WaveConfig `json:"wave3"`
	Wave4 WaveConfig `json:"wave4"`
}

// WaveConfig 波次配置
type WaveConfig struct {
	Duration    int       `json:"duration"`
	Enemies     []EnemySpawn `json:"enemies"`
	RewardCoins int       `json:"rewardCoins"`
}

// EnemySpawn 敌人生成
type EnemySpawn struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
	Delay int    `json:"delay"`
}

// GiftConfigs 礼物配置组
type GiftConfigs struct {
	Coin    GiftConfig `json:"coin"`
	Star    GiftConfig `json:"star"`
	Rocket  GiftConfig `json:"rocket"`
	Car     GiftConfig `json:"car"`
	Plane   GiftConfig `json:"plane"`
	Bang    GiftConfig `json:"bang"`
}

// GiftConfig 礼物配置
type GiftConfig struct {
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Coins     int    `json:"coins"`
	Effect    string `json:"effect"`
	Duration  int    `json:"duration"`
}

// DanmakuConfig 弹幕配置
type DanmakuConfig struct {
	MaxCount      int    `json:"maxCount"`
	MaxLength     int    `json:"maxLength"`
	Speed         int    `json:"speed"`
	DefaultColor  string `json:"defaultColor"`
	DefaultSize   int    `json:"defaultSize"`
	FilterEnabled bool   `json:"filterEnabled"`
	FilterWords   []string `json:"filterWords"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	RateLimit          RateLimitConfig  `json:"rateLimit"`
	IPBlackList       []string         `json:"ipBlackList"`
	EnableWAF          bool             `json:"enableWAF"`
	MaxRequestSize    int64            `json:"maxRequestSize"`
	AllowedOrigins    []string         `json:"allowedOrigins"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enable           bool          `json:"enable"`
	RequestsPerSec   int           `json:"requestsPerSec"`
	BurstSize        int           `json:"burstSize"`
	BlockDuration    time.Duration `json:"blockDuration"`
	WSMsgPerSec      int           `json:"wsMsgPerSec"`
	WSBurstSize      int           `json:"wsBurstSize"`
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	EnableMetrics    bool          `json:"enableMetrics"`
	EnableTracing    bool          `json:"enableTracing"`
	MetricsPort      int           `json:"metricsPort"`
	StatsdAddr       string        `json:"statsdAddr"`
	SamplingRate     float64       `json:"samplingRate"`
}

// 全局配置实例
var (
	defaultConfig *GameConfig
	configOnce    sync.Once
)

// Load 加载默认配置
func Load() *GameConfig {
	configOnce.Do(func() {
		defaultConfig = loadConfig("config.json")
	})
	return defaultConfig
}

// LoadFromFile 从文件加载配置
func LoadFromFile(path string) (*GameConfig, error) {
	return loadConfig(path), nil
}

// loadConfig 内部加载配置
func loadConfig(path string) *GameConfig {
	cfg := getDefaultConfig()

	// 尝试从文件加载
	if path != "" {
		if data, err := os.ReadFile(path); err == nil {
			if err := json.Unmarshal(data, cfg); err == nil {
				cfg.lastLoadAt = time.Now()
				return cfg
			}
		}
	}

	// 从环境变量覆盖
	loadFromEnv(cfg)

	cfg.lastLoadAt = time.Now()
	return cfg
}

// getDefaultConfig 返回默认配置
func getDefaultConfig() *GameConfig {
	return &GameConfig{
		Server: ServerConfig{
			Port:         "8080",
			Host:         "0.0.0.0",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
			MaxConns:     10000,
		},
		Database: DatabaseConfig{
			MongoURI:      "mongodb://localhost:27017",
			MongoDB:       "danmaku_game",
			MongoPoolSize: 100,
			UseSQLite:     false,
		},
		Redis: RedisConfig{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			PoolSize:     100,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		Log: LogConfig{
			Level:      "info",
			Path:       "./logs",
			MaxSize:    100,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   true,
			Console:    true,
		},
		Game: GameSettings{
			Battle: BattleConfig{
				FPS:            60,
				TickRate:       30,
				MaxPlayers:     4,
				PreloadTime:    5.0,
				AutoStartDelay: 10,
				WinScore:       100,
				LoseScore:      10,
			},
			Room: RoomConfig{
				MaxRooms:          1000,
				MaxPlayersPerRoom: 4,
				RoomTimeout:       30 * time.Minute,
				HeartbeatInterval: 10 * time.Second,
				ReconnectTimeout:  30 * time.Second,
				AutoMatchTimeout:  60 * time.Second,
			},
			AI: AIConfig{
				Enable:           true,
				Difficulty:       "normal",
				ThinkInterval:    100,
				MaxThinkTime:     50,
				AggressionFactor: 0.5,
				DefendFactor:     0.5,
			},
			Economy: EconomyConfig{
				InitialCoins:    1000,
				InitialDiamonds: 100,
				KillReward:      10,
				WaveClearReward: 50,
				WinReward:       100,
				LoseReward:      20,
				DailyBonus:      200,
				LoginBonus:      100,
			},
			Tower: TowerConfigs{
				Arrow: TowerConfig{
					Name: "箭塔", BaseDamage: 10, BaseRange: 150,
					BaseAttackSpeed: 1.0, BaseCost: 50,
					LevelFactor: 1.2, Special: "none",
				},
				Cannon: TowerConfig{
					Name: "炮塔", BaseDamage: 30, BaseRange: 120,
					BaseAttackSpeed: 0.5, BaseCost: 100,
					LevelFactor: 1.3, Special: "splash",
				},
				Ice: TowerConfig{
					Name: "冰塔", BaseDamage: 5, BaseRange: 130,
					BaseAttackSpeed: 1.5, BaseCost: 80,
					LevelFactor: 1.15, Special: "slow",
				},
				Lightning: TowerConfig{
					Name: "雷塔", BaseDamage: 20, BaseRange: 160,
					BaseAttackSpeed: 0.8, BaseCost: 120,
					LevelFactor: 1.25, Special: "chain",
				},
				Heal: TowerConfig{
					Name: "治疗塔", BaseDamage: -15, BaseRange: 100,
					BaseAttackSpeed: 1.0, BaseCost: 90,
					LevelFactor: 1.2, Special: "heal",
				},
			},
			Enemy: EnemyConfigs{
				Grunt: EnemyConfig{Name: "步兵", BaseHP: 100, BaseArmor: 5, BaseSpeed: 50, Bounty: 10, SpawnCost: 5},
				Ranger: EnemyConfig{Name: "弓手", BaseHP: 80, BaseArmor: 2, BaseSpeed: 60, Bounty: 15, SpawnCost: 8},
				Tank: EnemyConfig{Name: "坦克", BaseHP: 300, BaseArmor: 20, BaseSpeed: 30, Bounty: 30, SpawnCost: 15},
				Boss: EnemyConfig{Name: "Boss", BaseHP: 1000, BaseArmor: 50, BaseSpeed: 20, Bounty: 100, SpawnCost: 50},
			},
			Wave: WaveConfigs{
				Wave1: WaveConfig{
					Duration: 30, RewardCoins: 100,
					Enemies: []EnemySpawn{
						{Type: "grunt", Count: 5, Delay: 3},
					},
				},
				Wave2: WaveConfig{
					Duration: 45, RewardCoins: 150,
					Enemies: []EnemySpawn{
						{Type: "grunt", Count: 8, Delay: 2},
						{Type: "ranger", Count: 3, Delay: 5},
					},
				},
				Wave3: WaveConfig{
					Duration: 60, RewardCoins: 200,
					Enemies: []EnemySpawn{
						{Type: "grunt", Count: 10, Delay: 2},
						{Type: "ranger", Count: 5, Delay: 4},
						{Type: "tank", Count: 2, Delay: 8},
					},
				},
				Wave4: WaveConfig{
					Duration: 90, RewardCoins: 300,
					Enemies: []EnemySpawn{
						{Type: "grunt", Count: 15, Delay: 2},
						{Type: "ranger", Count: 8, Delay: 3},
						{Type: "tank", Count: 5, Delay: 5},
						{Type: "boss", Count: 1, Delay: 15},
					},
				},
			},
			Gift: GiftConfigs{
				Coin:   GiftConfig{Name: "金币", Price: 1, Coins: 10, Effect: "coins", Duration: 0},
				Star:   GiftConfig{Name: "星星", Price: 9, Coins: 100, Effect: "coins", Duration: 0},
				Rocket: GiftConfig{Name: "火箭", Price: 99, Coins: 1000, Effect: "damage", Duration: 60},
				Car:    GiftConfig{Name: "跑车", Price: 299, Coins: 3000, Effect: "damage", Duration: 120},
				Plane:  GiftConfig{Name: "飞机", Price: 999, Coins: 10000, Effect: "nuke", Duration: 30},
				Bang:  GiftConfig{Name: "轰炸", Price: 1999, Coins: 20000, Effect: "kill_all", Duration: 0},
			},
			Danmaku: DanmakuConfig{
				MaxCount:     50,
				MaxLength:    100,
				Speed:        200,
				DefaultColor: "#FFFFFF",
				DefaultSize:  16,
				FilterEnabled: true,
			},
		},
		Security: SecurityConfig{
			RateLimit: RateLimitConfig{
				Enable:         true,
				RequestsPerSec: 100,
				BurstSize:      200,
				BlockDuration:  5 * time.Minute,
				WSMsgPerSec:    50,
				WSBurstSize:   100,
			},
			EnableWAF:       true,
			MaxRequestSize:  1024 * 1024,
			AllowedOrigins:  []string{"*"},
		},
		Monitor: MonitorConfig{
			EnableMetrics: true,
			EnableTracing: false,
			MetricsPort:   9090,
			StatsdAddr:    "localhost:8125",
			SamplingRate:  0.1,
		},
	}
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(cfg *GameConfig) {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("SERVER_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("MONGO_URI"); v != "" {
		cfg.Database.MongoURI = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
}

// Get 获取当前配置 (线程安全)
func (c *GameConfig) Get() *GameConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c
}

// Reload 重新加载配置
func (c *GameConfig) Reload(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	newCfg := loadConfig(path)
	*c = *newCfg
	return nil
}

// Watch 监视配置文件变化自动重载
func (c *GameConfig) Watch(path string, interval time.Duration) error {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	watcher, err := newFileWatcher(dir, filename, interval)
	if err != nil {
		return err
	}

	c.watcher = watcher
	go func() {
		for range watcher.ch {
			if err := c.Reload(path); err != nil {
				fmt.Printf("Config reload failed: %v\n", err)
			}
		}
	}()

	return nil
}

// StopWatch 停止监视
func (c *GameConfig) StopWatch() {
	if c.watcher != nil {
		c.watcher.close()
	}
}

// fileWatcher 文件监视器
type fileWatcher struct {
	dir      string
	file     string
	interval time.Duration
	ch       chan struct{}
	closed   bool
}

func newFileWatcher(dir, file string, interval time.Duration) (*fileWatcher, error) {
	w := &fileWatcher{
		dir:      dir,
		file:     file,
		interval: interval,
		ch:       make(chan struct{}, 1),
	}

	if err := w.init(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *fileWatcher) init() error {
	// 确保目录存在
	if err := os.MkdirAll(w.dir, 0755); err != nil {
		return err
	}

	// 启动监视循环
	go w.watchLoop()
	return nil
}

func (w *fileWatcher) watchLoop() {
	var lastMod time.Time

	for !w.closed {
		path := filepath.Join(w.dir, w.file)
		info, err := os.Stat(path)
		if err == nil {
			if lastMod.IsZero() {
				lastMod = info.ModTime()
			} else if info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				select {
				case w.ch <- struct{}{}:
				default:
				}
			}
		}

		time.Sleep(w.interval)
	}
}

func (w *fileWatcher) close() {
	w.closed = true
	close(w.ch)
}
