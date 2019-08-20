package crypter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var plainInputs = []struct {
	in string
}{
	{"foo"},
	{"value2"},
	{"testing encrypt this thing 世界"},
}

func TestAESCrypter(t *testing.T) {
	crypto, err := NewAESCrypter([]byte("012345678901234567890123456789ab"))
	require.NoError(t, err)

	for _, v := range plainInputs {
		t.Logf("plain: '%s'", v.in)
		cipher, err := crypto.Encrypt(v.in)
		require.NoError(t, err)
		t.Logf("encrypt: '%s'", cipher)

		plain, err := crypto.Decrypt(cipher)
		require.NoError(t, err)
		t.Logf("decrypt: '%s'", plain)

		require.Equal(t, v.in, plain)
	}
}

func TestAESCrypterInvalidKey(t *testing.T) {
	// Too long
	_, err := NewAESCrypter([]byte("012345678901234567890123456789abcdefg"))
	require.Equal(t, ErrInvalidAESKeyLength, err)

	// Too short
	_, err = NewAESCrypter([]byte("0123456789"))
	require.Equal(t, ErrInvalidAESKeyLength, err)
}
