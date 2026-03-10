// 商城系统 - danmaku_game/server/mall.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// MallType 商城类型
type MallType int

const (
	MallTypeGift   MallType = 1 // 礼物商城
	MallTypeItem   MallType = 2 // 道具商城
	MallTypeSkin   MallType = 3 // 皮肤商城
	MallTypeRandom MallType = 4 // 随机商城
	MallTypeHonor  MallType = 5 // 荣誉商城
)

// CurrencyType 货币类型
type CurrencyType int

const (
	CurrencyGold  CurrencyType = 1 // 金币
	CurrencyGem   CurrencyType = 2 // 钻石
	CurrencyHonor CurrencyType = 3 // 荣誉点
	CurrencyCredit CurrencyType = 4 // 积分
)

// MallItem 商城商品
type MallItem struct {
	ItemID       string       `json:"item_id"`        // 商品ID
	ItemType     int          `json:"item_type"`      // 道具类型
	ItemName     string       `json:"item_name"`      // 商品名称
	Description  string       `json:"description"`     // 描述
	Price        int          `json:"price"`          // 价格
	Currency     CurrencyType `json:"currency"`       // 货币类型
	Discount     int          `json:"discount"`       // 折扣 (0-100)
	Stock        int          `json:"stock"`          // 库存 (-1=无限)
	LimitType    int          `json:"limit_type"`     // 限制类型: 0=无限制, 1=每日, 2=每周, 3=每月
	LimitCount   int          `json:"limit_count"`    // 限制数量
	RequireLevel int          `json:"require_level"`  // 需求等级
	RequireVIP   int          `json:"require_vip"`    // 需求VIP等级
	StartTime    int64        `json:"start_time"`     // 开始时间
	EndTime      int64        `json:"end_time"`       // 结束时间
	Sort         int          `json:"sort"`           // 排序
	Tag          string       `json:"tag"`            // 标签
	Icon         string       `json:"icon"`           // 图标
	Hot          bool         `json:"hot"`            // 热门
	New          bool         `json:"new"`            // 新品
	Show         bool         `json:"show"`           // 是否展示
}

// MallOrder 商城订单
type MallOrder struct {
	OrderID     string    `json:"order_id"`     // 订单ID
	PlayerID    int64     `json:"player_id"`    // 玩家ID
	ItemID      string    `json:"item_id"`      // 商品ID
	ItemName    string    `json:"item_name"`    // 商品名称
	Count       int       `json:"count"`        // 数量
	Price       int       `json:"price"`        // 单价
	TotalPrice  int       `json:"total_price"`  // 总价
	Currency    CurrencyType `json:"currency"`  // 货币类型
	Status      int       `json:"status"`       // 状态: 0=待支付, 1=已支付, 2=已取消, 3=已退款
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	PaidAt      time.Time `json:"paid_at"`      // 支付时间
}

// PurchaseRecord 购买记录
type PurchaseRecord struct {
	PlayerID   int64     `json:"player_id"`   // 玩家ID
	ItemID     string    `json:"item_id"`     // 商品ID
	Count      int       `json:"count"`       // 购买数量
	TotalPrice int       `json:"total_price"` // 总价
	Period     string    `json:"period"`      // 周期 (daily/weekly/monthly)
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}

// MallSystem 商城系统
type MallSystem struct {
	db           *Database
	cache        *Cache
	inventory    *Inventory
	playerSystem *PlayerSystem
	orderIDGen   *IDGenerator
	malls        map[MallType]map[string]*MallItem // 商城类型 -> 商品ID -> 商品
}

// NewMallSystem 创建商城系统
func NewMallSystem(db *Database, cache *Cache, inventory *Inventory, playerSystem *PlayerSystem) *MallSystem {
	m := &MallSystem{
		db:            db,
		cache:         cache,
		inventory:     inventory,
		playerSystem:  playerSystem,
		orderIDGen:    NewIDGenerator("order"),
		malls:         make(map[MallType]map[string]*MallItem),
	}

	// 初始化商城
	m.initMalls()

	return m
}

