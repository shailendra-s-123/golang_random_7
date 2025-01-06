// plugin_test.go
package main

import (
	"fuzz"
	"testing"
)

func TestValidateEmailFuzz(t *testing.T) {
	f := fuzz.New()
	for i := 0; i < 1000; i++ {
		var email string
		f.Fuzz(&email)
		err := ValidateEmail(email)
		if err != nil && err.Error() != "Invalid email format" {
			t.Errorf("Expected 'Invalid email format', got %s for email %s", err.Error(), email)
		}
	}
}
