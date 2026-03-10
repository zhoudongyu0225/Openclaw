package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"game-dev/server/internal/model"
)

// PlayerHandler 玩家数据处理器
type PlayerHandler struct {
	playerCache map[string]*model.Player
}

// NewPlayerHandler 创建玩家处理器
func NewPlayerHandler() *PlayerHandler {
	return &PlayerHandler{
		playerCache: make(map[string]*model.Player),
	}
}

// GetPlayer 获取玩家信息
func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	player, exists := h.playerCache[playerID]
	if !exists {
		// 返回默认玩家数据
		player = &model.Player{
			ID:        playerID,
			Level:     1,
			Exp:       0,
			Coins:     100,
			Gems:      0,
			CreatedAt: time.Now(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// UpdatePlayer 更新玩家数据
func (h *PlayerHandler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		Level    int    `json:"level"`
		Exp      int    `json:"exp"`
		Coins    int    `json:"coins"`
		Gems     int    `json:"gems"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	player := &model.Player{
		ID:        req.PlayerID,
		Level:     req.Level,
		Exp:       req.Exp,
		Coins:     req.Coins,
		Gems:      req.Gems,
		UpdatedAt: time.Now(),
	}

	h.playerCache[req.PlayerID] = player

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"player": req.PlayerID,
	})
}

// AddExp 玩家增加经验
func (h *PlayerHandler) AddExp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		Exp      int    `json:"exp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	player, exists := h.playerCache[req.PlayerID]
	if !exists {
		player = &model.Player{
			ID:        req.PlayerID,
			Level:     1,
			Exp:       0,
			Coins:     100,
			Gems:      0,
			CreatedAt: time.Now(),
		}
	}

	player.Exp += req.Exp
	player.UpdatedAt = time.Now()

	// 检查升级
	for player.Exp >= player.GetExpForNextLevel() {
		player.Exp -= player.GetExpForNextLevel()
		player.Level++
	}

	h.playerCache[req.PlayerID] = player

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// AddCoins 玩家增加金币
func (h *PlayerHandler) AddCoins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		Coins    int    `json:"coins"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	player, exists := h.playerCache[req.PlayerID]
	if !exists {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	player.Coins += req.Coins
	player.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// DeductCoins 扣除玩家金币
func (h *PlayerHandler) DeductCoins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		Coins    int    `json:"coins"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	player, exists := h.playerCache[req.PlayerID]
	if !exists {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	if player.Coins < req.Coins {
		http.Error(w, "Insufficient coins", http.StatusBadRequest)
		return
	}

	player.Coins -= req.Coins
	player.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// AddGems 玩家增加钻石
func (h *PlayerHandler) AddGems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		Gems     int    `json:"gems"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	player, exists := h.playerCache[req.PlayerID]
	if !exists {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	player.Gems += req.Gems
	player.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// GetPlayerStats 获取玩家统计
func (h *PlayerHandler) GetPlayerStats(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	stats := &model.PlayerStats{
		PlayerID:       playerID,
		TotalGames:     0,
		WinGames:       0,
		TotalKills:     0,
		TotalDeaths:    0,
		HighestScore:   0,
		HighestLevel:   1,
		TotalPlayTime:  0,
		LastLoginTime:  time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ResetPlayer 重置玩家数据
func (h *PlayerHandler) ResetPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	delete(h.playerCache, req.PlayerID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Player data reset",
	})
}
