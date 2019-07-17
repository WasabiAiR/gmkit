package errors

// NotFounder determines if an error exhibits the behavior of a resource not found err.
type NotFounder interface {
	NotFound() bool
}

// Exister determines if an error exhibits the behavior of a resource already exists err.
type Exister interface {
	Exists() bool
}

// Temporarier determines if an error exhibits the behavior of a system with services that are
// unreachable for whatever reason and are temporarily unavailable.
type Temporarier interface {
	Temporary() bool
}

// Conflicter determines if an error exhibits the behavior of a conflict err. This corresponds
// to times when you receive unexpected errors or an error of high severity.
type Conflicter interface {
	Conflict() bool
}

// Retrier determines if an error exhibits the behavior of an error that is safe to retry.
// TODO: a lot of times these types of errors may respond with a time/duration for when a retry
// is safe to commence. May be worth considering.
type Retrier interface {
	Retry() bool
}

// HTTP determines if an error exhibits the behavior of an error that provides an HTTP status code and
// the error message itself is safe for returning to a client.
type HTTP interface {
	StatusCode() int
}

// InternalErrorMessage provides a message separate from that returned by a call
// to Error() that is for internal (non-client) consumption.
type InternalErrorMessage interface {
	InternalErrorMessage() string
}
