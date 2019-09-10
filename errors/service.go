package errors

import "fmt"

var (
	// NewExistsSVCErrGen is a generator function for building ExistsSVCErrs with the specified
	// resource supplied to all exists errs created.
	NewExistsSVCErrGen = func(resource string) ExistsErrGenFn {
		return func(op, msg string) *ExistsSVCErr {
			return NewExistsSVCErr(resource, op, msg)
		}
	}

	// NewNotFoundSVCErrGen is a generator function for building NotFoundSVCErrs with the specified
	// resource supplied to all not found errs created.
	NewNotFoundSVCErrGen = func(resource string) NotFoundSVCErrGenFn {
		return func(op, msg string) *NotFoundSVCErr {
			return NewNotFoundSVCErr(resource, op, msg)
		}
	}

	// NewConflictSVCErrGen is a generator function for building ConflictSVCErrs with the specified
	// resource supplied to all conflict errs created.
	NewConflictSVCErrGen = func(resource string) ConflictSVCErrGenFn {
		return func(op, msg string) *ConflictSVCErr {
			return NewConflictSVCErr(resource, op, msg)
		}
	}

	// NewTemporarySVCErrGen is a generator function for building TemporarySVCErrs with the specified
	// resource supplied to all conflict errs created.
	NewTemporarySVCErrGen = func(resource string) TemporarySVCErrGenFn {
		return func(op, msg string) *TemporarySVCErr {
			return NewTemporarySVCErr(resource, op, msg)
		}
	}
)

// ExistsSVCErr is a service error for an exists behavior.
type ExistsSVCErr struct {
	*SVCErr
}

// ExistsErrGenFn is a function for generating an exists error from
// the supplied operation and errMsg paramets.
type ExistsErrGenFn func(op, msg string) *ExistsSVCErr

// NewExistsSVCErr is a constructor for an ExistsSVCErr.
func NewExistsSVCErr(resource, op, msg string) *ExistsSVCErr {
	return &ExistsSVCErr{
		SVCErr: NewSVCErr(resource, op, msg),
	}
}

// Exists provides the exists behavior.
func (*ExistsSVCErr) Exists() bool {
	return true
}

// NotFoundSVCErr is a service error for a not found behavior.
type NotFoundSVCErr struct {
	*SVCErr
}

// NotFoundSVCErrGenFn is a function for generating a not found error from
// the supplied operation and errMsg parameters.
type NotFoundSVCErrGenFn func(op, msg string) *NotFoundSVCErr

// NewNotFoundSVCErr is a constructor for a NotFoundSVC error.
func NewNotFoundSVCErr(resource, op, msg string) *NotFoundSVCErr {
	return &NotFoundSVCErr{
		SVCErr: NewSVCErr(resource, op, msg),
	}
}

// NotFound provides the not found behavior.
func (*NotFoundSVCErr) NotFound() bool {
	return true
}

// ConflictSVCErr is a service error for a conflict behavior.
type ConflictSVCErr struct {
	*SVCErr
}

// ConflictSVCErrGenFn is a function for generating a conflict error from
// the supplied operation and errMsg parameters.
type ConflictSVCErrGenFn func(op, msg string) *ConflictSVCErr

// NewConflictSVCErr is constructor for a ConflictSVC error.
func NewConflictSVCErr(resource, op, msg string) *ConflictSVCErr {
	return &ConflictSVCErr{
		SVCErr: NewSVCErr(resource, op, msg),
	}
}

// Conflict provides the conflict behavior.
func (*ConflictSVCErr) Conflict() bool {
	return true
}

// TemporarySVCErr is a service error for a temporary behavior.
type TemporarySVCErr struct {
	*SVCErr
}

// TemporarySVCErrGenFn is a function for generating a temporary error from
// the supplied operation and errMsg parameters.
type TemporarySVCErrGenFn func(op, msg string) *TemporarySVCErr

// NewTemporarySVCErr is constructor for a TemporarySVC error.
func NewTemporarySVCErr(resource, op, msg string) *TemporarySVCErr {
	return &TemporarySVCErr{
		SVCErr: NewSVCErr(resource, op, msg),
	}
}

// Temporary provides the conflict behavior.
func (*TemporarySVCErr) Temporary() bool {
	return true
}

// SVCErr is a service error for that provides the basic data for all service errors.
type SVCErr struct {
	resouce string
	op      string
	msg     string
}

// NewSVCErr is a constructor for the SVCErr.
func NewSVCErr(resource, op, msg string) *SVCErr {
	return &SVCErr{
		resouce: resource,
		op:      op,
		msg:     msg,
	}
}

// Error provides the error behavior.
func (svcErr *SVCErr) Error() string {
	if svcErr.op == "" {
		return fmt.Sprintf("resource=%q err=%q", svcErr.resouce, svcErr.msg)
	}
	return fmt.Sprintf("%s: resource=%q err=%q", svcErr.op, svcErr.resouce, svcErr.msg)
}
