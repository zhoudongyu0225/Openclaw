package main

import (
	"fmt"
	"sync"
	"time"
)

// ============================================
// 数据库层 (Database Layer - MongoDB 风格)
// ============================================

// 文档基类
type Document struct {
	ID        string    `json:"_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 用户文档
type UserDocument struct {
	Document
	UserID    string `json:"userId"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Level     int    `json:"level"`
	Score     int    `json:"score"`
	Money     int    `json:"money"`
	Gem       int    `json:"gem"`
	TotalGames int   `json:"totalGames"`
	WinCount  int    `json:"winCount"`
	LoseCount int    `json:"loseCount"`
}

// 房间文档
type RoomDocument struct {
	Document
	RoomID    string   `json:"roomId"`
	Name      string   `json:"name"`
	Mode      string   `json:"mode"`
	MapID     string   `json:"mapId"`
	Status    string   `json:"status"` // waiting, playing, finished
	HostID    string   `json:"hostId"`
	GuestIDs  []string `json:"guestIds"`
	Settings  RoomSettings
}

// 战斗记录文档
type BattleRecordDocument struct {
	Document
	RoomID     string         `json:"roomId"`
	WinnerID   string         `json:"winnerId"`
	Duration   int            `json:"duration"` // 秒
	Wave       int            `json:"wave"`
	Score      int            `json:"score"`
	Players    []BattlePlayer `json:"players"`
}

// 战斗玩家
type BattlePlayer struct {
	UserID     string `json:"userId"`
	Nickname   string `json:"nickname"`
	Score      int    `json:"score"`
	Kills      int    `json:"kills"`
	Damage     int    `json:"damage"`
	TowerKills int    `json:"towerKills"`
}

// 订单文档
type OrderDocument struct {
	Document
	OrderID    string         `json:"orderId"`
	UserID     string         `json:"userId"`
	ProductID  string         `json:"productId"`
	Amount     int            `json:"amount"`
	Currency   string         `json:"currency"`
	Status     PaymentStatus  `json:"status"`
	Channel    PaymentChannel `json:"channel"`
	TradeNo    string         `json:"tradeNo"`
	PaidAt     *time.Time     `json:"paidAt"`
}

// 数据库接口
type Database interface {
	// 用户
	FindUser(userID string) (*UserDocument, error)
	InsertUser(user *UserDocument) error
	UpdateUser(userID string, update func(*UserDocument)) error
	DeleteUser(userID string) error

	// 房间
	FindRoom(roomID string) (*RoomDocument, error)
	InsertRoom(room *RoomDocument) error
	UpdateRoom(roomID string, update func(*RoomDocument)) error
	DeleteRoom(roomID string) error
	FindRoomsByStatus(status string) ([]*RoomDocument, error)

	// 战斗记录
	InsertBattleRecord(record *BattleRecordDocument) error
	FindBattleRecords(userID string, limit int) ([]*BattleRecordDocument, error)

	// 订单
	FindOrder(orderID string) (*OrderDocument, error)
	InsertOrder(order *OrderDocument) error
	UpdateOrder(orderID string, update func(*OrderDocument)) error
}

// 内存数据库实现 (生产环境应替换为 MongoDB)
type MemoryDatabase struct {
	users    map[string]*UserDocument
	rooms    map[string]*RoomDocument
	battles  []*BattleRecordDocument
	orders   map[string]*OrderDocument
	mu       sync.RWMutex
}

func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{
		users:   make(map[string]*UserDocument),
		rooms:   make(map[string]*RoomDocument),
		battles: make([]*BattleRecordDocument, 0),
		orders:  make(map[string]*OrderDocument),
	}
}

// 用户操作
func (db *MemoryDatabase) FindUser(userID string) (*UserDocument, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if user, ok := db.users[userID]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (db *MemoryDatabase) InsertUser(user *UserDocument) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.users[user.UserID] = user
	return nil
}

func (db *MemoryDatabase) UpdateUser(userID string, update func(*UserDocument)) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	user, ok := db.users[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	update(user)
	user.UpdatedAt = time.Now()
	return nil
}

func (db *MemoryDatabase) DeleteUser(userID string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.users, userID)
	return nil
}

// 房间操作
func (db *MemoryDatabase) FindRoom(roomID string) (*RoomDocument, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if room, ok := db.rooms[roomID]; ok {
		return room, nil
	}
	return nil, fmt.Errorf("room not found")
}

func (db *MemoryDatabase) InsertRoom(room *RoomDocument) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.rooms[room.RoomID] = room
	return nil
}

func (db *MemoryDatabase) UpdateRoom(roomID string, update func(*RoomDocument)) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	room, ok := db.rooms[roomID]
	if !ok {
		return fmt.Errorf("room not found")
	}
	update(room)
	room.UpdatedAt = time.Now()
	return nil
}

func (db *MemoryDatabase) DeleteRoom(roomID string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.rooms, roomID)
	return nil
}

