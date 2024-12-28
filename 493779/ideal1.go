package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// Middleware structure to wrap handlers
type Middleware struct {
	next http.Handler
}

// NewMiddleware creates a new instance of the middleware
func NewMiddleware(next http.Handler) *Middleware {
	return &Middleware{next}
}

// Function to escape SQL-related characters
func escapeSQL(input string) string {
	input = strings.Replace(input, "'", "''", -1) // Escape single quotes for SQL
	input = strings.Replace(input, "\\", "\\\\", -1) // Escape backslashes
	return input
}

// Function to sanitize HTML content to prevent XSS
func sanitizeHTML(input string) string {
	p := bluemonday.NewPolicy()
	return p.Sanitize(input) // Remove dangerous HTML tags and attributes
}

// Log suspicious activity
func logSuspiciousActivity(r *http.Request, message string) {
	log.Printf("Suspicious activity detected: %s | URL: %s | IP: %s\n", message, r.URL.Path, r.RemoteAddr)
}

// Middleware to process incoming HTTP requests and filter query parameters
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	queryParams := r.URL.Query()

	// Loop through query parameters to check for potential SQL injection or XSS
	for key, values := range queryParams {
		for _, value := range values {
			// Check for SQL injection patterns
			if strings.Contains(value, "'") || strings.Contains(value, "\\") {
				escapedValue := escapeSQL(value)
				queryParams.Set(key, escapedValue)
				if escapedValue != value {
					logSuspiciousActivity(r, fmt.Sprintf("Potential SQL Injection in '%s' with value '%s'", key, value))
				}
			}

			// Check for XSS patterns
			if strings.Contains(value, "<") || strings.Contains(value, ">") || strings.Contains(value, "&") {
				sanitizedValue := sanitizeHTML(value)
				queryParams.Set(key, sanitizedValue)
				if sanitizedValue != value {
					logSuspiciousActivity(r, fmt.Sprintf("Potential XSS attack in '%s' with value '%s'", key, value))
				}
			}
		}
	}

	// Update the query string in the request
	r.URL.RawQuery = queryParams.Encode()

	// Call the next handler
	m.next.ServeHTTP(w, r)
}

// Simple handler for demonstration purposes
func myHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Request processed successfully!")
}

func main() {
	// Create a router and apply the middleware
	http.Handle("/", NewMiddleware(http.HandlerFunc(myHandler)))

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}