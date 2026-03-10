package item

import (
	"sync"
	"time"
)

// ItemType 道具类型
type ItemType int

const (
	ItemTypeWeapon ItemType = iota + 1
	ItemTypeArmor
	ItemTypeAccessory
	ItemTypeConsumable
	ItemTypeMaterial
	ItemTypeCurrency
	ItemTypeTicket
)

// Rarity 稀有度
type Rarity int

const (
	RarityCommon Rarity = iota + 1
	RarityUncommon
	RarityRare
	RarityEpic
	RarityLegendary
)

// Item 道具结构
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        ItemType  `json:"type"`
	Rarity      Rarity    `json:"rarity"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Level       int       `json:"level"`
	Stats       ItemStats `json:"stats"`
	Icon        string    `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
}

// ItemStats 道具属性
type ItemStats struct {
	Attack    int `json:"attack"`
	Defense   int `json:"defense"`
	HP        int `json:"hp"`
	Speed     int `json:"speed"`
	Critical  int `json:"critical"`
	Evasion   int `json:"evasion"`
	Recover   int `json:"recover"`
}

// WeaponStats 武器专属属性
type WeaponStats struct {
	Damage       int `json:"damage"`
	AttackSpeed int `json:"attack_speed"`
	Range        int `json:"range"`
	Penetration  int `json:"penetration"`
}

// ArmorStats 护甲专属属性
type ArmorStats struct {
	Armor       int `json:"armor"`
	Durability  int `json:"durability"`
	Weight      int `json:"weight"`
	Resistance  int `json:"resistance"`
}

// PlayerInventory 玩家背包
type PlayerInventory struct {
	PlayerID    string             `json:"player_id"`
	Items       map[string]*InventoryItem `json:"items"`
	Coins       int                `json:"coins"`
	Gems        int                `json:"gems"`
	Capacity    int                `json:"capacity"`
	Mu          sync.RWMutex       `json:"-"`
}

// InventoryItem 背包中的道具
type InventoryItem struct {
	ItemID     string    `json:"item_id"`
	Item       *Item     `json:"item"`
	Count      int       `json:"count"`
	Level      int       `json:"level"`
	Exp        int       `json:"exp"`
	AcquiredAt time.Time `json:"acquired_at"`
	IsEquipped bool      `json:"is_equipped"`
}

// NewPlayerInventory 创建新背包
func NewPlayerInventory(playerID string, capacity int) *PlayerInventory {
	return &PlayerInventory{
		PlayerID: playerID,
		Items:    make(map[string]*InventoryItem),
		Coins:    100,
		Gems:     0,
		Capacity: capacity,
	}
}

// AddItem 添加道具
func (inv *PlayerInventory) AddItem(item *Item, count int) error {
	inv.Mu.Lock()
	defer inv.Mu.Unlock()

	if len(inv.Items) >= inv.Capacity && item.Type != ItemTypeCurrency {
		return ErrInventoryFull
	}

	if existing, ok := inv.Items[item.ID]; ok {
		existing.Count += count
	} else {
		inv.Items[item.ID] = &InventoryItem{
			ItemID:     item.ID,
			Item:       item,
			Count:      count,
			Level:      1,
			Exp:        0,
			AcquiredAt: time.Now(),
			IsEquipped: false,
		}
	}

	return nil
}

// RemoveItem 移除道具
func (inv *PlayerInventory) RemoveItem(itemID string, count int) error {
	inv.Mu.Lock()
	defer inv.Mu.Unlock()

	item, ok := inv.Items[itemID]
	if !ok {
		return ErrItemNotFound
	}

	if item.Count < count {
		return ErrInsufficientItem
	}

	item.Count -= count
	if item.Count <= 0 {
		delete(inv.Items, itemID)
	}

	return nil
}

// EquipItem 装备道具
func (inv *PlayerInventory) EquipItem(itemID string) error {
	inv.Mu.Lock()
	defer inv.Mu.Unlock()

	item, ok := inv.Items[itemID]
	if !ok {
		return ErrItemNotFound
	}

	if item.Item.Type != ItemTypeWeapon && item.Item.Type != ItemTypeArmor && item.Item.Type != ItemTypeAccessory {
		return ErrCannotEquip
	}

	item.IsEquipped = true
	return nil
}

// UnequipItem 卸下装备
func (inv *PlayerInventory) UnequipItem(itemID string) error {
	inv.Mu.Lock()
	defer inv.Mu.Unlock()

	item, ok := inv.Items[itemID]
	if !ok {
		return ErrItemNotFound
	}

	item.IsEquipped = false
	return nil
}

