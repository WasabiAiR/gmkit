package bulk

import (
	"errors"
	"net/http"

	gmerrors "github.com/graymeta/gmkit/errors"
)

// ReadRequest is a request to read multiple IDs of a single resource type.
type ReadRequest struct {
	IDs []string `json:"ids"`
}

var errIDsNotSpecified = errors.New("one or more ids must be specified")
var errBlankIDs = errors.New("blank id(s) specified")

// OK validates the request.
func (r ReadRequest) OK() error {
	if len(r.IDs) == 0 {
		return errIDsNotSpecified
	}

	for _, v := range r.IDs {
		if v == "" {
			return errBlankIDs
		}
	}

	return nil
}

// ClassifyError takes an error and tries to bucket it into one of the error categories.
func ClassifyError(err error) string {
	if err == nil {
		return ""
	}

	if nf, ok := err.(gmerrors.NotFounder); ok && nf.NotFound() {
		return ErrorNotFound
	}

	if hErr, ok := err.(gmerrors.HTTP); ok {
		switch hErr.StatusCode() {
		case http.StatusNotFound:
			return ErrorNotFound
		case http.StatusForbidden:
			return ErrorForbidden
		default:
			return ErrorInternalServerError
		}
	}

	return ErrorInternalServerError
}

// ResponseError captures information about an error for a specific resource.
type ResponseError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// NewResponseError initializes a NewResponseError
func NewResponseError(err error) *ResponseError {
	if err == nil {
		return nil
	}

	return &ResponseError{
		Message: err.Error(),
		Type:    ClassifyError(err),
	}
}

// Error categories.
const (
	// ErrorNotFound indicates the given resource was not found.
	ErrorNotFound = "NotFound"
	//ErrorForbidden indicates the given resource was not accessible by the user making the request.
	ErrorForbidden = "Forbidden"
	// ErrorInternalServerError indicates a server side error occured.
	ErrorInternalServerError = "InternalServerError"
)
