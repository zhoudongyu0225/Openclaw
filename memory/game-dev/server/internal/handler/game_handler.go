package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"game-dev/server/internal/friend"
	"game-dev/server/internal/guild"
	"game-dev/server/internal/mail"
	"game-dev/server/internal/shop"
	"game-dev/server/internal/chat"
)

// FriendHandler 好友系统HTTP处理器
type FriendHandler struct {
	manager    *friend.FriendManager
	playerCache map[string]*friend.Friend
	mu         sync.RWMutex
}

// NewFriendHandler 创建好友处理器
func NewFriendHandler() *FriendHandler {
	return &FriendHandler{
		manager:    friend.NewFriendManager(),
		playerCache: make(map[string]*friend.Friend),
	}
}

// AddFriend 添加好友（通过发送好友请求）
func (h *FriendHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FromPlayerID string `json:"from_player_id"`
		ToPlayerID   string `json:"to_player_id"`
		ToPlayerName string `json:"to_player_name"`
		Message      string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := h.manager.SendFriendRequest(req.FromPlayerID, req.ToPlayerID, req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Friend request sent",
	})
}

// RemoveFriend 删除好友
func (h *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
		FriendID string `json:"friend_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.RemoveFriend(req.PlayerID, req.FriendID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Friend removed successfully",
	})
}

// GetFriends 获取好友列表
func (h *FriendHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	friends := h.manager.GetFriends(playerID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"friends": friends,
		"count": len(friends),
	})
}

// SendFriendRequest 发送好友请求
func (h *FriendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		FromPlayerID string `json:"from_player_id"`
		ToPlayerID   string `json:"to_player_id"`
		Message      string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	friendReq, err := h.manager.SendFriendRequest(req.FromPlayerID, req.ToPlayerID, req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"request": friendReq,
	})
}

// AcceptFriendRequest 接受好友请求
func (h *FriendHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID  string `json:"player_id"`
		RequestID string `json:"request_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.AcceptFriendRequest(req.RequestID, req.PlayerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Friend request accepted",
	})
}

// GetFriendRequests 获取好友请求列表
func (h *FriendHandler) GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	requests := h.manager.GetFriendRequests(playerID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"requests": requests,
		"count": len(requests),
	})
}

// GuildHandler 公会系统HTTP处理器
type GuildHandler struct {
	manager    *guild.GuildManager
	playerCache map[string]*guild.Guild
	mu         sync.RWMutex
}

// NewGuildHandler 创建公会处理器
func NewGuildHandler() *GuildHandler {
	return &GuildHandler{
		manager:    guild.NewGuildManager(),
		playerCache: make(map[string]*guild.Guild),
	}
}

// CreateGuild 创建公会
func (h *GuildHandler) CreateGuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		GuildName   string `json:"guild_name"`
		LeaderID    string `json:"leader_id"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	guildInfo, err := h.manager.CreateGuild(req.GuildName, req.LeaderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"guild": guildInfo,
	})
}

// JoinGuild 加入公会
func (h *GuildHandler) JoinGuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		GuildID  string `json:"guild_id"`
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.JoinGuild(req.GuildID, req.PlayerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Joined guild successfully",
	})
}

// LeaveGuild 离开公会
func (h *GuildHandler) LeaveGuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.LeaveGuild(req.PlayerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Left guild successfully",
	})
}

// GetGuildInfo 获取公会信息
func (h *GuildHandler) GetGuildInfo(w http.ResponseWriter, r *http.Request) {
	guildID := r.URL.Query().Get("guild_id")
	if guildID == "" {
		http.Error(w, "guild_id is required", http.StatusBadRequest)
		return
	}

	guildInfo, err := h.manager.GetGuild(guildID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"guild": guildInfo,
	})
}

// GetPlayerGuild 获取玩家公会
func (h *GuildHandler) GetPlayerGuild(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	guildInfo, err := h.manager.GetPlayerGuild(playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"guild": guildInfo,
	})
}

// DonateToGuild 贡献公会
func (h *GuildHandler) DonateToGuild(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		GuildID  string `json:"guild_id"`
		PlayerID string `json:"player_id"`
		Points   int    `json:"points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.AddContrib(req.GuildID, req.PlayerID, req.Points); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Donation successful",
	})
}

