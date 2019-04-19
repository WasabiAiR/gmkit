package testhelpers

import (
	"testing"

	"github.com/graymeta/gmkit/crypter"

	"github.com/stretchr/testify/require"
)

const testEncryptionKey = "012345678901234567890123456789ab"

// Crypter returns a crypter.Crypter that has been configured with a test key
func Crypter(t *testing.T) crypter.Crypter {
	enc, err := crypter.NewAESCrypter([]byte(testEncryptionKey))
	require.NoError(t, err)
	return enc
}
