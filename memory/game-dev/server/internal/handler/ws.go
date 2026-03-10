package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"barrage-game/internal/model"
	"barrage-game/internal/room"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 生产环境应该检查 Origin
		return true
	},
	HandshakeTimeout: 10 * time.Second,
}

// WSHandler WebSocket处理器
type WSHandler struct {
	roomMgr     *room.Manager
	playerMgr   *PlayerManager
	writeWait   time.Duration
	pongWait    time.Duration
	pingPeriod  time.Duration
	maxMsgSize  int64
}

// NewWSHandler 创建处理器
func NewWSHandler(roomMgr *room.Manager) *WSHandler {
	return &WSHandler{
		roomMgr:    roomMgr,
		playerMgr:  NewPlayerManager(),
		writeWait:  10 * time.Second,
		pongWait:   60 * time.Second,
		pingPeriod: 30 * time.Second,
		maxMsgSize: 512 * 1024, // 512KB
	}
}

// HandleWS 处理WebSocket连接
func (h *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	// 获取playerID (从URL参数或Token解析)
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		playerID = generatePlayerID()
	}

	// 升级为WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("升级WebSocket失败: %v", err)
		return
	}

	// 注册连接
	player := h.playerMgr.Add(playerID, conn)
	defer func() {
		h.playerMgr.Remove(playerID)
		if roomID := h.roomMgr.GetPlayerRoom(playerID); roomID != "" {
			h.roomMgr.LeaveRoom(playerID)
		}
		conn.Close()
	}()

	// 启动读/写协程
	done := make(chan struct{})
	go h.readPump(player, conn, done)
	go h.writePump(player, conn, done)

	<-done // 等待连接关闭
}

// readPump 读取客户端消息
func (h *WSHandler) readPump(player *Player, conn *websocket.Conn, done chan struct{}) {
	defer close(done)

	conn.SetReadLimit(h.maxMsgSize)
	conn.SetReadDeadline(time.Now().Add(h.pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(h.pongWait))
		player.LastActive = time.Now()
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		h.handleMessage(player, message)
		player.LastActive = time.Now()
	}
}

