package main

import (
	_"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// ... (existing code)

// Define recoverable errors for which retry should be attempted
var recoverableErrors = map[error]struct{}{
	&TimeoutError{}: {},
}

// Function that applies a callback to each item, retrying on transient errors
func ProcessItemsWithRetry(items []string, callback func(string) error, retries int) error {
	log.Info("Starting processing of items...")
	defer log.Info("Processing completed.")

	for _, item := range items {
		err := processItemWithRetry(item, callback, retries)
		if err != nil {
			return fmt.Errorf("error processing item %s: %w", item, err)
		}
	}
	return nil
}

func processItemWithRetry(item string, callback func(string) error, retries int) error {
	log.WithFields(log.Fields{"item": item}).Info("Processing item...")

	for attempt := 1; attempt <= retries+1; attempt++ {
		start := time.Now()
		err := callback(item)
		if err == nil {
			log.WithFields(log.Fields{
				"item":     item,
				"duration": time.Since(start),
				"attempts": attempt,
			}).Info("Item processed successfully.")
			return nil
		}

		// Check if the error is recoverable
		if _, ok := recoverableErrors[err]; ok {
			log.WithFields(log.Fields{
				"item":     item,
				"duration": time.Since(start),
				"attempts": attempt,
				"error":    err,
			}).Error("Error processing item, retrying...")

			// Introduce a delay before retrying
			backoff := time.Duration(rand.Intn(1000)) * time.Millisecond
			time.Sleep(backoff)
		} else {
			// Non-recoverable error, return immediately
			return err
		}
	}

	// All retries failed
	return fmt.Errorf("all retries failed for item %s", item)
}
 
// ... (existing code)
func main() {
    // ... (existing code)

    // Set up the logger
    setupLogging()
    
    // List of items to process
    items := []string{"item1", "timeout", "bad-format", "timeout", "unexpected", "timeout", "item2"}
    retries := 3  // Number of retry attempts for transient errors
    
    // Handle panic recovery
    defer func() {
        if r := recover(); r != nil {
            log.WithField("panic", r).Panic("Recovered from panic")
        }
    }()

    // Process items with retry
    err := ProcessItemsWithRetry(items, ExampleCallback, retries)
    if err != nil {
        log.WithError(err).Error("Processing items failed")
        os.Exit(1)
    }
    
    log.Info("All items processed successfully.")
}