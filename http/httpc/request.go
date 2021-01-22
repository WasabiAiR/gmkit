package httpc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/graymeta/gmkit/backoff"
	gmerrors "github.com/graymeta/gmkit/errors"
)

// ErrInvalidEncodeFn is an error that is returned when calling the Request Do and the
// encode function is not set.
var ErrInvalidEncodeFn = errors.New("no encode fn provided for body")

// ResponseErrorFn is a response error function that can be used to provide
// behavior when a response fails to "Do".
type ResponseErrorFn func(error) error

// ResponseHeadersFn is a function that can be used to retrieve headers from a response
type ResponseHeadersFn func(header http.Header)

type kvPair struct {
	key   string
	value string
}

type seekParams struct {
	offset int64
	whence int
}

// Request is built up to create an http request.
type Request struct {
	method, addr  string
	doer          Doer
	body          interface{}
	headers       []kvPair
	params        []kvPair
	logMeta       []kvPair
	seekParams    *seekParams
	contentLength int

	authFn            AuthFn
	encodeFn          EncodeFn
	decodeFn          DecodeFn
	onErrorFn         DecodeFn
	responseErrFn     ResponseErrorFn
	responseHeadersFn ResponseHeadersFn

	existsFns      []StatusFn
	notFoundFns    []StatusFn
	retryStatusFns []StatusFn
	successFns     []StatusFn
	backoff        backoff.Backoffer
}

// Auth sets the authorization for hte request, overriding the authFn set
// by the client.
func (r *Request) Auth(authFn AuthFn) *Request {
	r.authFn = authFn
	return r
}

// Backoff sets the backoff of the Request.
func (r *Request) Backoff(b backoff.Backoffer) *Request {
	r.backoff = b
	return r
}

// Body sets the body of the Request.
func (r *Request) Body(v interface{}) *Request {
	r.body = v
	return r
}

// ContentLength sets the content length on the request.
func (r *Request) ContentLength(t int) *Request {
	r.contentLength = t
	return r
}

// ContentType sets the content type on the request.
func (r *Request) ContentType(t string) *Request {
	return r.Header("Content-Type", t)
}

// Decode sets the decoder func for the Request.
func (r *Request) Decode(fn DecodeFn) *Request {
	r.decodeFn = fn
	return r
}

// Encode sets the encoder func for the Request.
func (r *Request) Encode(fn EncodeFn) *Request {
	r.encodeFn = fn
	return r
}

// SeekParams sets the seek params.
// this is useful in cases where the body is a seeker and needs
// to be reset on retry
func (r *Request) SeekParams(offset int64, whence int) *Request {
	r.seekParams = &seekParams{
		offset: offset,
		whence: whence,
	}
	return r
}

// DecodeJSON is is a short hand for decoding to JSON.
func (r *Request) DecodeJSON(v interface{}) *Request {
	return r.Decode(JSONDecode(v))
}

// Exists appends a exists func to the Request.
func (r *Request) Exists(fn StatusFn) *Request {
	r.existsFns = append(r.existsFns, fn)
	return r
}

// Header adds a header to the request.
func (r *Request) Header(key, value string) *Request {
	r.headers = append(r.headers, kvPair{key: key, value: value})
	return r
}

// Meta adds k/v pairs to the eror message be added in the event of an error.
func (r *Request) Meta(key, value string, pairs ...string) *Request {
	r.logMeta = append(r.logMeta, kvPair{key: key, value: value})
	r.logMeta = append(r.logMeta, toKVPairs(pairs)...)
	return r
}

// NotFound appends a not found func to the Request.
func (r *Request) NotFound(fn StatusFn) *Request {
	r.notFoundFns = append(r.notFoundFns, fn)
	return r
}

// OnError provides a decode hook to decode a responses body.
func (r *Request) OnError(fn DecodeFn) *Request {
	r.onErrorFn = fn
	return r
}

// ResponseHeaders provides a hook to get the headers of a response
func (r *Request) ResponseHeaders(fn ResponseHeadersFn) *Request {
	r.responseHeadersFn = fn
	return r
}

