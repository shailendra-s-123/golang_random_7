package main

import (
	"fmt"
	"log"
	"os"
	"plugin"
	"strings"
)

// Plugin interface defines the methods that plugins must implement.
type Plugin interface {
	Run() (string, error)
}

// loadPlugin loads a plugin from the specified path.
func loadPlugin(path string) (Plugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %v", err)
	}

	sym, err := p.Lookup("Run")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup symbol: %v", err)
	}

	plugin, ok := sym.(func() (string, error))
	if !ok {
		return nil, fmt.Errorf("symbol is not a function")
	}

	return &pluginWrapper{plugin}, nil
}

type pluginWrapper struct {
	run func() (string, error)
}

func (pw *pluginWrapper) Run() (string, error) {
	return pw.run()
}

// runTests executes the provided test cases.
func runTests(plugin Plugin, testCases []string) {
	for _, testCase := range testCases {
		output, err := plugin.Run()
		if err != nil {
			fmt.Printf("Test '%s' failed: %v\n", testCase, err)
		} else {
			fmt.Printf("Test '%s' passed: %s\n", testCase, output)
		}
	}
}

// main function handles plugin loading and testing.
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run test_plugin.go <plugin_path> [test_case1 test_case2...]")
	}

	pluginPath := os.Args[1]
	testCases := os.Args[2:]

	plugin, err := loadPlugin(pluginPath)
	if err != nil {
		log.Fatal(err)
	}

	// Default test cases
	defaultTestCases := []string{
		"functional_test_1",
		"functional_test_2",
		"security_test_1",
	}

	// Combine default and user-provided test cases
	allTestCases := append(defaultTestCases, testCases...)

	runTests(plugin, allTestCases)
}