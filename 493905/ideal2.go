
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Simulate an I/O-bound task by introducing a delay
func simulateIOTask(delay time.Duration, wg *sync.WaitGroup, index int, activeGoroutines *int64, ioWaitTimes chan<- time.Duration) {
	defer wg.Done()

	// Increment the active goroutines count
	atomic.AddInt64(activeGoroutines, 1)
	defer atomic.AddInt64(activeGoroutines, -1) // Decrement the active goroutines count when the task finishes

	// Record the start time of the I/O task
	beforeIO := time.Now()

	// Simulate I/O by sleeping for a random duration
	time.Sleep(delay)

	// Calculate how long this goroutine spent waiting for I/O
	ioWaitTime := time.Since(beforeIO)

	// Log the I/O wait time for this task
	fmt.Printf("Goroutine %d finished with I/O wait time: %s\n", index, ioWaitTime)

	// Send the I/O wait time to a channel for further aggregation
	ioWaitTimes <- ioWaitTime
}

func monitorActiveGoroutines(activeGoroutines *int64) {
	// This goroutine monitors and prints the number of active goroutines every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// This line ensures that we print the number of active goroutines
		fmt.Printf("Active Goroutines: %d\n", atomic.LoadInt64(activeGoroutines))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Number of I/O-bound tasks to simulate
	const numIOTasks = 100

	// Different levels of concurrency to test
	numConcurrent := []int{1, 5, 10, 20, 50, 100}

	// For each concurrency level, perform an experiment
	for _, concurrency := range numConcurrent {
		var wg sync.WaitGroup
		activeGoroutines := int64(0)
		ioWaitTimes := make(chan time.Duration, numIOTasks)

		// Start monitoring goroutines' activity in the background
		go monitorActiveGoroutines(&activeGoroutines)

		start := time.Now()

		// Limit concurrency by using a buffered channel to control the number of concurrently running goroutines
		limit := make(chan struct{}, concurrency)

		// Create and start goroutines for the tasks
		for i := 0; i < numIOTasks; i++ {
			// Random I/O delay between 0 and 1 second
			delay := time.Duration(rand.Intn(1000)) * time.Millisecond
			wg.Add(1)
			go func(index int) {
				// Ensure we don't exceed the concurrency limit
				limit <- struct{}{}
				defer func() { <-limit }() // Release the token after the task is finished

				// Perform the I/O task
				simulateIOTask(delay, &wg, index, &activeGoroutines, ioWaitTimes)
			}(i)
		}

		// Wait for all goroutines to finish
		wg.Wait()

		// Close the channel to stop collecting I/O wait times
		close(ioWaitTimes)

		// Calculate the total execution time
		elapsed := time.Since(start)
		fmt.Printf("Test with %d concurrent tasks completed. Total execution time = %s\n", concurrency, elapsed)

		// Analyze and report average I/O wait time
		var totalIOWait time.Duration
		for waitTime := range ioWaitTimes {
			totalIOWait += waitTime
		}
		fmt.Printf("Average I/O wait time per task: %s\n", totalIOWait/time.Duration(numIOTasks))
	}
}
