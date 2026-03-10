package shop

import (
	"sync"
	"time"
)

// ShopItem represents an item in the shop
type ShopItem struct {
	ID          string    `json:"id"`
	ItemID      string    `json:"item_id"`
	ItemType    string    `json:"item_type"` // weapon, armor, accessory, consumable, material
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	PriceType   string    `json:"price_type"` // coins, gems
	Stock       int       `json:"stock"`      // -1 for unlimited
	Sold        int       `json:"sold"`
	LevelReq    int       `json:"level_req"` // minimum player level
	Discount    float64   `json:"discount"`  // 0.0 - 1.0
	Category    string    `json:"category"`  // weapons, armor, accessories, consumables, materials
	Sort        int       `json:"sort"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsActive    bool      `json:"is_active"`
}

// PurchaseRecord represents a purchase record
type PurchaseRecord struct {
	ID        string    `json:"id"`
	PlayerID  string    `json:"player_id"`
	ItemID    string    `json:"item_id"`
	ShopItemID string   `json:"shop_item_id"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	PriceType string    `json:"price_type"`
	Time      time.Time `json:"time"`
}

// ShopManager manages the shop
type ShopManager struct {
	mu       sync.RWMutex
	items    map[string]*ShopItem // itemID -> shop item
	categories map[string][]string // category -> itemIDs
	purchases map[string][]*PurchaseRecord // playerID -> records
	dailyLimit map[string]map[string]int // playerID -> (itemID -> count)
}

// NewShopManager creates a new shop manager
func NewShopManager() *ShopManager {
	return &ShopManager{
		items:      make(map[string]*ShopItem),
		categories: make(map[string][]string),
		purchases:  make(map[string][]*PurchaseRecord),
		dailyLimit: make(map[string]map[string]int),
	}
}

// AddItem adds an item to the shop
func (sm *ShopManager) AddItem(item *ShopItem) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if item.ID == "" || item.ItemID == "" {
		return ErrInvalidItem
	}

	item.IsActive = true
	sm.items[item.ID] = item
	sm.categories[item.Category] = append(sm.categories[item.Category], item.ID)

	return nil
}

// AddItems adds multiple items to the shop
func (sm *ShopManager) AddItems(items []*ShopItem) error {
	for _, item := range items {
		if err := sm.AddItem(item); err != nil {
			return err
		}
	}
	return nil
}

// GetItem gets a shop item by ID
func (sm *ShopManager) GetItem(itemID string) (*ShopItem, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	item, ok := sm.items[itemID]
	if !ok {
		return nil, ErrItemNotFound
	}
	return item, nil
}

// GetItemsByCategory gets all items in a category
func (sm *ShopManager) GetItemsByCategory(category string) []*ShopItem {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	itemIDs := sm.categories[category]
	items := make([]*ShopItem, 0, len(itemIDs))
	now := time.Now()

	for _, id := range itemIDs {
		if item, exists := sm.items[id]; exists && item.IsActive {
			// Check time range
			if !item.StartTime.IsZero() && now.Before(item.StartTime) {
				continue
			}
			if !item.EndTime.IsZero() && now.After(item.EndTime) {
				continue
			}
			items = append(items, item)
		}
	}

	return items
}

// GetAllItems gets all active items
func (sm *ShopManager) GetAllItems() []*ShopItem {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	items := make([]*ShopItem, 0)
	now := time.Now()

	for _, item := range sm.items {
		if !item.IsActive {
			continue
		}
		// Check time range
		if !item.StartTime.IsZero() && now.Before(item.StartTime) {
			continue
		}
		if !item.EndTime.IsZero() && now.After(item.EndTime) {
			continue
		}
		items = append(items, item)
	}

	return items
}

// Purchase buys an item
func (sm *ShopManager) Purchase(playerID, itemID string, quantity int) (*PurchaseRecord, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	item, ok := sm.items[itemID]
	if !ok {
		return nil, ErrItemNotFound
	}

	if !item.IsActive {
		return nil, ErrItemNotAvailable
	}

	now := time.Now()
	// Check time range
	if !item.StartTime.IsZero() && now.Before(item.StartTime) {
		return nil, ErrItemNotAvailable
	}
	if !item.EndTime.IsZero() && now.After(item.EndTime) {
		return nil, ErrItemNotAvailable
	}

	// Check stock
	if item.Stock != -1 && item.Sold+quantity > item.Stock {
		return nil, ErrOutOfStock
	}

	// Calculate price
	price := int(float64(item.Price) * (1 - item.Discount)) * quantity

	// Check daily limit
	if sm.dailyLimit[playerID] == nil {
		sm.dailyLimit[playerID] = make(map[string]int)
	}
	// Reset daily limit at midnight (simplified)
	today := now.Format("2006-01-02")
	// In real implementation, you'd track date properly

	record := &PurchaseRecord{
		ID:        generateID(),
		PlayerID:  playerID,
		ItemID:    item.ItemID,
		ShopItemID: itemID,
		Quantity:  quantity,
		Price:     price,
		PriceType: item.PriceType,
		Time:      now,
	}

	sm.purchases[playerID] = append(sm.purchases[playerID], record)
	item.Sold += quantity

	return record, nil
}