func (db *MemoryDatabase) FindRoomsByStatus(status string) ([]*RoomDocument, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	result := make([]*RoomDocument, 0)
	for _, room := range db.rooms {
		if room.Status == status {
			result = append(result, room)
		}
	}
	return result, nil
}

// 战斗记录操作
func (db *MemoryDatabase) InsertBattleRecord(record *BattleRecordDocument) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.battles = append(db.battles, record)
	return nil
}

func (db *MemoryDatabase) FindBattleRecords(userID string, limit int) ([]*BattleRecordDocument, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	result := make([]*BattleRecordDocument, 0)
	for i := len(db.battles) - 1; i >= 0 && len(result) < limit; i-- {
		record := db.battles[i]
		for _, p := range record.Players {
			if p.UserID == userID {
				result = append(result, record)
				break
			}
		}
	}
	return result, nil
}

// 订单操作
func (db *MemoryDatabase) FindOrder(orderID string) (*OrderDocument, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if order, ok := db.orders[orderID]; ok {
		return order, nil
	}
	return nil, fmt.Errorf("order not found")
}

func (db *MemoryDatabase) InsertOrder(order *OrderDocument) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.orders[order.OrderID] = order
	return nil
}

func (db *MemoryDatabase) UpdateOrder(orderID string, update func(*OrderDocument)) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	order, ok := db.orders[orderID]
	if !ok {
		return fmt.Errorf("order not found")
	}
	update(order)
	order.UpdatedAt = time.Now()
	return nil
}

// ============================================
// Redis 缓存层 (简单实现)
// ============================================

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl int) error
	Del(key string) error
	Exists(key string) (bool, error)
	Inc(key string) (int64, error)
	Dec(key string) (int64, error)
	Expire(key string, ttl int) error
}

type MemoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	Value    string
	ExpireAt *time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]cacheItem),
	}
}

func (c *MemoryCache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.data[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}
	if item.ExpireAt != nil && time.Now().After(*item.ExpireAt) {
		return "", fmt.Errorf("key expired")
	}
	return item.Value, nil
}

func (c *MemoryCache) Set(key string, value string, ttl int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var expireAt *time.Time
	if ttl > 0 {
		t := time.Now().Add(time.Duration(ttl) * time.Second)
		expireAt = &t
	}
	c.data[key] = cacheItem{Value: value, ExpireAt: expireAt}
	return nil
}

func (c *MemoryCache) Del(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *MemoryCache) Exists(key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.data[key]
	return ok, nil
}

func (c *MemoryCache) Inc(key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.data[key]
	var val int64
	if ok {
		fmt.Sscanf(item.Value, "%d", &val)
	}
	val++
	c.data[key] = cacheItem{Value: fmt.Sprintf("%d", val)}
	return val, nil
}

func (c *MemoryCache) Dec(key string) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.data[key]
	var val int64
	if ok {
		fmt.Sscanf(item.Value, "%d", &val)
	}
	val--
	c.data[key] = cacheItem{Value: fmt.Sprintf("%d", val)}
	return val, nil
}

func (c *MemoryCache) Expire(key string, ttl int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.data[key]
	if !ok {
		return fmt.Errorf("key not found")
	}
	if ttl > 0 {
		t := time.Now().Add(time.Duration(ttl) * time.Second)
		item.ExpireAt = &t
	}
	c.data[key] = item
	return nil
}

// ============================================
// 游戏服务 (整合所有模块)
// ============================================

type GameService struct {
	DB           Database
	Cache        Cache
	PlayerMgr    *PlayerManager
	RoomMgr      *Manager
	PaymentMgr   *PaymentManager
	Leaderboards *LeaderboardManager
}

func NewGameService() *GameService {
	db := NewMemoryDatabase()
	return &GameService{
		DB:           db,
		Cache:        NewMemoryCache(),
		PlayerMgr:    NewPlayerManager(),
		RoomMgr:      NewManager(30 * time.Second),
		PaymentMgr:   NewPaymentManager(),
		Leaderboards: NewLeaderboardManager(),
	}
}

