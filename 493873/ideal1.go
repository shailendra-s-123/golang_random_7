package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"sync"
	//	"time"
)

// Data represents a simple data structure to be encoded and decoded
type Data struct {
	Value int32
}

// Encoder is a callback function type for encoding data
type Encoder func(d Data) []byte

// Decoder is a callback function type for decoding data
type Decoder func(b []byte) (Data, error)

// Function to encode Data to a byte slice
func encode(d Data) []byte {
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.LittleEndian, d); err != nil {
		log.Fatalf("Error encoding data: %v", err)
	}
	return buf.Bytes()
}

// Function to decode a byte slice back to Data
func decode(b []byte) (Data, error) {
	var d Data
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &d); err != nil {
		return Data{}, err
	}
	return d, nil
}

// Function to simulate real-time data generation (data stream)
func generateData(wg *sync.WaitGroup, ch chan Data) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		ch <- Data{Value: rand.Int31()}
	}
	close(ch)
}

// Function to process data using callbacks for encoding and decoding
func processData(wg *sync.WaitGroup, ch chan Data, encoder Encoder, decoder Decoder) {
	defer wg.Done()
	for data := range ch {
		encodedData := encoder(data)   // Encode the data
		_, err := decoder(encodedData) // Decode the data
		if err != nil {
			log.Printf("Error decoding data: %v", err)
		}
	}
}

func main() {
	// Channel for simulating the data stream
	dataChannel := make(chan Data, 1000)
	var wg sync.WaitGroup

	// Start data generation in a goroutine
	wg.Add(1)
	go generateData(&wg, dataChannel)

	// Start data processing in a goroutine with callbacks
	wg.Add(1)
	go processData(&wg, dataChannel, encode, decode)

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("Data processing completed.")
}

