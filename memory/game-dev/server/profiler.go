// Package profiler provides runtime performance profiling utilities
package profiler

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Profiler collects and reports runtime performance metrics
type Profiler struct {
	// Counters
	Goroutines    atomic.Int64
	Threads       atomic.Int64
	MemoryAlloc   atomic.Int64
	MemoryTotal   atomic.Int64
	MemorySys     atomic.Int64
	NumGC         atomic.Int64
	
	// Custom metrics
	RequestCount  atomic.Int64
	ErrorCount    atomic.Int64
	ResponseTimes atomic.Int64 // accumulated in microseconds
	
	// History
	history     []Snapshot
	historyMu   sync.RWMutex
	maxHistory  int
	
	// Sampling
	sampleInterval time.Duration
	lastSample     atomic.Int64
	
	mu sync.RWMutex
}

// Snapshot represents a point-in-time performance snapshot
type Snapshot struct {
	Timestamp    time.Time `json:"timestamp"`
	Goroutines   int64     `json:"goroutines"`
	Threads      int64     `json:"threads"`
	MemoryAlloc  int64     `json:"memory_alloc"`
	MemoryTotal  int64     `json:"memory_total"`
	MemorySys    int64     `json:"memory_sys"`
	NumGC        int64     `json:"num_gc"`
	RequestCount int64     `json:"request_count"`
	ErrorCount   int64     `json:"error_count"`
	AvgRespTime  float64   `json:"avg_resp_time_ms"`
}

// New creates a new profiler
func New(maxHistory int, sampleInterval time.Duration) *Profiler {
	p := &Profiler{
		maxHistory:     maxHistory,
		sampleInterval: sampleInterval,
		history:        make([]Snapshot, 0, maxHistory),
	}
	
	// Initial sample
	p.Sample()
	
	return p
}

// Sample takes a performance snapshot
func (p *Profiler) Sample() *Snapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	snap := Snapshot{
		Timestamp:    time.Now(),
		Goroutines:   int64(runtime.NumGoroutine()),
		Threads:      int64(runtime.GOMAXPROCS(0)),
		MemoryAlloc:  int64(m.Alloc),
		MemoryTotal:  int64(m.TotalAlloc),
		MemorySys:    int64(m.Sys),
		NumGC:        int64(m.NumGC),
		RequestCount: p.RequestCount.Load(),
		ErrorCount:   p.ErrorCount.Load(),
	}
	
	// Calculate average response time
	reqCount := p.RequestCount.Load()
	if reqCount > 0 {
		respTimes := p.ResponseTimes.Load()
		snap.AvgRespTime = float64(respTimes) / float64(reqCount) / 1000.0 // Convert to ms
	}
	
	// Store in history
	p.historyMu.Lock()
	if len(p.history) >= p.maxHistory {
		// Remove oldest
		p.history = p.history[1:]
	}
	p.history = append(p.history, snap)
	p.historyMu.Unlock()
	
	// Update atomic counters
	p.Goroutines.Store(snap.Goroutines)
	p.MemoryAlloc.Store(snap.MemoryAlloc)
	p.NumGC.Store(snap.NumGC)
	
	return &snap
}

// RecordRequest records a request completion
func (p *Profiler) RecordRequest(duration time.Duration) {
	p.RequestCount.Add(1)
	p.ResponseTimes.Add(duration.Microseconds())
}

// RecordError records an error occurrence
func (p *Profiler) RecordError() {
	p.ErrorCount.Add(1)
}

// GetSnapshot returns the latest snapshot
func (p *Profiler) GetSnapshot() *Snapshot {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return &Snapshot{
		Timestamp:    time.Now(),
		Goroutines:   int64(runtime.NumGoroutine()),
		Threads:      int64(runtime.GOMAXPROCS(0)),
		MemoryAlloc:  int64(m.Alloc),
		MemoryTotal:  int64(m.TotalAlloc),
		MemorySys:    int64(m.Sys),
		NumGC:        int64(m.NumGC),
		RequestCount: p.RequestCount.Load(),
		ErrorCount:   p.ErrorCount.Load(),
	}
}

// GetHistory returns the performance history
func (p *Profiler) GetHistory() []Snapshot {
	p.historyMu.RLock()
	defer p.historyMu.RUnlock()
	
	result := make([]Snapshot, len(p.history))
	copy(result, p.history)
	return result
}

