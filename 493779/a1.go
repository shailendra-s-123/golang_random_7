package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Define a middleware structure
type Middleware struct {
	next http.Handler
}

// Function to create a middleware instance
func NewMiddleware(next http.Handler) *Middleware {
	return &Middleware{next}
}

// Helper function to escape SQL characters
func escapeSQL(input string) string {
	input = strings.Replace(input, "'", "''", -1)
	input = strings.Replace(input, "\\", "\\\\", -1)
	return input
}

// Helper function to escape HTML characters
func escapeHTML(input string) string {
	input = strings.Replace(input, "&", "&amp;", -1)
	input = strings.Replace(input, "<", "&lt;", -1)
	input = strings.Replace(input, ">", "&gt;", -1)
	input = strings.Replace(input, "'", "&#039;", -1)
	input = strings.Replace(input, "\"", "&quot;", -1)
	return input
}

// Handler function that will be processed by the middleware
func myHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, world!")
}

// Middleware function to filter query parameters
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters
	queryParams := r.URL.Query()
	
	// Log query parameters
	fmt.Printf("Received query parameters: %v\n", queryParams)
	
	// Check and filter for SQL Injection
	for key, values := range queryParams {
		for i, value := range values {
			// Check if the value contains SQL-like characters
			if strings.Contains(value, "'") || strings.Contains(value, "\\") {
				fmt.Printf("Suspicious SQL Injection detected in parameter '%s' with value '%s'\n", key, value)
				// Escape the value
				queryParams.Set(key, escapeSQL(value))
			}
		}
	}
	
	// Check and filter for XSS
	for key, values := range queryParams {
		for i, value := range values {
			// Check if the value contains HTML-like characters
			if strings.Contains(value, "<") || strings.Contains(value, ">") || strings.Contains(value, "&") {
				fmt.Printf("Suspicious XSS detected in parameter '%s' with value '%s'\n", key, value)
				// Escape the value
				queryParams.Set(key, escapeHTML(value))
			}
		}
	}
	
	// Update the query string in the request
	r.URL.RawQuery = queryParams.Encode()
	
	// Call the next handler
	m.next.ServeHTTP(w, r)
}

func main() {
	// Create a middleware instance
	middleware := NewMiddleware(http.HandlerFunc(myHandler))
	
	// Define the HTTP server
	http.Handle("/", middleware)
	
	// Start the server
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}