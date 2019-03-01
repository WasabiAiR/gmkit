package icinga

import (
	"flag"
	"log"
	"os"

	"github.com/graymeta/env"
)

// Config setups the information needed to connect to icinga
// If the system has icinga installed then you should use the cets used in /etc/icinga2/pki
// Username/Password should only used if it does not have a icinga2 agent installed.
type Config struct {
	flag *flag.FlagSet

	BaseURL       string
	Username      string
	Password      string
	TLSClientCert string
	TLSClientKey  string
	TLSCACert     string
	TLSInsecure   bool
}

// NewConfig prepares a new configuration for the parsing icinga flags
func NewConfig(flag *flag.FlagSet) *Config {
	if flag != nil && flag.Parsed() {
		log.Fatal("flag must be unparsed when calling NewConfig")
	}

	cfg := Config{
		flag: flag,
	}

	flag.StringVar(&cfg.BaseURL, "icinga", os.Getenv("gm_icinga"), "icinga address")
	flag.StringVar(&cfg.Username, "icinga-username", os.Getenv("gm_icinga_username"), "icinga username (optional)")
	flag.StringVar(&cfg.Password, "icinga-password", os.Getenv("gm_icinga_password"), "icinga password (optional)")
	flag.StringVar(&cfg.TLSClientCert, "icinga-tls-client-cert", os.Getenv("gm_icinga_tls_client_cert"), "The TLS client certificate to use when connecting to Icinga")
	flag.StringVar(&cfg.TLSClientKey, "icinga-tls-client-key", os.Getenv("gm_icinga_tls_client_key"), "The TLS client key to use when connecting to Icinga")
	flag.StringVar(&cfg.TLSCACert, "icinga-tls-ca-cert", os.Getenv("gm_icinga_tls_ca_cert"), "The CA cert to use to validate the Icinga server's certificate")
	flag.BoolVar(&cfg.TLSInsecure, "icinga-tls-insecure", env.GetenvBoolWithDefault("gm_icinga_tls_insecure", false), "Whether or not to perform certificate hostname validation when connecting to the Nomad server")

	return &cfg
}
