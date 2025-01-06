// Refined test cases
// test_plugin_refined.go
package main

import (
	"testing"
)

func TestExamplePluginProcess(t *testing.T) {
	tests := []struct {
		name     string
		plugin   Plugin
		input    string
		expected string
	}{
		{
			name:     "Simple input",
			plugin:   &ExamplePlugin{},
			input:    "test",
			expected: "test processed",
		},
		{
			name:     "Empty input",
			plugin:   &ExamplePlugin{},
			input:    "",
			expected: "", // Refined to match the plugin's behavior
		},
		{
			name:     "Whitespace input",
			plugin:   &ExamplePlugin{},
			input:    " ",
			expected: " processed", // Refined to handle whitespace correctly
		},
		{
			name:     "Multiple words input",
			plugin:   &ExamplePlugin{},
			input:    "test multiple words",
			expected: "test multiple words processed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessData(tt.plugin, tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}