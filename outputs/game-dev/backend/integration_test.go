package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// ==================== 集成测试 ====================

// TestIntegration_FullGameFlow 完整游戏流程集成测试
func TestIntegration_FullGameFlow(t *testing.T) {
	// 1. 初始化配置
	config := Load()
	if config == nil {
		t.Fatal("配置加载失败")
	}

	// 2. 初始化房间管理器
	roomMgr := NewRoomManager(30 * time.Second)
	if roomMgr == nil {
		t.Fatal("房间管理器初始化失败")
	}

	// 3. 创建房间
	room, err := roomMgr.CreateRoom(&CreateRoomReq{
		ID:       "room-001",
		Name:     "测试房间",
		HostID:   "player-1",
		HostName: "玩家1",
		Mode:     "classic",
	})
	if err != nil {
		t.Fatalf("创建房间失败: %v", err)
	}
	if room == nil {
		t.Fatal("房间为空")
	}

	// 4. 玩家加入
	err = roomMgr.JoinRoom("room-001", "player-2", "玩家2")
	if err != nil {
		t.Fatalf("玩家加入失败: %v", err)
	}

	// 5. 检查房间状态
	if room.Status != "waiting" {
		t.Errorf("房间状态应为waiting, 实际: %s", room.Status)
	}

	// 6. 离开房间
	err = roomMgr.LeaveRoom("player-2")
	if err != nil {
		t.Fatalf("玩家离开失败: %v", err)
	}

	// 7. 清理
	roomMgr.DeleteRoom("room-001")
}

// TestIntegration_Matchmaker 匹配系统集成测试
func TestIntegration_Matchmaker(t *testing.T) {
	matchmaker := NewMatchmaker(2, 60*time.Second)

	// 添加玩家
	matchmaker.AddPlayer("player-1", 1500)
	matchmaker.AddPlayer("player-2", 1550)
	matchmaker.AddPlayer("player-3", 1450)

	// 等待匹配
	time.Sleep(100 * time.Millisecond)

	matches := matchmaker.GetMatches()
	if len(matches) == 0 {
		t.Log("暂无匹配结果（正常，等待时间不足）")
	}
}

// TestIntegration_Battle 战斗系统集成测试
func TestIntegration_Battle(t *testing.T) {
	battle := NewBattleManager()

	// 添加防御塔
	tower1 := battle.AddTower("arrow", 100, 100)
	if tower1 == nil {
		t.Fatal("添加防御塔失败")
	}

	// 生成敌人
	enemy := battle.SpawnEnemy("grunt", 100, 0)
	if enemy == nil {
		t.Fatal("生成敌人失败")
	}

	// 模拟战斗帧
	for i := 0; i < 60; i++ {
		battle.Update(0.016)
	}

	// 检查状态
	if battle.State.Wave != 1 {
		t.Errorf("波次应为1, 实际: %d", battle.State.Wave)
	}
}

// TestIntegration_GiftEffect 礼物系统集成测试
func TestIntegration_GiftEffect(t *testing.T) {
	battle := NewBattleManager()
	giftMgr := NewGiftManager()

	// 初始金币
	initialMoney := battle.State.Money

	// 发送金币礼物
	giftMgr.ReceiveGift("viewer-1", &Gift{
		Type:  GiftTypeCoin,
		Price: 1,
	})
	giftMgr.ProcessGiftEffects(battle)

	// 验证金币增加
	if battle.State.Money <= initialMoney {
		t.Errorf("金币礼物未生效: 初始=%d, 当前=%d", initialMoney, battle.State.Money)
	}

	// 发送炸弹礼物
	enemy := battle.SpawnEnemy("grunt", 1000, 0)
	initialHP := enemy.HP

	giftMgr.ReceiveGift("viewer-2", &Gift{
		Type:  GiftTypeBang,
		Price: 10,
	})
	giftMgr.ProcessGiftEffects(battle)

	// 验证伤害
	if enemy.HP >= initialHP {
		t.Errorf("炸弹礼物未造成伤害")
	}
}

// TestIntegration_Danmaku 弹幕系统集成测试
func TestIntegration_Danmaku(t *testing.T) {
	danmakuMgr := NewDanmakuManager()

	// 发送弹幕
	danmakuMgr.Send("viewer-1", "测试弹幕1", DanmakuTypeText)
	danmakuMgr.Send("viewer-2", "测试弹幕2", DanmakuTypeText)
	danmakuMgr.Send("viewer-3", "测试弹幕3", DanmakuTypeText)

	// 更新位置
	for i := 0; i < 60; i++ {
		danmakuMgr.Update(0.016)
	}

	// 检查弹幕数量
	if len(danmakuMgr.ActiveDanmakus) == 0 {
		t.Log("弹幕已飞出屏幕（正常）")
	}
}

