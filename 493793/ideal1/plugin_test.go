// plugin_test.go
package main

import "testing"

func TestHelloWorld(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal Input", "Alice", "Hello, Alice!"},
		{"Empty Input", "", "Hello, !"},
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
