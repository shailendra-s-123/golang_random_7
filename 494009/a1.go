package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Example of using os.Create and defer to close the file automatically.
	filePath := "example.txt"
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close() // Ensure the file is closed even if there's an error

	_, err = f.WriteString("Hello, world!")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File created successfully:", filePath)

	// Example of using exec.Command and defer to wait for the command to finish
	// and handle its output/error streams.
	cmd := exec.Command("echo", "Hello from the shell!")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Println("Shell command output:", string(out))

	// Demonstrate nested defer and resource cleanup with dependencies
	fmt.Println("Starting nested defer example...")
	func nestedDeferExample() {
		fmt.Println("Opening file...")
		f, err := os.Create("nested_example.txt")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer f.Close()

		fmt.Println("Writing to file...")
		_, err = f.WriteString("Nested example content.")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Println("File closed successfully.")
	}
	nestedDeferExample()

	fmt.Println("All resources cleaned up successfully.")
}