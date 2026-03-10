// Package bot provides AI bot functionality for single-player mode
package bot

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"game-server/server/proto"
)

// BotConfig bot configuration
type BotConfig struct {
	BotID         string        // Unique bot ID
	BotName       string        // Bot display name
	BotType       BotType       // Bot difficulty level
	ResponseDelay time.Duration // Min delay before responding
	ThinkTime     time.Duration // Time to "think" before action
	Aggressiveness float64      // 0-1, how aggressive the bot plays
}

// BotType bot difficulty types
type BotType int

const (
	BotTypeEasy BotType = iota
	BotTypeMedium
	BotTypeHard
	BotTypeExpert
	BotTypeChampion
)

// String returns bot type name
func (b BotType) String() string {
	switch b {
	case BotTypeEasy:
		return "easy"
	case BotTypeMedium:
		return "medium"
	case BotTypeHard:
		return "hard"
	case BotTypeExpert:
		return "expert"
	case BotTypeChampion:
		return "champion"
	default:
		return "unknown"
	}
}

// BotAction represents an action a bot can take
type BotAction struct {
	Type    ActionType
	Target  string
	Payload interface{}
	Score   float64 // Action evaluation score
}

// ActionType types of bot actions
type ActionType int

const (
	ActionSpawnDinosaur ActionType = iota
	ActionUpgradeTower
	ActionUseSkill
	ActionMoveUnit
	ActionDefend
	ActionAttack
	ActionWait
)

// Bot represents an AI bot player
type Bot struct {
	config         BotConfig
	currentState   *GameState
	isRunning     atomic.Bool
	stopChan      chan struct{}
	actionChan    chan *BotAction
	responseTimer *time.Timer
	mu            sync.RWMutex

	// Statistics
	actionsTaken   int64
	wins           int64
	losses         int64
	totalGames     int64
	averageScore   float64
}

// GameState represents current game state for bot decision making
type GameState struct {
	RoomID         string
	MySide         int // 0=red(attack), 1=blue(defend)
	MyBaseHealth   float64
	EnemyBaseHealth float64
	MyTowers       []TowerInfo
	EnemyTowers    []TowerInfo
	MyUnits        []UnitInfo
	EnemyUnits     []UnitInfo
	Energy         int
	Credit         int
	EnemyBotType   BotType
	GameTime       time.Duration
}

// TowerInfo information about a tower
type TowerInfo struct {
	ID       string
	Type     string
	Level    int
	Health   float64
	Position [2]float64
}

// UnitInfo information about a unit
type UnitInfo struct {
	ID       string
	Type     string
	Health   float64
	Position [2]float64
	Target   string
}

// BotManager manages multiple bots
type BotManager struct {
	bots      map[string]*Bot
	configs   map[string]BotConfig
	mu        sync.RWMutex
	stats     *ManagerStats
	stopChan  chan struct{}
}

// ManagerStats statistics for bot manager
type ManagerStats struct {
	ActiveBots    int64
	TotalGames    int64
	TotalWins     int64
	TotalLosses   int64
	AvgWinRate    float64
}

// NewBot creates a new bot instance
func NewBot(config BotConfig) *Bot {
	return &Bot{
		config:     config,
		stopChan:   make(chan struct{}),
		actionChan: make(chan *BotAction, 10),
	}
}

// NewBotManager creates a new bot manager
func NewBotManager() *BotManager {
	return &BotManager{
		bots:    make(map[string]*Bot),
		configs: make(map[string]BotConfig),
		stats:   &ManagerStats{},
		stopChan: make(chan struct{}),
	}
}

// RegisterBot registers a bot with the manager
func (bm *BotManager) RegisterBot(bot *Bot) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.bots[bot.config.BotID] = bot
	bm.configs[bot.config.BotID] = bot.config
	atomic.AddInt64(&bm.stats.ActiveBots, 1)
}

// UnregisterBot removes a bot from the manager
func (bm *BotManager) UnregisterBot(botID string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if bot, ok := bm.bots[botID]; ok {
		bot.Stop()
		delete(bm.bots, botID)
		delete(bm.configs, botID)
		atomic.AddInt64(&bm.stats.ActiveBots, -1)
	}
}

// GetBot retrieves a bot by ID
func (bm *BotManager) GetBot(botID string) (*Bot, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	bot, ok := bm.bots[botID]
	return bot, ok
}

