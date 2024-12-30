package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type pool struct {
	rdb *redis.Client
}

// NewPool creates a new connection pool to the Redis server.
func NewPool(addr string) (*pool, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // use default Addr
		Password: "",   // no password set
		DB:       0,   // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return &pool{rdb: rdb}, nil
}

// Get gets a new connection from the pool.
func (p *pool) Get() redis.Cmdable {
	return p.rdb.WithContext(context.Background())
}

// Close closes the connection pool.
func (p *pool) Close() error {
	return p.rdb.Close()
}

type operation struct {
	pool *pool
	key  string
	value string
	ctx    context.Context
}

func (op *operation) execute() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered: %v", r)
		}
	}()

	ctx, cancel := context.WithTimeout(op.ctx, 1*time.Second)
	defer cancel()

	redisConn := op.pool.Get()
	defer redisConn.Close()

	_, err := redisConn.Set(ctx, op.key, op.value, 0).Result()
	if err != nil {
		log.Printf("Error setting key %s: %v", op.key, err)
		return
	}

	val, err := redisConn.Get(ctx, op.key).Result()
	if err != nil {
		log.Printf("Error getting key %s: %v", op.key, err)
		return
	}

	log.Printf("Value for key %s is %s", op.key, val)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := NewPool("localhost:6379")
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	defer pool.Close() // Ensure the pool is closed

	var wg sync.WaitGroup
	var numOperations int32 = 100

	for i := 0; i < int(numOperations); i++ {
		op := &operation{
			pool: pool,
			key:  fmt.Sprintf("key_%s", uuid.New()),
			value: fmt.Sprintf("value_%s", uuid.New()),
			ctx: ctx,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			op.execute()
		}()
	}

	wg.Wait()
	log.Printf("Completed %d operations", numOperations)
}