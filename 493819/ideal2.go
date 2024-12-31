package main

import (
	"fmt"
	"log"
	"time"
)

type RateLimiter struct {
	UserID          string
	CurrentRate     int
	MaxRate         int
	Threshold       int
	ResetInterval   time.Duration
	BehaviorHistory []int // Stores the rate history for dynamic adjustments
	LastRequestTime time.Time
}

// Initialize a new rate limiter for a user
func NewRateLimiter(userID string, maxRate, threshold int, resetInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		UserID:        userID,
		CurrentRate:   0,
		MaxRate:       maxRate,
		Threshold:     threshold,
		ResetInterval: resetInterval,
		LastRequestTime: time.Now(),
	}
}

// Check if a user is within their rate limit
func (rl *RateLimiter) CheckLimit() bool {
	if rl.CurrentRate >= rl.MaxRate {
		log.Printf("User %s exceeded rate limit: %d/%d\n", rl.UserID, rl.CurrentRate, rl.MaxRate)
		return false
	}
	return true
}

// Simulate a request and adjust the rate accordingly
func (rl *RateLimiter) Request() string {
	// If enough time has passed to reset the rate limit, reset it
	if time.Since(rl.LastRequestTime) >= rl.ResetInterval {
		rl.CurrentRate = 0
		log.Printf("Rate limit reset for user %s.\n", rl.UserID)
	}

	// First check if the user is within the rate limit
	if !rl.CheckLimit() {
		// User exceeded rate limit, log error and stop processing further requests
		rl.LastRequestTime = time.Now() // Update the last request time to manage cooldown
		return fmt.Sprintf("Error: Rate limit exceeded. Please try again in %v.", rl.ResetInterval)
	}

	// Increment the current request count
	rl.CurrentRate++

	// Log the successful request
	log.Printf("User %s made a request: %d/%d\n", rl.UserID, rl.CurrentRate, rl.MaxRate)
	return "Request successful."
}

func main() {
	// Initialize rate limiter for a user with max rate 10 and a reset interval of 1 minute
	rl := NewRateLimiter("user123", 10, 5, time.Minute)

	// Simulate user behavior with a few requests
	for i := 0; i < 30; i++ { // Simulating 30 requests
		result := rl.Request()
		fmt.Println(result)

		// Stop further requests once the limit is exceeded
		if result == "Error: Rate limit exceeded. Please try again in 1m0s." {
			// Break the loop when the rate limit is exceeded
			break
		}

		time.Sleep(500 * time.Millisecond) // Simulate time delay between requests
	}
}