package main  
import (  
    "fmt"
    "net/http"
    "strconv"
    "sync"
    "time"
)

type userRateLimiter struct {
    limit       int
    burst       int
    userLimits  map[string]*userLimit
    userLimitsMu sync.Mutex
}

type userLimit struct {
    lastCheck time.Time
    requests  int
}

func newUserRateLimiter(limit, burst int) *userRateLimiter {
    return &userRateLimiter{
        limit:      limit,
        burst:      burst,
        userLimits: make(map[string]*userLimit),
    }
}

func (rl *userRateLimiter) allow(userID string) bool {
    rl.userLimitsMu.Lock()
    defer rl.userLimitsMu.Unlock()
    
    limit, ok := rl.userLimits[userID]
    if !ok {
        rl.userLimits[userID] = &userLimit{
            lastCheck: time.Now(),
            requests:  1,
        }
        return true // New user, allow immediately
    }
    
    now := time.Now()
    if now.Sub(limit.lastCheck) > time.Second {
        // If more than a second has passed, reset the counter
        limit.requests = 1
    } else {
        limit.requests++
    }
    
    limit.lastCheck = now
    allowed := limit.requests <= rl.limit+rl.burst
    
    return allowed
}
func main() {
    limiter := newUserRateLimiter(5, 1) // Limit of 5 requests per second with burst of 1

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        userID := r.FormValue("user")
        if userID == "" {
            http.Error(w, "Missing user parameter", http.StatusBadRequest)
            return
        }

        allowed := limiter.allow(userID)
        if !allowed {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        fmt.Fprintf(w, "Hello, %s!", userID)
    })
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
} 