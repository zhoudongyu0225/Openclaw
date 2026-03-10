// Package stress provides advanced stress testing scenarios
package stress

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"time"
)

// StressConfig holds stress test configuration
type StressConfig struct {
	TargetHost       string        // Target server host
	TargetPort       int           // Target server port
	NumClients       int           // Number of concurrent clients
	RampUpTime       time.Duration // Ramp up period
	StressDuration   time.Duration // Total stress duration
	ConnectionType   string        // tcp/http/ws
	PayloadSize      int           // Size of each message
	BurstMode        bool          // Enable burst traffic
	SlowStart        bool          // Enable slow start
	MaxConnections   int           // Max connections per client
	RequestTimeout   time.Duration // Timeout for each request
}

// StressResult holds stress test results
type StressResult struct {
	TotalBytes       uint64
	TotalPackets     uint64
	TotalConnections uint64
	FailedConnections uint64
	PeakConcurrent   int
	Errors           map[string]int
	StartTime        time.Time
	EndTime          time.Time
	Duration         time.Duration
	ThroughputMB     float64
	PacketsPerSec    float64
	AvgLatency       time.Duration
	MaxLatency       time.Duration
}

// StressRunner orchestrates stress testing
type StressRunner struct {
	config  StressConfig
	result  StressResult
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.Mutex
	
	activeConns atomic.Int32
	peakConns   atomic.Int32
	latencies   []time.Duration
}

// NewStressRunner creates a new stress runner
func NewStressRunner(config StressConfig) *StressRunner {
	ctx, cancel := context.WithCancel(context.Background())
	return &StressRunner{
		config:    config,
		ctx:       ctx,
		cancel:    cancel,
		latencies: make([]time.Duration, 0, 100000),
		result:    StressResult{Errors: make(map[string]int)},
	}
}

// RunTCPStress runs TCP stress test
func (sr *StressRunner) RunTCPStress() *StressResult {
	sr.result.StartTime = time.Now()
	
	// Start pprof server for diagnostics
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	
	// Ramp up clients
	clientsPerBatch := sr.config.NumClients / 10
	batchDelay := sr.config.RampUpTime / 10
	
	for batch := 0; batch < 10; batch++ {
		for i := 0; i < clientsPerBatch; i++ {
			sr.wg.Add(1)
			go sr.runTCPClient(batch*clientsPerBatch + i)
		}
		time.Sleep(batchDelay)
	}
	
	// Wait for stress duration
	select {
	case <-sr.ctx.Done():
	case <-time.After(sr.config.StressDuration):
	}
	
	sr.cancel()
	sr.wg.Wait()
	
	sr.result.EndTime = time.Now()
	sr.result.Duration = sr.result.EndTime.Sub(sr.result.StartTime)
	sr.result.PeakConcurrent = int(sr.peakConns.Load())
	sr.calculateStats()
	
	return &sr.result
}

func (sr *StressRunner) runTCPClient(clientID int) {
	defer sr.wg.Done()
	
	// Create payload
	payload := make([]byte, sr.config.PayloadSize)
	rng := rand.New(rand.NewSource(int64(clientID)))
	rng.Read(payload)
	
	addr := fmt.Sprintf("%s:%d", sr.config.TargetHost, sr.config.TargetPort)
	
	for {
		select {
		case <-sr.ctx.Done():
			return
		default:
		}
		
		conn, err := net.DialTimeout("tcp", addr, sr.config.RequestTimeout)
		if err != nil {
			sr.recordError("connection_failed")
			atomic.AddUint64(&sr.result.FailedConnections, 1)
			time.Sleep(time.Second)
			continue
		}
		
		atomic.AddUint64(&sr.result.TotalConnections, 1)
		sr.activeConns.Add(1)
		
		// Update peak
		current := sr.activeConns.Load()
		for {
			peak := sr.peakConns.Load()
			if current > peak {
				if sr.peakConns.CompareAndSwap(peak, current) {
					break
				}
			} else {
				break
			}
		}
		
		// Send data in bursts
		if sr.config.BurstMode {
			for bursts := 0; bursts < 10; bursts++ {
				select {
				case <-sr.ctx.Done():
					sr.activeConns.Add(-1)
					conn.Close()
					return
				default:
				}
				
				start := time.Now()
				n, err := conn.Write(payload)
				elapsed := time.Since(start)
				
				if err != nil {
					sr.recordError("write_failed")
					break
				}
				
				atomic.AddUint64(&sr.result.TotalBytes, uint64(n))
				atomic.AddUint64(&sr.result.TotalPackets, 1)
				sr.recordLatency(elapsed)
				
				// Small delay between bursts
				time.Sleep(time.Millisecond * 10)
			}
		} else {
			// Continuous mode
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			
			for {
				select {
				case <-sr.ctx.Done():
					sr.activeConns.Add(-1)
					conn.Close()
					return
				case <-ticker.C:
				}
				
				start := time.Now()
				n, err := conn.Write(payload)
				elapsed := time.Since(start)
				
				if err != nil {
					sr.recordError("write_failed")
					break
				}
				
				atomic.AddUint64(&sr.result.TotalBytes, uint64(n))
				atomic.AddUint64(&sr.result.TotalPackets, 1)
				sr.recordLatency(elapsed)
			}
		}
		
		sr.activeConns.Add(-1)
		conn.Close()
		
		// Reconnect delay
		if sr.config.SlowStart {
			time.Sleep(time.Duration(rng.Int63n(1000)) * time.Millisecond)
		}
	}
}