// writePump 发送消息到客户端
func (h *WSHandler) writePump(player *Player, conn *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(h.pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg, ok := <-player.Send:
			conn.SetWriteDeadline(time.Now().Add(h.writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// 添加队列中的其他消息
			n := len(player.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-player.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(h.writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理客户端消息
func (h *WSHandler) handleMessage(player *Player, data []byte) {
	var msg model.CSMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		h.sendError(player, 1, "消息格式错误")
		return
	}

	switch model.CSMsgID(msg.MsgID) {
	case model.CS_HEARTBEAT:
		h.handleHeartbeat(player, data)
	case model.CS_CREATE_ROOM:
		h.handleCreateRoom(player, data)
	case model.CS_JOIN_ROOM:
		h.handleJoinRoom(player, data)
	case model.CS_LEAVE_ROOM:
		h.handleLeaveRoom(player)
	case model.CS_ROOM_LIST:
		h.handleRoomList(player, data)
	case model.CS_READY:
		h.handleReady(player, data)
	case model.CS_START_GAME:
		h.handleStartGame(player)
	default:
		h.sendError(player, 2, "未知消息类型")
	}
}

// handleHeartbeat 心跳
func (h *WSHandler) handleHeartbeat(player *Player, data []byte) {
	var req model.CSHeartbeat
	json.Unmarshal(data, &req)

	player.LastActive = time.Now()

	resp := model.SCHeartbeatAck{
		ServerTime: uint64(time.Now().UnixMilli()),
		ClientTime: req.Timestamp,
	}
	h.send(player, model.SC_HEARTBEAT_ACK, resp)
}

// handleCreateRoom 创建房间
func (h *WSHandler) handleCreateRoom(player *Player, data []byte) {
	var req model.CSCreateRoom
	json.Unmarshal(data, &req)

	r, err := h.roomMgr.CreateRoom(&model.CreateRoomReq{
		Name:       req.Name,
		Mode:       req.Mode,
		MapID:      req.MapID,
		MaxPlayers: req.MaxPlayers,
	})

	if err != nil {
		h.sendError(player, 10, err.Error())
		return
	}

	// 自动加入
	h.roomMgr.JoinRoom(r.ID, player.ID, player.ID)

	resp := model.SCCreateRoomAck{
		Success: true,
		RoomID:  r.ID,
		Room:    convertRoomToInfo(r),
	}
	h.send(player, model.SC_CREATE_ROOM_ACK, resp)
}

// handleJoinRoom 加入房间
func (h *WSHandler) handleJoinRoom(player *Player, data []byte) {
	var req model.CSJoinRoom
	json.Unmarshal(data, &req)

	err := h.roomMgr.JoinRoom(req.RoomID, player.ID, req.PlayerName)
	if err != nil {
		h.sendError(player, 11, err.Error())
		return
	}

	r := h.roomMgr.GetRoom(req.RoomID)
	if r == nil {
		h.sendError(player, 11, "房间不存在")
		return
	}

	resp := model.SCJoinRoomAck{
		Success: true,
		Room:    convertRoomToInfo(r),
	}
	h.send(player, model.SC_JOIN_ROOM_ACK, resp)

	// 广播给房间内其他玩家
	h.roomMgr.Broadcast(req.RoomID, &model.WSMessage{
		MsgID: uint32(model.SC_PLAYER_JOIN),
		Data:  nil, // TODO: 添加玩家信息
	})
}

// handleLeaveRoom 离开房间
func (h *WSHandler) handleLeaveRoom(player *Player) {
	roomID := h.roomMgr.GetPlayerRoom(player.ID)
	if roomID != "" {
		h.roomMgr.LeaveRoom(player.ID)
		h.roomMgr.Broadcast(roomID, &model.WSMessage{
			MsgID: uint32(model.SC_PLAYER_LEAVE),
			Data:  nil,
		})
	}

	h.send(player, model.SC_LEAVE_ROOM_ACK, &model.SCLeaveRoomAck{Success: true})
}

// handleRoomList 房间列表
func (h *WSHandler) handleRoomList(player *Player, data []byte) {
	var req model.CSRoomList
	json.Unmarshal(data, &req)

	var rooms []*model.RoomInfo
	if req.Mode != "" {
		rooms = convertRoomsToInfo(h.roomMgr.ListRoomsByMode(req.Mode))
	} else {
		rooms = convertRoomsToInfo(h.roomMgr.ListRooms())
	}

	h.send(player, model.SC_ROOM_LIST_ACK, &model.SCRoomListAck{
		Rooms:   rooms,
		Total:   uint32(len(rooms)),
		Page:    req.Page,
	})
}

// handleReady 准备
func (h *WSHandler) handleReady(player *Player, data []byte) {
	var req model.CSReady
	json.Unmarshal(data, &req)

	roomID := h.roomMgr.GetPlayerRoom(player.ID)
	if roomID == "" {
		h.sendError(player, 12, "不在房间中")
		return
	}

	r := h.roomMgr.GetRoom(roomID)
	if r == nil {
		h.sendError(player, 12, "房间不存在")
		return
	}

	r.SetPlayerReady(player.ID, req.Ready)

	// 广播准备状态
	h.roomMgr.Broadcast(roomID, &model.WSMessage{
		MsgID: uint32(model.SC_ROOM_UPDATE),
		Data:  nil,
	})
}

// handleStartGame 开始游戏
func (h *WSHandler) handleStartGame(player *Player) {
	roomID := h.roomMgr.GetPlayerRoom(player.ID)
	if roomID == "" {
		h.sendError(player, 13, "不在房间中")
		return
	}

	r := h.roomMgr.GetRoom(roomID)
	if r == nil {
		h.sendError(player, 13, "房间不存在")
		return
	}

	if !r.IsOwner(player.ID) {
		h.sendError(player, 13, "只有房主可以开始游戏")
		return
	}

	if !r.AllPlayersReady() {
		h.sendError(player, 13, "并非所有玩家都准备好了")
		return
	}

	if err := r.StartGame(); err != nil {
		h.sendError(player, 13, err.Error())
		return
	}

	// 广播游戏开始
	h.roomMgr.Broadcast(roomID, &model.WSMessage{
		MsgID: uint32(model.SC_GAME_START),
		Data:  nil,
	})
}

// send 发送消息
func (h *WSHandler) send(player *Player, msgID model.SCMsgID, data interface{}) {
	msg := model.WSMessage{
		MsgID: uint32(msgID),
	}

	if data != nil {
		b, _ := json.Marshal(data)
		msg.Data = b
	}

	select {
	case player.Send <- b:
	default:
		log.Printf("发送队列已满, player=%s", player.ID)
	}
}

// sendError 发送错误
func (h *WSHandler) sendError(player *Player, code uint32, message string) {
	h.send(player, model.SC_ERROR, &model.SCError{
		Code:    code,
		Message: message,
	})
}
