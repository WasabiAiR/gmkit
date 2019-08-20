package sendgrid

import (
	"github.com/graymeta/gmkit/notification"

	"github.com/pkg/errors"
	sendgridlib "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EnvAPIKey is the environment variable that the API key is read from.
const EnvAPIKey = "sendgrid_api_key"

// Sender is a notificaiton.Sender that uses SendGrid to ship messages
type Sender struct {
	client *sendgridlib.Client
}

var _ notification.Sender = (*Sender)(nil)

// New accepts a SendGrid API key and initializes a Sender
func New(apiKey string) *Sender {
	return &Sender{
		client: sendgridlib.NewSendClient(apiKey),
	}
}

// Send sends the specified message
func (s *Sender) Send(msg notification.Message) error {
	from := mail.NewEmail(msg.FromPretty, msg.From)
	for _, r := range msg.To {
		to := mail.NewEmail(r, r)
		m := mail.NewSingleEmail(from, msg.Subject, to, msg.Body, msg.BodyHTML)
		_, err := s.client.Send(m)
		if err != nil {
			return errors.Wrap(err, "sending email with SendGrid")
		}
	}
	return nil
}
