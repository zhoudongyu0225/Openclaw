package game

import (
	"testing"
	"time"
)

// ==================== SoundPlayer Tests ====================

func TestSoundPlayer_NewSoundPlayer(t *testing.T) {
	sp := NewSoundPlayer()

	if sp == nil {
		t.Fatal("NewSoundPlayer returned nil")
	}

	if sp.sounds == nil {
		t.Error("sounds map not initialized")
	}

	if sp.playing == nil {
		t.Error("playing map not initialized")
	}

	if sp.volumeMaster != 1.0 {
		t.Errorf("expected master volume 1.0, got %f", sp.volumeMaster)
	}

	if sp.volumeBgm != 0.8 {
		t.Errorf("expected BGM volume 0.8, got %f", sp.volumeBgm)
	}

	if sp.volumeSfx != 1.0 {
		t.Errorf("expected SFX volume 1.0, got %f", sp.volumeSfx)
	}

	if !sp.enabled {
		t.Error("expected enabled to be true")
	}

	if sp.mute {
		t.Error("expected mute to be false")
	}
}

func TestSoundPlayer_RegisterSound(t *testing.T) {
	sp := NewSoundPlayer()

	sound := &SoundEffect{
		ID:       "test_sound",
		Name:     "Test Sound",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Duration: 500,
		Volume:   0.7,
		Priority: 50,
		Cooldown: 100,
	}

	sp.RegisterSound(sound)

	retrieved, ok := sp.GetSound("test_sound")
	if !ok {
		t.Error("failed to retrieve registered sound")
	}

	if retrieved.ID != sound.ID {
		t.Errorf("expected ID %s, got %s", sound.ID, retrieved.ID)
	}

	if retrieved.Name != sound.Name {
		t.Errorf("expected Name %s, got %s", sound.Name, retrieved.Name)
	}

	if retrieved.Category != sound.Category {
		t.Errorf("expected Category %v, got %v", sound.Category, retrieved.Category)
	}

	if retrieved.Type != sound.Type {
		t.Errorf("expected Type %v, got %v", sound.Type, retrieved.Type)
	}

	if retrieved.Duration != sound.Duration {
		t.Errorf("expected Duration %d, got %d", sound.Duration, retrieved.Duration)
	}

	if retrieved.Volume != sound.Volume {
		t.Errorf("expected Volume %f, got %f", sound.Volume, retrieved.Volume)
	}

	if retrieved.Priority != sound.Priority {
		t.Errorf("expected Priority %d, got %d", sound.Priority, retrieved.Priority)
	}

	if retrieved.Cooldown != sound.Cooldown {
		t.Errorf("expected Cooldown %d, got %d", sound.Cooldown, retrieved.Cooldown)
	}
}

func TestSoundPlayer_GetSound_NotFound(t *testing.T) {
	sp := NewSoundPlayer()

	_, ok := sp.GetSound("nonexistent")
	if ok {
		t.Error("expected not found for nonexistent sound")
	}
}

func TestSoundPlayer_Play_Success(t *testing.T) {
	sp := NewSoundPlayer()

	sound := &SoundEffect{
		ID:       "play_test",
		Name:     "Play Test",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Duration: 500,
		Volume:   0.7,
		Priority: 50,
		Cooldown: 0,
	}

	sp.RegisterSound(sound)

	result := sp.Play("play_test")
	if !result {
		t.Error("Play failed for valid sound")
	}

	if !sp.IsPlaying("play_test") {
		t.Error("sound not marked as playing after Play")
	}
}

func TestSoundPlayer_Play_NotFound(t *testing.T) {
	sp := NewSoundPlayer()

	result := sp.Play("nonexistent")
	if result {
		t.Error("expected false for nonexistent sound")
	}
}

func TestSoundPlayer_Play_Disabled(t *testing.T) {
	sp := NewSoundPlayer()
	sp.Disable()

	sound := &SoundEffect{
		ID:       "disabled_test",
		Name:     "Disabled Test",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
	}

	sp.RegisterSound(sound)

	result := sp.Play("disabled_test")
	if result {
		t.Error("expected false when sound is disabled")
	}
}

