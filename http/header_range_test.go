package http

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseHeaderRange(t *testing.T) {
	var tests = []struct {
		desc  string
		input string
		start int64
		end   int64
		err   string
	}{
		{"normal", "bytes=1234-4567", 1234, 4567, ""},
		{"no prefix", "1234-4567", 0, 0, errMissingBytesPrefix.Error()},
		{"invalid pieces", "bytes=1234-4567-7890", 0, 0, errInvalidNumPieces.Error()},
		{"invalid pieces", "bytes=1234-4567-7890", 0, 0, errInvalidNumPieces.Error()},
		{"invalid start", "bytes=foo-4567", 0, 0, "parsing start value"},
		{"invalid end", "bytes=1234-foo", 0, 0, "parsing end value"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			start, end, err := ParseHeaderRange(tt.input)

			require.Equal(t, tt.start, start)
			require.Equal(t, tt.end, end)
			if tt.err == "" {
				return
			}

			require.True(t, strings.Contains(err.Error(), tt.err))
		})
	}
}

func TestFormatHeaderRange(t *testing.T) {
	require.Equal(t, "bytes=1234-4567", FormatHeaderRange(1234, 4567))
}
