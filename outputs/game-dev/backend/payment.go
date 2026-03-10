package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================
// 支付系统 (Payment System)
// ============================================

type PaymentStatus int

const (
	PaymentStatusPending PaymentStatus = iota // 待支付
	PaymentStatusSuccess                      // 支付成功
	PaymentStatusFailed                       // 支付失败
	PaymentStatusCancelled                    // 已取消
	PaymentStatusRefunded                     // 已退款
)

// 支付渠道
type PaymentChannel int

const (
	PaymentChannelDouyin PaymentChannel = iota // 抖音支付
	PaymentChannelKuaishou                      // 快手支付
	PaymentChannelAlipay                       // 支付宝
	PaymentChannelWechat                       // 微信支付
	PaymentChannelApple                        // Apple Pay
)

// 支付方式
type PaymentMethod int

const (
	PaymentMethodGem PaymentMethod = iota // 购买钻石
	PaymentMethodVIP                      // 开通VIP
	PaymentMethodGift                     // 购买礼物
	PaymentMethodSkin                     // 购买皮肤
	PaymentMethodPass                     // 购买通行证
)

// 订单
type PaymentOrder struct {
	OrderID     string         `json:"orderId"`     // 订单号
	UserID      string         `json:"userId"`      // 用户ID
	ProductID   string         `json:"productId"`  // 商品ID
	ProductName string         `json:"productName"` // 商品名称
	Amount      int            `json:"amount"`      // 金额(分)
	Currency    string         `json:"currency"`    // 货币 (CNY/USD)
	Channel     PaymentChannel `json:"channel"`     // 支付渠道
	Method      PaymentMethod  `json:"method"`      // 支付方式
	Status      PaymentStatus  `json:"status"`     // 订单状态
	Extra       string         `json:"extra"`      // 附加数据 (JSON)
	CreatedAt   time.Time      `json:"createdAt"`  // 创建时间
	PaidAt       *time.Time     `json:"paidAt"`     // 支付时间
	ExpiredAt    time.Time      `json:"expiresAt"`  // 过期时间
	TradeNo      string         `json:"tradeNo"`     // 第三方交易号
	mu           sync.RWMutex
}

// 商品
type Product struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Type        PaymentMethod  `json:"type"`
	Price       int            `json:"price"`        // 价格(分)
	Currency    string         `json:"currency"`     // 货币
	GemValue    int            `json:"gemValue"`     // 钻石数量
	BonusGem    int            `json:"bonusGem"`     // 赠送钻石
	Tag         string         `json:"tag"`          // 标签 (hot/new/limit)
	Description string         `json:"description"` // 描述
	Icon        string         `json:"icon"`         // 图标URL
	IsActive    bool           `json:"isActive"`     // 是否上架
}

// 商品配置
var ProductConfigs = map[string]Product{
	// 钻石套餐
	"gem_60":    {ID: "gem_60", Name: "60钻石", Type: PaymentMethodGem, Price: 60, Currency: "CNY", GemValue: 60, BonusGem: 0, IsActive: true},
	"gem_300":   {ID: "gem_300", Name: "300钻石", Type: PaymentMethodGem, Price: 300, Currency: "CNY", GemValue: 300, BonusGem: 30, IsActive: true},
	"gem_680":   {ID: "gem_680", Name: "680钻石", Type: PaymentMethodGem, Price: 680, Currency: "CNY", GemValue: 680, BonusGem: 100, IsActive: true},
	"gem_1280":  {ID: "gem_1280", Name: "1280钻石", Type: PaymentMethodGem, Price: 1280, Currency: "CNY", GemValue: 1280, BonusGem: 300, IsActive: true},
	"gem_3280":  {ID: "gem_3280", Name: "3280钻石", Type: PaymentMethodGem, Price: 3280, Currency: "CNY", GemValue: 3280, BonusGem: 1000, IsActive: true},
	"gem_6480":  {ID: "gem_6480", Name: "6480钻石", Type: PaymentMethodGem, Price: 6480, Currency: "CNY", GemValue: 6480, BonusGem: 2500, IsActive: true},

	// VIP
	"vip_weekly":   {ID: "vip_weekly", Name: "周卡", Type: PaymentMethodVIP, Price: 600, Currency: "CNY", GemValue: 0, BonusGem: 100, Tag: "hot", Description: "每日100钻石", IsActive: true},
	"vip_monthly":  {ID: "vip_monthly", Name: "月卡", Type: PaymentMethodVIP, Price: 2500, Currency: "CNY", GemValue: 0, BonusGem: 500, Tag: "hot", Description: "每日500钻石", IsActive: true},
	"vip_yearly":   {ID: "vip_yearly", Name: "年卡", Type: PaymentMethodVIP, Price: 25000, Currency: "CNY", GemValue: 0, BonusGem: 8000, Description: "每日800钻石", IsActive: true},

	// 礼物包
	"gift_rocket":  {ID: "gift_rocket", Name: "火箭礼包", Type: PaymentMethodGift, Price: 500, Currency: "CNY", GemValue: 0, BonusGem: 0, Description: "火箭x10", IsActive: true},
	"gift_car":     {ID: "gift_car", Name: "跑车礼包", Type: PaymentMethodGift, Price: 1000, Currency: "CNY", GemValue: 0, BonusGem: 0, Description: "跑车x10", IsActive: true},

	// 通行证
	"pass_season1": {ID: "pass_season1", Name: "赛季通行证", Type: PaymentMethodPass, Price: 2500, Currency: "CNY", GemValue: 0, BonusGem: 0, Tag: "new", Description: "专属奖励", IsActive: true},
}

