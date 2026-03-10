// Package loadtest provides load testing utilities for the game server
package loadtest

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Config holds load test configuration
type Config struct {
	NumPlayers      int           // Number of concurrent players
	NumRooms        int           // Number of rooms to create
	RampUpDuration  time.Duration // Time to ramp up all players
	TestDuration    time.Duration // Total test duration
	MsgInterval     time.Duration // Interval between messages
	ThinkTimeMin    time.Duration // Min think time between actions
	ThinkTimeMax    time.Duration // Max think time between actions
	EnableBattle   bool          // Enable battle simulation
	EnableGifts    bool          // Enable gift sending
	EnableDanmaku  bool          // Enable danmaku sending
}

// Result holds load test results
type Result struct {
	TotalRequests    int64
	SuccessRequests  int64
	FailedRequests   int64
	TotalLatency     time.Duration
	MinLatency       time.Duration
	MaxLatency       time.Duration
	AvgLatency       time.Duration
	P50Latency       time.Duration
	P90Latency       time.Duration
	P99Latency       time.Duration
	RoomsCreated     int64
	PlayersJoined    int64
	GiftsSent        int64
	DanmakuSent      int64
	Errors           []string
	Duration         time.Duration
	Throughput       float64
}

// LoadTester is the main load testing struct
type LoadTester struct {
	config  Config
	result  Result
	ctx     context.Context
	cancel  context.CancelFunc
	latencies []time.Duration
	mu      sync.RWMutex
}

// New creates a new load tester
func New(config Config) *LoadTester {
	ctx, cancel := context.WithCancel(context.Background())
	return &LoadTester{
		config:   config,
		ctx:      ctx,
		cancel:   cancel,
		latencies: make([]time.Duration, 0, 10000),
	}
}

// Run starts the load test
func (lt *LoadTester) Run() *Result {
	startTime := time.Now()
	
	// Create workers
	var wg sync.WaitGroup
	
	// Ramp up players
	rampUpInterval := lt.config.RampUpDuration / time.Duration(lt.config.NumPlayers)
	
	for i := 0; i < lt.config.NumPlayers; i++ {
		wg.Add(1)
		go func(playerID int) {
			defer wg.Done()
			
			// Ramp up delay
			if i > 0 {
				time.Sleep(time.Duration(playerID) * rampUpInterval)
			}
			
			lt.runPlayer(playerID)
		}(i)
	}
	
	// Wait for test duration or cancellation
	select {
	case <-lt.ctx.Done():
	case <-time.After(lt.config.TestDuration):
	}
	
	lt.cancel()
	wg.Wait()
	
	// Calculate results
	lt.calculateResults(startTime)
	
	return &lt.result
}

func (lt *LoadTester) runPlayer(playerID int) {
	player := fmt.Sprintf("player_%d", playerID)
	rng := rand.New(rand.NewSource(int64(playerID)))
	
	// Simulate player lifecycle
	for {
		select {
		case <-lt.ctx.Done():
			return
		default:
		}
		
		// Random think time
		thinkTime := lt.config.ThinkTimeMin + 
			time.Duration(rng.Int63n(int64(lt.config.ThinkTimeMax - lt.config.ThinkTimeMin)))
		time.Sleep(thinkTime)
		
		// Random action
		action := rng.Intn(100)
		
		switch {
		case action < 30:
			// Send message
			lt.simulateMessage(player)
		case action < 50 && lt.config.EnableBattle:
			// Battle action
			lt.simulateBattleAction(player)
		case action < 60 && lt.config.EnableGifts:
			// Send gift
			lt.simulateGift(player)
		case action < 70 && lt.config.EnableDanmaku:
			// Send danmaku
			lt.simulateDanmaku(player)
		default:
			// Heartbeat/ping
			lt.simulateHeartbeat(player)
		}
		
		atomic.AddInt64(&lt.result.TotalRequests, 1)
	}
}

func (lt *LoadTester) simulateMessage(player string) {
	start := time.Now()
	
	// Simulate message send latency
	latency := 5*time.Millisecond + time.Duration(rand.Int63n(20))
	time.Sleep(latency)
	
	lt.recordLatency(time.Since(start))
	atomic.AddInt64(&lt.result.SuccessRequests, 1)
}

func (lt *LoadTester) simulateBattleAction(player string) {
	start := time.Now()
	
	// Simulate battle action
	latency := 10*time.Millisecond + time.Duration(rand.Int63n(30))
	time.Sleep(latency)
	
	lt.recordLatency(time.Since(start))
}

func (lt *LoadTester) simulateGift(player string) {
	start := time.Now()
	
	// Simulate gift send
	latency := 20*time.Millisecond + time.Duration(rand.Int63n(50))
	time.Sleep(latency)
	
	lt.recordLatency(time.Since(start))
	atomic.AddInt64(&lt.result.GiftsSent, 1)
}

func (lt *LoadTester) simulateDanmaku(player string) {
	start := time.Now()
	
	// Simulate danmaku send
	latency := 3*time.Millisecond + time.Duration(rand.Int63n(10))
	time.Sleep(latency)
	
	lt.recordLatency(time.Since(start))
	atomic.AddInt64(&lt.result.DanmakuSent, 1)
}

