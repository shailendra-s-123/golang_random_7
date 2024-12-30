package main

import (
	"fmt"
	"sync"
	"time"
)

type Data struct {
	Chunk []byte
}

type MemoryPool struct {
	pool sync.Pool
}

func NewMemoryPool() *MemoryPool {
	return &MemoryPool{
		pool: sync.Pool{
			New: func() interface{} {
				// Allocating memory for data chunk size
				return &Data{Chunk: make([]byte, 1024)} // Adjust size based on expected chunk size
			},
		},
	}
}

func (mp *MemoryPool) Get() *Data {
	return mp.pool.Get().(*Data)
}

func (mp *MemoryPool) Put(data *Data) {
	data.Chunk = data.Chunk[:0] // Clear the data slice
	mp.pool.Put(data)
}

func encode(data *Data) {
	// Simulate encoding logic
	// (for now just append "_encoded" to the data)
	data.Chunk = append(data.Chunk, []byte("_encoded")...)
}

func decode(data *Data) {
	// Simulate decoding logic
	// (for now just strip "_encoded" from the data)
	data.Chunk = data.Chunk[:len(data.Chunk)-8]
}

func processData(pool *MemoryPool, parallel bool, wg *sync.WaitGroup) {
	data := pool.Get()
	defer pool.Put(data)

	// Simulate real-time streaming data
	data.Chunk = []byte("raw data")

	// Perform encoding and decoding
	if parallel {
		var encodeWg sync.WaitGroup
		encodeWg.Add(1)
		go func() {
			defer encodeWg.Done()
			encode(data)
		}()

		var decodeWg sync.WaitGroup
		decodeWg.Add(1)
		go func() {
			defer decodeWg.Done()
			decode(data)
		}()

		encodeWg.Wait()
		decodeWg.Wait()
	} else {
		encode(data)
		decode(data)
	}

	// Print the processed data (in reality, this could be saved or further processed)
	fmt.Printf("Processed Data: %s\n", string(data.Chunk))
}

func benchmark(pool *MemoryPool, parallel bool, iterations int) time.Duration {
	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processData(pool, parallel, &wg)
		}()
	}

	wg.Wait()
	return time.Since(start)
}

func main() {
	pool := NewMemoryPool()

	// Benchmark serial processing (without parallelism)
	fmt.Println("Benchmarking serial processing...")
	serialDuration := benchmark(pool, false, 10)
	fmt.Printf("Time taken (serial): %v\n", serialDuration)

	// Benchmark parallel processing
	fmt.Println("\nBenchmarking parallel processing...")
	parallelDuration := benchmark(pool, true, 10)
	fmt.Printf("Time taken (parallel): %v\n", parallelDuration)

	// You can adjust the number of iterations based on the load you want to simulate
}