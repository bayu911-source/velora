
package integrations

import "fmt"

// CRM is a mock CRM client.
type CRM struct {
	apiKey string
}

// NewCRM creates a new CRM client.
func NewCRM(apiKey string) *CRM {
	return &CRM{
		apiKey: apiKey,
	}
}

// CreateLead creates a new lead.
func (c *CRM) CreateLead(name, email string) error {
	if c.apiKey == "" {
		return fmt.Errorf("CRM API key not set")
	}

	fmt.Printf("Creating lead: %s, %s\n", name, email)
	return nil
}

// UpdateLead updates a lead.
func (c *CRM) UpdateLead(id, status string) error {
	if c.apiKey == "" {
		return fmt.Errorf("CRM API key not set")
	}

	fmt.Printf("Updating lead %s to %s\n", id, status)
	return nil
}
