package crypter

import (
	"encoding/base64"

	"github.com/gtank/cryptopasta"
	"github.com/pkg/errors"
)

// ErrInvalidAESKeyLength is the error returned when the encryption key isn't
// the correct length
var ErrInvalidAESKeyLength = errors.New("Invalid key length")

type aesCrypter struct {
	key *[32]byte
}

// NewAESCrypter initializes a new Crypter
func NewAESCrypter(key []byte) (Crypter, error) {
	if len(key) != 32 {
		return nil, ErrInvalidAESKeyLength
	}

	c := &aesCrypter{}
	var dest [32]byte
	copy(dest[:], key[:32])
	c.key = &dest

	return c, nil
}

// Encrypt encrypts the data string and returns the base64 encoded ciphertext
func (c *aesCrypter) Encrypt(data string) (string, error) {
	ciphertext, err := cryptopasta.Encrypt([]byte(data), c.key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt accepts a base64 encoded ciphertext string and attempts to decode and decrypt it
func (c *aesCrypter) Decrypt(ciphertext string) (string, error) {
	unencoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errors.Wrap(err, "decoding base64 ciphertext")
	}

	plaintext, err := cryptopasta.Decrypt([]byte(unencoded), c.key)
	if err != nil {
		return "", errors.Wrap(err, "decrypting ciphertext")
	}

	return string(plaintext), nil
}
