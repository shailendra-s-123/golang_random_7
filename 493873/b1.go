package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// DataItem represents a single item in the data stream
type DataItem struct {
	Value int
}

// DataStream represents a stream of data
type DataStream chan DataItem

// EncoderCallback is a function type for encoding data
type EncoderCallback func(DataItem) DataItem

// DecoderCallback is a function type for decoding data
type DecoderCallback func(DataItem) DataItem

// NewDataStream creates a new DataStream with a given capacity
func NewDataStream(capacity int) DataStream {
	return make(DataStream, capacity)
}

// GenerateData generates random data items and sends them to the stream
func GenerateData(stream DataStream, wg *sync.WaitGroup, count int) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		item := DataItem{Value: rand.Int()}
		stream <- item
	}
	close(stream)
}

// ProcessData processes data items using callbacks for encoding and decoding
func ProcessData(stream DataStream, encoder EncoderCallback, decoder DecoderCallback, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range stream {
		encodedItem := encoder(item)
		decodedItem := decoder(encodedItem)
		// Do something with the decoded item
		_ = decodedItem
	}
}

func main() {
	const dataCount = 1000000
	const streamCapacity = 1000

	var wg sync.WaitGroup

	encodeCallback := func(item DataItem) DataItem {
		// Simple encoding example: add 1 to the value
		return DataItem{Value: item.Value + 1}
	}

	decodeCallback := func(item DataItem) DataItem {
		// Simple decoding example: subtract 1 from the value
		return DataItem{Value: item.Value - 1}
	}

	dataStream := NewDataStream(streamCapacity)

	wg.Add(2)
	go GenerateData(dataStream, &wg, dataCount)
	go ProcessData(dataStream, encodeCallback, decodeCallback, &wg)

	wg.Wait()

	fmt.Println("Data processing completed.")
}