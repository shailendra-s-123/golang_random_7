package main

import (
	"fmt"
	"log"
	"net/http"
)

// Callback is a function that will be executed when a specific event occurs.
type Callback func(http.ResponseWriter, *http.Request)

// EventDispatcher holds the registered callbacks and dispatches events.
type EventDispatcher struct {
	routes map[string]Callback
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		routes: make(map[string]Callback),
	}
}

// RegisterCallback registers a new callback for a given path.
func (d *EventDispatcher) RegisterCallback(path string, cb Callback) {
	d.routes[path] = cb
}

// HandleEvent dispatches an event (HTTP request) to the registered callback.
func (d *EventDispatcher) HandleEvent(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	cb, ok := d.routes[path]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	cb(w, r)
}

// Main function to set up the server and handle events.
func main() {
	// Create a new event dispatcher
	dispatcher := NewEventDispatcher()

	// Register some callback functions for specific paths
	dispatcher.RegisterCallback("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
	})

	dispatcher.RegisterCallback("/goodbye", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Goodbye, %s!", r.URL.Query().Get("name"))
	})

	// Start the HTTP server
	http.HandleFunc("/", dispatcher.HandleEvent)
	log.Fatal(http.ListenAndServe(":8080", nil))
}