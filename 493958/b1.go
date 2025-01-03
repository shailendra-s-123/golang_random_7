package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	// Demonstrate safe string formatting using different functions
	safeStringFormattingDemo()

	// Handle user input and format the result
	handleUserInput()
}

func safeStringFormattingDemo() {
	// Use Sprintf for safe string formatting
	name := "Alice"
	age := 25
	result := fmt.Sprintf("Hello, my name is %s and I am %d years old.", name, age)
	fmt.Println(result)

	// Use Sscanf for safe parsing of user input
	var input string
	fmt.Print("Enter your name and age: ")
	_, err := fmt.Sscanf(input, "%s %d", &name, &age)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	fmt.Printf("Hello, %s! You are %d years old.\n", name, age)

	// Use Join for safe concatenation of strings
	parts := []string{"Hello", ", ", "world", "!"}
	joined := strings.Join(parts, "")
	fmt.Println(joined)
}

func handleUserInput() {
	var input string
	fmt.Print("Enter a number: ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Validate user input as a number
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input:", err)
		return
	}

	// Perform calculations and format the result
	result := fmt.Sprintf("The square of %d is %d.\n", num, num*num)
	fmt.Print(result)
}