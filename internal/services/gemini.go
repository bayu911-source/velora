
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"velora/config"
	"velora/internal/retry"
)

// GeminiService handles interactions with the Google Gemini API.
type GeminiService struct {
	generativeClient *genai.Client
}

// NewGeminiService creates a new GeminiService.
func NewGeminiService(cfg config.Config) (*GeminiService, error) {
	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiAPIKey), option.WithEndpoint(cfg.GeminiAPIURL))
	if err != nil {
		return nil, err
	}

	return &GeminiService{
		generativeClient: client,
	}, nil
}

// Generate generates content from a text prompt with configurable parameters.
func (s *GeminiService) Generate(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
	var resp *genai.GenerateContentResponse
	err := retry.Do(
		func() error {
			model := s.generativeClient.GenerativeModel(modelName)
			model.SetTemperature(temperature)
			model.SetMaxOutputTokens(int32(maxOutputTokens))

			ctx := context.Background()
			var err error
			resp, err = model.GenerateContent(ctx, genai.Text(prompt))
			return err
		},
		3,
		2*time.Second,
	)

	if err != nil {
		return "", err
	}

	return s.extractText(resp), nil
}

// GenerateStream generates content from a text prompt and streams the response.
func (s *GeminiService) GenerateStream(prompt, modelName string, temperature float32, maxOutputTokens int) (<-chan string, <-chan error) {
	model := s.generativeClient.GenerativeModel(modelName)
	model.SetTemperature(temperature)
	model.SetMaxOutputTokens(int32(maxOutputTokens))

	ctx := context.Background()
	iter := model.GenerateContentStream(ctx, genai.Text(prompt))

	out := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errChan)
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				errChan <- err
				return
			}
			out <- s.extractText(resp)
		}
	}()

	return out, errChan
}

// ListModels lists the available models.
func (s *GeminiService) ListModels(ctx context.Context) *genai.ModelInfoIterator {
	return s.generativeClient.ListModels(ctx)
}

// Close closes the Gemini client.
func (s *GeminiService) Close() error {
	return s.generativeClient.Close()
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
