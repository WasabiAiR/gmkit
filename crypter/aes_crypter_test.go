package crypter

import (
	"testing"

	"github.com/cheekybits/is"
)

var plainInputs = []struct {
	in string
}{
	{"foo"},
	{"value2"},
	{"testing encrypt this thing 世界"},
}

func TestAESCrypter(t *testing.T) {
	is := is.New(t)
	crypto, err := NewAESCrypter([]byte("012345678901234567890123456789ab"))
	is.NoErr(err)

	for _, v := range plainInputs {
		t.Logf("plain: '%s'", v.in)
		cipher, err := crypto.Encrypt(v.in)
		is.NoErr(err)
		t.Logf("encrypt: '%s'", cipher)

		plain, err := crypto.Decrypt(cipher)
		is.NoErr(err)
		t.Logf("decrypt: '%s'", plain)

		is.Equal(v.in, plain)
	}
}

func TestAESCrypterInvalidKey(t *testing.T) {
	is := is.New(t)

	// Too long
	_, err := NewAESCrypter([]byte("012345678901234567890123456789abcdefg"))
	is.Err(err)
	is.Equal(err, ErrInvalidAESKeyLength)

	// Too short
	_, err = NewAESCrypter([]byte("0123456789"))
	is.Err(err)
	is.Equal(err, ErrInvalidAESKeyLength)
}
