package game

import (
	"encoding/json"
	"fmt"
	"time"
)

// ============================================
// 称号系统 (Title System)
// ============================================

// TitleRarity 称号稀有度
type TitleRarity int

const (
	TitleRarityCommon    TitleRarity = iota // 普通
	TitleRarityUncommon                     // 优秀
	TitleRarityRare                         // 稀有
	TitleRarityEpic                         // 史诗
	TitleRarityLegendary                    // 传说
)

// TitleSource 称号来源
type TitleSource int

const (
	TitleSourceAchievement TitleSource = iota // 成就解锁
	TitleSourceQuest                          // 任务解锁
	TitleSourceLevel                          // 等级解锁
	TitleSourceBattle                        // 战斗解锁
	TitleSourceEvent                          // 活动解锁
	TitleSourcePurchase                       // 购买解锁
	TitleSourceSpecial                        // 特殊解锁
)

// Title 称号
type Title struct {
	TitleID     string       `json:"title_id"`     // 称号ID
	Name        string       `json:"name"`         // 称号名称
	Description string       `json:"description"`  // 称号描述
	Rarity      TitleRarity  `json:"rarity"`       // 稀有度
	Source      TitleSource `json:"source"`       // 来源类型
	SourceID    string      `json:"source_id"`     // 来源ID
	Icon        string      `json:"icon"`          // 图标
	Attributes  TitleAttribute `json:"attributes"` // 属性加成
	IsPermanent bool        `json:"is_permanent"`  // 是否永久
	ExpireTime  int64       `json:"expire_time"`   // 过期时间
	RequiredLevel int      `json:"required_level"`// 所需等级
	SortOrder   int        `json:"sort_order"`    // 排序
}

// TitleAttribute 称号属性
type TitleAttribute struct {
	Attack    float64 `json:"attack"`    // 攻击加成
	Defense   float64 `json:"defense"`   // 防御加成
	HP        float64 `json:"hp"`         // 生命加成
	Critical  float64 `json:"critical"`  // 暴击加成
	Dodge     float64 `json:"dodge"`     // 闪避加成
	MoveSpeed float64 `json:"move_speed"`// 移速加成
}

// PlayerTitle 玩家称号
type PlayerTitle struct {
	PlayerID     string            `json:"player_id"`      // 玩家ID
	OwnedTitles  map[string]*Title `json:"owned_titles"`  // 拥有的称号
	EquippedTitle string           `json:"equipped_title"` // 当前装备的称号
	ActiveUntil  int64             `json:"active_until"`   // 有效期限
}

// TitleManager 称号管理器
type TitleManager struct {
	titles      map[string]*Title  // 称号模板
	categories  map[string][]string // 称号分类
}

// NewTitleManager 创建称号管理器
func NewTitleManager() *TitleManager {
	mgr := &TitleManager{
	titles:     make(map[string]*Title),
	categories: make(map[string][]string),
	}
	mgr.initTitles()
	mgr.initCategories()
	return mgr
}