// QueryParam allows a user to set query params on their request. This can be
// called numerous times. Will add keys for each value that is passed in here.
// In the case of duplicate query param values, the last pair that is entered
// will be set and the former will not be available.
func (r *Request) QueryParam(key, value string) *Request {
	r.params = append(r.params, kvPair{key: key, value: value})
	return r
}

// QueryParams allows a user to set multiple query params at one time. Following
// the same rules as QueryParam. If a key is provided without a value, it will
// not be added to the request. If it is desired, pass in "" to add a query param
// with no string field.
func (r *Request) QueryParams(key, value string, pairs ...string) *Request {
	paramed := r.QueryParam(key, value)
	paramed.params = append(paramed.params, toKVPairs(pairs)...)
	return paramed
}

// Retry sets the retry policy(s) on the request.
func (r *Request) Retry(fn RetryFn) *Request {
	return fn(r)
}

// RetryStatus sets the retry policy(s) on the request.
func (r *Request) RetryStatus(fn StatusFn) *Request {
	return RetryStatus(fn)(r)
}

// RetryStatusNotIn sets the retry policy(s) to retry if the response status
// matches that of the given args.
func (r *Request) RetryStatusNotIn(status int, rest ...int) *Request {
	return RetryStatus(StatusNotIn(status, rest...))(r)
}

// Success appends a success func to the Request.
func (r *Request) Success(fn StatusFn) *Request {
	r.successFns = append(r.successFns, fn)
	return r
}

