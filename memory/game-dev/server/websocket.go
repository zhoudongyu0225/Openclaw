package game

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// WebSocketMessage WebSocket消息
type WebSocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// WebSocketHandler WebSocket处理器接口
type WebSocketHandler interface {
	Handle(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error
}

// WebSocketConnection WebSocket连接
type WebSocketConnection struct {
	ID        string
	WriteChan chan []byte
	ReadChan  chan *WebSocketMessage
	Context   context.Context
	Cancel    context.CancelFunc
	
	mu       sync.RWMutex
	closed   bool
	playerID string
	roomID   string
}

// NewWebSocketConnection 创建WebSocket连接
func NewWebSocketConnection(id string, writeBufSize, readBufSize int) *WebSocketConnection {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WebSocketConnection{
		ID:        id,
		WriteChan: make(chan []byte, writeBufSize),
		ReadChan:  make(chan *WebSocketMessage, readBufSize),
		Context:   ctx,
		Cancel:    cancel,
	}
}

// Send 发送消息
func (c *WebSocketConnection) Send(msgType string, payload interface{}) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return errors.New("connection closed")
	}
	c.mu.RUnlock()
	
	data, err := json.Marshal(WebSocketMessage{
		Type:    msgType,
		Payload: payload.(json.RawMessage),
	})
	if err != nil {
		return err
	}
	
	select {
	case c.WriteChan <- data:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("send timeout")
	case <-c.Context.Done():
		return errors.New("connection closed")
	}
}

// SendRaw 发送原始数据
func (c *WebSocketConnection) SendRaw(data []byte) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return errors.New("connection closed")
	}
	c.mu.RUnlock()
	
	select {
	case c.WriteChan <- data:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("send timeout")
	case <-c.Context.Done():
		return errors.New("connection closed")
	}
}

// Close 关闭连接
func (c *WebSocketConnection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.closed {
		return
	}
	
	c.closed = true
	c.Cancel()
	close(c.WriteChan)
	close(c.ReadChan)
}

// IsClosed 检查是否关闭
func (c *WebSocketConnection) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// SetPlayer 设置玩家ID
func (c *WebSocketConnection) SetPlayer(playerID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.playerID = playerID
}

// GetPlayer 获取玩家ID
func (c *WebSocketConnection) GetPlayer() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.playerID
}

// SetRoom 设置房间ID
func (c *WebSocketConnection) SetRoom(roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.roomID = roomID
}

// GetRoom 获取房间ID
func (c *WebSocketConnection) GetRoom() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.roomID
}

// WebSocketServer WebSocket服务器
type WebSocketServer struct {
	addr            string
	handler         WebSocketHandler
 upgrader        *WebSocketUpgrader
	connectionMgr  *ConnectionManager
	maxConnections  int
	idleTimeout    time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	
	mu             sync.RWMutex
	server         *http.Server
	running        bool
}

// WebSocketUpgrader WebSocket升级器配置
type WebSocketUpgrader struct {
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin    func(r *http.Request) bool
}

// DefaultUpgrader 默认升级器
var DefaultUpgrader = &WebSocketUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:    func(r *http.Request) bool { return true },
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*WebSocketConnection
	mu           sync.RWMutex
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*WebSocketConnection),
	}
}

// Add 添加连接
func (m *ConnectionManager) Add(conn *WebSocketConnection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[conn.ID] = conn
}

// Remove 移除连接
func (m *ConnectionManager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.connections, id)
}

// Get 获取连接
func (m *ConnectionManager) Get(id string) (*WebSocketConnection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, ok := m.connections[id]
	return conn, ok
}

// GetByPlayer 根据玩家ID获取连接
func (m *ConnectionManager) GetByPlayer(playerID string) *WebSocketConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for _, conn := range m.connections {
		if conn.GetPlayer() == playerID {
			return conn
		}
	}
	
	return nil
}

// GetByRoom 根据房间ID获取连接
func (m *ConnectionManager) GetByRoom(roomID string) []*WebSocketConnection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var conns []*WebSocketConnection
	for _, conn := range m.connections {
		if conn.GetRoom() == roomID {
			conns = append(conns, conn)
		}
	}
	
	return conns
}

