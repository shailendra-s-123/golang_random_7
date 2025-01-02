package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Event represents an event that triggers an API endpoint
type Event struct {
	Path string
	R    *http.Request
}

// EndpointCallback is a function type that represents an API endpoint callback
type EndpointCallback func(http.ResponseWriter, *http.Request)

var endpointMap map[string]EndpointCallback
var eventChans map[string]chan Event
var mu sync.Mutex

const eventChannelBufferSize = 100

func init() {
	endpointMap = make(map[string]EndpointCallback)
	eventChans = make(map[string]chan Event)
}

// RegisterEndpoint registers an API endpoint with a callback and starts a Goroutine to handle events asynchronously
func RegisterEndpoint(path string, callback EndpointCallback) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := endpointMap[path]; exists {
		log.Printf("Endpoint %s already registered.", path)
		return
	}
	endpointMap[path] = callback
	eventChans[path] = make(chan Event, eventChannelBufferSize)

	go handleEventsAsync(path, callback)
}

// handleEventsAsync processes events asynchronously for the given endpoint
func handleEventsAsync(path string, callback EndpointCallback) {
	for event := range eventChans[path] {
		// Process the event asynchronously
		go func(e Event) {
			w := &dummyResponseWriter{}
			callback(w, e.R)
			log.Printf("Event processed for endpoint %s: %s", path, e.R.URL.Path)
		}(event)
	}
}

// HandlerFunction handles incoming HTTP requests and triggers the corresponding event channel
func HandlerFunction(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	callback, exists := endpointMap[r.URL.Path]
	if exists {
		// Send the event to the corresponding event channel
		eventChans[r.URL.Path] <- Event{Path: r.URL.Path, R: r}
		w.WriteHeader(http.StatusAccepted)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// dummyResponseWriter is a simple struct to implement http.ResponseWriter for dummy requests
type dummyResponseWriter struct {
}

func (d *dummyResponseWriter) Header() http.Header {
	return nil
}

func (d *dummyResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (d *dummyResponseWriter) WriteHeader(code int) {
}

func main() {
	// Register sample endpoints
	RegisterEndpoint("/hello", HelloHandler)
	RegisterEndpoint("/goodbye", GoodbyeHandler)

	// Start the HTTP server
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(HandlerFunction))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
 
// Rest of the handlers remain the same...