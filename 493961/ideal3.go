
package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// BaseError represents a generic application error with a message and code.
type BaseError struct {
	Code    int    // HTTP status code
	Message string // Error message
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

// ValidationError represents a validation error.
type ValidationError struct {
	BaseError
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{BaseError{Code: http.StatusBadRequest, Message: message}}
}

// AuthenticationError represents an authentication error.
type AuthenticationError struct {
	BaseError
}

func NewAuthenticationError(message string) *AuthenticationError {
	return &AuthenticationError{BaseError{Code: http.StatusUnauthorized, Message: message}}
}

// InternalServerError represents an internal server error.
type InternalServerError struct {
	BaseError
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{BaseError{Code: http.StatusInternalServerError, Message: message}}
}

// Middleware type that accepts and returns an HTTP handler.
type Middleware func(http.Handler) http.Handler

// EventDispatcher holds the registered handlers and error handlers for different error types.
type EventDispatcher struct {
	routes        map[string]http.Handler
	errorHandlers map[interface{}]func(w http.ResponseWriter, r *http.Request, err error)
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		routes:        make(map[string]http.Handler),
		errorHandlers: make(map[interface{}]func(w http.ResponseWriter, r *http.Request, err error)),
	}
}

// RegisterEndpoint registers a new endpoint with optional middleware.
func (d *EventDispatcher) RegisterEndpoint(path string, handler http.Handler, middlewares ...Middleware) {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	d.routes[path] = handler
}

// RegisterErrorHandler registers an error handler for a specific error type.
func (d *EventDispatcher) RegisterErrorHandler(errType interface{}, handler func(w http.ResponseWriter, r *http.Request, err error)) {
	d.errorHandlers[errType] = handler
}

// HandleEvent dispatches the event (HTTP request) to the registered handler.
func (d *EventDispatcher) HandleEvent(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	handler, ok := d.routes[path]
	if !ok {
		d.handleError(w, r, NewValidationError("Not found"))
		return
	}
	handler.ServeHTTP(w, r)
}

// handleError handles the error and invokes the appropriate error handler.
func (d *EventDispatcher) handleError(w http.ResponseWriter, r *http.Request, err error) {
	logError(err)
	errorHandler, ok := d.errorHandlers[err]
	if !ok {
		errorHandler = d.errorHandlers[nil]
	}
	errorHandler(w, r, err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Logger is a middleware function that logs each request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// ErrorHandler logs error details (type, message, and stack trace) to a log file.
func logError(err error) {
	stack := make([]byte, 1<<16)
	n := runtime.Stack(stack, false)
	log.Printf("Error: %v\nStack trace:\n%s", err, stack[:n])
}

// DefaultErrorHandler handles errors when no specific handler is registered.
func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	logError(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Authenticator is a middleware function for simple authentication.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer validtoken" {
			// Call handleError on the EventDispatcher instance that is passed through the middleware
			dispatcher.handleError(w, r, NewAuthenticationError("Invalid token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Initialize global dispatcher
var dispatcher = NewEventDispatcher()

func main() {
	// Registering routes and middleware
	dispatcher.RegisterEndpoint("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Login successful"))
	}), Logger)

	// Registering error handlers for different error types
	dispatcher.RegisterErrorHandler((*ValidationError)(nil), func(w http.ResponseWriter, r *http.Request, err error) {
		logError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	})
	dispatcher.RegisterErrorHandler((*AuthenticationError)(nil), func(w http.ResponseWriter, r *http.Request, err error) {
		logError(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
	})
	dispatcher.RegisterErrorHandler((*InternalServerError)(nil), func(w http.ResponseWriter, r *http.Request, err error) {
		logError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	})

	// Handle events (requests)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dispatcher.HandleEvent(w, r)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}