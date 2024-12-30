//main_test.go
package main

import (
	"errors"
	"strings"
	"testing"
)

// TestCommandExecution tests the execution of commands using MockCommandExecutor.
func TestCommandExecution(t *testing.T) {
	mock := NewMockCommandExecutor()
	mock.SetMockResponse("echo hello", "hello\n", nil)
	mock.SetMockResponse("ls invalidDir", "", errors.New("directory not found"))

	tests := []struct {
		name       string
		command    string
		args       []string
		wantOutput string
		wantErr    bool
	}{
		{"ValidCommand", "echo", []string{"hello"}, "hello", false}, // Remove the newline from wantOutput
		{"InvalidCommand", "ls", []string{"invalidDir"}, "", true},
		{"EmptyCommand", "", []string{}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := mock.ExecuteCommand(tt.command, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, want error: %v", err, tt.wantErr)
			}
			if strings.TrimSpace(output) != tt.wantOutput { // Trim newline characters from output
				t.Errorf("got output %q, want output %q", output, tt.wantOutput)
			}
		})
	}
}

// BenchmarkCommandExecution benchmarks single command executions.
func BenchmarkCommandExecution(b *testing.B) {
	mock := NewMockCommandExecutor()
	mock.SetMockResponse("echo stress", "stress\n", nil)

	for i := 0; i < b.N; i++ {
		_, err := mock.ExecuteCommand("echo", "stress")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkConcurrentExecutions benchmarks concurrent command executions.
func BenchmarkConcurrentExecutions(b *testing.B) {
	mock := NewMockCommandExecutor()
	mock.SetMockResponse("echo loadtest", "loadtest\n", nil)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := mock.ExecuteCommand("echo", "loadtest")
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})
}