// Count 获取连接数
func (m *ConnectionManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

// Broadcast 广播消息
func (m *ConnectionManager) Broadcast(roomID string, msgType string, payload interface{}) error {
	conns := m.GetByRoom(roomID)
	
	var lastErr error
	for _, conn := range conns {
		if err := conn.Send(msgType, payload); err != nil {
			lastErr = err
		}
	}
	
	return lastErr
}

// NewWebSocketServer 创建WebSocket服务器
func NewWebSocketServer(addr string, handler WebSocketHandler, opts ...WebSocketServerOption) *WebSocketServer {
	srv := &WebSocketServer{
		addr:            addr,
		handler:         handler,
		upgrader:        DefaultUpgrader,
		connectionMgr:   NewConnectionManager(),
		maxConnections:  10000,
		idleTimeout:     5 * time.Minute,
		readTimeout:     30 * time.Second,
		writeTimeout:    30 * time.Second,
	}
	
	for _, opt := range opts {
		opt(srv)
	}
	
	return srv
}

// WebSocketServerOption WebSocket服务器配置选项
type WebSocketServerOption func(*WebSocketServer)

// WithMaxConnections 设置最大连接数
func WithMaxConnections(max int) WebSocketServerOption {
	return func(s *WebSocketServer) {
		s.maxConnections = max
	}
}

// WithIdleTimeout 设置空闲超时
func WithIdleTimeout(timeout time.Duration) WebSocketServerOption {
	return func(s *WebSocketServer) {
		s.idleTimeout = timeout
	}
}

// WithReadTimeout 设置读取超时
func WithReadTimeout(timeout time.Duration) WebSocketServerOption {
	return func(s *WebSocketServer) {
		s.readTimeout = timeout
	}
}

// WithWriteTimeout 设置写入超时
func WithWriteTimeout(timeout time.Duration) WebSocketServerOption {
	return func(s *WebSocketServer) {
		s.writeTimeout = timeout
	}
}

// Start 启动服务器
func (s *WebSocketServer) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return errors.New("server already running")
	}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	
	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		IdleTimeout:  s.idleTimeout,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}
	
	s.running = true
	s.mu.Unlock()
	
	return s.server.ListenAndServe()
}

// StartTLS 启动TLS服务器
func (s *WebSocketServer) StartTLS(certFile, keyFile string) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return errors.New("server already running")
	}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWebSocket)
	
	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		IdleTimeout:  s.idleTimeout,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}
	
	s.running = true
	s.mu.Unlock()
	
	return s.server.ListenAndServeTLS(certFile, keyFile)
}

// Stop 停止服务器
func (s *WebSocketServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if !s.running {
		return nil
	}
	
	s.running = false
	
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	
	return nil
}

// handleWebSocket 处理WebSocket连接
func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 检查连接数
	if s.connectionMgr.Count() >= s.maxConnections {
		http.Error(w, "too many connections", http.StatusServiceUnavailable)
		return
	}
	
	// 升级连接
	// 这里简化处理，实际需要使用gorilla/websocket或nhooyr.io/websocket
	// conn, err := s.upgrader.Upgrade(w, r, nil)
	
	// 生成连接ID
	connID := generateConnID()
	
	// 创建连接
	conn := NewWebSocketConnection(connID, 256, 256)
	s.connectionMgr.Add(conn)
	
	// 启动读 goroutine
	go s.readLoop(conn, r)
	
	// 启动写 goroutine
	go s.writeLoop(conn)
	
	// 等待连接关闭
	<-conn.Context.Done()
	
	// 清理
	s.connectionMgr.Remove(connID)
}