// Reset resets all counters
func (p *Profiler) Reset() {
	p.RequestCount.Store(0)
	p.ErrorCount.Store(0)
	p.ResponseTimes.Store(0)
	
	p.historyMu.Lock()
	p.history = p.history[:0]
	p.historyMu.Unlock()
}

// Report generates a human-readable performance report
func (p *Profiler) Report() string {
	snap := p.GetSnapshot()
	
	return fmt.Sprintf(`
=== Performance Report ===
Timestamp: %s

Goroutines: %d
Threads: %d
Memory:
  Alloc: %s
  Total: %s
  Sys:   %s
GC:
  Count: %d

Requests:
  Total: %d
  Errors: %d
  Avg Response: %.2fms

Errors Rate: %.2f%%
`,
		snap.Timestamp.Format("2006-01-02 15:04:05"),
		snap.Goroutines,
		snap.Threads,
		formatBytes(snap.MemoryAlloc),
		formatBytes(snap.MemoryTotal),
		formatBytes(snap.MemorySys),
		snap.NumGC,
		snap.RequestCount,
		snap.ErrorCount,
		snap.AvgRespTime,
		errorRate(snap.RequestCount, snap.ErrorCount),
	)
}

// String implements fmt.Stringer
func (p *Profiler) String() string {
	return p.Report()
}

// formatBytes formats bytes to human readable string
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// errorRate calculates error rate percentage
func errorRate(total, errors int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(errors) / float64(total) * 100
}

// ============== Function Profiler ==============

// FuncProfile records function-level execution time
type FuncProfile struct {
	Name      string
	Count     atomic.Int64
	TotalTime atomic.Int64 // nanoseconds
	MinTime   atomic.Int64
	MaxTime   atomic.Int64
}

// NewFuncProfile creates a new function profiler
func NewFuncProfile(name string) *FuncProfile {
	return &FuncProfile{Name: name}
}

// FuncProfiler manages multiple function profilers
type FuncProfiler struct {
	profiles map[string]*FuncProfile
	mu        sync.RWMutex
}

// NewFuncProfiler creates a new function profiler manager
func NewFuncProfiler() *FuncProfiler {
	return &FuncProfiler{
		profiles: make(map[string]*FuncProfile),
	}
}

// Get returns or creates a function profiler
func (fp *FuncProfiler) Get(name string) *FuncProfile {
	fp.mu.RLock()
	p, ok := fp.profiles[name]
	fp.mu.RUnlock()
	
	if ok {
		return p
	}
	
	fp.mu.Lock()
	defer fp.mu.Unlock()
	
	// Double-check after acquiring write lock
	if p, ok = fp.profiles[name]; ok {
		return p
	}
	
	p = NewFuncProfile(name)
	fp.profiles[name] = p
	return p
}

// Record records function execution time
func (fp *FuncProfiler) Record(name string, duration time.Duration) {
	p := fp.Get(name)
	p.Count.Add(1)
	ns := duration.Nanoseconds()
	p.TotalTime.Add(ns)
	
	// Update min (atomic compare-and-swap)
	for {
		min := p.MinTime.Load()
		if min == 0 || ns < min {
			if p.MinTime.CompareAndSwap(min, ns) {
				break
			}
		} else {
			break
		}
	}
	
	// Update max
	for {
		max := p.MaxTime.Load()
		if ns > max {
			if p.MaxTime.CompareAndSwap(max, ns) {
				break
			}
		} else {
			break
		}
	}
}

// Report generates a human-readable function profiling report
func (fp *FuncProfiler) Report() string {
	fp.mu.RLock()
	defer fp.mu.RUnlock()
	
	result := "=== Function Profile Report ===\n"
	result += fmt.Sprintf("%-30s %12s %12s %12s %12s\n", 
		"Function", "Calls", "Total(ms)", "Min(μs)", "Max(ms)")
	result += "------------------------------------------------------------------------\n"
	
	for name, p := range fp.profiles {
		count := p.Count.Load()
		total := float64(p.TotalTime.Load()) / 1e6
		min := float64(p.MinTime.Load()) / 1e3
		max := float64(p.MaxTime.Load()) / 1e6
		
		result += fmt.Sprintf("%-30s %12d %12.2f %12.2f %12.2f\n",
			name, count, total, min, max)
	}
	
	return result
}

