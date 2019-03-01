package icinga

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestClientBad(t *testing.T) {
	icingaCfg := Config{
		BaseURL:       "shitola",
		TLSClientCert: "Bogus",
		TLSClientKey:  "Bogus",
		TLSCACert:     "Bogus",
	}

	is := is.New(t)
	_, err := icingaCfg.Client()
	is.Err(err)
}

func TestClientBadBase(t *testing.T) {
	icingaCfg := Config{
		BaseURL:       "",
		TLSClientCert: "Bogus",
		TLSClientKey:  "Bogus",
		TLSCACert:     "Bogus",
	}

	is := is.New(t)
	_, err := icingaCfg.Client()
	is.Err(err)
}

func TestClientBadAuth(t *testing.T) {
	icingaCfg := Config{
		BaseURL:       "shitola",
		TLSClientCert: "",
		TLSClientKey:  "",
		TLSCACert:     "",
		Username:      "",
		Password:      "",
	}

	is := is.New(t)
	_, err := icingaCfg.Client()
	is.Err(err)
}