// DoAndGetReader makes the http request and does not close the body in the http.Response that is returned
func (r *Request) DoAndGetReader(ctx context.Context) (*http.Response, error) {
	var resp *http.Response
	err := r.backoff.BackoffCtx(ctx, func(ctx context.Context) error {
		body, err := r.getReqBody()
		if err != nil {
			return gmerrors.NewClientErr("encode body", err, nil, r.metaErrOpts()...)
		}

		req, err := http.NewRequest(r.method, r.addr, body)
		if err != nil {
			return gmerrors.NewClientErr("new req", err, nil, r.metaErrOpts()...)
		}
		req = req.WithContext(ctx)

		if len(r.headers) > 0 {
			for _, pair := range r.headers {
				req.Header.Set(pair.key, pair.value)
			}
		}

		if len(r.params) > 0 {
			params := req.URL.Query()
			for _, kv := range r.params {
				params.Set(kv.key, kv.value)
			}
			req.URL.RawQuery = params.Encode()
		}

		if r.authFn != nil {
			req = r.authFn(req)
		}

		if r.contentLength > 0 {
			req.ContentLength = int64(r.contentLength)
		}

		resp, err = r.doer.Do(req)
		if err != nil {
			return r.responseErr(resp, err)
		}
		if r.responseHeadersFn != nil {
			r.responseHeadersFn(resp.Header)
		}

		status := resp.StatusCode
		if !statusMatches(status, r.successFns) {
			defer func() {
				drain(resp.Body)
			}()
			var err error
			if r.onErrorFn != nil {
				var buf bytes.Buffer
				tee := io.TeeReader(resp.Body, &buf)
				err = r.onErrorFn(tee)
				resp.Body = ioutil.NopCloser(&buf)
			}
			return gmerrors.NewClientErr("status code", err, resp, r.statusErrOpts(status)...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Do makes the http request and applies the backoff.
func (r *Request) Do(ctx context.Context) error {
	return r.backoff.BackoffCtx(ctx, r.do)
}

func (r *Request) do(ctx context.Context) error {
	body, err := r.getReqBody()
	if err != nil {
		return gmerrors.NewClientErr("encode body", err, nil, r.metaErrOpts()...)
	}

	req, err := http.NewRequest(r.method, r.addr, body)
	if err != nil {
		return gmerrors.NewClientErr("new req", err, nil, r.metaErrOpts()...)
	}
	req = req.WithContext(ctx)

	if len(r.headers) > 0 {
		for _, pair := range r.headers {
			req.Header.Set(pair.key, pair.value)
		}
	}

	if len(r.params) > 0 {
		params := req.URL.Query()
		for _, kv := range r.params {
			params.Set(kv.key, kv.value)
		}
		req.URL.RawQuery = params.Encode()
	}

	if r.authFn != nil {
		req = r.authFn(req)
	}

	if r.contentLength > 0 {
		req.ContentLength = int64(r.contentLength)
	}

	resp, err := r.doer.Do(req)
	if err != nil {
		return r.responseErr(resp, err)
	}
	if r.responseHeadersFn != nil {
		r.responseHeadersFn(resp.Header)
	}
	defer func() {
		drain(resp.Body)
	}()

	status := resp.StatusCode
	if !statusMatches(status, r.successFns) {
		var err error
		if r.onErrorFn != nil {
			var buf bytes.Buffer
			tee := io.TeeReader(resp.Body, &buf)
			err = r.onErrorFn(tee)
			resp.Body = ioutil.NopCloser(&buf)
		}
		return gmerrors.NewClientErr("status code", err, resp, r.statusErrOpts(status)...)
	}

	if r.decodeFn == nil {
		return nil
	}

	if err := r.decodeFn(resp.Body); err != nil {
		var opts []gmerrors.ClientOptFn
		if isRetryErr(err) {
			opts = append(opts, gmerrors.Retry())
		}
		opts = append(opts, r.metaErrOpts()...)
		return gmerrors.NewClientErr("decode", err, resp, opts...)
	}
	return nil
}

func (r *Request) getReqBody() (io.Reader, error) {
	if r.body == nil {
		return nil, nil
	}

	if reader, ok := r.body.(io.Reader); ok {
		if seeker, ok2 := r.body.(io.Seeker); ok2 && r.seekParams != nil {
			seeker.Seek(r.seekParams.offset, r.seekParams.whence)
		}
		return reader, nil
	}

	if r.encodeFn == nil {
		return nil, ErrInvalidEncodeFn
	}

	encodedBody, err := r.encodeFn(r.body)
	if err != nil || encodedBody == nil {
		return nil, err
	}
	return encodedBody, nil
}

func (r *Request) metaErrOpts() []gmerrors.ClientOptFn {
	if len(r.logMeta) == 0 {
		return nil
	}
	var pairs []string
	for _, p := range r.logMeta {
		pairs = append(pairs, p.key, p.value)
	}
	return []gmerrors.ClientOptFn{gmerrors.Meta(pairs[0], pairs[1], pairs[2:]...)}
}

func (r *Request) statusErrOpts(status int) []gmerrors.ClientOptFn {
	var opts []gmerrors.ClientOptFn
	if statusMatches(status, r.retryStatusFns) {
		opts = append(opts, gmerrors.Retry())
	}
	if statusMatches(status, r.notFoundFns) {
		opts = append(opts, gmerrors.NotFound())
	}
	if statusMatches(status, r.existsFns) {
		opts = append(opts, gmerrors.Exists())
	}
	opts = append(opts, r.metaErrOpts()...)
	return opts
}

func (r *Request) responseErr(resp *http.Response, err error) error {
	if r.responseErrFn != nil {
		err = r.responseErrFn(err)
	}
	var opts []gmerrors.ClientOptFn
	if isRetryErr(err) {
		opts = append(opts, gmerrors.Retry())
	}
	opts = append(opts, r.metaErrOpts()...)
	return gmerrors.NewClientErr("do", err, resp, opts...)
}

func toKVPairs(pairs []string) []kvPair {
	var kvPairs []kvPair
	for index := 0; index < len(pairs)/2; index++ {
		i := index * 2
		pair := pairs[i : i+2]
		if len(pair) != 2 {
			return kvPairs
		}
		kvPairs = append(kvPairs, kvPair{key: pair[0], value: pair[1]})
	}
	return kvPairs
}

// drain reads everything from the ReadCloser and closes it
func drain(r io.ReadCloser) error {
	if _, err := io.Copy(ioutil.Discard, r); err != nil {
		return err
	}
	return r.Close()
}

func isRetryErr(err error) bool {
	if err == nil {
		return false
	}
	r, ok := err.(gmerrors.Retrier)
	return ok && r.Retry()
}