// initMalls 初始化商城商品
func (m *MallSystem) initMalls() {
	// 礼物商城
	m.malls[MallTypeGift] = map[string]*MallItem{
		"gift_rose": {
			ItemID: "gift_rose", ItemType: 1, ItemName: "玫瑰",
			Description: "送给心爱的主播", Price: 10, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "鲜花",
		},
		"gift_rocket": {
			ItemID: "gift_rocket", ItemType: 1, ItemName: "火箭",
			Description: "超级火箭", Price: 2000, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "豪礼", Hot: true,
		},
		"gift_car": {
			ItemID: "gift_car", ItemType: 1, ItemName: "跑车",
			Description: "豪华跑车", Price: 5000, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "豪礼", Hot: true,
		},
		"gift_plane": {
			ItemID: "gift_plane", ItemType: 1, ItemName: "飞机",
			Description: "私人飞机", Price: 10000, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "豪礼",
		},
		"gift_heart": {
			ItemID: "gift_heart", ItemType: 1, ItemName: "爱心",
			Description: "小心心", Price: 1, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "小心意",
		},
		"gift_barrage": {
			ItemID: "gift_barrage", ItemType: 1, ItemName: "弹幕雨",
			Description: "满屏弹幕", Price: 100, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "特效",
		},
	}

	// 道具商城
	m.malls[MallTypeItem] = map[string]*MallItem{
		"item_exp_card": {
			ItemID: "item_exp_card", ItemType: 2, ItemName: "经验卡(双倍)",
			Description: "双倍经验持续1小时", Price: 100, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "增益",
		},
		"item_coin_card": {
			ItemID: "item_coin_card", ItemType: 2, ItemName: "金币卡(双倍)",
			Description: "双倍金币持续1小时", Price: 100, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "增益",
		},
		"item_skill_point": {
			ItemID: "item_skill_point", ItemType: 2, ItemName: "技能点",
			Description: "获得50点技能点", Price: 500, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "资源",
		},
		"item_random_box": {
			ItemID: "item_random_box", ItemType: 2, ItemName: "随机宝箱",
			Description: "随机获得道具", Price: 200, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "惊喜", New: true,
		},
		"item_rename_card": {
			ItemID: "item_rename_card", ItemType: 2, ItemName: "改名卡",
			Description: "修改角色名称", Price: 1000, Currency: CurrencyGem,
			Discount: 0, Stock: 1, Show: true, Tag: "功能",
		},
		"item_color_card": {
			ItemID: "item_color_card", ItemType: 2, ItemName: "改名卡(彩色)",
			Description: "修改彩色名称", Price: 5000, Currency: CurrencyGem,
			Discount: 0, Stock: 1, Show: true, Tag: "功能",
		},
	}

	// 皮肤商城
	m.malls[MallTypeSkin] = map[string]*MallItem{
		"skin_red": {
			ItemID: "skin_red", ItemType: 3, ItemName: "红色战衣",
			Description: "炫酷红色皮肤", Price: 2888, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "皮肤",
		},
		"skin_gold": {
			ItemID: "skin_gold", ItemType: 3, ItemName: "金色战衣",
			Description: "炫酷金色皮肤", Price: 5888, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "皮肤", Hot: true,
		},
		"skin_dragon": {
			ItemID: "skin_dragon", ItemType: 3, ItemName: "龙纹战甲",
			Description: "龙纹特效皮肤", Price: 9888, Currency: CurrencyGem,
			Discount: 0, Stock: -1, Show: true, Tag: "皮肤", New: true,
		},
	}

	// 荣誉商城
	m.malls[MallTypeHonor] = map[string]*MallItem{
		"honor_title_1": {
			ItemID: "honor_title_1", ItemType: 4, ItemName: "初级荣誉称号",
			Description: "荣耀的象征", Price: 1000, Currency: CurrencyHonor,
			Discount: 0, Stock: -1, Show: true, Tag: "称号",
		},
		"honor_title_2": {
			ItemID: "honor_title_2", ItemType: 4, ItemName: "中级荣誉称号",
			Description: "尊贵的象征", Price: 5000, Currency: CurrencyHonor,
			Discount: 0, Stock: -1, Show: true, Tag: "称号",
		},
		"honor_title_3": {
			ItemID: "honor_title_3", ItemType: 4, ItemName: "高级荣誉称号",
			Description: "荣耀的象征", Price: 10000, Currency: CurrencyHonor,
			Discount: 0, Stock: -1, Show: true, Tag: "称号",
		},
	}
}