// ListBots lists all registered bots
func (bm *BotManager) ListBots() []*Bot {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	bots := make([]*Bot, 0, len(bm.bots))
	for _, bot := range bm.bots {
		bots = append(bots, bot)
	}
	return bots
}

// GetStats returns manager statistics
func (bm *BotManager) GetStats() ManagerStats {
	return ManagerStats{
		ActiveBots:  atomic.LoadInt64(&bm.stats.ActiveBots),
		TotalGames:  atomic.LoadInt64(&bm.stats.TotalGames),
		TotalWins:   atomic.LoadInt64(&bm.stats.TotalWins),
		TotalLosses: atomic.LoadInt64(&bm.stats.TotalLosses),
		AvgWinRate:  bm.calculateWinRate(),
	}
}

func (bm *BotManager) calculateWinRate() float64 {
	total := atomic.LoadInt64(&bm.stats.TotalGames)
	if total == 0 {
		return 0
	}
	wins := atomic.LoadInt64(&bm.stats.TotalWins)
	return float64(wins) / float64(total) * 100
}

// Start starts the bot's decision loop
func (b *Bot) Start(ctx context.Context) {
	if !b.isRunning.CompareAndSwap(false, true) {
		return
	}
	go b.decisionLoop(ctx)
}

// Stop stops the bot
func (b *Bot) Stop() {
	if !b.isRunning.CompareAndSwap(true, false) {
		return
	}
	close(b.stopChan)
}

func (b *Bot) decisionLoop(ctx context.Context) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-b.stopChan:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if b.currentState != nil {
				b.evaluateAndAct()
			}
		}
	}
}

func (b *Bot) evaluateAndAct() {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Add thinking delay based on bot type
	thinkTime := b.getThinkTime()
	time.Sleep(thinkTime)

	// Evaluate possible actions
	actions := b.evaluateActions()
	if len(actions) == 0 {
		return
	}

	// Select best action based on difficulty
	selectedAction := b.selectAction(actions)

	// Execute action
	if selectedAction != nil {
		select {
		case b.actionChan <- selectedAction:
			atomic.AddInt64(&b.actionsTaken, 1)
		default:
			// Channel full, skip action
		}
	}
}

func (b *Bot) evaluateActions() []*BotAction {
	var actions []*BotAction
	state := b.currentState

	// Evaluate spawn dinosaur actions
	if state.Energy >= 10 {
		for _, unitType := range b.getPreferredUnits() {
			cost := b.getUnitCost(unitType)
			if state.Energy >= cost {
				action := &BotAction{
					Type:   ActionSpawnDinosaur,
					Target: unitType,
					Score:  b.evaluateSpawnAction(unitType, cost),
				}
				actions = append(actions, action)
			}
		}
	}

	// Evaluate tower upgrade actions
	for _, tower := range state.MyTowers {
		if state.Credit >= 50 && tower.Level < 5 {
			action := &BotAction{
				Type:   ActionUpgradeTower,
				Target: tower.ID,
				Score:  b.evaluateUpgradeAction(tower),
			}
			actions = append(actions, action)
		}
	}

	// Evaluate skill usage
	if state.Energy >= 30 {
		action := &BotAction{
			Type:   ActionUseSkill,
			Target: b.selectSkillTarget(),
			Score:  b.evaluateSkillAction(),
		}
		actions = append(actions, action)
	}

	// Evaluate defense actions
	action := &BotAction{
		Type:  ActionDefend,
		Score: b.evaluateDefenseAction(),
	}
	actions = append(actions, action)

	return actions
}

func (b *Bot) evaluateSpawnAction(unitType string, cost int) float64 {
	state := b.currentState
	score := 0.0

	// Base score from unit type
	baseScore := map[string]float64{
		"velociraptor":  70,
		"triceratops":   80,
		"pterodactyl":   60,
		"tyrannosaurus": 100,
		"stegosaurus":   75,
		"ankylosaurus":  65,
	}[unitType]

	score += baseScore

	// Adjust based on game state
	if state.EnemyBaseHealth < 30 {
		score += 50 // Finish off!
	}

	if len(state.EnemyUnits) > len(state.MyUnits) {
		score += 20 // Need more units
	}

	// Adjust based on aggressiveness
	score *= (0.5 + b.config.Aggressiveness*0.5)

	// Random factor based on difficulty
	randomFactor := b.getRandomFactor()
	score += randomFactor * 20

	return score
}

