package http

import (
	"net/http"
)

// BasicAuthTransport is an http.RoundTripper that adds basic auth credentials to each request.
type BasicAuthTransport struct {
	username  string
	password  string
	transport http.RoundTripper
}

// NewBasicAuthTransport sets up a transport that wraps the provided transport.
func NewBasicAuthTransport(t http.RoundTripper, username, password string) *BasicAuthTransport {
	return &BasicAuthTransport{
		username:  username,
		password:  password,
		transport: t,
	}
}

// RoundTrip adds the basic auth creds to the request and passes it along to the wrapped transport.
func (t *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.username, t.password)
	return t.transport.RoundTrip(req)
}
