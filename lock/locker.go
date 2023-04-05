package lock

import (
	"time"
)

// Locker defines the interface for the backend lock system.
type Locker interface {
	// Lock attempts to acquire a lock identified by name for the given duration.
	// uniqueID needs to be some sort of unique identifier that is different among
	// all hosts attempting to acquire the lock.
	Lock(name, uniqueID string, duration time.Duration) (bool, error)

	// Unlock attempts to unlock a lock, but only if the uniqueID is the one that
	// has the lock.
	Unlock(name, uniqueID string) error
}