// 支付管理器
type PaymentManager struct {
	Orders    map[string]*PaymentOrder // orderID -> Order
	Products  map[string]*Product     // productID -> Product
	Callbacks map[string]func(*PaymentOrder) // 回调函数
	mu        sync.RWMutex
}

func NewPaymentManager() *PaymentManager {
	pm := &PaymentManager{
		Orders:    make(map[string]*PaymentOrder),
		Products:  make(map[string]*Product),
		Callbacks: make(map[string]func(*PaymentOrder)),
	}

	// 加载商品配置
	for id, p := range ProductConfigs {
		pm.Products[id] = &p
	}

	return pm
}

// 创建订单
func (pm *PaymentManager) CreateOrder(userID, productID string, channel PaymentChannel) (*PaymentOrder, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 验证商品
	product, ok := pm.Products[productID]
	if !ok || !product.IsActive {
		return nil, fmt.Errorf("商品不存在或已下架")
	}

	// 生成订单号
	orderID := pm.generateOrderID()

	order := &PaymentOrder{
		OrderID:     orderID,
		UserID:      userID,
		ProductID:   productID,
		ProductName: product.Name,
		Amount:      product.Price,
		Currency:    product.Currency,
		Channel:     channel,
		Method:      product.Type,
		Status:      PaymentStatusPending,
		CreatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(30 * time.Minute), // 30分钟过期
	}

	pm.Orders[orderID] = order
	return order, nil
}

// 生成订单号
func (pm *PaymentManager) generateOrderID() string {
	return fmt.Sprintf("ORD%d%s", time.Now().Unix(), generateID()[:8])
}

// 支付回调 (模拟第三方支付回调)
func (pm *PaymentManager) HandleCallback(orderID, tradeNo string, success bool) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	order, ok := pm.Orders[orderID]
	if !ok {
		return fmt.Errorf("订单不存在")
	}

	if order.Status != PaymentStatusPending {
		return fmt.Errorf("订单状态异常")
	}

	now := time.Now()
	if now.After(order.ExpiredAt) {
		order.Status = PaymentStatusCancelled
		return fmt.Errorf("订单已过期")
	}

	if success {
		order.Status = PaymentStatusSuccess
		order.PaidAt = &now
		order.TradeNo = tradeNo

		// 触发回调
		if cb, ok := pm.Callbacks[order.UserID]; ok {
			cb(order)
		}
	} else {
		order.Status = PaymentStatusFailed
	}

	return nil
}

// 发放商品 (需要在业务逻辑中调用)
func (pm *PaymentManager) DeliverProduct(order *PaymentOrder, playerMgr *PlayerManager) error {
	if order.Status != PaymentStatusSuccess {
		return fmt.Errorf("订单未支付成功")
	}

	product, ok := pm.Products[order.ProductID]
	if !ok {
		return fmt.Errorf("商品不存在")
	}

	// 根据商品类型发放
	switch product.Type {
	case PaymentMethodGem:
		// 发放钻石
		playerMgr.mu.Lock()
		player, ok := playerMgr.Players[order.UserID]
		if !ok {
			playerMgr.mu.Unlock()
			return fmt.Errorf("玩家不存在")
		}
		totalGem := product.GemValue + product.BonusGem
		player.Gem += totalGem
		player.UpdatedAt = time.Now()
		playerMgr.mu.Unlock()
		fmt.Printf("Delivered %d gems to player %s\n", totalGem, order.UserID)

	case PaymentMethodVIP:
		// TODO: VIP逻辑
		fmt.Printf("Activated VIP for player %s\n", order.UserID)

	case PaymentMethodGift, PaymentMethodSkin, PaymentMethodPass:
		// TODO: 其他商品
		fmt.Printf("Delivered product %s to player %s\n", product.ID, order.UserID)
	}

	return nil
}

