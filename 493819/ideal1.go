package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// userRateLimiter struct to manage rate limiting for users
type userRateLimiter struct {
	limit       int               // Maximum requests allowed
	window      time.Duration     // Time window in which requests are counted
	userLimits  map[string]*userLimit
	userLimitsMu sync.Mutex
}

// userLimit struct to store individual user request data
type userLimit struct {
	lastCheck time.Time // Last time user made a request
	requests  int       // Number of requests in the current window
}

// NewUserRateLimiter creates a new rate limiter for users with specified limit and window size
func newUserRateLimiter(limit int, window time.Duration) *userRateLimiter {
	return &userRateLimiter{
		limit:      limit,
		window:     window,
		userLimits: make(map[string]*userLimit),
	}
}

// allow checks if a user's request is within the rate limit
func (rl *userRateLimiter) allow(userID string) bool {
	rl.userLimitsMu.Lock()
	defer rl.userLimitsMu.Unlock()

	// Retrieve or initialize the user's rate limit data
	limit, ok := rl.userLimits[userID]
	if !ok {
		rl.userLimits[userID] = &userLimit{
			lastCheck: time.Now(),
			requests:  1,
		}
		return true // New user, allow immediately
	}

	// Check the time difference to determine if we need to reset the counter
	now := time.Now()
	if now.Sub(limit.lastCheck) > rl.window {
		// Reset request count if the window has expired
		limit.requests = 1
	} else {
		limit.requests++
	}

	limit.lastCheck = now
	allowed := limit.requests <= rl.limit

	return allowed
}

func main() {
	// Initialize rate limiter with a limit of 5 requests per 10 seconds
	limiter := newUserRateLimiter(5, time.Second*10)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.FormValue("user") // Extract userID from the query parameter
		if userID == "" {
			http.Error(w, "Missing user parameter", http.StatusBadRequest)
			return
		}

		// Check if the user is allowed to make a request
		allowed := limiter.allow(userID)
		if !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		fmt.Fprintf(w, "Hello, %s!", userID) // Successful request
	})

	// Start the HTTP server
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}