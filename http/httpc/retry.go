package httpc

import (
	"strings"
)

// RetryFn will apply a retry policy to a request.
type RetryFn func(*Request) *Request

// RetryStatus appends a retry func to the Request.
func RetryStatus(fn StatusFn) RetryFn {
	return func(req *Request) *Request {
		req.retryStatusFns = append(req.retryStatusFns, fn)
		return req
	}
}

// RetryResponseError applies a retry on all response errors. The errors
// typically associated with request timeouts or oauth token error.
// This option useful when the oauth auth made me invalid or a request timeout
// is an issue.
func RetryResponseError(fn ResponseErrorFn) RetryFn {
	return func(r *Request) *Request {
		r.responseErrFn = fn
		return r
	}
}

// RetryClientTimeout will retry the request if the request had been canceled
// by the http client.
func RetryClientTimeout() RetryFn {
	return func(r *Request) *Request {
		r.responseErrFn = retryClientTimeout
		return r
	}
}

func retryClientTimeout(err error) error {
	if err == nil {
		return nil
	}
	// can be found in net/http/transport.go within the go src code.
	// there is no exported type for this unfortunately.
	reqCanceledMsg := "net/http: request canceled while waiting for connection"
	if strings.Contains(err.Error(), reqCanceledMsg) {
		return &retryErr{err}
	}
	return err
}
