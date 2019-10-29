package httpc

import (
	"context"
	"errors"
	"testing"

	"github.com/graymeta/gmkit/http/httpc/httpcfakes"

	"github.com/graymeta/gmkit/backoff"
	"github.com/graymeta/gmkit/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryResponseErrors(t *testing.T) {
	_, addr := testhelpers.NewRandomPort(t)

	t.Run("without RetryResponseErrors", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoReturns(nil, errors.New("some error"))

		client := New(
			doer,
			WithBackoff(backoff.New(backoff.MaxCalls(2))),
			WithBaseURL(addr),
		)

		err := client.
			GET("/foo").
			Do(context.TODO())

		require.Error(t, err)
		require.Equal(t, 1, doer.DoCallCount())
	})

	t.Run("with RetryResponseErrors", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoReturns(nil, errors.New("some error"))

		client := New(
			doer,
			WithBackoff(backoff.New(backoff.MaxCalls(2))),
			WithBaseURL(addr),
		)

		err := client.
			GET("/foo").
			RetryResponseErrors().
			Do(context.TODO())

		require.Error(t, err)
		require.Equal(t, 2, doer.DoCallCount())
	})

	t.Run("WithRetryResponseErrors", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoReturns(nil, errors.New("some error"))

		client := New(
			doer,
			WithBackoff(backoff.New(backoff.MaxCalls(2))),
			WithBaseURL(addr),
			WithRetryResponseErrors(),
		)

		err := client.
			GET("/foo").
			Do(context.TODO())

		require.Error(t, err)
		require.Equal(t, 2, doer.DoCallCount())
	})

	t.Run("with seek params set on request", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		doer.DoReturns(nil, errors.New("some error"))

		readerSeeker := new(httpcfakes.FakeReadSeeker)
		readerSeeker.SeekStub = func(offset int64, whence int) (i int64, e error) {
			assert.Equal(t, int64(13), offset)
			assert.Equal(t, 37, whence)
			return 0, nil
		}

		client := New(
			doer,
			WithBackoff(backoff.New(backoff.MaxCalls(2))),
			WithBaseURL(addr),
			WithRetryResponseErrors(),
		)

		err := client.
			POST("/foo").
			Body(readerSeeker).
			SeekParams(13, 37).
			Do(context.TODO())

		require.Error(t, err)
		require.Equal(t, 2, doer.DoCallCount())
		require.Equal(t, 2, readerSeeker.SeekCallCount())
	})

	t.Run("WithResetSeekerToZero", func(t *testing.T) {
		doer := new(httpcfakes.FakeDoer)
		readerSeeker := new(httpcfakes.FakeReadSeeker)
		doer.DoReturns(nil, errors.New("some error"))

		client := New(
			doer,
			WithBackoff(backoff.New(backoff.MaxCalls(2))),
			WithBaseURL(addr),
			WithRetryResponseErrors(),
			WithResetSeekerToZero(),
		)

		err := client.
			POST("/foo").
			Body(readerSeeker).
			Do(context.TODO())

		require.Error(t, err)
		require.Equal(t, 2, doer.DoCallCount())
		require.Equal(t, 2, readerSeeker.SeekCallCount())
	})
}
