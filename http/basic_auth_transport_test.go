package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicAuthTransport(t *testing.T) {
	client := &http.Client{
		Transport: NewBasicAuthTransport(http.DefaultTransport, "someuser", "somepassword"),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		require.True(t, ok)
		require.Equal(t, "someuser", user)
		require.Equal(t, "somepassword", pass)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, err := client.Get(ts.URL)
	require.NoError(t, err)
}
