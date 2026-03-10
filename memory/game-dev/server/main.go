package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"barrage-game/internal/config"
	"barrage-game/internal/handler"
	"barrage-game/internal/room"
	"barrage-game/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	// 初始化日志系统
	initLogger()

	// 加载配置
	cfg := config.Load()
	logger.Info("游戏服务器启动中...")

	// 创建房间管理器
	roomMgr := room.NewManager(30 * time.Second)
	logger.RoomLogger.Info("房间管理器初始化完成")

	// 注册路由
	registerRoutes(roomMgr, cfg)

	// 优雅关闭
	go gracefulShutdown(roomMgr)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Port)
	logger.Info("服务器监听: %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Fatal("服务器启动失败: %v", err)
	}
}

// initLogger 初始化日志系统
func initLogger() {
	// 使用日志轮转
	rotator, err := logger.NewRotatingLogger(logger.Config{
		Level:      "debug",
		Path:       "./logs/game.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
		Caller:     true,
	})
	if err != nil {
		// 回退到标准日志
		logger.Init(logger.WithLevel(logger.InfoLevel))
		logger.Warn("日志轮转初始化失败，使用标准输出: %v", err)
	} else {
		logger.Info("日志系统初始化完成")
		_ = rotator // 保持引用
	}
}

// registerRoutes 注册路由
func registerRoutes(roomMgr *room.Manager, cfg *config.Config) {
	// WebSocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		logger.NetworkLogger.Debug("WebSocket 连接: %s", r.RemoteAddr)
		handler.HandleWS(w, r, roomMgr)
	})

	// 房间 API
	http.HandleFunc("/api/room/create", handler.CreateRoom)
	http.HandleFunc("/api/room/join", handler.JoinRoom)
	http.HandleFunc("/api/room/list", handler.ListRooms)

	// 游戏 API
	http.HandleFunc("/api/game/start", handler.StartGame)
	http.HandleFunc("/api/game/action", handler.GameAction)
	http.HandleFunc("/api/game/end", handler.EndGame)

	// 支付 API
	http.HandleFunc("/api/payment/create", handler.CreatePayment)
	http.HandleFunc("/api/payment/notify", handler.PaymentNotify)

	// 排行榜
	http.HandleFunc("/api/leaderboard", handler.GetLeaderboard)

	// 健康检查
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}

// gracefulShutdown 优雅关闭
func gracefulShutdown(roomMgr *room.Manager) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	logger.Warn("收到退出信号: %v", sig)
	logger.Info("开始关闭服务器...")

	// 关闭所有房间
	roomMgr.Close()

	logger.Info("服务器已关闭")
	os.Exit(0)
}
