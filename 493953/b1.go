package main  
import (  
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/justinas/alice"
    "github.com/xslavik/go-httpbench/runner"
)

// ServerStats holds performance metrics for the server
type ServerStats struct {
    TotalRequests int
    LatencySum    time.Duration
    LatencyMin    time.Duration
    LatencyMax    time.Duration
    ErrorCount   int
}

// benchHandler is a sample HTTP handler that sleeps for a random duration
func benchHandler(w http.ResponseWriter, r *http.Request) {
    // Simulate some work by sleeping for a random duration between 0 and 10ms
    time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

    // Write a simple response
    fmt.Fprintf(w, "OK")
}

// logMiddleware is an HTTP middleware to log requests and responses
func logMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        h.ServeHTTP(w, r)
    })
}

// measureMiddleware is an HTTP middleware to measure latency and count errors
func measureMiddleware(stats *ServerStats) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        defer func() {
            duration := time.Since(start)
            stats.TotalRequests++
            stats.LatencySum += duration
            if stats.LatencyMin == 0 || duration < stats.LatencyMin {
                stats.LatencyMin = duration
            }
            if duration > stats.LatencyMax {
                stats.LatencyMax = duration
            }
            if err := recover(); err != nil {
                stats.ErrorCount++
            }
        }()
        h.ServeHTTP(w, r)
    })
}

func main() {
    // Initialize server stats
    var stats ServerStats

    // Create an HTTP chain with middlewares
    chain := alice.New(logMiddleware, measureMiddleware(&stats)).ThenFunc(benchHandler)

    // Create an HTTP server
    server := &http.Server{
        Addr:    ":8080",
        Handler: chain,
    }

    // Start the server in a goroutine
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Define the workload parameters
    workloads := []runner.Workload{
        {
            Name:          "Constant Load",
            Concurrency:    100,
            Rate:           100,
            Duration:       time.Minute,
        },
        {
            Name:          "Burst Load",
            Concurrency:    100,
            Rate:           1000,
            Duration:       time.Second * 10,
            BurstSize:      200,
            BurstFrequency: 5,
        },
        {
            Name:          "Spike Load",
            Concurrency:    100,
            Rate:           100,
            Duration:       time.Minute,
            Spikes: []runner.Spike{
                {
                    Duration:  time.Second * 10,
                    Rate:      1000,
                    Delay:     time.Second * 30,
                },
            },
        },
    }

    // Start the benchmark
    var wg sync.WaitGroup
    for _, workload := range workloads {
        wg.Add(1)
        go func(workload runner.Workload) {
            defer wg.Done()
            log.Printf("Starting workload: %s", workload.Name)
            runner.Run(&workload, server.URL)
            log.Printf("Workload %s completed", workload.Name)