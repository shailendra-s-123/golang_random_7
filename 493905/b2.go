package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func simulateCPUBoundTask(wg *sync.WaitGroup, id int, duration time.Duration) {
	defer wg.Done()
	fmt.Printf("CPU-bound Goroutine %d: Starting...\n", id)

	// Simulate CPU-bound work by performing a series of calculations
	startTime := time.Now()
	for startTime.Add(duration) .After(time.Now()) {
		for i := 0; i < 1000000; i++ {
			// Do some CPU work
		}
	}

	fmt.Printf("CPU-bound Goroutine %d: Finished.\n", id)
}

func simulateIOBoundTask(wg *sync.WaitGroup, id int, duration time.Duration) {
	defer wg.Done()
	fmt.Printf("I/O-bound Goroutine %d: Starting...\n", id)

	// Simulate I/O operation by introducing random sleeps
	startTime := time.Now()
	for startTime.Add(duration) .After(time.Now()) {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}

	fmt.Printf("I/O-bound Goroutine %d: Finished.\n", id)
}

func main() {
	// Number of CPU-bound goroutines
	numCPUBoundGoroutines := 5
	// Number of I/O-bound goroutines
	numIOBoundGoroutines := 5
	// Duration of the CPU-bound task
	cpuDuration := time.Second * 3 // 3 seconds for each CPU-bound task
	// Duration of the I/O-bound task
	ioDuration := time.Second * 5 // 5 seconds for each I/O-bound task

	// Create wait groups for CPU-bound and I/O-bound tasks
	var cpuWg sync.WaitGroup
	var ioWg sync.WaitGroup

	// Start CPU-bound goroutines
	fmt.Println("Starting CPU-bound tasks...")
	startTime := time.Now()
	for i := 0; i < numCPUBoundGoroutines; i++ {
		cpuWg.Add(1)
		go simulateCPUBoundTask(&cpuWg, i+1, cpuDuration)
	}
	
	// Start I/O-bound goroutines
	fmt.Println("Starting I/O-bound tasks...")
	for i := 0; i < numIOBoundGoroutines; i++ {
		ioWg.Add(1)
		go simulateIOBoundTask(&ioWg, i+1, ioDuration)
	}

	// Wait for all CPU-bound goroutines to finish
	fmt.Println("Waiting for CPU-bound tasks to complete...")
	cpuWg.Wait()

	// Wait for all I/O-bound goroutines to finish
	fmt.Println("Waiting for I/O-bound tasks to complete...")
	ioWg.Wait()

	// Calculate total execution time
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	fmt.Printf("\nTotal execution time: %v\n", totalTime)
	fmt.Printf("CPU-bound tasks: %d, I/O-bound tasks: %d\n", numCPUBoundGoroutines, numIOBoundGoroutines)
	fmt.Printf("Simulated CPU-bound duration per task: %v\n", cpuDuration)
	fmt.Printf("Simulated I/O-bound duration per task: %v\n", ioDuration)
}