func (b *Bot) evaluateUpgradeAction(tower TowerInfo) float64 {
	score := 30.0

	// Higher score for low health towers
	if tower.Health < 50 {
		score += 30
	}

	// Higher score for towers under attack
	enemyCount := 0
	for _, unit := range b.currentState.EnemyUnits {
		dist := distance(unit.Position, tower.Position)
		if dist < 100 {
			enemyCount++
		}
	}
	score += float64(enemyCount * 10)

	return score
}

func (b *Bot) evaluateSkillAction() float64 {
	state := b.currentState
	score := 40.0

	// Higher score if many enemy units
	score += float64(len(state.EnemyUnits) * 5)

	// Higher score if enemy base is weak
	if state.EnemyBaseHealth < 50 {
		score += 30
	}

	return score
}

func (b *Bot) evaluateDefenseAction() float64 {
	state := b.currentState
	score := 20.0

	// Higher score if base is under threat
	if state.MyBaseHealth < 50 {
		score += 40
	}

	// Adjust based on difficulty
	score *= (1.0 - b.config.Aggressiveness*0.3)

	return score
}

func (b *Bot) selectAction(actions []*BotAction) *BotAction {
	if len(actions) == 0 {
		return nil
	}

	// Sort by score descending
	for i := 0; i < len(actions)-1; i++ {
		for j := i + 1; j < len(actions); j++ {
			if actions[j].Score > actions[i].Score {
				actions[i], actions[j] = actions[j], actions[i]
			}
		}
	}

	// Select based on difficulty
	switch b.config.BotType {
	case BotTypeEasy:
		// 60% best, 40% random
		if rand.Float64() < 0.6 {
			return actions[0]
		}
		return actions[rand.Intn(len(actions))]
	case BotTypeMedium:
		// 75% best, 25% random
		if rand.Float64() < 0.75 {
			return actions[0]
		}
		return actions[rand.Intn(len(actions))]
	case BotTypeHard:
		// 90% best, 10% random
		if rand.Float64() < 0.9 {
			return actions[0]
		}
		return actions[rand.Intn(len(actions))]
	case BotTypeExpert, BotTypeChampion:
		// Always best
		return actions[0]
	default:
		return actions[0]
	}
}

func (b *Bot) getThinkTime() time.Duration {
	base := b.config.ThinkTime
	variation := time.Duration(rand.Int63n(int64(base)))
	return base + variation
}

func (b *Bot) getRandomFactor() float64 {
	// Lower difficulty = more randomness
	factors := map[BotType]float64{
		BotTypeEasy:     0.8,
		BotTypeMedium:   0.5,
		BotTypeHard:     0.3,
		BotTypeExpert:   0.1,
		BotTypeChampion: 0.05,
	}
	return factors[b.config.BotType] * rand.Float64()
}

func (b *Bot) getPreferredUnits() []string {
	preferences := map[BotType][]string{
		BotTypeEasy:     {"velociraptor", "stegosaurus"},
		BotTypeMedium:   {"velociraptor", "triceratops", "pterodactyl"},
		BotTypeHard:     {"triceratops", "pterodactyl", "tyrannosaurus"},
		BotTypeExpert:   {"pterodactyl", "tyrannosaurus"},
		BotTypeChampion: {"tyrannosaurus", "pterodactyl", "triceratops"},
	}
	return preferences[b.config.BotType]
}

func (b *Bot) getUnitCost(unitType string) int {
	costs := map[string]int{
		"velociraptor":  10,
		"triceratops":   15,
		"pterodactyl":  20,
		"tyrannosaurus": 30,
		"stegosaurus":   12,
		"ankylosaurus": 18,
	}
	return costs[unitType]
}

func (b *Bot) selectSkillTarget() string {
	state := b.currentState

	// Target enemy base if low health
	if state.EnemyBaseHealth < 30 {
		return "enemy_base"
	}

	// Target strongest enemy unit
	var strongestUnit string
	maxHealth := 0.0
	for _, unit := range state.EnemyUnits {
		if unit.Health > maxHealth {
			maxHealth = unit.Health
			strongestUnit = unit.ID
		}
	}

	if strongestUnit != "" {
		return strongestUnit
	}

	return "enemy_base"
}

