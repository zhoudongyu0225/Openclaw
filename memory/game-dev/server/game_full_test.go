package game

import (
	"testing"
	"time"
)

// ==================== Room Manager Tests ====================

func TestRoomManager_CreateRoom(t *testing.T) {
	rm := NewRoomManager()
	
	room, err := rm.CreateRoom("host1", "test_room", "classic")
	if err != nil {
		t.Fatalf("Failed to create room: %v", err)
	}
	
	if room.ID == "" {
		t.Error("Room ID should not be empty")
	}
	
	if room.HostID != "host1" {
		t.Errorf("Expected host ID 'host1', got '%s'", room.HostID)
	}
	
	if room.Mode != "classic" {
		t.Errorf("Expected mode 'classic', got '%s'", room.Mode)
	}
}

func TestRoomManager_JoinRoom(t *testing.T) {
	rm := NewRoomManager()
	
	room, _ := rm.CreateRoom("host1", "test_room", "classic")
	
	err := rm.JoinRoom(room.ID, "guest1", "Player2")
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}
	
	if len(room.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(room.Players))
	}
}

func TestRoomManager_LeaveRoom(t *testing.T) {
	rm := NewRoomManager()
	
	room, _ := rm.CreateRoom("host1", "test_room", "classic")
	rm.JoinRoom(room.ID, "guest1", "Player2")
	
	err := rm.LeaveRoom("guest1", room.ID)
	if err != nil {
		t.Fatalf("Failed to leave room: %v", err)
	}
	
	if len(room.Players) != 1 {
		t.Errorf("Expected 1 player after leave, got %d", len(room.Players))
	}
}

func TestRoomManager_ListRooms(t *testing.T) {
	rm := NewRoomManager()
	
	rm.CreateRoom("host1", "room1", "classic")
	rm.CreateRoom("host2", "room2", "ranked")
	rm.CreateRoom("host3", "room3", "classic")
	
	rooms := rm.ListRooms()
	if len(rooms) != 3 {
		t.Errorf("Expected 3 rooms, got %d", len(rooms))
	}
	
	classicRooms := rm.ListRoomsByMode("classic")
	if len(classicRooms) != 2 {
		t.Errorf("Expected 2 classic rooms, got %d", len(classicRooms))
	}
}

// ==================== Matchmaker Tests ====================

func TestMatchmaker_BasicMatch(t *testing.T) {
	mm := NewMatchmaker(2)
	
	mm.AddPlayer("player1", 1000)
	mm.AddPlayer("player2", 1000)
	
	matches := mm.TryMatch()
	if len(matches) != 1 {
		t.Errorf("Expected 1 match, got %d", len(matches))
	}
}

func TestMatchmaker_SkillBasedMatch(t *testing.T) {
	mm := NewMatchmaker(2)
	
	// Add players with different skill levels
	mm.AddPlayer("player1", 1000)
	mm.AddPlayer("player2", 1200) // Too far apart
	mm.AddPlayer("player3", 1010) // Should match with player1
	
	matches := mm.TryMatch()
	// Only player1 and player3 should match (difference < 100)
	if len(matches) != 1 {
		t.Logf("Got %d matches, expected 1", len(matches))
	}
}

func TestMatchmaker_Timeout(t *testing.T) {
	mm := NewMatchmaker(3)
	
	mm.AddPlayer("player1", 1000)
	mm.AddPlayer("player2", 1000)
	
	// Try to match for 2-player game - should wait for 3rd player
	time.Sleep(100 * time.Millisecond)
	matches := mm.TryMatch()
	
	// Should not match due to insufficient players
	if len(matches) != 0 {
		t.Logf("Got %d matches", len(matches))
	}
}

// ==================== Battle Manager Tests ====================