// readLoop 读取循环
func (s *WebSocketServer) readLoop(conn *WebSocketConnection, r *http.Request) {
	defer conn.Close()
	
	// 这里简化处理，实际需要循环读取WebSocket帧
	// for {
	//     msgType, reader, err := conn.NextReader()
	//     if err != nil {
	//         break
	//     }
	//     
	//     // 解码消息
	//     var msg WebSocketMessage
	//     if err := json.NewDecoder(reader).Decode(&msg); err != nil {
	//         continue
	//     }
	//     
	//     // 处理消息
	//     if err := s.handler.Handle(conn.Context, conn, &msg); err != nil {
	//         // 处理错误
	//     }
	// }
}

// writeLoop 写入循环
func (s *WebSocketServer) writeLoop(conn *WebSocketConnection) {
	defer conn.Close()
	
	for {
		select {
		case data, ok := <-conn.WriteChan:
			if !ok {
				return
			}
			// 写入WebSocket帧
			// conn.WriteMessage(websocket.TextMessage, data)
		case <-conn.Context.Done():
			return
		}
	}
}

// GetConnectionManager 获取连接管理器
func (s *WebSocketServer) GetConnectionManager() *ConnectionManager {
	return s.connectionMgr
}

// GetStats 获取统计信息
func (s *WebSocketServer) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"connections":    s.connectionMgr.Count(),
		"max_connections": s.maxConnections,
		"running":        s.running,
	}
}

// generateConnID 生成连接ID
func generateConnID() string {
	return fmt.Sprintf("conn-%d-%d", time.Now().UnixNano(), time.Now().Unix()%10000)
}

// GameWebSocketHandler 游戏WebSocket处理器
type GameWebSocketHandler struct {
	roomManager    *RoomManager
	matchmaker    *Matchmaker
	battleManager *BattleManager
	danmakuMgr    *DanmakuManager
}

// NewGameWebSocketHandler 创建游戏WebSocket处理器
func NewGameWebSocketHandler(
	roomManager *RoomManager,
	matchmaker *Matchmaker,
	battleManager *BattleManager,
	danmakuMgr *DanmakuManager,
) *GameWebSocketHandler {
	return &GameWebSocketHandler{
		roomManager:    roomManager,
		matchmaker:    matchmaker,
		battleManager: battleManager,
		danmakuMgr:    danmakuMgr,
	}
}

