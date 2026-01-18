
package services

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiService handles interactions with the Google Gemini API.
type GeminiService struct {
	generativeClient *genai.Client
}

// NewGeminiService creates a new GeminiService.
func NewGeminiService(apiKey string) (*GeminiService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey), option.WithEndpoint("generativelanguage.googleapis.com"))
	if err != nil {
		return nil, err
	}

	return &GeminiService{
		generativeClient: client,
	}, nil
}

// Generate generates content from a text prompt with configurable parameters.
func (s *GeminiService) Generate(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
	model := s.generativeClient.GenerativeModel(modelName)
	model.SetTemperature(temperature)
	model.SetMaxOutputTokens(int32(maxOutputTokens))

	ctx := context.Background()
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	return s.extractText(resp), nil
}


// ListModels lists the available models.
func (s *GeminiService) ListModels(ctx context.Context) *genai.ModelInfoIterator {
	return s.generativeClient.ListModels(ctx)
}

// Close closes the Gemini client.
func (s *GeminiService) Close() {
	s.generativeClient.Close()
}

// extractText extracts the text from a GenerateContentResponse.
func (s *GeminiService) extractText(resp *genai.GenerateContentResponse) string {
	var text string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					text += string(txt)
				}
			}
		}
	}
	return text
}
