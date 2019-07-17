package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFrameOptions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	})
	srv := httptest.NewServer(FrameOptions("foobar")(handler))

	req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, "foobar", resp.Header.Get("X-Frame-Options"))
}
