package main

import (
    "fmt"
    "net/http"
    "strconv"


    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

type UserRateLimiter struct {
    limiter *rate.Limiter
    burst   int
}

func (l *UserRateLimiter) Allow() bool {
    return l.limiter.Allow()
}

func NewUserRateLimiter(limit float64, burst int) *UserRateLimiter {
    return &UserRateLimiter{
        limiter: rate.NewLimiter(rate.Limit(limit), burst),
        burst:   burst,
    }
}

var rateLimits = make(map[string]*UserRateLimiter)

func main() {
    r := gin.Default()

    r.GET("/api/data", rateLimitMiddleware, func(c *gin.Context) {
        // Handle API requests here
        c.JSON(http.StatusOK, gin.H{"message": "Success"})
    })

    r.Run(":8080")
}

func rateLimitMiddleware(c *gin.Context) {
    // Extract user ID from query parameter
    userID := c.Query("user_id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        c.Abort()
        return
    }

    // Extract rate limit from query parameter (default to 10 requests per second)
    rateLimitStr := c.Query("rate_limit")
    rateLimit, err := strconv.ParseFloat(rateLimitStr, 64)
    if err != nil {
        rateLimit = 10.0
    }

    // Extract burst from query parameter (default to 10)
    burstStr := c.Query("burst")
    burst, err := strconv.Atoi(burstStr)
    if err != nil {
        burst = 10
    }

    // Get or create rate limiter for the user
    limiter, ok := rateLimits[userID]
    if !ok {
        limiter = NewUserRateLimiter(rateLimit, burst)
        rateLimits[userID] = limiter
    }

    if !limiter.Allow() {
        c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
        c.Abort()
        return
    }

    c.Next()
}