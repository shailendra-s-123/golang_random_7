
// plugin_test.go
package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommandExecutor provides a mock implementation of command execution.
type MockCommandExecutor struct {
	mock.Mock
}

// Run is a mock method that returns an error.
func (m *MockCommandExecutor) Run(cmd *exec.Cmd) error {
	args := m.Called(cmd)
	return args.Error(0)
}

// Example function to check if a string contains a substring (helper function).
func containsSubstring(str, substr string) bool {
	return strings.Contains(str, substr)
}

// TestContainsSubstring tests the containsSubstring function for different cases.
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
			assert.Equal(t, tt.expected, result, "Expected %v, but got %v", tt.expected, result)
		})
	}
}

// TestProcessCommand tests the processCommand function with the MockCommandExecutor.
func TestProcessCommand(t *testing.T) {
	tests := []struct {
		command  string
		expected string
		err      error
	}{
		{"ls -l", "", nil},
		{"invalid command", "", fmt.Errorf("failed to execute command: %v", &exec.ExitError{})},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing command '%s'", tt.command), func(t *testing.T) {
			mockExecutor := &MockCommandExecutor{}
			mockExecutor.On("Run", mock.Anything).Return(tt.err)

			err := processCommandWithExecutor(tt.command, mockExecutor)
			assert.Equal(t, tt.err, err, "Expected error %v, but got %v", tt.err, err)

			mockExecutor.AssertExpectations(t)
		})
	}
}

// BenchmarkProcessCommand benchmarks the performance of the processCommand function under varying load.
func BenchmarkProcessCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		command := "ls -l"
		_ = processCommand(command)
	}
}

// BenchmarkProcessCommandWithExecutor benchmarks the performance of the processCommand function with a mock executor under varying load.
func BenchmarkProcessCommandWithExecutor(b *testing.B) {
	mockExecutor := &MockCommandExecutor{}
	mockExecutor.On("Run", mock.Anything).Return(nil)

	for i := 0; i < b.N; i++ {
		command := "ls -l"
		_ = processCommandWithExecutor(command, mockExecutor)
	}

	mockExecutor.AssertExpectations(b)
}

// processCommand executes a given command.
func processCommand(command string) error {
	var cmd *exec.Cmd

	// Choose the appropriate shell based on the operating system.
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command) // On Windows, use cmd.exe.
	} else {
		cmd = exec.Command("sh", "-c", command) // On Unix-like systems, use sh.
	}

	// Execute the command.
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}
	return nil
}

// processCommandWithExecutor executes a given command with the provided CommandExecutor.
func processCommandWithExecutor(command string, executor CommandExecutor) error {
	var cmd *exec.Cmd

	// Choose the appropriate shell based on the operating system.
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command) // On Windows, use cmd.exe.
	} else {
		cmd = exec.Command("sh", "-c", command) // On Unix-like systems, use sh.
	}

	// Execute the command using the provided executor.
	err := executor.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}
	return nil
}

// CommandExecutor interface defines the method for executing commands.
type CommandExecutor interface {
	Run(cmd *exec.Cmd) error
}

