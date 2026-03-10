// Package game provides the core game functionality for the bullet hell game
package game

import (
	"math"
	"sync"
	"time"
)

// ==================== Sound System ====================

// SoundType represents the type of sound effect
type SoundType int

const (
	SoundTypeBgm SoundType = iota // Background music
	SoundTypeSfx                  // Sound effects
	SoundTypeVoice                // Voice lines
	SoundTypeAmbient              // Ambient sounds
)

// SoundCategory represents the category of sound
type SoundCategory int

const (
	SoundCategoryUI SoundCategory = iota // UI sounds
	SoundCategoryBattle                  // Battle sounds
	SoundCategorySkill                   // Skill sounds
	SoundCategoryItem                    // Item sounds
	SoundCategoryEnemy                   // Enemy sounds
	SoundCategoryEnvironment             // Environment sounds
	SoundCategorySystem                  // System sounds
)

// SoundEffect represents a specific sound effect
type SoundEffect struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Category    SoundCategory `json:"category"`
	Type        SoundType   `json:"type"`
	FilePath    string      `json:"file_path"`
	Duration    int         `json:"duration"` // in milliseconds
	Volume      float32     `json:"volume"`    // 0.0 - 1.0
	Loop        bool        `json:"loop"`
	Priority    int         `json:"priority"` // 0-100, higher plays first
	Cooldown    int         `json:"cooldown"`  // in milliseconds
	LastPlay    int64       `json:"last_play"` // timestamp
}

// SoundPlayer manages sound playback
type SoundPlayer struct {
	mu           sync.RWMutex
	sounds       map[string]*SoundEffect
	playing      map[string]bool
	volumeMaster float32 // Master volume 0.0 - 1.0
	volumeBgm    float32 // BGM volume
	volumeSfx    float32 // SFX volume
	enabled      bool    // Sound enabled
	mute         bool    // Mute all sounds
}

// NewSoundPlayer creates a new sound player
func NewSoundPlayer() *SoundPlayer {
	return &SoundPlayer{
		sounds:       make(map[string]*SoundEffect),
		playing:      make(map[string]bool),
		volumeMaster: 1.0,
		volumeBgm:    0.8,
		volumeSfx:    1.0,
		enabled:      true,
		mute:         false,
	}
}

// RegisterSound registers a new sound effect
func (sp *SoundPlayer) RegisterSound(sound *SoundEffect) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.sounds[sound.ID] = sound
}

// GetSound retrieves a sound by ID
func (sp *SoundPlayer) GetSound(id string) (*SoundEffect, bool) {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	sound, ok := sp.sounds[id]
	return sound, ok
}

// Play plays a sound effect by ID
func (sp *SoundPlayer) Play(id string) bool {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if !sp.enabled || sp.mute {
		return false
	}

	sound, ok := sp.sounds[id]
	if !ok {
		return false
	}

	// Check cooldown
	now := time.Now().UnixMilli()
	if now-sound.LastPlay < int64(sound.Cooldown) {
		return false
	}

	// Check if already playing (for non-looping sounds)
	if !sound.Loop && sp.playing[id] {
		return false
	}

	sound.LastPlay = now
	sp.playing[id] = true
	return true
}

// Stop stops a playing sound
func (sp *SoundPlayer) Stop(id string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.playing[id] = false
}

// StopAll stops all playing sounds
func (sp *SoundPlayer) StopAll() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	for id := range sp.playing {
		sp.playing[id] = false
	}
}

// IsPlaying checks if a sound is playing
func (sp *SoundPlayer) IsPlaying(id string) bool {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	return sp.playing[id]
}

// SetMasterVolume sets the master volume
func (sp *SoundPlayer) SetMasterVolume(volume float32) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.volumeMaster = clampFloat32(volume, 0, 1)
}

// SetBgmVolume sets the BGM volume
func (sp *SoundPlayer) SetBgmVolume(volume float32) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.volumeBgm = clampFloat32(volume, 0, 1)
}

// SetSfxVolume sets the SFX volume
func (sp *SoundPlayer) SetSfxVolume(volume float32) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.volumeSfx = clampFloat32(volume, 0, 1)
}

