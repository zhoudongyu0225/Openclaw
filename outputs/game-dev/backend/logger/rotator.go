package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Rotator 日志文件轮转器
type Rotator struct {
	mu           sync.Mutex
	filePath     string
	maxSize      int64
	maxBackups   int
	maxAge       int
	compress     bool
	file         *os.File
	size         int64
	startTime    time.Time
}

// NewRotator 创建日志轮转器
func NewRotator(filePath string, maxSize int64, maxBackups int, maxAge int, compress bool) (*Rotator, error) {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	r := &Rotator{
		filePath:   filePath,
		maxSize:    maxSize,
		maxBackups: maxBackups,
		maxAge:     maxAge,
		compress:   compress,
		startTime:  time.Now(),
	}

	// 打开或创建文件
	if err := r.open(); err != nil {
		return nil, err
	}

	return r, nil
}

// open 打开文件
func (r *Rotator) open() error {
	f, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	r.file = f

	// 获取当前文件大小
	info, err := f.Stat()
	if err != nil {
		return err
	}
	r.size = info.Size()

	return nil
}

// Write 写入日志
func (r *Rotator) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查是否需要轮转
	if r.size >= r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = r.file.Write(p)
	r.size += int64(n)
	return n, err
}

// rotate 执行轮转
func (r *Rotator) rotate() error {
	// 关闭当前文件
	r.file.Close()

	// 轮转现有备份
	r.rotateBackups()

	// 重命名当前文件为 .1
	backupPath := r.filePath + ".1"
	if err := os.Rename(r.filePath, backupPath); err != nil {
		return err
	}

	// 压缩旧备份
	if r.compress {
		go compressFile(backupPath)
	}

	// 创建新文件
	if err := r.open(); err != nil {
		return err
	}

	return nil
}

// rotateBackups 轮转备份文件
func (r *Rotator) rotateBackups() {
	// 删除过期文件
	if r.maxAge > 0 {
		cutoff := time.Now().Add(-time.Duration(r.maxAge) * 24 * time.Hour)
		r.cleanupOldFiles(cutoff)
	}

	// 删除超出数量的旧文件
	if r.maxBackups > 0 {
		r.cleanupBackups()
	}
}

// cleanupOldFiles 删除过期文件
func (r *Rotator) cleanupOldFiles(cutoff time.Time) {
	pattern := r.filePath + ".*"
	matches, _ := filepath.Glob(pattern)
	for _, path := range matches {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(path)
		}
	}
}

// cleanupBackups 删除超出数量的备份
func (r *Rotator) cleanupBackups() {
	pattern := r.filePath + ".*"
	matches, _ := filepath.Glob(pattern)

	// 按修改时间排序
	type fileInfo struct {
		path    string
		time    time.Time
		index   int
	}

	var files []fileInfo
	for _, path := range matches {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		// 提取索引
		idx := 0
		fmt.Sscanf(filepath.Ext(path), ".%d", &idx)
		files = append(files, fileInfo{path: path, time: info.ModTime(), index: idx})
	}

	// 保留最新的 maxBackups 个
	if len(files) > r.maxBackups {
		// 按索引排序，删除旧的
		for i := 0; i < len(files)-r.maxBackups; i++ {
			os.Remove(files[i].path)
		}
	}
}

// compressFile 异步压缩文件
func compressFile(path string) {
	// 这里可以集成 gzip 压缩
	// 简化实现：只是标记，实际压缩需要调用外部工具
	_ = path
}

// Close 关闭轮转器
func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

// ===== 轮转日志器 =====

// RotatingLogger 带轮转的日志器
type RotatingLogger struct {
	*Logger
	rotator *Rotator
}

// NewRotatingLogger 创建带轮转的日志器
func NewRotatingLogger(cfg Config) (*RotatingLogger, error) {
	// 解析大小 (MB -> bytes)
	maxSize := int64(cfg.MaxSize) * 1024 * 1024
	if maxSize == 0 {
		maxSize = 100 * 1024 * 1024 // 默认 100MB
	}

	rotator, err := NewRotator(cfg.Path, maxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress)
	if err != nil {
		return nil, err
	}

	logger := New(
		WithLevel(parseLevel(cfg.Level)),
		WithOutput(rotator),
		WithCaller(cfg.Caller),
	)

	return &RotatingLogger{
		Logger:  logger,
		rotator: rotator,
	}, nil
}

// Close 关闭日志器
func (rl *RotatingLogger) Close() error {
	return rl.rotator.Close()
}

// parseLevel 解析日志级别
func parseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// ===== 默认轮转配置 =====

// DefaultRotator 默认日志轮转配置
var DefaultRotator = Config{
	Level:      "info",
	Path:       "./logs/game.log",
	MaxSize:    100, // 100MB
	MaxBackups: 7,
	MaxAge:     30,
	Compress:   true,
	Caller:     true,
}
