package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type PerformanceData struct {
	Timestamp time.Time        `json:"timestamp"`
	Latency    time.Duration    `json:"latency"`
	Requests   int64            `json:"requests"`
	Errors     int64            `json:"errors"`
	LoadRate   int              `json:"loadRate"`
}

type PerformanceMetrics struct {
	TotalRequests     int64
	TotalLatency      time.Duration
	MinLatency        time.Duration
	MaxLatency        time.Duration
	Successes         int64
	Errors            int64
	LastLatencySample time.Duration
	LoadThreshold     int
}

// sendRequest sends a GET request and updates performance metrics
func sendRequest(url string, metrics *PerformanceMetrics, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	// Update performance metrics
	atomic.AddInt64(&metrics.TotalRequests, 1)
	metrics.TotalLatency += duration

	if err != nil {
		atomic.AddInt64(&metrics.Errors, 1)
		log.Printf("Error: %v", err)
		return
	}
	atomic.AddInt64(&metrics.Successes, 1)

	// Update min and max latency
	if metrics.MinLatency == 0 || duration < metrics.MinLatency {
		metrics.MinLatency = duration
	}
	if duration > metrics.MaxLatency {
		metrics.MaxLatency = duration
	}

	// Update the last latency sample for adaptive throttling
	metrics.LastLatencySample = duration

	resp.Body.Close()
}

// adaptiveLoad generates traffic with adaptive throttling
func adaptiveLoad(url string, initialRate int, duration time.Duration, metrics *PerformanceMetrics, historicalData []PerformanceData) {
	timeout := time.After(duration)
	rate := initialRate

	// Predictive logic based on heuristic trend analysis
	for {
		select {
		case <-timeout:
			return
		default:
			var wg sync.WaitGroup
			currentLoadRate := rate

			// Predict future load based on historical data (basic linear trend)
			// Note: This is a simple heuristic. In a real scenario, use ML.
			predictedLoad := rate
			if len(historicalData) > 2 {
				trend := (historicalData[len(historicalData)-1].LoadRate - historicalData[len(historicalData)-2].LoadRate) / (historicalData[len(historicalData)-1].Timestamp.Unix() - historicalData[len(historicalData)-2].Timestamp.Unix())
				predictedLoad += trend * int(time.Since(historicalData[len(historicalData)-1].Timestamp).Seconds())
			}

			// Adjust the rate based on the server's current performance and predictions
			if metrics.LastLatencySample > time.Millisecond*150 {
				// Reduce load if latency is high
				currentLoadRate = max(currentLoadRate/2, 1)
				log.Printf("Adjusted rate: %d (Latency high)", currentLoadRate)
			} else if predictedLoad > rate {
				// Increase load if predicted load is higher than current rate
				currentLoadRate = min(currentLoadRate*2, predictedLoad)
				log.Printf("Adjusted rate: %d (Predicted load higher)", currentLoadRate)
			}

			// Send requests at the adjusted rate
			for i := 0; i < currentLoadRate; i++ {
				wg.Add(1)
				go sendRequest(url, metrics, &wg)
			}

			// Wait for all requests in this cycle to finish
			wg.Wait()
			time.Sleep(time.Second)

			// Record performance data
			historicalData = append(historicalData, PerformanceData{
				Timestamp: time.Now(),
				Latency:    metrics.LastLatencySample,
				Requests:   metrics.TotalRequests,
				Errors:     metrics.Errors,
				LoadRate:   currentLoadRate,
			})

			// Print performance data for dashboard visualization
			data, err := json.MarshalIndent(historicalData, "", "  ")
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("Performance Data:\n%s\n", data)
		}
	}
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Setup server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond) // Simulate random work load
		w.Write([]byte("Hello, world!"))
	})
	go func() {
		log.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Setup performance metrics
	metrics := &PerformanceMetrics{LoadThreshold: 100} // Adjust load threshold as needed

	// Historical data for predicting traffic patterns
	historicalData := []PerformanceData{
		{Timestamp: time.Now().Add(-10 * time.Second), Latency: 0, Requests: 0, Errors: 0, LoadRate: 5},
		{Timestamp: time.Now().Add(-5 * time.Second), Latency: 0, Requests: 0, Errors: 0, LoadRate: 5},
	}

	// Test with adaptive load
	go adaptiveLoad("http://localhost:8080", 5, 30*time.Second, metrics, historicalData) // Adaptive load

	// Run for a total of 30 seconds
	time.Sleep(30 * time.Second)

	// Analyze performance after the test
	analyzeMetrics(metrics)
}