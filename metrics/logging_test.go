package metrics

import (
	"bytes"
	"strings"
	"testing"

	"github.com/graymeta/gmkit/logger"

	"github.com/stretchr/testify/require"
)

func TestLoggingClient(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf, "somelog", "all")
	c := NewLoggingClient(l, "warn")

	c.Incr("stat", 1234)
	c.Incr("stat2", 2345)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")

	require.Len(t, lines, 2)
	require.Contains(t, lines[0], "Incr stat=1234")
	require.Contains(t, lines[1], "Incr stat2=2345")
	require.Contains(t, lines[0], "warn")
}
