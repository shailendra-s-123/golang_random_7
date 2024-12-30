package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CommandExecutor interface defines the operations for executing commands.
type CommandExecutor interface {
	Execute(command string, ctx context.Context) error
}

// ProcessCommand executes a command with basic safety checks.
func ProcessCommand(command string, executor CommandExecutor) error {
	if command == "" {
		return errors.New("command cannot be empty")
	}

	if strings.Contains(command, ";") || strings.Contains(command, "&") {
		return errors.New("command contains potentially unsafe characters")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return executor.Execute(command, ctx)
}
