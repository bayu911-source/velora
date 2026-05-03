
package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"velora/internal/services"
)

// AppBuilderAgent builds simple applications based on a description.
type AppBuilderAgent struct {
	LLM *services.LLM
}

// NewAppBuilderAgent creates a new AppBuilderAgent.
func NewAppBuilderAgent(llm *services.LLM) *AppBuilderAgent {
	return &AppBuilderAgent{
		LLM: llm,
	}
}

// Name returns the name of the agent.
func (a *AppBuilderAgent) Name() string {
	return "app_builder"
}

// Description returns a brief description of the agent.
func (a *AppBuilderAgent) Description() string {
	return "Builds a simple application based on a description, including technology stack, file structure, and code."
}

// Execute executes the agent's primary function: building an application.
func (a *AppBuilderAgent) Execute(ctx context.Context, input string) (string, error) {
	if a.LLM == nil {
		return "", fmt.Errorf("LLM service is not initialized")
	}

	prompt := fmt.Sprintf(`You are a 10x software engineer. Create the files for the following application. Each file should be separated by '---' and the file path should be the first line. Do not include any explanations, just the files themselves.

Application Description: %s`, input)

	// Generate the plan using the LLM service.
	resp, err := a.LLM.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to build application: %w", err)
	}

	// Parse the response and create the files.
	files := strings.Split(resp, "---")
	for _, file := range files {
		if strings.TrimSpace(file) == "" {
			continue
		}

		parts := strings.SplitN(file, "\n", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid file format: %s", file)
		}

		path := strings.TrimSpace(parts[0])
		content := parts[1]

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return "", fmt.Errorf("failed to create directory for %s: %w", path, err)
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return "", fmt.Errorf("failed to write file %s: %w", path, err)
		}
	}

	return "Application built successfully!", nil
}
