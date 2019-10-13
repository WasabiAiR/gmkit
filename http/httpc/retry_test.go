package httpc

import (
	"context"
	"errors"
	"testing"

	"github.com/graymeta/gmkit/http/httpc/httpcfakes"

	"github.com/graymeta/gmkit/backoff"
	"github.com/graymeta/gmkit/testhelpers"
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
}
