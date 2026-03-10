package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ==================== HTTP API Handler ====================

// APIHandler HTTP API 处理器
type APIHandler struct {
	roomMgr      *RoomManager
	matchmaker   *Matchmaker
	leaderboard  *Leaderboard
	statsMgr     *StatsManager
	achMgr       *AchievementManager
	questMgr     *QuestManager
	validator    *InputValidator
	limiter      *RateLimiter
}

// NewAPIHandler 创建API处理器
func NewAPIHandler(roomMgr *RoomManager) *APIHandler {
	return &APIHandler{
		roomMgr:     roomMgr,
		matchmaker:  NewMatchmaker(2, 60*time.Second),
		leaderboard: NewLeaderboard(100),
		statsMgr:    NewStatsManager(),
		achMgr:      NewAchievementManager(),
		questMgr:    NewQuestManager(),
		validator:   NewInputValidator(),
		limiter:     NewRateLimiter(100, time.Minute),
	}
}

// APIResponse API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ==================== 房间管理 API ====================

// HandleRoomCreate 处理房间创建
func (h *APIHandler) HandleRoomCreate(w http.ResponseWriter, r *http.Request) {
	if !h.limiter.Allow(r.RemoteAddr) {
		h.sendError(w, http.StatusTooManyRequests, "请求过于频繁")
		return
	}

	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		HostID   string `json:"host_id"`
		HostName string `json:"host_name"`
		Mode     string `json:"mode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	// 验证输入
	if err := h.validator.ValidateRoomName(req.Name); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	room, err := h.roomMgr.CreateRoom(&CreateRoomReq{
		ID:       req.ID,
		Name:     req.Name,
		HostID:   req.HostID,
		HostName: req.HostName,
		Mode:     req.Mode,
	})

	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, room)
}

// HandleRoomList 处理房间列表
func (h *APIHandler) HandleRoomList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	mode := r.URL.Query().Get("mode")
	var rooms []*Room

	if mode != "" {
		rooms = h.roomMgr.ListRoomsByMode(mode)
	} else {
		rooms = h.roomMgr.ListRooms()
	}

	h.sendSuccess(w, rooms)
}

// HandleRoomGet 处理获取房间详情
func (h *APIHandler) HandleRoomGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	roomID := r.URL.Query().Get("id")
	if roomID == "" {
		h.sendError(w, http.StatusBadRequest, "缺少房间ID")
		return
	}

	room := h.roomMgr.GetRoom(roomID)
	if room == nil {
		h.sendError(w, http.StatusNotFound, "房间不存在")
		return
	}

	h.sendSuccess(w, room)
}

// HandleRoomJoin 处理加入房间
func (h *APIHandler) HandleRoomJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		RoomID   string `json:"room_id"`
		PlayerID string `json:"player_id"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	// 验证输入
	if err := h.validator.ValidateUsername(req.Name); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.roomMgr.JoinRoom(req.RoomID, req.PlayerID, req.Name)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, map[string]string{"status": "joined"})
}

// HandleRoomLeave 处理离开房间
func (h *APIHandler) HandleRoomLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	err := h.roomMgr.LeaveRoom(req.PlayerID)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, map[string]string{"status": "left"})
}

// ==================== 匹配系统 API ====================

// HandleMatchMake 处理匹配请求
func (h *APIHandler) HandleMatchMake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string  `json:"player_id"`
		Rating   float64 `json:"rating"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	h.matchmaker.AddPlayer(req.PlayerID, req.Rating)
	h.sendSuccess(w, map[string]string{"status": "matching"})
}

// HandleMatchCancel 处理取消匹配
func (h *APIHandler) HandleMatchCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	h.matchmaker.RemovePlayer(req.PlayerID)
	h.sendSuccess(w, map[string]string{"status": "cancelled"})
}

// HandleMatchStatus 处理匹配状态查询
func (h *APIHandler) HandleMatchStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		h.sendError(w, http.StatusBadRequest, "缺少玩家ID")
		return
	}

	matches := h.matchmaker.GetMatches()
	var matchInfo interface{}

	for _, match := range matches {
		for _, p := range match.Players {
			if p == playerID {
				matchInfo = match
				break
			}
		}
	}

	if matchInfo == nil {
		h.sendSuccess(w, map[string]string{"status": "waiting"})
	} else {
		h.sendSuccess(w, matchInfo)
	}
}

// ==================== 排行榜 API ====================

// HandleLeaderboardUpdate 处理排行榜更新
func (h *APIHandler) HandleLeaderboardUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		Score    int64  `json:"score"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	h.leaderboard.UpdateScore(req.PlayerID, req.Score)
	h.sendSuccess(w, map[string]string{"status": "updated"})
}

