package logger

import (
	"fmt"
	"log/slog"
)

// StripeLogger is a thin wrapper around our the Logger that bends it to
// the interface uses by the stripe library
type StripeLogger struct {
	fn func(msg any, keyvals ...any)
}

// NewStripeLogger initializes a new Logger
func NewStripeLogger(l *L, level string) *StripeLogger {
	sl := &StripeLogger{}
	switch ParseLevel(level) {
	case slog.LevelError:
		sl.fn = l.Err
	case slog.LevelDebug:
		sl.fn = l.Warn
	case slog.LevelInfo:
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
