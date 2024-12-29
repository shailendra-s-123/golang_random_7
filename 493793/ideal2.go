
package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommandExecutor mocks the CommandExecutor interface for testing.
type MockCommandExecutor struct {
	mock.Mock
}

// Execute mocks the execution of a command.
func (m *MockCommandExecutor) Execute(command string) error {
	args := m.Called(command)
	return args.Error(0)
}

// TestProcessCommand tests processCommand for various scenarios.
func TestProcessCommand(t *testing.T) {
	mockExecutor := new(MockCommandExecutor)

	tests := []struct {
		name        string
		command     string
		mockReturn  error
		expectedErr error
		expectExec  bool
	}{
		{"Valid Command", "ls", nil, nil, true},
		{"Empty Command", "", nil, errors.New("command cannot be empty"), false},
		{"Unsafe Command", "ls; rm -rf /", nil, errors.New("command contains potentially unsafe characters"), false},
		{"Execution Error", "invalid_command", errors.New("execution failed"), errors.New("execution failed"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectExec {
				mockExecutor.On("Execute", tt.command).Return(tt.mockReturn)
			}

			err := processCommand(tt.command, mockExecutor)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectExec {
				mockExecutor.AssertExpectations(t)
			}
		})
	}
}

// BenchmarkProcessCommand benchmarks the processCommand function under load.
func BenchmarkProcessCommand(b *testing.B) {
	mockExecutor := new(MockCommandExecutor)
	mockExecutor.On("Execute", "ls").Return(nil)

	for i := 0; i < b.N; i++ {
		_ = processCommand("ls", mockExecutor)
	}
	mockExecutor.AssertExpectations(b)
}

// TestProcessCommand_Stress simulates a stress test with multiple iterations.
func TestProcessCommand_Stress(t *testing.T) {
	mockExecutor := new(MockCommandExecutor)
	mockExecutor.On("Execute", "ls").Return(nil)

	for i := 0; i < 1000; i++ {
		assert.NoError(t, processCommand("ls", mockExecutor))
	}
	mockExecutor.AssertExpectations(t)
}

