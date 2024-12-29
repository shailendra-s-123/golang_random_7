package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Middleware is a function that takes another function as an argument and returns a function.
type Middleware func(http.Handler) http.Handler

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

// RegisterCallback registers a new callback for a given path, optionally with middleware.
func (d *EventDispatcher) RegisterCallback(path string, cb Callback, middlewares ...Middleware) {
	h := cb
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	d.routes[path] = h
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

// LoggingMiddleware logs each request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// AuthenticationMiddleware checks for an API key in the query parameters.
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.URL.Query().Get("api_key")
		if apiKey != "secret_api_key" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ErrorHandlingMiddleware catches any error and returns a 500 status code.
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[ERROR] %+v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Main function to set up the server and handle events.
func main() {
	// Create a new event dispatcher
	dispatcher := NewEventDispatcher()

	// Register some callback functions for specific paths with middleware
	dispatcher.RegisterCallback("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Hello, %s!", name)
	}, LoggingMiddleware, AuthenticationMiddleware, ErrorHandlingMiddleware)

	dispatcher.RegisterCallback("/goodbye", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Goodbye, %s!", name)
	}, LoggingMiddleware, ErrorHandlingMiddleware)

	// Start the HTTP server
	http.HandleFunc("/", dispatcher.HandleEvent)
	log.Fatal(http.ListenAndServe(":8080", nil))
}