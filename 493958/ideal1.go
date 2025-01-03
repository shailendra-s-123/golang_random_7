package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

// User struct represents a simple user with basic details.
type User struct {
	Name  string
	Email string
	Age   int
}

// FormatUserDetails safely formats user details for output.
func FormatUserDetails(user User) (string, error) {
	if user.Name == "" || user.Email == "" {
		return "", errors.New("user name and email cannot be empty")
	}
	if user.Age <= 0 {
		return "", errors.New("user age must be greater than zero")
	}

	formatted := fmt.Sprintf("User Details:\nName: %s\nEmail: %s\nAge: %d", user.Name, user.Email, user.Age)
	return formatted, nil
}

// LogFormattedMessage safely formats a message for logging.
func LogFormattedMessage(level, message string) string {
	level = strings.ToUpper(strings.TrimSpace(level))
	if level != "INFO" && level != "WARN" && level != "ERROR" {
		level = "INFO"
	}
	return fmt.Sprintf("[%s] %s", level, message)
}

// HandleDynamicFormatting demonstrates a function that handles different dynamic contexts.
func HandleDynamicFormatting(template string, args ...interface{}) (string, error) {
	if template == "" {
		return "", errors.New("template string cannot be empty")
	}
	if !strings.Contains(template, "%") {
		return "", errors.New("template must contain placeholders for arguments")
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovering from panic during formatting: %v", r)
		}
	}()
	formatted := fmt.Sprintf(template, args...)
	return formatted, nil
}

func main() {
	// Example 1: Formatting user details
	user := User{Name: "Alice", Email: "alice@example.com", Age: 25}
	userDetails, err := FormatUserDetails(user)
	if err != nil {
		log.Fatalf("Error formatting user details: %v", err)
	}
	fmt.Println(userDetails)

	// Example 2: Logging messages
	logMessage := LogFormattedMessage("info", "Application started successfully")
	fmt.Println(logMessage)

	// Example 3: Dynamic string formatting with error handling
	template := "Hello, %s! Your account balance is $%.2f."
	dynamicMessage, err := HandleDynamicFormatting(template, "Alice", 1234.56)
	if err != nil {
		log.Fatalf("Error during dynamic formatting: %v", err)
	}
	fmt.Println(dynamicMessage)

	// Example 4: Error handling for invalid template
	invalidTemplate := "This string has no placeholders"
	_, err = HandleDynamicFormatting(invalidTemplate)
	if err != nil {
		log.Printf("Expected error: %v", err)
	}
}