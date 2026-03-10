package main

import (
	"fmt"
	"testing"
	"time"
)

// ============================================
// 成就系统测试
// ============================================

func TestAchievementManager_Register(t *testing.T) {
	am := NewAchievementManager()
	
	// 测试注册成就
	if len(am.Achievements) == 0 {
		t.Error("成就管理器应该已注册成就")
	}
	
	// 测试获取成就
	achievement := am.GetAchievement("killer_100")
	if achievement == nil {
		t.Error("应该能获取已注册的成就")
	}
	
	if achievement.Name != "初出茅庐" {
		t.Errorf("期望成就名称为 '初出茅庐', 实际为 '%s'", achievement.Name)
	}
}

func TestAchievementManager_InitPlayer(t *testing.T) {
	am := NewAchievementManager()
	playerID := "test_player_1"
	
	// 初始化玩家
	am.InitPlayer(playerID)
	
	// 检查玩家成就
	achievements := am.GetPlayerAchievements(playerID)
	if len(achievements) == 0 {
		t.Error("玩家应该有成就进度")
	}
}

func TestAchievementManager_UpdateProgress(t *testing.T) {
	am := NewAchievementManager()
	playerID := "test_player_2"
	
	// 初始化玩家
	am.InitPlayer(playerID)
	
	// 更新击杀进度
	completed := am.UpdateProgress(playerID, AchievementTypeKill, 50)
	if completed != nil {
		t.Error("50杀不应该完成100杀成就")
	}
	
	// 继续更新
	completed = am.UpdateProgress(playerID, AchievementTypeKill, 60)
	if completed == nil {
		t.Error("110杀应该完成100杀成就")
	}
	
	// 检查完成状态
	achievements := am.GetPlayerAchievements(playerID)
	var found bool
	for _, a := range achievements {
		if a.AchievementID == "killer_100" && a.Completed {
			found = true
			break
		}
	}
	if !found {
		t.Error("成就应该已标记为完成")
	}
}

func TestAchievementManager_GetCompletedCount(t *testing.T) {
	am := NewAchievementManager()
	playerID := "test_player_3"
	
	am.InitPlayer(playerID)
	am.UpdateProgress(playerID, AchievementTypeKill, 100)
	am.UpdateProgress(playerID, AchievementTypeScore, 10000)
	
	count := am.GetCompletedCount(playerID)
	if count < 2 {
		t.Errorf("期望至少2个成就完成, 实际为 %d", count)
	}
}

func TestAchievementManager_ClaimReward(t *testing.T) {
	am := NewAchievementManager()
	playerID := "test_player_4"
	
	am.InitPlayer(playerID)
	am.UpdateProgress(playerID, AchievementTypeKill, 100)
	
	// 领取奖励
	reward, err := am.ClaimReward(playerID, "killer_100")
	if err != nil {
		t.Errorf("领取奖励失败: %v", err)
	}
	
	if reward.Exp != 100 || reward.Coins != 50 {
		t.Errorf("奖励数据不正确:%+v", reward)
	}
	
	// 再次领取应该失败
	_, err = am.ClaimReward(playerID, "killer_100")
	if err == nil {
		t.Error("重复领取应该失败")
	}
}

// ============================================
// 任务系统测试
// ============================================

func TestQuestManager_Register(t *testing.T) {
	qm := NewQuestManager()
	
	// 测试注册任务
	if len(qm.Quests) == 0 {
		t.Error("任务管理器应该已注册任务")
	}
	
	// 测试获取任务
	quest := qm.GetDailyQuests("test_player")
	if len(quest) == 0 {
		t.Error("应该有每日任务")
	}
}

func TestQuestManager_InitPlayer(t *testing.T) {
	qm := NewQuestManager()
	playerID := "test_player_5"
	
	// 初始化玩家
	qm.InitPlayer(playerID)
	
	// 检查玩家任务
	quests := qm.GetPlayerQuests(playerID)
	if len(quests) == 0 {
		t.Error("玩家应该有任务进度")
	}
}

func TestQuestManager_AcceptQuest(t *testing.T) {
	qm := NewQuestManager()
	playerID := "test_player_6"
	
	qm.InitPlayer(playerID)
	
	// 接受任务
	err := qm.AcceptQuest(playerID, "daily_kill_10")
	if err != nil {
		t.Errorf("接受任务失败: %v", err)
	}
}

