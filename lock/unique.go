package lock

import (
	"context"
	"time"

	"github.com/graymeta/gmkit/backoff"
)

// UniqueLocker wraps a Locker and uniqueID and calls the Locker's lock/unlock
// functions with the uniqueID to simplify passing around the uniqueID (pass the
// UniqueLocker instead).
type UniqueLocker struct {
	locker Locker
	unique string
}

// NewUniqueLocker initializes a UniqueLocker.
func NewUniqueLocker(l Locker, uniqueID string) *UniqueLocker {
	return &UniqueLocker{
		locker: l,
		unique: uniqueID,
	}
}

// Lock attempts to acquire the lock identified by name. Returns true if the lock
// was acquired.
func (l *UniqueLocker) Lock(name string, duration time.Duration) (bool, error) {
	return l.locker.Lock(name, l.unique, duration)
}

// Unlock attempts to unlock the lock identified by name.
func (l *UniqueLocker) Unlock(name string) error {
	return l.locker.Unlock(name, l.unique)

}

// Refresh will periodically refresh the lock (attempts to refresh at refreshPeriod)
// until the context is cancelled. This assumes the UniqueLocker has already
// acquired the lock. This method should be run in a goroutine.
func (l *UniqueLocker) Refresh(ctx context.Context, boff backoff.Backoffer, lockName string, lockDuration, refreshPeriod time.Duration) {
	ticker := time.NewTicker(refreshPeriod)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			boff.BackoffCtx(ctx, func(context.Context) error {
				_, err := l.Lock(lockName, lockDuration)
				return err
			})
		}
	}
}
