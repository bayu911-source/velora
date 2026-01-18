
package agents

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"velora/internal/services"
)

// MultimodalAgent is an agent that can process text and images.
type MultimodalAgent struct {
	llm services.LLM
}

// NewMultimodalAgent creates a new MultimodalAgent.
func NewMultimodalAgent(llm services.LLM) *MultimodalAgent {
	return &MultimodalAgent{llm: llm}
}

// Name returns the name of the agent.
func (a *MultimodalAgent) Name() string {
	return "multimodal"
}

// Description returns the description of the agent.
func (a *MultimodalAgent) Description() string {
	return "Analyzes an image from a URL or file path and provides information about it."
}

// Run executes the agent with the given input.
func (a *MultimodalAgent) Run(ctx context.Context, input string) (string, error) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid input: must be in the format '<image_path_or_url> <prompt>'")
	}

	imagePath := parts[0]
	prompt := parts[1]

	var imageData []byte
	var err error

	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		imageData, err = readImageFromURL(imagePath)
		if err != nil {
			return "", fmt.Errorf("failed to read image from URL: %w", err)
		}
	} else {
		imageData, err = os.ReadFile(imagePath)
		if err != nil {
			return "", fmt.Errorf("failed to read image file: %w", err)
		}
	}

	// TODO: Infer mime type from image data or file extension.
	resp, err := a.llm.GenerateContent(ctx, genai.Text(prompt), genai.ImageData("jpeg", imageData))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return resp, nil
}

func readImageFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
