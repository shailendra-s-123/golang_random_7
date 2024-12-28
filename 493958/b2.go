package main

import (
	"fmt"
	"log"
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
		sf.logger.Printf("Error dry-running format: %v", err)
		return "", err
	}
	return fmt.Sprintf(template, args...), nil
}

type User struct {
	ID   int
	Name string
}

type Account struct {
	Balance float64
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

func (gs *GreetingService) UserAccountSummary(user *User, account *Account) (string, error) {
	return gs.formatter.Format("User %s (ID: %d) has an account balance of %.2f", user.Name, user.ID, account.Balance)
}

type MockFormatter struct {
	errors chan error
}

func NewMockFormatter() *MockFormatter {
	return &MockFormatter{errors: make(chan error, 1)}
}

func (mf *MockFormatter) Format(template string, args ...interface{}) (string, error) {
	if len(mf.errors) > 0 {
		return "", <-mf.errors
	}
	return fmt.Sprintf(template, args...), nil
}

func (mf *MockFormatter) InjectError(err error) {
	mf.errors <- err
}

func main() {
	mainLogger := log.New(log.Writer(), "", log.LstdFlags)
	formatter := NewSafeFormatter(mainLogger)
	gs := NewGreetingService(formatter)

	user := &User{ID: 1, Name: "Alice"}
	account := &Account{Balance: 100.50}

	greeting, err := gs.GreetUser(user)
	if err != nil {
		log.Fatalf("Error greeting user: %v", err)
	}
	fmt.Println(greeting)

	summary, err := gs.UserAccountSummary(user, account)
	if err != nil {
		log.Fatalf("Error generating account summary: %v", err)
	}
	fmt.Println(summary)
	mockFormatter := NewMockFormatter()
	mockGs := NewGreetingService(mockFormatter)

	mockFormatter.InjectError(fmt.Errorf("invalid type for format specifier"))
	_, err = mockGs.GreetUser(user)
	if err != nil {
		log.Printf("Expected error: %v", err)
	}
}