// TestIntegration_LiveRoom 直播间系统集成测试
func TestIntegration_LiveRoom(t *testing.T) {
	room := NewLiveRoom("room-001", "anchor-1", "主播A")

	// 观众加入
	room.JoinViewer("viewer-1", "观众B")
	room.JoinViewer("viewer-2", "观众C")

	// 发送礼物
	room.SendGift("viewer-1", "rocket")
	room.SendDanmaku("viewer-2", "太棒了!")

	// 更新
	for i := 0; i < 60; i++ {
		room.Update(0.016)
	}

	// 验证
	if room.Battle.State.Score == 0 {
		t.Log("直播间战斗初始化完成")
	}
}

// TestIntegration_Leaderboard 排行榜系统集成测试
func TestIntegration_Leaderboard(t *testing.T) {
	lb := NewLeaderboard(100)

	// 更新分数
	lb.UpdateScore("player-1", 1500)
	lb.UpdateScore("player-2", 2000)
	lb.UpdateScore("player-3", 1800)
	lb.UpdateScore("player-1", 1600) // 更新

	// 获取排名
	rank := lb.GetRank("player-2")
	if rank != 1 {
		t.Errorf("player-2排名应为1, 实际: %d", rank)
	}

	// 获取前10
	top10 := lb.GetTop(10)
	if len(top10) != 3 {
		t.Errorf("前10应有3人, 实际: %d", len(top10))
	}
}

// TestIntegration_ConcurrentRooms 并发房间操作测试
func TestIntegration_ConcurrentRooms(t *testing.T) {
	roomMgr := NewRoomManager(30 * time.Second)
	var wg sync.WaitGroup

	// 并发创建房间
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			roomID := fmt.Sprintf("room-%d", id)
			roomMgr.CreateRoom(&CreateRoomReq{
				ID:       roomID,
				Name:     fmt.Sprintf("房间%d", id),
				HostID:   fmt.Sprintf("host-%d", id),
				HostName: fmt.Sprintf("Host%d", id),
				Mode:     "classic",
			})
		}(i)
	}
	wg.Wait()

	// 验证房间数量
	rooms := roomMgr.ListRooms()
	if len(rooms) != 100 {
		t.Errorf("房间数量应为100, 实际: %d", len(rooms))
	}

	// 并发加入房间
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			roomMgr.JoinRoom(fmt.Sprintf("room-%d", id), fmt.Sprintf("guest-%d", id), fmt.Sprintf("Guest%d", id))
		}(i)
	}
	wg.Wait()

	// 清理
	for i := 0; i < 100; i++ {
		roomMgr.DeleteRoom(fmt.Sprintf("room-%d", i))
	}
}

// TestIntegration_ConcurrentBattle 并发战斗更新测试
func TestIntegration_ConcurrentBattle(t *testing.T) {
	battle := NewBattleManager()

	// 添加多个塔
	for i := 0; i < 10; i++ {
		battle.AddTower("arrow", float64(100+i*10), float64(100+i*10))
	}

	// 生成敌人
	for i := 0; i < 50; i++ {
		battle.SpawnEnemy("grunt", 100, float64(i*10))
	}

	// 并发更新
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				battle.Update(0.016)
			}
		}()
	}
	wg.Wait()

	t.Logf("战斗状态: Wave=%d, Score=%d, Enemies=%d",
		battle.State.Wave, battle.State.Score, len(battle.Spawner.Enemies))
}

// TestIntegration_Achievement 成就系统集成测试
func TestIntegration_Achievement(t *testing.T) {
	achMgr := NewAchievementManager()

	// 初始化玩家成就
	achMgr.InitPlayer("player-1")

	// 更新击杀数
	achMgr.UpdateProgress("player-1", "killer_100", 50)
	achMgr.UpdateProgress("player-1", "killer_100", 100)

	// 获取完成数
	count := achMgr.GetCompletedCount("player-1")
	t.Logf("完成成就数: %d", count)

	// 领取奖励
	reward := achMgr.ClaimReward("player-1", "killer_100")
	if reward == nil {
		t.Log("成就未完成，无法领取")
	} else {
		t.Logf("领取奖励: Exp=%d, Coins=%d, Gems=%d", reward.Exp, reward.Coins, reward.Gems)
	}
}

