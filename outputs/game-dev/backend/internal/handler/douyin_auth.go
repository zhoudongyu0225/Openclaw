package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ========== 抖音授权配置 ==========

type DouyinAuthConfig struct {
	ClientKey    string // 应用 Key
	ClientSecret string // 应用 Secret
	RedirectURI  string // 授权回调地址
	Scope        string // 权限范围
}

// 抖音 OAuth 接口
const (
	DouyinAuthURL   = "https://open.douyin.com/oauth/access_token/"
	DouyinUserInfoURL = "https://open.douyin.com/oauth/userinfo/"
)

// ========== 授权流程 ==========

// Step 1: 生成授权链接
func GenerateAuthURL(config DouyinAuthConfig, state string) string {
	return fmt.Sprintf("https://open.douyin.com/oauth/authorize/?client_key=%s&redirect_uri=%s&scope=%s&state=%s&response_type=code",
		config.ClientKey,
		config.RedirectURI,
		config.Scope,
		state,
	)
}

// Step 2: 通过 code 获取 access_token
type DouyinTokenResponse struct {
	ErrNo   int    `json:"err_no"`
	ErrMsg  string `json:"err_msg"`
	Data    TokenData `json:"data"`
}

type TokenData struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"open_id"`
}

func ExchangeToken(clientKey, clientSecret, code string) (*TokenData, error) {
	url := fmt.Sprintf("%s?client_key=%s&client_secret=%s&code=%s&grant_type=authorization_code",
		DouyinAuthURL, clientKey, clientSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result DouyinTokenResponse
	json.Unmarshal(body, &result)

	if result.ErrNo != 0 {
		return nil, fmt.Errorf("抖音授权失败: %s", result.ErrMsg)
	}

	return &result.Data, nil
}

// Step 3: 刷新 access_token
func RefreshToken(clientKey, refreshToken string) (*TokenData, error) {
	url := fmt.Sprintf("%s?client_key=%s&refresh_token=%s&grant_type=refresh_token",
		DouyinAuthURL, clientKey, refreshToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result DouyinTokenResponse
	json.Unmarshal(body, &result)

	if result.ErrNo != 0 {
		return nil, fmt.Errorf("刷新Token失败: %s", result.ErrMsg)
	}

	return &result.Data, nil
}

// Step 4: 获取用户信息
type DouyinUserInfoResponse struct {
	ErrNo  int       `json:"err_no"`
	ErrMsg string    `json:"err_msg"`
	Data   UserInfoData `json:"data"`
}

type UserInfoData struct {
	OpenID    string `json:"open_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Gender    int    `json:"gender"` // 0:未知 1:男 2:女
	Country   string `json:"country"`
	Province  string `json:"province"`
	City      string `json:"city"`
}

func GetUserInfo(accessToken, openID string) (*UserInfoData, error) {
	url := fmt.Sprintf("%s?access_token=%s&open_id=%s", DouyinUserInfoURL, accessToken, openID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result DouyinUserInfoResponse
	json.Unmarshal(body, &result)

	if result.ErrNo != 0 {
		return nil, fmt.Errorf("获取用户信息失败: %s", result.ErrMsg)
	}

	return &result.Data, nil
}

// ========== 抖音登录管理器 ==========

type DouyinLoginManager struct {
	config     DouyinAuthConfig
	tokenCache map[string]*TokenData // openID -> TokenData
}

func NewDouyinLoginManager(config DouyinAuthConfig) *DouyinLoginManager {
	return &DouyinLoginManager{
		config:     config,
		tokenCache: make(map[string]*TokenData),
	}
}

// 处理授权回调
func (m *DouyinLoginManager) HandleCallback(code, state string) (*UserInfoData, error) {
	// 1. 换取 Token
	token, err := ExchangeToken(m.config.ClientKey, m.config.ClientSecret, code)
	if err != nil {
		return nil, err
	}

	// 2. 缓存 Token
	m.tokenCache[token.OpenID] = token

	// 3. 获取用户信息
	userInfo, err := GetUserInfo(token.AccessToken, token.OpenID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// 获取已登录用户信息
func (m *DouyinLoginManager) GetUserInfo(openID string) (*UserInfoData, error) {
	token, ok := m.tokenCache[openID]
	if !ok {
		return nil, fmt.Errorf("用户未登录")
	}

	// 检查 Token 是否过期
	// 实际应用中需要存储 Token 获取时间
	return GetUserInfo(token.AccessToken, openID)
}

// 刷新用户 Token
func (m *DouyinLoginManager) RefreshUserToken(openID string) error {
	token, ok := m.tokenCache[openID]
	if !ok {
		return fmt.Errorf("用户未登录")
	}

	newToken, err := RefreshToken(m.config.ClientKey, token.RefreshToken)
	if err != nil {
		return err
	}

	m.tokenCache[openID] = newToken
	return nil
}

// ========== HTTP 处理器 ==========

type AuthHandler struct {
	loginManager *DouyinLoginManager
}

func NewAuthHandler(config DouyinAuthConfig) *AuthHandler {
	return &AuthHandler{
		loginManager: NewDouyinLoginManager(config),
	}
}

// 授权页面跳转
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		state = time.Now().Format("20060102150405")
	}

	authURL := GenerateAuthURL(h.loginManager.config, state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// 授权回调
func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing code"))
		return
	}

	userInfo, err := h.loginManager.HandleCallback(code, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// 返回用户信息（实际应该生成 JWT 或 Session）
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    userInfo,
	})
}

// 验证登录状态
func (h *AuthHandler) HandleCheckLogin(w http.ResponseWriter, r *http.Request) {
	openID := r.URL.Query().Get("open_id")
	if openID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo, err := h.loginManager.GetUserInfo(openID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
