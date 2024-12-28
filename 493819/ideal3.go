package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// DynamicSlidingWindow manages a dynamically adjusting rate limiting window for users.
type DynamicSlidingWindow struct {
	baseLimit    float64         // Base rate limit (requests per second)
	windowSize   time.Duration  // Time window for rate limiting
	requests     []time.Time    // List of request timestamps
	mu           sync.Mutex     // Mutex for concurrency safety
	rateEstimate float64        // Estimated request rate based on EWMA
	alpha        float64        // EWMA smoothing factor
	lastRequest  time.Time      // Last request time
}

// NewDynamicSlidingWindow creates a new dynamic sliding window with the specified parameters.
func NewDynamicSlidingWindow(baseLimit float64, windowSize time.Duration, alpha float64) *DynamicSlidingWindow {
	return &DynamicSlidingWindow{
		baseLimit:  baseLimit,
		windowSize: windowSize,
		alpha:      alpha,
	}
}

// Allow evaluates if a request is allowed based on the current rate and the dynamic adjustment.
func (dsw *DynamicSlidingWindow) Allow() bool {
	dsw.mu.Lock()
	defer dsw.mu.Unlock()

	now := time.Now()
	if dsw.lastRequest.IsZero() {
		dsw.lastRequest = now
		return true
	}

	// Remove requests that fall outside the current window
	var validRequests []time.Time
	for _, reqTime := range dsw.requests {
		if now.Sub(reqTime) <= dsw.windowSize {
			validRequests = append(validRequests, reqTime)
		}
	}
	dsw.requests = validRequests

	// Calculate the current request rate over the time window
	elapsed := now.Sub(dsw.lastRequest)
	currentRate := float64(len(dsw.requests)) / elapsed.Seconds()

	// Update the rate estimate using an exponentially weighted moving average
	dsw.rateEstimate = dsw.alpha*currentRate + (1-dsw.alpha)*dsw.rateEstimate

	// Allow the request if the rate is within the base limit
	if dsw.rateEstimate <= dsw.baseLimit {
		dsw.requests = append(dsw.requests, now)
		dsw.lastRequest = now
		return true
	}

	return false
}

// DynamicRateLimiter manages rate limits for multiple users.
type DynamicRateLimiter struct {
	userLimits map[string]*DynamicSlidingWindow
	mu         sync.Mutex
}

// NewDynamicRateLimiter creates a new rate limiter with a specified base limit and window size.
func NewDynamicRateLimiter(baseLimit float64, windowSize time.Duration, alpha float64) *DynamicRateLimiter {
	return &DynamicRateLimiter{
		userLimits: make(map[string]*DynamicSlidingWindow),
	}
}

// Allow checks if a user's request is allowed based on their individual rate limiter.
func (dr *DynamicRateLimiter) Allow(userID string) bool {
	dr.mu.Lock()
	defer dr.mu.Unlock()

	// Initialize the rate limiter for the user if not already done
	if _, exists := dr.userLimits[userID]; !exists {
		dr.userLimits[userID] = NewDynamicSlidingWindow(10.0, time.Minute, 0.1) // 10 RPS, window of 1 minute, alpha = 0.1
	}

	// Evaluate if the request is allowed
	return dr.userLimits[userID].Allow()
}

func main() {
	limiter := NewDynamicRateLimiter(10.0, time.Minute, 0.1) // 10 RPS over a 1-minute window

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.FormValue("user")
		if userID == "" {
			http.Error(w, "Missing user parameter", http.StatusBadRequest)
			return
		}

		// Check if the request is allowed for this user
		if !limiter.Allow(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Successful request
		fmt.Fprintf(w, "Hello, %s!", userID)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}