package main

import "fmt"

func Run() (string, error) {
	switch strings.ToLower(os.Args[1]) {
	case "functional_test_1":
		return "Functional Test 1 Passed", nil
	case "functional_test_2":
		return "Functional Test 2 Passed", nil
	case "security_test_1":
		// Simulate a security vulnerability
		return "Security Test 1 Failed: Input validation error", fmt.Errorf("invalid input")
	default:
		return "Test not found", nil
	}
}