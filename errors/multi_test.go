package errors_test

import (
	"errors"
	"testing"

	gmerrors "github.com/graymeta/gmkit/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMulti(t *testing.T) {
	t.Run("handles nil old error", func(t *testing.T) {
		err := gmerrors.Append(nil, errors.New("something wrong"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("handles nil new error", func(t *testing.T) {
		err := gmerrors.Append(errors.New("something wrong"), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("handles nil old and new errors", func(t *testing.T) {
		err := gmerrors.Append(nil, nil)
		require.NoError(t, err)
	})

	t.Run("existing non-multi conflict error", func(t *testing.T) {
		old := &testConflictErr{}
		err := gmerrors.Append(old, errors.New("something wrong"))
		require.Error(t, err)
		require.Implements(t, (*gmerrors.Conflicter)(nil), err)
		assert.True(t, err.(gmerrors.Conflicter).Conflict())
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("existing non-multi exists error", func(t *testing.T) {
		old := &testExistsErr{}
		err := gmerrors.Append(old, errors.New("something wrong"))
		require.Error(t, err)
		require.Implements(t, (*gmerrors.Exister)(nil), err)
		assert.True(t, err.(gmerrors.Exister).Exists())
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("existing non-multi notFound error", func(t *testing.T) {
		old := &testNotFoundErr{}
		err := gmerrors.Append(old, errors.New("something wrong"))
		require.Error(t, err)
		require.Implements(t, (*gmerrors.NotFounder)(nil), err)
		assert.True(t, err.(gmerrors.NotFounder).NotFound())
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("existing non-multi retry error", func(t *testing.T) {
		old := &testRetryErr{}
		err := gmerrors.Append(old, errors.New("something wrong"))
		require.Error(t, err)
		require.Implements(t, (*gmerrors.Retrier)(nil), err)
		assert.True(t, err.(gmerrors.Retrier).Retry())
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("existing non-multi temporary error", func(t *testing.T) {
		old := &testTemporaryErr{}
		err := gmerrors.Append(old, errors.New("something wrong"))
		require.Error(t, err)
		require.Implements(t, (*gmerrors.Temporarier)(nil), err)
		assert.True(t, err.(gmerrors.Temporarier).Temporary())
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("append conflict error", func(t *testing.T) {
		new := &testConflictErr{}
		err := gmerrors.Append(errors.New("something wrong"), new)
		require.Error(t, err)
		require.Implements(t, (*gmerrors.Conflicter)(nil), err)
		assert.True(t, err.(gmerrors.Conflicter).Conflict())
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("append exists error", func(t *testing.T) {
		new := &testExistsErr{}
		err := gmerrors.Append(errors.New("something wrong"), new)
		require.Error(t, err)
		require.Implements(t, (*gmerrors.Exister)(nil), err)
		assert.True(t, err.(gmerrors.Exister).Exists())
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("append notFound error", func(t *testing.T) {
		new := &testNotFoundErr{}
		err := gmerrors.Append(errors.New("something wrong"), new)
		require.Error(t, err)
		require.Implements(t, (*gmerrors.NotFounder)(nil), err)
		assert.True(t, err.(gmerrors.NotFounder).NotFound())
		assert.Contains(t, err.Error(), "something wrong")
	})

	t.Run("append retry error", func(t *testing.T) {
		new := &testRetryErr{}
		err := gmerrors.Append(errors.New("something wrong"), new)
		require.Error(t, err)
		require.Implements(t, (*gmerrors.Retrier)(nil), err)
		assert.True(t, err.(gmerrors.Retrier).Retry())
	})

	t.Run("append temporary error", func(t *testing.T) {
		new := &testTemporaryErr{}
		err := gmerrors.Append(errors.New("something wrong"), new)
		require.Error(t, err)
		require.Implements(t, (*gmerrors.Temporarier)(nil), err)
		assert.True(t, err.(gmerrors.Temporarier).Temporary())
		assert.Contains(t, err.Error(), "something wrong")
	})
}

type testConflictErr struct{}

func (err *testConflictErr) Error() string {
	return "conflict"
}

func (err *testConflictErr) Conflict() bool {
	return true
}

type testExistsErr struct{}

func (err *testExistsErr) Error() string {
	return "exists"
}

func (err *testExistsErr) Exists() bool {
	return true
}

type testNotFoundErr struct{}

func (err *testNotFoundErr) Error() string {
	return "not found"
}

func (err *testNotFoundErr) NotFound() bool {
	return true
}

type testRetryErr struct{}

func (err *testRetryErr) Error() string {
	return "retry me!"
}

func (err *testRetryErr) Retry() bool {
	return true
}

type testTemporaryErr struct{}

func (err *testTemporaryErr) Error() string {
	return "retry me!"
}

func (err *testTemporaryErr) Temporary() bool {
	return true
}