func (sr *StressRunner) recordError(err string) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.result.Errors[err]++
}

func (sr *StressRunner) recordLatency(lat time.Duration) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	if len(sr.latencies) < 100000 {
		sr.latencies = append(sr.latencies, lat)
	}
}

func (sr *StressRunner) calculateStats() {
	// Calculate throughput
	if sr.result.Duration.Seconds() > 0 {
		sr.result.ThroughputMB = float64(sr.result.TotalBytes) / 1024 / 1024 / sr.result.Duration.Seconds()
		sr.result.PacketsPerSec = float64(sr.result.TotalPackets) / sr.result.Duration.Seconds()
	}
	
	// Calculate latency stats
	if len(sr.latencies) > 0 {
		var total time.Duration
		var max time.Duration
		
		for _, lat := range sr.latencies {
			total += lat
			if lat > max {
				max = lat
			}
		}
		
		sr.result.AvgLatency = total / time.Duration(len(sr.latencies))
		sr.result.MaxLatency = max
	}
}

// PrintResults prints the stress test results
func (sr *StressResult) PrintResults() {
	fmt.Println("\n========== Stress Test Results ==========")
	fmt.Printf("Duration:          %v\n", sr.Duration)
	fmt.Printf("Total Bytes:       %d (%.2f MB)\n", sr.TotalBytes, float64(sr.TotalBytes)/1024/1024)
	fmt.Printf("Total Packets:     %d\n", sr.TotalPackets)
	fmt.Printf("Total Connections: %d\n", sr.TotalConnections)
	fmt.Printf("Failed Connections:%d\n", sr.FailedConnections)
	fmt.Printf("Peak Concurrent:   %d\n", sr.PeakConcurrent)
	fmt.Println()
	fmt.Printf("Throughput:        %.2f MB/s\n", sr.ThroughputMB)
	fmt.Printf("Packets/sec:       %.2f\n", sr.PacketsPerSec)
	fmt.Println()
	fmt.Printf("Avg Latency:       %v\n", sr.AvgLatency)
	fmt.Printf("Max Latency:       %v\n", sr.MaxLatency)
	
	if len(sr.Errors) > 0 {
		fmt.Println()
		fmt.Println("Errors:")
		for err, count := range sr.Errors {
			fmt.Printf("  - %s: %d\n", err, count)
		}
	}
	
	fmt.Println("==========================================\n")
}

// DoSConfig holds DoS attack simulation config
type DoSConfig struct {
	TargetHost    string
	TargetPort    int
	AttackType    string // syn_flood/udp_flood/http_flood
	Rate          int    // packets per second
	Duration      time.Duration
	PacketSize    int
	SourceIPRange string
}

// DoSRunner simulates DoS attacks for testing
type DoSRunner struct {
	config DoSConfig
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewDoSRunner creates a new DoS runner
func NewDoSRunner(config DoSConfig) *DoSRunner {
	ctx, cancel := context.WithCancel(context.Background())
	return &DoSRunner{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run starts the DoS simulation
func (dr *DoSRunner) Run() {
	fmt.Printf("Starting %s attack simulation...\n", dr.AttackType)
	
	attackers := dr.config.Rate / 10
	if attackers < 1 {
		attackers = 1
	}
	
	for i := 0; i < attackers; i++ {
		dr.wg.Add(1)
		go dr.attacker(i)
	}
	
	<-time.After(dr.config.Duration)
	dr.cancel()
	dr.wg.Wait()
	
	fmt.Println("DoS simulation completed")
}

func (dr *DoSRunner) attacker(id int) {
	defer dr.wg.Done()
	
	addr := fmt.Sprintf("%s:%d", dr.config.TargetHost, dr.config.TargetPort)
	
	ticker := time.NewTicker(time.Second / time.Duration(dr.config.Rate/10))
	defer ticker.Stop()
	
	for {
		select {
		case <-dr.ctx.Done():
			return
		case <-ticker.C:
		}
		
		switch dr.config.AttackType {
		case "udp_flood":
			dr.udpFlood(addr)
		case "http_flood":
			dr.httpFlood()
		}
	}
}

func (dr *DoSRunner) udpFlood(addr string) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(dr.config.TargetHost),
		Port: dr.config.TargetPort,
	})
	if err != nil {
		return
	}
	defer conn.Close()
	
	payload := make([]byte, dr.config.PacketSize)
	rand.Read(payload)
	
	conn.Write(payload)
}

func (dr *DoSRunner) httpFlood() {
	url := fmt.Sprintf("http://%s:%d/", dr.config.TargetHost, dr.config.TargetPort)
	http.Get(url)
}
