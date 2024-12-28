package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// slidingWindow struct to manage rate limiting with a sliding window
type slidingWindow struct {
	limit    int        // Maximum requests allowed
	window   time.Duration // Time window in which requests are counted
	buckets  map[time.Time][]string
	bucketsMu sync.Mutex
}

// NewSlidingWindow creates a new sliding window rate limiter for users with specified limit and window size
func newSlidingWindow(limit int, window time.Duration) *slidingWindow {
	return &slidingWindow{
		limit:    limit,
		window:   window,
		buckets:  make(map[time.Time][]string),
	}
}

// cleanOldBuckets removes buckets that are older than the sliding window
func (sw *slidingWindow) cleanOldBuckets() {
	now := time.Now()
	for key := range sw.buckets {
		if key.Add(sw.window).Before(now) {
			delete(sw.buckets, key)
		}
	}
}

// allow checks if a user's request is within the rate limit using a sliding window
func (sw *slidingWindow) allow(userID string) bool {
	sw.bucketsMu.Lock()
	defer sw.bucketsMu.Unlock()

	now := time.Now()
	sw.cleanOldBuckets()

	// Get bucket for current second
	for key, _ := range sw.buckets {
		if key.Year() == now.Year() && key.Month() == now.Month() && key.Day() == now.Day() && key.Hour() == now.Hour() && key.Minute() == now.Minute() && key.Second() == now.Second() {
			bucket, ok := sw.buckets[key]
			if ok {
				// Check if the user is already in the bucket
				if contains(bucket, userID) {
					return false
				}
				// Add user to the bucket
				sw.buckets[key] = append(bucket, userID)
				return true
			}
			// Create new bucket if it doesn't exist
			sw.buckets[key] = []string{userID}
			return true
		}
	}

	// Create new bucket for current second if none exists
	sw.buckets[now] = []string{userID}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	// Initialize rate limiter with a limit of 5 requests per 10 seconds
	limiter := newSlidingWindow(5, time.Second*10)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.FormValue("user") // Extract userID from the query parameter
		limitStr := r.FormValue("limit")
		windowStr := r.FormValue("window")

		if userID == "" {
			http.Error(w, "Missing user parameter", http.StatusBadRequest)
			return
		}

		if limitStr != "" {
			limit, err := strconv.Atoi(limitStr)
			if err != nil {
				http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
				return
			}
			limiter.limit = limit
		}

		if windowStr != "" {
			duration, err := strconv.ParseInt(windowStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid window parameter", http.StatusBadRequest)
				return
			}
			limiter.window = time.Duration(duration) * time.Second
		}

		// Check if the user is allowed to make a request
		allowed := limiter.allow(userID)
		if !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		fmt.Fprintf(w, "Hello, %s! Rate Limit: %d/%ds\n", userID, limiter.limit, limiter.window.Seconds()) // Successful request
	})

	// Start the HTTP server
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}