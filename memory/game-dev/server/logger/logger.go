package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Level 日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// Logger 日志记录器
type Logger struct {
	mu         sync.Mutex
	level      Level
	output     io.Writer
	prefix     string
	timeFormat string
	caller     bool
}

// Config 日志配置
type Config struct {
	Level      string // debug, info, warn, error, fatal
	Path       string // 日志文件路径
	MaxSize    int64  // 单文件最大 MB
	MaxBackups int    // 保留旧文件数量
	MaxAge     int    // 保留天数
	Compress   bool   // 压缩旧文件
	Caller     bool   // 是否显示调用者信息
}

// Option 日志选项
type Option func(*Logger)

// WithLevel 设置日志级别
func WithLevel(level Level) Option {
	return func(l *Logger) { l.level = level }
}

// WithOutput 设置输出
func WithOutput(w io.Writer) Option {
	return func(l *Logger) { l.output = w }
}

// WithPrefix 设置前缀
func WithPrefix(prefix string) Option {
	return func(l *Logger) { l.prefix = prefix }
}

// WithTimeFormat 设置时间格式
func WithTimeFormat(format string) Option {
	return func(l *Logger) { l.timeFormat = format }
}

// WithCaller 设置是否显示调用者
func WithCaller(caller bool) Option {
	return func(l *Logger) { l.caller = caller }
}

// New 创建日志记录器
func New(opts ...Option) *Logger {
	l := &Logger{
		level:      InfoLevel,
		output:     os.Stdout,
		timeFormat: "2006-01-02 15:04:05",
		caller:     true,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// NewFromConfig 从配置创建日志记录器
func NewFromConfig(cfg Config) (*Logger, error) {
	level := InfoLevel
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn", "warning":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	}

	var output io.Writer = os.Stdout
	if cfg.Path != "" {
		// 确保目录存在
		dir := filepath.Dir(cfg.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %w", err)
		}
		f, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("打开日志文件失败: %w", err)
		}
		output = f
	}

	return New(
		WithLevel(level),
		WithOutput(output),
		WithCaller(cfg.Caller),
	), nil
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// log 内部日志方法
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format(l.timeFormat)
	levelName := levelNames[level]

	// 构建日志行
	var line string
	if l.caller {
		_, file, lineNo, ok := runtime.Caller(2)
		if ok {
			caller := fmt.Sprintf("%s:%d", filepath.Base(file), lineNo)
			line = fmt.Sprintf("[%s] [%s] [%s] %s\n", timestamp, levelName, caller, msg)
		} else {
			line = fmt.Sprintf("[%s] [%s] %s\n", timestamp, levelName, msg)
		}
	} else {
		line = fmt.Sprintf("[%s] [%s] %s\n", timestamp, levelName, msg)
	}

	// 写入输出
	l.output.Write([]byte(line))

	// Fatal 级别退出
	if level == FatalLevel {
		os.Exit(1)
	}
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
}

// With 创建带前缀的子日志器
func (l *Logger) With(prefix string) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	return &Logger{
		level:      l.level,
		output:     l.output,
		prefix:     l.prefix + prefix,
		timeFormat: l.timeFormat,
		caller:     l.caller,
	}
}

// ===== 全局日志器 =====

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init 初始化全局日志器
func Init(opts ...Option) {
	once.Do(func() {
		defaultLogger = New(opts...)
	})
}

// InitFromConfig 从配置初始化
func InitFromConfig(cfg Config) error {
	var err error
	once.Do(func() {
		defaultLogger, err = NewFromConfig(cfg)
	})
	return err
}

// Default 获取默认日志器
func Default() *Logger {
	if defaultLogger == nil {
		defaultLogger = New()
	}
	return defaultLogger
}

// Debug 调试日志
func Debug(format string, args ...interface{}) {
	Default().Debug(format, args...)
}

// Info 信息日志
func Info(format string, args ...interface{}) {
	Default().Info(format, args...)
}

// Warn 警告日志
func Warn(format string, args ...interface{}) {
	Default().Warn(format, args...)
}

// Error 错误日志
func Error(format string, args ...interface{}) {
	Default().Error(format, args...)
}

// Fatal 致命错误日志
func Fatal(format string, args ...interface{}) {
	Default().Fatal(format, args...)
}

// ===== 游戏专用日志器 =====

// GameLogger 游戏日志工具
type GameLogger struct {
	*Logger
	module string
}

// NewGameLogger 创建游戏模块日志器
func NewGameLogger(module string) *GameLogger {
	return &GameLogger{
		Logger: Default().With("[" + module + "] "),
		module: module,
	}
}

// 预设模块日志器
var (
	RoomLogger    = NewGameLogger("ROOM")
	BattleLogger  = NewGameLogger("BATTLE")
	NetworkLogger = NewGameLogger("NET")
	DBLogger      = NewGameLogger("DB")
	AuthLogger    = NewGameLogger("AUTH")
	PayLogger     = NewGameLogger("PAY")
)

// ===== 结构化日志 =====

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// String 创建字符串字段
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int 创建整数字段
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 创建64位整数字段
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Bool 创建布尔字段
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Any 创建任意类型字段
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// WithFields 创建带字段的日志
func (l *Logger) WithFields(fields ...Field) string {
	var sb strings.Builder
	for i, f := range fields {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(f.Key)
		sb.WriteString("=")
		sb.WriteString(fmt.Sprintf("%v", f.Value))
	}
	return sb.String()
}
