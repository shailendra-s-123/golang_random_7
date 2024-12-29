package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

func benchmarkFmtSprintf(iterations int) string {
	var result string
	for i := 0; i < iterations; i++ {
		result = fmt.Sprintf("Number: %d, Name: %s, Value: %f", i, "Test", 3.14)
	}
	return result
}

func benchmarkStringBuilder(iterations int) string {
	var builder strings.Builder
	for i := 0; i < iterations; i++ {
		builder.WriteString(fmt.Sprintf("Number: %d, Name: %s, Value: %f", i, "Test", 3.14))
	}
	return builder.String()
}

func main() {
	// Set up CPU profiling
	f, err := os.Create("cpu_profile.out")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Start CPU profiling
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}

	// Benchmark fmt.Sprintf
	iterations := 1000000
	start := time.Now()
	benchmarkFmtSprintf(iterations)
	fmt.Println("fmt.Sprintf took:", time.Since(start))

	// Benchmark strings.Builder
	start = time.Now()
	benchmarkStringBuilder(iterations)
	fmt.Println("strings.Builder took:", time.Since(start))

	// Stop CPU profiling
	pprof.StopCPUProfile()

	// Memory profiling could also be added here, using pprof.WriteHeapProfile(f) for example.
}
