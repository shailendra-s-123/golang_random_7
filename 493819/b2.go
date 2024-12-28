package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type slidingWindow struct {
	limit     int
	window    time.Duration
	requests  []time.Time
	requestsMu sync.Mutex
}

func newSlidingWindow(limit int, window time.Duration) *slidingWindow {
	return &slidingWindow{
		limit:  limit,
		window: window,
	}
}

func (w *slidingWindow) allow() bool {
	w.requestsMu.Lock()
	defer w.requestsMu.Unlock()

	// Remove requests older than the window duration
	now := time.Now()
	for i, reqTime := range w.requests {
		if now.Sub(reqTime) > w.window {
			w.requests = w.requests[i+1:]
			break
		}
	}

	// Add the current request
	w.requests = append(w.requests, now)

	// Check if the limit is exceeded
	return len(w.requests) <= w.limit
}

type userRateLimiter struct {
	userLimits  map[string]*slidingWindow
	userLimitsMu sync.Mutex
}

func newUserRateLimiter() *userRateLimiter {
	return &userRateLimiter{
		userLimits: make(map[string]*slidingWindow),
	}
}

func (rl *userRateLimiter) allow(userID string, limit int, window time.Duration) bool {
	rl.userLimitsMu.Lock()
	defer rl.userLimitsMu.Unlock()

	// Retrieve or initialize the user's rate limit window
	window, ok := rl.userLimits[userID]
	if !ok {
		rl.userLimits[userID] = newSlidingWindow(limit, window)
	}

	return window.allow()
}

func main() {
	limiter := newUserRateLimiter()

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

		allowed := limiter.allow(userID, limit, window)
		if !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		fmt.Fprintf(w, "Hello, %s!", userID)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}  