// initTitles 初始化称号
func (m *TitleManager) initTitles() {
	// 普通称号
	m.titles["title_newbie"] = &Title{
		TitleID:     "title_newbie",
		Name:        "初入江湖",
		Description: "完成新手引导",
		Rarity:      TitleRarityCommon,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_newbie",
		Icon:        "icon_title_newbie",
		Attributes:  TitleAttribute{Attack: 0.01},
		IsPermanent: true,
		RequiredLevel: 1,
		SortOrder:   1,
	}

	m.titles["title_warrior"] = &Title{
		TitleID:     "title_warrior",
		Name:        "初级战士",
		Description: "累计击杀100个敌人",
		Rarity:      TitleRarityCommon,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_kill_100",
		Icon:        "icon_title_warrior",
		Attributes:  TitleAttribute{Attack: 0.02},
		IsPermanent: true,
		RequiredLevel: 5,
		SortOrder:   2,
	}

	m.titles["title_veteran"] = &Title{
		TitleID:     "title_veteran",
		Name:        "老练战士",
		Description: "累计击杀1000个敌人",
		Rarity:      TitleRarityUncommon,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_kill_1000",
		Icon:        "icon_title_veteran",
		Attributes:  TitleAttribute{Attack: 0.05, Defense: 0.02},
		IsPermanent: true,
		RequiredLevel: 15,
		SortOrder:   3,
	}

	// 稀有称号
	m.titles["title_combo_master"] = &Title{
		TitleID:     "title_combo_master",
		Name:        "连击大师",
		Description: "达成100连击",
		Rarity:      TitleRarityRare,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_combo_100",
		Icon:        "icon_title_combo",
		Attributes:  TitleAttribute{Attack: 0.08, Critical: 0.05},
		IsPermanent: true,
		RequiredLevel: 10,
		SortOrder:   10,
	}

	m.titles["title_boss_killer"] = &Title{
		TitleID:     "title_boss_killer",
		Name:        "Boss杀手",
		Description: "累计击杀10个Boss",
		Rarity:      TitleRarityRare,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_boss_10",
		Icon:        "icon_title_boss",
		Attributes:  TitleAttribute{Attack: 0.10, HP: 0.05},
		IsPermanent: true,
		RequiredLevel: 20,
		SortOrder:   11,
	}

	m.titles["title_survivor"] = &Title{
		TitleID:     "title_survivor",
		Name:        "绝境求生",
		Description: "在1%血量下获胜",
		Rarity:      TitleRarityRare,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_survival",
		Icon:        "icon_title_survivor",
		Attributes:  TitleAttribute{HP: 0.10, Defense: 0.05},
		IsPermanent: true,
		RequiredLevel: 15,
		SortOrder:   12,
	}

	// 史诗称号
	m.titles["title_legend"] = &Title{
		TitleID:     "title_legend",
		Name:        "传奇英雄",
		Description: "获得最强王者段位",
		Rarity:      TitleRarityEpic,
		Source:      TitleSourceBattle,
		SourceID:    "rank_king",
		Icon:        "icon_title_legend",
		Attributes:  TitleAttribute{Attack: 0.15, Defense: 0.10, HP: 0.10, Critical: 0.08},
		IsPermanent: true,
		RequiredLevel: 30,
		SortOrder:   20,
	}

	m.titles["title_collector"] = &Title{
		TitleID:     "title_collector",
		Name:        "收藏家",
		Description: "收集全部装备类型",
		Rarity:      TitleRarityEpic,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_collector",
		Icon:        "icon_title_collector",
		Attributes:  TitleAttribute{Attack: 0.12, Defense: 0.12},
		IsPermanent: true,
		RequiredLevel: 25,
		SortOrder:   21,
	}

	// 传说称号
	m.titles["title_god"] = &Title{
		TitleID:     "title_god",
		Name:        "弹幕之神",
		Description: "同时触发10种弹幕效果",
		Rarity:      TitleRarityLegendary,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_danmaku_god",
		Icon:        "icon_title_god",
		Attributes:  TitleAttribute{
			Attack: 0.20, Defense: 0.15, HP: 0.15,
			Critical: 0.10, Dodge: 0.10, MoveSpeed: 0.05,
		},
		IsPermanent: true,
		RequiredLevel: 40,
		SortOrder:   30,
	}

	m.titles["title_millionaire"] = &Title{
		TitleID:     "title_millionaire",
		Name:        "千万富翁",
		Description: "累计拥有10000000金币",
		Rarity:      TitleRarityLegendary,
		Source:      TitleSourceAchievement,
		SourceID:    "achievement_millionaire",
		Icon:        "icon_title_millionaire",
		Attributes:  TitleAttribute{Attack: 0.18, Defense: 0.18},
		IsPermanent: true,
		RequiredLevel: 35,
		SortOrder:   31,
	}

	// 限时称号
	m.titles["title_event_2024"] = &Title{
		TitleID:     "title_event_2024",
		Name:        "2024周年庆",
		Description: "参与2024周年庆活动",
		Rarity:      TitleRarityEpic,
		Source:      TitleSourceEvent,
		SourceID:    "event_2024_anniversary",
		Icon:        "icon_title_event",
		Attributes:  TitleAttribute{Attack: 0.10, HP: 0.10},
		IsPermanent: false,
		ExpireTime:  time.Now().Add(30 * 24 * time.Hour).Unix(),
		RequiredLevel: 1,
		SortOrder:   50,
	}
}

// initCategories 初始化称号分类
func (m *TitleManager) initCategories() {
	m.categories["combat"] = []string{
		"title_warrior", "title_veteran", "title_boss_killer", "title_survivor",
	}
	m.categories["skill"] = []string{
		"title_combo_master", "title_god",
	}
	m.categories["rank"] = []string{
		"title_legend",
	}
	m.categories["collection"] = []string{
		"title_collector", "title_millionaire",
	}
	m.categories["event"] = []string{
		"title_event_2024",
	}
}

// GetTitle 获取称号
func (m *TitleManager) GetTitle(titleID string) (*Title, bool) {
	title, ok := m.titles[titleID]
	return title, ok
}

