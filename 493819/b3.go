package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// AdaptiveSlidingWindow struct to manage adaptive rate limiting with a sliding window
type AdaptiveSlidingWindow struct {
	baseLimit    int               // Base rate limit before adaptation
	window       time.Duration     // Window size
	requests     []time.Time
	mu           sync.Mutex
	adaptFactor  float64           // Factor for adapting rate limit dynamically
	maxAdaptRate float64           // Maximum rate of adaptation (10% per minute by default)
	lastAdapted   time.Time         // Time of the last adaptation
}

// NewAdaptiveSlidingWindow creates a new adaptive sliding window rate limiter for users with specified base limit and window size
func NewAdaptiveSlidingWindow(baseLimit int, window time.Duration) *AdaptiveSlidingWindow {
	return &AdaptiveSlidingWindow{
		baseLimit:    baseLimit,
		window:       window,
		adaptFactor:  1.0,
		maxAdaptRate: 0.1, // 10% per minute
		lastAdapted:  time.Now(),
	}
}

// calculateRequestsPerMinute calculates the requests per minute based on the requests in the current window.
func (sw *AdaptiveSlidingWindow) calculateRequestsPerMinute() float64 {
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
	sw.requests = validRequests

	// Calculate requests per minute
	return float64(len(validRequests)) / time.Minute.Seconds()
}

// Allow checks if a user's request is within the adaptive rate limit using a sliding window
func (sw *AdaptiveSlidingWindow) Allow() bool {
	rps := sw.calculateRequestsPerMinute()
	currentLimit := int(float64(sw.baseLimit) * sw.adaptFactor)

	// If the request per minute exceeds the limit, deny the request
	if rps >= float64(currentLimit) {
		return false
	}

	// Allow the request
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.requests = append(sw.requests, time.Now())
	return true
}

// Adapt dynamically adjusts the rate limit based on the user's activity over time.
func (sw *AdaptiveSlidingWindow) Adapt() {
	rps := sw.calculateRequestsPerMinute()
	now := time.Now()

	// Prevent too frequent adaptations
	if now.Sub(sw.lastAdapted) < time.Minute {
		return
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.lastAdapted = now

	if rps >= float64(sw.baseLimit) {
		// Reduce the rate limit if it is exceeded
		targetAdaptRate := math.Min(sw.maxAdaptRate, (float64(sw.baseLimit)-rps)/float64(sw.baseLimit))
		sw.adaptFactor = math.Max(1.0, sw.adaptFactor*(1-targetAdaptRate))
	} else {
		// Increase the rate limit if traffic is low
		targetAdaptRate := math.Min(sw.maxAdaptRate, (rps/float64(sw.baseLimit) - 1))
		sw.adaptFactor = math.Min(float64(sw.baseLimit), sw.adaptFactor*(1+targetAdaptRate))
	}

	fmt.Printf("Adaptive Factor: %.2f\n", sw.adaptFactor)
}

func main() {
	limiter := NewRateLimiter()
	baseLimit := 10   // Base rate limit (e.g., 10 requests per minute)
	adaptionInterval := 1 * time.Minute // Adapt rate limit every minute

	go func() {
		for range time.Tick(adaptionInterval) {
			limiter.AdaptAll()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ... (User authentication and rate limit logic remains the same)

		// Before allowing the request, adapt the rate limit for the user
		limiter.Adapt(userID)

		if !limiter.Allow(userID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
	})