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
	Throughput    float64
}

// sendRequest sends a GET request and updates performance metrics
func sendRequest(url string, metrics *PerformanceMetrics, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()
	resp, err := http.Get(url)
	duration := time.Since(start)
	metrics.TotalRequests++
	metrics.TotalLatency += duration
	if err != nil {
		metrics.Errors++
		log.Printf("Error: %v", err)
		return
	}
	metrics.Successes++
	if metrics.MinLatency == 0 || duration < metrics.MinLatency {
		metrics.MinLatency = duration
	}
	if duration > metrics.MaxLatency {
		metrics.MaxLatency = duration
	}
	resp.Body.Close()
}

// adaptiveLoadGenerator generates traffic with adjustable rate based on performance metrics
func adaptiveLoadGenerator(url string, rate *int, metrics *PerformanceMetrics) {
	maxRate := 100
	thresholdLatency := time.Millisecond * 200 // Adjust the latency threshold as needed
	for {
		var wg sync.WaitGroup
		for i := 0; i < *rate; i++ {
			wg.Add(1)
			go sendRequest(url, metrics, &wg)
		}
		wg.Wait()
		metrics.Throughput = float64(metrics.TotalRequests) / metrics.TotalLatency.Seconds()
		// Adaptive throttling: Reduce load if latency exceeds threshold
		if metrics.MaxLatency > thresholdLatency {
			if *rate > 1 {
				*rate--
			}
			log.Printf("Throttling load to %d requests per second due to high latency.\n", *rate)
		} else if metrics.Throughput < float64(maxRate)/2 {
			if *rate < maxRate {
				*rate++
			}
			log.Printf("Increasing load to %d requests per second.\n", *rate)
		}
		time.Sleep(time.Second)
	}
}

// analyzeMetrics analyzes the performance metrics and suggests optimizations
func analyzeMetrics(metrics *PerformanceMetrics) {
	// Calculate throughput as requests per second
	throughput := float64(metrics.TotalRequests) / metrics.TotalLatency.Seconds()
	fmt.Printf("Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("Total Latency: %v\n", metrics.TotalLatency)
	fmt.Printf("Min Latency: %v\n", metrics.MinLatency)
	fmt.Printf("Max Latency: %v\n", metrics.MaxLatency)
	fmt.Printf("Successes: %d\n", metrics.Successes)
	fmt.Printf("Errors: %d\n", metrics.Errors)
	fmt.Printf("Throughput (requests/sec): %.2f\n", throughput)
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