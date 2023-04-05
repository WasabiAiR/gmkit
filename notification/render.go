package notification

import (
	"fmt"

	"github.com/matcornic/hermes"
)

// Renderer is an interface for rendering emails.
type Renderer interface {
	// Render returns the plain text email body, html email body, and an error.
	Render(tmpl hermes.Email) (string, string, error)
}

// HermesRenderer takes a hermes template and uses that to render a message.
type HermesRenderer struct {
	Template *hermes.Hermes
}

// Render returns the plain text email body, html email body, and an error.
func (r *HermesRenderer) Render(tmpl hermes.Email) (string, string, error) {
	html, err := r.Template.GenerateHTML(tmpl)
	if err != nil {
		return "", "", fmt.Errorf("generating html email: %w", err)
	}
	text, err := r.Template.GeneratePlainText(tmpl)
	if err != nil {
		return "", "", fmt.Errorf("generating plaintext email: %w", err)
	}

	return text, html, nil
}
