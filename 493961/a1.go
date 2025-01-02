package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
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

func main() {
	// Register a sample endpoint callback
	RegisterEndpoint("/hello", HelloHandler)
	RegisterEndpoint("/goodbye", GoodbyeHandler)

	// Start the HTTP server
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(HandlerFunction))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// HelloHandler is a callback function for the "/hello" endpoint
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
}

// GoodbyeHandler is a callback function for the "/goodbye" endpoint
func GoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goodbye, %s!", r.URL.Query().Get("name"))
}