// 注册支付回调
func (pm *PaymentManager) RegisterCallback(userID string, callback func(*PaymentOrder)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.Callbacks[userID] = callback
}

// 获取订单
func (pm *PaymentManager) GetOrder(orderID string) *PaymentOrder {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.Orders[orderID]
}

// 获取商品列表
func (pm *PaymentManager) GetProducts(method PaymentMethod) []*Product {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make([]*Product, 0)
	for _, p := range pm.Products {
		if p.IsActive && (method == 0 || p.Type == method) {
			result = append(result, p)
		}
	}
	return result
}

// 取消订单
func (pm *PaymentManager) CancelOrder(orderID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	order, ok := pm.Orders[orderID]
	if !ok {
		return fmt.Errorf("订单不存在")
	}

	if order.Status != PaymentStatusPending {
		return fmt.Errorf("只有待支付订单可取消")
	}

	order.Status = PaymentStatusCancelled
	return nil
}

// 退款
func (pm *PaymentManager) RefundOrder(orderID, reason string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	order, ok := pm.Orders[orderID]
	if !ok {
		return fmt.Errorf("订单不存在")
	}

	if order.Status != PaymentStatusSuccess {
		return fmt.Errorf("只有已支付订单可退款")
	}

	order.Status = PaymentStatusRefunded
	order.Extra = reason

	// TODO: 实际退款逻辑
	fmt.Printf("Refunded order %s: %s\n", orderID, reason)
	return nil
}

// ============================================
// 支付网关接口 (模拟)
// ============================================

type PaymentGateway interface {
	CreatePayment(order *PaymentOrder) (paymentURL string, err error)
	QueryPayment(orderID string) (status PaymentStatus, err error)
	CancelPayment(orderID string) error
}

// 抖音支付网关
type DouyinPaymentGateway struct {
	AppID     string
	AppSecret string
	MchID     string
}

func NewDouyinPaymentGateway(appID, appSecret, mchID string) *DouyinPaymentGateway {
	return &DouyinPaymentGateway{
		AppID:     appID,
		AppSecret: appSecret,
		MchID:     mchID,
	}
}

// 创建支付 (模拟)
func (g *DouyinPaymentGateway) CreatePayment(order *PaymentOrder) (string, error) {
	// 实际应该调用抖音支付API
	paymentURL := fmt.Sprintf("douyin://pay?orderId=%s&amount=%d", order.OrderID, order.Amount)
	return paymentURL, nil
}

// 查询支付状态 (模拟)
func (g *DouyinPaymentGateway) QueryPayment(orderID string) (PaymentStatus, error) {
	// 实际应该调用抖音支付API查询
	return PaymentStatusSuccess, nil
}

// 取消支付 (模拟)
func (g *DouyinPaymentGateway) CancelPayment(orderID string) error {
	// 实际应该调用抖音支付API
	return nil
}

// ============================================
// 示例代码
// ============================================

/*
func main() {
	// 创建支付管理器
	pm := NewPaymentManager()
	playerMgr := NewPlayerManager()

	// 创建测试玩家
	playerMgr.CreatePlayer("user001", "测试用户", "")

	// 注册支付回调
	pm.RegisterCallback("user001", func(order *PaymentOrder) {
		fmt.Printf("Payment callback: Order %s status changed to %d\n", order.OrderID, order.Status)
		// 发放商品
		pm.DeliverProduct(order, playerMgr)
	})

	// 获取商品列表
	gems := pm.GetProducts(PaymentMethodGem)
	fmt.Println("=== 钻石商品 ===")
	for _, p := range gems {
		fmt.Printf("%s: %d元 (%d+%d钻石)\n", p.Name, p.Price/100, p.GemValue, p.BonusGem)
	}

	// 创建订单
	order, err := pm.CreateOrder("user001", "gem_680", PaymentChannelDouyin)
	if err != nil {
		fmt.Printf("Create order failed: %v\n", err)
		return
	}
	fmt.Printf("Created order: %s, amount: %d\n", order.OrderID, order.Amount)

	// 模拟支付回调
	err = pm.HandleCallback(order.OrderID, "trade_123456", true)
	if err != nil {
		fmt.Printf("Handle callback failed: %v\n", err)
		return
	}
	fmt.Printf("Order status: %d\n", order.Status)

	// 查看玩家钻石
	player := playerMgr.GetPlayer("user001")
	fmt.Printf("Player gems: %d\n", player.Gem)
}
*/
