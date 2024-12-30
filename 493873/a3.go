package main

import (
	"fmt"
	"sync"
	"time"
)

type Data struct {
	Chunk []byte
}

type ObjectPool struct {
	pool sync.Pool
}

func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Data{Chunk: make([]byte, 1024)} // Adjust size based on expected chunk size
			},
		},
	}
}

func (mp *ObjectPool) Get() *Data {
	return mp.pool.Get().(*Data)
}

func (mp *ObjectPool) Put(data *Data) {
	data.Chunk = data.Chunk[:0] // Clear the data slice
	mp.pool.Put(data)
}

type Worker struct {
	pool    *ObjectPool
	channel chan *Data
}

func NewWorker(pool *ObjectPool) *Worker {
	return &Worker{pool: pool, channel: make(chan *Data)}
}

func (w *Worker) Start() {
	go func() {
		for data := range w.channel {
			encode(data)
			decode(data)
			fmt.Printf("Processed Data: %s\n", string(data.Chunk))
			w.pool.Put(data)
		}
	}()
}

func (w *Worker) Stop() {
	close(w.channel)
}

func encode(data *Data) {
	// Simulate encoding logic
	data.Chunk = append(data.Chunk, []byte("_encoded")...)
}

func decode(data *Data) {
	// Simulate decoding logic
	data.Chunk = data.Chunk[:len(data.Chunk)-8]
}

func processData(workers []*Worker, data *Data) {
	// Send data to a random worker for processing
	randomWorker := workers[rand.Intn(len(workers))]
	randomWorker.channel <- data
}

func benchmark(workers []*Worker, iterations int) time.Duration {
	start := time.Now()

	for i := 0; i < iterations; i++ {
		data := workers[0].pool.Get()
		data.Chunk = []byte("raw data")
		processData(workers, data)
	}

	// Wait for all workers to finish
	for _, worker := range workers {
		worker.Stop()
	}
	time.Sleep(100 * time.Millisecond) // Allow workers to finish processing

	return time.Since(start)
}

func main() {
	objectPool := NewObjectPool()

	// Create a bounded worker pool
	const numWorkers = 4
	workers := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = NewWorker(objectPool)
		workers[i].Start()
	}

	// Benchmark processing with a bounded worker pool
	fmt.Println("Benchmarking with a bounded worker pool...")
	workerPoolDuration := benchmark(workers, 10000)
	fmt.Printf("Time taken (worker pool): %v\n", workerPoolDuration)

	// Clean up workers
	for _, worker := range workers {
		worker.Stop()
	}
}