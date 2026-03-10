package game

import (
	"testing"
)

// TestEnhancementManager_Enhance 测试装备强化
func TestEnhancementManager_Enhance(t *testing.T) {
	// 创建测试玩家
	player := &Player{
		ID:     "test_player_1",
		Name:   "TestPlayer",
		Level:  10,
		Gold:   10000,
		Gem:    1000,
		Exp:    0,
		Health: 1000,
		Inventory: &Inventory{
			Items: make([]*Item, 0),
			Capacity: 100,
		},
	}

	// 创建测试装备
	item := &Item{
		ID:            "weapon_test",
		Name:          "测试武器",
		Type:          ItemTypeWeapon,
		Quality:       1,
		Level:         1,
		EnhanceLevel:  0,
		Durability:    100,
		MaxDurability: 100,
		Attack:        10,
		Defense:       5,
		HP:            0,
	}

	player.AddItem(item)

	// 测试强化
	result, err := enhancementMgr.Enhance(player, item, "atk_1")
	if err != nil {
		t.Errorf("强化失败: %v", err)
	}

	t.Logf("强化结果: %+v", result)
}

// TestSynthesisManager_Synthesize 测试装备合成
func TestSynthesisManager_Synthesize(t *testing.T) {
	player := &Player{
		ID:     "test_player_2",
		Name:   "TestPlayer2",
		Level:  5,
		Gold:   1000,
		Inventory: &Inventory{
			Items:    make([]*Item, 0),
			Capacity: 100,
		},
	}

	// 添加测试材料
	materials := []struct {
		id    string
		name  string
		count int
	}{
		{"iron_ingot", "铁锭", 5},
		{"wood", "木材", 3},
	}

	for _, m := range materials {
		item := &Item{
			ID:     m.id,
			Name:   m.name,
			Type:   ItemTypeMaterial,
			Count:  m.count,
		}
		player.AddItem(item)
	}

	// 测试合成
	result, err := synthesisMgr.Synthesize(player, "recipe_weapon_1")
	if err != nil {
		t.Errorf("合成失败: %v", err)
	}

	t.Logf("合成结果: %+v", result)
}

// TestEquipmentEnhancementManager_RefineEquipment 测试装备精炼
func TestEquipmentEnhancementManager_RefineEquipment(t *testing.T) {
	player := &Player{
		ID:     "test_player_3",
		Name:   "TestPlayer3",
		Level:  10,
		Gold:   5000,
		Gem:    500,
		Inventory: &Inventory{
			Items:    make([]*Item, 0),
			Capacity: 100,
		},
	}

	item := &Item{
		ID:            "armor_test",
		Name:          "测试装甲",
		Type:          ItemTypeArmor,
		Quality:       1,
		Level:         5,
		Durability:    100,
		MaxDurability: 100,
		Attack:        0,
		Defense:       20,
		HP:            50,
	}

	player.AddItem(item)

	result, err := enhancementMgr.RefineEquipment(player, item)
	if err != nil {
		t.Errorf("精炼失败: %v", err)
	}

	t.Logf("精炼结果: %+v", result)
}

// TestEquipmentEnhancementManager_Enchant 测试装备附魔
func TestEquipmentEnhancementManager_Enchant(t *testing.T) {
	player := &Player{
		ID:     "test_player_4",
		Name:   "TestPlayer4",
		Level:  10,
		Gold:   5000,
		Gem:    2000,
		Inventory: &Inventory{
			Items:    make([]*Item, 0),
			Capacity: 100,
		},
	}

	item := &Item{
		ID:            "weapon_enchant_test",
		Name:          "附魔测试武器",
		Type:          ItemTypeWeapon,
		Quality:       3,
		Level:         10,
		Durability:    100,
		MaxDurability: 100,
		Attack:        50,
		Enchantments:  []string{},
	}

	player.AddItem(item)

	result, err := enhancementMgr.Enchant(player, item, "ench_fire")
	if err != nil {
		t.Errorf("附魔失败: %v", err)
	}

	t.Logf("附魔结果: %+v", result)
}

// TestEquipmentSynthesisManager_DisassembleEquipment 测试装备分解
func TestEquipmentSynthesisManager_DisassembleEquipment(t *testing.T) {
	player := &Player{
		ID:     "test_player_5",
		Name:   "TestPlayer5",
		Level:  5,
		Gold:   0,
		Inventory: &Inventory{
			Items:    make([]*Item, 0),
			Capacity: 100,
		},
	}

	item := &Item{
		ID:            "weapon_disassemble",
		Name:          "待分解武器",
		Type:          ItemTypeWeapon,
		Quality:       3,
		Level:         5,
		Durability:    50,
		MaxDurability: 100,
		Attack:        30,
	}

	player.AddItem(item)

	result, err := synthesisMgr.DisassembleEquipment(player, item)
	if err != nil {
		t.Errorf("分解失败: %v", err)
	}

	t.Logf("分解结果: %+v", result)
}

