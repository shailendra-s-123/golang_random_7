package main

import (
	"fmt"
	"sync"
	"time"
)

var sharedCounter int
var sharedCounterMutex sync.Mutex

func simulateIO(wg *sync.WaitGroup, id int, duration time.Duration) {
	defer wg.Done()
	// Simulate I/O operation by sleeping
	fmt.Printf("Goroutine %d (I/O): Simulating I/O operation...\n", id)
	time.Sleep(duration)
	fmt.Printf("Goroutine %d (I/O): I/O operation completed.\n", id)
}

func simulateCPU(wg *sync.WaitGroup, id int, iterations int) {
	defer wg.Done()
	// Simulate CPU-bound operation by performing busy work
	fmt.Printf("Goroutine %d (CPU): Performing CPU-bound operations...\n", id)
	for i := 0; i < iterations; i++ {
		// Lock the mutex for shared resource access
		sharedCounterMutex.Lock()
		sharedCounter++
		sharedCounterMutex.Unlock()
	}
	fmt.Printf("Goroutine %d (CPU): CPU-bound operations completed.\n", id)
}

func main() {
	// Number of goroutines to run for each type
	numIOGoroutines := 10
	numCPUGoroutines := 10
	// Duration of the I/O operation
	ioDuration := time.Second * 2 // 2 seconds for each I/O operation
	// Iterations for CPU-bound operations
	cpuIterations := 1_000_000 // 1 million iterations

	// Create wait groups to track the goroutines
	var ioWG sync.WaitGroup
	var cpuWG sync.WaitGroup

	// Start I/O-bound goroutines
	startTime := time.Now()
	for i := 0; i < numIOGoroutines; i++ {
		ioWG.Add(1)
		go simulateIO(&ioWG, i+1, ioDuration)
	}

	// Start CPU-bound goroutines
	for i := 0; i < numCPUGoroutines; i++ {
		cpuWG.Add(1)
		go simulateCPU(&cpuWG, i+1, cpuIterations)
	}

	// Wait for all I/O-bound and CPU-bound goroutines to finish
	ioWG.Wait()
	cpuWG.Wait()

	// Calculate total execution time
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	fmt.Printf("Total execution time: %v\n", totalTime)
	fmt.Printf("Number of I/O-bound goroutines: %d\n", numIOGoroutines)
	fmt.Printf("Number of CPU-bound goroutines: %d\n", numCPUGoroutines)
	fmt.Printf("Simulated I/O duration per goroutine: %v\n", ioDuration)
	fmt.Printf("Simulated CPU iterations per goroutine: %d\n", cpuIterations)
	fmt.Printf("Shared counter value after completion: %d\n", sharedCounter)
}