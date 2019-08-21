package backoff

import (
	"context"
	"math/rand"
	"time"

	"github.com/graymeta/gmkit/errors"
	"github.com/graymeta/gmkit/logger"
	"github.com/graymeta/gmkit/metrics"
)

var (
	defaultInitialBackoff = 1 * time.Second
	defaultMaxBackoff     = 1 * time.Minute
	defaultMaxCalls       = 10
)

// Backoffer is an interface to abstract the Runner type into
// just a the backoff call.
type Backoffer interface {
	Backoff(fn func() error) error
	BackoffCtx(ctx context.Context, fn func(context.Context) error) error
}

// Runner is a backoff runner that will run a backoff func, safe for concurrent use,
// and provides DI for a backoff callsites.
type Runner struct {
	initDur, maxDur time.Duration
	maxCalls        int
	jitter          bool
	logger          *logger.L
}

// New returns a runner with the defined options. If no options are given,
// then the new runner is given the defaults for initial and max backoff as well
// as max calls. The defaults being:
//		InitBackoff: 1 second
//		MaxBackoff:	 1 minute
//		MaxCalls:	 10
//		Jitter: 	 false
func New(opts ...RunnerOptFn) Runner {
	r := Runner{
		initDur:  defaultInitialBackoff,
		maxDur:   defaultMaxBackoff,
		maxCalls: defaultMaxCalls,
	}

	for _, o := range opts {
		r = o(r)
	}

	return r
}

// RunnerOptFn is a functional option to set the
type RunnerOptFn func(r Runner) Runner

// InitBackoff sets the runners initial backoff time.
func InitBackoff(i time.Duration) RunnerOptFn {
	return func(r Runner) Runner {
		if i < 0 {
			return r
		}
		r.initDur = i
		return r
	}
}

// MaxBackoff sets the runners max backoff time.
func MaxBackoff(m time.Duration) RunnerOptFn {
	return func(r Runner) Runner {
		if m < 0 {
			return r
		}
		r.maxDur = m
		return r
	}
}

// MaxCalls sets the runners max call count.
func MaxCalls(m int) RunnerOptFn {
	return func(r Runner) Runner {
		if m < 0 {
			return r
		}
		r.maxCalls = m
		return r
	}
}

// Jitter sets the Backoff method to jitter the backoff duration.
func Jitter() RunnerOptFn {
	return func(r Runner) Runner {
		r.jitter = true
		return r
	}
}

// Logger sets the Logger on the Runner type.
func Logger(l *logger.L) RunnerOptFn {
	return func(r Runner) Runner {
		r.logger = l
		return r
	}
}

// New allows one to create a new runner, with any options, from an existing
// runner type.
func (r Runner) New(opts ...RunnerOptFn) Runner {
	newRunner := r
	for _, o := range opts {
		newRunner = o(newRunner)
	}
	return newRunner
}

// Backoff runs the given func in a backoff loop defined by the runner type.
func (r Runner) Backoff(fn func() error) error {
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))

	backoff := time.Duration(0)
	calls := 0
	for {
		err := fn()
		if err == nil {
			return nil
		}
		if retrier, ok := err.(errors.Retrier); ok && !retrier.Retry() {
			return err
		}

		calls++
		if r.maxCalls != 0 && calls >= r.maxCalls {
			return err
		}
		if backoff == 0 {
			backoff = r.initDur
		} else {
			backoff *= 2
		}
		if r.maxDur != 0 && backoff >= r.maxDur {
			backoff = r.maxDur
		}

		sleep := backoff
		if r.jitter {
			sleep = time.Duration(rander.Int63n(int64(backoff)))
		}
		metrics.Incr("backoffs", 1)

		time.Sleep(sleep)

		if r.logger != nil {
			errMsg := err.Error()
			if clientErr, ok := err.(*errors.ClientErr); ok {
				errMsg = clientErr.BackoffMessage()
			}

			r.logger.Warn("backoff", "calls", calls, "retry_after", backoff, "error", errMsg)
		}
	}
}

// BackoffCtx runs the given func in a backoff loop defined by the runner type.
func (r Runner) BackoffCtx(ctx context.Context, fn func(context.Context) error) error {
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))

	backoff := time.Duration(0)
	calls := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}
		if retrier, ok := err.(errors.Retrier); ok && !retrier.Retry() {
			return err
		}

		calls++
		if r.maxCalls != 0 && calls >= r.maxCalls {
			return err
		}
		if backoff == 0 {
			backoff = r.initDur
		} else {
			backoff *= 2
		}
		if r.maxDur != 0 && backoff > r.maxDur {
			backoff = r.maxDur
		}

		sleep := backoff
		if r.jitter {
			sleep = time.Duration(rander.Int63n(int64(backoff)))
		}
		metrics.Incr("backoffs", 1)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
		}

		if r.logger != nil {
			errMsg := err.Error()
			if clientErr, ok := err.(*errors.ClientErr); ok {
				errMsg = clientErr.BackoffMessage()
			}

			r.logger.Warn("backoff", "calls", calls, "retry_after", backoff, "error", errMsg)
		}
	}
}

// NoopBackoff is a backoff that is a noop for the backoff. Only calls the func provided to the backoff calls.
type NoopBackoff struct{}

var _ Backoffer = NoopBackoff{}

// Backoff is a backoff runner.
func (n NoopBackoff) Backoff(fn func() error) error {
	return fn()
}

// BackoffCtx is a backoff runner that may be interrupted via the provided context.
func (n NoopBackoff) BackoffCtx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