// GetEffectiveVolume gets the effective volume for a sound
func (sp *SoundPlayer) GetEffectiveVolume(soundType SoundType) float32 {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	if sp.mute || !sp.enabled {
		return 0
	}

	baseVolume := sp.volumeMaster
	if soundType == SoundTypeBgm {
		baseVolume *= sp.volumeBgm
	} else {
		baseVolume *= sp.volumeSfx
	}

	return baseVolume
}

// Enable enables sound
func (sp *SoundPlayer) Enable() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.enabled = true
}

// Disable disables sound
func (sp *SoundPlayer) Disable() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.enabled = false
}

// Mute mutes all sounds
func (sp *SoundPlayer) Mute() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.mute = true
}

// Unmute unmutes sounds
func (sp *SoundPlayer) Unmute() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.mute = false
}

// GetSoundsByCategory gets all sounds in a category
func (sp *SoundPlayer) GetSoundsByCategory(category SoundCategory) []*SoundEffect {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	result := make([]*SoundEffect, 0)
	for _, sound := range sp.sounds {
		if sound.Category == category {
			result = append(result, sound)
		}
	}
	return result
}

// ==================== Sound Presets ====================

// SoundPresets returns default sound effects
func SoundPresets() map[string]*SoundEffect {
	sounds := make(map[string]*SoundEffect)

	// UI Sounds
	sounds["ui_click"] = &SoundEffect{
		ID:       "ui_click",
		Name:     "UI Click",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Duration: 100,
		Volume:   0.5,
		Priority: 50,
		Cooldown: 50,
	}

	sounds["ui_hover"] = &SoundEffect{
		ID:       "ui_hover",
		Name:     "UI Hover",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Duration: 80,
		Volume:   0.3,
		Priority: 30,
		Cooldown: 100,
	}

	sounds["ui_confirm"] = &SoundEffect{
		ID:       "ui_confirm",
		Name:     "UI Confirm",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Duration: 200,
		Volume:   0.6,
		Priority: 70,
		Cooldown: 200,
	}

	sounds["ui_cancel"] = &SoundEffect{
		ID:       "ui_cancel",
		Name:     "UI Cancel",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Duration: 150,
		Volume:   0.5,
		Priority: 60,
		Cooldown: 150,
	}

	sounds["ui_error"] = &SoundEffect{
		ID:       "ui_error",
		Name:     "UI Error",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Duration: 300,
		Volume:   0.6,
		Priority: 80,
		Cooldown: 300,
	}

	sounds["ui_success"] = &SoundEffect{
		ID:       "ui_success",
		Name:     "UI Success",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Duration: 250,
		Volume:   0.6,
		Priority: 70,
		Cooldown: 250,
	}

	// Battle Sounds
	sounds["battle_start"] = &SoundEffect{
		ID:       "battle_start",
		Name:     "Battle Start",
		Category: SoundCategoryBattle,
		Type:     SoundTypeSfx,
		Duration: 1000,
		Volume:   0.8,
		Priority: 90,
		Cooldown: 0,
	}

	sounds["battle_win"] = &SoundEffect{
		ID:       "battle_win",
		Name:     "Battle Win",
		Category: SoundCategoryBattle,
		Type:     SoundTypeSfx,
		Duration: 2000,
		Volume:   0.9,
		Priority: 95,
		Cooldown: 0,
	}

	sounds["battle_lose"] = &SoundEffect{
		ID:       "battle_lose",
		Name:     "Battle Lose",
		Category: SoundCategoryBattle,
		Type:     SoundTypeSfx,
		Duration: 1500,
		Volume:   0.7,
		Priority: 90,
		Cooldown: 0,
	}

	sounds["battle_countdown"] = &SoundEffect{
		ID:       "battle_countdown",
		Name:     "Battle Countdown",
		Category: SoundCategoryBattle,
		Type:     SoundTypeSfx,
		Duration: 500,
		Volume:   0.7,
		Priority: 85,
		Cooldown: 500,
	}

	sounds["battle_overtime"] = &SoundEffect{
		ID:       "battle_overtime",
		Name:     "Battle Overtime",
		Category: SoundCategoryBattle,
		Type:     SoundTypeSfx,
		Duration: 800,
		Volume:   0.8,
		Priority: 90,
		Cooldown: 0,
	}

	// Skill Sounds
	sounds["skill_charge"] = &SoundEffect{
		ID:       "skill_charge",
		Name:     "Skill Charge",
		Category: SoundCategorySkill,
		Type:     SoundTypeSfx,
		Duration: 500,
		Volume:   0.6,
		Priority: 60,
		Cooldown: 0,
	}

	sounds["skill_release"] = &SoundEffect{
		ID:       "skill_release",
		Name:     "Skill Release",
		Category: SoundCategorySkill,
		Type:     SoundTypeSfx,
		Duration: 800,
		Volume:   0.8,
		Priority: 70,
		Cooldown: 0,
	}

	sounds["skill_cooldown"] = &SoundEffect{
		ID:       "skill_cooldown",
		Name:     "Skill Cooldown",
		Category: SoundCategorySkill,
		Type:     SoundTypeSfx,
		Duration: 300,
		Volume:   0.5,
		Priority: 50,
		Cooldown: 0,
	}

	sounds["skill_upgrade"] = &SoundEffect{
		ID:       "skill_upgrade",
		Name:     "Skill Upgrade",
		Category: SoundCategorySkill,
		Type:     SoundTypeSfx,
		Duration: 600,
		Volume:   0.7,
		Priority: 75,
		Cooldown: 0,
	}

	// Item Sounds
	sounds["item_pickup"] = &SoundEffect{
		ID:       "item_pickup",
		Name:     "Item Pickup",
		Category: SoundCategoryItem,
		Type:     SoundTypeSfx,
		Duration: 300,
		Volume:   0.6,
		Priority: 60,
		Cooldown: 100,
	}

	sounds["item_use"] = &SoundEffect{
		ID:       "item_use",
		Name:     "Item Use",
		Category: SoundCategoryItem,
		Type:     SoundTypeSfx,
		Duration: 400,
		Volume:   0.7,
		Priority: 65,
		Cooldown: 200,
	}

	sounds["item_equip"] = &SoundEffect{
		ID:       "item_equip",
		Name:     "Item Equip",
		Category: SoundCategoryItem,
		Type:     SoundTypeSfx,
		Duration: 350,
		Volume:   0.6,
		Priority: 55,
		Cooldown: 150,
	}

	sounds["item_unequip"] = &SoundEffect{
		ID:       "item_unequip",
		Name:     "Item Unequip",
		Category: SoundCategoryItem,
		Type:     SoundTypeSfx,
		Duration: 250,
		Volume:   0.5,
		Priority: 50,
		Cooldown: 100,
	}

	// Enemy Sounds
	sounds["enemy_spawn"] = &SoundEffect{
		ID:       "enemy_spawn",
		Name:     "Enemy Spawn",
		Category: SoundCategoryEnemy,
		Type:     SoundTypeSfx,
		Duration: 400,
		Volume:   0.6,
		Priority: 55,
		Cooldown: 200,
	}

	sounds["enemy_death"] = &SoundEffect{
		ID:       "enemy_death",
		Name:     "Enemy Death",
		Category: SoundCategoryEnemy,
		Type:     SoundTypeSfx,
		Duration: 500,
		Volume:   0.7,
		Priority: 70,
		Cooldown: 100,
	}

	sounds["enemy_hit"] = &SoundEffect{
		ID:       "enemy_hit",
		Name:     "Enemy Hit",
		Category: SoundCategoryEnemy,
		Type:     SoundTypeSfx,
		Duration: 150,
		Volume:   0.5,
		Priority: 40,
		Cooldown: 50,
	}

	sounds["enemy_attack"] = &SoundEffect{
		ID:       "enemy_attack",
		Name:     "Enemy Attack",
		Category: SoundCategoryEnemy,
		Type:     SoundTypeSfx,
		Duration: 300,
		Volume:   0.6,
		Priority: 50,
		Cooldown: 200,
	}

	sounds["boss_appear"] = &SoundEffect{
		ID:       "boss_appear",
		Name:     "Boss Appear",
		Category: SoundCategoryEnemy,
		Type:     SoundTypeSfx,
		Duration: 1500,
		Volume:   0.9,
		Priority: 95,
		Cooldown: 0,
	}

	sounds["boss_death"] = &SoundEffect{
		ID:       "boss_death",
		Name:     "Boss Death",
		Category: SoundCategoryEnemy,
		Type:     SoundTypeSfx,
		Duration: 2500,
		Volume:   1.0,
		Priority: 100,
		Cooldown: 0,
	}

	// Environment Sounds
	sounds["env_rain"] = &SoundEffect{
		ID:       "env_rain",
		Name:     "Rain",
		Category: SoundCategoryEnvironment,
		Type:     SoundTypeAmbient,
		Duration: 5000,
		Volume:   0.4,
		Priority: 20,
		Cooldown: 0,
		Loop:     true,
	}

	sounds["env_thunder"] = &SoundEffect{
		ID:       "env_thunder",
		Name:     "Thunder",
		Category: SoundCategoryEnvironment,
		Type:     SoundTypeAmbient,
		Duration: 2000,
		Volume:   0.7,
		Priority: 50,
		Cooldown: 5000,
	}

	sounds["env_wind"] = &SoundEffect{
		ID:       "env_wind",
		Name:     "Wind",
		Category: SoundCategoryEnvironment,
		Type:     SoundTypeAmbient,
		Duration: 4000,
		Volume:   0.3,
		Priority: 15,
		Cooldown: 0,
		Loop:     true,
	}

	// System Sounds
	sounds["sys_notification"] = &SoundEffect{
		ID:       "sys_notification",
		Name:     "System Notification",
		Category: SoundCategorySystem,
		Type:     SoundTypeSfx,
		Duration: 500,
		Volume:   0.5,
		Priority: 40,
		Cooldown: 1000,
	}

	sounds["sys_achievement"] = &SoundEffect{
		ID:       "sys_achievement",
		Name:     "Achievement Unlocked",
		Category: SoundCategorySystem,
		Type:     SoundTypeSfx,
		Duration: 1200,
		Volume:   0.8,
		Priority: 85,
		Cooldown: 0,
	}

	sounds["sys_levelup"] = &SoundEffect{
		ID:       "sys_levelup",
		Name:     "Level Up",
		Category: SoundCategorySystem,
		Type:     SoundTypeSfx,
		Duration: 1500,
		Volume:   0.8,
		Priority: 90,
		Cooldown: 0,
	}

	sounds["sys_save"] = &SoundEffect{
		ID:       "sys_save",
		Name:     "Game Saved",
		Category: SoundCategorySystem,
		Type:     SoundTypeSfx,
		Duration: 400,
		Volume:   0.5,
		Priority: 45,
		Cooldown: 500,
	}

	sounds["sys_load"] = &SoundEffect{
		ID:       "sys_load",
		Name:     "Game Loaded",
		Category: SoundCategorySystem,
		Type:     SoundTypeSfx,
		Duration: 600,
		Volume:   0.6,
		Priority: 50,
		Cooldown: 500,
	}

	return sounds
}

