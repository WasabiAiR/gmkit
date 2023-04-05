package errors

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/graymeta/gmkit/http/middleware"
)

// ClientErr is an error type that provides useful error messages that include
// both request and response bodies, status code of response and valid request
// parameters.
type ClientErr struct {
	op         string
	u          url.URL
	method     string
	err        error
	respBody   string
	reqBody    string
	reqID      string
	meta       []string
	StatusCode int
	retry      bool
	notFound   bool
	exists     bool
}

// NewClientErr is a constructor for a client error. The provided options
// allow the caller to set an optional retry.
func NewClientErr(op string, err error, resp *http.Response, opts ...ClientOptFn) error {
	newClientErr := &ClientErr{
		op:  op,
		err: errors.New("received unexpected response"),
	}
	for _, o := range opts {
		newClientErr = o(newClientErr)
	}

	if err != nil {
		newClientErr.err = err
	}

	if resp == nil {
		return newClientErr
	}

	if req := resp.Request; req != nil {
		newClientErr.u = *req.URL
		newClientErr.method = req.Method

		if req.Header != nil && strings.Contains(req.Header.Get("Content-Type"), "application/json") {
			if body, err := io.ReadAll(req.Body); err == nil {
				newClientErr.respBody = string(body)
			}
		}
	}
	newClientErr.StatusCode = resp.StatusCode
	newClientErr.reqID = resp.Header.Get(middleware.RequestHeader)

	if body, err := io.ReadAll(resp.Body); err == nil {
		newClientErr.respBody = string(body)
	}

	return newClientErr
}

// Error returns the full client error message.
func (e *ClientErr) Error() string {
	parts := []string{e.errorBase()}

	if respBody := e.respBody; respBody != "" {
		parts = append(parts, fmt.Sprintf("response_body=%q", respBody))
	}

	if reqBody := e.reqBody; reqBody != "" {
		parts = append(parts, fmt.Sprintf("request_body=%q", reqBody))
	}

	return strings.Join(parts, " ")
}

// Unwrap returns the wrapped error.
func (e *ClientErr) Unwrap() error {
	return e.err
}

// BackoffMessage provides a condensed error message that can be consumed during
// a backoff loop.
func (e *ClientErr) BackoffMessage() string {
	return e.errorBase()
}

// Retry provides the retry behavior.
func (e *ClientErr) Retry() bool {
	return e.retry
}

// NotFound provides the NotFounder behavior.
func (e *ClientErr) NotFound() bool {
	return e.notFound
}

// Exists provides the Exister behavior.
func (e *ClientErr) Exists() bool {
	return e.exists
}

func (e *ClientErr) errorBase() string {
	var parts []string

	if e.StatusCode > 0 {
		parts = append(parts, fmt.Sprintf("status=%d", e.StatusCode))
	}

	if e.method != "" {
		parts = append(parts, fmt.Sprintf("method=%s", e.method))
	}

	if e.reqID != "" {
		parts = append(parts, fmt.Sprintf("response_http_req_id=%q", e.reqID))
	}

	if len(e.meta) > 0 {
		parts = append(parts, getPairPrint(e.meta))
	}

	parts = append(parts, fmt.Sprintf("err=%q", e.err.Error()))

	if e.u.String() != "" {
		q := e.u.Query()
		if q.Get("access_token") != "" {
			q.Set("access_token", "REDACTED")
		}
		if q.Get("secret") != "" {
			q.Set("secret", "REDACTED")
		}
		e.u.RawQuery = q.Encode()
		parts = append(parts, fmt.Sprintf("url=%q", e.u.String()))
	}

	return strings.Join(parts, " ")
}

// ClientOptFn is a optional parameter that allows one to extend
// a client error.
type ClientOptFn func(o *ClientErr) *ClientErr

// Exists sets the client error to Exists, exists=true.
func Exists() ClientOptFn {
	return func(o *ClientErr) *ClientErr {
		o.exists = true
		return o
	}
}

// Meta sets key value pairs to add context to the error output.
func Meta(key, value string, pairs ...string) ClientOptFn {
	return func(e *ClientErr) *ClientErr {
		e.meta = append(e.meta, key, value)
		e.meta = append(e.meta, pairs...)
		return e
	}
}

// NotFound sets the client error to NotFound, notFound=true.
func NotFound() ClientOptFn {
	return func(o *ClientErr) *ClientErr {
		o.notFound = true
		return o
	}
}

// Retry sets the option and subsequent client error to retriable, retry=true.
func Retry() ClientOptFn {
	return func(o *ClientErr) *ClientErr {
		o.retry = true
		return o
	}
}

func getPairPrint(pairs []string) string {
	var paired []string
	for index := 0; index < len(pairs)/2; index++ {
		i := index * 2
		pair := pairs[i : i+2]
		if len(pair) != 2 {
			break
		}
		paired = append(paired, fmt.Sprintf("%s=%q", pair[0], pair[1]))
	}
	return strings.Join(paired, " ")
}