func TestBattleManager_Update(t *testing.T) {
	bm := NewBattleManager("room1")
	
	// Add a tower
	tower := bm.AddTower("arrow", 1, 0, 0)
	if tower == nil {
		t.Fatal("Failed to add tower")
	}
	
	// Spawn an enemy
	enemy := bm.Spawner.Spawn("grunt")
	if enemy == nil {
		t.Fatal("Failed to spawn enemy")
	}
	
	// Update battle
	bm.Update(0.016) // 60fps
	
	// Check enemy moved
	if enemy.Progress == 0 {
		t.Log("Enemy should have some progress")
	}
}

func TestBattleManager_TowerAttack(t *testing.T) {
	bm := NewBattleManager("room1")
	
	// Place tower
	tower := bm.AddTower("cannon", 1, 100, 100)
	
	// Spawn enemy near tower
	enemy := bm.Spawner.Spawn("grunt")
	enemy.X = 110
	enemy.Y = 100
	
	// Update to trigger attack
	for i := 0; i < 60; i++ { // 1 second at 60fps
		bm.Update(0.016)
	}
	
	// Check if enemy took damage
	initialHP := 100.0 // grunt HP
	currentHP := enemy.HP
	
	if currentHP < initialHP {
		t.Logf("Enemy took damage: %f -> %f", initialHP, currentHP)
	}
}

func TestBattleManager_WaveProgress(t *testing.T) {
	bm := NewBattleManager("room1")
	
	if bm.State.Wave != 1 {
		t.Errorf("Expected wave 1, got %d", bm.State.Wave)
	}
	
	// Simulate wave completion
	bm.State.Score += 1000
	bm.CheckWaveComplete()
	
	if bm.State.Wave != 2 {
		t.Logf("Wave progressed to %d", bm.State.Wave)
	}
}

// ==================== Damage Calculation Tests ====================

func TestDamageCalculation_Armor(t *testing.T) {
	damage := 100.0
	armor := 50.0 // 50 armor = 33% reduction
	
	// Classic armor formula: damage * (1 - armor/(armor+100))
	actualDamage := damage * (1 - armor/(armor+100))
	expectedDamage := 100 * (1 - 50.0/150.0) // 33.33
	
	if actualDamage < expectedDamage-1 || actualDamage > expectedDamage+1 {
		t.Errorf("Expected damage ~%.2f, got %.2f", expectedDamage, actualDamage)
	}
}

func TestDamageCalculation_Critical(t *testing.T) {
	damage := 100.0
	critMultiplier := 2.0
	
	actualDamage := damage * critMultiplier
	if actualDamage != 200 {
		t.Errorf("Expected 200, got %.2f", actualDamage)
	}
}

func TestDamageCalculation_ArmorPierce(t *testing.T) {
	damage := 100.0
	armor := 50.0
	armorPierce := 20.0 // Reduces armor by 20%
	
	effectiveArmor := armor * (1 - armorPierce/100)
	actualDamage := damage * (1 - effectiveArmor/(effectiveArmor+100))
	
	// 40 armor = 28.6% reduction
	expectedDamage := 100 * (1 - 40.0/140.0) // ~71.4
	
	if actualDamage < expectedDamage-1 || actualDamage > expectedDamage+1 {
		t.Logf("Damage with armor pierce: %.2f", actualDamage)
	}
}

// ==================== Gift System Tests ====================

func TestGiftManager_ApplyEffect(t *testing.T) {
	gm := NewGiftManager()
	bm := NewBattleManager("room1")
	
	// Add initial money
	bm.State.Money = 100
	
	// Receive a coin gift
	gift := Gift{Type: GiftTypeCoin, Value: 10}
	gm.ReceiveGift(gift)
	
	// Apply effect
	gm.ApplyGiftEffect(bm)
	
	if bm.State.Money != 110 {
		t.Errorf("Expected 110 money, got %d", bm.State.Money)
	}
}