// GetEquippedItems 获取已装备道具
func (inv *PlayerInventory) GetEquippedItems() []*InventoryItem {
	inv.Mu.RLock()
	defer inv.Mu.RUnlock()

	var equipped []*InventoryItem
	for _, item := range inv.Items {
		if item.IsEquipped {
			equipped = append(equipped, item)
		}
	}

	return equipped
}

// GetTotalStats 获取总属性
func (inv *PlayerInventory) GetTotalStats() ItemStats {
	inv.Mu.RLock()
	defer inv.Mu.RUnlock()

	var total ItemStats
	for _, item := range inv.Items {
		if item.IsEquipped {
			total.Attack += item.Item.Stats.Attack
			total.Defense += item.Item.Stats.Defense
			total.HP += item.Item.Stats.HP
			total.Speed += item.Item.Stats.Speed
			total.Critical += item.Item.Stats.Critical
			total.Evasion += item.Item.Stats.Evasion
			total.Recover += item.Item.Stats.Recover
		}
	}

	return total
}

// HasItem 检查是否拥有道具
func (inv *PlayerInventory) HasItem(itemID string) bool {
	inv.Mu.RLock()
	defer inv.Mu.RUnlock()

	_, ok := inv.Items[itemID]
	return ok
}

// GetItemCount 获取道具数量
func (inv *PlayerInventory) GetItemCount(itemID string) int {
	inv.Mu.RLock()
	defer inv.Mu.RUnlock()

	if item, ok := inv.Items[itemID]; ok {
		return item.Count
	}
	return 0
}

// AddCoins 添加金币
func (inv *PlayerInventory) AddCoins(amount int) {
	inv.Mu.Lock()
	defer inv.Mu.Unlock()
	inv.Coins += amount
}

// DeductCoins 扣除金币
func (inv *PlayerInventory) DeductCoins(amount int) error {
	inv.Mu.Lock()
	defer inv.Mu.Unlock()

	if inv.Coins < amount {
		return ErrInsufficientCoins
	}

	inv.Coins -= amount
	return nil
}

// AddGems 添加钻石
func (inv *PlayerInventory) AddGems(amount int) {
	inv.Mu.Lock()
	defer inv.Mu.Unlock()
	inv.Gems += amount
}

// DeductGems 扣除钻石
func (inv *PlayerInventory) DeductGems(amount int) error {
	inv.Mu.Lock()
	defer inv.Mu.Unlock()

	if inv.Gems < amount {
		return ErrInsufficientGems
	}

	inv.Gems -= amount
	return nil
}

// InventoryError 背包错误
type InventoryError string

func (e InventoryError) Error() string {
	return string(e)
}

const (
	ErrInventoryFull       InventoryError = "inventory is full"
	ErrItemNotFound         InventoryError = "item not found"
	ErrInsufficientItem     InventoryError = "insufficient item quantity"
	ErrCannotEquip          InventoryError = "this item cannot be equipped"
	ErrInsufficientCoins    InventoryError = "insufficient coins"
	ErrInsufficientGems    InventoryError = "insufficient gems"
)

// ItemManager 道具管理器
type ItemManager struct {
	items map[string]*Item
	mu    sync.RWMutex
}

// NewItemManager 创建道具管理器
func NewItemManager() *ItemManager {
	m := &ItemManager{
		items: make(map[string]*Item),
	}
	m.initDefaultItems()
	return m
}

