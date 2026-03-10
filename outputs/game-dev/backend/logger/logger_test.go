package logger

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := New(
		WithLevel(DebugLevel),
		WithCaller(true),
	)

	logger.Debug("调试信息: %s", "test")
	logger.Info("信息: %d", 123)
	logger.Warn("警告: %v", []int{1, 2, 3})
	logger.Error("错误: %s", "something wrong")

	// 测试带前缀
	sub := logger.With("[SUB]")
	sub.Info("子日志")
}

func TestWithFields(t *testing.T) {
	logger := New(WithLevel(InfoLevel))

	fields := []Field{
		String("player_id", "player_001"),
		Int("room_id", 123),
		Int64("damage", 1500),
		Bool("critical", true),
	}

	logger.Info(logger.WithFields(fields...))
}

func TestRotator(t *testing.T) {
	rotator, err := NewRotator("./test.log", 1024*1024, 3, 7, false)
	if err != nil {
		t.Fatal(err)
	}
	defer rotator.Close()

	// 写入测试数据
	for i := 0; i < 100; i++ {
		rotator.Write([]byte("Test log entry\n"))
	}
}

func BenchmarkLogger(b *testing.B) {
	logger := New(WithLevel(InfoLevel))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark test %d", i)
	}
}

func BenchmarkLoggerWithLock(b *testing.B) {
	logger := New(WithLevel(InfoLevel))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Info("Parallel test %d", i)
			i++
		}
	})
}

// 示例代码（不会被执行）
func Example() {
	// 基础用法
	logger := New(
		WithLevel(InfoLevel),
		WithCaller(true),
	)
	logger.Info("游戏开始")

	// 使用默认日志器
	Init(WithLevel(DebugLevel))
	Info("服务器启动")
	Debug("调试信息: %s", "value")

	// 使用模块日志器
	RoomLogger.Info("玩家加入房间")
	BattleLogger.Info("战斗开始: 波次=%d", 1)
	NetworkLogger.Error("连接失败: %s", "timeout")

	// 使用字段
	fields := []Field{
		String("player_id", "p001"),
		Int("score", 1000),
		Int64("gold", 5000),
	}
	Default().Info(Default().WithFields(fields...))

	// 使用轮转日志
	rotator, err := NewRotatingLogger(Config{
		Level:      "info",
		Path:       "./logs/game.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	})
	if err != nil {
		panic(err)
	}
	defer rotator.Close()

	rotator.Info("服务器启动")
	rotator.Close()

	// 从配置初始化
	InitFromConfig(DefaultRotator)
	Info("日志系统初始化完成")
	time.Sleep(time.Millisecond)
}
