// +build int

package lock

import (
	"testing"
	"time"

	"github.com/graymeta/gmkit/testhelpers/redis"

	"github.com/stretchr/testify/require"
)

func TestUniqueLocker(t *testing.T) {
	expiration := 200 * time.Millisecond
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
}
