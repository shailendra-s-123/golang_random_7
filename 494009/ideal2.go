
package main

import (
	"context"
	"fmt"
	"log"
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
	pool *RedisPool
	wg   *sync.WaitGroup
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

func main() {
	const (
		redisAddr  = "localhost:6379"
		numWorkers = 10
		numTasks   = 100
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

	op := &Operation{pool: pool, wg: &wg}

	// Run tasks concurrently.
	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go op.performTask(ctx, fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i))
	}

	// Wait for all tasks to complete.
	wg.Wait()

	log.Println("All tasks completed successfully.")
}