// TestEquipmentEnhancementManager_AddSocket 测试装备打孔
func TestEquipmentEnhancementManager_AddSocket(t *testing.T) {
	player := &Player{
		ID:     "test_player_6",
		Name:   "TestPlayer6",
		Level:  10,
		Gold:   5000,
		Gem:    1000,
		Inventory: &Inventory{
			Items:    make([]*Item, 0),
			Capacity: 100,
		},
	}

	item := &Item{
		ID:            "weapon_socket_test",
		Name:          "打孔测试武器",
		Type:          ItemTypeWeapon,
		Quality:       2,
		Level:         5,
		Durability:    100,
		MaxDurability: 100,
		Sockets:       []Socket{},
	}

	player.AddItem(item)

	result, err := enhancementMgr.AddSocket(player, item)
	if err != nil {
		t.Errorf("打孔失败: %v", err)
	}

	t.Logf("打孔结果: %+v", result)
}

// TestEquipmentEnhancementManager_RepairEquipment 测试装备修理
func TestEquipmentEnhancementManager_RepairEquipment(t *testing.T) {
	player := &Player{
		ID:     "test_player_7",
		Name:   "TestPlayer7",
		Level:  5,
		Gold:   1000,
		Inventory: &Inventory{
			Items:    make([]*Item, 0),
			Capacity: 100,
		},
	}

	item := &Item{
		ID:            "weapon_repair_test",
		Name:          "待修理武器",
		Type:          ItemTypeWeapon,
		Quality:       1,
		Level:         1,
		Durability:    30,
		MaxDurability: 100,
		Attack:        10,
	}

	player.AddItem(item)

	cost, err := enhancementMgr.RepairEquipment(player, item)
	if err != nil {
		t.Errorf("修理失败: %v", err)
	}

	t.Logf("修理费用: %d, 装备耐久度: %d/%d", cost, item.Durability, item.MaxDurability)
}

// TestEquipmentEnhancementManager_CalculateRefineCost 测试精炼费用计算
func TestEquipmentEnhancementManager_CalculateRefineCost(t *testing.T) {
	testCases := []struct {
		quality int
		expected int
	}{
		{1, 100},
		{2, 200},
		{3, 300},
		{5, 500},
		{10, 1000},
	}

	for _, tc := range testCases {
		item := &Item{Quality: tc.quality}
		cost := enhancementMgr.CalculateRefineCost(item)
		if cost != tc.expected {
			t.Errorf("品质 %d 的精炼费用应为 %d，实际为 %d", tc.quality, tc.expected, cost)
		}
	}
}

// TestEquipmentEnhancementManager_GetRefineSuccessRate 测试精炼成功率
func TestEquipmentEnhancementManager_GetRefineSuccessRate(t *testing.T) {
	testCases := []struct {
		quality int
		expected float64
	}{
		{1, 0.75},
		{2, 0.7},
		{3, 0.65},
		{5, 0.55},
		{10, 0.3},
	}

	for _, tc := range testCases {
		item := &Item{Quality: tc.quality}
		rate := enhancementMgr.GetRefineSuccessRate(item)
		if rate != tc.expected {
			t.Errorf("品质 %d 的精炼成功率应为 %f，实际为 %f", tc.quality, tc.expected, rate)
		}
	}
}

// TestEquipmentSynthesisManager_CanSynthesize 测试合成条件检查
func TestEquipmentSynthesisManager_CanSynthesize(t *testing.T) {
	player := &Player{
		ID:     "test_player_8",
		Name:   "TestPlayer8",
		Level:  10,
		Gold:   500,
		Inventory: &Inventory{
			Items:    make([]*Item, 0),
			Capacity: 100,
		},
	}

	// 材料不足
	can, msg := synthesisMgr.CanSynthesize(player, "recipe_weapon_1")
	if can {
		t.Errorf("材料不足时应该返回 false")
	}
	t.Logf("检查结果: %s - %s", can, msg)

	// 添加材料
	player.AddItem(&Item{ID: "iron_ingot", Name: "铁锭", Type: ItemTypeMaterial, Count: 5})
	player.AddItem(&Item{ID: "wood", Name: "木材", Type: ItemTypeMaterial, Count: 3})

	// 金币不足
	can, msg = synthesisMgr.CanSynthesize(player, "recipe_weapon_1")
	if can {
		t.Errorf("金币不足时应该返回 false")
	}
	t.Logf("检查结果: %s - %s", can, msg)

	// 满足条件
	player.Gold = 1000
	can, msg = synthesisMgr.CanSynthesize(player, "recipe_weapon_1")
	if !can {
		t.Errorf("满足条件时应该返回 true")
	}
	t.Logf("检查结果: %s - %s", can, msg)
}

// TestListEnhancements 测试列出所有强化类型
func TestListEnhancements(t *testing.T) {
	enhancements := enhancementMgr.ListEnhancements()
	t.Logf("强化类型数量: %d", len(enhancements))
	for _, e := range enhancements {
		t.Logf("  - %s: %s (%s)", e.ID, e.Name, e.Description)
	}
}

// TestListEnchantments 测试列出所有附魔类型
func TestListEnchantments(t *testing.T) {
	enchantments := enhancementMgr.ListEnchantments()
	t.Logf("附魔类型数量: %d", len(enchantments))
	for _, e := range enchantments {
		t.Logf("  - %s: %s (%s) - 稀有度: %d", e.ID, e.Name, e.Description, e.Rarity)
	}
}

// TestListRecipes 测试列出所有合成配方
func TestListRecipes(t *testing.T) {
	recipes := synthesisMgr.ListRecipes()
	t.Logf("合成配方数量: %d", len(recipes))
	for _, r := range recipes {
		t.Logf("  - %s: %s (%s)", r.ID, r.Name, r.Description)
	}
}
