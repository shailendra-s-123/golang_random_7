package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Struct to hold user information
type User struct {
	Name    string
	Age     int
	Email   string
	Phone   string
	Address string
}

// Function to safely format a string with user information
func formatUserInfo(user User) string {
	// Using fmt.Printf to format the string safely
	return fmt.Sprintf("Name: %s\nAge: %d\nEmail: %s\nPhone: %s\nAddress: %s\n",
		user.Name,
		user.Age,
		user.Email,
		user.Phone,
		user.Address,
	)
}

// Function to format a string with an optional argument
func formatOptional(value interface{}) string {
	if value == nil {
		return "No value provided"
	}
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("Value: %s", v)
	case int:
		return fmt.Sprintf("Value: %d", v)
	case float64:
		return fmt.Sprintf("Value: %.2f", v)
	default:
		return "Unsupported type"
	}
}

// Function to format a string with error handling
func formatWithError(input string) (string, error) {
	// Simulate a situation where formatting could fail
	if strings.Contains(input, "error") {
		return "", fmt.Errorf("Input contains error: %s", input)
	}
	return fmt.Sprintf("Formatted string: %s", input), nil
}

// Function to format a number with precision
func formatNumber(number float64) string {
	return fmt.Sprintf("%.2f", number)
}

func main() {
	user := User{
		Name:    "John Doe",
		Age:     30,
		Email:   "johndoe@example.com",
		Phone:   "+123456789",
		Address: "123 Main St, Anytown, USA",
	}

	// Format user information
	formattedUserInfo := formatUserInfo(user)
	fmt.Println(formattedUserInfo)

	// Format an optional argument
	fmt.Println(formatOptional("Hello, World!")) // Should print Value: Hello, World!
	fmt.Println(formatOptional(42))              // Should print Value: 42
	fmt.Println(formatOptional(3.14159))         // Should print Value: 3.14

	// Format a string with error handling
	inputString := "This is a test"
	formattedString, err := formatWithError(inputString)
	if err != nil {
		fmt.Println("Error formatting string:", err)
	} else {
		fmt.Println(formattedString)
	}

	// Format a number with precision
	number := 3.14159
	fmt.Println(formatNumber(number)) // Should print 3.14
}