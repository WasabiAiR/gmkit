package testhelpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

// NewRandomPort identifies a port on the localhost interface for use during tests
// and returns the string in host:port format as well as a url with an http scheme.
// It uses similar methodology to how the net/http/httptest server chooses a port.
func NewRandomPort(t *testing.T) (string, string) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr()
	l.Close()
	return addr.String(), "http://" + addr.String()
}

// GET is a test utility to test a svr handles a GET call.
func GET(t *testing.T, svr http.Handler, addr string, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodGet, addr, nil, assertFns...)
}

// HEAD is a test utility to test a svr handles a HEAD call.
func HEAD(t *testing.T, svr http.Handler, addr string, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodHead, addr, nil, assertFns...)
}

// POST is a test utility to test a svr handles a POST call.
func POST(t *testing.T, svr http.Handler, addr string, body io.Reader, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodPost, addr, body, assertFns...)
}

// PATCH is a test utility to test a svr handles a PATCH call.
func PATCH(t *testing.T, svr http.Handler, addr string, body io.Reader, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodPatch, addr, body, assertFns...)
}

// PUT is a test utility to test a svr handles a PUT call.
func PUT(t *testing.T, svr http.Handler, addr string, body io.Reader, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodPut, addr, body, assertFns...)
}

// DELETE is a test utility to test a svr handles a DELETE call.
func DELETE(t *testing.T, svr http.Handler, addr string, assertFns ...RecTestFn) {
	t.Helper()
	HTTP(t, svr, http.MethodDelete, addr, nil, assertFns...)
}

// HTTP is a test utility to test a svr handles whatever call you provide it.
func HTTP(
	t *testing.T, svr http.Handler,
	method, addr string, body io.Reader,
	assertFns ...RecTestFn,
) {
	t.Helper()
	req, w := httptest.NewRequest(method, addr, body), httptest.NewRecorder()

	svr.ServeHTTP(w, req)

	for _, fn := range assertFns {
		fn(t, w)
	}
}

// RecTestFn is a functional option to run assertions against the response of the http request
// being made by the HTTP func or its relatives.
type RecTestFn func(t *testing.T, w *httptest.ResponseRecorder)

// Resp allows the body to be asserted and viewed per the function.
func Resp(assertFn func(t *testing.T, w *httptest.ResponseRecorder)) RecTestFn {
	return func(t *testing.T, w *httptest.ResponseRecorder) {
		t.Helper()
		assertFn(t, w)
	}
}

// Status verifies the status code of the response matches the status provided.
func Status(status int) RecTestFn {
	return func(t *testing.T, w *httptest.ResponseRecorder) {
		t.Helper()
		if w.Code == status {
			return
		}

		require.Failf(t, "received incorrect status code", "want: %d\tgot: %d", status, w.Code)
	}
}

// StatusOK verifies the status code is 200 (Status OK).
func StatusOK() RecTestFn {
	return Status(http.StatusOK)
}

// StatusCreated verifies the status code is 201 (Status Created).
func StatusCreated() RecTestFn {
	return Status(http.StatusCreated)
}

// StatusAccepted verifies the status code is 202 (Status Accepted).
func StatusAccepted() RecTestFn {
	return Status(http.StatusAccepted)
}

// StatusNoContent verifies the status code is 204 (Status No Content).
func StatusNoContent() RecTestFn {
	return Status(http.StatusNoContent)
}

// StatusBadRequest verifies the sttus code is 400 (Bad Request).
func StatusBadRequest() RecTestFn {
	return Status(http.StatusBadRequest)
}

// StatusNotFound verifies the status code is 404 (Status Not Found).
func StatusNotFound() RecTestFn {
	return Status(http.StatusNotFound)
}

// StatusUnprocessableEntity verifies the status code is 422 (Status Unprocessable Entity).
func StatusUnprocessableEntity() RecTestFn {
	return Status(http.StatusUnprocessableEntity)
}

// StatusInternalServerError verifies the status code is 500 (Status Internal Server Error).
func StatusInternalServerError() RecTestFn {
	return Status(http.StatusInternalServerError)
}

// StatusNotImplemented verifies the status code is 501 (Status Not Implemented).
func StatusNotImplemented() RecTestFn {
	return Status(http.StatusNotImplemented)
}

// Header verifies that the header of the response with key matches the val
// provided.
func Header(key, val string) RecTestFn {
	return func(t *testing.T, w *httptest.ResponseRecorder) {
		t.Helper()
		actual := w.Header().Get(key)
		if actual == val {
			return
		}

		require.Failf(t, "received incorrect value for header", "key: %q; want: %q, got: %q", key, val, actual)
	}
}

// ContentType verifies that the Content-Type of the response matches the
// contentType provided.
func ContentType(contentType string) RecTestFn {
	return Header("Content-Type", contentType)
}

// Do is a test helper for creating a req and sending that req to a test server.
func Do(t *testing.T, method, addr string, body io.Reader) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, addr, body)
	require.NoError(t, errors.Wrap(err, "new request"))

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Do(req)
	require.NoError(t, errors.Wrap(err, "do"))

	return resp
}

// GetURL is a helper for parsing a string address.
func GetURL(t *testing.T, addr string) *url.URL {
	t.Helper()

	u, err := url.Parse(addr)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

// EncodeBody is a helper func for encoding a type into an bytes.Buffer type.
// Can be used for a response body or whatever else. The returned buffer can
// be written or read from as needed.
func EncodeBody(t *testing.T, v interface{}) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		t.Fatal(err)
	}
	return &buf
}

// DecodeBody is a helper for decoding an io.reader (i.e. response body)
// into the destination type (v).
func DecodeBody(t *testing.T, r io.Reader, v interface{}) {
	t.Helper()

	if err := json.NewDecoder(r).Decode(v); err != nil {
		t.Fatal(err)
	}
}

// JSONPrettyPrint pretty prints a json body in all its glory.
func JSONPrettyPrint(t *testing.T, v interface{}) {
	t.Helper()

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + string(b))
}
