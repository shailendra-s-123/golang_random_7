package main

import (
	"errors"
	"fmt"
	"log"
	"text/template"
)

// User struct represents a simple user with basic details.
type User struct {
	Name  string
	Email string
	Age   int
}

// FormatUserDetails uses a template to format user details for output.
func FormatUserDetails(user User) (string, error) {
	if user.Name == "" || user.Email == "" {
		return "", errors.New("user name and email cannot be empty")
	}
	if user.Age <= 0 {
		return "", errors.New("user age must be greater than zero")
	}

	templateStr := `User Details:
Name: {{ .Name }}
Email: {{ .Email }}
Age: {{ .Age }}`

	tmpl, err := template.New("userDetails").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var output string
	err = tmpl.Execute(&output, user)
	if err != nil {
		return "", err
	}

	return output, nil
}

// LogFormattedMessage uses a template to format a message for logging.
func LogFormattedMessage(level, message string) string {
	level = strings.ToUpper(strings.TrimSpace(level))
	if level != "INFO" && level != "WARN" && level != "ERROR" {
		level = "INFO"
	}

	templateStr := `[{{ .Level }}] {{ .Message }}`

	tmpl, err := template.New("logMessage").Parse(templateStr)
	if err != nil {
		log.Fatalf("Error parsing log message template: %v", err)
	}

	var output string
	err = tmpl.Execute(&output, map[string]string{"Level": level, "Message": message})
	if err != nil {
		log.Fatalf("Error executing log message template: %v", err)
	}

	return output
}

// HandleDynamicFormatting uses a template to handle dynamic contexts.
func HandleDynamicFormatting(templateStr string, data interface{}) (string, error) {
	if templateStr == "" {
		return "", errors.New("template string cannot be empty")
	}

	tmpl, err := template.New("dynamicFormatting").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var output string
	err = tmpl.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output, nil
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
	dynamicTemplate := "Hello, {{ .Name }}! Your account balance is ${{ .Balance | printf '%.2f' }}."
	dynamicData := map[string]interface{}{"Name": "Alice", "Balance": 1234.56}
	dynamicMessage, err := HandleDynamicFormatting(dynamicTemplate, dynamicData)
	if err != nil {
		log.Fatalf("Error during dynamic formatting: %v", err)
	}
	fmt.Println(dynamicMessage)

	// Example 4: Error handling for invalid template
	invalidTemplate := "This string has no placeholders"
	_, err = HandleDynamicFormatting(invalidTemplate, map[string]interface{}{})
	if err != nil {
		log.Printf("Expected error: %v", err)
	}
}