package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Simple HTTP handler that simulates a database query or other work.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulate random work
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	// Set up HTTP server
	http.HandleFunc("/", helloHandler)
	fmt.Println("Starting HTTP server on :8080")
	go http.ListenAndServe(":8080", nil)

	// Wait a moment for the server to start
	time.Sleep(1 * time.Second)

	// Initialize rand seed
	rand.Seed(time.Now().UnixNano())

	// Load generator
	var wg sync.WaitGroup
	totalRequests := 100
	timeout := 30 * time.Second
	const concurrentRequests = 10
	var successes, errors int

	// Define traffic patterns: constant, burst, spike
	patterns := []func(int, int){
		constantLoad, burstLoad, spikeLoad,
	}

	// Constant load pattern
	constantLoad := func(duration int, rate int) {
		for i := 0; i < totalRequests; i++ {
			time.Sleep(time.Duration(1000/rate) * time.Millisecond)
			sendRequest()
		}
	}

	// Burst load pattern
	burstLoad := func(duration int, rate int) {
		for i := 0; i < totalRequests; i++ {
			if rand.Intn(2) == 0 { // Sending twice as fast during burst intervals
				time.Sleep(time.Duration(500/rate) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(1000/rate) * time.Millisecond)
			}
			sendRequest()
		}
	}

	// Spike load pattern
	spikeLoad := func(duration int, rate int) {
		for i := 0; i < totalRequests; i++ {
			// Slow for 2/3 of the time, fast for 1/3
			if rand.Intn(3) == 0 {
				time.Sleep(time.Duration(100/rate) * time.Millisecond)
			} else {
				time.Sleep(time.Duration(1000/rate) * time.Millisecond)
			}
			sendRequest()
		}
	}

	// Function to send a single HTTP request
	sendRequest := func() {
		startTime := time.Now()
		resp, err := http.Get("http://localhost:8080/")
		if err != nil {
			errors++
		} else {
			successes++
			resp.Body.Close()
		}
		endTime := time.Now()
		duration := endTime.Sub(startTime)

		// Log request outcome
		if err != nil {
			fmt.Printf("Error: %v, Duration: %v\n", err, duration)
		} else {
			fmt.Printf("Success: Duration: %v\n", duration)
		}
	}

	// Choose a pattern and run the load test
	patterns[0](20, concurrentRequests) // Constant load for 20 seconds with 10 concurrent requests
	patterns[1](20, concurrentRequests) // Burst load for 20 seconds with 10 concurrent requests
	patterns[2](20, concurrentRequests) // Spike load for 20 seconds with 10 concurrent requests

	// Wait for all goroutines to finish
	wg.Wait()

	// Analyze results
	throughput := float64(successes) / (time.Now().Unix() - time.Now().Unix()-float64(timeout.Seconds()))
	fmt.Printf("\nTotal Requests: %d, Successes: %d, Errors: %d\n", totalRequests, successes, errors)
	fmt.Printf("Throughput: %.2f requests/second\n", throughput)
	fmt.Printf("Error Rate: %.2f%%\n", float64(errors)*100/float64(successes+errors))
}