// GetMallItems 获取商城商品列表
func (m *MallSystem) GetMallItems(mallType MallType) ([]*MallItem, error) {
	items, ok := m.malls[mallType]
	if !ok {
		return nil, errors.New("商城不存在")
	}

	var result []*MallItem
	now := time.Now().Unix()

	for _, item := range items {
		if !item.Show {
			continue
		}

		// 检查时间
		if item.StartTime > 0 && now < item.StartTime {
			continue
		}
		if item.EndTime > 0 && now > item.EndTime {
			continue
		}

		result = append(result, item)
	}

	// 排序
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Sort > result[j].Sort {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// GetItem 获取单个商品
func (m *MallSystem) GetItem(mallType MallType, itemID string) (*MallItem, error) {
	items, ok := m.malls[mallType]
	if !ok {
		return nil, errors.New("商城不存在")
	}

	item, ok := items[itemID]
	if !ok {
		return nil, errors.New("商品不存在")
	}

	return item, nil
}

// GetItemPrice 获取商品价格(折扣后)
func (m *MallSystem) GetItemPrice(item *MallItem) int {
	if item.Discount > 0 {
		return item.Price * (100 - item.Discount) / 100
	}
	return item.Price
}

// CheckPurchaseLimit 检查购买限制
func (m *MallSystem) CheckPurchaseLimit(playerID int64, item *MallItem) (bool, error) {
	if item.LimitType == 0 {
		return true, nil
	}

	period := "daily"
	switch item.LimitType {
	case 2:
		period = "weekly"
	case 3:
		period = "monthly"
	}

	key := fmt.Sprintf("purchase:%d:%s:%s", playerID, period, item.ItemID)
	
	// 从数据库查询购买记录
	query := `SELECT count FROM purchase_records WHERE player_id = ? AND item_id = ? AND period = ?`
	var count int
	err := m.db.QueryRow(query, playerID, item.ItemID, period).Scan(&count)
	if err != nil {
		return true, nil // 没有记录表示可以购买
	}

	return count < item.LimitCount, nil
}

// Purchase 购买商品
func (m *MallSystem) Purchase(playerID int64, mallType MallType, itemID string, count int) (*MallOrder, error) {
	if count <= 0 {
		return nil, errors.New("购买数量必须大于0")
	}

	// 获取商品
	item, err := m.GetItem(mallType, itemID)
	if err != nil {
		return nil, err
	}

	// 检查等级限制
	player, err := m.playerSystem.GetPlayer(playerID)
	if err != nil {
		return nil, err
	}
	if player.Level < item.RequireLevel {
		return nil, fmt.Errorf("需要等级%d", item.RequireLevel)
	}

	// 检查VIP限制
	if item.RequireVIP > 0 && player.VIPLevel < item.RequireVIP {
		return nil, fmt.Errorf("需要VIP%d", item.RequireVIP)
	}

	// 检查库存
	if item.Stock >= 0 && item.Stock < count {
		return nil, errors.New("库存不足")
	}

	// 检查购买限制
	canBuy, err := m.CheckPurchaseLimit(playerID, item)
	if err != nil {
		return nil, err
	}
	if !canBuy {
		return nil, errors.New("已达购买上限")
	}

	// 计算价格
	price := m.GetItemPrice(item)
	totalPrice := price * count

	// 扣除货币
	switch item.Currency {
	case CurrencyGold:
		if !m.playerSystem.DeductCoins(playerID, totalPrice) {
			return nil, errors.New("金币不足")
		}
	case CurrencyGem:
		if !m.playerSystem.DeductGems(playerID, totalPrice) {
			return nil, errors.New("钻石不足")
		}
	case CurrencyHonor:
		if !m.playerSystem.DeductHonor(playerID, totalPrice) {
			return nil, errors.New("荣誉点不足")
		}
	case CurrencyCredit:
		if !m.playerSystem.DeductCredit(playerID, totalPrice) {
			return nil, errors.New("积分不足")
		}
	}

	// 发放道具
	m.inventory.AddItem(itemID, count)

	// 创建订单
	order := &MallOrder{
		OrderID:    m.orderIDGen.Generate(),
		PlayerID:   playerID,
		ItemID:     itemID,
		ItemName:   item.ItemName,
		Count:      count,
		Price:      price,
		TotalPrice: totalPrice,
		Currency:   item.Currency,
		Status:     1, // 已支付
		CreatedAt:  time.Now(),
		PaidAt:     time.Now(),
	}

	// 保存订单
	m.saveOrder(order)

	// 更新购买记录
	m.updatePurchaseRecord(playerID, itemID, count, totalPrice)

	// 扣减库存
	if item.Stock > 0 {
		m.updateStock(mallType, itemID, item.Stock-count)
	}

	return order, nil
}

// GetPurchaseHistory 获取购买历史
func (m *MallSystem) GetPurchaseHistory(playerID int64, limit int) ([]*MallOrder, error) {
	query := `SELECT order_id, player_id, item_id, item_name, count, price, total_price, 
			  currency, status, created_at, paid_at 
			  FROM mall_orders WHERE player_id = ? ORDER BY created_at DESC LIMIT ?`

	rows, err := m.db.Query(query, playerID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*MallOrder
	for rows.Next() {
		order := &MallOrder{}
		err := rows.Scan(&order.OrderID, &order.PlayerID, &order.ItemID, &order.ItemName,
			&order.Count, &order.Price, &order.TotalPrice, &order.Currency, &order.Status,
			&order.CreatedAt, &order.PaidAt)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// GetTodaySpend 获取今日消费
func (m *MallSystem) GetTodaySpend(playerID int64, currency CurrencyType) (int, error) {
	startOfDay := time.Now().Format("2006-01-02") + " 00:00:00"
	
	query := `SELECT COALESCE(SUM(total_price), 0) FROM mall_orders 
			  WHERE player_id = ? AND currency = ? AND status = 1 AND created_at >= ?`
	
	var total int
	err := m.db.QueryRow(query, playerID, currency, startOfDay).Scan(&total)
	return total, err
}

// saveOrder 保存订单
func (m *MallSystem) saveOrder(order *MallOrder) error {
	query := `INSERT INTO mall_orders (order_id, player_id, item_id, item_name, count, price, 
			  total_price, currency, status, created_at, paid_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := m.db.Exec(query, order.OrderID, order.PlayerID, order.ItemID, order.ItemName,
		order.Count, order.Price, order.TotalPrice, order.Currency, order.Status,
		order.CreatedAt, order.PaidAt)
	return err
}

// updatePurchaseRecord 更新购买记录
func (m *MallSystem) updatePurchaseRecord(playerID int64, itemID string, count int, totalPrice int) {
	period := "daily"
	
	query := `INSERT INTO purchase_records (player_id, item_id, count, total_price, period, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE 
			  count = count + ?, total_price = total_price + ?, updated_at = ?`
	
	_, err := m.db.Exec(query, playerID, itemID, count, totalPrice, period, count, totalPrice, time.Now())
	if err != nil {
		fmt.Printf("更新购买记录失败: %v\n", err)
	}
}

// updateStock 更新库存
func (m *MallSystem) updateStock(mallType MallType, itemID string, newStock int) {
	if items, ok := m.malls[mallType]; ok {
		if item, ok := items[itemID]; ok {
			item.Stock = newStock
		}
	}
}

// AddMallItem 添加商品(运营后台)
func (m *MallSystem) AddMallItem(mallType MallType, item *MallItem) error {
	if _, ok := m.malls[mallType]; !ok {
		m.malls[mallType] = make(map[string]*MallItem)
	}

	m.malls[mallType][item.ItemID] = item
	return nil
}

// RemoveMallItem 删除商品
func (m *MallSystem) RemoveMallItem(mallType MallType, itemID string) error {
	if items, ok := m.malls[mallType]; ok {
		delete(items, itemID)
	}
	return nil
}

// InitMallTable 初始化商城表
func (m *MallSystem) InitMallTable() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS mall_orders (
			order_id VARCHAR(64) PRIMARY KEY,
			player_id BIGINT NOT NULL,
			item_id VARCHAR(64) NOT NULL,
			item_name VARCHAR(128) NOT NULL,
			count INT NOT NULL DEFAULT 1,
			price INT NOT NULL,
			total_price INT NOT NULL,
			currency TINYINT NOT NULL DEFAULT 1,
			status TINYINT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL,
			paid_at DATETIME,
			INDEX idx_player (player_id),
			INDEX idx_created (created_at)
		)`,
		`CREATE TABLE IF NOT EXISTS purchase_records (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			player_id BIGINT NOT NULL,
			item_id VARCHAR(64) NOT NULL,
			count INT NOT NULL DEFAULT 0,
			total_price INT NOT NULL DEFAULT 0,
			period VARCHAR(32) NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE KEY uk_player_item_period (player_id, item_id, period),
			INDEX idx_player (player_id)
		)`,
	}

	for _, q := range queries {
		if _, err := m.db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

// GetHotItems 获取热门商品
func (m *MallSystem) GetHotItems(limit int) ([]*MallItem, error) {
	var result []*MallItem
	
	for _, items := range m.malls {
		for _, item := range items {
			if item.Hot && item.Show {
				result = append(result, item)
			}
		}
	}

	// 限制数量
	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// GetNewItems 获取新品
func (m *MallSystem) GetNewItems(limit int) ([]*MallItem, error) {
	var result []*MallItem
	
	for _, items := range m.malls {
		for _, item := range items {
			if item.New && item.Show {
				result = append(result, item)
			}
		}
	}

	// 限制数量
	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// GetItemsByTag 按标签获取商品
func (m *MallSystem) GetItemsByTag(mallType MallType, tag string) ([]*MallItem, error) {
	items, err := m.GetMallItems(mallType)
	if err != nil {
		return nil, err
	}

	var result []*MallItem
	for _, item := range items {
		if item.Tag == tag {
			result = append(result, item)
		}
	}

	return result, nil
}

// JSON序列化辅助
func (m *MallSystem) MarshalJSON() (string, error) {
	data := make(map[string]map[string]*MallItem)
	for k, v := range m.malls {
		data[fmt.Sprintf("%d", k)] = v
	}
	return json.Marshal(data)
}
