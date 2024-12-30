package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/cpu"
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
	pool         *RedisPool
	wg           *sync.WaitGroup
	workChan     chan<- int
	maxConcurrency int
}

// performTask performs a single Redis operation with proper error handling and resource cleanup.
func (op *Operation) performTask(ctx context.Context, key, value string) {
	defer op.wg.Done()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Use defer for cleanup and error recovery.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Set a key-value pair in Redis.
	err := op.pool.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
		return
	}

	// Retrieve the value for verification.
	result, err := op.pool.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Failed to get key %s: %v", key, err)
		return
	}

	log.Printf("Successfully set and retrieved key %s with value: %s", key, result)
}

// adjustConcurrency adjusts the number of active goroutines based on CPU usage.
func (op *Operation) adjustConcurrency() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Monitor CPU usage.
			cpuUsage, _ := cpu.Percent(time.Second, false)
			log.Printf("Current CPU usage: %.2f%%", cpuUsage)

			// Adjust concurrency based on CPU usage.
			if cpuUsage < 70.0 {
				// Increase concurrency if CPU usage is low.
				op.maxConcurrency += 5
				log.Printf("Increasing concurrency to %d", op.maxConcurrency)
			} else if cpuUsage > 90.0 {
				// Decrease concurrency if CPU usage is high.
				op.maxConcurrency -= 5
				if op.maxConcurrency < 1 {
					op.maxConcurrency = 1
				}
				log.Printf("Decreasing concurrency to %d", op.maxConcurrency)
			}

			// Limit the number of workers.
			if cap(op.workChan) < op.maxConcurrency {
				op.workChan = make(chan int, op.maxConcurrency)
			}

		case <-op.workChan:
		}
	}
}

func main() {
	const (
		redisAddr  = "localhost:6379"
		initialWorkers = 10
		numTasks   = 1000
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

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	op := &Operation{
		pool:         pool,
		wg:           &wg,
		workChan:     make(chan int, initialWorkers),
		maxConcurrency: initialWorkers,
	}

	// Start goroutine to adjust concurrency.
	go op.adjustConcurrency()

	// Run tasks concurrently.
	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		select {
		case <-op.workChan:
			go op.performTask(ctx, fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i))
		case <-ctx.Done():
			log.Println("Context canceled, stopping tasks.")
			return
		}
	}

	// Wait for all tasks to complete.
	wg.Wait()

	log.Println("All tasks completed successfully.")
}