// initDefaultItems 初始化默认道具
func (m *ItemManager) initDefaultItems() {
	defaultItems := []*Item{
		// 武器
		{ID: "weapon_001", Name: "新手剑", Type: ItemTypeWeapon, Rarity: RarityCommon, Price: 100, Level: 1, Stats: ItemStats{Attack: 10}},
		{ID: "weapon_002", Name: "铁剑", Type: ItemTypeWeapon, Rarity: RarityUncommon, Price: 300, Level: 5, Stats: ItemStats{Attack: 25}},
		{ID: "weapon_003", Name: "钢剑", Type: ItemTypeWeapon, Rarity: RarityRare, Price: 800, Level: 10, Stats: ItemStats{Attack: 50}},
		{ID: "weapon_004", Name: "魔法剑", Type: ItemTypeWeapon, Rarity: RarityEpic, Price: 2000, Level: 20, Stats: ItemStats{Attack: 100, Critical: 10}},
		{ID: "weapon_005", Name: "传奇之剑", Type: ItemTypeWeapon, Rarity: RarityLegendary, Price: 10000, Level: 50, Stats: ItemStats{Attack: 250, Critical: 25, Speed: 10}},
		// 护甲
		{ID: "armor_001", Name: "布衣", Type: ItemTypeArmor, Rarity: RarityCommon, Price: 80, Level: 1, Stats: ItemStats{Defense: 5}},
		{ID: "armor_002", Name: "皮甲", Type: ItemTypeArmor, Rarity: RarityUncommon, Price: 200, Level: 5, Stats: ItemStats{Defense: 15, HP: 50}},
		{ID: "armor_003", Name: "锁甲", Type: ItemTypeArmor, Rarity: RarityRare, Price: 600, Level: 10, Stats: ItemStats{Defense: 35, HP: 100}},
		{ID: "armor_004", Name: "板甲", Type: ItemTypeArmor, Rarity: RarityEpic, Price: 1500, Level: 20, Stats: ItemStats{Defense: 75, HP: 250}},
		{ID: "armor_005", Name: "龙鳞甲", Type: ItemTypeArmor, Rarity: RarityLegendary, Price: 8000, Level: 50, Stats: ItemStats{Defense: 200, HP: 500, Evasion: 10}},
		// 饰品
		{ID: "acc_001", Name: "力量戒指", Type: ItemTypeAccessory, Rarity: RarityCommon, Price: 150, Level: 1, Stats: ItemStats{Attack: 5}},
		{ID: "acc_002", Name: "防御戒指", Type: ItemTypeAccessory, Rarity: RarityCommon, Price: 150, Level: 1, Stats: ItemStats{Defense: 5}},
		{ID: "acc_003", Name: "生命戒指", Type: ItemTypeAccessory, Rarity: RarityUncommon, Price: 400, Level: 5, Stats: ItemStats{HP: 100}},
		{ID: "acc_004", Name: "暴击戒指", Type: ItemTypeAccessory, Rarity: RarityRare, Price: 1000, Level: 15, Stats: ItemStats{Critical: 15}},
		{ID: "acc_005", Name: "极速戒指", Type: ItemTypeAccessory, Rarity: RarityEpic, Price: 2500, Level: 25, Stats: ItemStats{Speed: 20, Evasion: 10}},
		// 消耗品
		{ID: "potion_hp_001", Name: "小血瓶", Type: ItemTypeConsumable, Rarity: RarityCommon, Price: 10, Stats: ItemStats{Recover: 50}},
		{ID: "potion_hp_002", Name: "中血瓶", Type: ItemTypeConsumable, Rarity: RarityUncommon, Price: 50, Stats: ItemStats{Recover: 200}},
		{ID: "potion_hp_003", Name: "大血瓶", Type: ItemTypeConsumable, Rarity: RarityRare, Price: 200, Stats: ItemStats{Recover: 500}},
		{ID: "potion_exp_001", Name: "经验药水", Type: ItemTypeConsumable, Rarity: RarityUncommon, Price: 100, Stats: ItemStats{Exp: 100}},
		// 材料
		{ID: "material_001", Name: "铁矿石", Type: ItemTypeMaterial, Rarity: RarityCommon, Price: 5},
		{ID: "material_002", Name: "魔法水晶", Type: ItemTypeMaterial, Rarity: RarityRare, Price: 50},
		{ID: "material_003", Name: "龙之心", Type: ItemTypeMaterial, Rarity: RarityLegendary, Price: 500},
		// 货币
		{ID: "gold_pack_001", Name: "金币包(100)", Type: ItemTypeCurrency, Rarity: RarityCommon, Price: 10},
		{ID: "gem_pack_001", Name: "钻石包(10)", Type: ItemTypeCurrency, Rarity: RarityCommon, Price: 100},
	}

	for _, item := range defaultItems {
		m.items[item.ID] = item
	}
}

// GetItem 获取道具
func (m *ItemManager) GetItem(itemID string) (*Item, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[itemID]
	return item, ok
}

// GetAllItems 获取所有道具
func (m *ItemManager) GetAllItems() []*Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := make([]*Item, 0, len(m.items))
	for _, item := range m.items {
		items = append(items, item)
	}

	return items
}

// GetItemsByType 按类型获取道具
func (m *ItemManager) GetItemsByType(itemType ItemType) []*Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var items []*Item
	for _, item := range m.items {
		if item.Type == itemType {
			items = append(items, item)
		}
	}

	return items
}

// GetItemsByRarity 按稀有度获取道具
func (m *ItemManager) GetItemsByRarity(rarity Rarity) []*Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var items []*Item
	for _, item := range m.items {
		if item.Rarity == rarity {
			items = append(items, item)
		}
	}

	return items
}

// AddCustomItem 添加自定义道具
func (m *ItemManager) AddCustomItem(item *Item) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[item.ID] = item
}
