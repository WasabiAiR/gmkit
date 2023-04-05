package logger

import (
	"fmt"

	"gopkg.in/olivere/elastic.v5"
)

// ElasticLogger wraps a L and satisfies the interface required by
// https://godoc.org/github.com/olivere/elastic#Logger
type ElasticLogger struct {
	l *L
}

var _ (elastic.Logger) = (*ElasticLogger)(nil)

// NewElasticLogger creates a new ElasticSearch compatible logger that wraps l
func NewElasticLogger(l *L) *ElasticLogger {
	return &ElasticLogger{l: l}
}

// Printf logs the message to the wrapped logger at the debug level
func (el *ElasticLogger) Printf(format string, v ...any) {
	// The ES library logs lots of stuff with newlines and other junk...this escapes
	// all that so it's a single line with the message
	el.l.Debug(fmt.Sprintf("%q", fmt.Sprintf(format, v...)))
}
