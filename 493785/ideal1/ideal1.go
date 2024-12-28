package main

import (
	"fmt"
	"errors"
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
	for _, item := range items {
		err := callback(item)
		if err != nil {
			// Propagating the error with context (item that failed)
			return fmt.Errorf("error processing item %s: %w", item, err)
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
	fmt.Printf("Processing item: %s\n", item)
	return nil
}

func main() {
	items := []string{"item1", "bad-format", "timeout", "unexpected", "item2"}

	err := ProcessItems(items, ExampleCallback)
	if err != nil {
		// Enhanced error handling, based on the error type
		if timeoutErr, ok := err.(*TimeoutError); ok {
			// Specific handling for TimeoutError
			fmt.Println("Error:", timeoutErr)
		} else {
			// General handling for other errors
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Processing completed successfully.")
	}
}
