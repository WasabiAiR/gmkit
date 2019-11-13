package errors

import (
	"errors"
	"strings"
)

// multiErr is a container for multiple errors which preserves our gmkit/errors
// error behaviors.
type multiErr struct {
	errs []error

	conflict  bool
	exists    bool
	notFound  bool
	retry     bool
	temporary bool
}

// Append adds a new error to an existing (possibly nil) error.
func Append(old, new error) error {
	if old == nil && new == nil {
		return nil
	}

	var merr *multiErr
	if old == nil {
		merr = &multiErr{}
	}
	if !errors.As(old, &merr) {
		merr = &multiErr{}
		merr = Append(merr, old).(*multiErr)
	}

	if new == nil {
		return merr
	}

	merr.errs = append(merr.errs, new)
	var cerr Conflicter
	if errors.As(new, &cerr) {
		merr.conflict = true
	}
	var exerr Exister
	if errors.As(new, &exerr) {
		merr.exists = true
	}
	var nferr NotFounder
	if errors.As(new, &nferr) {
		merr.notFound = true
	}
	var rerr Retrier
	if errors.As(new, &rerr) {
		merr.retry = true
	}
	var terr Temporarier
	if errors.As(new, &terr) {
		merr.temporary = true
	}

	return merr
}

func (m multiErr) Error() string {
	var errs []string
	for i := range m.errs {
		errs = append(errs, m.errs[i].Error())
	}

	return strings.Join(errs, ": ")
}

func (m *multiErr) Conflict() bool {
	return m.conflict
}

func (m *multiErr) Exists() bool {
	return m.exists
}

func (m *multiErr) NotFound() bool {
	return m.notFound
}

func (m *multiErr) Retry() bool {
	return m.retry
}

func (m *multiErr) Temporary() bool {
	return m.temporary
}
