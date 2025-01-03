package main

import (
 "errors"
 "fmt"
 "log"
 "strings"
)

// User struct remains the same

// FormatUserDetails safely formats user details for output using strings.Builder.
func FormatUserDetails(user User) (string, error) {
 if user.Name == "" || user.Email == "" {
  return "", errors.New("user name and email cannot be empty")
 }
 if user.Age <= 0 {
  return "", errors.New("user age must be greater than zero")
 }

 var sb strings.Builder
 sb.WriteString("User Details:\n")
 fmt.Fprintf(&sb, "Name: %s\n", user.Name)
 fmt.Fprintf(&sb, "Email: %s\n", user.Email)
 fmt.Fprintf(&sb, "Age: %d\n", user.Age)
 return sb.String(), nil
}

// Rest of the code remains the same

func main() {
 // Example usage remains the same
}