package uuid

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"
)

// TimestampUUID generates timestamp UUID, 32 bits are a timestamp in hexadecimal and 96 random
// for example:
//
//   56d43a23c379031f51b2bb406b8703dc
//   --------|-----------------------
//     time   			random
//
func TimestampUUID() (string, error) {
	now := uint32(time.Now().UTC().Unix())
	b := make([]byte, 12)
	count, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	if count != len(b) {
		err = errors.New("not enough random bytes")
	}
	return fmt.Sprintf("%08x%04x%04x%04x%04x%08x", now, b[0:2], b[2:4], b[4:6], b[6:8], b[8:]), err
}
