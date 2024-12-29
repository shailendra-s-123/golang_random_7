package main

import (
	"testing"
	"io/ioutil"
	"os/exec"
	"bytes"
)

func TestCommandProcessor(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect error
	}{
		{
			name:   "normal usage",
			input:  "ls -l",
			expect: nil,
		},
		{
			name:   "command injection",
			input:  "ls; rm -rf /",
			expect: errCommandInjection,
		},
		{
			name:   "unexpected command",
			input:  "non-existing-command",
			expect: errCommandNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			cmd := exec.Command("sh", "-c", test.input)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			if err != test.expect {
				t.Errorf("got error %v, want %v\nstdout: %s\nstderr: %s", err, test.expect, stdout.String(), stderr.String())
			}
		})
	}
}

var errCommandInjection = errors.New("command injection detected")
var errCommandNotFound = errors.New("command not found")