// GetStats returns function profiling stats
func (fp *FuncProfiler) GetStats() map[string]map[string]float64 {
	fp.mu.RLock()
	defer fp.mu.RUnlock()
	
	stats := make(map[string]map[string]float64)
	
	for name, p := range fp.profiles {
		count := float64(p.Count.Load())
		total := float64(p.TotalTime.Load()) / 1e6
		min := float64(p.MinTime.Load()) / 1e3
		max := float64(p.MaxTime.Load()) / 1e6
		avg := total / count * 1000 // Convert to μs
		
		stats[name] = map[string]float64{
			"count":      count,
			"total_ms":   total,
			"min_us":     min,
			"max_ms":     max,
			"avg_us":     avg,
		}
	}
	
	return stats
}

// ============== HTTP Middleware ==============

// HTTPProfiler provides HTTP request profiling middleware
type HTTPProfiler struct {
	profiler *Profiler
	funcProf *FuncProfiler
}

// NewHTTPProfiler creates a new HTTP profiler
func NewHTTPProfiler() *HTTPProfiler {
	return &HTTPProfiler{
		profiler: New(100, 10*time.Second),
		funcProf: NewFuncProfiler(),
	}
}

// RecordRequest records an HTTP request
func (hp *HTTPProfiler) RecordRequest(name string, duration time.Duration, err error) {
	hp.profiler.RecordRequest(duration)
	hp.funcProf.Record(name, duration)
	
	if err != nil {
		hp.profiler.RecordError()
	}
}

// GetProfiler returns the underlying profiler
func (hp *HTTPProfiler) GetProfiler() *Profiler {
	return hp.profiler
}

// GetFuncProfiler returns the function profiler
func (hp *HTTPProfiler) GetFuncProfiler() *FuncProfiler {
	return hp.funcProf
}

// ============== Block Profiler ==============

// BlockProfile records blocking operation times
type BlockProfile struct {
	mu           sync.Mutex
	operations   map[string][]time.Duration
	maxSamples   int
}

// NewBlockProfile creates a new block profiler
func NewBlockProfile(maxSamples int) *BlockProfile {
	return &BlockProfile{
		operations: make(map[string][]time.Duration),
		maxSamples: maxSamples,
	}
}

// Record records a blocking operation
func (bp *BlockProfile) Record(name string, duration time.Duration) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	
	samples := bp.operations[name]
	samples = append(samples, duration)
	
	if len(samples) > bp.maxSamples {
		samples = samples[1:]
	}
	
	bp.operations[name] = samples
}

// GetStats returns statistics for a blocking operation
func (bp *BlockProfile) GetStats(name string) (count int, avg, p50, p95, p99 time.Duration) {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	
	samples := bp.operations[name]
	if len(samples) == 0 {
		return
	}
	
	count = len(samples)
	var total time.Duration
	for _, d := range samples {
		total += d
	}
	avg = total / time.Duration(count)
	
	// Sort for percentiles
	sorted := make([]time.Duration, count)
	copy(sorted, samples)
	
	p50 = sorted[len(sorted)/2]
	p95 = sorted[int(float64(len(sorted))*0.95)]
	if len(sorted) > 1 {
		p99 = sorted[int(float64(len(sorted))*0.99)]
	}
	
	return
}

// ============== Global Profiler ==============

var (
	globalProfiler   *Profiler
	globalFuncProf   *FuncProfiler
	profilerOnce     sync.Once
	funcProfOnce     sync.Once
)

// GetGlobalProfiler returns the global profiler (singleton)
func GetGlobalProfiler() *Profiler {
	profilerOnce.Do(func() {
		globalProfiler = New(1000, 5*time.Second)
	})
	return globalProfiler
}

// GetGlobalFuncProfiler returns the global function profiler (singleton)
func GetGlobalFuncProfiler() *FuncProfiler {
	funcProfOnce.Do(func() {
		globalFuncProf = NewFuncProfiler()
	})
	return globalFuncProf
}

// RecordFunction records function execution time to global profiler
func RecordFunction(name string, duration time.Duration) {
	GetGlobalFuncProf().Record(name, duration)
}
