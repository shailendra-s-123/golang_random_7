package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// PerformanceMetrics holds the server's performance data
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
func adaptiveLoad(url string, initialRate int, duration time.Duration, metrics *PerformanceMetrics) {
	timeout := time.After(duration)
	rate := initialRate

	for {
		select {
		case <-timeout:
			return
		default:
			var wg sync.WaitGroup

			// Adjust the rate based on the server's current performance
			if metrics.LastLatencySample > time.Millisecond*150 {
				// Reduce load if latency is high
				rate = max(rate/2, 1) // Never go below 1 request per second
				log.Printf("Adjusted rate: %d (Latency high)", rate)
			} else {
				// Increase load if latency is within acceptable range
				rate = min(rate*2, initialRate) // Limit increase to initial rate
				log.Printf("Adjusted rate: %d (Latency within range)", rate)
			}

			// Send requests at the adjusted rate
			for i := 0; i < rate; i++ {
				wg.Add(1)
				go sendRequest(url, metrics, &wg)
			}

			// Wait for all requests in this cycle to finish
			wg.Wait()
			time.Sleep(time.Second) // Sleep for one second before sending another burst
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

// analyzeMetrics analyzes the performance metrics and suggests optimizations
func analyzeMetrics(metrics *PerformanceMetrics) {
	// Calculate throughput as requests per second
	throughput := float64(metrics.TotalRequests) / metrics.TotalLatency.Seconds()

	// Output the performance analysis
	fmt.Printf("Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("Total Latency: %v\n", metrics.TotalLatency)
	fmt.Printf("Min Latency: %v\n", metrics.MinLatency)
	fmt.Printf("Max Latency: %v\n", metrics.MaxLatency)
	fmt.Printf("Successes: %d\n", metrics.Successes)
	fmt.Printf("Errors: %d\n", metrics.Errors)
	fmt.Printf("Throughput (requests/sec): %.2f\n", throughput)

	// Identify potential bottlenecks
	if metrics.Errors > 0 {
		fmt.Println("Warning: Errors occurred. Investigate the cause.")
	}
	if throughput < 50 {
		fmt.Println("Optimization Suggestion: Throughput is low. Consider scaling up or optimizing handler performance.")
	}
	if metrics.MaxLatency > time.Millisecond*100 {
		fmt.Println("Optimization Suggestion: Latency is high. Investigate server or network performance.")
	}
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

	// Test with adaptive load
	go adaptiveLoad("http://localhost:8080", 5, 30*time.Second, metrics) // Adaptive load

	// Run for a total of 30 seconds
	time.Sleep(30 * time.Second)

	// Analyze performance after the test
	analyzeMetrics(metrics)
}