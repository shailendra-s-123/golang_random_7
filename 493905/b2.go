package main  
import (  
    "fmt"
    "math/rand"
    "runtime"
    "sync"
    "time"
)

func simulateIOTask(delay time.Duration, start time.Time, wg *sync.WaitGroup, index int) {  
    ioWaitStart := time.Now()
    time.Sleep(delay) // Simulate I/O wait
    ioWaitDuration := time.Since(ioWaitStart)
    elapsed := time.Since(start)
    fmt.Printf("Goroutine %d finished in: %s (I/O wait: %s)\n", index, elapsed, ioWaitDuration)
    wg.Done()
}

func main() {  
    rand.Seed(time.Now().UnixNano())
    const numIOTasks = 10000 // Increased number of tasks for better analysis
    var numConcurrent = []int{1, 5, 10, 20, 50, 100, 200, 500}

    for _, concurrency := range numConcurrent {
        var wg sync.WaitGroup
        var activeGoroutines int
        var ioWaitTotal time.Duration

        start := time.Now()
        go func() {
            for {
                runtime.Gosched() // Yield control back to the scheduler
                activeGoroutines = runtime.NumGoroutine()
                fmt.Printf("Active Goroutines: %d\n", activeGoroutines)
                time.Sleep(10 * time.Millisecond)
            }
        }()

        for i := 0; i < numIOTasks; i++ {
            delay := time.Duration(rand.Intn(1000)) * time.Millisecond
            wg.Add(1)
            go simulateIOTask(delay, start, &wg, i)
        }
        wg.Wait()
        elapsed := time.Since(start)
        fmt.Printf("Experienced %d concurrent tasks: Total execution time = %s, Total I/O wait = %s\n", concurrency, elapsed, ioWaitTotal)
    }
}  