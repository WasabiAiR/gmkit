package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/graymeta/gmkit/notification"
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
		return fmt.Errorf("couldn't dial SMTP server: %w", err)
	}

	smtpClient, err := smtp.NewClient(conn, c.host)
	if err != nil {
		return fmt.Errorf("cannot create new SMTP client: %w", err)
	}

	err = smtpClient.Auth(auth)
	if err != nil {
		return fmt.Errorf("cannot authenticate against SMTP server: %w", err)
	}

	err = smtpClient.Mail(msg.From)
	if err != nil {
		return fmt.Errorf("preparing new mail: %w", err)
	}

	for _, r := range msg.To {
		if err = smtpClient.Rcpt(r); err != nil {
			return fmt.Errorf("setting receivers mails: %w", err)
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
	return fmt.Errorf("sending email: %w", err)
}

func buildEmail(msg notification.Message) []byte {
	header := ""

	header += fmt.Sprintf("From: %s\r\n", msg.From)

	header += fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ";"))

	header += fmt.Sprintf("Subject: %s\r\n", msg.Subject)

	header += "\r\n" + msg.Body

	return []byte(header)
}
