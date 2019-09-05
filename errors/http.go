package errors

// HTTPErr is a error that provides a status code and a message that
// is safe for clients to view.
type HTTPErr struct {
	Msg  string
	Code int
}

var _ HTTP = (*HTTPErr)(nil)

// Error provides an error message.
func (h *HTTPErr) Error() string {
	return h.Msg
}

// StatusCode provides the status code associated with the error message.
func (h *HTTPErr) StatusCode() int {
	return h.Code
}

// HTTPInternalMessageError provides a status code, message safe for clients to
// view, and wraps an error for logging internally.
type HTTPInternalMessageError struct {
	err             error
	code            int
	friendlyMessage string
}

var (
	_ HTTP                 = (*HTTPInternalMessageError)(nil)
	_ InternalErrorMessage = (*HTTPInternalMessageError)(nil)
)

// NewHTTPInternalMessageError initializes an HTTPInternalMessageError.
func NewHTTPInternalMessageError(err error, msg string, code int) *HTTPInternalMessageError {
	if err == nil {
		return nil
	}
	return &HTTPInternalMessageError{
		err:             err,
		code:            code,
		friendlyMessage: msg,
	}
}

// Error provides the client-safe error message
func (e *HTTPInternalMessageError) Error() string {
	return e.friendlyMessage
}

// InternalErrorMessage provides the internal-only wrapped error message.
func (e *HTTPInternalMessageError) InternalErrorMessage() string {
	return e.err.Error()
}

// StatusCode provides the http status code associated with this message.
func (e *HTTPInternalMessageError) StatusCode() int {
	return e.code
}
