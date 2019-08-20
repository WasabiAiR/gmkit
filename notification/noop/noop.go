package noop

import (
	"strings"

	"github.com/graymeta/gmkit/notification"

	"github.com/graymeta/gmkit/logger"
)

// Sender is a notificaiton.Sender that just logs the body of the messages
type Sender struct {
	logger *logger.L
}

var _ notification.Sender = (*Sender)(nil)

// New constructs a new Sender
func New(l *logger.L) *Sender {
	return &Sender{
		logger: l,
	}
}

// Send sends the specified message
func (s *Sender) Send(msg notification.Message) error {
	s.logger.Info(
		"notification_send_message",
		"from", msg.From,
		"to", strings.Join(msg.To, ","),
		"body", msg.Body,
		"body_html", msg.BodyHTML,
	)
	return nil
}
