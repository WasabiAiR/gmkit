package icinga

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHostExistSuccess(t *testing.T) {
	jsonData := []byte(`
	{
    "results": [
        {
            "attrs": {
                "address": "10.21.8.13",
                "name": "client1.example.com"
            },
            "joins": {},
            "meta": {},
            "name": "cleint1.example.com",
            "type": "Host"
        }
    ]
	}`)

	hler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	s := httptest.NewServer(hler)
	defer s.Close()

	icingaCfg := Config{
		BaseURL:  s.URL,
		Username: "test",
		Password: "test",
	}

	ic, err := icingaCfg.Client()
	require.NoError(t, err)

	result, err := ic.HostExist("client1.example.com")
	require.NoError(t, err)
	require.True(t, result)
}

func TestHostExistFailure(t *testing.T) {
	jsonData := []byte(`
	{
		"results": []
	}`)

	hler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	s := httptest.NewServer(hler)
	defer s.Close()

	icingaCfg := Config{
		BaseURL:  s.URL,
		Username: "test",
		Password: "test",
	}

	ic, err := icingaCfg.Client()
	require.NoError(t, err)

	result, err := ic.HostExist("client1.example.com")
	require.NoError(t, err)
	require.False(t, result)
}
