package httpc

import (
	"github.com/graymeta/gmkit/backoff"
)

// ClientOptFn sets keys on a client type.
type ClientOptFn func(Client) Client

// WithAuth sets the authorization func on the client type,
// and will be used as the default authFn for all requests
// from this client unless overwritten atn the request lvl.
func WithAuth(authFn AuthFn) ClientOptFn {
	return func(c Client) Client {
		c.authFn = authFn
		return c
	}
}

// WithBackoff sets the backoff on the client.
func WithBackoff(b backoff.Backoffer) ClientOptFn {
	return func(c Client) Client {
		c.backoff = b
		return c
	}
}

// WithBaseURL sets the base url for all requests. Any path provided will be
// appended to this WithBaseURL.
func WithBaseURL(baseURL string) ClientOptFn {
	return func(c Client) Client {
		c.baseURL = baseURL
		return c
	}
}

// WithEncode sets the encode func for the client.
func WithEncode(fn EncodeFn) ClientOptFn {
	return func(c Client) Client {
		c.encodeFn = fn
		return c
	}
}

// WithRetryClientTimeouts sets the response retry mechanism. Useful if you want to
// retry on a client timeout or something of that nature.
func WithRetryClientTimeouts() ClientOptFn {
	return func(c Client) Client {
		c.respRetryFn = retryClientTimeout
		return c
	}
}

// WithRetryResponseErrors sets the response retry mechanism.
func WithRetryResponseErrors() ClientOptFn {
	return func(c Client) Client {
		c.respRetryFn = makeRetrierError
		return c
	}
}

// WithResetSeekerToZero sets the seek params to zero for all future requests
// Useful if Body param is a ReadSeeker and should be reset on retry
func WithResetSeekerToZero() ClientOptFn {
	return func(c Client) Client {
		c.seekParams = &seekParams{}
		return c
	}
}
