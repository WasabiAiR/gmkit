package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	gmerrors "github.com/graymeta/gmkit/errors"
	"github.com/graymeta/gmkit/http/middleware"
	"github.com/graymeta/gmkit/logger"

	"github.com/pkg/errors"
)

// Responder writes API responses.
// nil Responder is safe to use.
type Responder struct {
	log     *logger.L
	version string
}

// NewResponder makes a new Responder object.
func NewResponder(logger *logger.L, version string) *Responder {
	return &Responder{log: logger, version: version}
}

// Option functions provide a means for manipulating things like adding additional
// headers to responses.
type Option func(w http.ResponseWriter)

// Deprecated adds a header to indicate the particular endpoint is deprecated.
func Deprecated(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Warning", `299 - "Deprecated"`)
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// Header adds a Header to the response.
func Header(key, value string) Option {
	return func(w http.ResponseWriter) {
		w.Header().Set(key, value)
	}
}

// HeaderAPIVersion is the canonical name of the HTTP header key that reports the API version.
const HeaderAPIVersion = `X-Api-Version`

// With responds with the specified data.
func (r *Responder) With(w http.ResponseWriter, req *http.Request, status int, data interface{}, opts ...Option) {
	var buf bytes.Buffer
	// cannot write to buf if data is nil, in case of StatusNoContent, this write will fail
	// so we need an escape hatch here.
	if data != nil {
		enc := json.NewEncoder(&buf)
		enc.SetIndent("", "\t")
		err := enc.Encode(data)
		if err != nil {
			err = errors.Wrap(err, "failed to encode response object")
			r.Err(w, req, err)
			return
		}
	}
	r.log.Debug("api_response",
		"status", status,
		"body", buf.String(),
	)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set(HeaderAPIVersion, r.version)
	for _, fn := range opts {
		fn(w)
	}
	w.WriteHeader(status)
	if _, err := io.Copy(w, &buf); err != nil {
		err = errors.Wrap(err, "failed to copy response bytes")
		r.log.Err("api_response", err)
	}
}

// ErrorResponse is the body of an Error served up by Err
type ErrorResponse struct {
	Error struct {
		Message   string `json:"message"`
		RequestID string `json:"request_id,omitempty"`
	} `json:"error"`
}

func newErr(requestID, msg string) ErrorResponse {
	return ErrorResponse{
		Error: struct {
			Message   string `json:"message"`
			RequestID string `json:"request_id,omitempty"`
		}{
			Message:   msg,
			RequestID: requestID,
		},
	}
}

// Err responds with an error that corresponds to the behavior the type is illustrating.
func (r *Responder) Err(w http.ResponseWriter, req *http.Request, err error) {
	reqID := middleware.GetReqID(req.Context())
	logMsg := true
	defer func() {
		sanitizedQuery := middleware.SanitizeQuery(req.URL.Query())
		if logMsg {
			msg := err.Error()
			if ieErr, ok := err.(gmerrors.InternalErrorMessage); ok {
				msg = ieErr.InternalErrorMessage()
			}
			r.log.Err("api_response_error",
				"method", req.Method,
				"path", req.URL.Path,
				"query", sanitizedQuery.Encode(),
				"err", msg,
				"http_request_id", reqID,
			)
		}
	}()

	// not using a type switch here as we may have types that satisfy the behavior but
	// are purposely not a valid error type as they return false, which indicates this.
	// most useful when receiving a PGErr type that does fulfill more than 1 behavior potentially.
	if httpErr, ok := err.(gmerrors.HTTP); ok {
		r.With(w, req, httpErr.StatusCode(), newErr(reqID, err.Error()))
		return
	}

	if nfErr, ok := err.(gmerrors.NotFounder); ok && nfErr.NotFound() {
		logMsg = false
		r.With(w, req, http.StatusNotFound, newErr(reqID, "resource not found"))
		return
	}

	if exErr, ok := err.(gmerrors.Exister); ok && exErr.Exists() {
		r.With(w, req, http.StatusUnprocessableEntity, newErr(reqID, "resource exists"))
		return
	}

	if cfErr, ok := err.(gmerrors.Conflicter); ok && cfErr.Conflict() {
		r.With(w, req, http.StatusUnprocessableEntity, newErr(reqID, "resource conflict"))
		return
	}

	if tmpErr, ok := err.(gmerrors.Temporarier); ok && tmpErr.Temporary() {
		r.With(w, req, http.StatusServiceUnavailable, newErr(reqID, "service unavailable"))
		return
	}

	r.With(w, req, http.StatusInternalServerError, newErr(reqID, "internal server error"))
}

// OK responds with status code OK and JSON response indicating the operation succeeded
func (r *Responder) OK(w http.ResponseWriter, req *http.Request, opts ...Option) {
	r.With(w, req, http.StatusOK, struct {
		OK bool `json:"ok"`
	}{
		OK: true,
	}, opts...)
}
