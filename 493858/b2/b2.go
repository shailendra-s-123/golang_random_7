package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
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

	numRuns := 100000

	// Compare fmt.Sprintf with strings.Builder for string formatting
	fmt.Println("\nComparing fmt.Sprintf and strings.Builder:")

	// **fmt.Sprintf Approach**
	// Start CPU profiling
	f, err := os.Create("profile_fmt.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	startTimeFmt := time.Now()
	for i := 0; i < numRuns; i++ {
		_ = fmt.Sprintf("Hello, my name is %s. I am %d years old and %.2f feet tall. I have visited %d cities: %v. My most visited city is %s with %d visits.",
			person.Name, person.Age, person.Height, len(cities), cities,
			visits[fmt.Sprintf("%s", cities[0])], visits[fmt.Sprintf("%s", cities[0])])
	}
	endTimeFmt := time.Now()
	durationFmt := endTimeFmt.Sub(startTimeFmt)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("fmt.Sprintf: Time taken for %d runs: %s\n", numRuns, durationFmt)
	fmt.Printf("fmt.Sprintf: Allocated memory: %d MB\n", m.Alloc/1024/1024)
	fmt.Printf("fmt.Sprintf: Total allocated memory: %d MB\n", m.TotalAlloc/1024/1024)

	// **strings.Builder Approach**
	// Start CPU profiling
	f, err = os.Create("profile_builder.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	startTimeBuilder := time.Now()
	var b strings.Builder
	for i := 0; i < numRuns; i++ {
		b.Reset()
		fmt.Fprintf(&b, "Hello, my name is %s. I am %d years old and %.2f feet tall. I have visited %d cities: %v. My most visited city is %s with %d visits.",
			person.Name, person.Age, person.Height, len(cities), cities,
			visits[fmt.Sprintf("%s", cities[0])], visits[fmt.Sprintf("%s", cities[0])])
		_ = b.String()
	}
	endTimeBuilder := time.Now()
	durationBuilder := endTimeBuilder.Sub(startTimeBuilder)

	runtime.ReadMemStats(&m)
	fmt.Printf("strings.Builder: Time taken for %d runs: %s\n", numRuns, durationBuilder)
	fmt.Printf("strings.Builder: Allocated memory: %d MB\n", m.Alloc/1024/1024)
	fmt.Printf("strings.Builder: Total allocated memory: %d MB\n", m.TotalAlloc/1024/1024) 