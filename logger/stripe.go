package logger

import (
	"fmt"

	"github.com/stripe/stripe-go"
)

// StripeLogger is a thin wrapper around our the Logger that bends it to
// the interface uses by the stripe library
type StripeLogger struct {
	fn func(msg interface{}, keyvals ...interface{}) error
}

var _ stripe.Printfer = (*StripeLogger)(nil)

// NewStripeLogger initializes a new Logger
func NewStripeLogger(l *L, level string) *StripeLogger {
	sl := &StripeLogger{}
	switch ParseLevel(level) {
	case Err:
		sl.fn = l.Err
	case Warn:
		sl.fn = l.Warn
	case Info:
		sl.fn = l.Info
	default:
		sl.fn = l.Debug
	}
	return sl
}

// Printf prints a message to the logs
func (l *StripeLogger) Printf(format string, v ...interface{}) {
	l.fn(fmt.Sprintf(format, v...))
}
