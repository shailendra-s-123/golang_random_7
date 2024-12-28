package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// Cache to store request counts
var cache = make(map[string]*counter, 100)
var cacheLock sync.RWMutex

// Counter structure to hold request counts and reset time
type counter struct {
	count  int
	reset  time.Time
	mutex  sync.Mutex
}

// RateLimitMiddleware middleware to enforce rate limiting
func RateLimitMiddleware(limit int, period time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the API key from the query parameter
			apiKey := r.URL.Query().Get("apiKey")

			// Check if the API key is present
			if apiKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			// Lock the cache for safe access
			cacheLock.RLock()
			c, exists := cache[apiKey]
			cacheLock.RUnlock()

			// Initialize counter if not exists
			if !exists {
				cacheLock.Lock()
				defer cacheLock.Unlock()
				c = &counter{
					count:  0,
					reset:  time.Now(),
					mutex:  sync.Mutex{},
				}
				cache[apiKey] = c
			}

			// Increment request count
			c.mutex.Lock()
			c.count++
			c.mutex.Unlock()

			// Check against rate limit
			now := time.Now()
			if now.Before(c.reset) {
				if c.count > limit {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
					return
				}
			} else {
				// Reset counter if the period is over
				cacheLock.Lock()
				c.reset = now.Add(period)
				cacheLock.Unlock()
			}

			// Proceed with the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	// Define the rate limit and period
	limit := 10
	period := time.Minute

	// Create a handler for your API endpoint
	http.HandleFunc("/api/resource", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	// Apply rate limit middleware
	http.Handle("/api/resource", RateLimitMiddleware(limit, period)(http.DefaultServeMux))

	// Start the HTTP server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}