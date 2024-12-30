package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// PerformanceMetrics holds the server's performance data
type PerformanceMetrics struct {
	TotalRequests int
	TotalLatency  time.Duration
	MinLatency    time.Duration
	MaxLatency    time.Duration
	Successes     int
	Errors        int
}

// sendRequest sends a GET request and updates performance metrics
func sendRequest(url string, metrics *PerformanceMetrics, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure Done is always called to avoid blocking

	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)

	// Update performance metrics
	metrics.TotalRequests++
	metrics.TotalLatency += duration

	if err != nil {
		metrics.Errors++
		log.Printf("Error: %v", err)
		return
	}
	metrics.Successes++

	// Update min and max latency
	if metrics.MinLatency == 0 || duration < metrics.MinLatency {
		metrics.MinLatency = duration
	}
	if duration > metrics.MaxLatency {
		metrics.MaxLatency = duration
	}

	resp.Body.Close()
}

// constantLoad generates a constant load with fixed rate
func constantLoad(url string, rate int, duration time.Duration, metrics *PerformanceMetrics) {
	timeout := time.After(duration)
	for {
		select {
		case <-timeout:
			return
		default:
			var wg sync.WaitGroup

			// Send requests at a constant rate
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

// burstLoad generates a burst load where multiple requests are sent at once
func burstLoad(url string, rate int, burstSize int, duration time.Duration, metrics *PerformanceMetrics) {
	timeout := time.After(duration)
	for {
		select {
		case <-timeout:
			return
		default:
			var wg sync.WaitGroup

			// Send requests in bursts
			for i := 0; i < burstSize; i++ {
				wg.Add(1) // Add to the WaitGroup before spawning a goroutine
				go sendRequest(url, metrics, &wg)
			}

			// Wait for all requests in the burst to complete before continuing
			wg.Wait()

			time.Sleep(time.Second / time.Duration(rate)) // Control the rate of burst generation
		}
	}
}

// spikeLoad generates spike loads by increasing traffic suddenly for short periods
func spikeLoad(url string, rate int, spikeSize int, duration time.Duration, metrics *PerformanceMetrics) {
	timeout := time.After(duration)
	for {
		select {
		case <-timeout:
			return
		default:
			var wg sync.WaitGroup

			// Increase load suddenly (spike)
			for i := 0; i < spikeSize; i++ {
				wg.Add(1)
				go sendRequest(url, metrics, &wg)
			}

			// Wait for all spike requests to finish
			wg.Wait()

			time.Sleep(time.Second / time.Duration(rate)) // Sleep for one second to simulate periodic spikes
		}
	}
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
	// Setup server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond) // Simulate work being done
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
	metrics := &PerformanceMetrics{}

	// Test different load patterns
	go constantLoad("http://localhost:8080", 5, 10*time.Second, metrics) // Constant load
	go burstLoad("http://localhost:8080", 2, 20, 10*time.Second, metrics)  // Burst load
	go spikeLoad("http://localhost:8080", 1, 50, 10*time.Second, metrics)   // Spike load

	// Run for a total of 30 seconds
	time.Sleep(30 * time.Second)

	// Analyze performance after the test
	analyzeMetrics(metrics)
}