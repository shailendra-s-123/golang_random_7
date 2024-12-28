package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// SlidingWindow struct to manage rate limiting with a sliding window
type SlidingWindow struct {
	limit    int
	window   time.Duration
	requests []time.Time
	mu       sync.Mutex
}

// NewSlidingWindow creates a new sliding window rate limiter for users with a specified limit and window size
func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		limit:  limit,
		window: window,
	}
}

// Allow checks if a user's request is within the rate limit using a sliding window
func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()

	// Remove requests that are outside the window
	var validRequests []time.Time
	for _, reqTime := range sw.requests {
		if now.Sub(reqTime) <= sw.window {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Update the list of requests to only include those within the window
	sw.requests = validRequests

	// If the number of requests exceeds the limit, deny the request
	if len(sw.requests) >= sw.limit {
		return false
	}

	// Record the current request
	sw.requests = append(sw.requests, now)
	return true
}

type RateLimiter struct {
	userLimits map[string]*SlidingWindow
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter for tracking different user limits
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		userLimits: make(map[string]*SlidingWindow),
	}
}

// Allow checks if a user's request is allowed, considering their individual limit and window
func (rl *RateLimiter) Allow(userID string, limit int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// If the user doesn't have a rate limiter, create one
	if _, exists := rl.userLimits[userID]; !exists {
		rl.userLimits[userID] = NewSlidingWindow(limit, window)
	}

	// Check if the request is allowed based on the sliding window
	return rl.userLimits[userID].Allow()
}

func main() {
	limiter := NewRateLimiter()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.FormValue("user")
		limitStr := r.FormValue("limit")
		windowStr := r.FormValue("window")

		if userID == "" || limitStr == "" || windowStr == "" {
			http.Error(w, "Missing required parameters: user, limit, window", http.StatusBadRequest)
			return
		}

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}

		window, err := time.ParseDuration(windowStr)
		if err != nil {
			http.Error(w, "Invalid window parameter", http.StatusBadRequest)
			return
		}

		// Check if the request is allowed for this user
		if !limiter.Allow(userID, limit, window) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Successful request
		fmt.Fprintf(w, "Hello, %s! You have %d requests left in the %v window.\n", userID, limit-len(limiter.userLimits[userID].requests), window)
	})

	// Start the HTTP server
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}