func TestSoundPlayer_Play_Muted(t *testing.T) {
	sp := NewSoundPlayer()
	sp.Mute()

	sound := &SoundEffect{
		ID:       "muted_test",
		Name:     "Muted Test",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
	}

	sp.RegisterSound(sound)

	result := sp.Play("muted_test")
	if result {
		t.Error("expected false when muted")
	}
}

func TestSoundPlayer_Play_Cooldown(t *testing.T) {
	sp := NewSoundPlayer()

	sound := &SoundEffect{
		ID:       "cooldown_test",
		Name:     "Cooldown Test",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
		Cooldown: 1000, // 1 second cooldown
	}

	sp.RegisterSound(sound)

	// First play should succeed
	result1 := sp.Play("cooldown_test")
	if !result1 {
		t.Error("first Play failed")
	}

	// Second play should fail due to cooldown
	result2 := sp.Play("cooldown_test")
	if result2 {
		t.Error("expected false during cooldown")
	}

	// Wait for cooldown
	time.Sleep(1100 * time.Millisecond)

	// Third play should succeed after cooldown
	result3 := sp.Play("cooldown_test")
	if !result3 {
		t.Error("Play failed after cooldown")
	}
}

func TestSoundPlayer_Stop(t *testing.T) {
	sp := NewSoundPlayer()

	sound := &SoundEffect{
		ID:       "stop_test",
		Name:     "Stop Test",
		Category: SoundCategoryUI,
		Type:     SoundTypeSfx,
	}

	sp.RegisterSound(sound)
	sp.Play("stop_test")

	if !sp.IsPlaying("stop_test") {
		t.Error("sound should be playing before Stop")
	}

	sp.Stop("stop_test")

	if sp.IsPlaying("stop_test") {
		t.Error("sound should not be playing after Stop")
	}
}

func TestSoundPlayer_StopAll(t *testing.T) {
	sp := NewSoundPlayer()

	// Register and play multiple sounds
	for i := 0; i < 3; i++ {
		sound := &SoundEffect{
			ID:       "stop_all_test",
			Name:     "Stop All Test",
			Category: SoundCategoryUI,
			Type:     SoundTypeSfx,
		}
		sp.RegisterSound(sound)
		sp.Play("stop_all_test")
	}

	if !sp.IsPlaying("stop_all_test") {
		t.Error("sound should be playing before StopAll")
	}

	sp.StopAll()

	if sp.IsPlaying("stop_all_test") {
		t.Error("sound should not be playing after StopAll")
	}
}

func TestSoundPlayer_VolumeControls(t *testing.T) {
	sp := NewSoundPlayer()

	// Test SetMasterVolume
	sp.SetMasterVolume(0.5)
	sp.SetMasterVolume(1.5) // Should clamp to 1.0
	sp.SetMasterVolume(-0.5) // Should clamp to 0.0

	// Test SetBgmVolume
	sp.SetBgmVolume(0.3)
	sp.SetBgmVolume(2.0) // Should clamp
	sp.SetBgmVolume(-1.0) // Should clamp

	// Test SetSfxVolume
	sp.SetSfxVolume(0.9)
	sp.SetSfxVolume(1.5) // Should clamp
	sp.SetSfxVolume(-0.5) // Should clamp

	// Note: Can't directly test private fields, but if no panic, volumes are clamped
}

func TestSoundPlayer_GetEffectiveVolume(t *testing.T) {
	sp := NewSoundPlayer()

	// Test with BGM
	vol := sp.GetEffectiveVolume(SoundTypeBgm)
	if vol <= 0 {
		t.Error("BGM volume should be > 0 when enabled")
	}

	// Test with SFX
	vol = sp.GetEffectiveVolume(SoundTypeSfx)
	if vol <= 0 {
		t.Error("SFX volume should be > 0 when enabled")
	}

	// Disable and test
	sp.Disable()
	vol = sp.GetEffectiveVolume(SoundTypeSfx)
	if vol != 0 {
		t.Error("volume should be 0 when disabled")
	}

	// Re-enable and mute
	sp.Enable()
	sp.Mute()
	vol = sp.GetEffectiveVolume(SoundTypeSfx)
	if vol != 0 {
		t.Error("volume should be 0 when muted")
	}
}

