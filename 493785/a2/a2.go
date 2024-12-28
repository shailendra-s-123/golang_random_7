package main

import (
	"fmt"
	"errors"
	"os"
	"github.com/sirupsen/logrus"
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
	logrus.Info("Starting processing of items...")
	defer logrus.Info("Processing completed.")

	for _, item := range items {
		logrus.WithFields(logrus.Fields{"item": item}).Info("Processing item...")

		err := callback(item)
		if err != nil {
			// Propagating the error with context (item that failed)
			return fmt.Errorf("error processing item %s: %w", item, err)
		}

		logrus.WithFields(logrus.Fields{"item": item}).Info("Item processed successfully.")
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
	logrus.WithFields(logrus.Fields{"item": item}).Info("Simulating successful processing")
	return nil
}

// Setup logrus with fields
func setupLogging() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
	logrus.WithField("app", "example-app").Info("Application started.")
}

func main() {
	setupLogging()

	items := []string{"item1", "bad-format", "timeout", "unexpected", "item2"}

	defer func() {
		if r := recover(); r != nil {
			// Log panic errors
			logrus.WithField("panic", r).Panic("Recovered from panic")
		}
	}()

	err := ProcessItems(items, ExampleCallback)
	if err != nil {
		// Enhanced error handling, based on the error type
		if timeoutErr, ok := err.(*TimeoutError); ok {
			// Specific handling for TimeoutError
			logrus.WithError(timeoutErr).Error("TimeoutError occurred")
		} else {
			// General handling for other errors
			logrus.WithError(err).Error("Error occurred")
		}
	}
}