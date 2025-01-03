package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Simulate an I/O-bound task by introducing a delay
func simulateIOTask(delay time.Duration, start time.Time, wg *sync.WaitGroup, index int) {
	time.Sleep(delay)
	elapsed := time.Since(start)
	fmt.Printf("Goroutine %d finished in: %s\n", index, elapsed)
	wg.Done()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const numIOTasks = 100 // Number of I/O-bound tasks
	var numConcurrent = []int{1, 5, 10, 20, 50, 100} // Different levels of concurrency

	// For each level of concurrency, perform an experiment
	for _, concurrency := range numConcurrent {
		var wg sync.WaitGroup

		start := time.Now()

		// Create and start goroutines for the tasks
		for i := 0; i < numIOTasks; i++ {
			delay := time.Duration(rand.Intn(1000)) * time.Millisecond // Random delay up to 1 second
			wg.Add(1)
			go simulateIOTask(delay, start, &wg, i)
		}

		wg.Wait()

		elapsed := time.Since(start)
		fmt.Printf("Experienced %d concurrent tasks: Total execution time = %s\n", concurrency, elapsed)
	}
}