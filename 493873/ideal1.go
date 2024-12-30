package main

import (
	"fmt"
	"sync"
	"time"
)

// Define the callback type
type DataCallback func(data []byte) ([]byte, error)

// Simulate an encoding function
func encodeData(data []byte) ([]byte, error) {
	// Simulate encoding delay
	time.Sleep(10 * time.Millisecond)
	return append([]byte("encoded:"), data...), nil
}

// Simulate a decoding function
func decodeData(data []byte) ([]byte, error) {
	// Simulate decoding delay
	time.Sleep(10 * time.Millisecond)
	return data[8:], nil
}

// Function to process data with callbacks for encoding/decoding
func processData(stream [][]byte, callback DataCallback, wg *sync.WaitGroup, resultChan chan<- []byte) {
	defer wg.Done()
	for _, data := range stream {
		transformedData, err := callback(data)
		if err != nil {
			fmt.Println("Error processing data:", err)
			continue
		}
		resultChan <- transformedData
	}
}

// Function to optimize using concurrency
func processStreamConcurrently(stream [][]byte, callback DataCallback) ([]byte, error) {
	var wg sync.WaitGroup
	resultChan := make(chan []byte, len(stream))
	
	// Split work into chunks and process concurrently
	chunkSize := len(stream) / 4 // Process in 4 concurrent workers
	for i := 0; i < 4; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == 3 {
			end = len(stream)
		}
		wg.Add(1)
		go processData(stream[start:end], callback, &wg, resultChan)
	}

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Collect results from the channel
	var result []byte
	for transformedData := range resultChan {
		result = append(result, transformedData...)
	}
	return result, nil
}

func main() {
	// Simulating a large data stream
	dataStream := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		dataStream[i] = []byte(fmt.Sprintf("data-%d", i))
	}

	// Measure start time
	startTime := time.Now()

	// Process data concurrently with encoding
	encodedData, err := processStreamConcurrently(dataStream, encodeData)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

	// Re-structure the data to pass as [][]byte for decoding
	// Each element of encodedData should be split into a single item slice for decoding
	encodedDataChunks := [][]byte{encodedData}

	// Process data concurrently with decoding
	decodedData, err := processStreamConcurrently(encodedDataChunks, decodeData)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		return
	}

	// Measure end time
	endTime := time.Now()
	fmt.Printf("Total processing time: %v\n", endTime.Sub(startTime))
	fmt.Printf("Decoded Data: %s\n", string(decodedData[:100])) // Display part of decoded data for verification
}