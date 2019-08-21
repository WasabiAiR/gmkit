package metrics

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeMetrics(t *testing.T) {
	var tests = []struct {
		method   string
		route    string
		expected string
	}{
		{http.MethodPut, "/foo", "PUT_foo"},
		{http.MethodPut, "foo", "PUT_foo"},
		{"put", "foo", "PUT_foo"},
		{http.MethodGet, "/foo/{id}/{bar}", "GET_foo_id_bar"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s:%s", tt.method, tt.route), func(t *testing.T) {
			require.Equal(t, tt.expected, SanitizeRoute(tt.method, tt.route))
		})
	}
}
