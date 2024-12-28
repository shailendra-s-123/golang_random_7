package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// dynamicSlidingWindow manages a dynamically adjusting rate limiting window for users
type dynamicSlidingWindow struct {
	targetLimit   float64 // Maximum allowed requests per second
	windowSize    time.Duration
	requests      []time.Time
	mu            sync.Mutex
	initialized    bool
	rateEstimate   float64
	lastRequest    time.Time
	alpha         float64 // Exponentially weighted moving average smoothing parameter
}

// newDynamicSlidingWindow creates a dynamic rate limiter
func newDynamicSlidingWindow(targetLimit float64, windowSize time.Duration, alpha float64) *dynamicSlidingWindow {
	return &dynamicSlidingWindow{
		targetLimit:   targetLimit,
		windowSize:    windowSize,
		mu:            sync.Mutex{},
		alpha:         alpha,
	}
}

// allow evaluates whether a request is within the adjusted rate limit
func (dsw *dynamicSlidingWindow) allow() bool {
	dsw.mu.Lock()
	defer dsw.mu.Unlock()

	// Calculate the current request rate based on the time window
	now := time.Now()

	if !dsw.initialized {
		dsw.lastRequest = now
		dsw.initialized = true
		return true
	}

	// Clean out expired requests
	var validRequests []time.Time
	for _, reqTime := range dsw.requests {
		if now.Sub(reqTime) <= dsw.windowSize {
			validRequests = append(validRequests, reqTime)
		}
	}
	dsw.requests = validRequests

	// Calculate the current request rate
	elapsed := now.Sub(dsw.lastRequest)
	currentRate := float64(len(dsw.requests)) / elapsed.Seconds()

	// Update estimate using exponentially weighted moving average (EWMA)
	dsw.rateEstimate = dsw.alpha*currentRate + (1-dsw.alpha)*dsw.rateEstimate

	// Allow if the rate estimate is within or below the target limit
	allowed := dsw.rateEstimate <= dsw.targetLimit
	dsw.lastRequest = now
	dsw.requests = append(dsw.requests, now)
	return allowed
}

type dynamicRateLimiter struct {
	userLimits   map[string]*dynamicSlidingWindow
	mu           sync.Mutex
	targetLimit  float64 // Required: Average requests per second
	windowSize    time.Duration
	alpha        float64 // Smoothing parameter for EWMA
	initialRPS   float64 // Alternative way to specify targeted activity
}

// NewDynamicRateLimiter creates a dynamic rate limiter based on average rate and window
func NewDynamicRateLimiter(targetLimit float64, windowSize time.Duration, alpha float64) *dynamicRateLimiter {
	return &dynamicRateLimiter{
		userLimits:   make(map[string]*dynamicSlidingWindow),
		mu:           sync.Mutex{},
		targetLimit:  targetLimit,
		windowSize:    windowSize,
		alpha:        alpha,
	}
}

// Allow determines if a request should be processed based on each user's limits
func (drl *dynamicRateLimiter) Allow(userID string) bool {
	drl.mu.Lock()
	defer drl.mu.Unlock()

	// Ensure user has a limiter
	if _, exists := drl.userLimits[userID]; !exists {
		drl.userLimits[userID] = newDynamicSlidingWindow(drl.targetLimit, drl.windowSize, drl.alpha)
	}

	// Evaluate if the request is allowed
	return drl.userLimits[userID].allow()
}

func main() {
	limiter := NewDynamicRateLimiter(10.0, time.Second*60, 0.02) // Target average: 10 RPS, over a 1-minute window

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