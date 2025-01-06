package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// Define a custom error type for clarity and traceability
type ProcessError struct {
	message string
	cause   error // to support chaining of errors
}

func (e *ProcessError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("ProcessError: %s (cause: %v)", e.message, e.cause)
	}
	return fmt.Sprintf("ProcessError: %s", e.message)
}

// Function that applies a callback to each item in a list with retry mechanism
func ProcessItemsWithRetry(items []string, callback func(context.Context, string) error, retries int) error {
	for _, item := range items {
		for retryCount := 0; retryCount <= retries; retryCount++ {
			ctx := context.WithValue(context.Background(), "retryCount", retryCount)
			err := callback(ctx, item)
			if err == nil {
				log.Printf("Processing item %s succeeded after %d retries.", item, retryCount)
				break
			}
			log.Printf("Processing item %s failed (retry %d): %v. Retrying...", item, retryCount, err)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		if err != nil {
			return fmt.Errorf("error processing item %s after all retries: %w", item, err)
		}
	}
	return nil
}

// Example callback function with error handling
func ExampleCallback(ctx context.Context, item string) error {
	retryCount := ctx.Value("retryCount").(int)

	// Simulate a temporary error that can be retried
	if item == "temporary-error" && retryCount < 2 {
		return &ProcessError{message: "Temporary error occurred", cause: fmt.Errorf("internal processing error")}
	}

	// Simulate a permanent error
	if item == "permanent-error" {
		return fmt.Errorf("item has a permanent error")
	}

	log.Printf("Processing item: %s, retry count: %d", item, retryCount)
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	items := []string{"good", "temporary-error", "permanent-error", "good"}

	err := ProcessItemsWithRetry(items, ExampleCallback, 2)
	if err != nil {
		log.Printf("Processing failed: %v\n", err)
	} else {
		log.Println("Processing completed successfully.")
	}
}