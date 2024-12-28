package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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
	auditLog  *AuditLog
	lock      sync.Mutex
}

// Config struct to hold configuration options
type Config struct {
	LogFile   string   `json:"log_file"`
	RateLimit int      `json:"rate_limit"` // Max requests per minute
	BlockIPs  []string `json:"block_ips"`  // List of blocked IPs
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
	f.file.Println(message)
}

// AuditLog struct to hold audit log entries
type AuditLog struct {
	entries []AuditLogEntry
	lock    sync.Mutex
}

// AuditLogEntry struct to hold an audit log entry
type AuditLogEntry struct {
	Time   time.Time `json:"time"`
	Action string    `json:"action"`
	User   string    `json:"user"`   // To be implemented: User authentication and authorization
	Data   string    `json:"data"`   // Data changed (e.g., rule name, IP address)
}

// NewAuditLog creates a new AuditLog
func NewAuditLog() *AuditLog {
	return &AuditLog{entries: []AuditLogEntry{}}
}

// Log logs an audit log entry
func (a *AuditLog) Log(entry AuditLogEntry) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.entries = append(a.entries, entry)
}

// Export exports the audit log entries to a JSON file
func (a *AuditLog) Export(filePath string) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	data, err := json.MarshalIndent(a.entries, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, data, 0644)
}

// NewMiddleware creates a new instance of the middleware
func NewMiddleware(next http.Handler, config Config, auditLog *AuditLog) *Middleware {
	sqlInject, err := regexp.Compile(`(?i)(select|insert|update|delete|drop|union|exec|create|alter|truncate|declare)`)
	if err != nil {
		log.Fatalf("Error compiling SQL injection regex: %v", err)
	}
	xss, err := regexp.Compile(`(?i)<script.*?>.*?</script>`)
	if err != nil {
		log.Fatalf("Error compiling XSS regex: %v", err)
	}
	return &Middleware{
		next:      next,
		sqlInject: sqlInject,
		xss:       xss,