func TestQuestManager_UpdateProgress(t *testing.T) {
	qm := NewQuestManager()
	playerID := "test_player_7"
	
	qm.InitPlayer(playerID)
	qm.AcceptQuest(playerID, "daily_kill_10")
	
	// 更新进度
	qm.UpdateProgress(playerID, "kill_enemy", 5)
	
	// 检查进度
	quests := qm.GetPlayerQuests(playerID)
	var found *PlayerQuest
	for _, q := range quests {
		if q.QuestID == "daily_kill_10" {
			found = q
			break
		}
	}
	
	if found == nil {
		t.Error("应该找到任务")
	}
	
	if found.Progress != 5 {
		t.Errorf("期望进度为5, 实际为 %d", found.Progress)
	}
	
	// 继续更新到完成
	qm.UpdateProgress(playerID, "kill_enemy", 5)
	
	for _, q := range quests {
		if q.QuestID == "daily_kill_10" {
			found = q
			break
		}
	}
	
	if found.Status != QuestStatusCompleted {
		t.Error("任务应该已完成")
	}
}

func TestQuestManager_CompleteQuest(t *testing.T) {
	qm := NewQuestManager()
	playerID := "test_player_8"
	
	qm.InitPlayer(playerID)
	qm.AcceptQuest(playerID, "daily_kill_10")
	qm.UpdateProgress(playerID, "kill_enemy", 10)
	
	// 完成任务
	reward, err := qm.CompleteQuest(playerID, "daily_kill_10")
	if err != nil {
		t.Errorf("完成任务失败: %v", err)
	}
	
	if reward == nil {
		t.Error("应该返回奖励")
	}
}

// ============================================
// 玩家统计测试
// ============================================

func TestStatsManager_GetOrCreateStats(t *testing.T) {
	sm := NewStatsManager()
	playerID := "test_player_9"
	
	// 获取统计
	stats := sm.GetOrCreateStats(playerID)
	if stats == nil {
		t.Error("应该返回统计对象")
	}
	
	if stats.PlayerID != playerID {
		t.Errorf("玩家ID不匹配: %s", stats.PlayerID)
	}
}

func TestStatsManager_RecordGameStart(t *testing.T) {
	sm := NewStatsManager()
	playerID := "test_player_10"
	
	sm.RecordGameStart(playerID)
	
	stats := sm.GetStats(playerID)
	if stats.TotalGames != 1 {
		t.Errorf("期望游戏数为1, 实际为 %d", stats.TotalGames)
	}
}

func TestStatsManager_RecordGameWin(t *testing.T) {
	sm := NewStatsManager()
	playerID := "test_player_11"
	
	sm.RecordGameStart(playerID)
	sm.RecordGameWin(playerID, 50, 10000, 10)
	
	stats := sm.GetStats(playerID)
	if stats.WinGames != 1 {
		t.Error("应该记录胜利")
	}
	
	if stats.MaxKills != 50 {
		t.Errorf("期望最高击杀为50, 实际为 %d", stats.MaxKills)
	}
	
	if stats.MaxWave != 10 {
		t.Errorf("期望最高波次为10, 实际为 %d", stats.MaxWave)
	}
	
	if stats.MaxScore != 10000 {
		t.Errorf("期望最高分为10000, 实际为 %d", stats.MaxScore)
	}
}

func TestStatsManager_RecordGameLose(t *testing.T) {
	sm := NewStatsManager()
	playerID := "test_player_12"
	
	sm.RecordGameStart(playerID)
	sm.RecordGameLose(playerID, 30, 5000, 5)
	
	stats := sm.GetStats(playerID)
	if stats.LoseGames != 1 {
		t.Error("应该记录失败")
	}
	
	if stats.TotalGames != 1 {
		t.Errorf("期望总场次为1, 实际为 %d", stats.TotalGames)
	}
}

func TestStatsManager_RecordTowerOperations(t *testing.T) {
	sm := NewStatsManager()
	playerID := "test_player_13"
	
	sm.RecordTowerPlaced(playerID)
	sm.RecordTowerPlaced(playerID)
	sm.RecordTowerUpgraded(playerID)
	sm.RecordTowerSold(playerID)
	
	stats := sm.GetStats(playerID)
	if stats.TowersPlaced != 2 {
		t.Errorf("期望放置2座塔, 实际为 %d", stats.TowersPlaced)
	}
	
	if stats.TowersUpgraded != 1 {
		t.Errorf("期望升级1座塔, 实际为 %d", stats.TowersUpgraded)
	}
	
	if stats.TowersSold != 1 {
		t.Errorf("期望出售1座塔, 实际为 %d", stats.TowersSold)
	}
}

func TestStatsManager_RecordMoneyChange(t *testing.T) {
	sm := NewStatsManager()
	playerID := "test_player_14"
	
	sm.RecordMoneyChange(playerID, 100, 50)
	
	stats := sm.GetStats(playerID)
	if stats.TotalMoneyEarned != 100 {
		t.Errorf("期望获得100金币, 实际为 %d", stats.TotalMoneyEarned)
	}
	
	if stats.TotalMoneySpent != 50 {
		t.Errorf("期望花费50金币, 实际为 %d", stats.TotalMoneySpent)
	}
}

