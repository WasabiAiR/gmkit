package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// BCryptCost is the cost factor used for generating bcrypted passwords
var BCryptCost = bcrypt.DefaultCost

func init() {
	assertAvailablePRNG()
}

func assertAvailablePRNG() {
	// Assert that a cryptographically secure PRNG is available.
	// Panic otherwise.
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

// HashPassword bcrypts a password using bcrypt.DefaultCost as the cost factor
func HashPassword(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), BCryptCost)
}

// CompareHashAndPassword compares a password to a bcrypted hash of the password
func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

// GenerateRandomKey generates a random key of length strength
func GenerateRandomKey(strength int) []byte {
	k := make([]byte, strength)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}

// GenerateRandomURLSafeKey generates a random key of length strength and returns the base64 encoded string
func GenerateRandomURLSafeKey(strength int) string {
	return base64.RawURLEncoding.EncodeToString(GenerateRandomKey(strength))
}

// Hash computes the md5 checksum of the given strings
func Hash(args ...string) string {
	h := md5.New()
	for _, arg := range args {
		io.WriteString(h, arg)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