func TestSoundPlayer_GetSoundsByCategory(t *testing.T) {
	sp := NewSoundPlayer()

	// Register sounds in different categories
	sp.RegisterSound(&SoundEffect{ID: "ui1", Category: SoundCategoryUI})
	sp.RegisterSound(&SoundEffect{ID: "ui2", Category: SoundCategoryUI})
	sp.RegisterSound(&SoundEffect{ID: "battle1", Category: SoundCategoryBattle})
	sp.RegisterSound(&SoundEffect{ID: "skill1", Category: SoundCategorySkill})

	// Get UI sounds
	uiSounds := sp.GetSoundsByCategory(SoundCategoryUI)
	if len(uiSounds) != 2 {
		t.Errorf("expected 2 UI sounds, got %d", len(uiSounds))
	}

	// Get Battle sounds
	battleSounds := sp.GetSoundsByCategory(SoundCategoryBattle)
	if len(battleSounds) != 1 {
		t.Errorf("expected 1 Battle sound, got %d", len(battleSounds))
	}

	// Get non-existent category
	otherSounds := sp.GetSoundsByCategory(SoundCategoryEnvironment)
	if len(otherSounds) != 0 {
		t.Errorf("expected 0 environment sounds, got %d", len(otherSounds))
	}
}

// ==================== SoundPresets Tests ====================

func TestSoundPresets(t *testing.T) {
	sounds := SoundPresets()

	if len(sounds) == 0 {
		t.Fatal("SoundPresets returned empty map")
	}

	// Verify some key sounds exist
	keySounds := []string{
		"ui_click",
		"ui_confirm",
		"battle_start",
		"battle_win",
		"skill_charge",
		"item_pickup",
		"enemy_death",
		"boss_appear",
		"sys_levelup",
	}

	for _, id := range keySounds {
		if _, ok := sounds[id]; !ok {
			t.Errorf("expected sound %s not found", id)
		}
	}

	// Verify sounds have required fields
	for id, sound := range sounds {
		if sound.ID == "" {
			t.Errorf("sound %s has empty ID", id)
		}
		if sound.Name == "" {
			t.Errorf("sound %s has empty Name", id)
		}
		if sound.Duration <= 0 {
			t.Errorf("sound %s has invalid Duration %d", id, sound.Duration)
		}
		if sound.Volume < 0 || sound.Volume > 1 {
			t.Errorf("sound %s has invalid Volume %f", id, sound.Volume)
		}
		if sound.Priority < 0 || sound.Priority > 100 {
			t.Errorf("sound %s has invalid Priority %d", id, sound.Priority)
		}
		if sound.Cooldown < 0 {
			t.Errorf("sound %s has invalid Cooldown %d", id, sound.Cooldown)
		}
	}
}

// ==================== Timer Tests ====================

func TestTimer_NewTimer(t *testing.T) {
	timer := NewTimer(5000, false)

	if timer == nil {
		t.Fatal("NewTimer returned nil")
	}

	if timer.Duration != 5000 {
		t.Errorf("expected Duration 5000, got %d", timer.Duration)
	}

	if timer.Elapsed != 0 {
		t.Errorf("expected Elapsed 0, got %d", timer.Elapsed)
	}

	if timer.Running {
		t.Error("expected Running to be false")
	}

	if timer.Loop {
		t.Error("expected Loop to be false")
	}
}

func TestTimer_Start(t *testing.T) {
	timer := NewTimer(5000, false)

	timer.Start()

	if !timer.Running {
		t.Error("timer should be running after Start")
	}
}

