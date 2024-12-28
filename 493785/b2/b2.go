package main

import (
	"fmt"
	"time"
	"os"

	log "github.com/sirupsen/logrus"
)

// ... (Previous code remains the same)

func ProcessItems(items []string, callback func(string) error) error {
	for _, item := range items {
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
		}).Info("Item processed successfully")
	}
	return nil
}

func ExampleCallback(item string) error {
	switch item {
	case "timeout":
		return &TimeoutError{msg: "Callback timed out"}
	case "bad-format":
		return errors.New("Item has bad format")
	case "unexpected":
		panic("Unexpected error occurred")
	}

	fmt.Printf("Processing item: %s\n", item)
	return nil
}

func main() {
	// Configure logrus to log to stdout with JSON format
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	items := []string{"item1", "bad-format", "timeout", "unexpected", "item2"}

	err := ProcessItems(items, ExampleCallback)
	if err != nil {
		log.WithField("error", err).Error("Processing failed")
	} else {
		log.Info("Processing completed successfully.")
	}
}