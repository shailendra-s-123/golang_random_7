
package main

import (
	"errors"
	"testing"
	"time"
)

// CommandExecutor defines an interface for executing commands.
type CommandExecutor interface {
	Execute(command string) error
}

// processCommand processes a command string, ensuring it's safe and executable.
func processCommand(command string, executor CommandExecutor) error {
	if command == "" {
		return errors.New("command cannot be empty")
	}

	// Basic security check for unsafe characters.
	if strings.Contains(command, ";") || strings.Contains(command, "&") {
		return errors.New("command contains potentially unsafe characters")
	}

	return executor.Execute(command)
}

type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Execute(command string) error {
	args := m.Called(command)
	return args.Error(0)
}

func TestProcessCommand(t *testing.T) {
	mockExecutor := new(MockCommandExecutor)

	tests := []struct {
		name        string
		command     string
		mockReturn  error
		expectedErr error
	}{
		{"Valid Command", "ls", nil, nil},
		{"Empty Command", "", nil, errors.New("command cannot be empty")},
		{"Unsafe Command", "ls; rm -rf /", nil, errors.New("command contains potentially unsafe characters")},
		{"Command Error", "invalid_command", errors.New("execution failed"), errors.New("execution failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor.On("Execute", tt.command).Return(tt.mockReturn)
			err := processCommand(tt.command, mockExecutor)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockExecutor.AssertExpectations(t)
		})
	}
}

func TestProcessCommand_Stress(t *testing.T) {
	mockExecutor := new(MockCommandExecutor)
	mockExecutor.On("Execute", "ls").Return(nil)

	for i := 0; i < 1000; i++ {
		assert.NoError(t, processCommand("ls", mockExecutor))
	}
	mockExecutor.AssertExpectations(t)
}

func BenchmarkProcessCommand_Sequential(b *testing.B) {
	mockExecutor := new(MockCommandExecutor)
	mockExecutor.On("Execute", "ls").Return(nil)

	for i := 0; i < b.N; i++ {
		_ = processCommand("ls", mockExecutor)
	}
	mockExecutor.AssertExpectations(b)
}

func BenchmarkProcessCommand_Concurrent(b *testing.B) {
	mockExecutor := new(MockCommandExecutor)
	mockExecutor.On("Execute", "ls").Return(nil)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = processCommand("ls", mockExecutor)
		}
	})
	mockExecutor.AssertExpectations(b)
}

