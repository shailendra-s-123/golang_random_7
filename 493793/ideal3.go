//main.go

package main

import (
	"errors"
	"os/exec"
	"strings"
	"sync"
)

// CommandExecutor defines the interface for executing commands.
type CommandExecutor interface {
	ExecuteCommand(name string, args ...string) (string, error)
}

// RealCommandExecutor is the production implementation of CommandExecutor.
type RealCommandExecutor struct{}

// ExecuteCommand executes the given command with arguments.
func (e *RealCommandExecutor) ExecuteCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// MockCommandExecutor is a mock implementation of CommandExecutor.
type MockCommandExecutor struct {
	mu         sync.Mutex
	callCount  int
	mockReturn map[string]struct {
		output string
		err    error
	}
}

// NewMockCommandExecutor creates a new MockCommandExecutor instance.
func NewMockCommandExecutor() *MockCommandExecutor {
	return &MockCommandExecutor{
		mockReturn: make(map[string]struct {
			output string
			err    error
		}),
	}
}

// SetMockResponse sets mock responses for a command.
func (m *MockCommandExecutor) SetMockResponse(command string, output string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mockReturn[command] = struct {
		output string
		err    error
	}{output, err}
}

// ExecuteCommand simulates the execution of a command based on predefined responses.
func (m *MockCommandExecutor) ExecuteCommand(name string, args ...string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callCount++
	cmd := strings.Join(append([]string{name}, args...), " ")
	if response, found := m.mockReturn[cmd]; found {
		return response.output, response.err
	}
	return "", errors.New("command not found in mock responses")
}