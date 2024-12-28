package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus" // Import for logrus
)

// Custom error types for better error handling
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

// Function to process an item with retry logic for transient errors
func ProcessItemWithRetry(item string, callback func(string) error, maxRetries int) error {
	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		start := time.Now()

		// Log the attempt
		log.WithFields(log.Fields{
			"item":      item,
			"attempt":   attempt + 1,
			"maxRetries": maxRetries,
		}).Info("Attempting to process item...")

		// Call the callback function to process the item
		err = callback(item)
		if err == nil {
			log.WithFields(log.Fields{
				"item":     item,
				"duration": time.Since(start),
				"attempt":  attempt + 1,
				"status":   "success",
			}).Info("Item processed successfully.")
			return nil
		}

		// If the error is transient (e.g., timeout), retry the operation
		if isTransientError(err) {
			log.WithFields(log.Fields{
				"item":     item,
				"duration": time.Since(start),
				"attempt":  attempt + 1,
				"error":    err,
				"status":   "retry",
			}).Warning("Transient error occurred. Retrying...")

			// Exponential backoff between retries
			time.Sleep(time.Second * time.Duration(attempt+1))
			continue
		}

		// For permanent errors, log and return
		log.WithFields(log.Fields{
			"item":     item,
			"duration": time.Since(start),
			"attempt":  attempt + 1,
			"error":    err,
			"status":   "failure",
		}).Error("Permanent error occurred.")
		return err
	}

	// Return error if all retry attempts fail
	return fmt.Errorf("retry limit exceeded for item %s", item)
}

// Function to check if an error is transient
func isTransientError(err error) bool {
	if _, ok := err.(*TimeoutError); ok {
		return true
	}
	// Add more transient error types as needed
	return false
}

// Function that processes multiple items with retry logic
func ProcessItems(items []string, callback func(string) error, maxRetries int) error {
	log.Info("Starting processing of items...")
	defer log.Info("Processing completed.")

	for _, item := range items {
		err := ProcessItemWithRetry(item, callback, maxRetries)
		if err != nil {
			log.WithError(err).Error("Failed to process item")
			// Continue processing remaining items even if one fails
			continue
		}
	}
	return nil
}

// Example callback function to simulate item processing
func ExampleCallback(item string) error {
	// Simulate errors based on item names
	switch item {
	case "timeout":
		return &TimeoutError{msg: "Request timed out"}
	case "bad-format":
		return errors.New("Item has bad format")
	case "unexpected":
		panic("Unexpected error occurred")
	}

	// Simulate successful processing
	log.WithFields(log.Fields{"item": item}).Info("Item processed successfully")
	return nil
}

// Setup logging configuration
func setupLogging() {
	log.SetFormatter(&log.JSONFormatter{}) // Set JSON format for logs
	log.SetLevel(log.InfoLevel)             // Set log level to Info
	log.SetOutput(os.Stdout)                // Output to standard output
	log.WithField("app", "retry-example").Info("Application started.")
}

func main() {
	// Set up logging
	setupLogging()

	// List of items to process
	items := []string{"item1", "timeout", "bad-format", "unexpected", "item2"}

	// Handle panic recovery to avoid application crash
	defer func() {
		if r := recover(); r != nil {
			log.WithField("panic", r).Panic("Recovered from panic")
		}
	}()

	// Process items with retry
	err := ProcessItems(items, ExampleCallback, 3) // Retry up to 3 times
	if err != nil {
		log.WithError(err).Error("Processing items failed")
		os.Exit(1)
	}

	log.Info("All items processed successfully.")
}
