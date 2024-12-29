
package main
import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// BaseError represents an application-level error with a message.
type BaseError struct {
	Message string
}

func (e BaseError) Error() string {
	return e.Message
}

// ValidationError represents a validation error.
type ValidationError struct {
	BaseError
}

// AuthenticationError represents an authentication error.
type AuthenticationError struct {
	BaseError
}

// InternalServerError represents an internal server error.
type InternalServerError struct {
	BaseError
}

// Middleware is a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// EventDispatcher holds the registered handlers and error handlers for different error types.
type EventDispatcher struct {
	routes         map[string]http.Handler
	errorHandlers map[interface{}]func(w http.ResponseWriter, r *http.Request, err error)
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		routes:         make(map[string]http.Handler),
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
		d.handleError(w, r, ValidationError{BaseError{Message: "Not found"}}, http.StatusNotFound)
		return
	}
	handler.ServeHTTP(w, r)
}

func (d *EventDispatcher) handleError(w http.ResponseWriter, r *http.Request, err error, code int) {
	logError(err)
	errorHandler, ok := d.errorHandlers[typeOf(err)]
	if !ok {
		errorHandler = d.errorHandlers[nil]
	}
	errorHandler(w, r, err)
	http.Error(w, err.Error(), code)
}

func typeOf(v interface{}) interface{} {
	return fmt.Sprintf("%T", v)
}

// Logger is a middleware function that logs each request and error.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
			}
		}()
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// errorHandler logs error details (type, message, and stack trace) to a log file.
func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	stack := make([]byte, 1<<16)
	n := runtime.Stack(stack, false)
	log.Printf("Error: %s\nStack trace:\n%s", err, stack[:n])
}

// Authenticator is a middleware function for simple authentication.