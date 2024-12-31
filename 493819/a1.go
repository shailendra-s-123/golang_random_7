package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var requestLimits map[string]int // Map to store request limits for each user
var requestCount map[string]int  // Map to store number of requests per user
var windowDuration = time.Minute  // Time window for counting requests

// Initialize the request limit and request count maps
func init() {
	requestLimits = make(map[string]int)
	requestCount = make(map[string]int)
}

// GetUserID retrieves the user ID from the request URL query parameters
func getUserID(r *http.Request) string {
	return r.URL.Query().Get("user_id")
}

// IsRateLimited checks if the user has exceeded the rate limit
func isRateLimited(userID string, limit int) bool {
	// Check if the limit is set
	if limit <= 0 {
		return false // No limit means no restriction
	}

	// Get the current time and reset the request count map every windowDuration
	currentTime := time.Now()
	keysToRemove := make([]string, 0)
	for key, count := range requestCount {
		if currentTime.Sub(requestCount[key].time) > windowDuration {
			keysToRemove = append(keysToRemove, key)
		}
	}
	for _, key := range keysToRemove {
		delete(requestCount, key)
	}

	// Check if the user has already reached the limit in the current window
	if count, ok := requestCount[userID]; ok && count >= limit {
		return true
	}

	// Update request count
	requestCount[userID] = requestCount[userID]+1 if ok else 1
	requestCount[userID].time = currentTime
	return false
}

// HandleOverload manages the response when a user exceeds the rate limit
func handleOverload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, `{"message": "Rate limit exceeded for user_id: %s", "limit": %d}`, getUserID(r), requestLimits[getUserID(r)])
}

// ExampleHandler is a simple handler function to demonstrate rate limiting
func exampleHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	limit, ok := requestLimits[userID]
	if !ok {
		// Set a default limit if user_id not found
		limit = 10
		requestLimits[userID] = limit
	}

	if isRateLimited(userID, limit) {
		handleOverload(w, r)
		return
	}

	fmt.Fprintf(w, "Hello, user_id: %s! Your limit is %d requests per %v.\n", userID, limit, windowDuration)
}

func main() {
	// Define some example routes
	http.HandleFunc("/example", exampleHandler)

	// Start the server
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}