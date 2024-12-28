package main

import (
	"fmt"
	"sync"
	"time"
)

func simulateIO(wg *sync.WaitGroup, id int, duration time.Duration) {
	defer wg.Done()
	// Simulate I/O operation by sleeping
	fmt.Printf("Goroutine %d: Simulating I/O operation...\n", id)
	time.Sleep(duration)
	fmt.Printf("Goroutine %d: I/O operation completed.\n", id)
}

func main() {
	// Number of goroutines to run
	numGoroutines := 10
	// Duration of the I/O operation
	ioDuration := time.Second * 2 // 2 seconds for each I/O operation

	// Create a wait group to track the goroutines
	var wg sync.WaitGroup

	// Start goroutines
	startTime := time.Now()
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go simulateIO(&wg, i+1, ioDuration)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Calculate total execution time
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	fmt.Printf("Total execution time: %v\n", totalTime)
	fmt.Printf("Number of goroutines: %d\n", numGoroutines)
	fmt.Printf("Simulated I/O duration per goroutine: %v\n", ioDuration)
}