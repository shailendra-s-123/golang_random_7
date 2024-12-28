package main

import (
	"fmt"
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

// Function that applies a callback to each item in a list
func ProcessItems(items []string, callback func(string) error) error {
	for _, item := range items {
		err := callback(item)
		if err != nil {
			return fmt.Errorf("error processing item %s: %w", item, err)
		}
	}
	return nil
}

// Example callback function with error handling
func ExampleCallback(item string) error {
	// Simulate an error
	if item == "error-prone" {
		return &ProcessError{message: "Item cannot be processed", cause: fmt.Errorf("internal processing error")}
	}

	// Simulate a common error
	if item == "bad-format" {
		return fmt.Errorf("item has bad format")
	}

	// Simulate an unexpected error
	if item == "unexpected" {
		panic("Unexpected error occurred")
	}

	fmt.Printf("Processing item: %s\n", item)
	return nil
}

func main() {
	items := []string{"good", "bad-format", "error-prone", "unexpected", "good"}

	err := ProcessItems(items, ExampleCallback)
	if err != nil {
		fmt.Printf("Processing failed: %v\n", err)
	} else {
		fmt.Println("Processing completed successfully.")
	}
}