// HandleLeaderboardGet 处理获取排行榜
func (h *APIHandler) HandleLeaderboardGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	leaderboard := h.leaderboard.GetTop(limit)
	h.sendSuccess(w, leaderboard)
}

// HandleLeaderboardRank 处理获取玩家排名
func (h *APIHandler) HandleLeaderboardRank(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		h.sendError(w, http.StatusBadRequest, "缺少玩家ID")
		return
	}

	rank := h.leaderboard.GetRank(playerID)
	h.sendSuccess(w, map[string]int{"rank": rank})
}

// ==================== 成就系统 API ====================

// HandleAchievementList 处理获取玩家成就列表
func (h *APIHandler) HandleAchievementList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		h.sendError(w, http.StatusBadRequest, "缺少玩家ID")
		return
	}

	// 获取所有成就定义
	achievements := h.achMgr.achievements

	// 获取玩家进度
	progress := h.achMgr.GetPlayerProgress(playerID)

	h.sendSuccess(w, map[string]interface{}{
		"achievements": achievements,
		"progress":     progress,
	})
}

// HandleAchievementUpdate 处理成就进度更新
func (h *APIHandler) HandleAchievementUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID     string `json:"player_id"`
		AchievementID string `json:"achievement_id"`
		Value        int   `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	h.achMgr.UpdateProgress(req.PlayerID, req.AchievementID, req.Value)
	h.sendSuccess(w, map[string]string{"status": "updated"})
}

// HandleAchievementClaim 处理成就奖励领取
func (h *APIHandler) HandleAchievementClaim(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID     string `json:"player_id"`
		AchievementID string `json:"achievement_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	reward := h.achMgr.ClaimReward(req.PlayerID, req.AchievementID)
	if reward == nil {
		h.sendError(w, http.StatusBadRequest, "成就未完成或已领取")
		return
	}

	h.sendSuccess(w, reward)
}

// ==================== 任务系统 API ====================

// HandleQuestList 处理获取玩家任务列表
func (h *APIHandler) HandleQuestList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		h.sendError(w, http.StatusBadRequest, "缺少玩家ID")
		return
	}

	// 初始化玩家任务
	h.questMgr.InitPlayer(playerID)

	// 获取所有任务
	quests := h.questMgr.quests

	// 获取玩家任务状态
	playerQuests := h.questMgr.GetPlayerQuests(playerID)

	h.sendSuccess(w, map[string]interface{}{
		"quests":      quests,
		"playerQuests": playerQuests,
	})
}

// HandleQuestAccept 处理接受任务
func (h *APIHandler) HandleQuestAccept(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		QuestID  string `json:"quest_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	err := h.questMgr.AcceptQuest(req.PlayerID, req.QuestID)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.sendSuccess(w, map[string]string{"status": "accepted"})
}

// HandleQuestUpdate 处理任务进度更新
func (h *APIHandler) HandleQuestUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		QuestID  string `json:"quest_id"`
		Value    int   `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	h.questMgr.UpdateProgress(req.PlayerID, req.QuestID, req.Value)
	h.sendSuccess(w, map[string]string{"status": "updated"})
}