func (lt *LoadTester) simulateHeartbeat(player string) {
	start := time.Now()
	
	// Simulate heartbeat
	latency := 2*time.Millisecond + time.Duration(rand.Int63n(5))
	time.Sleep(latency)
	
	lt.recordLatency(time.Since(start))
}

func (lt *LoadTester) recordLatency(latency time.Duration) {
	lt.mu.Lock()
	lt.latencies = append(lt.latencies, latency)
	lt.mu.Unlock()
	
	atomic.AddInt64(&lt.result.TotalLatency, int64(latency))
	
	// Update min/max atomically
	for {
		current := atomic.LoadInt64(&lt.result.MinLatency.Nanoseconds())
		if current == 0 || int64(latency) < current {
			if atomic.CompareAndSwapInt64((*int64)(&lt.result.MinLatency), current, int64(latency)) {
				break
			}
		} else {
			break
		}
	}
	
	for {
		current := atomic.LoadInt64(&lt.result.MaxLatency.Nanoseconds())
		if int64(latency) > current {
			if atomic.CompareAndSwapInt64((*int64)(&lt.result.MaxLatency), current, int64(latency)) {
				break
			}
		} else {
			break
		}
	}
}

func (lt *LoadTester) calculateResults(startTime time.Time) {
	lt.result.Duration = time.Since(startTime)
	
	if lt.result.TotalRequests > 0 {
		lt.result.AvgLatency = time.Duration(atomic.LoadInt64(&lt.result.TotalLatency)) / time.Duration(lt.result.TotalRequests)
	}
	
	// Calculate percentiles
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	
	if len(lt.latencies) > 0 {
		// Sort latencies
		sorted := make([]time.Duration, len(lt.latencies))
		copy(sorted, lt.latencies)
		
		lt.result.P50Latency = sorted[len(sorted)*50/100]
		lt.result.P90Latency = sorted[len(sorted)*90/100]
		lt.result.P99Latency = sorted[len(sorted)*99/100]
	}
	
	// Calculate throughput
	if lt.result.Duration.Seconds() > 0 {
		lt.result.Throughput = float64(lt.result.TotalRequests) / lt.result.Duration.Seconds()
	}
}

// AddError records an error during testing
func (lt *LoadTester) AddError(err string) {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	lt.result.Errors = append(lt.result.Errors, err)
	atomic.AddInt64(&lt.result.FailedRequests, 1)
}

// Stop stops the load test
func (lt *LoadTester) Stop() {
	lt.cancel()
}

// PrintResults prints the test results
func (r *Result) PrintResults() {
	fmt.Println("\n========== Load Test Results ==========")
	fmt.Printf("Duration:           %v\n", r.Duration)
	fmt.Printf("Total Requests:     %d\n", r.TotalRequests)
	fmt.Printf("Success:           %d\n", r.SuccessRequests)
	fmt.Printf("Failed:            %d\n", r.FailedRequests)
	fmt.Printf("Throughput:        %.2f req/s\n", r.Throughput)
	fmt.Println()
	fmt.Printf("Min Latency:       %v\n", r.MinLatency)
	fmt.Printf("Avg Latency:       %v\n", r.AvgLatency)
	fmt.Printf("Max Latency:       %v\n", r.MaxLatency)
	fmt.Printf("P50 Latency:       %v\n", r.P50Latency)
	fmt.Printf("P90 Latency:       %v\n", r.P90Latency)
	fmt.Printf("P99 Latency:       %v\n", r.P99Latency)
	fmt.Println()
	fmt.Printf("Rooms Created:     %d\n", r.RoomsCreated)
	fmt.Printf("Players Joined:    %d\n", r.PlayersJoined)
	fmt.Printf("Gifts Sent:        %d\n", r.GiftsSent)
	fmt.Printf("Danmaku Sent:      %d\n", r.DanmakuSent)
	
	if len(r.Errors) > 0 {
		fmt.Println()
		fmt.Println("Errors:")
		errorCounts := make(map[string]int)
		for _, e := range r.Errors {
			errorCounts[e]++
		}
		for err, count := range errorCounts {
			fmt.Printf("  - %s: %d\n", err, count)
		}
	}
	
	fmt.Println("========================================\n")
}

// BenchmarkRunner runs micro-benchmarks
type BenchmarkRunner struct {
	results map[string]time.Duration
	mu      sync.Mutex
}

// NewBenchmarkRunner creates a new benchmark runner
func NewBenchmarkRunner() *BenchmarkRunner {
	return &BenchmarkRunner{
		results: make(map[string]time.Duration),
	}
}

// RunBenchmark runs a named benchmark
func (br *BenchmarkRunner) RunBenchmark(name string, fn func()) {
	// Warm up
	for i := 0; i < 3; i++ {
		fn()
	}
	
	// Actual benchmark
	start := time.Now()
	iterations := 1000
	for i := 0; i < iterations; i++ {
		fn()
	}
	elapsed := time.Since(start)
	
	br.mu.Lock()
	br.results[name] = elapsed / time.Duration(iterations)
	br.mu.Unlock()
}

// PrintResults prints benchmark results
func (br *BenchmarkRunner) PrintResults() {
	fmt.Println("\n========== Benchmark Results ==========")
	for name, duration := range br.results {
		fmt.Printf("%-30s %v\n", name, duration)
	}
	fmt.Println("========================================\n")
}
