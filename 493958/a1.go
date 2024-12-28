package main

import (
	"fmt"
	"log"
)

// Define a simple interface for our service
type Formatter interface {
	Format(template string, args ...interface{}) string
}

// A struct that implements the Formatter interface
type SprintfFormatter struct {
}

// Implements the Format method using fmt.Sprintf
func (sf SprintfFormatter) Format(template string, args ...interface{}) string {
	formatted, err := fmt.Sprintf(template, args...)
	if err != nil {
		// Handle formatting errors
		log.Fatalf("Error formatting string: %v", err)
	}
	return formatted
}

// A struct that represents a service using the Formatter
type LoggerService struct {
	formatter Formatter
}

// Constructor for LoggerService that takes a Formatter as a dependency
func NewLoggerService(formatter Formatter) *LoggerService {
	return &LoggerService{
		formatter: formatter,
	}
}

// Logs a message using the provided formatter
func (ls *LoggerService) Log(template string, args ...interface{}) {
	formatted := ls.formatter.Format(template, args...)
	fmt.Println("Log:", formatted)
}

// Main function to demonstrate the application
func main() {
	// Create an instance of SprintfFormatter
	formatter := SprintfFormatter{}

	// Create an instance of LoggerService with SprintfFormatter as a dependency
	loggerService := NewLoggerService(formatter)

	// Log a message using the logger service
	loggerService.Log("User %s has %d posts", "Alice", 10)

	// Attempt to log a message with incorrect argument types
	loggerService.Log("User %d has %s posts", "Bob", 20)
}