// GetAllTitles 获取所有称号
func (m *TitleManager) GetAllTitles() []*Title {
	titles := make([]*Title, 0, len(m.titles))
	for _, title := range m.titles {
		titles = append(titles, title)
	}
	return titles
}

// GetTitleByCategory 按分类获取称号
func (m *TitleManager) GetTitleByCategory(category string) []*Title {
	titleIDs := m.categories[category]
	titles := make([]*Title, 0, len(titleIDs))
	for _, id := range titleIDs {
		if title, ok := m.titles[id]; ok {
			titles = append(titles, title)
		}
	}
	return titles
}

// UnlockTitle 解锁称号
func (m *TitleManager) UnlockTitle(playerTitle *PlayerTitle, titleID string) error {
	title, ok := m.titles[titleID]
	if !ok {
		return fmt.Errorf("title not found: %s", titleID)
	}

	// 检查是否已拥有
	if _, exists := playerTitle.OwnedTitles[titleID]; exists {
		return fmt.Errorf("title already owned: %s", titleID)
	}

	// 添加称号
	newTitle := *title
	playerTitle.OwnedTitles[titleID] = &newTitle

	return nil
}

// EquipTitle 装备称号
func (m *TitleManager) EquipTitle(playerTitle *PlayerTitle, titleID string) error {
	// 检查是否拥有该称号
	title, exists := playerTitle.OwnedTitles[titleID]
	if !exists {
		return fmt.Errorf("title not owned: %s", titleID)
	}

	// 检查是否过期
	if !title.IsPermanent && title.ExpireTime < time.Now().Unix() {
		return fmt.Errorf("title expired: %s", titleID)
	}

	playerTitle.EquippedTitle = titleID
	return nil
}

// UnequipTitle 卸下称号
func (m *TitleManager) UnequipTitle(playerTitle *PlayerTitle) {
	playerTitle.EquippedTitle = ""
}

// GetEquippedTitle 获取已装备称号
func (m *TitleManager) GetEquippedTitle(playerTitle *PlayerTitle) *Title {
	if playerTitle.EquippedTitle == "" {
		return nil
	}
	return playerTitle.OwnedTitles[playerTitle.EquippedTitle]
}

// GetEquippedAttributes 获取装备称号的属性加成
func (m *TitleManager) GetEquippedAttributes(playerTitle *PlayerTitle) TitleAttribute {
	equipped := m.GetEquippedTitle(playerTitle)
	if equipped == nil {
		return TitleAttribute{}
	}
	return equipped.Attributes
}

// CreatePlayerTitle 创建玩家称号数据
func (m *TitleManager) CreatePlayerTitle(playerID string) *PlayerTitle {
	return &PlayerTitle{
		PlayerID:     playerID,
		OwnedTitles:  make(map[string]*Title),
		EquippedTitle: "",
		ActiveUntil:  0,
	}
}

// GetAvailableTitles 获取可解锁的称号
func (m *TitleManager) GetAvailableTitles(playerTitle *PlayerTitle, level int) []*Title {
	available := make([]*Title, 0)
	for _, title := range m.titles {
		// 排除已拥有的
		if _, owned := playerTitle.OwnedTitles[title.TitleID]; owned {
			continue
		}
		// 检查等级
		if title.RequiredLevel > level {
			continue
		}
		// 排除已过期的限时称号
		if !title.IsPermanent && title.ExpireTime < time.Now().Unix() {
			continue
		}
		available = append(available, title)
	}
	return available
}

// CheckTitleExpiration 检查称号过期
func (m *TitleManager) CheckTitleExpiration(playerTitle *PlayerTitle) {
	now := time.Now().Unix()
	for titleID, title := range playerTitle.OwnedTitles {
		if !title.IsPermanent && title.ExpireTime < now {
			delete(playerTitle.OwnedTitles, titleID)
			// 如果装备的称号过期，卸下
			if playerTitle.EquippedTitle == titleID {
				playerTitle.EquippedTitle = ""
			}
		}
	}
}

// GetTitleList 获取玩家称号列表
func (m *TitleManager) GetTitleList(playerTitle *PlayerTitle) []*Title {
	titles := make([]*Title, 0, len(playerTitle.OwnedTitles))
	for _, title := range playerTitle.OwnedTitles {
		titles = append(titles, title)
	}
	return titles
}

// ============================================
// 序列化方法
// ============================================

// MarshalJSON 序列化
func (t *Title) MarshalJSON() ([]byte, error) {
	type Alias Title
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})
}
