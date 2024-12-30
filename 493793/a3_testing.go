package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Execute(command string, ctx context.Context) error {
	args := m.Called(command, ctx)
	return args.Error(0)
}

func TestProcessCommand(t *testing.T) {
	mockExecutor := new(MockCommandExecutor)

	tests := []struct {
		name         string
		command      string
		mockReturn  error
		expectedErr error
		expectExec  bool
	}{
		{"Valid Command", "ls", nil, nil, true},
		{"Empty Command", "", nil, errors.New("command cannot be empty"), false},
		{"Unsafe Command", "ls; rm -rf /", nil, errors.New("command contains potentially unsafe characters"), false},
		{"Execution Error", "invalid_command", errors.New("execution failed"), errors.New("execution failed"), true},
		{"Context Timeout", "", context.DeadlineExceeded, errors.New("command execution timed out"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectExec {
				mockExecutor.On("Execute", tt.command, mock.Anything).Return(tt.mockReturn)
			}

			err := ProcessCommand(tt.command, mockExecutor)

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

func BenchmarkProcessCommand(b *testing.B) {
	mockExecutor := new(MockCommandExecutor)
	mockExecutor.On("Execute", "ls", mock.Anything).Return(nil)

	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			if err := ProcessCommand("ls", mockExecutor); err != nil {
				b.Errorf("failed to process command: %v", err)
			}
		}()
	}
	wg.Wait()
	mockExecutor.AssertExpectations(b)
}

func TestProcessCommand_Stress(t *testing.T) {
	mockExecutor := new(MockCommandExecutor)
	mockExecutor.On("Execute", "ls", mock.Anything).Return(nil)

	const numIterations = 1000
	for i := 0; i < numIterations; i++ {
		assert.NoError(t, ProcessCommand("ls", mockExecutor))
	}
	mockExecutor.AssertExpectations(t)
}
