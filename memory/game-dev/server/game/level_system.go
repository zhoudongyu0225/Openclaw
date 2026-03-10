// Copyright 2026 GameDev.Daily - Bullet Hell Game
// Level Chapter System - Manages game progression through chapters and levels

package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ChapterLevel represents a single level in a chapter
type ChapterLevel struct {
	ID            int                    `json:"id"`
	ChapterID     int                    `json:"chapter_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Difficulty    Difficulty             `json:"difficulty"`
	RecommendedLevel int                  `json:"recommended_level"`
	Stars         [3]bool                `json:"stars"` // 3-star rating
	BestScore     int64                  `json:"best_score"`
	BestTime      time.Duration          `json:"best_time"`
	IsUnlocked    bool                   `json:"is_unlocked"`
	IsCompleted   bool                   `json:"is_completed"`
	UnlockCost    int                    `json:"unlock_cost"` // Coins to unlock
	EntryFee      int                    `json:"entry_fee"`   // Coins to play
	Rewards       LevelRewards           `json:"rewards"`
	Enemies       []LevelEnemySpawn      `json:"enemies"`
	BossID        int                    `json:"boss_id"`
	TimeLimit     time.Duration          `json:"time_limit"`
	KillTarget    int                    `json:"kill_target"` // Enemies to kill
	ScoreTarget   int64                  `json:"score_target"`
}

// LevelRewards represents rewards for completing a level
type LevelRewards struct {
	Coins       int            `json:"coins"`
	Gems        int            `json:"gems"`
	EXP         int            `json:"exp"`
	Items       []RewardItem   `json:"items"`
	StarRewards []StarReward   `json:"star_rewards"` // Rewards per star
}

// RewardItem represents an item reward
type RewardItem struct {
	ItemID   string `json:"item_id"`
	ItemType string `json:"item_type"`
	Quantity int    `json:"quantity"`
	DropRate float64 `json:"drop_rate"`
}

// StarReward represents reward for achieving a star rating
type StarReward struct {
	Star      int    `json:"star"`
	Coins     int    `json:"coins"`
	Gems      int    `json:"gems"`
	Items     []RewardItem `json:"items"`
}

// LevelEnemySpawn defines enemy spawn configuration
type LevelEnemySpawn struct {
	EnemyType  EnemyType     `json:"enemy_type"`
	Wave       int           `json:"wave"`
	Count      int           `json:"count"`
	Interval   time.Duration `json:"interval"`
	Position   SpawnPosition `json:"position"`
	Formation  string       `json:"formation"` // line, circle, random
}

// Difficulty level for chapters and levels
type Difficulty int

const (
	DifficultyNormal Difficulty = iota
	DifficultyHard
	DifficultyExpert
	DifficultyNightmare
)

// Chapter represents a chapter containing multiple levels
type Chapter struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Difficulty  Difficulty    `json:"difficulty"`
	Levels      []*ChapterLevel `json:"levels"`
	IsUnlocked  bool          `json:"is_unlocked"`
	IsCompleted bool          `json:"is_completed"`
	Stars       int           `json:"stars"`        // Total stars
	MaxStars    int           `json:"max_stars"`    // Max possible stars
	UnlockCost  int           `json:"unlock_cost"`
	EntryCost   int           `json:"entry_cost"`
	ChapterEXP  int           `json:"chapter_exp"`  // EXP reward for completing chapter
	Background  string        `json:"background"`  // Background asset
	Music       string        `json:"music"`        // Background music
	Thumbnail   string        `json:"thumbnail"`    // Chapter thumbnail
}

// ChapterManager manages all chapters and levels
type ChapterManager struct {
	Chapters    []*Chapter
	Levels      map[int]*ChapterLevel // levelID -> level
	PlayerProgress map[int]*PlayerLevelProgress
	rand        *rand.Rand
}

// PlayerLevelProgress tracks player's progress in a level
type PlayerLevelProgress struct {
	PlayerID    string    `json:"player_id"`
	LevelID     int       `json:"level_id"`
	Stars       int       `json:"stars"`
	BestScore   int64     `json:"best_score"`
	BestTime    time.Duration `json:"best_time"`
	PlayCount   int       `json:"play_count"`
	WinCount    int       `json:"win_count"`
	LastPlayAt  time.Time `json:"last_play_at"`
	FirstWinAt  time.Time `json:"first_win_at"`
	BestCombo   int       `json:"best_combo"`
	TotalKills  int       `json:"total_kills"`
	TotalDamage int64     `json:"total_damage"`
}

// NewChapterManager creates a new chapter manager
func NewChapterManager() *ChapterManager {
	cm := &ChapterManager{
		Chapters:        make([]*Chapter, 0),
		Levels:          make(map[int]*ChapterLevel),
		PlayerProgress:  make(map[int]*PlayerLevelProgress),
		rand:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	cm.initializeChapters()
	return cm
}

// initializeChapters sets up the default chapter structure
func (cm *ChapterManager) initializeChapters() {
	// Chapter 1: Tutorial - The Beginning
	chapter1 := &Chapter{
		ID:          1,
		Name:        "序章：启程",
		Description: "Learn the basics of bullet hell combat",
		Difficulty:  DifficultyNormal,
		IsUnlocked:  true,
		UnlockCost:  0,
		EntryCost:   0,
		ChapterEXP: 100,
		Background: "tutorial_bg",
		Music:       "tutorial_theme",
		Thumbnail:   "chapter1_thumb",
		Levels:      make([]*ChapterLevel, 0),
	}

	// Level 1-1: First Steps
	level1_1 := &ChapterLevel{
		ID:              101,
		ChapterID:      1,
		Name:            "初试身手",
		Description:     "Your first battle. Defeat the basic enemies!",
		Difficulty:     DifficultyNormal,
		RecommendedLevel: 1,
		IsUnlocked:     true,
		UnlockCost:      0,
		EntryFee:        0,
		TimeLimit:       3 * time.Minute,
		KillTarget:      10,
		ScoreTarget:     1000,
		Rewards: LevelRewards{
			Coins:    100,
			Gems:     5,
			EXP:      20,
			Items:    []RewardItem{{ItemID: "weapon_001", ItemType: "Weapon", Quantity: 1}},
			StarRewards: []StarReward{
				{Star: 1, Coins: 50},
				{Star: 2, Coins: 100},
				{Star: 3, Coins: 200},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeBasic, Wave: 1, Count: 5, Interval: 2 * time.Second, Position: SpawnTop, Formation: "line"},
			{EnemyType: EnemyTypeBasic, Wave: 2, Count: 5, Interval: 2 * time.Second, Position: SpawnTop, Formation: "line"},
		},
	}

	// Level 1-2: Dodging Practice
	level1_2 := &ChapterLevel{
		ID:              102,
		ChapterID:      1,
		Name:            "躲避训练",
		Description:     "Learn to dodge bullet patterns",
		Difficulty:     DifficultyNormal,
		RecommendedLevel: 3,
		IsUnlocked:     false,
		UnlockCost:      0,
		EntryFee:        0,
		TimeLimit:       2 * time.Minute,
		KillTarget:      15,
		ScoreTarget:     2000,
		Rewards: LevelRewards{
			Coins:    150,
			Gems:     10,
			EXP:      30,
			Items:    []RewardItem{},
			StarRewards: []StarReward{
				{Star: 1, Coins: 75},
				{Star: 2, Coins: 150},
				{Star: 3, Coins: 300},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeShooter, Wave: 1, Count: 3, Interval: 3 * time.Second, Position: SpawnTop, Formation: "circle"},
			{EnemyType: EnemyTypeShooter, Wave: 2, Count: 5, Interval: 2 * time.Second, Position: SpawnTop, Formation: "random"},
		},
	}

	// Level 1-3: First Boss
	level1_3 := &ChapterLevel{
		ID:              103,
		ChapterID:      1,
		Name:            "初次 boss 战",
		Description:     "Face your first boss - The Guardian",
		Difficulty:     DifficultyNormal,
		RecommendedLevel: 5,
		IsUnlocked:     false,
		UnlockCost:      0,
		EntryFee:        0,
		TimeLimit:       5 * time.Minute,
		KillTarget:      1,
		ScoreTarget:     5000,
		BossID:          1001,
		Rewards: LevelRewards{
			Coins:    500,
			Gems:     50,
			EXP:      100,
			Items:    []RewardItem{{ItemID: "skill_001", ItemType: "Skill", Quantity: 1}},
			StarRewards: []StarReward{
				{Star: 1, Coins: 250},
				{Star: 2, Coins: 500},
				{Star: 3, Coins: 1000},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeBasic, Wave: 1, Count: 5, Interval: 2 * time.Second, Position: SpawnTop, Formation: "line"},
		},
	}

	chapter1.Levels = []*ChapterLevel{level1_1, level1_2, level1_3}
	chapter1.MaxStars = 9 // 3 stars * 3 levels

	// Chapter 2: The Forest of Echoes
	chapter2 := &Chapter{
		ID:          2,
		Name:        "回声之森",
		Description: "Navigate through the mysterious forest",
		Difficulty:  DifficultyNormal,
		IsUnlocked:  false,
		UnlockCost:  500,
		EntryCost:   100,
		ChapterEXP:  200,
		Background:  "forest_bg",
		Music:       "forest_theme",
		Thumbnail:   "chapter2_thumb",
		Levels:      make([]*ChapterLevel, 0),
	}

	level2_1 := &ChapterLevel{
		ID:              201,
		ChapterID:      2,
		Name:            "林间小路",
		Description:     "The path through the forest",
		Difficulty:     DifficultyNormal,
		RecommendedLevel: 7,
		IsUnlocked:     false,
		UnlockCost:      200,
		EntryFee:        50,
		TimeLimit:       3 * time.Minute,
		KillTarget:      20,
		ScoreTarget:     3000,
		Rewards: LevelRewards{
			Coins:    200,
			Gems:     15,
			EXP:      40,
			Items:    []RewardItem{},
			StarRewards: []StarReward{
				{Star: 1, Coins: 100},
				{Star: 2, Coins: 200},
				{Star: 3, Coins: 400},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeShooter, Wave: 1, Count: 8, Interval: 1.5 * time.Second, Position: SpawnTop, Formation: "random"},
			{EnemyType: EnemyTypeFast, Wave: 2, Count: 4, Interval: 2 * time.Second, Position: SpawnSides, Formation: "line"},
		},
	}

	level2_2 := &ChapterLevel{
		ID:              202,
		ChapterID:      2,
		Name:            "蘑菇领地",
		Description:     "Defeat the mushroom guardians",
		Difficulty:     DifficultyNormal,
		RecommendedLevel: 9,
		IsUnlocked:     false,
		UnlockCost:      250,
		EntryFee:        75,
		TimeLimit:       4 * time.Minute,
		KillTarget:      25,
		ScoreTarget:     4000,
		Rewards: LevelRewards{
			Coins:    300,
			Gems:     20,
			EXP:      50,
			Items:    []RewardItem{{ItemID: "material_001", ItemType: "Material", Quantity: 5}},
			StarRewards: []StarReward{
				{Star: 1, Coins: 150},
				{Star: 2, Coins: 300},
				{Star: 3, Coins: 600},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeTank, Wave: 1, Count: 3, Interval: 4 * time.Second, Position: SpawnTop, Formation: "line"},
			{EnemyType: EnemyTypeShooter, Wave: 2, Count: 6, Interval: 2 * time.Second, Position: SpawnTop, Formation: "circle"},
		},
	}

	level2_3 := &ChapterLevel{
		ID:              203,
		ChapterID:      2,
		Name:            "树妖 boss 战",
		Description:     "The Ancient Tree Guardian",
		Difficulty:     DifficultyHard,
		RecommendedLevel: 12,
		IsUnlocked:     false,
		UnlockCost:      300,
		EntryFee:        100,
		TimeLimit:       6 * time.Minute,
		KillTarget:      1,
		ScoreTarget:     8000,
		BossID:          1002,
		Rewards: LevelRewards{
			Coins:    800,
			Gems:     100,
			EXP:      200,
			Items:    []RewardItem{{ItemID: "armor_001", ItemType: "Armor", Quantity: 1}},
			StarRewards: []StarReward{
				{Star: 1, Coins: 400},
				{Star: 2, Coins: 800},
				{Star: 3, Coins: 1600},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeShooter, Wave: 1, Count: 10, Interval: 1 * time.Second, Position: SpawnTop, Formation: "random"},
		},
	}

	chapter2.Levels = []*ChapterLevel{level2_1, level2_2, level2_3}
	chapter2.MaxStars = 9

	// Chapter 3: Volcanic Wasteland
	chapter3 := &Chapter{
		ID:          3,
		Name:        "火山废土",
		Description: "Survive the molten landscape",
		Difficulty:  DifficultyHard,
		IsUnlocked:  false,
		UnlockCost:  1000,
		EntryCost:   200,
		ChapterEXP:  400,
		Background:  "volcano_bg",
		Music:       "volcano_theme",
		Thumbnail:   "chapter3_thumb",
		Levels:      make([]*ChapterLevel, 0),
	}

	level3_1 := &ChapterLevel{
		ID:              301,
		ChapterID:      3,
		Name:            "熔岩边缘",
		Description:     "Edge of the lava fields",
		Difficulty:     DifficultyHard,
		RecommendedLevel: 15,
		IsUnlocked:     false,
		UnlockCost:      500,
		EntryFee:        150,
		TimeLimit:       4 * time.Minute,
		KillTarget:      30,
		ScoreTarget:     6000,
		Rewards: LevelRewards{
			Coins:    400,
			Gems:     30,
			EXP:      80,
			Items:    []RewardItem{},
			StarRewards: []StarReward{
				{Star: 1, Coins: 200},
				{Star: 2, Coins: 400},
				{Star: 3, Coins: 800},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeFast, Wave: 1, Count: 8, Interval: 1.5 * time.Second, Position: SpawnSides, Formation: "line"},
			{EnemyType: EnemyTypeShooter, Wave: 2, Count: 6, Interval: 2 * time.Second, Position: SpawnTop, Formation: "circle"},
			{EnemyType: EnemyTypeTank, Wave: 3, Count: 2, Interval: 5 * time.Second, Position: SpawnTop, Formation: "line"},
		},
	}

	level3_2 := &ChapterLevel{
		ID:              302,
		ChapterID:      3,
		Name:            "火焰洞穴",
		Description:     "Deep within the volcanic caves",
		Difficulty:     DifficultyHard,
		RecommendedLevel: 18,
		IsUnlocked:     false,
		UnlockCost:      600,
		EntryFee:        200,
		TimeLimit:       5 * time.Minute,
		KillTarget:      40,
		ScoreTarget:     8000,
		Rewards: LevelRewards{
			Coins:    500,
			Gems:     40,
			EXP:      100,
			Items:    []RewardItem{{ItemID: "material_002", ItemType: "Material", Quantity: 10}},
			StarRewards: []StarReward{
				{Star: 1, Coins: 250},
				{Star: 2, Coins: 500},
				{Star: 3, Coins: 1000},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeShooter, Wave: 1, Count: 10, Interval: 1 * time.Second, Position: SpawnTop, Formation: "random"},
			{EnemyType: EnemyTypeFast, Wave: 2, Count: 8, Interval: 1.5 * time.Second, Position: SpawnSides, Formation: "random"},
			{EnemyType: EnemyTypeBossMinion, Wave: 3, Count: 3, Interval: 3 * time.Second, Position: SpawnTop, Formation: "circle"},
		},
	}

	level3_3 := &ChapterLevel{
		ID:              303,
		ChapterID:      3,
		Name:            "火山巨龙 boss 战",
		Description:     "The Fire Dragon Awakens",
		Difficulty:     DifficultyExpert,
		RecommendedLevel: 25,
		IsUnlocked:     false,
		UnlockCost:      800,
		EntryFee:        300,
		TimeLimit:       8 * time.Minute,
		KillTarget:      1,
		ScoreTarget:     15000,
		BossID:          1003,
		Rewards: LevelRewards{
			Coins:    1500,
			Gems:     200,
			EXP:      500,
			Items:    []RewardItem{{ItemID: "weapon_002", ItemType: "Weapon", Quantity: 1}},
			StarRewards: []StarReward{
				{Star: 1, Coins: 750},
				{Star: 2, Coins: 1500},
				{Star: 3, Coins: 3000},
			},
		},
		Enemies: []LevelEnemySpawn{
			{EnemyType: EnemyTypeShooter, Wave: 1, Count: 15, Interval: 0.8 * time.Second, Position: SpawnTop, Formation: "random"},
		},
	}

	chapter3.Levels = []*ChapterLevel{level3_1, level3_2, level3_3}
	chapter3.MaxStars = 9

	// Add chapters to manager
	cm.Chapters = []*Chapter{chapter1, chapter2, chapter3}

	// Add levels to lookup
	cm.Levels[101] = level1_1
	cm.Levels[102] = level1_2
	cm.Levels[103] = level1_3
	cm.Levels[201] = level2_1
	cm.Levels[202] = level2_2
	cm.Levels[203] = level2_3
	cm.Levels[301] = level3_1
	cm.Levels[302] = level3_2
	cm.Levels[303] = level3_3
}

// GetChapter returns a chapter by ID
func (cm *ChapterManager) GetChapter(chapterID int) *Chapter {
	for _, ch := range cm.Chapters {
		if ch.ID == chapterID {
			return ch
		}
	}
	return nil
}

// GetLevel returns a level by ID
func (cm *ChapterManager) GetLevel(levelID int) *ChapterLevel {
	return cm.Levels[levelID]
}

// UnlockLevel unlocks a level for a player
func (cm *ChapterManager) UnlockLevel(levelID int, playerID string) error {
	level, ok := cm.Levels[levelID]
	if !ok {
		return fmt.Errorf("level %d not found", levelID)
	}

	if level.IsUnlocked {
		return nil // Already unlocked
	}

	// Check if previous level is completed
	prevLevelID := levelID - 1
	if prevLevelID >= 100 {
		prevLevel, ok := cm.Levels[prevLevelID]
		if !ok || !prevLevel.IsCompleted {
			return fmt.Errorf("previous level not completed")
		}
	}

	level.IsUnlocked = true
	return nil
}

// CompleteLevel records level completion
func (cm *ChapterManager) CompleteLevel(levelID int, playerID string, score int64, timeSpent time.Duration, combo int, kills int, damage int64) (*LevelRewards, error) {
	level, ok := cm.Levels[levelID]
	if !ok {
		return nil, fmt.Errorf("level %d not found", levelID)
	}

	// Calculate stars
	stars := cm.calculateStars(level, score, timeSpent, combo)

	// Update level progress
	level.IsCompleted = true
	for i := 0; i < stars; i++ {
		level.Stars[i] = true
	}

	if score > level.BestScore {
		level.BestScore = score
	}

	if timeSpent < level.BestTime || level.BestTime == 0 {
		level.BestTime = timeSpent
	}

	// Update player progress
	progressKey := cm.getProgressKey(playerID, levelID)
	if _, exists := cm.PlayerProgress[progressKey]; !exists {
		cm.PlayerProgress[progressKey] = &PlayerLevelProgress{
			PlayerID:   playerID,
			LevelID:    levelID,
			FirstWinAt: time.Now(),
		}
	}

	progress := cm.PlayerProgress[progressKey]
	progress.Stars = max(progress.Stars, stars)
	progress.BestScore = max(progress.BestScore, score)
	if timeSpent < progress.BestTime || progress.BestTime == 0 {
		progress.BestTime = timeSpent
	}
	progress.PlayCount++
	progress.WinCount++
	progress.LastPlayAt = time.Now()
	progress.BestCombo = max(progress.BestCombo, combo)
	progress.TotalKills += kills
	progress.TotalDamage += damage

	// Calculate rewards
	rewards := cm.calculateRewards(level, stars)

	// Unlock next level
	nextLevelID := levelID + 1
	if nextLevel, ok := cm.Levels[nextLevelID]; ok {
		nextLevel.IsUnlocked = true
	}

	// Check chapter completion
	chapter := cm.GetLevel(levelID).ChapterID
	if chapter != 0 {
		cm.checkChapterCompletion(chapter)
	}

	return rewards, nil
}

// calculateStars calculates star rating based on performance
func (cm *ChapterManager) calculateStars(level *ChapterLevel, score int64, timeSpent time.Duration, combo int) int {
	stars := 0

	// Star 1: Complete the level
	if level.IsCompleted || score > 0 {
		stars = 1
	}

	// Star 2: Meet score target
	if score >= level.ScoreTarget {
		stars = 2
	}

	// Star 3: Complete under time limit with high combo
	if timeSpent < level.TimeLimit && combo >= 50 {
		stars = 3
	}

	return stars
}

// calculateRewards calculates rewards based on stars earned
func (cm *ChapterManager) calculateRewards(level *ChapterLevel, stars int) *LevelRewards {
	rewards := &LevelRewards{
		Coins: level.Rewards.Coins,
		Gems:  level.Rewards.Gems,
		EXP:   level.Rewards.EXP,
		Items: make([]RewardItem, 0),
	}

	// Add star-specific rewards
	for _, sr := range level.Rewards.StarRewards {
		if sr.Star <= stars {
			rewards.Coins += sr.Coins
			rewards.Gems += sr.Gems
			rewards.Items = append(rewards.Items, sr.Items...)
		}
	}

	// Add random item drops
	for _, item := range level.Rewards.Items {
		if cm.rand.Float64() < item.DropRate {
			rewards.Items = append(rewards.Items, item)
		}
	}

	return rewards
}

// checkChapterCompletion checks and updates chapter completion status
func (cm *ChapterManager) checkChapterCompletion(chapterID int) {
	chapter := cm.GetChapter(chapterID)
	if chapter == nil {
		return
	}

	completedLevels := 0
	totalStars := 0

	for _, level := range chapter.Levels {
		if level.IsCompleted {
			completedLevels++
			for _, star := range level.Stars {
				if star {
					totalStars++
				}
			}
		}
	}

	chapter.Stars = totalStars

	if completedLevels == len(chapter.Levels) {
		chapter.IsCompleted = true

		// Unlock next chapter
		nextChapterID := chapterID + 1
		if nextChapter := cm.GetChapter(nextChapterID); nextChapter != nil {
			nextChapter.IsUnlocked = true
		}
	}
}

// GetPlayerProgress returns player's progress for a specific level
func (cm *ChapterManager) GetPlayerProgress(playerID string, levelID int) *PlayerLevelProgress {
	progressKey := cm.getProgressKey(playerID, levelID)
	return cm.PlayerProgress[progressKey]
}

// GetChapterProgress returns player's overall progress for a chapter
func (cm *ChapterManager) GetChapterProgress(playerID string, chapterID int) (completedLevels, totalStars, maxStars int) {
	chapter := cm.GetChapter(chapterID)
	if chapter == nil {
		return 0, 0, 0
	}

	for _, level := range chapter.Levels {
		if level.IsCompleted {
			completedLevels++
		}
		for _, star := range level.Stars {
			if star {
				totalStars++
			}
		}
	}

	return completedLevels, totalStars, chapter.MaxStars
}

// GetAllChapters returns all chapters with player progress
func (cm *ChapterManager) GetAllChapters(playerID string) []*Chapter {
	result := make([]*Chapter, len(cm.Chapters))
	copy(result, cm.Chapters)

	// Enrich with player progress
	for i, chapter := range result {
		_, stars, maxStars := cm.GetChapterProgress(playerID, chapter.ID)
		chapter.Stars = stars
		chapter.MaxStars = maxStars
	}

	return result
}

// GetLevelsForChapter returns all levels in a chapter
func (cm *ChapterManager) GetLevelsForChapter(chapterID int) []*ChapterLevel {
	chapter := cm.GetChapter(chapterID)
	if chapter == nil {
		return nil
	}
	return chapter.Levels
}

// GenerateLevelEnemies generates enemy spawn data for a level session
func (cm *ChapterManager) GenerateLevelEnemies(levelID int, wave int) []Enemy {
	level, ok := cm.Levels[levelID]
	if !ok {
		return nil
	}

	var enemies []Enemy
	for _, spawn := range level.Enemies {
		if spawn.Wave == wave {
			for i := 0; i < spawn.Count; i++ {
				enemy := Enemy{
					Type:     spawn.EnemyType,
					HP:       cm.getEnemyHPByType(spawn.EnemyType),
					Position: cm.calculateSpawnPosition(spawn.Position, spawn.Formation, i, spawn.Count),
					Speed:    cm.getEnemySpeedByType(spawn.EnemyType),
				}
				enemies = append(enemies, enemy)
			}
		}
	}

	return enemies
}

func (cm *ChapterManager) getProgressKey(playerID string, levelID int) int {
	// Simple hash combining playerID and levelID
	hash := 0
	for _, c := range playerID {
		hash = hash*31 + int(c)
	}
	return hash*1000 + levelID
}

func (cm *ChapterManager) getEnemyHPByType(enemyType EnemyType) int64 {
	switch enemyType {
	case EnemyTypeBasic:
		return 100
	case EnemyTypeFast:
		return 50
	case EnemyTypeShooter:
		return 80
	case EnemyTypeTank:
		return 500
	case EnemyTypeBossMinion:
		return 200
	default:
		return 100
	}
}

func (cm *ChapterManager) getEnemySpeedByType(enemyType EnemyType) float64 {
	switch enemyType {
	case EnemyTypeFast:
		return 300
	case EnemyTypeShooter:
		return 50
	case EnemyTypeTank:
		return 30
	default:
		return 100
	}
}

func (cm *ChapterManager) calculateSpawnPosition(pos SpawnPosition, formation string, index, total int) Vector2 {
	var basePos Vector2

	switch pos {
	case SpawnTop:
		basePos = Vector2{X: float64(cm.rand.Intn(800)), Y: -50}
	case SpawnBottom:
		basePos = Vector2{X: float64(cm.rand.Intn(800)), Y: 850}
	case SpawnSides:
		if cm.rand.Float64() > 0.5 {
			basePos = Vector2{X: -50, Y: float64(cm.rand.Intn(600) + 100)}
		} else {
			basePos = Vector2{X: 850, Y: float64(cm.rand.Intn(600) + 100)}
		}
	default:
		basePos = Vector2{X: 400, Y: -50}
	}

	// Apply formation offset
	switch formation {
	case "line":
		basePos.X += float64(index * 60)
	case "circle":
		angle := 2 * math.Pi * float64(index) / float64(total)
		radius := 100.0
		basePos.X += math.Cos(angle) * radius
		basePos.Y += math.Sin(angle) * radius
	case "random":
		basePos.X += float64(cm.rand.Intn(200) - 100)
		basePos.Y += float64(cm.rand.Intn(200) - 100)
	}

	return basePos
}

// SpawnPosition defines where enemies spawn
type SpawnPosition int

const (
	SpawnTop SpawnPosition = iota
	SpawnBottom
	SpawnSides
	SpawnCenter
)

// EnemyType enum
type EnemyType int

const (
	EnemyTypeBasic EnemyType = iota
	EnemyTypeFast
	EnemyTypeShooter
	EnemyTypeTank
	EnemyTypeBossMinion
)