// ListGuilds 获取公会列表
func (h *GuildHandler) ListGuilds(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10
	
	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	guilds := h.manager.ListGuilds(page, pageSize)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"guilds": guilds,
		"count": len(guilds),
		"page": page,
		"page_size": pageSize,
	})
}

// MailHandler 邮件系统HTTP处理器
type MailHandler struct {
	manager *mail.MailManager
	mu      sync.RWMutex
}

// NewMailHandler 创建邮件处理器
func NewMailHandler() *MailHandler {
	return &MailHandler{
		manager: mail.NewMailManager(),
	}
}

// SendMail 发送邮件
func (h *MailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SenderID   string `json:"sender_id"`
		SenderName string `json:"sender_name"`
		ReceiverID string `json:"receiver_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mailInfo, err := h.manager.SendMail(req.ReceiverID, req.Title, req.Content, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"mail": mailInfo,
	})
}

// GetMail 获取邮件
func (h *MailHandler) GetMail(w http.ResponseWriter, r *http.Request) {
	mailID := r.URL.Query().Get("mail_id")
	if mailID == "" {
		http.Error(w, "mail_id is required", http.StatusBadRequest)
		return
	}

	mailInfo, err := h.manager.GetMail(mailID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"mail": mailInfo,
	})
}

// GetPlayerMails 获取玩家邮件列表
func (h *MailHandler) GetPlayerMails(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	mails, err := h.manager.GetInbox(playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"mails": mails,
		"count": len(mails),
	})
}

// ClaimMailAttachment 领取邮件附件
func (h *MailHandler) ClaimMailAttachment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		MailID   string `json:"mail_id"`
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	items, err := h.manager.ClaimAttachments(req.MailID, req.PlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"items": items,
	})
}

// DeleteMail 删除邮件
func (h *MailHandler) DeleteMail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		MailID   string `json:"mail_id"`
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.DeleteMail(req.MailID, req.PlayerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Mail deleted",
	})
}

// MarkAsRead 标记邮件为已读
func (h *MailHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		MailID   string `json:"mail_id"`
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.MarkAsRead(req.MailID, req.PlayerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Mail marked as read",
	})
}

// ShopHandler 商店系统HTTP处理器
type ShopHandler struct {
	manager *shop.ShopManager
	mu      sync.RWMutex
}

// NewShopHandler 创建商店处理器
func NewShopHandler() *ShopHandler {
	h := &ShopHandler{
		manager: shop.NewShopManager(),
	}
	h.manager.InitDefaultItems()
	return h
}

// GetShopItems 获取商店物品列表
func (h *ShopHandler) GetShopItems(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	
	var items []*shop.ShopItem
	
	if category != "" {
		items = h.manager.GetItemsByCategory(category)
	} else {
		items = h.manager.GetAllItems()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"items": items,
		"count": len(items),
	})
}

// PurchaseItem 购买物品
func (h *ShopHandler) PurchaseItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PlayerID  string `json:"player_id"`
		ItemID    string `json:"item_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	record, err := h.manager.Purchase(req.PlayerID, req.ItemID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"purchase": record,
	})
}

// GetPurchaseHistory 获取购买历史
func (h *ShopHandler) GetPurchaseHistory(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	history := h.manager.GetPurchaseHistory(playerID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"history": history,
		"count": len(history),
	})
}

// GetShopCategories 获取商店分类
func (h *ShopHandler) GetShopCategories(w http.ResponseWriter, r *http.Request) {
	categories := h.manager.GetCategories()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"categories": categories,
	})
}

