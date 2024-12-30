package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Number of worker goroutines
	workerPoolSize = 4

	// Buffer size for worker queue
	workerQueueSize = 100
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
				return &Data{Chunk: make([]byte, 1024)}
			},
		},
	}
}

func (mp *MemoryPool) Get() *Data {
	return mp.pool.Get().(*Data)
}

func (mp *MemoryPool) Put(data *Data) {
	data.Chunk = data.Chunk[:0]
	mp.pool.Put(data)
}

type ObjectPool struct {
	pool sync.Pool
}

func NewObjectPool() *ObjectPool {
	return &ObjectPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &DataProcessor{}
			},
		},
	}
}

func (op *ObjectPool) Get() *DataProcessor {
	return op.pool.Get().(*DataProcessor)
}

func (op *ObjectPool) Put(dp *DataProcessor) {
	dp.Reset()
	op.pool.Put(dp)
}

type DataProcessor struct {
	memoryPool *MemoryPool
	data       *Data
}

func (dp *DataProcessor) Init(memoryPool *MemoryPool) {
	dp.memoryPool = memoryPool
	dp.data = memoryPool.Get()
}

func (dp *DataProcessor) ProcessData() {
	// Simulate encoding logic
	dp.data.Chunk = append(dp.data.Chunk, []byte("_encoded")...)

	// Simulate decoding logic
	dp.data.Chunk = dp.data.Chunk[:len(dp.data.Chunk)-8]
}

func (dp *DataProcessor) Reset() {
	dp.memoryPool.Put(dp.data)
	dp.data = nil
}

func processDataWorker(workerID int, workerQueue <-chan *DataProcessor, wg *sync.WaitGroup) {
	defer wg.Done()
	for dp := range workerQueue {
		dp.ProcessData()
		fmt.Printf("Worker %d processed data: %s\n", workerID, string(dp.data.Chunk))
	}
}

func startWorkers(memoryPool *MemoryPool, objectPool *ObjectPool, wg *sync.WaitGroup) chan<- *DataProcessor {
	workerQueue := make(chan *DataProcessor, workerQueueSize)

	for i := 0; i < workerPoolSize; i++ {
		dp := objectPool.Get()
		dp.Init(memoryPool)

		wg.Add(1)
		go processDataWorker(i+1, workerQueue, wg)
	}

	return workerQueue
}

func benchmark(memoryPool *MemoryPool, objectPool *ObjectPool, iterations int) time.Duration {
	start := time.Now()

	var wg sync.WaitGroup
	workerQueue := startWorkers(memoryPool, objectPool, &wg)

	for i := 0; i < iterations; i++ {
		dp := objectPool.Get()
		dp.Init(memoryPool)
		dp.data.Chunk = []byte("raw data")
		workerQueue <- dp
	}

	close(workerQueue)

	wg.Wait()
	return time.Since(start)
}

func main() {
	memoryPool := NewMemoryPool()
	objectPool := NewObjectPool()

	// Benchmark using bounded worker pool
	fmt.Println("Benchmarking with bounded worker pool...")
	boundedPoolDuration := benchmark(memoryPool, objectPool, 1000)