// UpdateState updates the bot's knowledge of game state
func (b *Bot) UpdateState(state *GameState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.currentState = state
}

// GetActionChan returns the channel for bot actions
func (b *Bot) GetActionChan() <-chan *BotAction {
	return b.actionChan
}

// GetStats returns bot statistics
func (b *Bot) GetStats() BotStats {
	return BotStats{
		BotID:         b.config.BotID,
		BotName:       b.config.BotName,
		BotType:       b.config.BotType.String(),
		ActionsTaken:  atomic.LoadInt64(&b.actionsTaken),
		Wins:          atomic.LoadInt64(&b.wins),
		Losses:        atomic.LoadInt64(&b.losses),
		TotalGames:    atomic.LoadInt64(&b.totalGames),
		AverageScore:  b.averageScore,
		IsRunning:     b.isRunning.Load(),
	}
}

// RecordWin records a win for the bot
func (b *Bot) RecordWin() {
	atomic.AddInt64(&b.wins, 1)
	atomic.AddInt64(&b.totalGames, 1)
	b.updateAverageScore(1.0)
}

// RecordLoss records a loss for the bot
func (b *Bot) RecordLoss() {
	atomic.AddInt64(&b.losses, 1)
	atomic.AddInt64(&b.totalGames, 1)
	b.updateAverageScore(0.0)
}

func (b *Bot) updateAverageScore(score float64) {
	oldAvg := b.averageScore
	total := float64(atomic.LoadInt64(&b.totalGames))
	b.averageScore = oldAvg + (score-oldAvg)/total
}

// BotStats represents bot statistics
type BotStats struct {
	BotID        string
	BotName      string
	BotType      string
	ActionsTaken int64
	Wins         int64
	Losses       int64
	TotalGames   int64
	AverageScore float64
	IsRunning    bool
}

// CreateBotFromType creates a bot with predefined difficulty settings
func CreateBotFromType(botID, botName string, botType BotType) *Bot {
	configs := map[BotType]BotConfig{
		BotTypeEasy: {
			BotID:          botID,
			BotName:        botName,
			BotType:        BotTypeEasy,
			ResponseDelay:  500 * time.Millisecond,
			ThinkTime:      2 * time.Second,
			Aggressiveness: 0.3,
		},
		BotTypeMedium: {
			BotID:          botID,
			BotName:        botName,
			BotType:        BotTypeMedium,
			ResponseDelay:  300 * time.Millisecond,
			ThinkTime:      1 * time.Second,
			Aggressiveness: 0.5,
		},
		BotTypeHard: {
			BotID:          botID,
			BotName:        botName,
			BotType:        BotTypeHard,
			ResponseDelay:  200 * time.Millisecond,
			ThinkTime:      500 * time.Millisecond,
			Aggressiveness: 0.7,
		},
		BotTypeExpert: {
			BotID:          botID,
			BotName:        botName,
			BotType:        BotTypeExpert,
			ResponseDelay: 100 * time.Millisecond,
			ThinkTime:      200 * time.Millisecond,
			Aggressiveness: 0.85,
		},
		BotTypeChampion: {
			BotID:          botID,
			BotName:        botName,
			BotType:        BotTypeChampion,
			ResponseDelay:  50 * time.Millisecond,
			ThinkTime:      100 * time.Millisecond,
			Aggressiveness: 1.0,
		},
	}

	return NewBot(configs[botType])
}

// CreateRandomBot creates a bot with random difficulty
func CreateRandomBot(botID, botName string) *Bot {
	botTypes := []BotType{BotTypeEasy, BotTypeMedium, BotTypeHard, BotTypeExpert, BotTypeChampion}
	weights := []float64{0.3, 0.3, 0.2, 0.15, 0.05}

	r := rand.Float64()
	cumulative := 0.0
	var selectedType BotType
	for i, w := range weights {
		cumulative += w
		if r <= cumulative {
			selectedType = botTypes[i]
			break
		}
	}

	return CreateBotFromType(botID, botName, selectedType)
}

// distance calculates distance between two points
func distance(p1, p2 [2]float64) float64 {
	dx := p2[0] - p1[0]
	dy := p2[1] - p1[1]
	return math.Sqrt(dx*dx + dy*dy)
}

// BotMatchmaker handles matching players with bots
type BotMatchmaker struct {
	manager     *BotManager
	playerQueue chan *PlayerRequest
	stopChan    chan struct{}
}

