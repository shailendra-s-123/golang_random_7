package main

import (
	"fmt"
	"log"
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

const (
	defaultLimit      = 10         // Default max requests per user
	defaultResetTime  = 60 * time.Second // Default reset time
	increaseThreshold = 0.8         // If requests <= threshold * maxRequests, increase limit
	decreaseThreshold = 1.2         // If requests >= threshold * maxRequests, decrease limit
	limitIncreaseStep = 1         // How much to increase the limit
	limitDecreaseStep = -1        // How much to decrease the limit
)

var userLimits = map[string]*RateLimiter{}
var logger *log.Logger

func init() {
	// Initialize logger to a standard error
	logger = log.New(log.Writer(), "rate-limiter: ", log.LstdFlags|log.Lshortfile)
}

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

// Apply the rate limiting check and dynamically adjust the limit based on user behavior
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

	// User has exceeded the limit, decide whether to increase or decrease limit dynamically
	if float64(limiter.requests)/float64(limiter.maxRequests) < increaseThreshold {
		// User is under the threshold, increase the limit
		limiter.maxRequests += limitIncreaseStep
		logger.Printf("User %q's rate limit increased to %d", userID, limiter.maxRequests)
	} else if float64(limiter.requests)/float64(limiter.maxRequests) > decreaseThreshold {
		// User is over the threshold, decrease the limit
		limiter.maxRequests += limitDecreaseStep
		logger.Printf("User %q's rate limit decreased to %d", userID, limiter.maxRequests)
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

	// Get the rate limiter for the user
	limiter := getRateLimiter(userID)

	// Check if the user request is allowed
	if !limiter.isAllowed() {
		// User has exceeded the limit, respond with a clear message
		http.Error(w, "Rate limit exceeded. Please try again later or contact support for temporary rate increase.", http.StatusTooManyRequests)
		return
	}

	// If allowed, process the request
	w.WriteHeader(http.StatusOK)