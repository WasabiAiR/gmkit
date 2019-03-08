package uuid

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestUUIDok(t *testing.T) {
	is := is.New(t)
	uuid, err := TimestampUUID()
	t.Log(uuid)

	is.NoErr(err)
	is.Equal(len(uuid), 32)
}
