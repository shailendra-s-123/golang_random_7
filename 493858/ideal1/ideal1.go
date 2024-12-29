package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// Define a Person struct
type Person struct {
	Name   string
	Age    int
	Height float64
}

func main() {
	// Create a Person instance
	person := Person{Name: "Alice", Age: 30, Height: 5.9}

	// Create a slice of cities
	cities := []string{"New York", "Los Angeles", "Chicago"}
	// Create a map with city names and visit counts
	visits := map[string]int{"New York": 5, "Los Angeles": 3, "Chicago": 2}

	// Start CPU profiling to measure computational overhead
	f, err := os.Create("profile.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Measure the time and memory usage of fmt.Sprintf with complex data types
	numRuns := 10000
	for i := 0; i < numRuns; i++ {
		// Format a string using fmt.Sprintf and multiple data types
		str := fmt.Sprintf("Hello, my name is %s. I am %d years old and %.2f feet tall. I have visited %d cities: %v. My most visited city is %s with %d visits.",
			person.Name, person.Age, person.Height, len(cities), cities,
			visits[fmt.Sprintf("%s", cities[0])], visits[fmt.Sprintf("%s", cities[0])])

		// Discard the string to avoid allocation in benchmark
		_ = str
	}

	// Measure the time taken for the tests
	startTime := time.Now()
	for i := 0; i < numRuns; i++ {
		// Format a string using fmt.Sprintf and multiple data types
		str := fmt.Sprintf("Hello, my name is %s. I am %d years old and %.2f feet tall. I have visited %d cities: %v. My most visited city is %s with %d visits.",
			person.Name, person.Age, person.Height, len(cities), cities,
			visits[fmt.Sprintf("%s", cities[0])], visits[fmt.Sprintf("%s", cities[0])])
		_ = str
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Time taken for %d runs: %s\n", numRuns, duration)

	// Measure memory allocation
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated memory: %d MB\n", m.Alloc/1024/1024)
	fmt.Printf("Total allocated memory: %d MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("System memory: %d MB\n", m.Sys/1024/1024)
	fmt.Printf("Number of garbage collections: %d\n", m.NumGC)
}
