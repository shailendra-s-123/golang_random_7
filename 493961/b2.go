package main  
import (  
    "fmt"
    "log"
    "net/http"
    "strings"
)

// Middleware is a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// HandlerFunc is a function that takes an http.ResponseWriter and *http.Request and responds to HTTP requests.
type HandlerFunc func(http.ResponseWriter, *http.Request)

// EventDispatcher holds the registered handlers and dispatches events.
type EventDispatcher struct {  
    handlers map[string]HandlerFunc
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {  
    return &EventDispatcher{  
        handlers: make(map[string]HandlerFunc),
    }
}

// RegisterHandler registers a new handler for a given path.
func (d *EventDispatcher) RegisterHandler(path string, h HandlerFunc) {  
    d.handlers[path] = h
}

// HandleEvent dispatches an event (HTTP request) to the registered handler, applying the specified middleware.
func (d *EventDispatcher) HandleEvent(w http.ResponseWriter, r *http.Request) {  
    path := r.URL.Path
    h, ok := d.handlers[path]
    if !ok {  
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }
    chain := newHandlerChain(h)
    for _, m := range []Middleware{Logger, Authenticator, ErrorHandler} {
        chain = m(chain)
    }
    chain.ServeHTTP(w, r)
}

// newHandlerChain creates a new http.Handler from the given HandlerFunc.
func newHandlerChain(h HandlerFunc) http.Handler {  
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  
        h(w, r)
    })
}

// Logger is a middleware function for logging requests.
func Logger(next http.Handler) http.Handler {  
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  
        log.Printf("%s %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// Authenticator is a middleware function for basic authentication.
func Authenticator(next http.Handler) http.Handler {  
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Basic ") {  
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// ErrorHandler is a middleware function for handling errors.
func ErrorHandler(next http.Handler) http.Handler {  
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  
        defer func() {  
            if err := recover(); err != nil {  
                log.Printf("Error: %v\n", err)
                http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// Main function to set up the server and handle events.
func main() {  
    dispatcher := NewEventDispatcher()

    dispatcher.RegisterHandler("/hello", func(w http.ResponseWriter, r *http.Request) {  
        fmt.Fprintf(w, "Hello, %s!", r.URL.Query().Get("name"))
    })

    http.HandleFunc("/", dispatcher.HandleEvent)
    log.Fatal(http.ListenAndServe(":8080", nil))
}