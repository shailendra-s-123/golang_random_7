package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Initialize logging
var log = logrus.New()
log.SetLevel(logrus.DebugLevel)

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
var adjustmentFactor = 1.2 // Increase factor for good behavior
var decreaseFactor = 0.8  // Decrease factor for hitting limit

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

// Adjust rate limits dynamically
func (limiter *RateLimiter) adjustLimit(allowed bool) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	if allowed {
		limiter.maxRequests = int(float64(limiter.maxRequests) * adjustmentFactor)
	} else {
		limiter.maxRequests = int(float64(limiter.maxRequests) * decreaseFactor)
	}

	log.WithFields(logrus.Fields{
		"userID":     limiter.maxRequests,
		"newLimit":   limiter.maxRequests,
		"allowed":    allowed,
		"resetTime":  limiter.resetTime,
		"requests":   limiter.requests,
		"lastReset":  limiter.lastReset,
	}).Debug("Adjusted rate limit")
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
			log.WithFields(logrus.Fields{
				"userID":     userID,
				"newLimit":   limiter.maxRequests,
				"adjusted":   true,
			}).Debug("User-defined limit set")
		}
	}

	// Get the rate limiter for the user
	limiter := getRateLimiter(userID)

	// Check if the user request is allowed
	if !limiter.isAllowed() {
		limiter.adjustLimit(false)
		http.Error(w, "Rate limit exceeded. Please try again later or contact support for assistance.", http.StatusTooManyRequests)
		return
	}

	// If allowed, process the request (for demonstration, just return success)
	limiter.adjustLimit(true)
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