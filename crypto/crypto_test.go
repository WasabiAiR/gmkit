package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomURLSafeKey(t *testing.T) {
	key := GenerateRandomURLSafeKey(32)
	require.True(t, len(key) > 32)
}
