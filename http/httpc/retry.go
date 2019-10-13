package httpc

import (
	"strings"

	gmerrors "github.com/graymeta/gmkit/errors"
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

// RetryResponseError sets a function to be called on all response errors.
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

// makeRetrierError transforms any error into a retryErr if it doesn't exhibit
// the gmerrors.Retrier behaviour already.
func makeRetrierError(err error) error {
	if _, ok := err.(gmerrors.Retrier); ok {
		return err
	}

	return &retryErr{err}
}

// RetryResponseErrors tells the request to retry execution errors.
func (r *Request) RetryResponseErrors() *Request {
	return RetryResponseError(makeRetrierError)(r)
}
