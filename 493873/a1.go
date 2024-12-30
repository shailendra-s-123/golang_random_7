package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

// Sample data structure
type Data struct {
	Value int32
}

// Callback functions
type Encoder func(d Data) []byte
type Decoder func(b []byte) (Data, error)

func main() {
	// Define simple encoding and decoding functions
	encoder := func(d Data) []byte {
		buf := bytes.NewBuffer(nil)
		if err := binary.Write(buf, binary.LittleEndian, d); err != nil {
			panic(err)
		}
		return buf.Bytes()
	}

	decoder := func(b []byte) (Data, error) {
		var d Data
		if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &d); err != nil {
			return Data{}, err
		}
		return d, nil
	}

	// Create a channel to handle processed data
	processedData := make(chan Data)

	// Worker function to process data using callbacks
	processData := func(w io.Writer, callback Encoder, decoder Decoder) {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, isPrefix, err := reader.ReadLine()
			if err != nil {
				log.Fatalf("Error reading from input: %v", err)
			}
			if isPrefix {
				log.Fatalf("Error reading from input: unexpected prefix")
			}

			// Simulate data processing
			data := Data{Value: int32(time.Now().UnixNano())}
			encodedData := callback(data)
			decodedData, err := decoder(encodedData)
			if err != nil {
				log.Fatalf("Error decoding data: %v", err)
			}

			// Write processed data to the channel
			processedData <- decodedData
		}
	}

	// Start processing in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		processData(os.Stdout, encoder, decoder)
	}()

	// Consume processed data
	for {
		data := <-processedData
		fmt.Printf("Processed data: %d\n", data.Value)
	}
}