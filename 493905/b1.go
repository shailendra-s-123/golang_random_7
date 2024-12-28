package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func ioBoundGoroutine(id int, wg *sync.WaitGroup, numIterations int, waitTimeMillis int) {
	defer wg.Done()

	for i := 0; i < numIterations; i++ {
		// Simulate I/O operation with a random wait
		time.Sleep(time.Duration(rand.Intn(waitTimeMillis)) * time.Millisecond)
		// Record goroutine completion time
		fmt.Printf("Goroutine %d: Iteration %d completed at %v\n", id, i, time.Now())
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	numGoroutines := 10
	numIterations := 5
	waitTimeMillis := 100

	var wg sync.WaitGroup

	startTime := time.Now()

	// Create and run I/O-bound goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go ioBoundGoroutine(i, &wg, numIterations, waitTimeMillis)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	endTime := time.Now()

	totalExecutionTime := endTime.Sub(startTime)

	fmt.Printf("\nTotal execution time: %v\n", totalExecutionTime)
}  