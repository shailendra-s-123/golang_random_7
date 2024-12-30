package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

// pool represents a connection pool to a Redis server.
type pool struct {
	rdb *redis.Client
	mu  sync.Mutex
}

// NewPool creates a new connection pool to the Redis server.
func NewPool(addr string) (*pool, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // use default Addr
		Password: "",   // no password set
		DB:       0,   // use default DB
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return &pool{rdb: rdb}, nil
}

// Get gets a new connection from the pool.
func (p *pool) Get() redis.Cmdable {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.rdb.WithContext(context.Background())
}

// Close closes the connection pool.
func (p *pool) Close() error {
	return p.rdb.Close()
}

// Example usage:
func main() {
	// Connect to Redis using our connection pool
	pool, err := NewPool("localhost:6379")
	if err != nil {
		fmt.Printf("Error connecting to Redis: %v\n", err)
		return
	}
	defer pool.Close() // Clean up pool on program exit

    // No need for multiple defer statements, single Get call is enough
	redisConn := pool.Get()

	// Now, you can use the redisConn for your operations
	_, err = redisConn.Set(context.Background(), "key", "value", 0).Result()
	if err != nil {
		fmt.Printf("Error setting value: %v\n", err)
		return
	}

	val, err := redisConn.Get(context.Background(), "key").Result()
	if err != nil {
		fmt.Printf("Error getting value: %v\n", err)
		return
	}

	fmt.Printf("Value: %s\n", val)
}