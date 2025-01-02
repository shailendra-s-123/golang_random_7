package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// EndpointCallback is a function type that represents an API endpoint callback
type EndpointCallback func(http.ResponseWriter, *http.Request)

var endpointMap map[string]EndpointCallback
var mu sync.Mutex

func init() {
	endpointMap = make(map[string]EndpointCallback)
}

// RegisterEndpoint registers an API endpoint with a callback
func RegisterEndpoint(path string, callback EndpointCallback) {
	mu.Lock()
	defer mu.Unlock()
	endpointMap[path] = callback
}

// HandlerFunction handles incoming HTTP requests and triggers the corresponding endpoint callback
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

// EventHandler processes an event asynchronously
func EventHandler(event interface{}) {
	atomic.AddInt32(&eventsProcessed, 1)
	time.Sleep(time.Duration(atomic.LoadInt32(&eventDelay)) * time.Millisecond)
	fmt.Println("Processed event:", event)
}

var (
	events       chan interface{}
	workerPool   *sync.WorkerPool
	eventsProcessed int32
	eventDelay   int32
)

func main() {
	// Configure worker pool size and event delay
	workerPoolSize := 10
	workerPool = sync.NewWorkerPool(workerPoolSize)
	eventDelay = 50

	// Register sample endpoints
	RegisterEndpoint("/hello", HelloHandler)
	RegisterEndpoint("/goodbye", GoodbyeHandler)

	// Start event worker goroutines
	go func() {
		for event := range events {
			workerPool.Submit(func() { EventHandler(event) })
		}
	}()

	// Start HTTP server
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(HandlerFunction))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	workerPool.Wait()
}

// HelloHandler is a callback function for the "/hello" endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
	events <- r.URL.Query().Get("name")
}

// GoodbyeHandler is a callback function for the "/goodbye" endpoint
func GoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goodbye, %s!", r.URL.Query().Get("name"))
	events <- r.URL.Query().Get("name")
}