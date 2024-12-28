package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus" // Correct import for logrus
)

// Custom error types for better traceability
type TimeoutError struct {
	msg string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("Timeout error: %s", e.msg)
}

type ProcessingError struct {
	message string
	cause   error
}

func (e *ProcessingError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("ProcessingError: %s (cause: %v)", e.message, e.cause)
	}
	return fmt.Sprintf("ProcessingError: %s", e.message)
}

// Function to apply a callback with retry logic for transient errors
func ProcessItemWithRetry(item string, callback func(string) error, maxRetries int) error {
	ctx := context.Background()
	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		start := time.Now()
		log.WithFields(log.Fields{
			"item":     item,
			"attempt":  attempt,
			"maxAttempts": maxRetries,
		}).Info("Attempting to process item...")

		err = callback(item)
		if err == nil {
			log.WithFields(log.Fields{
				"item":       item,
				"duration":   time.Since(start),
				"attempts":   attempt + 1,
				"status":     "success",
			}).Info("Item processed successfully.")
			return nil
		}

		// Check if it's a transient error and retry
		if isTransientError(err) {
			log.WithFields(log.Fields{
				"item":       item,
				"duration":   time.Since(start),
				"attempts":   attempt + 1,
				"status":     "error",
				"error":      err,
			}).Warning("Transient error occurred. Retrying...")

			// Introduce a backoff between retry attempts
			time.Sleep(time.Second * time.Duration(attempt)*2) // Exponential backoff
			continue
		}

		// Log permanent error
		log.WithFields(log.Fields{
			"item":       item,
			"duration":   time.Since(start),
			"attempts":   attempt + 1,
			"status":     "error",
			"error":      err,
		}).Error("Permanent error occurred.")
		return err
	}

	// If all retries fail, return a retry failed error
	return fmt.Errorf("retry limit exceeded for item %s", item)
}

// Function to determine if an error is transient
func isTransientError(err error) bool {
	// Consider TimeoutError as transient for example
	if _, ok := err.(*TimeoutError); ok {
		return true
	}
	return false
}

// Function that applies a callback to each item, propagating errors
func ProcessItems(items []string, callback func(string) error, maxRetries int) error {
	log.Info("Starting processing of items...")
	defer log.Info("Processing completed.")

	for _, item := range items {
		err := ProcessItemWithRetry(item, callback, maxRetries)
		if err != nil {
			log.WithError(err).Error("Failed to process item")
			return err
		}
	}
	return nil
}

// Example callback function demonstrating error handling strategies
func ExampleCallback(item string) error {
	// Simulate different error scenarios based on the item name
	switch item {
	case "timeout":
		return &TimeoutError{msg: "Callback timed out"}
	case "bad-format":
		return errors.New("Item has bad format")
	case "unexpected":
		panic("Unexpected error occurred")
	}

	// If no errors, simulate successful processing
	log.WithFields(log.Fields{"item": item}).Info("Simulating successful processing")
	return nil
}

// Setup logrus with fields
func setupLogging() {
	log.SetFormatter(&log.JSONFormatter{}) // Set log format to JSON
	log.SetLevel(log.InfoLevel)             // Set log level to Info
	log.SetOutput(os.Stdout)                // Output to standard output
	log.WithField("app", "example-app").Info("Application started.")
}

func main() {
	// Set up logging
	setupLogging()

	// List of items to process
	items := []string{"item1", "bad-format", "timeout", "unexpected", "item2"}

	// Handle panic recovery
	defer func() {
		if r := recover(); r != nil {
			// Log panic errors
			log.WithField("panic", r).Panic("Recovered from panic")
		}
	}()

	// Process items with retry
	err := ProcessItems(items, ExampleCallback, 2) // Specify max retries
	if err != nil {
		// Enhanced error handling, based on the error type
		if timeoutErr, ok := err.(*TimeoutError); ok {
			// Specific handling for TimeoutError
			log.WithError(timeoutErr).Error("TimeoutError occurred")
		} else {
			// General handling for other errors
			log.WithError(err).Error("Error occurred")
		}
	}
}