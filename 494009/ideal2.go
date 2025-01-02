//main.go
package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestConcurrentRedisOperations tests the concurrent Redis operations.
func TestConcurrentRedisOperations(t *testing.T) {
	// Create a context with timeout to handle cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	// Number of concurrent goroutines to test
	numGoroutines := 10
	wg.Add(numGoroutines)

	// Start concurrent operations
	for i := 0; i < numGoroutines; i++ {
		go concurrentRedisOperations(ctx, &wg, numGoroutines)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	t.Log("All operations completed successfully.")
}

// BenchmarkRedisOperation benchmarks Redis operations.
func BenchmarkRedisOperation(b *testing.B) {
	client := createRedisClient()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Set(ctx, "foo", "bar", 10*time.Second).Err()
		client.Get(ctx, "foo").Result()
	}
}
