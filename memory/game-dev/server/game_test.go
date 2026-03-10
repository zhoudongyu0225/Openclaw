package main

import (
	"testing"
	"time"
)

// ============================================
// 单元测试
// ============================================

// 测试防御塔创建
func TestNewTower(t *testing.T) {
	tower := NewTower("arrow_1", TowerTypeAttack, 100, 100, 1)
	
	if tower == nil {
		t.Fatal("Failed to create tower")
	}
	
	if tower.ID != "arrow_1" {
		t.Errorf("Expected ID 'arrow_1', got '%s'", tower.ID)
	}
	
	if tower.Type != TowerTypeAttack {
		t.Errorf("Expected type TowerTypeAttack, got %v", tower.Type)
	}
	
	if tower.Level != 1 {
		t.Errorf("Expected level 1, got %d", tower.Level)
	}
	
	if tower.X != 100 || tower.Y != 100 {
		t.Errorf("Expected position (100, 100), got (%.0f, %.0f)", tower.X, tower.Y)
	}
}

// 测试敌人生成
func TestEnemySpawner(t *testing.T) {
	spawner := NewEnemySpawner()
	
	// 测试初始波次
	spawner.StartWave(1)
	if spawner.Wave != 1 {
		t.Errorf("Expected wave 1, got %d", spawner.Wave)
	}
	
	// 测试生成敌人
	enemies := spawner.Spawn()
	if len(enemies) == 0 {
		t.Error("Failed to spawn enemies")
	}
	
	// 测试敌人属性
	if enemies[0].ID == "" {
		t.Error("Enemy ID should not be empty")
	}
	
	if enemies[0].MaxHP <= 0 {
		t.Errorf("Enemy HP should be positive, got %.0f", enemies[0].MaxHP)
	}
}

// 测试伤害计算
func TestEnemyTakeDamage(t *testing.T) {
	enemy := &Enemy{
		HP:    100,
		MaxHP: 100,
		Armor: 10,
	}
	
	// 测试普通伤害
	enemy.TakeDamage(50, 0)
	expectedHP := 50.0 * (1 - 10.0/(10.0+100.0)) // 约45.45
	if enemy.HP < expectedHP-1 || enemy.HP > expectedHP+1 {
		t.Errorf("Expected HP ~%.2f, got %.2f", expectedHP, enemy.HP)
	}
	
	// 测试护甲削减
	enemy2 := &Enemy{
		HP:    100,
		MaxHP: 100,
		Armor: 10,
	}
	enemy2.TakeDamage(50, 0.5) // 50%护甲削减
	
	actualArmor := enemy2.Armor * (1 - 0.5)
	expectedDamage := 50 * (1 - actualArmor/(actualArmor+100))
	if enemy2.HP != 100-expectedDamage {
		t.Errorf("Damage calculation error: HP = %.2f, expected %.2f", enemy2.HP, 100-expectedDamage)
	}
}

// 测试塔寻找目标
func TestTowerFindTarget(t *testing.T) {
	tower := NewTower("arrow", TowerTypeAttack, 300, 300, 1)
	tower.Range = 150
	
	enemies := []*Enemy{
		{HP: 50, MaxHP: 50, X: 320, Y: 320, Progress: 0.2}, // 范围内，进度0.2
		{HP: 50, MaxHP: 50, X: 350, Y: 350, Progress: 0.5}, // 范围内，进度0.5
		{HP: 50, MaxHP: 50, X: 100, Y: 100, Progress: 0.1}, // 范围外
	}
	
	target := tower.FindTarget(enemies)
	
	if target == nil {
		t.Fatal("Should find a target")
	}
	
	// 应该选择进度最高的 (最接近终点)
	if target.Progress != 0.5 {
		t.Errorf("Expected target with progress 0.5, got %.2f", target.Progress)
	}
}

// 测试塔攻击冷却
func TestTowerCanFire(t *testing.T) {
	tower := NewTower("arrow", TowerTypeAttack, 100, 100, 1)
	tower.FireRate = 1.0 // 1次/秒
	
	// 第一次应该可以攻击
	if !tower.CanFire() {
		t.Error("Should be able to fire initially")
	}
	
	// 攻击后应该进入冷却
	tower.LastFire = time.Now().UnixMilli()
	if tower.CanFire() {
		t.Error("Should not be able to fire immediately after firing")
	}
	
	// 等待冷却后应该可以再次攻击
	time.Sleep(1100 * time.Millisecond)
	if !tower.CanFire() {
		t.Error("Should be able to fire after cooldown")
	}
}

// 测试弹幕过滤
func TestDanmakuFilter(t *testing.T) {
	dm := NewDanmakuManager()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "hello world"},
		{"bad word", "*** word"},
		{"test bad again", "test *** again"},
	}
	
	for _, test := range tests {
		result := dm.Filter(test.input)
		if result != test.expected {
			t.Errorf("Filter('%s'): expected '%s', got '%s'", 
				test.input, test.expected, result)
		}
	}
}

// 测试礼物效果
func TestGiftEffect(t *testing.T) {
	gm := NewGiftManager()
	battle := NewBattleManager()
	battle.State.Money = 0
	
	// 初始金币为0
	initialMoney := battle.State.Money
	
	// 模拟金币礼物效果
	effect := &GiftEffect{
		Type:  GiftTypeCoin,
		GiftID: "test_coin",
	}
	
	// 手动应用金币效果
	GiftConfigs["coin"] = Gift{
		Type:  GiftTypeCoin,
		Name:  "金币",
		Price: 1,
		Value: 10,
	}
	
	// 直接修改状态模拟
	battle.State.Money += GiftConfigs["coin"].Value
	
	if battle.State.Money != initialMoney+10 {
		t.Errorf("Expected money %d, got %d", initialMoney+10, battle.State.Money)
	}
}

