// plugin_test.go
package main

import (
	"fmt"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	// Test with a normal case
	result := HelloWorld("Alice")
	if result != "Hello, Alice!" {
		t.Errorf("Expected Hello, Alice!, got %s", result)
	}

	// Test with an empty name
	result = HelloWorld("")
	if result != "Hello, !" {
		t.Errorf("Expected Hello, !, got %s", result)
	}
}