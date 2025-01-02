package main

import (
	"context"
	"fmt"
	"log"
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

	// Simulating a cleanup operation that might happen after using resources
	// For example: clearing temporary files, closing open file handles, etc.
	fmt.Println("Performing additional cleanup tasks...")
}

// performRedisOperations demonstrates using the defer keyword for Redis operations.
func performRedisOperations() {
	client := createRedisClient()
	defer cleanup(client) // Ensures Redis client is closed once all operations are done

	// Using the Redis client to set and get values
	err := client.Set(context.Background(), "foo", "bar", 10*time.Second).Err()
	if err != nil {
		log.Fatalf("Error setting value in Redis: %v", err)
	}
	fmt.Println("Set 'foo' to 'bar' in Redis.")

	val, err := client.Get(context.Background(), "foo").Result()
	if err != nil {
		log.Fatalf("Error getting value from Redis: %v", err)
	}
	fmt.Printf("Got value from Redis: %s\n", val)
}

// main function that demonstrates the efficient use of defer for resource management
func main() {
	// Performing Redis operations with proper defer usage
	performRedisOperations()

	// Additional tasks can be added here, for example, working with other third-party resources
	// and ensuring proper cleanup through defer.
}
