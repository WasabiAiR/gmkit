package crypto

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestGenerateRandomURLSafeKey(t *testing.T) {
	is := is.New(t)

	key := GenerateRandomURLSafeKey(32)
	is.True(len(key) > 32)
}
