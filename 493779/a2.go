package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

// ValidateEmail uses a regular expression to validate email format
func ValidateEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

// ValidateNumber uses a regular expression to validate numeric fields
func ValidateNumber(value string) bool {
	const numberRegex = `^[-+]?\d+(\.\d+)?$`
	match, _ := regexp.MatchString(numberRegex, value)
	return match
}

func main() {
	// Example usage
	email := "user@example.com"
	number := "123.45"

	if ValidateEmail(email) {
		fmt.Println("Email is valid.")
	} else {
		fmt.Println("Email is invalid.")
	}

	if ValidateNumber(number) {
		fmt.Println("Number is valid.")
	} else {
		fmt.Println("Number is invalid.")
	}
}