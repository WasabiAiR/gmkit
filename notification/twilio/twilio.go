package twilio

import (
	"fmt"

	"github.com/graymeta/gmkit/notification"

	twiliolib "github.com/carlosdp/twiliogo"
)

// Envioronment variable keys for configuration.
const (
	ConfigAccountSID = "twilio_account_sid"
	ConfigAuthToken  = "twilio_auth_token"
)

// New constructs a new notification sender that sends messages via Twilio.
func New(accountSID, authToken string) *Sender {
	return &Sender{
		client: twiliolib.NewClient(accountSID, authToken),
	}
}

// Sender is a Twilio based sender.
type Sender struct {
	client twiliolib.Client
}

var _ notification.Sender = (*Sender)(nil)

// Send sends the message.
func (t *Sender) Send(msg notification.Message) error {
	for _, r := range msg.To {
		_, err := twiliolib.NewMessage(t.client, msg.From, r, twiliolib.Body(msg.Body))
		if err != nil {
			return fmt.Errorf("sending Twilio SMS message: %w", err)
		}
	}

	return nil
}