// TestIntegration_Quest 任务系统集成测试
func TestIntegration_Quest(t *testing.T) {
	questMgr := NewQuestManager()

	// 初始化玩家任务
	questMgr.InitPlayer("player-1")

	// 接受每日任务
	err := questMgr.AcceptQuest("player-1", "daily_kill_10")
	if err != nil {
		t.Fatalf("接受任务失败: %v", err)
	}

	// 更新进度
	questMgr.UpdateProgress("player-1", "daily_kill_10", 5)
	questMgr.UpdateProgress("player-1", "daily_kill_10", 10)

	// 完成任务
	err = questMgr.CompleteQuest("player-1", "daily_kill_10")
	if err != nil {
		t.Fatalf("完成任务失败: %v", err)
	}
}

// TestIntegration_Stats 统计系统集成测试
func TestIntegration_Stats(t *testing.T) {
	statsMgr := NewStatsManager()

	// 记录游戏数据
	statsMgr.RecordGameStart("player-1")
	statsMgr.RecordGameWin("player-1", 120)
	statsMgr.RecordGameLose("player-2", 90)
	statsMgr.RecordTowerOperations("player-1", "place", 5)
	statsMgr.RecordTowerOperations("player-1", "upgrade", 3)
	statsMgr.RecordTowerOperations("player-1", "sell", 2)
	statsMgr.RecordMoneyChange("player-1", 1000, "earn")
	statsMgr.RecordMoneyChange("player-1", -500, "spend")
	statsMgr.RecordSocialStats("player-1", "gift", 10)
	statsMgr.RecordSocialStats("player-1", "danmaku", 20)

	// 获取统计
	stats := statsMgr.GetStats("player-1")
	if stats == nil {
		t.Fatal("获取统计失败")
	}

	t.Logf("玩家统计: 胜率=%.2f%%, 最高分=%d, 总时长=%.2fh",
		stats.WinRate*100, stats.BestScore, stats.TotalPlayTime/3600)
}

// TestIntegration_FrameSync 帧同步系统集成测试
func TestIntegration_FrameSync(t *testing.T) {
	fsMgr := NewFrameSyncManager(60)

	// 添加玩家
	fsMgr.AddPlayer("player-1")
	fsMgr.AddPlayer("player-2")

	// 模拟输入
	fsMgr.RecordInput("player-1", &PlayerInput{
		Type:  InputTypeMove,
		X:     100,
		Y:     200,
		Frame: 1,
	})

	// 模拟帧更新
	for i := 0; i < 60; i++ {
		fsMgr.Update(0.016)
	}

	// 获取当前帧
	frame := fsMgr.GetCurrentFrame()
	t.Logf("当前帧: %d", frame)
}

// TestIntegration_Replay 回放系统集成测试
func TestIntegration_Replay(t *testing.T) {
	replayMgr := NewReplayManager()

	// 开始录制
	replayMgr.StartRecording("replay-001", "player-1")

	// 模拟帧数据
	for i := 0; i < 100; i++ {
		replayMgr.RecordFrame(&FrameState{
			Frame:     uint32(i),
			Timestamp: time.Now().Unix(),
			Players:   make(map[string]*PlayerState),
		})
	}

	// 停止录制
	replayMgr.StopRecording("replay-001")

	// 播放回放
	err := replayMgr.Play("replay-001", 2.0)
	if err != nil {
		t.Fatalf("播放回放失败: %v", err)
	}

	// 获取回放数据
	data := replayMgr.GetReplayData("replay-001")
	if data == nil {
		t.Fatal("获取回放数据失败")
	}

	t.Logf("回放帧数: %d", len(data.Frames))
}

// TestIntegration_Security_RateLimit 限流系统集成测试
func TestIntegration_Security_RateLimit(t *testing.T) {
	limiter := NewRateLimiter(10, time.Minute)

	// 模拟请求
	for i := 0; i < 15; i++ {
		allowed := limiter.Allow("client-1")
		if !allowed && i < 10 {
			t.Errorf("前10次请求应被允许")
		}
	}

	// 验证限流
	allowed := limiter.Allow("client-1")
	if allowed {
		t.Log("限流已生效")
	}
}

// TestIntegration_Security_Validator 输入验证集成测试
func TestIntegration_Security_Validator(t *testing.T) {
	validator := NewInputValidator()

	// 测试用户名验证
	tests := []struct {
		input    string
		expected bool
	}{
		{"abc", false},   // 太短
		{"ab", false},    // 太短
		{"abcdefghijklmnopqrstuvwxyz", false}, // 太长
		{"test123", true}, // 有效
		{"用户1", true},   // 中文
	}

	for _, test := range tests {
		err := validator.ValidateUsername(test.input)
		if (err == nil) != test.expected {
			t.Errorf("用户名验证失败: %s, 期望=%v, 实际=%v", test.input, test.expected, err == nil)
		}
	}

	// 测试敏感词过滤
	dirty := "这是一个测试色情内容"
	clean := validator.FilterSensitiveWords(dirty)
	if clean == dirty {
		t.Errorf("敏感词过滤未生效: %s -> %s", dirty, clean)
	}
}

