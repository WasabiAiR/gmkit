package amazonses

import (
	"fmt"

	"github.com/graymeta/gmkit/notification"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

// The environment variable names for the various SES configuration keys.
const (
	EnvAccessKeyID = "ses_id"
	EnvSecretKey   = "ses_secret"
	EnvRegion      = "ses_region"
)

// New creates a new Amazon SES notification sender pulling the configuration from
// environment variables.
func New(accessKeyID, secretKey, region string) (notification.Sender, error) {
	config := aws.NewConfig().WithRegion(region)
	if accessKeyID != "" && secretKey != "" {
		config = config.WithCredentials(credentials.NewStaticCredentials(accessKeyID, secretKey, ""))
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("creating SES session: %w", err)
	}

	service := ses.New(sess)
	return &Sender{sess: service}, nil
}

// Sender is an SES based sender
type Sender struct {
	sess sesiface.SESAPI
}

var _ notification.Sender = (*Sender)(nil)

// GetSES returns the SES service so it can be used elsewhere.
func (c *Sender) GetSES() sesiface.SESAPI {
	return c.sess
}

// Send atempts to send the message.
func (c *Sender) Send(msg notification.Message) error {
	var toAddr []*string
	for _, to := range msg.To {
		toAddr = append(toAddr, &to)
	}
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddr,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(msg.BodyHTML),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(msg.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(msg.Subject),
			},
		},
		Source: aws.String(msg.From),
	}

	_, err := c.sess.SendEmail(input)
	if err != nil {
		return fmt.Errorf("sending email with Amazon SES: %w", err)
	}
	return nil
}
