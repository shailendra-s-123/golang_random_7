package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type ProcessError struct {
	message string
	cause   error
}

func (e *ProcessError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("ProcessError: %s (cause: %v)", e.message, e.cause)
	}
	return fmt.Sprintf("ProcessError: %s", e.message)
}

func ProcessItems(ctx context.Context, items []string, callback func(context.Context, string) error) error {
	for _, item := range items {
		err := callback(ctx, item)
		if err != nil {
			return fmt.Errorf("error processing item %s: %w", item, err)
		}
	}
	return nil
}

func ExampleCallback(ctx context.Context, item string) error {
	// Simulate different types of errors
	switch item {
	case "temporary-error":
		// Simulate a temporary error that can be retried
		return &ProcessError{message: "Temporary error occurred", cause: errors.New("retryable error")}
	case "permanent-error":
		// Simulate a permanent error that cannot be retried
		return &ProcessError{message: "Permanent error occurred", cause: errors.New("non-retryable error")}
	case "unexpected-error":
		// Simulate an unexpected error
		panic("Unexpected error occurred")
	default:
		fmt.Printf("Processing item: %s\n", item)
		return nil
	}
}

func retry(ctx context.Context, fn func(context.Context) error, attempts int, delay time.Duration) error {
	for i := 0; i < attempts; i++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			log.Printf("Retrying in %v: %v", delay, err)
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("exhausted retry attempts")
}

func main() {
	items := []string{"good", "temporary-error", "permanent-error", "unexpected-error", "good"}

	// Set up a logger with structured output
	logger := log.New(log.Writer(), "", log.LstdFlags|log.Lshortfile)

	// Wrap the ExampleCallback with a retry mechanism
	retryableCallback := func(ctx context.Context, item string) error {
		return retry(ctx, func(ctx context.Context) error {
			return ExampleCallback(ctx, item)
		}, 3, 1*time.Second)
	}

	// Process items with retry
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := ProcessItems(ctx, items, retryableCallback)
	if err != nil {
		logger.Printf("Processing failed: %v", err)
	} else {
		logger.Println("Processing completed successfully.")
	}
}