// 用户登录/注册
func (s *GameService) LoginOrRegister(userID, nickname, avatar string) *PlayerProfile {
	// 先从数据库查找
	user, err := s.DB.FindUser(userID)
	if err == nil {
		// 同步到内存
		profile := s.PlayerMgr.CreatePlayer(user.UserID, user.Nickname, user.Avatar)
		profile.Score = user.Score
		profile.Money = user.Money
		profile.Gem = user.Gem
		profile.Level = user.Level
		profile.WinCount = user.WinCount
		profile.LoseCount = user.LoseCount
		return profile
	}

	// 创建新用户
	profile := s.PlayerMgr.CreatePlayer(userID, nickname, avatar)

	// 存入数据库
	s.DB.InsertUser(&UserDocument{
		Document: Document{
			ID:        userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID:    userID,
		Nickname:  nickname,
		Avatar:    avatar,
		Level:     1,
		Score:     1000,
		Money:     100,
		Gem:       0,
		TotalGames: 0,
		WinCount:  0,
		LoseCount: 0,
	})

	return profile
}

// 保存玩家数据到数据库
func (s *GameService) SavePlayerProfile(userID string) error {
	profile := s.PlayerMgr.GetPlayer(userID)
	if profile == nil {
		return fmt.Errorf("player not found")
	}

	return s.DB.UpdateUser(userID, func(doc *UserDocument) {
		doc.Nickname = profile.Nickname
		doc.Avatar = profile.Avatar
		doc.Level = profile.Level
		doc.Score = profile.Score
		doc.Money = profile.Money
		doc.Gem = profile.Gem
		doc.TotalGames = profile.TotalGames
		doc.WinCount = profile.WinCount
		doc.LoseCount = profile.LoseCount
	})
}

// 创建支付订单
func (s *GameService) CreatePayment(userID, productID string, channel PaymentChannel) (*PaymentOrder, error) {
	order, err := s.PaymentMgr.CreateOrder(userID, productID, channel)
	if err != nil {
		return nil, err
	}

	// 存入数据库
	s.DB.InsertOrder(&OrderDocument{
		Document: Document{
			ID:        order.OrderID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		OrderID:   order.OrderID,
		UserID:    order.UserID,
		ProductID: order.ProductID,
		Amount:    order.Amount,
		Currency:  order.Currency,
		Status:    order.Status,
		Channel:   order.Channel,
	})

	return order, nil
}

// 处理支付回调
func (s *GameService) HandlePaymentCallback(orderID, tradeNo string, success bool) error {
	err := s.PaymentMgr.HandleCallback(orderID, tradeNo, success)
	if err != nil {
		return err
	}

	order := s.PaymentMgr.GetOrder(orderID)
	if order != nil && order.Status == PaymentStatusSuccess {
		// 发放商品
		s.PaymentMgr.DeliverProduct(order, s.PlayerMgr)
		// 保存玩家数据
		s.SavePlayerProfile(order.UserID)
	}

	return nil
}

// 记录战斗结果
func (s *GameService) RecordBattle(roomID, winnerID string, duration int, wave int, score int, players []BattlePlayer) error {
	record := &BattleRecordDocument{
		Document: Document{
			ID:        generateID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		RoomID:   roomID,
		WinnerID: winnerID,
		Duration: duration,
		Wave:     wave,
		Score:    score,
		Players:  players,
	}

	// 存入数据库
	s.DB.InsertBattleRecord(record)

	// 更新玩家数据
	for _, p := range players {
		if p.UserID == winnerID {
			s.PlayerMgr.PlayerWin(p.UserID, score, 100, p.Kills)
		} else {
			s.PlayerMgr.PlayerLose(p.UserID, 50)
		}
	}

	return nil
}

// 获取玩家战斗历史
func (s *GameService) GetBattleHistory(userID string, limit int) ([]*BattleRecordDocument, error) {
	return s.DB.FindBattleRecords(userID, limit)
}

// 获取排行榜
func (s *GameService) GetLeaderboard(lbType LeaderboardType, limit int) []*LeaderboardEntry {
	return s.Leaderboards.GetLeaderboard(lbType, limit)
}

// ============================================
// 示例
// ============================================

/*
func main() {
	// 创建游戏服务
	service := NewGameService()

	// 用户登录
	user := service.LoginOrRegister("user001", "主播A", "avatar.png")
	fmt.Printf("User: %s, Score: %d, Money: %d\n", user.Nickname, user.Score, user.Money)

	// 创建支付订单
	order, err := service.CreatePayment("user001", "gem_680", PaymentChannelDouyin)
	if err != nil {
		fmt.Printf("Create payment failed: %v\n", err)
		return
	}
	fmt.Printf("Order created: %s, Amount: %d\n", order.OrderID, order.Amount)

	// 模拟支付成功
	service.HandlePaymentCallback(order.OrderID, "trade_123456", true)

	// 验证充值结果
	user = service.GetPlayer("user001")
	fmt.Printf("After payment: Gem: %d\n", user.Gem)

	// 记录战斗
	service.RecordBattle("room001", "user001", 180, 5, 1000, []BattlePlayer{
		{UserID: "user001", Nickname: "主播A", Score: 1000, Kills: 10, Damage: 5000},
		{UserID: "user002", Nickname: "观众B", Score: 500, Kills: 3, Damage: 2000},
	})

	// 查看排行榜
	board := service.GetLeaderboard(LeaderboardTypeScore, 10)
	fmt.Println("\n=== Leaderboard ===")
	for _, e := range board {
		fmt.Printf("#%d %s: %d\n", e.Rank, e.Nickname, e.Value)
	}
}
*/
