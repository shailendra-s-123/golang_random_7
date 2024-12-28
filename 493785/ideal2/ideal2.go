package main

import (
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

// Function that applies a callback to each item, propagating errors
func ProcessItems(items []string, callback func(string) error) error {
	log.Info("Starting processing of items...")
	defer log.Info("Processing completed.")

	for _, item := range items {
		log.WithFields(log.Fields{"item": item}).Info("Processing item...")

		start := time.Now()
		err := callback(item)
		if err != nil {
			log.WithFields(log.Fields{
				"item":     item,
				"duration": time.Since(start),
				"error":    err,
			}).Error("Error processing item")
			return fmt.Errorf("error processing item %s: %w", item, err)
		}

		log.WithFields(log.Fields{
			"item":     item,
			"duration": time.Since(start),
		}).Info("Item processed successfully.")
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

	// Process items
	err := ProcessItems(items, ExampleCallback)
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