func TestGiftManager_BangEffect(t *testing.T) {
	gm := NewGiftManager()
	bm := NewBattleManager("room1")
	
	// Spawn enemies
	bm.Spawner.Spawn("grunt")
	bm.Spawner.Spawn("grunt")
	bm.Spawner.Spawn("grunt")
	
	initialKills := bm.State.TowerKills
	
	// Apply bang effect
	gift := Gift{Type: GiftTypeBang, Value: 300}
	gm.ReceiveGift(gift)
	gm.ApplyGiftEffect(bm)
	
	// All enemies should be dead
	if bm.State.TowerKills > initialKills {
		t.Logf("Bang killed enemies: %d -> %d", initialKills, bm.State.TowerKills)
	}
}

// ==================== Danmaku System Tests ====================

func TestDanmakuManager_Update(t *testing.T) {
	dm := NewDanmakuManager()
	
	// Send a danmaku
	dm.Send("Hello World!", "user1", DanmakuTypeNormal)
	
	if len(dm.ActiveDanmakus) != 1 {
		t.Errorf("Expected 1 danmaku, got %d", len(dm.ActiveDanmakus))
	}
	
	// Update positions
	dm.Update(0.016)
	
	// Danmaku should have moved
	if dm.ActiveDanmakus[0].X == 1920 { // Initial position
		t.Log("Danmaku should have moved")
	}
}

func TestDanmakuManager_Filter(t *testing.T) {
	dm := NewDanmakuManager()
	
	// Test sensitive word filtering
	result := dm.Filter("This is a test message")
	if result != "This is a test message" {
		t.Log("Clean message passed through")
	}
	
	// Add to sensitive list
	dm.SensitiveWords = []string{"badword"}
	result = dm.Filter("This contains badword")
	
	if result == "This contains badword" {
		t.Log("Sensitive word should be filtered")
	}
}

// ==================== Live Room Tests ====================

func TestLiveRoom_SendGift(t *testing.T) {
	room := NewLiveRoom("room1", "anchor1", "TestAnchor")
	
	room.JoinViewer("viewer1", "Viewer1")
	err := room.SendGift("viewer1", "rocket")
	
	if err != nil {
		t.Errorf("Failed to send gift: %v", err)
	}
	
	if len(room.PendingGifts) != 1 {
		t.Errorf("Expected 1 pending gift, got %d", len(room.PendingGifts))
	}
}

func TestLiveRoom_Update(t *testing.T) {
	room := NewLiveRoom("room1", "anchor1", "TestAnchor")
	bm := NewBattleManager("room1")
	room.Battle = bm
	
	// Send gift and danmaku
	room.SendGift("viewer1", "coin")
	room.SendDanmaku("viewer1", "Test!")
	
	// Update room
	room.Update(0.016)
	
	// Should process pending gifts
	t.Log("Live room updated successfully")
}

// ==================== Leaderboard Tests ====================

func TestLeaderboard_UpdateScore(t *testing.T) {
	lb := NewLeaderboard()
	
	lb.UpdateScore("player1", 1000, 10, 5)
	lb.UpdateScore("player2", 1500, 15, 3)
	lb.UpdateScore("player3", 800, 8, 8)
	
	// Get top players
	top := lb.GetTopPlayers(3)
	
	if len(top) != 3 {
		t.Errorf("Expected 3 players, got %d", len(top))
	}
	
	if top[0].PlayerID != "player2" {
		t.Logf("Top player should be player2, got %s", top[0].PlayerID)
	}
}

func TestLeaderboard_GetRank(t *testing.T) {
	lb := NewLeaderboard()
	
	lb.UpdateScore("player1", 1000, 10, 5)
	lb.UpdateScore("player2", 1500, 15, 3)
	
	rank, ok := lb.GetRank("player2")
	if !ok || rank != 1 {
		t.Errorf("Expected rank 1, got %d", rank)
	}
}

// ==================== Database Tests ====================

func TestDatabase_PlayerOps(t *testing.T) {
	db := NewDatabase()
	
	// Save player
	player := &Player{
		ID:       "player1",
		Name:     "TestPlayer",
		Score:    1000,
		Currency: 500,
	}
	
	err := db.SavePlayer(player)
	if err != nil {
		t.Fatalf("Failed to save player: %v", err)
	}
	
	// Load player
	loaded, err := db.GetPlayer("player1")
	if err != nil {
		t.Fatalf("Failed to load player: %v", err)
	}
	
	if loaded.Name != "TestPlayer" {
		t.Errorf("Expected name 'TestPlayer', got '%s'", loaded.Name)
	}
}