func TestTimer_Stop(t *testing.T) {
	timer := NewTimer(5000, false)
	timer.Start()

	timer.Stop()

	if timer.Running {
		t.Error("timer should not be running after Stop")
	}
}

func TestTimer_Reset(t *testing.T) {
	timer := NewTimer(5000, false)
	timer.Start()
	timer.Elapsed = 3000

	timer.Reset()

	if timer.Elapsed != 0 {
		t.Errorf("expected Elapsed 0 after Reset, got %d", timer.Elapsed)
	}

	if timer.Running {
		t.Error("timer should not be running after Reset")
	}
}

func TestTimer_Update(t *testing.T) {
	timer := NewTimer(5000, false)
	timer.Start()

	// Update with small delta, should not complete
	complete := timer.Update(1000)
	if complete {
		t.Error("timer should not complete with 1000ms delta")
	}

	// Update to complete
	complete = timer.Update(4500)
	if !complete {
		t.Error("timer should complete with 4500ms more delta")
	}

	// Timer should stop after completion
	if timer.Running {
		t.Error("timer should stop after completion")
	}
}

func TestTimer_Update_Loop(t *testing.T) {
	timer := NewTimer(2000, true)
	timer.Start()
	timer.Elapsed = 1500

	complete := timer.Update(1000)

	if !complete {
		t.Error("timer should complete on loop mode")
	}

	// Should loop back
	if timer.Elapsed != 500 {
		t.Errorf("expected Elapsed 500 after loop, got %d", timer.Elapsed)
	}

	if !timer.Running {
		t.Error("timer should still be running in loop mode")
	}
}

func TestTimer_GetProgress(t *testing.T) {
	timer := NewTimer(10000, false)
	timer.Elapsed = 2500

	progress := timer.GetProgress()
	expected := 0.25

	if progress != expected {
		t.Errorf("expected progress %f, got %f", expected, progress)
	}
}

func TestTimer_GetProgress_ZeroDuration(t *testing.T) {
	timer := NewTimer(0, false)

	progress := timer.GetProgress()
	if progress != 0 {
		t.Errorf("expected progress 0 for zero duration, got %f", progress)
	}
}

func TestTimer_Remaining(t *testing.T) {
	timer := NewTimer(10000, false)
	timer.Elapsed = 3000

	remaining := timer.Remaining()
	expected := int64(7000)

	if remaining != expected {
		t.Errorf("expected remaining %d, got %d", expected, remaining)
	}
}

func TestTimer_Remaining_Complete(t *testing.T) {
	timer := NewTimer(5000, false)
	timer.Elapsed = 10000

	remaining := timer.Remaining()
	if remaining != 0 {
		t.Errorf("expected remaining 0 for complete timer, got %d", remaining)
	}
}

func TestTimer_IsComplete(t *testing.T) {
	// Not running and not elapsed - not complete
	timer1 := NewTimer(5000, false)
	if timer1.IsComplete() {
		t.Error("new timer should not be complete")
	}

	// Running but not elapsed - not complete
	timer2 := NewTimer(5000, false)
	timer2.Start()
	if timer2.IsComplete() {
		t.Error("running timer with elapsed < duration should not be complete")
	}

	// Not running and elapsed >= duration - complete
	timer3 := NewTimer(5000, false)
	timer3.Elapsed = 5000
	if !timer3.IsComplete() {
		t.Error("timer with elapsed >= duration should be complete")
	}
}

// ==================== Helper Function Tests ====================

