// +build int

package lock

import (
	"context"
	"testing"
	"time"

	"github.com/graymeta/gmkit/backoff"
	"github.com/graymeta/gmkit/testhelpers/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniqueLocker(t *testing.T) {
	const expiration = 200 * time.Millisecond
	const name = "foo"

	t.Run("testing expiration", func(t *testing.T) {
		pool := redis.Setup(t)
		defer pool.Close()
		l := NewRedisLocker(pool, "somePrefix:")

		u1 := NewUniqueLocker(l, "l1")
		u2 := NewUniqueLocker(l, "l2")

		result, err := u1.Lock(name, expiration)
		require.NoError(t, err)
		require.True(t, result)

		result, err = u2.Lock(name, expiration)
		require.NoError(t, err)
		require.False(t, result)

		time.Sleep(100*time.Millisecond + expiration)

		// should be able to get the lock now
		result, err = u2.Lock(name, expiration)
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("unlock", func(t *testing.T) {
		pool := redis.Setup(t)
		defer pool.Close()
		l := NewRedisLocker(pool, "somePrefix:")
		u1 := NewUniqueLocker(l, "l1")
		u2 := NewUniqueLocker(l, "l2")

		result, err := u1.Lock(name, expiration)
		require.NoError(t, err)
		require.True(t, result)

		result, err = u2.Lock(name, expiration)
		require.NoError(t, err)
		require.False(t, result)

		// try a host2 unlock - should fail
		require.Equal(t, ErrUnlock, u2.Unlock(name))

		// host1 should be able to unlock it
		require.NoError(t, u1.Unlock(name))

		// host2 should be able to immedately acquire the lock
		result, err = u2.Lock(name, expiration)
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("Refresh", func(t *testing.T) {
		l := &LockerMock{
			LockFunc: func(lockName string, uniqueID string, duration time.Duration) (bool, error) {
				assert.Equal(t, name, lockName)
				assert.Equal(t, 10*time.Minute, duration)
				return true, nil
			},
		}
		uLock := NewUniqueLocker(l, "l1")

		ctx, cancel := context.WithCancel(context.Background())

		go uLock.Refresh(ctx, backoff.New(), name, 10*time.Minute, 100*time.Millisecond)

		time.Sleep(250 * time.Millisecond)
		cancel()

		assert.Len(t, l.LockCalls(), 2)
	})
}
