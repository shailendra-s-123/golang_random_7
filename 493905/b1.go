package main  
import (  
    "fmt"
    "math/rand"
    "sync"
    "time"
)

const (  
    numTasks      = 10000
    numGoroutines = 10
    maxDelay      = 10 * time.Millisecond
)

func ioBoundTask(wg *sync.WaitGroup, id int, startTime time.Time) {  
    defer wg.Done()
    delay := time.Duration(rand.Intn(int(maxDelay)))
    time.Sleep(delay)
    endTime := time.Now()
    fmt.Printf("Goroutine %d: Delay: %v, Execution time: %v\n", id, delay, endTime.Sub(startTime))
}

func main() {  
    rand.Seed(time.Now().UnixNano())
    var wg sync.WaitGroup
    startTime := time.Now()
    wg.Add(numTasks)
    for i := 0; i < numTasks; i++ {
        go ioBoundTask(&wg, i, startTime)
    }
    wg.Wait()
    endTime := time.Now()
    totalExecutionTime := endTime.Sub(startTime)
    fmt.Printf("Total execution time: %v\n", totalExecutionTime)
}