// ChatHandler 聊天系统HTTP处理器
type ChatHandler struct {
	manager *chat.ChatManager
	mu      sync.RWMutex
}

// NewChatHandler 创建聊天处理器
func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		manager: chat.NewChatManager(),
	}
}

// SendMessage 发送消息
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SenderID    string          `json:"sender_id"`
		SenderName  string          `json:"sender_name"`
		SenderAvatar string         `json:"sender_avatar"`
		ChannelID   string          `json:"channel_id"`
		Content     string          `json:"content"`
		MessageType chat.MessageType `json:"message_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	message, err := h.manager.SendMessage(req.SenderID, req.SenderName, req.SenderAvatar, req.ChannelID, req.Content, req.MessageType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": message,
	})
}

// GetMessages 获取消息历史
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel_id")
	limit := 50

	if channelID == "" {
		http.Error(w, "channel_id is required", http.StatusBadRequest)
		return
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	messages := h.manager.GetChannelMessages(channelID, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"messages": messages,
		"count": len(messages),
	})
}

// CreateChannel 创建频道
func (h *ChatHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ChannelName string          `json:"channel_name"`
		ChannelType chat.ChannelType `json:"channel_type"`
		OwnerID     string          `json:"owner_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	channel, err := h.manager.CreateChannel(req.ChannelName, req.ChannelType, req.OwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"channel": channel,
	})
}

// JoinChannel 加入频道
func (h *ChatHandler) JoinChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ChannelID string `json:"channel_id"`
		PlayerID  string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.JoinChannel(req.PlayerID, req.ChannelID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Joined channel",
	})
}

// LeaveChannel 离开频道
func (h *ChatHandler) LeaveChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ChannelID string `json:"channel_id"`
		PlayerID  string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.LeaveChannel(req.PlayerID, req.ChannelID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Left channel",
	})
}

// DeleteMessage 删除消息
func (h *ChatHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		MessageID  string `json:"message_id"`
		ChannelID  string `json:"channel_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.manager.DeleteMessage(req.MessageID, req.ChannelID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Message deleted",
	})
}

// GetChannels 获取频道列表
func (h *ChatHandler) GetChannels(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	
	var channels []*chat.ChatChannel
	
	if playerID != "" {
		channels = h.manager.GetUserChannels(playerID)
	} else {
		// Return all channels
		channelID := r.URL.Query().Get("channel_id")
		if channelID != "" {
			ch, _ := h.manager.GetChannel(channelID)
			if ch != nil {
				channels = []*chat.ChatChannel{ch}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"channels": channels,
		"count": len(channels),
	})
}

// SendPrivateMessage 发送私信
func (h *ChatHandler) SendPrivateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SenderID   string `json:"sender_id"`
		ReceiverID string `json:"receiver_id"`
		Content    string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	message, err := h.manager.SendPrivateMessage(req.SenderID, req.ReceiverID, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": message,
	})
}

// InitHandler 初始化所有处理器
func InitHandler() {
	// 注册所有HTTP处理器
	// 这些处理器将在HTTP服务器启动时注册
}

// GetAllHandlers 获取所有处理器（用于注册路由）
func GetAllHandlers() map[string]interface{} {
	return map[string]interface{}{
		"friend": NewFriendHandler(),
		"guild": NewGuildHandler(),
		"mail": NewMailHandler(),
		"shop": NewShopHandler(),
		"chat": NewChatHandler(),
	}
}

var (
	FriendHandlerInstance *FriendHandler
	GuildHandlerInstance *GuildHandler
	MailHandlerInstance  *MailHandler
	ShopHandlerInstance  *ShopHandler
	ChatHandlerInstance  *ChatHandler
)

func init() {
	FriendHandlerInstance = NewFriendHandler()
	GuildHandlerInstance = NewGuildHandler()
	MailHandlerInstance = NewMailHandler()
	ShopHandlerInstance = NewShopHandler()
	ChatHandlerInstance = NewChatHandler()
}