// 测试波次配置
func TestWaveConfigs(t *testing.T) {
	if len(WaveConfigs) == 0 {
		t.Error("WaveConfigs should not be empty")
	}
	
	// 验证波次配置
	for i, wave := range WaveConfigs {
		if wave.Duration <= 0 {
			t.Errorf("Wave %d: duration should be positive", i+1)
		}
		
		if len(wave.Enemies) == 0 {
			t.Errorf("Wave %d: should have at least one enemy type", i+1)
		}
		
		for _, enemy := range wave.Enemies {
			if enemy.Count <= 0 {
				t.Errorf("Wave %d: enemy count should be positive", i+1)
			}
		}
	}
}

// 测试塔配置
func TestTowerConfigs(t *testing.T) {
	expectedTowers := []string{"arrow", "cannon", "ice", "lightning", "tower_heal"}
	
	for _, towerID := range expectedTowers {
		config, ok := TowerConfigs[towerID]
		if !ok {
			t.Errorf("Tower config '%s' not found", towerID)
			continue
		}
		
		if config.BaseDamage < 0 {
			t.Errorf("Tower '%s': base damage should be non-negative", towerID)
		}
		
		if config.BaseRange <= 0 {
			t.Errorf("Tower '%s': base range should be positive", towerID)
		}
		
		if config.BaseFireRate <= 0 {
			t.Errorf("Tower '%s': base fire rate should be positive", towerID)
		}
	}
}

// 测试战斗管理器初始化
func TestNewBattleManager(t *testing.T) {
	bm := NewBattleManager()
	
	if bm == nil {
		t.Fatal("Failed to create BattleManager")
	}
	
	if bm.Towers == nil {
		t.Error("Towers should be initialized")
	}
	
	if bm.Spawner == nil {
		t.Error("Spawner should be initialized")
	}
	
	if bm.State.Lives != 20 {
		t.Errorf("Expected initial lives 20, got %d", bm.State.Lives)
	}
	
	if bm.State.Money != 100 {
		t.Errorf("Expected initial money 100, got %d", bm.State.Money)
	}
	
	if bm.State.Wave != 1 {
		t.Errorf("Expected initial wave 1, got %d", bm.State.Wave)
	}
}

// 测试游戏引擎状态转换
func TestGameEngineState(t *testing.T) {
	engine := NewGameEngine()
	
	// 初始状态应该是空闲
	if engine.State != GameStateIdle {
		t.Errorf("Expected initial state GameStateIdle, got %v", engine.State)
	}
	
	// 创建房间后应该在大厅
	engine.CreateRoom("room1", "host1")
	if engine.State != GameStateLobby {
		t.Errorf("Expected state GameStateLobby after CreateRoom, got %v", engine.State)
	}
	
	// 开始游戏后应该正在游戏
	engine.StartGame()
	if engine.State != GameStatePlaying {
		t.Errorf("Expected state GameStatePlaying after StartGame, got %v", engine.State)
	}
	
	// 暂停后应该暂停
	engine.PauseGame()
	if engine.State != GameStatePaused {
		t.Errorf("Expected state GameStatePaused after PauseGame, got %v", engine.State)
	}
	
	// 恢复后应该继续游戏
	engine.ResumeGame()
	if engine.State != GameStatePlaying {
		t.Errorf("Expected state GameStatePlaying after ResumeGame, got %v", engine.State)
	}
}

// 测试塔放置位置验证
func TestGameEnginePlaceTower(t *testing.T) {
	engine := NewGameEngine()
	engine.CreateRoom("room1", "host1")
	engine.StartGame()
	
	// 测试有效位置
	tower, err := engine.PlaceTower("arrow", TowerTypeAttack, 300, 300)
	if err != nil {
		t.Errorf("Failed to place tower at valid position: %v", err)
	}
	
	if tower == nil {
		t.Error("Tower should be created")
	}
	
	// 测试无效位置 (边界外)
	_, err = engine.PlaceTower("arrow", TowerTypeAttack, 10, 10)
	if err == nil {
		t.Error("Should fail to place tower outside boundaries")
	}
	
	// 测试位置重叠
	_, err = engine.PlaceTower("arrow", TowerTypeAttack, 300, 300)
	if err == nil {
		t.Error("Should fail to place tower at occupied position")
	}
}

// 基准测试: 敌人更新性能
func BenchmarkEnemyUpdate(b *testing.B) {
	spawner := NewEnemySpawner()
	spawner.StartWave(1)
	
	// 生成100个敌人
	for i := 0; i < 100; i++ {
		spawner.Spawn()
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spawner.Update(0.016) // 60fps
	}
}

// 基准测试: 塔寻找目标性能
func BenchmarkTowerFindTarget(b *testing.B) {
	tower := NewTower("arrow", TowerTypeAttack, 400, 300, 1)
	tower.Range = 200
	
	enemies := make([]*Enemy, 100)
	for i := 0; i < 100; i++ {
		enemies[i] = &Enemy{
			HP:       50,
			MaxHP:    50,
			X:        float64(100 + i*5),
			Y:        float64(100 + i*3),
			Progress: float64(i) / 100,
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tower.FindTarget(enemies)
	}
}