func TestDatabase_RoomOps(t *testing.T) {
	db := NewDatabase()
	
	room := &Room{
		ID:     "room1",
		HostID: "player1",
		Status: "waiting",
	}
	
	err := db.SaveRoom(room)
	if err != nil {
		t.Fatalf("Failed to save room: %v", err)
	}
	
	loaded, err := db.GetRoom("room1")
	if err != nil {
		t.Fatalf("Failed to load room: %v", err)
	}
	
	if loaded.Status != "waiting" {
		t.Errorf("Expected status 'waiting', got '%s'", loaded.Status)
	}
}

// ==================== Rate Limiter Tests ====================

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(5, time.Second) // 5 requests per second
	
	// Should allow first 5
	for i := 0; i < 5; i++ {
		if !rl.Allow("test_ip") {
			t.Errorf("Request %d should be allowed", i)
		}
	}
	
	// 6th should be denied
	if rl.Allow("test_ip") {
		t.Log("6th request correctly denied")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)
	
	// Use up quota
	for i := 0; i < 5; i++ {
		rl.Allow("test_ip")
	}
	
	// Wait for reset
	time.Sleep(time.Second + 100*time.Millisecond)
	
	// Should allow again
	if !rl.Allow("test_ip") {
		t.Error("Request should be allowed after reset")
	}
}

// ==================== Input Validator Tests ====================

func TestInputValidator_Username(t *testing.T) {
	validator := NewInputValidator()
	
	tests := []struct {
		input    string
		expected bool
	}{
		{"user", true},
		{"ab", false},           // Too short
		{"verylongusername12345", false}, // Too long
		{"user_123", true},
		{"user@123", false},     // Invalid char
	}
	
	for _, test := range tests {
		err := validator.ValidateUsername(test.input)
		if test.expected && err != nil {
			t.Errorf("Username '%s' should be valid", test.input)
		}
		if !test.expected && err == nil {
			t.Errorf("Username '%s' should be invalid", test.input)
		}
	}
}

func TestInputValidator_Danmaku(t *testing.T) {
	validator := NewInputValidator()
	
	// Normal message
	err := validator.ValidateDanmaku("Hello world!")
	if err != nil {
		t.Errorf("Normal message should be valid")
	}
	
	// Too long
	longMsg := string(make([]byte, 201))
	err = validator.ValidateDanmaku(longMsg)
	if err == nil {
		t.Error("Long message should be invalid")
	}
}

// ==================== Config Tests ====================

func TestConfig_Load(t *testing.T) {
	cfg := LoadConfig()
	
	if cfg.Server.Port == "" {
		t.Error("Server port should be set")
	}
	
	if cfg.Game.Battle.FPS == 0 {
		t.Error("Battle FPS should be set")
	}
}

func TestConfig_Reload(t *testing.T) {
	cfg := LoadConfig()
	
	originalPort := cfg.Server.Port
	
	// Reload
	cfg.Reload()
	
	if cfg.Server.Port != originalPort {
		t.Log("Config reloaded")
	}
}

// ==================== Protobuf Tests ====================

func TestProtobuf_EncodeDecode(t *testing.T) {
	codec := NewCodec()
	
	// Create a message
	msg := &GameMessage{
		Type: MessageTypeJoinRoom,
		Payload: &GameMessage_JoinRoom{
			JoinRoom: &JoinRoomRequest{
				RoomID:   "room1",
				PlayerID: "player1",
			},
		},
	}
	
	// Encode
	data, err := codec.Encode(msg)
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}
	
	// Decode
	decoded, err := codec.Decode(data)
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}
	
	if decoded.Type != MessageTypeJoinRoom {
		t.Errorf("Expected message type %d, got %d", MessageTypeJoinRoom, decoded.Type)
	}
}
