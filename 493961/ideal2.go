
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Middleware is a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// EventDispatcher holds the registered handlers and dispatches events.
type EventDispatcher struct {
	routes map[string]http.Handler
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		routes: make(map[string]http.Handler),
	}
}

// RegisterEndpoint registers a new endpoint with optional middleware.
func (d *EventDispatcher) RegisterEndpoint(path string, handler http.Handler, middlewares ...Middleware) {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	d.routes[path] = handler
}

// HandleEvent dispatches the event (HTTP request) to the registered handler.
func (d *EventDispatcher) HandleEvent(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	handler, ok := d.routes[path]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	handler.ServeHTTP(w, r)
}

// Logger is a middleware function that logs each request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Authenticator is a middleware function for simple authentication.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") || authHeader[len("Bearer "):] != "secret_token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ErrorHandler is a middleware function that handles errors.
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Error: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// HelloHandler is an example handler for the "/hello" endpoint.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

// GoodbyeHandler is an example handler for the "/goodbye" endpoint.
func GoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Goodbye, %s!", name)
}

func main() {
	// Create a new event dispatcher
	dispatcher := NewEventDispatcher()

	// Register the "/hello" and "/goodbye" endpoints with middleware
	dispatcher.RegisterEndpoint("/hello", http.HandlerFunc(HelloHandler), Logger, Authenticator, ErrorHandler)
	dispatcher.RegisterEndpoint("/goodbye", http.HandlerFunc(GoodbyeHandler), Logger, ErrorHandler)

	// Set up the HTTP server
	http.HandleFunc("/", dispatcher.HandleEvent)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


