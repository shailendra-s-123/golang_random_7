

package main

import (
	"fmt"
	"regexp"
)

// Precompiled regex patterns for better performance
var (
	emailPattern  = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	numberPattern = regexp.MustCompile(`^[-+]?\d+(\.\d+)?$`)
	datePattern   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	urlPattern    = regexp.MustCompile(`^(http|https):\/\/[a-zA-Z0-9\-]+\.[a-zA-Z0-9\-]+(?:\/[^\s]*)?$`)
)

// ValidateEmail checks if the email is in a valid format
func ValidateEmail(email string) bool {
	return emailPattern.MatchString(email)
}

// ValidateNumber checks if the input is a valid number (integer or float)
func ValidateNumber(value string) bool {
	return numberPattern.MatchString(value)
}

// ValidateDate checks if the input date is in the "YYYY-MM-DD" format
func ValidateDate(date string) bool {
	return datePattern.MatchString(date)
}

// ValidateURL checks if the input URL is valid and starts with http/https
func ValidateURL(url string) bool {
	return urlPattern.MatchString(url)
}

// Custom query parameter validation for phone number (just an example)
func ValidatePhoneNumber(phone string) bool {
	// A simple pattern for phone numbers: 10 digits, optionally with dashes or spaces
	phonePattern := regexp.MustCompile(`^(\+?\d{1,2}\s?)?(\(?\d{3}\)?[\s\-]?)?\d{3}[\s\-]?\d{4}$`)
	return phonePattern.MatchString(phone)
}

func main() {
	// Sample inputs for validation
	email := "user@example.com"
	number := "12345"
	date := "2024-12-31"
	url := "https://www.example.com"
	phone := "+1 (123) 456-7890"

	// Validate email
	if ValidateEmail(email) {
		fmt.Println("Email is valid.")
	} else {
		fmt.Println("Email is invalid.")
	}

	// Validate number
	if ValidateNumber(number) {
		fmt.Println("Number is valid.")
	} else {
		fmt.Println("Number is invalid.")
	}

	// Validate date
	if ValidateDate(date) {
		fmt.Println("Date is valid.")
	} else {
		fmt.Println("Date is invalid.")
	}

	// Validate URL
	if ValidateURL(url) {
		fmt.Println("URL is valid.")
	} else {
		fmt.Println("URL is invalid.")
	}

	// Validate phone number
	if ValidatePhoneNumber(phone) {
		fmt.Println("Phone number is valid.")
	} else {
		fmt.Println("Phone number is invalid.")
	}
}
