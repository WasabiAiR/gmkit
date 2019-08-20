package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/graymeta/gmkit/notification"

	"github.com/pkg/errors"
)

// Environment variable keys controlling configuration
const (
	ConfigHost               = "smtp_host"
	ConfigPort               = "smtp_port"
	ConfigPassword           = "smtp_password"
	ConfigInsecureSkipVerify = "smtp_insecure_skip_verify"
)

// New constructs a new SMTP notification sender
func New(host, port, password string, tlsSkipVerify bool) notification.Sender {
	return &Client{
		host:     host,
		port:     port,
		password: password,
		insecure: tlsSkipVerify,
	}
}

var _ notification.Sender = (*Client)(nil)

// Client is a SMTP client and implements Sender
type Client struct {
	host     string
	port     string
	password string
	insecure bool
}

// Send sends the notification
func (c Client) Send(msg notification.Message) error {
	auth := smtp.PlainAuth("", msg.From, c.password, c.host)
	tlsConfig := tls.Config{ServerName: c.host, InsecureSkipVerify: c.insecure}

	conn, err := tls.Dial("tcp", c.host+":"+c.port, &tlsConfig)
	if err != nil {
		return errors.Wrap(err, "couldn't dial SMTP server")
	}

	smtpClient, err := smtp.NewClient(conn, c.host)
	if err != nil {
		return errors.Wrap(err, "cannot create new SMTP client")
	}

	err = smtpClient.Auth(auth)
	if err != nil {
		return errors.Wrap(err, "cannot authenticate against SMTP server")
	}

	err = smtpClient.Mail(msg.From)
	if err != nil {
		return errors.Wrap(err, "preparing new mail")
	}

	for _, r := range msg.To {
		if err = smtpClient.Rcpt(r); err != nil {
			return errors.Wrap(err, "setting receivers mails")
		}
	}

	w, err := smtpClient.Data()
	if err != nil {
		return err
	}

	defer func() {
		smtpClient.Close()
		w.Close()
	}()

	m := buildEmail(msg)

	_, err = w.Write(m)
	return errors.Wrap(err, "sending email")
}

func buildEmail(msg notification.Message) []byte {
	header := ""

	header += fmt.Sprintf("From: %s\r\n", msg.From)

	header += fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ";"))

	header += fmt.Sprintf("Subject: %s\r\n", msg.Subject)

	header += "\r\n" + msg.Body

	return []byte(header)
}
