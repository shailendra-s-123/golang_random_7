package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockShell interface {
	Exec(cmd string) error
}

type Shell struct{}

func (s *Shell) Exec(cmd string) error {
	return exec.Command("sh", "-c", cmd).Run()
}

// processes a given command using the Shell interface
func processCommand(command string, shell MockShell) error {
	return shell.Exec(command)
}

func TestProcessCommand(t *testing.T) {
	mockShell := &MockShell{}

	tests := []struct {
		command     string
		expectedErr bool
	}{
		{command: "ls", expectedErr: false},
		{command: "; ls", expectedErr: true}, // Security check for injection
		{command: "", expectedErr: true},      // Edge case: empty command
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing command: %s", tt.command), func(t *testing.T) {
			if tt.expectedErr {
				mockShell.On("Exec", tt.command).Return(fmt.Errorf("command failed"))
			} else {
				mockShell.On("Exec", tt.command).Return(nil)
			}

			err := processCommand(tt.command, mockShell)
			mockShell.AssertExpectations(t)

			if err != nil {
				if !tt.expectedErr {
					t.Errorf("Unexpected error: %v", err)
				}
			} else if tt.expectedErr {
				t.Errorf("Expected error, but got none")
			}
		})
	}
}

func BenchmarkProcessCommand(b *testing.B) {
	mockShell := &MockShell{}
	mockShell.On("Exec", "ls").Return(nil)

	for i := 0; i < b.N; i++ {
		processCommand("ls", mockShell)
	}

	mockShell.AssertExpectations(b)
}

func main() {
	testing.Main()
}