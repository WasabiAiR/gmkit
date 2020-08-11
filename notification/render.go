package notification

import (
	"github.com/matcornic/hermes/v2"
	"github.com/pkg/errors"
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
		return "", "", errors.Wrap(err, "generating html email")
	}
	text, err := r.Template.GeneratePlainText(tmpl)
	if err != nil {
		return "", "", errors.Wrap(err, "generating plaintext email")
	}

	return text, html, nil
}
