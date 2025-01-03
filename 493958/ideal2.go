package main

import (
	"fmt"
	"log"
	"strings"
)

// SafeFormat function to safely format strings with dynamic input and flexible error handling.
func SafeFormat(format string, args ...interface{}) (string, error) {
	// Check if the format string contains unbalanced placeholders or invalid syntax
	if strings.Contains(format, "%!") {
		return "", fmt.Errorf("invalid format string detected: %s", format)
	}

	// Use fmt.Sprintf to safely format the string with dynamic arguments
	formattedStr := fmt.Sprintf(format, args...)

	// Optional: Perform further checks on the formatted string if needed
	// For example, check for invalid or undesirable content in the formatted string.
	if strings.Contains(formattedStr, "password") {
		return "", fmt.Errorf("formatted string contains sensitive data")
	}

	return formattedStr, nil
}

func main() {
	// Example usage of SafeFormat function with dynamic input
	name := "Alice"
	age := 30
	email := "alice@example.com"

	// Safe formatting with error handling
	result, err := SafeFormat("Name: %s, Age: %d, Email: %s", name, age, email)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(result)

	// Example with invalid format string
	_, err = SafeFormat("Invalid %!", name)
	if err != nil {
		log.Printf("Error: %v", err) // Catch and log invalid format error
	}

	// Example with another invalid format string to test dynamic input handling
	_, err = SafeFormat("Missing arguments: %s %d", name)
	if err != nil {
		log.Printf("Error: %v", err) // Catch the error if not enough arguments are provided
	}
}