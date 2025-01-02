package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// EndpointCallback represents a callback for an API endpoint
type EndpointCallback func(http.ResponseWriter, *http.Request)

// Struct to store endpoint information
type Endpoint struct {
	Path     string
	Callback EndpointCallback
}

// EventHandler processes events asynchronously
func EventHandler(event interface{}) {
	// Simulate some asynchronous work
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Processed event:", event)
}

var endpointMap map[string]EndpointCallback
var mu sync.Mutex

// Worker Pool to manage concurrent processing
type WorkerPool struct {
	Workers     int
	EventQueue  chan interface{}
	wg          sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		Workers:    workers,
		EventQueue: make(chan interface{}, 100),
	}
}

// Start starts the worker pool to process events
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.Workers; i++ {
		go func() {
			for event := range wp.EventQueue {
				wp.wg.Add(1)
				EventHandler(event)
				wp.wg.Done()
			}
		}()
	}
}

// Wait waits for all workers to finish
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

func init() {
	endpointMap = make(map[string]EndpointCallback)
}

func RegisterEndpoint(path string, callback EndpointCallback) {
	mu.Lock()
	defer mu.Unlock()
	endpointMap[path] = callback
}

func HandlerFunction(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	callback, exists := endpointMap[r.URL.Path]
	if exists {
		callback(w, r)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "Hello, %s!", name)
	event := fmt.Sprintf("Hello event for %s", name)
	// Add event to the worker pool
	workerPool.EventQueue <- event
}

func GoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "Goodbye, %s!", name)
	event := fmt.Sprintf("Goodbye event for %s", name)
	// Add event to the worker pool
	workerPool.EventQueue <- event
}

var workerPool *WorkerPool

func main() {
	// Configure the worker pool size
	workerPool = NewWorkerPool(10)
	workerPool.Start()

	// Register the endpoints
	RegisterEndpoint("/hello", HelloHandler)
	RegisterEndpoint("/goodbye", GoodbyeHandler)

	// Start the HTTP server
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(HandlerFunction))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	// Wait for all events to be processed before shutting down
	workerPool.Wait()
}
