package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// pool represents a connection pool to a Redis server.
type pool struct {
	rdb   *redis.Client
	mu    sync.Mutex
	ctx   context.Context
	cancel context.CancelFunc
}

// NewPool creates a new connection pool to the Redis server with concurrency control.
func NewPool(addr string, concurrency int) (*pool, error) {
	if concurrency <= 0 {
		return nil, errors.New("concurrency should be a positive integer")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // use default Addr
		Password: "",   // no password set
		DB:       0,   // use default DB
		PoolSize: concurrency,
	})

	ctx, cancel := context.WithCancel(context.Background())
	err := rdb.Ping(ctx).Err()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &pool{rdb: rdb, ctx: ctx, cancel: cancel}, nil
}

// Get gets a new connection from the pool.
func (p *pool) Get() redis.Cmdable {
	// Let's introduce a simulated random delay to showcase thread safety and concurrency
	delay := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(delay)

	p.mu.Lock()
	defer p.mu.Unlock()
	return p.rdb.WithContext(p.ctx)
}

// Close closes the connection pool.
func (p *pool) Close() error {
	p.cancel() // Cancel the context to gracefully shutdown the pool
	return p.rdb.Close()
}

// worker performs tasks using a Redis connection from the pool.
func worker(wg *sync.WaitGroup, pool *pool, id int) {
	defer wg.Done()

	redisConn := pool.Get()
	defer redisConn.Close() // Properly close the connection when the goroutine exits

	// Simulate some work with Redis
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("worker_%d_%d", id, i)
		_, err := redisConn.Set(pool.ctx, key, "value", 0).Result()
		if err != nil {
			log.Printf("Error setting key %q: %v", key, err)
			return
		}
	}
}

// Example usage:
func main() {
	const (
		numWorkers  = 10
		concurrency = 10 // Redis pool size should be at least equal to concurrency
	)

	// Connect to Redis using our connection pool with concurrency control
	pool, err := NewPool("localhost:6379", concurrency)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}
	defer pool.Close() // Clean up pool on program exit

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start workers in parallel
	for i := 0; i < numWorkers; i++ {
		go worker(&wg, pool, i)
	}

	wg.Wait()
	fmt.Println("All workers completed.")
}
