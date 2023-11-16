package logger

import "log/slog"

// NSQLogger is a thin wrapper around our the Logger that bends it to
// the interface uses by the go-nsq library
type NSQLogger struct {
	fn func(msg string, keyvals ...any)
}

// NewNSQLogger initializes a new Logger
func NewNSQLogger(l *L, level string) *NSQLogger {
	hl := &NSQLogger{}
	switch ParseLevel(level) {
	case slog.LevelError:
		hl.fn = l.Err
	case slog.LevelWarn:
		hl.fn = l.Warn
	case slog.LevelInfo:
		hl.fn = l.Info
	default:
		hl.fn = l.Debug
	}
	return hl
}

// Output writes the message to our logger
func (l *NSQLogger) Output(calldepth int, s string) error {
	l.fn(s)
	return nil
}
