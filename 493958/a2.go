package main

import (
	"fmt"
	"log"
	"testing"
)

type Formatter interface {
	Format(template string, args ...interface{}) (string, error)
}

type SafeFormatter struct {
	logger *log.Logger
}

func NewSafeFormatter(logger *log.Logger) *SafeFormatter {
	return &SafeFormatter{logger: logger}
}

func (sf *SafeFormatter) Format(template string, args ...interface{}) (string, error) {
	_, err := fmt.Printf(template, args...)
	if err != nil {
		sf.logger.Printf("Error in formatting: %s\n", err)
		return "", fmt.Errorf("format error: %v", err)
	}
	return fmt.Sprintf(template, args...), nil
}

type User struct {
	ID   int
	Name string
}

type GreetingService struct {
	formatter Formatter
}

func NewGreetingService(formatter Formatter) *GreetingService {
	return &GreetingService{formatter: formatter}
}

func (gs *GreetingService) GreetUser(user *User) (string, error) {
	return gs.formatter.Format("Hello, %s (ID: %d)!", user.Name, user.ID)
}

// MockFormatter for testing purposes
type MockFormatter struct {
	formatCalled int
	mockError    error
}

func (mf *MockFormatter) Format(template string, args ...interface{}) (string, error) {
	mf.formatCalled++
	return "", mf.mockError
}

func TestGreetingService(t *testing.T) {
	formatter := &SafeFormatter{logger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)}
	gs := NewGreetingService(formatter)

	user := &User{ID: 1, Name: "Alice"}
	_, err := gs.GreetUser(user)
	if err != nil {
		t.Errorf("GreetUser failed: %v", err)
	}

	// Mock test: simulate error during formatting
	mockFormatter := &MockFormatter{mockError: fmt.Errorf("mock error")}
	gsMock := NewGreetingService(mockFormatter)
	_, errMock := gsMock.GreetUser(user)
	if errMock == nil || errMock.Error() != "mock error" {
		t.Errorf("Expected mock error: %v, got %v", "mock error", errMock)
	}
	if mockFormatter.formatCalled != 1 {
		t.Errorf("Expected format to be called once, got %d", mockFormatter.formatCalled)
	}
}

func main() {
	formatter := &SafeFormatter{logger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)}
	gs := NewGreetingService(formatter)

	user := &User{ID: 1, Name: "Alice"}
	greeting, err := gs.GreetUser(user)
	if err != nil {
		log.Fatalf("Error greeting user: %v", err)
	}
	fmt.Println(greeting)

	user2 := &User{ID: 2, Name: "Bob"}
	greeting2, err := gs.GreetUser(user2)
	if err != nil {
		log.Fatalf("Error greeting user: %v", err)
	}
	fmt.Println(greeting2)
}