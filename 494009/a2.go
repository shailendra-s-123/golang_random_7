package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// createRedisClient creates a Redis client with the given parameters and returns it.
func createRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Default DB
	})

	// Ensure Redis server is reachable
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}

// cleanup handles the cleanup of resources, ensuring safe deallocation.
func cleanup(client *redis.Client) {
	// Using defer to ensure that Redis connection is closed after all operations
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}()

	// Simulating additional cleanup tasks...
	fmt.Println("Performing additional cleanup tasks...")
}

// performRedisOperation performs a single Redis operation with error handling.
func performRedisOperation(ctx context.Context, client *redis.Client) {
	// Using defer to close connection in case of error
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Operation canceled due to context")
		return
	default:
	}

	err := client.Set(ctx, "foo", "bar", 10*time.Second).Err()
	if err != nil {
		log.Printf("Error setting value in Redis: %v", err)
		return
	}
	fmt.Println("Set 'foo' to 'bar' in Redis.")

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		log.Printf("Error getting value from Redis: %v", err)
		return
	}
	fmt.Printf("Got value from Redis: %s\n", val)
}

// concurrentRedisOperations performs Redis operations concurrently.
func concurrentRedisOperations(ctx context.Context, wg *sync.WaitGroup, numGoroutines int) {
	defer wg.Done()
	
	for i := 0; i < numGoroutines; i++ {
		// Create a new Redis client for each goroutine
		client := createRedisClient()

		go func() {
			// Perform Redis operation
			performRedisOperation(ctx, client)
			// Cleanup resources
			cleanup(client)
		}()
	}
}

// main function that demonstrates concurrent Redis operations with defer
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Cancel the context on completion

	var wg sync.WaitGroup
	// Set number of concurrent goroutines
	numGoroutines := 10
	wg.Add(numGoroutines)

	// Start concurrent operations
	for i := 0; i < numGoroutines; i++ {
		go concurrentRedisOperations(ctx, &wg, numGoroutines)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All operations completed.")
}