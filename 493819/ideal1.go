package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RateLimiter struct {
	maxRequests int
	resetTime   time.Duration
	requests    int
	lastReset   time.Time
	mu          sync.Mutex
}

var userLimits = map[string]*RateLimiter{}
var defaultLimit = 10 // Default max requests per user
var defaultResetTime = 60 * time.Second // Default reset time

// Initialize the rate limiter for a user
func getRateLimiter(userID string) *RateLimiter {
	// Check if user rate limiter exists
	if limiter, exists := userLimits[userID]; exists {
		return limiter
	}

	// Otherwise, create a new rate limiter
	limiter := &RateLimiter{
		maxRequests: defaultLimit,
		resetTime:   defaultResetTime,
		lastReset:   time.Now(),
	}
	userLimits[userID] = limiter
	return limiter
}

// Apply the rate limiting check
func (limiter *RateLimiter) isAllowed() bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	// If reset time has passed, reset the counter
	if time.Since(limiter.lastReset) > limiter.resetTime {
		limiter.requests = 0
		limiter.lastReset = time.Now()
	}

	// Check if the user has exceeded the max requests
	if limiter.requests < limiter.maxRequests {
		limiter.requests++
		return true
	}
	return false
}

// Rate limiter middleware function
func rateLimitHandler(w http.ResponseWriter, r *http.Request) {
	// Get userID from query parameters
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	// Parse the dynamic limit from query parameters if provided
	maxRequestsStr := r.URL.Query().Get("maxRequests")
	if maxRequestsStr != "" {
		maxRequests, err := strconv.Atoi(maxRequestsStr)
		if err == nil && maxRequests > 0 {
			// Update the limit for the user
			limiter := getRateLimiter(userID)
			limiter.maxRequests = maxRequests
		}
	}

	// Get the rate limiter for the user
	limiter := getRateLimiter(userID)

	// Check if the user request is allowed
	if !limiter.isAllowed() {
		http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
		return
	}

	// If allowed, process the request (for demonstration, just return success)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Request for user %s is allowed.\n", userID)
}

func main() {
	http.HandleFunc("/api", rateLimitHandler)

	// Start the server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}