// HandleQuestComplete 处理完成任务
func (h *APIHandler) HandleQuestComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		QuestID  string `json:"quest_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	err := h.questMgr.CompleteQuest(req.PlayerID, req.QuestID)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.sendSuccess(w, map[string]string{"status": "completed"})
}

// ==================== 统计系统 API ====================

// HandleStatsGet 处理获取玩家统计
func (h *APIHandler) HandleStatsGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		h.sendError(w, http.StatusBadRequest, "缺少玩家ID")
		return
	}

	stats := h.statsMgr.GetStats(playerID)
	h.sendSuccess(w, stats)
}

// HandleStatsRecordGameStart 处理记录游戏开始
func (h *APIHandler) HandleStatsRecordGameStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	h.statsMgr.RecordGameStart(req.PlayerID)
	h.sendSuccess(w, map[string]string{"status": "recorded"})
}

// HandleStatsRecordGameEnd 处理记录游戏结束
func (h *APIHandler) HandleStatsRecordGameEnd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		Win      bool   `json:"win"`
		Duration int    `json:"duration"` // 秒
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	if req.Win {
		h.statsMgr.RecordGameWin(req.PlayerID, req.Duration)
	} else {
		h.statsMgr.RecordGameLose(req.PlayerID, req.Duration)
	}

	h.sendSuccess(w, map[string]string{"status": "recorded"})
}

// ==================== 弹幕系统 API ====================

// HandleDanmakuSend 处理发送弹幕
func (h *APIHandler) HandleDanmakuSend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		Content  string `json:"content"`
		Type     string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "请求格式错误")
		return
	}

	// 验证弹幕内容
	if err := h.validator.ValidateDanmaku(req.Content); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 敏感词过滤
	content := h.validator.FilterSensitiveWords(req.Content)

	// 创建弹幕管理器并发送
	danmakuMgr := NewDanmakuManager()
	danmakuMgr.Send(req.PlayerID, content, DanmakuTypeText)

	h.sendSuccess(w, map[string]string{"status": "sent"})
}

// ==================== 工具方法 ====================

// sendSuccess 发送成功响应
func (h *APIHandler) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// sendError 发送错误响应
func (h *APIHandler) sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIResponse{
		Code:    code,
		Message: message,
	})
}

// ==================== 路由注册 ====================

// RegisterRoutes 注册路由
func (h *APIHandler) RegisterRoutes(mux *http.ServeMux) {
	// 房间管理
	mux.HandleFunc("/api/room/create", h.HandleRoomCreate)
	mux.HandleFunc("/api/room/list", h.HandleRoomList)
	mux.HandleFunc("/api/room/get", h.HandleRoomGet)
	mux.HandleFunc("/api/room/join", h.HandleRoomJoin)
	mux.HandleFunc("/api/room/leave", h.HandleRoomLeave)

	// 匹配系统
	mux.HandleFunc("/api/match/make", h.HandleMatchMake)
	mux.HandleFunc("/api/match/cancel", h.HandleMatchCancel)
	mux.HandleFunc("/api/match/status", h.HandleMatchStatus)

	// 排行榜
	mux.HandleFunc("/api/leaderboard/update", h.HandleLeaderboardUpdate)
	mux.HandleFunc("/api/leaderboard/get", h.HandleLeaderboardGet)
	mux.HandleFunc("/api/leaderboard/rank", h.HandleLeaderboardRank)

	// 成就系统
	mux.HandleFunc("/api/achievement/list", h.HandleAchievementList)
	mux.HandleFunc("/api/achievement/update", h.HandleAchievementUpdate)
	mux.HandleFunc("/api/achievement/claim", h.HandleAchievementClaim)

	// 任务系统
	mux.HandleFunc("/api/quest/list", h.HandleQuestList)
	mux.HandleFunc("/api/quest/accept", h.HandleQuestAccept)
	mux.HandleFunc("/api/quest/update", h.HandleQuestUpdate)
	mux.HandleFunc("/api/quest/complete", h.HandleQuestComplete)

	// 统计系统
	mux.HandleFunc("/api/stats/get", h.HandleStatsGet)
	mux.HandleFunc("/api/stats/game/start", h.HandleStatsRecordGameStart)
	mux.HandleFunc("/api/stats/game/end", h.HandleStatsRecordGameEnd)

	// 弹幕系统
	mux.HandleFunc("/api/danmaku/send", h.HandleDanmakuSend)

	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		h.sendSuccess(w, map[string]string{"status": "ok"})
	})
}

// StartHTTPServer 启动HTTP服务器
func StartHTTPServer(addr string, roomMgr *RoomManager) error {
	handler := NewAPIHandler(roomMgr)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// CORS 中间件
	enhancedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		mux.ServeHTTP(w, r)
	})

	return http.ListenAndServe(addr, enhancedMux)
}

// ==================== 路由简写 ====================

// 简化路由处理函数
func handleAPI(w http.ResponseWriter, r *http.Request, handler func(*APIHandler, http.ResponseWriter, *http.Request)) {
	// 从上下文或全局获取handler
	// 实际使用时需要通过依赖注入传递
	h := &APIHandler{}
	handler(h, w, r)
}

// HTTP路由处理器的便捷函数
func handleRoomList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Room List API")
}

func handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/update") {
		fmt.Fprintf(w, "Leaderboard Update API")
	} else if strings.HasSuffix(path, "/rank") {
		fmt.Fprintf(w, "Leaderboard Rank API")
	} else {
		fmt.Fprintf(w, "Leaderboard Get API")
	}
}

func handleAchievement(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/list") {
		fmt.Fprintf(w, "Achievement List API")
	} else if strings.HasSuffix(path, "/claim") {
		fmt.Fprintf(w, "Achievement Claim API")
	} else {
		fmt.Fprintf(w, "Achievement Update API")
	}
}

func handleQuest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/list") {
		fmt.Fprintf(w, "Quest List API")
	} else if strings.HasSuffix(path, "/accept") {
		fmt.Fprintf(w, "Quest Accept API")
	} else if strings.HasSuffix(path, "/complete") {
		fmt.Fprintf(w, "Quest Complete API")
	} else {
		fmt.Fprintf(w, "Quest Update API")
	}
}
