package http

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	errMissingBytesPrefix = errors.New("missing bytes= prefix")
	errInvalidNumPieces   = errors.New("invalid num pieces")
)

// ParseHeaderRange parses the Range http header. Only supports the expresion of
// a single range. Returns the start and end values of the range.
// Example: "bytes=1234-4567" would return 1234 and 4567
func ParseHeaderRange(h string) (int64, int64, error) {
	if !strings.HasPrefix(h, "bytes=") {
		return 0, 0, errMissingBytesPrefix
	}

	h = strings.TrimPrefix(h, "bytes=")

	pieces := strings.Split(h, "-")
	if len(pieces) != 2 {
		return 0, 0, errInvalidNumPieces
	}

	start, err := strconv.ParseInt(pieces[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing start value: %w", err)
	}

	end, err := strconv.ParseInt(pieces[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing end value: %w", err)
	}

	return start, end, nil
}

// FormatHeaderRange generates the Range header string given start and end values
func FormatHeaderRange(start, end int64) string {
	return fmt.Sprintf("bytes=%d-%d", start, end)
}
