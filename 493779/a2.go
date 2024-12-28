package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Middleware struct to wrap handlers
type Middleware struct {
	next      http.Handler
	sqlInject *regexp.Regexp
	xss       *regexp.Regexp
	rateLimit map[string]time.Time
	ipBlock   map[string]bool
	config    Config
	logger    Logger
}

// NewMiddleware creates a new instance of the middleware
func NewMiddleware(next http.Handler, config Config) *Middleware {
	sqlInject, err := regexp.Compile("(select|union|drop|create|alter|insert|update)")
	if err != nil {
		log.Fatalf("Error compiling SQL injection regex: %v", err)
	}
	xss, err := regexp.Compile("<[^>]+>")
	if err != nil {
		log.Fatalf("Error compiling XSS regex: %v", err)
	}
	return &Middleware{
		next:      next,
		sqlInject: sqlInject,
		xss:       xss,
		rateLimit: make(map[string]time.Time),
		ipBlock:   make(map[string]bool),
		config:    config,
		logger:    NewFileLogger(config.LogFile),
	}
}

// Function to log suspicious activity
func (m *Middleware) logSuspiciousActivity(r *http.Request, message string) {
	m.logger.Log(fmt.Sprintf("Suspicious activity detected: %s | URL: %s | IP: %s\n", message, r.URL.Path, r.RemoteAddr))
}

// Middleware to process incoming HTTP requests and filter query parameters
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.isIPBlocked(r.RemoteAddr) {
		m.handleBlockedRequest(w, r, "IP is blocked")
		return
	}

	if m.isRateLimited(r.RemoteAddr) {
		m.handleBlockedRequest(w, r, "Rate limit exceeded")
		return
	}

	queryParams := r.URL.Query()

	for key, values := range queryParams {
		for _, value := range values {
			if m.sqlInject.MatchString(strings.ToLower(value)) {
				m.logSuspiciousActivity(r, fmt.Sprintf("Potential SQL Injection in '%s' with value '%s'", key, value))
				m.handleBlockedRequest(w, r, "SQL Injection detected")
				return
			}
			if m.xss.MatchString(value) {
				m.logSuspiciousActivity(r, fmt.Sprintf("Potential XSS attack in '%s' with value '%s'", key, value))
				m.handleBlockedRequest(w, r, "XSS attack detected")
				return
			}
		}
	}

	// Update the query string in the request
	r.URL.RawQuery = queryParams.Encode()

	// Call the next handler
	m.next.ServeHTTP(w, r)
}

// Function to check if an IP is blocked
func (m *Middleware) isIPBlocked(ip string) bool {
	_, exists := m.ipBlock[ip]
	return exists
}

// Function to check if an IP has exceeded the rate limit
func (m *Middleware) isRateLimited(ip string) bool {
	now := time.Now()
	if t, exists := m.rateLimit[ip]; exists && now.Before(t) {
		return true
	}
	return false
}

// Function to handle blocked requests
func (m *Middleware) handleBlockedRequest(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, "Access denied: %s\n", message)
}

// Config struct to hold configuration options
type Config struct {
	LogFile string `json:"log_file"`
	RateLimit int   `json:"rate_limit"`
	BlockIPs []string `json:"block_ips"`
}

// Logger interface for flexible logging
type Logger interface {
	Log(message string)
}

// FileLogger implements the Logger interface to log to a file
type FileLogger struct {
	file *log.Logger
}

// NewFileLogger creates a new FileLogger
func NewFileLogger(filePath string) *FileLogger {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	return &FileLogger{file: log.New(file, "", log.LstdFlags)}
}

// Log logs a message to the file
func (f *FileLogger) Log(message string) {
	f.file.Print(message)
}

func main() {
	// Load configuration from a file or environment variables
	config := Config{
		LogFile:   "application_firewall.log",
		RateLimit: 100,
		BlockIPs:  []string{"127.0.0.1"},
	}

	// Create a router and apply the middleware
	http.Handle("/", NewMiddleware(http.HandlerFunc(myHandler), config))

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Simple handler for demonstration purposes
func myHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Request processed successfully!")
}