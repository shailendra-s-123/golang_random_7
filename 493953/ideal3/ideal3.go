package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	Latency     []time.Duration `json:"latency"`
	Throughput  int             `json:"throughput"`
	ErrorRate   float64         `json:"error_rate"`
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

// Predictive Feedback Loop (Mock Example)
func predictiveLoadAdjustment() LoadPattern {
	// Here we use a basic random predictor as a placeholder for ML algorithms.
	patterns := []LoadPattern{Constant, Burst, Spike}
	return patterns[rand.Intn(len(patterns))]
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
			generateLoad(pattern, 10*time.Second, adaptiveLoadAdjuster)
			time.Sleep(5 * time.Second)
		}
	}()

	// Run for a specific duration
	time.Sleep(1 * time.Minute)
	generateReport()
}