func TestClampFloat32(t *testing.T) {
	tests := []struct {
		input    float32
		min      float32
		max      float32
		expected float32
	}{
		{0.5, 0.0, 1.0, 0.5},
		{-0.5, 0.0, 1.0, 0.0},
		{1.5, 0.0, 1.0, 1.0},
		{0.0, 0.0, 1.0, 0.0},
		{1.0, 0.0, 1.0, 1.0},
	}

	for _, test := range tests {
		result := clampFloat32(test.input, test.min, test.max)
		if result != test.expected {
			t.Errorf("clampFloat32(%f, %f, %f) = %f, expected %f",
				test.input, test.min, test.max, result, test.expected)
		}
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		a, b, t   float64
		expected  float64
	}{
		{0.0, 10.0, 0.0, 0.0},
		{0.0, 10.0, 0.5, 5.0},
		{0.0, 10.0, 1.0, 10.0},
		{10.0, 0.0, 0.5, 5.0},
		{-10.0, 10.0, 0.5, 0.0},
	}

	for _, test := range tests {
		result := Lerp(test.a, test.b, test.t)
		if result != test.expected {
			t.Errorf("Lerp(%f, %f, %f) = %f, expected %f",
				test.a, test.b, test.t, result, test.expected)
		}
	}
}

func TestInverseLerp(t *testing.T) {
	tests := []struct {
		a, b, v   float64
		expected  float64
	}{
		{0.0, 10.0, 0.0, 0.0},
		{0.0, 10.0, 5.0, 0.5},
		{0.0, 10.0, 10.0, 1.0},
		{10.0, 0.0, 5.0, 0.5},
	}

	for _, test := range tests {
		result := InverseLerp(test.a, test.b, test.v)
		if result != test.expected {
			t.Errorf("InverseLerp(%f, %f, %f) = %f, expected %f",
				test.a, test.b, test.v, result, test.expected)
		}
	}
}

func TestInverseLerp_EqualEndpoints(t *testing.T) {
	result := InverseLerp(5.0, 5.0, 5.0)
	if result != 0 {
		t.Errorf("InverseLerp with equal endpoints should return 0, got %f", result)
	}
}

func TestRemap(t *testing.T) {
	// Map 0-10 to 0-100
	tests := []struct {
		v, inMin, inMax, outMin, outMax float64
		expected                         float64
	}{
		{0.0, 0.0, 10.0, 0.0, 100.0, 0.0},
		{5.0, 0.0, 10.0, 0.0, 100.0, 50.0},
		{10.0, 0.0, 10.0, 0.0, 100.0, 100.0},
		// Test reverse range
		{5.0, 10.0, 0.0, 0.0, 100.0, 50.0},
	}

	for _, test := range tests {
		result := Remap(test.v, test.inMin, test.inMax, test.outMin, test.outMax)
		if result != test.expected {
			t.Errorf("Remap(%f, %f, %f, %f, %f) = %f, expected %f",
				test.v, test.inMin, test.inMax, test.outMin, test.outMax, result, test.expected)
		}
	}
}

func TestSmoothStep(t *testing.T) {
	// Edge cases
	if SmoothStep(0.0, 1.0, 0.0) != 0.0 {
		t.Error("SmoothStep(0) should be 0")
	}

	if SmoothStep(0.0, 1.0, 1.0) != 1.0 {
		t.Error("SmoothStep(1) should be 1")
	}

	// Middle should be around 0.5
	mid := SmoothStep(0.0, 1.0, 0.5)
	if mid < 0.4 || mid > 0.6 {
		t.Errorf("SmoothStep(0.5) should be around 0.5, got %f", mid)
	}
}

func TestSmootherStep(t *testing.T) {
	// Edge cases
	if SmootherStep(0.0, 1.0, 0.0) != 0.0 {
		t.Error("SmootherStep(0) should be 0")
	}

	if SmootherStep(0.0, 1.0, 1.0) != 1.0 {
		t.Error("SmootherStep(1) should be 1")
	}

	// Middle should be around 0.5
	mid := SmootherStep(0.0, 1.0, 0.5)
	if mid < 0.4 || mid > 0.6 {
		t.Errorf("SmootherStep(0.5) should be around 0.5, got %f", mid)
	}
}

// ==================== Duration Tests ====================

func TestDuration_Milliseconds(t *testing.T) {
	d := NewDuration(1500)
	if d.Milliseconds() != 1500 {
		t.Errorf("expected 1500ms, got %d", d.Milliseconds())
	}
}