// PlayerRequest represents a player's match request
type PlayerRequest struct {
	PlayerID    string
	PlayerName  string
	Difficulty BotType
	ResultChan chan *MatchResult
}

// MatchResult represents match result
type MatchResult struct {
	RoomID   string
	Bot      *Bot
	PlayerID string
	Success  bool
	Error    error
}

// NewBotMatchmaker creates a new bot matchmaker
func NewBotMatchmaker(manager *BotManager) *BotMatchmaker {
	return &BotMatchmaker{
		manager:     manager,
		playerQueue: make(chan *PlayerRequest, 100),
		stopChan:    make(chan struct{}),
	}
}

// Start starts the matchmaker
func (bm *BotMatchmaker) Start(ctx context.Context) {
	go bm.matchLoop(ctx)
}

func (bm *BotMatchmaker) matchLoop(ctx context.Context) {
	for {
		select {
		case <-bm.stopChan:
			return
		case <-ctx.Done():
			return
		case request := <-bm.playerQueue:
			bm.processMatch(request)
		}
	}
}

func (bm *BotMatchmaker) processMatch(request *PlayerRequest) {
	// Create a new bot for this match
	botID := fmt.Sprintf("bot_%s_%d", request.PlayerID, time.Now().UnixNano())
	bot := CreateBotFromType(botID, "AI对手", request.Difficulty)

	// Register with manager
	bm.manager.RegisterBot(bot)

	// Create room (in real implementation, this would call room manager)
	roomID := fmt.Sprintf("room_%d", rand.Int63n(1000000))

	// Return result
	result := &MatchResult{
		RoomID:   roomID,
		Bot:      bot,
		PlayerID: request.PlayerID,
		Success:  true,
	}

	select {
	case request.ResultChan <- result:
	case <-time.After(5 * time.Second):
		// Timeout
	}
}

// RequestMatch requests a match against a bot
func (bm *BotMatchmaker) RequestMatch(playerID, playerName string, difficulty BotType) (*MatchResult, error) {
	request := &PlayerRequest{
		PlayerID:    playerID,
		PlayerName:  playerName,
		Difficulty:  difficulty,
		ResultChan:  make(chan *MatchResult, 1),
	}

	select {
	case bm.playerQueue <- request:
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("matchmaking timeout")
	}

	select {
	case result := <-request.ResultChan:
		return result, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("matchmaking timeout")
	}
}

// Stop stops the matchmaker
func (bm *BotMatchmaker) Stop() {
	close(bm.stopChan)
}

// ConvertProtoToGameState converts protobuf state to bot game state
func ConvertProtoToGameState(protoState *proto.GameState) *GameState {
	state := &GameState{
		RoomID:          protoState.RoomId,
		MySide:          int(protoState.CurrentSide),
		MyBaseHealth:    protoState.RedBase.Health,
		EnemyBaseHealth: protoState.BlueBase.Health,
		Energy:          int(protoState.RedEnergy),
		Credit:          int(protoState.RedCredit),
	}

	// Convert towers
	for _, tower := range protoState.RedTowers {
		state.MyTowers = append(state.MyTowers, TowerInfo{
			ID:       tower.Id,
			Type:     tower.Type,
			Level:    int(tower.Level),
			Health:   tower.Health,
			Position: [2]float64{tower.X, tower.Y},
		})
	}

	for _, tower := range protoState.BlueTowers {
		state.EnemyTowers = append(state.EnemyTowers, TowerInfo{
			ID:       tower.Id,
			Type:     tower.Type,
			Level:    int(tower.Level),
			Health:   tower.Health,
			Position: [2]float64{tower.X, tower.Y},
		})
	}

	// Convert units
	for _, unit := range protoState.RedUnits {
		state.MyUnits = append(state.MyUnits, UnitInfo{
			ID:       unit.Id,
			Type:     unit.Type,
			Health:   unit.Health,
			Position: [2]float64{unit.X, unit.Y},
			Target:   unit.TargetId,
		})
	}

	for _, unit := range protoState.BlueUnits {
		state.EnemyUnits = append(state.EnemyUnits, UnitInfo{
			ID:       unit.Id,
			Type:     unit.Type,
			Health:   unit.Health,
			Position: [2]float64{unit.X, unit.Y},
			Target:   unit.TargetId,
		})
	}

	return state
}