// Handle 处理消息
func (h *GameWebSocketHandler) Handle(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	switch msg.Type {
	case "login":
		return h.handleLogin(ctx, conn, msg)
	case "join_room":
		return h.handleJoinRoom(ctx, conn, msg)
	case "leave_room":
		return h.handleLeaveRoom(ctx, conn, msg)
	case "ready":
		return h.handleReady(ctx, conn, msg)
	case "start_game":
		return h.handleStartGame(ctx, conn, msg)
	case "input":
		return h.handleInput(ctx, conn, msg)
	case "send_danmaku":
		return h.handleDanmaku(ctx, conn, msg)
	case "match_make":
		return h.handleMatchMake(ctx, conn, msg)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

// handleLogin 处理登录
func (h *GameWebSocketHandler) handleLogin(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	var payload struct {
		PlayerID string `json:"player_id"`
		Token    string `json:"token"`
	}
	
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	
	conn.SetPlayer(payload.PlayerID)
	
	return conn.Send("login_success", map[string]interface{}{
		"player_id": payload.PlayerID,
	})
}

// handleJoinRoom 处理加入房间
func (h *GameWebSocketHandler) handleJoinRoom(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	var payload struct {
		RoomID string `json:"room_id"`
	}
	
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	
	playerID := conn.GetPlayer()
	if playerID == "" {
		return errors.New("not logged in")
	}
	
	room, err := h.roomManager.JoinRoom(payload.RoomID, playerID)
	if err != nil {
		return err
	}
	
	conn.SetRoom(room.ID)
	
	// 通知房间内其他玩家
	for _, otherConn := range h.getRoomConnections(room.ID) {
		if otherConn.ID != conn.ID {
			otherConn.Send("player_joined", map[string]interface{}{
				"player_id": playerID,
			})
		}
	}
	
	return conn.Send("join_room_success", room)
}

// handleLeaveRoom 处理离开房间
func (h *GameWebSocketHandler) handleLeaveRoom(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	roomID := conn.GetRoom()
	if roomID == "" {
		return errors.New("not in room")
	}
	
	playerID := conn.GetPlayer()
	h.roomManager.LeaveRoom(roomID, playerID)
	conn.SetRoom("")
	
	// 通知房间内其他玩家
	for _, otherConn := range h.getRoomConnections(roomID) {
		otherConn.Send("player_left", map[string]interface{}{
			"player_id": playerID,
		})
	}
	
	return conn.Send("leave_room_success", nil)
}

// handleReady 处理准备
func (h *GameWebSocketHandler) handleReady(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	roomID := conn.GetRoom()
	playerID := conn.GetPlayer()
	
	if roomID == "" || playerID == "" {
		return errors.New("not in room")
	}
	
	// 广播准备状态
	return h.connectionMgr.Broadcast(roomID, "player_ready", map[string]interface{}{
		"player_id": playerID,
	})
}

// handleStartGame 处理开始游戏
func (h *GameWebSocketHandler) handleStartGame(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	roomID := conn.GetRoom()
	playerID := conn.GetPlayer()
	
	if roomID == "" || playerID == "" {
		return errors.New("not in room")
	}
	
	// 检查是否是房主
	room, err := h.roomManager.GetRoom(roomID)
	if err != nil {
		return err
	}
	
	if room.HostPlayerID != playerID {
		return errors.New("only host can start game")
	}
	
	// 开始游戏
	battleID, err := h.battleManager.StartBattle(room)
	if err != nil {
		return err
	}
	
	// 广播游戏开始
	return h.connectionMgr.Broadcast(roomID, "game_started", map[string]interface{}{
		"battle_id": battleID,
	})
}

// handleInput 处理玩家输入
func (h *GameWebSocketHandler) handleInput(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	roomID := conn.GetRoom()
	playerID := conn.GetPlayer()
	
	if roomID == "" || playerID == "" {
		return errors.New("not in room")
	}
	
	var payload struct {
		InputType string          `json:"input_type"`
		Data      json.RawMessage `json:"data"`
	}
	
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	
	// 转发到战斗管理器
	return h.battleManager.HandleInput(roomID, playerID, payload.InputType, payload.Data)
}

// handleDanmaku 处理弹幕
func (h *GameWebSocketHandler) handleDanmaku(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	var payload struct {
		Content string `json:"content"`
	}
	
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	
	playerID := conn.GetPlayer()
	roomID := conn.GetRoom()
	
	if roomID == "" {
		return errors.New("not in room")
	}
	
	// 发送弹幕
	danmaku := &Danmaku{
		ID:        generateID(),
		PlayerID:  playerID,
		Content:   payload.Content,
		Timestamp: time.Now(),
		RoomID:    roomID,
	}
	
	h.danmakuMgr.AddDanmaku(roomID, danmaku)
	
	// 广播弹幕
	return h.connectionMgr.Broadcast(roomID, "danmaku", danmaku)
}

// handleMatchMake 处理匹配
func (h *GameWebSocketHandler) handleMatchMake(ctx context.Context, conn *WebSocketConnection, msg *WebSocketMessage) error {
	var payload struct {
		GameMode string `json:"game_mode"`
	}
	
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	
	playerID := conn.GetPlayer()
	if playerID == "" {
		return errors.New("not logged in")
	}
	
	// 请求匹配
	result, err := h.matchmaker.MatchMake(playerID, payload.GameMode)
	if err != nil {
		return err
	}
	
	if result.Success {
		conn.SetRoom(result.RoomID)
		return conn.Send("match_success", result)
	}
	
	return conn.Send("match_pending", result)
}

// getRoomConnections 获取房间内的连接
func (h *GameWebSocketHandler) getRoomConnections(roomID string) []*WebSocketConnection {
	// 这里需要访问WebSocketServer的connectionMgr
	// 实际实现中可以通过依赖注入
	return nil
}

// connectionMgr 连接管理器（用于handler）
func (h *GameWebSocketHandler) connectionMgr *ConnectionManager {
	return nil
}
