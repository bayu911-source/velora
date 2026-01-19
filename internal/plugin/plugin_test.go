
package plugin_test

import (
	"os"
	"path/filepath"
	"testing"

	"velora/internal/plugin"
)

func TestManager(t *testing.T) {
	// Create a dummy plugin for testing.
	pluginDir := t.TempDir()
	pluginPath := filepath.Join(pluginDir, "test_plugin.so")

	// In a real test, you would compile a dummy plugin. For now, we'll
	// just create an empty file to simulate the plugin's existence.
	if _, err := os.Create(pluginPath); err != nil {
		t.Fatalf("failed to create dummy plugin file: %v", err)
	}

	m := plugin.NewManager()
	err := m.LoadAgentsFromDir(pluginDir)

	// Since we're not loading a real plugin, we expect an error here.
	// The purpose of this test is to ensure the directory walk and
	// plugin loading logic is exercised.
	if err == nil {
		t.Error("expected an error when loading a non-plugin file, but got nil")
	}
}
