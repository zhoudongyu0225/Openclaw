// Package health provides health checking and monitoring utilities
package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// HealthStatus represents the health status of the server
type HealthStatus struct {
	Status      string            `json:"status"` // healthy/degraded/unhealthy
	Version     string            `json:"version"`
	Uptime      time.Duration    `json:"uptime"`
	Timestamp   time.Time        `json:"timestamp"`
	Checks      map[string]Check `json:"checks"`
	System      SystemInfo        `json:"system"`
	GameStats   GameStats        `json:"game_stats"`
}

// Check represents a single health check
type Check struct {
	Status  string        `json:"status"` // pass/fail/warn
	Latency time.Duration `json:"latency"`
	Message string        `json:"message,omitempty"`
}

// SystemInfo contains system-level information
type SystemInfo struct {
	GoVersion   string  `json:"go_version"`
	NumCPU      int     `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
	MemAlloc    uint64  `json:"mem_alloc"`
	MemTotal    uint64  `json:"mem_total"`
	MemSys      uint64  `json:"mem_sys"`
	GCEnabled   bool    `json:"gc_enabled"`
}

// GameStats contains game-specific statistics
type GameStats struct {
	ActiveRooms    int `json:"active_rooms"`
	ActivePlayers  int `json:"active_players"`
	TotalPlayers   int `json:"total_players"`
	TotalGames     int `json:"total_games"`
	AvgMatchTime   float64 `json:"avg_match_time_ms"`
	PeakPlayers    int `json:"peak_players"`
	MessagesPerSec float64 `json:"messages_per_sec"`
}

// HealthChecker performs health checks
type HealthChecker struct {
	startTime   time.Time
	checks      map[string]CheckFunc
	mu          sync.RWMutex
	stats       GameStats
	version     string
	lastMetrics time.Time
	metrics     []float64
}

// CheckFunc is a function that performs a health check
type CheckFunc func() Check

// NewHealthChecker creates a new health checker
func NewHealthChecker(version string) *HealthChecker {
	return &HealthChecker{
		startTime: time.Now(),
		checks:    make(map[string]CheckFunc),
		version:   version,
		metrics:   make([]float64, 0, 60),
	}
}

// RegisterCheck registers a health check
func (hc *HealthChecker) RegisterCheck(name string, fn CheckFunc) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.checks[name] = fn
}

// UpdateStats updates game statistics
func (hc *HealthChecker) UpdateStats(stats GameStats) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.stats = stats
	
	// Track metrics for last 60 seconds
	hc.metrics = append(hc.metrics, stats.MessagesPerSec)
	if len(hc.metrics) > 60 {
		hc.metrics = hc.metrics[1:]
	}
	hc.lastMetrics = time.Now()
}

// GetStatus returns the current health status
func (hc *HealthChecker) GetStatus() HealthStatus {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	
	status := HealthStatus{
		Status:    "healthy",
		Version:   hc.version,
		Uptime:    time.Since(hc.startTime),
		Timestamp: time.Now(),
		Checks:    make(map[string]Check),
		System:    hc.getSystemInfo(),
		GameStats: hc.stats,
	}
	
	// Run all checks
	for name, fn := range hc.checks {
		check := fn()
		status.Checks[name] = check
		
		if check.Status == "fail" && status.Status == "healthy" {
			status.Status = "degraded"
		}
		if check.Status == "fail" {
			status.Status = "unhealthy"
		}
	}
	
	// Calculate average messages per second
	if len(hc.metrics) > 0 {
		var sum float64
		for _, m := range hc.metrics {
			sum += m
		}
		status.GameStats.MessagesPerSec = sum / float64(len(hc.metrics))
	}
	
	return status
}

func (hc *HealthChecker) getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return SystemInfo{
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		MemAlloc:    m.Alloc,
		MemTotal:    m.TotalAlloc,
		MemSys:      m.Sys,
		GCEnabled:   true,
	}
}

// Handler returns an HTTP handler for health checks
func (hc *HealthChecker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := hc.GetStatus()
		
		w.Header().Set("Content-Type", "application/json")
		
		if status.Status == "unhealthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else if status.Status == "degraded" {
			w.WriteHeader(http.StatusOK) // Still OK, just degraded
		} else {
			w.WriteHeader(http.StatusOK)
		}
		
		json.NewEncoder(w).Encode(status)
	}
}

// ReadyHandler returns a handler for readiness probes
func (hc *HealthChecker) ReadyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check critical dependencies
		hc.mu.RLock()
		_, hasChecks := hc.checks["database"]
		hc.mu.RUnlock()
		
		if hasChecks {
			check := hc.checks["database"]()
			if check.Status == "fail" {
				http.Error(w, "Not ready", http.StatusServiceUnavailable)
				return
			}
		}
		
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ready")
	}
}

// LiveHandler returns a handler for liveness probes
func (hc *HealthChecker) LiveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Alive")
	}
}

// DefaultHealthChecker creates a health checker with default checks
func DefaultHealthChecker(version string) *HealthChecker {
	hc := NewHealthChecker(version)
	
	// Register default checks
	hc.RegisterCheck("memory", func() Check {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		// Warn if using more than 80% of available memory
		limit := uint64(1024 * 1024 * 1024) // 1GB default limit
		usage := float64(m.Alloc) / float64(limit)
		
		if usage > 0.8 {
			return Check{
				Status:  "warn",
				Message: fmt.Sprintf("Memory usage: %.1f%%", usage*100),
			}
		}
		
		return Check{
			Status:  "pass",
			Message: fmt.Sprintf("Memory usage: %.1f%%", usage*100),
		}
	})
	
	hc.RegisterCheck("goroutines", func() Check {
		count := runtime.NumGoroutine()
		
		if count > 1000 {
			return Check{
				Status:  "warn",
				Message: fmt.Sprintf("High goroutine count: %d", count),
			}
		}
		
		return Check{
			Status:  "pass",
			Message: fmt.Sprintf("Goroutines: %d", count),
		}
	})
	
	return hc
}

// StartMetricsCollector starts a background metrics collector
func (hc *HealthChecker) StartMetricsCollector(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			
			hc.mu.Lock()
			hc.metrics = append(hc.metrics, float64(m.Alloc))
			if len(hc.metrics) > 300 { // 5 minutes at 1s interval
				hc.metrics = hc.metrics[1:]
			}
			hc.mu.Unlock()
		}
	}()
}

// PrintStatus prints the health status to console
func (hc *HealthChecker) PrintStatus() {
	status := hc.GetStatus()
	
	fmt.Println("\n========== Health Status ==========")
	fmt.Printf("Status:     %s\n", status.Status)
	fmt.Printf("Version:    %s\n", status.Version)
	fmt.Printf("Uptime:     %v\n", status.Uptime)
	fmt.Println()
	fmt.Println("--- System ---")
	fmt.Printf("Go Version:    %s\n", status.System.GoVersion)
	fmt.Printf("CPUs:          %d\n", status.System.NumCPU)
	fmt.Printf("Goroutines:   %d\n", status.System.NumGoroutine)
	fmt.Printf("Memory Alloc: %.2f MB\n", float64(status.System.MemAlloc)/1024/1024)
	fmt.Println()
	fmt.Println("--- Game Stats ---")
	fmt.Printf("Active Rooms:   %d\n", status.GameStats.ActiveRooms)
	fmt.Printf("Active Players: %d\n", status.GameStats.ActivePlayers)
	fmt.Printf("Peak Players:   %d\n", status.GameStats.PeakPlayers)
	fmt.Printf("Messages/sec:  %.2f\n", status.GameStats.MessagesPerSec)
	fmt.Println()
	fmt.Println("--- Checks ---")
	for name, check := range status.Checks {
		fmt.Printf("[%s] %s: %s\n", check.Status, name, check.Message)
	}
	fmt.Println("====================================\n")
}
