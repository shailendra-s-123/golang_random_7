// main.go
package main

import (
	"fmt"
	"testing"
)

type Plugin interface {
	Process(data string) string
}

func ProcessData(p Plugin, data string) string {
	return p.Process(data)
}

// Example Plugin implementation
type ExamplePlugin struct {
}

func (p *ExamplePlugin) Process(data string) string {
	// Add more processing logic here
	return data + " processed"
}

// Example test file
// test_plugin.go
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