
package agents

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"velora/internal/services"
)

// WebResearchAgent is an agent that researches topics on the web.
type WebResearchAgent struct {
	gemini *services.GeminiService
}

// NewWebResearchAgent creates a new WebResearchAgent.
func NewWebResearchAgent(gemini *services.GeminiService) *WebResearchAgent {
	return &WebResearchAgent{
		gemini: gemini,
	}
}

// Name returns the name of the agent.
func (a *WebResearchAgent) Name() string {
	return "web_research"
}

// Run runs the agent. It takes a URL as input, scrapes the content,
// and uses Gemini to summarize it.
func (a *WebResearchAgent) Run(url string) (string, error) {
	// 1. Create a new HTTP request and set a realistic User-Agent
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	// 2. Fetch the URL
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("request failed with status: %d %s", res.StatusCode, res.Status)
	}

	// 3. Parse the HTML using goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// 4. Extract text from the body
	textContent := doc.Find("body").Text()
	textContent = strings.Join(strings.Fields(textContent), " ") // Clean up whitespace

	if len(textContent) == 0 {
		return "Could not extract any text content from the URL.", nil
	}

	// Truncate for very long content to avoid exceeding token limits
	maxChars := 15000
	if len(textContent) > maxChars {
		textContent = textContent[:maxChars]
	}

	// 5. Use Gemini to summarize the text
	prompt := fmt.Sprintf("Please summarize the following text extracted from the URL %s:\n\n%s", url, textContent)

	return a.gemini.Generate(prompt, "gemini-1.5-flash", 0.7, 2048)
}
