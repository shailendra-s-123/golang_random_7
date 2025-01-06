package main

import (
	"errors"
	"fmt"
)

// ErrorInterface defines the interface that all custom errors must implement.
type ErrorInterface interface {
	Error() string
	Code() int
	Cause() error
}

// CustomError is a struct that implements ErrorInterface.
type CustomError struct {
	msg    string
	code   int
	cause  error
}

// Error implements the Error interface.
func (e *CustomError) Error() string {
	return fmt.Sprintf("error code: %d, message: %s", e.code, e.msg)
}

// Code returns the error code.
func (e *CustomError) Code() int {
	return e.code
}

// Cause returns the underlying cause of the error.
func (e *CustomError) Cause() error {
	return e.cause
}

// NewCustomError creates a new CustomError with the given message, code, and cause.
func NewCustomError(msg string, code int, cause error) error {
	return &CustomError{msg: msg, code: code, cause: cause}
}

// Is checks if the error implements ErrorInterface and matches the given code.
func Is(err error, code int) bool {
	if ce, ok := err.(ErrorInterface); ok {
		return ce.Code() == code
	}
	return false
}

func main() {
	// Create a custom error
	err := NewCustomError("Invalid input", 400, errors.New("data is missing"))

	// Handle the error
	switch err {
	case nil:
		fmt.Println("No error")
	default:
		if Is(err, 400) {
			fmt.Println("Invalid input error occurred.")
		} else {
			fmt.Printf("Unknown error: %v\n", err)
		}
	}
}