func TestStatsManager_RecordSocialStats(t *testing.T) {
	sm := NewStatsManager()
	playerID := "test_player_15"
	
	sm.RecordGiftReceived(playerID, 10)
	sm.RecordDanmakuReceived(playerID, 50)
	
	stats := sm.GetStats(playerID)
	if stats.TotalGiftsReceived != 10 {
		t.Errorf("期望收到10个礼物, 实际为 %d", stats.TotalGiftsReceived)
	}
	
	if stats.TotalDanmakuReceived != 50 {
		t.Errorf("期望收到50条弹幕, 实际为 %d", stats.TotalDanmakuReceived)
	}
}

// ============================================
// 基准测试
// ============================================

func BenchmarkAchievementManager_UpdateProgress(b *testing.B) {
	am := NewAchievementManager()
	playerID := "bench_player"
	am.InitPlayer(playerID)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		am.UpdateProgress(playerID, AchievementTypeKill, 1)
	}
}

func BenchmarkQuestManager_UpdateProgress(b *testing.B) {
	qm := NewQuestManager()
	playerID := "bench_player"
	qm.InitPlayer(playerID)
	qm.AcceptQuest(playerID, "daily_kill_10")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qm.UpdateProgress(playerID, "kill_enemy", 1)
	}
}

func BenchmarkStatsManager_RecordGameWin(b *testing.B) {
	sm := NewStatsManager()
	playerID := "bench_player"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.RecordGameStart(playerID)
		sm.RecordGameWin(playerID, 50, 10000, 10)
	}
}

// ============================================
// 示例代码
// ============================================

func ExampleAchievementManager() {
	// 创建成就管理器
	am := NewAchievementManager()
	
	// 初始化玩家
	playerID := "player_001"
	am.InitPlayer(playerID)
	
	// 模拟游戏过程 - 击杀敌人
	for i := 0; i < 120; i++ {
		completed := am.UpdateProgress(playerID, AchievementTypeKill, 1)
		if completed != nil {
			fmt.Printf("成就解锁: %s\n", completed.AchievementID)
		}
	}
	
	// 模拟游戏过程 - 获得分数
	am.UpdateProgress(playerID, AchievementTypeScore, 15000)
	
	// 查看完成数
	count := am.GetCompletedCount(playerID)
	fmt.Printf("已完成成就数: %d\n", count)
	
	// 领取奖励
	reward, _ := am.ClaimReward(playerID, "killer_100")
	fmt.Printf("领取奖励: Exp=%d, Coins=%d, Gems=%d\n", reward.Exp, reward.Coins, reward.Gems)
}

func ExampleQuestManager() {
	// 创建任务管理器
	qm := NewQuestManager()
	
	// 初始化玩家
	playerID := "player_002"
	qm.InitPlayer(playerID)
	
	// 接受每日任务
	qm.AcceptQuest(playerID, "daily_kill_10")
	
	// 模拟游戏过程 - 击杀敌人
	qm.UpdateProgress(playerID, "kill_enemy", 5)
	fmt.Println("任务进度: 5/10")
	
	// 继续击杀
	qm.UpdateProgress(playerID, "kill_enemy", 5)
	fmt.Println("任务进度: 10/10 (完成)")
	
	// 完成任务并领取奖励
	reward, _ := qm.CompleteQuest(playerID, "daily_kill_10")
	fmt.Printf("领取奖励: Exp=%d, Coins=%d\n", reward.Exp, reward.Coins)
}

func ExampleStatsManager() {
	// 创建统计管理器
	sm := NewStatsManager()
	
	// 玩家开始游戏
	playerID := "player_003"
	sm.RecordGameStart(playerID)
	
	// 模拟游戏过程
	sm.RecordGameWin(playerID, 50, 10000, 10)
	sm.RecordTowerPlaced(playerID)
	sm.RecordTowerPlaced(playerID)
	sm.RecordMoneyChange(playerID, 500, 200)
	sm.RecordGiftReceived(playerID, 5)
	sm.RecordDanmakuReceived(playerID, 20)
	
	// 获取统计
	stats := sm.GetStats(playerID)
	fmt.Printf("总场次: %d, 胜利: %d, 胜率: %.1f%%\n", 
		stats.TotalGames, stats.WinGames, stats.WinRate)
	fmt.Printf("总击杀: %d, 最高击杀: %d\n", stats.TotalKills, stats.MaxKills)
	fmt.Printf("最高波次: %d, 最高分: %d\n", stats.MaxWave, stats.MaxScore)
	fmt.Printf("放置塔数: %d\n", stats.TowersPlaced)
}
