// +build int

package lock

import (
	"testing"
	"time"

	"github.com/graymeta/gmkit/testhelpers/redis"

	"github.com/stretchr/testify/require"
)

func TestRedisLocker(t *testing.T) {
	expiration := 200 * time.Millisecond
	const name = "foo"

	t.Run("normal lock expiration", func(t *testing.T) {
		pool := redis.Setup(t)
		defer pool.Close()
		l := NewRedisLocker(pool, "somePrefix:")

		result, err := l.Lock(name, "host1", expiration)
		require.NoError(t, err)
		require.True(t, result)

		result, err = l.Lock(name, "host2", expiration)
		require.NoError(t, err)
		require.False(t, result)

		time.Sleep(100*time.Millisecond + expiration)

		// should be able to get the lock now
		result, err = l.Lock(name, "host2", expiration)
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("refresh lock", func(t *testing.T) {
		pool := redis.Setup(t)
		defer pool.Close()
		l := NewRedisLocker(pool, "somePrefix:")

		result, err := l.Lock(name, "host1", expiration)
		require.NoError(t, err)
		require.True(t, result)

		result, err = l.Lock(name, "host2", expiration)
		require.NoError(t, err)
		require.False(t, result)

		// set the lock again, this time 3x the expiration
		result, err = l.Lock(name, "host1", 3*expiration)
		require.NoError(t, err)
		require.True(t, result)

		time.Sleep(100*time.Millisecond + expiration)

		// should not be able to get the lock
		result, err = l.Lock(name, "host2", expiration)
		require.NoError(t, err)
		require.False(t, result)

		// should be able to get the lock now
		time.Sleep(2 * expiration)
		result, err = l.Lock(name, "host2", expiration)
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("unlock", func(t *testing.T) {
		pool := redis.Setup(t)
		defer pool.Close()
		l := NewRedisLocker(pool, "somePrefix:")

		result, err := l.Lock(name, "host1", expiration)
		require.NoError(t, err)
		require.True(t, result)

		result, err = l.Lock(name, "host2", expiration)
		require.NoError(t, err)
		require.False(t, result)

		// try a host2 unlock - should fail
		require.Equal(t, ErrUnlock, l.Unlock(name, "host2"))

		// host1 should be able to unlock it
		require.NoError(t, l.Unlock(name, "host1"))

		// host2 should be able to immedately acquire the lock
		result, err = l.Lock(name, "host2", expiration)
		require.NoError(t, err)
		require.True(t, result)
	})
}
