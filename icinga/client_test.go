package icinga

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientBad(t *testing.T) {
	icingaCfg := Config{
		BaseURL:       "shitola",
		TLSClientCert: "Bogus",
		TLSClientKey:  "Bogus",
		TLSCACert:     "Bogus",
	}

	_, err := icingaCfg.Client()
	require.Error(t, err)
}

func TestClientBadBase(t *testing.T) {
	icingaCfg := Config{
		BaseURL:       "",
		TLSClientCert: "Bogus",
		TLSClientKey:  "Bogus",
		TLSCACert:     "Bogus",
	}

	_, err := icingaCfg.Client()
	require.Error(t, err)
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

	_, err := icingaCfg.Client()
	require.Error(t, err)
}
