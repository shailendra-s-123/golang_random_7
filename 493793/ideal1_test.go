
// main_test.go
package main

import (
	"fmt"
	"strings"
	"testing"
)

// Example function to check if a string contains a substring (helper function)
func containsSubstring(str, substr string) bool {
	return strings.Contains(str, substr)
}

// Test function to test containsSubstring
func TestContainsSubstring(t *testing.T) {
	tests := []struct {
		str      string
		substr   string
		expected bool
	}{
		{"Hello World", "World", true},
		{"Go is awesome", "awesome", true},
		{"Testing", "Go", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing if '%s' contains '%s'", tt.str, tt.substr), func(t *testing.T) {
			result := containsSubstring(tt.str, tt.substr)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

