package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"runtime"
	"os"
	"io"
)

type CustomError struct {
	error
	Type    string
	Message string
	Stack   string
}

func newCustomError(msg string, t string, err error) *CustomError {
	// Generate the stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, true)
	stack := string(buf[:n])

	return &CustomError{
		error:  err,
		Type:   t,
		Message: msg,
		Stack: stack,
	}
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("[%s] %s\n%s", e.Type, e.Message, e.Stack)
}

// Middleware is a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// EventDispatcher holds the registered handlers and dispatches events.
type EventDispatcher struct {
	routes map[string]http.Handler
	errorHandlers map[string]func(error)
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		routes: make(map[string]http.Handler),
		errorHandlers: make(map[string]func(error)),
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
		d.triggerErrorHandler("Not Found", http.StatusNotFound)
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
			d.triggerErrorHandler("Unauthorized", http.StatusUnauthorized)
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
				d.triggerErrorHandler(fmt.Sprintf("Internal Error: %v", err), http.StatusInternalServerError)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (d *EventDispatcher) triggerErrorHandler(msg string, code int) {
	errorType := fmt.Sprintf("HTTP-%d", code)
	err := newCustomError(msg, errorType, nil)
	if handler, ok := d.errorHandlers[errorType]; ok {
		handler(err)
	}
}

// RegisterErrorHandler allows dynamic registration of error-handling functions.
func (d *EventDispatcher) RegisterErrorHandler(errorType string, handler func(error)) {
	d.errorHandlers[errorType] = handler
}

// HelloHandler is an example handler for the "/hello" endpoint.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		d.triggerErrorHandler("Name is required", http.StatusBadRequest)
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

// GoodbyeHandler is an example handler for the "/goodbye" endpoint.
func GoodbyeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		d.triggerErrorHandler("Name is required", http.StatusBadRequest)
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Goodbye, %s!", name)
}

func logErrorsToFile(err error) {
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Failed to open error log file:", err)
		return
	}
	defer logFile.Close()

	_, err = io.WriteString(logFile, err.Error()+"\n")
	if err != nil {
		log.Println("Failed to write error to log file:", err)
	}
}

func main() {
	// Create a new event dispatcher
	dispatcher := NewEventDispatcher()

	// Register the "/hello" and "/goodbye" endpoints with middleware
	dispatcher.RegisterEndpoint("/hello", http.HandlerFunc(HelloHandler), Logger, Authenticator, ErrorHandler)
	dispatcher.RegisterEndpoint("/goodbye", http.HandlerFunc(GoodbyeHandler), Logger, ErrorHandler)

	// Register custom error handlers
	dispatcher.RegisterErrorHandler("HTTP-400", logErrorsToFile)
	dispatcher.RegisterErrorHandler("HTTP-401", logErrorsToFile)
	dispatcher.RegisterErrorHandler("HTTP-500", logErrorsToFile)

	// Set up the HTTP server
	http.HandleFunc("/", dispatcher.HandleEvent)
	log.Fatal(http.ListenAndServe(":8080", nil))
}