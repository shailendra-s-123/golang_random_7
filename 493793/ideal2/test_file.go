//test file 
package main

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal Input", "Alice", "Hello, Alice!"},
		{"Empty Input", "", "Hello, !"},
		{"Long Input", "A" + "A" + "A" + "A" + "A", "Hello, AAAAA!"}, // Simplified input
		{"Unicode Input", "こんにちは", "Hello, こんにちは!"},
		{"Special Characters", "@#$%^&*()", "Hello, @#$%^&*()!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HelloWorld(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