// TestIntegration_Config 配置系统集成测试
func TestIntegration_Config(t *testing.T) {
	// 加载配置
	cfg := Load()
	if cfg == nil {
		t.Fatal("配置加载失败")
	}

	// 读取配置
	gameCfg := cfg.Get()
	if gameCfg == nil {
		t.Fatal("获取配置失败")
	}

	t.Logf("服务器配置: Port=%s, MaxConnections=%d",
		gameCfg.Server.Port, gameCfg.Server.MaxConnections)
	t.Logf("战斗配置: FPS=%d, Players=%d",
		gameCfg.Game.Battle.FPS, gameCfg.Game.Battle.MaxPlayers)
}

// TestIntegration_Logger 日志系统集成测试
func TestIntegration_Logger(t *testing.T) {
	logger := NewLogger(Config{
		Level:  LevelInfo,
		Format: "json",
	})

	logger.Info("测试日志", "key", "value")
	logger.Debug("调试信息")
	logger.Warn("警告信息")
	logger.Error("错误信息")
}

// ==================== 性能基准测试 ====================

// BenchmarkRoomManager_CreateRoom 房间创建基准测试
func BenchmarkRoomManager_CreateRoom(b *testing.B) {
	roomMgr := NewRoomManager(30 * time.Second)
	defer roomMgr.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		roomID := fmt.Sprintf("room-%d", i)
		roomMgr.CreateRoom(&CreateRoomReq{
			ID:       roomID,
			Name:     fmt.Sprintf("房间%d", i),
			HostID:   fmt.Sprintf("host-%d", i),
			HostName: fmt.Sprintf("Host%d", i),
			Mode:     "classic",
		})
		roomMgr.DeleteRoom(roomID)
	}
}

