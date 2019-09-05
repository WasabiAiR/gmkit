package middleware

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeQuery(t *testing.T) {
	v := make(url.Values)

	v.Set("foo", "bar")
	v.Set("access_token", "some token")

	v = SanitizeQuery(v)

	require.Equal(t, "bar", v.Get("foo"))
	require.Equal(t, "REDACTED", v.Get("access_token"))
}
