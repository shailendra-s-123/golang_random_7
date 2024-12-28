package main

import (
	"fmt"
	"log"
)

type Formatter interface {
	Format(template string, args ...interface{}) (string, error)
}

type SafeFormatter struct{}

func (sf *SafeFormatter) Format(template string, args ...interface{}) (string, error) {
	_, err := fmt.Printf(template, args...)
	if err != nil {
		return "", err
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

func main() {
	formatter := &SafeFormatter{}
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