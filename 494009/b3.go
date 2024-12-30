
package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// RedisPool represents a connection pool to the Redis server.
type RedisPool struct {
	client *redis.Client
}

// NewRedisPool initializes a new Redis connection pool.
func NewRedisPool(addr string) (*RedisPool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // No password
		DB:       0,  // Default database
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "failed to connect to Redis")
	}
	return &RedisPool{client: client}, nil
}

// Close cleans up the Redis connection pool.
func (p *RedisPool) Close() error {
	return p.client.Close()
}

// Operation encapsulates the logic for interacting with Redis.
type Operation struct {
	pool      *RedisPool
	maxWorkers int
	wg        sync.WaitGroup
	semaphore chan struct{}
}

// performTask performs a single Redis operation with proper error handling and resource cleanup.
func (op *Operation) performTask(ctx context.Context, key, value string) {
	defer op.wg.Done()

	<-op.semaphore // Acquire semaphore to limit concurrency
	defer func() {
		op.semaphore <- struct{}{} // Release semaphore after completion
	}()

	// The rest of the task logic remains the same
	// ...
}

// AdjustConcurrency dynamically adjusts the number of active goroutines based on system load.
func (op *Operation) AdjustConcurrency() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			numCpus := runtime.NumCPU()
			newMaxWorkers := numCpus * 2

			if newMaxWorkers != op.maxWorkers {
				log.Printf("Adjusting concurrency from %d to %d workers", op.maxWorkers, newMaxWorkers)
				op.maxWorkers = newMaxWorkers
				op.semaphore = make(chan struct{}, op.maxWorkers)
			}
		}
	}
}

func main() {
	const (
		redisAddr  = "localhost:6379"
		numTasks   = 100000
	)

	// Initialize Redis connection pool.
	pool, err := NewRedisPool(redisAddr)
	if err != nil {
		log.Fatalf("Error initializing Redis pool: %v", err)
	}
	defer func() {
		if closeErr := pool.Close(); closeErr != nil {
			log.Printf("Error closing Redis pool: %v", closeErr)
		}
	}()

	op := &Operation{
		pool:      pool,
		maxWorkers: runtime.NumCPU() * 2,
		wg:        sync.WaitGroup{},
		semaphore:  make(chan struct{}, runtime.NumCPU()*2),
	}

	// Start dynamic concurrency adjustment loop
	go op.AdjustConcurrency()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run tasks concurrently with adaptive throttling
	for i := 0; i < numTasks; i++ {
		op.wg.Add(1)
		go op.performTask(ctx, fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i))