// GetPurchaseHistory gets player's purchase history
func (sm *ShopManager) GetPurchaseHistory(playerID string) []*PurchaseRecord {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	records := sm.purchases[playerID]
	if records == nil {
		return []*PurchaseRecord{}
	}

	// Return last 50 purchases
	if len(records) > 50 {
		return records[len(records)-50:]
	}
	return records
}

// UpdateItem updates a shop item
func (sm *ShopManager) UpdateItem(itemID string, updates map[string]interface{}) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	item, ok := sm.items[itemID]
	if !ok {
		return ErrItemNotFound
	}

	// Apply updates
	if v, ok := updates["price"]; ok {
		item.Price = int(v.(float64))
	}
	if v, ok := updates["discount"]; ok {
		item.Discount = v.(float64)
	}
	if v, ok := updates["stock"]; ok {
		item.Stock = int(v.(float64))
	}
	if v, ok := updates["is_active"]; ok {
		item.IsActive = v.(bool)
	}
	if v, ok := updates["end_time"]; ok {
		item.EndTime = v.(time.Time)
	}

	return nil
}

// RemoveItem removes an item from shop
func (sm *ShopManager) RemoveItem(itemID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	item, ok := sm.items[itemID]
	if !ok {
		return ErrItemNotFound
	}

	item.IsActive = false
	return nil
}

// RefreshDailyLimits refreshes daily purchase limits
func (sm *ShopManager) RefreshDailyLimits() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.dailyLimit = make(map[string]map[string]int)
}

// GetCategories gets all shop categories
func (sm *ShopManager) GetCategories() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	categories := make([]string, 0, len(sm.categories))
	for cat := range sm.categories {
		categories = append(categories, cat)
	}
	return categories
}

// InitDefaultItems initializes default shop items
func (sm *ShopManager) InitDefaultItems() {
	items := []*ShopItem{
		{
			ID: "shop_weapon_1", ItemID: "weapon_001", ItemType: "weapon",
			Name: "闪电激光", Description: "高品质激光武器",
			Price: 1000, PriceType: "gems", Category: "weapons",
			LevelReq: 1, Discount: 0,
		},
		{
			ID: "shop_weapon_2", ItemID: "weapon_002", ItemType: "weapon",
			Name: "等离子炮", Description: "强力等离子武器",
			Price: 5000, PriceType: "gems", Category: "weapons",
			LevelReq: 10, Discount: 0.1,
		},
		{
			ID: "shop_armor_1", ItemID: "armor_001", ItemType: "armor",
			Name: "能量护盾", Description: "基础能量护盾",
			Price: 500, PriceType: "gems", Category: "armor",
			LevelReq: 1, Discount: 0,
		},
		{
			ID: "shop_consumable_1", ItemID: "hp_potion", ItemType: "consumable",
			Name: "生命药水", Description: "恢复500生命值",
			Price: 100, PriceType: "coins", Category: "consumables",
			LevelReq: 1, Stock: -1, Discount: 0,
		},
		{
			ID: "shop_consumable_2", ItemID: "mp_potion", ItemType: "consumable",
			Name: "能量药水", Description: "恢复300能量值",
			Price: 80, PriceType: "coins", Category: "consumables",
			LevelReq: 1, Stock: -1, Discount: 0,
		},
		{
			ID: "shop_material_1", ItemID: "enhance_stone", ItemType: "material",
			Name: "强化石", Description: "装备强化材料",
			Price: 50, PriceType: "coins", Category: "materials",
			LevelReq: 1, Stock: -1, Discount: 0,
		},
		{
			ID: "shop_gem_1", ItemID: "gem_pack_small", ItemType: "currency",
			Name: "钻石小包", Description: "100钻石",
			Price: 100, PriceType: "rmb", Category: "currency",
			LevelReq: 1, Stock: -1, Discount: 0,
		},
	}

	sm.AddItems(items)
}

// Shop errors
var (
	ErrInvalidItem      = &ShopError{"invalid item"}
	ErrItemNotFound     = &ShopError{"item not found"}
	ErrItemNotAvailable = &ShopError{"item not available"}
	ErrOutOfStock       = &ShopError{"out of stock"}
	ErrInsufficientBalance = &ShopError{"insufficient balance"}
)

type ShopError struct {
	msg string
}

func (e *ShopError) Error() string {
	return e.msg
}

func generateID() string {
	return time.Now().Format("20060102150405") + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
