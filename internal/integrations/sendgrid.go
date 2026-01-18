
package integrations

import (
	"fmt"
)

// SendGrid is a mock SendGrid client.
type SendGrid struct {
	apiKey string
}

// NewSendGrid creates a new SendGrid client.
func NewSendGrid(apiKey string) *SendGrid {
	return &SendGrid{
		apiKey: apiKey,
	}
}

// SendEmail sends an email.
func (s *SendGrid) SendEmail(to, from, subject, body string) error {
	if s.apiKey == "" {
		return fmt.Errorf("SendGrid API key not set")
	}

	fmt.Printf("Sending email to %s from %s with subject %s\n", to, from, subject)
	fmt.Printf("Body: %s\n", body)

	// In a real application, you would use the SendGrid API here.
	return nil
}