// ==================== Helper Functions ====================

func clampFloat32(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampFloat64(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func absFloat64(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func roundFloat64(v float64) float64 {
	return math.Round(v)
}

func floorFloat64(v float64) float64 {
	return math.Floor(v)
}

func ceilFloat64(v float64) float64 {
	return math.Ceil(v)
}

func sqrtFloat64(v float64) float64 {
	return math.Sqrt(v)
}

func powFloat64(v, exp float64) float64 {
	return math.Pow(v, exp)
}

func logFloat64(v float64) float64 {
	return math.Log(v)
}

func log10Float64(v float64) float64 {
	return math.Log10(v)
}

func expFloat64(v float64) float64 {
	return math.Exp(v)
}

func sinFloat64(v float64) float64 {
	return math.Sin(v)
}

func cosFloat64(v float64) float64 {
	return math.Cos(v)
}

func tanFloat64(v float64) float64 {
	return math.Tan(v)
}

func asinFloat64(v float64) float64 {
	return math.Asin(v)
}

func acosFloat64(v float64) float64 {
	return math.Acos(v)
}

func atanFloat64(v float64) float64 {
	return math.Atan(v)
}

func atan2Float64(y, x float64) float64 {
	return math.Atan2(y, x)
}

func degToRad(d float64) float64 {
	return d * math.Pi / 180
}

func radToDeg(r float64) float64 {
	return r * 180 / math.Pi
}

// Lerp linearly interpolates between a and b
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// InverseLerp returns the interpolation factor between a and b
func InverseLerp(a, b, v float64) float64 {
	if a == b {
		return 0
	}
	return (v - a) / (b - a)
}

// Remap remaps a value from one range to another
func Remap(v, inMin, inMax, outMin, outMax float64) float64 {
	t := InverseLerp(inMin, inMax, v)
	return Lerp(outMin, outMax, t)
}

// SmoothStep smoothstep interpolation
func SmoothStep(edge0, edge1, x float64) float64 {
	t := clampFloat64((x-edge0)/(edge1-edge0), 0, 1)
	return t * t * (3 - 2*t)
}

// SmootherStep smootherstep interpolation
func SmootherStep(edge0, edge1, x float64) float64 {
	t := clampFloat64((x-edge0)/(edge1-edge0), 0, 1)
	return t * t * t * (t*(t*6-15) + 10)
}

// ==================== Time Utilities ====================

// Duration represents a time duration with useful methods
type Duration struct {
	milliseconds int64
}

func NewDuration(ms int64) *Duration {
	return &Duration{milliseconds: ms}
}

func (d *Duration) Milliseconds() int64 {
	return d.milliseconds
}

func (d *Duration) Seconds() float64 {
	return float64(d.milliseconds) / 1000
}

func (d *Duration) Minutes() float64 {
	return float64(d.milliseconds) / 60000
}

func (d *Duration) Add(other *Duration) *Duration {
	return &Duration{milliseconds: d.milliseconds + other.milliseconds}
}

func (d *Duration) Sub(other *Duration) *Duration {
	return &Duration{milliseconds: d.milliseconds - other.milliseconds}
}

func (d *Duration) Multiply(factor float64) *Duration {
	return &Duration{milliseconds: int64(float64(d.milliseconds) * factor)}
}

// ==================== Timer ====================

// Timer represents a game timer
type Timer struct {
	Duration  int64   `json:"duration"`  // Total duration in milliseconds
	Elapsed   int64   `json:"elapsed"`   // Elapsed time in milliseconds
	Running   bool    `json:"running"`   // Is timer running
	Loop      bool    `json:"loop"`      // Should timer loop
	OnComplete func() `json:"-"`         // Callback on complete
	mu        sync.RWMutex
}

// NewTimer creates a new timer
func NewTimer(duration int64, loop bool) *Timer {
	return &Timer{
		Duration: duration,
		Elapsed:  0,
		Running:  false,
		Loop:     loop,
	}
}

// Start starts the timer
func (t *Timer) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Running = true
}

// Stop stops the timer
func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Running = false
}

// Reset resets the timer
func (t *Timer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Elapsed = 0
	t.Running = false
}

// Update updates the timer by delta time
func (t *Timer) Update(deltaMs int64) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Running {
		return false
	}

	t.Elapsed += deltaMs

	if t.Elapsed >= t.Duration {
		if t.Loop {
			t.Elapsed = t.Elapsed % t.Duration
			if t.OnComplete != nil {
				t.OnComplete()
			}
			return true
		}
		t.Running = false
		if t.OnComplete != nil {
			t.OnComplete()
		}
		return true
	}

	return false
}

// GetProgress returns the progress (0.0 - 1.0)
func (t *Timer) GetProgress() float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.Duration == 0 {
		return 0
	}
	return float64(t.Elapsed) / float64(t.Duration)
}

// Remaining returns remaining time in milliseconds
func (t *Timer) Remaining() int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	remaining := t.Duration - t.Elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsComplete checks if timer is complete
func (t *Timer) IsComplete() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return !t.Running && t.Elapsed >= t.Duration
}
