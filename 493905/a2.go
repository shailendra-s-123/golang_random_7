package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Simulate an I/O-bound task by introducing a delay
func simulateIOTask(delay time.Duration, start time.Time, wg *sync.WaitGroup, index int, activeGoroutines *int64, ioWaitTimes chan<- time.Duration) {
	defer wg.Done()
	defer func() {
		atomic.AddInt64(activeGoroutines, -1)
	}()
	atomic.AddInt64(activeGoroutines, 1)

	// Record the time before starting the sleep
	beforeIO := time.Now()
	time.Sleep(delay)
	afterIO := time.Now()

	// Calculate I/O wait time
	ioWaitTime := afterIO.Sub(beforeIO)

	// Log the number of active goroutines
	fmt.Printf("Goroutine %d active, I/O wait time: %s\n", index, ioWaitTime)

	// Send I/O wait time to the channel for further analysis
	ioWaitTimes <- ioWaitTime
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const numIOTasks = 100 // Number of I/O-bound tasks
	var numConcurrent = []int{1, 5, 10, 20, 50, 100} // Different levels of concurrency

	// For each level of concurrency, perform an experiment
	for _, concurrency := range numConcurrent {
		var wg sync.WaitGroup
		activeGoroutines := int64(0)
		ioWaitTimes := make(chan<- time.Duration, numIOTasks)

		start := time.Now()

		// Create a limited channel to control concurrency
		goroutineChan := make(chan struct{}, concurrency)

		// Create and start goroutines for the tasks
		for i := 0; i < numIOTasks; i++ {
			delay := time.Duration(rand.Intn(1000)) * time.Millisecond // Random delay up to 1 second
			wg.Add(1)
			go func() {
				// Wait to enter the goroutine channel to limit concurrency
				goroutineChan <- struct{}{}
				defer func() { <-goroutineChan }()

				simulateIOTask(delay, start, &wg, i, &activeGoroutines, ioWaitTimes)
			}()
		}

		wg.Wait()

		elapsed := time.Since(start)
		fmt.Printf("Experienced %d concurrent tasks: Total execution time = %s\n", concurrency, elapsed)

		// Analyze I/O wait times
		var totalIOWait time.Duration
		for waitTime := range ioWaitTimes {
			totalIOWait += waitTime
		}
		fmt.Printf("Average I/O wait time per task: %s\n", totalIOWait/time.Duration(numIOTasks))
	}
}