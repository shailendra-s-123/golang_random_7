package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Custom error types
type TimeoutError struct {
	msg string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("Timeout error: %s", e.msg)
}

// A simple callback function that could fail
func callback(name string, duration time.Duration) (string, error) {
	// Simulate a timeout after the specified duration
	if duration > 0 {
		time.Sleep(duration)
	}

	if duration > 5*time.Second {
		return "", &TimeoutError{msg: "Callback timed out"}
	}
	
	if name == "foo" {
		// Simulate an unexpected error
		return "", errors.New("Unexpected error occurred")
	}

	return fmt.Sprintf("Hello, %s!", name), nil
}

// Simple example using the callback function
func main() {
	name := os.Args[1]
	duration, err := time.ParseDuration(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing duration:", err)
		os.Exit(1)
	}

	result, err := callback(name, duration)
	if err != nil {
		// Handle specific custom error types
		if timeoutErr, ok := err.(*TimeoutError); ok {
			fmt.Fprintln(os.Stderr, timeoutErr.Error())
			os.Exit(1)
		}

		// Handle unexpected errors
		fmt.Fprintln(os.Stderr, "Unexpected error:", err)
		os.Exit(1)
	}

	fmt.Println(result)
}