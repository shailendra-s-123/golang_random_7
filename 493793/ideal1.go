// main.go
package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// processes a given command
func processCommand(command string) error {
	var cmd *exec.Cmd

	// Choose the appropriate shell based on the operating system
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command) // On Windows, use cmd.exe
	} else {
		cmd = exec.Command("sh", "-c", command) // On Unix-like systems, use sh
	}

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}
	return nil
}

func main() {
	// Example usage of processCommand in main
	command := "ls -l" // You can replace this with other commands
	err := processCommand(command)
	if err != nil {
		fmt.Println("Error executing command:", err)
	} else {
		fmt.Println("Command executed successfully!")
	}
}


