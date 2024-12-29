package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
)

// Define a Person struct
type Person struct {
	Name   string
	Age    int
	Height float64
}

func fmtSprintf(person Person, cities []string, visits map[string]int) string {
	return fmt.Sprintf("Hello, my name is %s. I am %d years old and %.2f feet tall. I have visited %d cities: %v. My most visited city is %s with %d visits.",
		person.Name, person.Age, person.Height, len(cities), cities,
		visits[fmt.Sprintf("%s", cities[0])], visits[fmt.Sprintf("%s", cities[0])])
}

func stringsBuilder(person Person, cities []string, visits map[string]int) string {
	buf := &bytes.Buffer{}
	buf.WriteString("Hello, my name is ")
	buf.WriteString(person.Name)
	buf.WriteString(". I am ")
	buf.WriteString(strconv.Itoa(person.Age))
	buf.WriteString(" years old and ")
	buf.WriteString(fmt.Sprintf("%.2f", person.Height))
	buf.WriteString(" feet tall. I have visited ")
	buf.WriteString(strconv.Itoa(len(cities)))
	buf.WriteString(" cities: ")
	buf.WriteString(fmt.Sprint(cities...))
	buf.WriteString(". My most visited city is ")
	buf.WriteString(cities[0])
	buf.WriteString(" with ")
	buf.WriteString(strconv.Itoa(visits[cities[0]]))
	buf.WriteString(" visits.")
	return buf.String()
}

func main() {
	// Create a Person instance
	person := Person{Name: "Alice", Age: 30, Height: 5.9}

	// Create a slice of cities
	cities := []string{"New York", "Los Angeles", "Chicago"}
	// Create a map with city names and visit counts
	visits := map[string]int{"New York": 5, "Los Angeles": 3, "Chicago": 2}

	// Start CPU profiling to measure computational overhead
	f, err := os.Create("profile_fmt.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Measure the time and memory usage of fmt.Sprintf with complex data types
	numRuns := 1000000
	startTime := time.Now()
	for i := 0; i < numRuns; i++ {
		fmt.Sprintf("%v", fmtSprintf(person, cities, visits))
	}
	endTime := time.Now()
	durationFmt := endTime.Sub(startTime)
	fmt.Printf("Time taken for %d runs using fmt.Sprintf: %s\n", numRuns, durationFmt)

	// Measure memory allocation
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated memory using fmt.Sprintf: %d MB\n", m.Alloc/1024/1024)
	fmt.Printf("Total allocated memory using fmt.Sprintf: %d MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("System memory using fmt.Sprintf: %d MB\n", m.Sys/1024/1024)
	fmt.Printf("Number of garbage collections using fmt.Sprintf: %d\n", m.NumGC)

	// Start CPU profiling to measure computational overhead
	f, err = os.Create("profile_builder.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Measure the time and memory usage of strings.Builder with complex data types
	startTime = time.Now()
	for i := 0; i < numRuns; i++ {
		fmt.Sprintf("%v", stringsBuilder(person, cities, visits))
	}
	endTime = time.Now()
	durationBuilder := endTime.Sub(startTime)
	fmt.Printf("Time taken for %d runs using strings.Builder: %s\n", numRuns, durationBuilder)

	// Measure memory allocation
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated memory using strings.Builder: %d MB\n", m.Alloc/1024/1024)
	fmt.Printf("Total allocated memory using strings.Builder: %d MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("System memory using strings.Builder: %d MB\n", m.Sys/1024/1024)
	fmt.Printf("Number of garbage collections using strings.Builder: %d\n", m.NumGC)
}