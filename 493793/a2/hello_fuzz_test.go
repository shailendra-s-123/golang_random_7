// hello_fuzz_test.go
package main

import (
	"fmt"
	"testing"
)

func FuzzHelloWorld(f *testing.F) {
	for _, tc := range []struct {
		name string
	}{
		{"Normal Input"},
		{"Empty Input"},
	} {
		f.Add(tc.name)
	}
	f.Fuzz(func(f *testing.F, b []byte) {
		input := string(b)
		HelloWorld(input)
	})
}