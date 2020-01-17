package bulk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadRequest(t *testing.T) {
	tests := []struct {
		description string
		ids         []string
		expected    error
	}{
		{"normal", []string{"abc123", "cde1234"}, nil},
		{"nil ids", nil, errIDsNotSpecified},
		{"blank ids", []string{"abc123", ""}, errBlankIDs},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require.Equal(t, tt.expected, ReadRequest{IDs: tt.ids}.OK())
		})
	}
}

func TestResponseError(t *testing.T) {
	tests := []struct {
		description string
		err         error
		expected    *ResponseError
	}{
		{"nil error", nil, nil},
		{"not found", notFound{true}, &ResponseError{Message: "not found", Type: ErrorNotFound}},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			res := NewResponseError(tt.err)
			if tt.expected == nil {
				require.Nil(t, res)
				return
			}

			require.Equal(t, tt.expected.Message, res.Message)
			require.Equal(t, tt.expected.Type, res.Type)
		})
	}
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		description string
		err         error
		expected    string
	}{
		{"nil error", nil, ""},
		{"not found interface and not found", notFound{true}, ErrorNotFound},
		{"not found interface and found", notFound{false}, ErrorInternalServerError},
		{"http error - 404", httpErr{http.StatusNotFound}, ErrorNotFound},
		{"http error - 403", httpErr{http.StatusForbidden}, ErrorForbidden},
		{"http error - 500", httpErr{http.StatusInternalServerError}, ErrorInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			require.Equal(t, tt.expected, ClassifyError(tt.err))
		})
	}
}

type notFound struct {
	notFound bool
}

func (e notFound) Error() string {
	return "not found"
}

func (e notFound) NotFound() bool {
	return e.notFound
}

type httpErr struct {
	code int
}

func (e httpErr) Error() string {
	return fmt.Sprintf("%d", e.code)
}

func (e httpErr) StatusCode() int {
	return e.code
}
