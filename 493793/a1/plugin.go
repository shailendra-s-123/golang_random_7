// plugin.go
package main

import (
	"fmt"
)

// HelloWorld is the function exposed by the plugin.
func HelloWorld(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}