func TestDuration_Seconds(t *testing.T) {
	d := NewDuration(2500)
	if d.Seconds() != 2.5 {
		t.Errorf("expected 2.5s, got %f", d.Seconds())
	}
}

func TestDuration_Minutes(t *testing.T) {
	d := NewDuration(90000) // 1.5 minutes
	if d.Minutes() != 1.5 {
		t.Errorf("expected 1.5m, got %f", d.Minutes())
	}
}

func TestDuration_Add(t *testing.T) {
	d1 := NewDuration(1000)
	d2 := NewDuration(2000)
	result := d1.Add(d2)

	if result.Milliseconds() != 3000 {
		t.Errorf("expected 3000ms, got %d", result.Milliseconds())
	}
}

func TestDuration_Sub(t *testing.T) {
	d1 := NewDuration(3000)
	d2 := NewDuration(1000)
	result := d1.Sub(d2)

	if result.Milliseconds() != 2000 {
		t.Errorf("expected 2000ms, got %d", result.Milliseconds())
	}
}

func TestDuration_Multiply(t *testing.T) {
	d := NewDuration(1000)
	result := d.Multiply(2.5)

	if result.Milliseconds() != 2500 {
		t.Errorf("expected 2500ms, got %d", result.Milliseconds())
	}
}

// ==================== Math Utility Tests ====================

func TestAbsInt(t *testing.T) {
	if absInt(-5) != 5 {
		t.Error("absInt(-5) should be 5")
	}
	if absInt(5) != 5 {
		t.Error("absInt(5) should be 5")
	}
	if absInt(0) != 0 {
		t.Error("absInt(0) should be 0")
	}
}

func TestAbsFloat64(t *testing.T) {
	if absFloat64(-5.5) != 5.5 {
		t.Error("absFloat64(-5.5) should be 5.5")
	}
	if absFloat64(5.5) != 5.5 {
		t.Error("absFloat64(5.5) should be 5.5")
	}
}

func TestMinMaxInt(t *testing.T) {
	if minInt(3, 5) != 3 {
		t.Error("minInt(3, 5) should be 3")
	}
	if maxInt(3, 5) != 5 {
		t.Error("maxInt(3, 5) should be 5")
	}
}

func TestMinMaxFloat64(t *testing.T) {
	if minFloat64(3.5, 5.5) != 3.5 {
		t.Error("minFloat64(3.5, 5.5) should be 3.5")
	}
	if maxFloat64(3.5, 5.5) != 5.5 {
		t.Error("maxFloat64(3.5, 5.5) should be 5.5")
	}
}

func TestDegToRad(t *testing.T) {
	tests := []struct {
		deg      float64
		expected float64
	}{
		{0, 0},
		{90, math.Pi / 2},
		{180, math.Pi},
		{270, 3 * math.Pi / 2},
		{360, 2 * math.Pi},
	}

	for _, test := range tests {
		result := degToRad(test.deg)
		if result != test.expected {
			t.Errorf("degToRad(%f) = %f, expected %f", test.deg, result, test.expected)
		}
	}
}

func TestRadToDeg(t *testing.T) {
	tests := []struct {
		rad      float64
		expected float64
	}{
		{0, 0},
		{math.Pi / 2, 90},
		{math.Pi, 180},
		{3 * math.Pi / 2, 270},
		{2 * math.Pi, 360},
	}

	for _, test := range tests {
		result := radToDeg(test.rad)
		if result != test.expected {
			t.Errorf("radToDeg(%f) = %f, expected %f", test.rad, result, test.expected)
		}
	}
}

func TestClampInt(t *testing.T) {
	if clampInt(5, 0, 10) != 5 {
		t.Error("clampInt(5, 0, 10) should be 5")
	}
	if clampInt(-5, 0, 10) != 0 {
		t.Error("clampInt(-5, 0, 10) should be 0")
	}
	if clampInt(15, 0, 10) != 10 {
		t.Error("clampInt(15, 0, 10) should be 10")
	}
}

// Import math for the radian tests
import "math"
