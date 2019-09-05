package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Decode decodes the body of a request into the specified object.
// If the object implements the OK interface, that method is called
// to validate the object.
func Decode(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return &ErrJSON{err}
	}
	if obj, ok := v.(OK); ok {
		if err := obj.OK(); err != nil {
			return &ErrValidationFailed{err}
		}
	}
	return nil
}

// OK checks to see if an object is valid or not.
type OK interface {
	OK() error
}

// ErrValidationFailed is returned by Decode if the object is not
// valid.
type ErrValidationFailed struct {
	err error
}

// Error reports the error string.
func (e *ErrValidationFailed) Error() string {
	return fmt.Sprintf("validation of decoded object failed: %v", e.err)
}

// StatusCode returns the http status code that would refer to a validation error.
func (e *ErrValidationFailed) StatusCode() int {
	return http.StatusUnprocessableEntity
}

// ErrJSON error is returned by Decode when the incoming JSON
// unmarshaling fails.
type ErrJSON struct {
	err error
}

// Error returns the error that unmarshal failed.
func (e *ErrJSON) Error() string {
	return fmt.Sprintf("failed to unmarshal JSON: %v", e.err)
}

// StatusCode returns the http status code that would refer to json payload that fails to decode error.
func (e *ErrJSON) StatusCode() int {
	return http.StatusBadRequest
}