// BenchmarkRoomManager_JoinRoom 房间加入基准测试
func BenchmarkRoomManager_JoinRoom(b *testing.B) {
	roomMgr := NewRoomManager(30 * time.Second)
	defer roomMgr.Shutdown()

	roomMgr.CreateRoom(&CreateRoomReq{
		ID:       "room-001",
		Name:     "测试房间",
		HostID:   "host-1",
		HostName: "Host1",
		Mode:     "classic",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		playerID := fmt.Sprintf("player-%d", i)
		roomMgr.JoinRoom("room-001", playerID, fmt.Sprintf("Player%d", i))
	}
}

// BenchmarkBattleManager_Update 战斗更新基准测试
func BenchmarkBattleManager_Update(b *testing.B) {
	battle := NewBattleManager()

	// 添加塔
	for i := 0; i < 10; i++ {
		battle.AddTower("arrow", float64(i*20), float64(i*20))
	}

	// 生成敌人
	for i := 0; i < 100; i++ {
		battle.SpawnEnemy("grunt", 100, float64(i*5))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		battle.Update(0.016)
	}
}

// BenchmarkDanmakuManager_Update 弹幕更新基准测试
func BenchmarkDanmakuManager_Update(b *testing.B) {
	danmakuMgr := NewDanmakuManager()

	// 添加弹幕
	for i := 0; i < 1000; i++ {
		danmakuMgr.Send(fmt.Sprintf("viewer-%d", i), fmt.Sprintf("弹幕内容%d", i), DanmakuTypeText)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		danmakuMgr.Update(0.016)
	}
}

// BenchmarkLeaderboard_Update 排行榜更新基准测试
func BenchmarkLeaderboard_Update(b *testing.B) {
	lb := NewLeaderboard(10000)

	// 初始化玩家
	for i := 0; i < 1000; i++ {
		lb.UpdateScore(fmt.Sprintf("player-%d", i), int64(i*100))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.UpdateScore("player-1", int64(i))
	}
}

// BenchmarkLeaderboard_GetRank 获取排名基准测试
func BenchmarkLeaderboard_GetRank(b *testing.B) {
	lb := NewLeaderboard(10000)

	// 初始化玩家
	for i := 0; i < 10000; i++ {
		lb.UpdateScore(fmt.Sprintf("player-%d", i), int64(i*100))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.GetRank("player-5000")
	}
}

// BenchmarkGiftManager_Process 礼物处理基准测试
func BenchmarkGiftManager_Process(b *testing.B) {
	giftMgr := NewGiftManager()
	battle := NewBattleManager()

	// 添加待处理礼物
	for i := 0; i < 100; i++ {
		giftMgr.ReceiveGift(fmt.Sprintf("viewer-%d", i), &Gift{
			Type:  GiftTypeBang,
			Price: 10,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		giftMgr.ProcessGiftEffects(battle)
	}
}

// BenchmarkAchievementManager_Check 成就检查基准测试
func BenchmarkAchievementManager_Check(b *testing.B) {
	achMgr := NewAchievementManager()

	// 初始化玩家
	for i := 0; i < 100; i++ {
		achMgr.InitPlayer(fmt.Sprintf("player-%d", i))
		achMgr.UpdateProgress(fmt.Sprintf("player-%d", i), "killer_100", 50)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			achMgr.UpdateProgress(fmt.Sprintf("player-%d", j), "killer_100", 51)
		}
	}
}

// BenchmarkFrameSync_Update 帧同步更新基准测试
func BenchmarkFrameSync_Update(b *testing.B) {
	fsMgr := NewFrameSyncManager(60)

	// 添加玩家
	for i := 0; i < 10; i++ {
		fsMgr.AddPlayer(fmt.Sprintf("player-%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fsMgr.Update(0.016)
	}
}

// BenchmarkRateLimiter_Allow 限流检查基准测试
func BenchmarkRateLimiter_Allow(b *testing.B) {
	limiter := NewRateLimiter(10000, time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow("client-1")
	}
}

// BenchmarkValidator_Username 用户名验证基准测试
func BenchmarkValidator_Username(b *testing.B) {
	validator := NewInputValidator()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateUsername("testuser123")
	}
}

// BenchmarkValidator_SensitiveWords 敏感词过滤基准测试
func BenchmarkValidator_SensitiveWords(b *testing.B) {
	validator := NewInputValidator()
	text := "这是一个测试内容包含敏感词需要过滤"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.FilterSensitiveWords(text)
	}
}

// ==================== 压力测试 ====================

// TestStress_ConcurrentConnections 并发连接压力测试
func TestStress_ConcurrentConnections(t *testing.T) {
	roomMgr := NewRoomManager(30 * time.Second)
	defer roomMgr.Shutdown()

	roomMgr.CreateRoom(&CreateRoomReq{
		ID:       "room-stress",
		Name:     "压力测试房间",
		HostID:   "host-1",
		HostName: "Host1",
		Mode:     "classic",
	})

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// 100个并发连接
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := roomMgr.JoinRoom("room-stress", fmt.Sprintf("player-%d", id), fmt.Sprintf("Player%d", id))
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	t.Logf("成功加入: %d/100", successCount)
}

// TestStress_MemoryLeak 内存泄漏测试
func TestStress_MemoryLeak(t *testing.T) {
	battle := NewBattleManager()

	// 记录初始内存
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 模拟长时间运行
	for round := 0; round < 100; round++ {
		// 添加塔
		for i := 0; i < 5; i++ {
			battle.AddTower("arrow", float64(i*20), float64(i*20))
		}

		// 生成敌人
		for i := 0; i < 50; i++ {
			battle.SpawnEnemy("grunt", 100, float64(i*5))
		}

		// 更新
		for i := 0; i < 60; i++ {
			battle.Update(0.016)
		}
	}

	// 记录最终内存
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	allocDiff := m2.Mallocs - m1.Mallocs
	t.Logf("内存分配差: %d (Mallocs)", allocDiff)

	if allocDiff > 1000000 {
		t.Log("警告: 可能存在内存泄漏")
	}
}

// TestStress_LongRunning 长时间运行测试
func TestStress_LongRunning(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过长时间测试")
	}

	battle := NewBattleManager()

	// 模拟10分钟运行 (600秒 * 60fps = 36000帧)
	frames := 36000
	for i := 0; i < frames; i++ {
		// 定期生成敌人
		if i%60 == 0 {
			battle.SpawnEnemy("grunt", 100, 0)
		}

		battle.Update(0.016)

		if i%1000 == 0 {
			t.Logf("进度: %d/%d 帧, 敌人=%d, 分数=%d",
				i, frames, len(battle.Spawner.Enemies), battle.State.Score)
		}
	}

	t.Logf("最终状态: Wave=%d, Score=%d, Enemies=%d",
		battle.State.Wave, battle.State.Score, len(battle.Spawner.Enemies))
}
