package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Object struct {
	ID int
}

type ObjectPool struct {
	pool chan *Object
}

func NewObjectPool(size int) *ObjectPool {
	pool := &ObjectPool{
		pool: make(chan *Object, size),
	}

	// Pre-populate the pool with objects
	for i := 0; i < size; i++ {
		pool.pool <- &Object{ID: i}
	}
	return pool
}

func (op *ObjectPool) Acquire() *Object {
	select {
	case obj := <-op.pool:
		return obj
	default:
		// Pool is empty, could expand if needed
		return nil
	}
}

func (op *ObjectPool) Release(obj *Object) {
	select {
	case op.pool <- obj:
	default:
		// Pool is full, discard or handle appropriately
	}
}

func processCallback(pool *ObjectPool, wg *sync.WaitGroup) {
	defer wg.Done()

	// Simulate processing an object
	obj := pool.Acquire()
	if obj == nil {
		fmt.Println("No object available to process")
		return
	}
	// Simulate work being done
	fmt.Printf("Processing Object ID: %d\n", obj.ID)
	time.Sleep(time.Millisecond * 100)

	// Release the object back to the pool
	pool.Release(obj)
}

func workerPool(pool *ObjectPool, numWorkers int, totalTasks int, wg *sync.WaitGroup) {
	// Create a worker pool with bounded concurrency
	sem := make(chan struct{}, numWorkers) // Semaphore to limit concurrent workers

	// Distribute tasks across workers
	for i := 0; i < totalTasks; i++ {
		sem <- struct{}{} // Ensure max concurrency limit
		go func() {
			defer func() { <-sem }() // Release the semaphore when done

			// Randomize object selection for better load balancing
			processCallback(pool, wg)
		}()
	}
}

func main() {
	// Object pool and worker pool configurations
	objectPool := NewObjectPool(10) // Pool size of 10
	numWorkers := 5                // Max concurrent workers
	totalTasks := 50               // Number of tasks to process
	var wg sync.WaitGroup

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Simulate processing with a bounded worker pool
	wg.Add(totalTasks) // Add total tasks to the wait group for processing
	workerPool(objectPool, numWorkers, totalTasks, &wg)

	// Wait for all tasks to complete
	wg.Wait()
	fmt.Println("All tasks completed.")
}
