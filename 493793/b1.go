package main

import (
	"fmt"
	"testing"
)

// Sample plugin function that checks if a string contains a substring
func containsSubstring(text, substring string) bool {
	return strings.Contains(text, substring)
}

func TestContainsSubstring(t *testing.T) {
	testCases := []struct {
		text       string
		substring  string
		want       bool
	}{
		{text: "Hello, world!", substring: "world", want: true},
		{text: "Hello, world!", substring: "python", want: false},
		{text: "", substring: "", want: true},
		// Based on feedback from a workshop: Test for empty substring case
		{text: "Hello, world!", substring: "", want: true},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf("text: %q, substring: %q", tt.text, tt.substring), func(t *testing.T) {
			got := containsSubstring(tt.text, tt.substring)
			if got != tt.want {
				t.Errorf("containsSubstring(%q, %q) = %v, want %v", tt.text, tt.substring, got, tt.want)
			}
		})
	}
}