package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Global variable for Redis client pool
var redisClient *redis.Client

// initializeRedisClient initializes the Redis client pool.
func initializeRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Check connection status
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

// cleanupRedisClient closes the Redis client pool.
func cleanupRedisClient() {
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis client: %v", err)
	}
}

// worker function that performs Redis operations concurrently.
func worker(id int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	// Acquire a lock to ensure mutual exclusion while using Redis client
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	// Use the defer keyword for resource cleanup within the worker
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Worker %d: Error closing Redis client: %v", id, err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Cancellation received. Exiting...\n", id)
			return
		default:
			// Simulate work
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

			// Perform Redis operations
			key := fmt.Sprintf("worker_%d_key", id)
			err := redisClient.Set(ctx, key, fmt.Sprintf("Worker %d's value", id), 10*time.Second).Err()
			if err != nil {
				log.Printf("Worker %d: Error setting value in Redis: %v", id, err)
				continue
			}

			val, err := redisClient.Get(ctx, key).Result()
			if err != nil {
				log.Printf("Worker %d: Error getting value from Redis: %v", id, err)
				continue
			}
			fmt.Printf("Worker %d: Got value from Redis: %s\n", id, val)
		}
	}
}

func main() {
	initializeRedisClient()
	defer cleanupRedisClient()

	// Number of workers to spawn
	numWorkers := 10
	var wg sync.WaitGroup

	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure context is canceled upon return

	// Spawn workers
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(i, &wg, ctx)
	}

	// Wait for all workers to complete
	wg.Wait()
	fmt.Println("All workers finished. Program exiting...")
} 