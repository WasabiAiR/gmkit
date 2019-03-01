package icinga

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Client holds the icinga2 http client
type Client struct {
	cfg        *Config
	httpClient *http.Client
	tlsConfig  *tls.Config
}

// Client checks and applies the config and return a icinga client
func (cfg *Config) Client() (*Client, error) {
	if cfg.flag != nil && !cfg.flag.Parsed() {
		return nil, errors.New("must parse flags before calling Client")
	}

	// Check to see if the configuration if valid.  We must have a address
	if cfg.BaseURL == "" {
		return nil, errors.New("icinga address is missing")
	}

	// Check to see if the configuration if valid.  We must have a tls or usernames/password
	if (cfg.TLSClientCert == "" || cfg.TLSClientKey == "" || cfg.TLSCACert == "") && (cfg.Username == "" || cfg.Password == "") {
		return nil, errors.New("icinga TLS or username/password not set")
	}

	tlsConfig, err := cfg.setupTLSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "setupTLSConfig error")
	}

	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
			Timeout: time.Second * 60,
		},
	}, nil
}

func (cfg *Config) setupTLSConfig() (*tls.Config, error) {
	if cfg.TLSClientCert == "" || cfg.TLSClientKey == "" || cfg.TLSCACert == "" {
		return &tls.Config{
			InsecureSkipVerify: cfg.TLSInsecure,
		}, nil
	}
	// Load client cert
	cert, err := tls.LoadX509KeyPair(cfg.TLSClientCert, cfg.TLSClientKey)
	if err != nil {
		return nil, errors.Wrap(err, "load tls cert and key")
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(cfg.TLSCACert)
	if err != nil {
		return nil, errors.Wrap(err, "Read ca cert")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: cfg.TLSInsecure,
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
