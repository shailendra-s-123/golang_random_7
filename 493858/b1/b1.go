package main
import (
        "fmt"
        "log"
        "runtime/pprof"
        "time"
)
// define a Person struct
type Person struct {
        Name string
        Age int
}
func main() {
        // Create a Person instance
        person := Person{"Alice", 25}
        // Create a slice of Person instances
        people := make([]Person, 10000)
        for i := 0; i < len(people); i++ {
                people[i].Name = fmt.Sprintf("Person %d", i+1)
                people[i].Age = i + 1
        }
        // Profile the computation overhead and memory usage of fmt.Sprintf
        f, err := os.Create("profile.out")
        if err != nil {
                log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
        // Run the performance test multiple times
        numRuns := 10000
        for i := 0; i < numRuns; i++ {
                // Format a string using fmt.Sprintf and multiple data types
                str := fmt.Sprintf("Hello, my name is %s, I am %d years old, and my friends are: %v", person.Name, person.Age, people)
                // Discard the string to avoid allocation in benchmark
                _ = str
        }
        // Measure the time taken for the tests
        startTime := time.Now()
        for i := 0; i < numRuns; i++ {
                // Format a string using fmt.Sprintf and multiple data types
                str := fmt.Sprintf("Hello, my name is %s, I am %d years old, and my friends are: %v", person.Name, person.Age, people)
                // Discard the string to avoid allocation in benchmark
                _ = str
        }
        endTime := time.Now()
        duration := endTime.Sub(startTime)
        fmt.Printf("Time taken for %d runs: %s\n", numRuns, duration)
        // Measure the memory allocation using runtime.MemStats
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        fmt.Printf("Allocated memory: %d MB\n", m.Alloc/1024/1024)
        fmt.Printf("Total allocated memory: %d MB\n", m.TotalAlloc/1024/1024)
        fmt.Printf("System memory: %d MB\n", m.Sys/1024/1024)
        fmt.Printf("Number of garbage collections: %d\n", m.NumGC)
} 