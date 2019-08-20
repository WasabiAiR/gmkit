package uuid

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	uuid, err := TimestampUUID()
	require.NoError(t, err)
	require.Len(t, uuid, 32)
}
