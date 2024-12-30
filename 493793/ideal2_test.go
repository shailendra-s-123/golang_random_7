


package main

import (
	"errors"
	"strings"
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

