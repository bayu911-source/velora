
package main

import (
	"fmt"

	"github.com/velora-chat/velora/pkg/config"
	"github.com/velora-chat/velora/pkg/memory"
)

func main() {
	// Create a new memory manager
	memoryManager := memory.NewMemoryManager()

	// Create a new config manager
	configManager := config.NewConfigManager(memoryManager)

	// Set the config values
	configManager.Set("name", "Velora")
	configManager.Set("version", "1.0.0")

	// Get the config values
	name, _ := configManager.Get("name")
	version, _ := configManager.Get("version")

	// Print the config values
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Version: %s\n", version)
}
