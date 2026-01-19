
package main

import (
	"fmt"
	"os"
	"velora/cmd"
	"velora/internal/plugin"
)

func main() {
	// Create the plugin directory if it doesn't exist.
	if err := os.MkdirAll("plugins", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create plugins directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize the plugin manager and load plugins.
	pluginManager := plugin.NewManager()
	if err := pluginManager.LoadAgentsFromDir("plugins"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load plugins: %v\n", err)
	}

	cmd.Execute(pluginManager)
}
