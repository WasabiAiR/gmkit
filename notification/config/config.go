package config

import (
	"flag"
	"os"

	"github.com/graymeta/gmkit/logger"
	"github.com/graymeta/gmkit/notification"
	ses "github.com/graymeta/gmkit/notification/amazonses"
	"github.com/graymeta/gmkit/notification/noop"
	"github.com/graymeta/gmkit/notification/sendgrid"

	"github.com/pkg/errors"
)

// Config used to setup the sender
type Config struct {
	flag *flag.FlagSet

	sendgridAPIKey string
	sesAccessKeyID string
	sesSecretKey   string
	sesRegion      string
}

// New initializes the configuration, registering command line flags.
func New(l *logger.L, flag *flag.FlagSet) *Config {
	if flag != nil && flag.Parsed() {
		l.Fatal("error", "flag must be unparsed when calling errs.NewConfig")
	}
	cfg := Config{
		flag: flag,
	}

	flag.StringVar(&cfg.sendgridAPIKey, "sendgrid-api-key", os.Getenv(sendgrid.EnvAPIKey), "The SendGrid API Key.")
	flag.StringVar(&cfg.sesAccessKeyID, "ses-id", os.Getenv(ses.EnvAccessKeyID), "The AWS SES access key ID.")
	flag.StringVar(&cfg.sesSecretKey, "ses-key", os.Getenv(ses.EnvSecretKey), "The AWS SES secret key.")
	flag.StringVar(&cfg.sesRegion, "ses-region", os.Getenv(ses.EnvRegion), "The AWS SES region.")

	return &cfg
}

// GetSender will look at the configuration and return a Sender
func (cfg *Config) GetSender(l *logger.L) (notification.Sender, error) {
	if cfg.flag != nil && !cfg.flag.Parsed() {
		return nil, errors.New("must parse flags before calling Setup")
	}

	switch {
	case cfg.sendgridAPIKey != "":
		return sendgrid.New(cfg.sendgridAPIKey), nil
	case cfg.sesRegion != "":
		return ses.New(cfg.sesAccessKeyID, cfg.sesSecretKey, cfg.sesRegion)
	default:
		return noop.New(l), nil
	}
}
