package main

import "fmt"

// Error is a common interface for errors in the system.
type Error interface {
	Error() string
	Code() int
}

// CustomError is a custom error type that conforms to the Error interface.
type CustomError struct {
	message string
	code    int
}

func (e *CustomError) Error() string {
	return e.message
}

func (e *CustomError) Code() int {
	return e.code
}

func main() {
	// Function that may return an error.
	doSomething := func() error {
		// Simulate an error.
		return &CustomError{
			message: "Something went wrong",
			code:    500,
		}
	}

	// Call the function and handle the error.
	if err := doSomething(); err != nil {
		// Check the error type and code.
		if customError, ok := err.(*CustomError); ok {
			fmt.Printf("Error: %s (Code: %d)\n", customError.Error(), customError.Code())
		} else {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}
}