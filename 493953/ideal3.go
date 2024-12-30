package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	Latency     []time.Duration `json:"latency"`
	Throughput  int             `json:"throughput"`
	ErrorRate   float64         `json:"error_rate"`
	Timestamp   []time.Time     `json:"timestamp"`
}

type LoadPattern string

const (
	Constant LoadPattern = "constant"
	Burst    LoadPattern = "burst"
	Spike    LoadPattern = "spike"
)

// Adaptive Throttling Configuration
const (
	LatencyThreshold = 200 * time.Millisecond
	AnomalyThreshold = 1.5
)

var (
	metrics Metrics
	mu      sync.Mutex
)

// Simulates HTTP handler with artificial delay
func handler(w http.ResponseWriter, r *http.Request) {
	// Simulate variable latency
	delay := time.Duration(rand.Intn(100)+50) * time.Millisecond
	time.Sleep(delay)

	if rand.Float64() < 0.1 { // Simulate 10% error rate
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "error")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")

	// Record metrics
	mu.Lock()
	metrics.Latency = append(metrics.Latency, delay)
	metrics.Throughput++
	metrics.Timestamp = append(metrics.Timestamp, time.Now())
	mu.Unlock()
}

// Load Generator
func generateLoad(pattern LoadPattern, duration time.Duration, adjustLoad func()) {
	end := time.Now().Add(duration)
	var wg sync.WaitGroup

	for time.Now().Before(end) {
		adjustLoad()

		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get("http://localhost:8080")
			if err != nil {
				log.Printf("Request error: %v", err)
				return
			}
			defer resp.Body.Close()
		}()

		if pattern == Burst {
			time.Sleep(10 * time.Millisecond) // Higher frequency
		} else if pattern == Spike {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		} else {
			time.Sleep(50 * time.Millisecond)
		}
	}

	wg.Wait()
}

// Adaptive Load Adjuster
func adaptiveLoadAdjuster() {
	mu.Lock()
	defer mu.Unlock()

	if len(metrics.Latency) > 0 {
		avgLatency := calculateAverageLatency()
		if avgLatency > LatencyThreshold {
			log.Printf("Latency too high (%v). Reducing load.", avgLatency)
		}
	}
}

// Calculates average latency
func calculateAverageLatency() time.Duration {
	totalLatency := time.Duration(0)
	for _, l := range metrics.Latency {
		totalLatency += l
	}
	return totalLatency / time.Duration(len(metrics.Latency))
}

// Anomaly Detection using Z-score
func detectAnomalies() bool {
	mu.Lock()
	defer mu.Unlock()

	if len(metrics.Latency) < 2 {
		return false // Not enough data
	}

	mean, stdDev := calculateStats(metrics.Latency)
	latestLatency := float64(metrics.Latency[len(metrics.Latency)-1])
	zScore := math.Abs((latestLatency - mean) / stdDev)

	if zScore > AnomalyThreshold {
		log.Printf("Anomaly detected! Z-score: %.2f", zScore)
		return true
	}

	return false
}

// Calculates mean and standard deviation
func calculateStats(latencies []time.Duration) (float64, float64) {
	mean := 0.0
	for _, l := range latencies {
		mean += float64(l)
	}
	mean /= float64(len(latencies))

	variance := 0.0
	for _, l := range latencies {
		variance += math.Pow(float64(l)-mean, 2)
	}
	variance /= float64(len(latencies))

	return mean, math.Sqrt(variance)
}

// Reporting
func generateReport() {
	mu.Lock()
	defer mu.Unlock()

	report, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		log.Fatalf("Error generating report: %v", err)
	}
	fmt.Println("Performance Metrics:")
	fmt.Println(string(report))
}

// Predictive Feedback Loop using Moving Average
func predictiveLoadAdjustment() LoadPattern {
	mu.Lock()
	defer mu.Unlock()

	if len(metrics.Latency) < 3 {
		return Constant // Not enough data
	}

	avgLatency := calculateAverageLatency()
	if avgLatency > LatencyThreshold {
		return Burst
	}

	return Spike
}

func main() {
	// Set up HTTP server
	http.HandleFunc("/", handler)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	log.Println("Server started on port 8080")

	rand.Seed(time.Now().UnixNano())

	// Load testing
	go func() {
		for {
			pattern := predictiveLoadAdjustment()
			log.Printf("Adjusting to pattern: %s", pattern)
			generateLoad(pattern, 10*time.Second, func() {
				if detectAnomalies() {
					log.Println("Anomaly detected during load generation.")
				}
				adaptiveLoadAdjuster()
			})
			time.Sleep(5 * time.Second)
		}
	}()

	// Run for a specific duration
	time.Sleep(1 * time.Minute)
	generateReport()
}
