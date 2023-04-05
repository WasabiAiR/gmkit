package logger

// NSQLogger is a thin wrapper around our the Logger that bends it to
// the interface uses by the go-nsq library
type NSQLogger struct {
	fn func(msg any, keyvals ...any) error
}

// NewNSQLogger initializes a new Logger
func NewNSQLogger(l *L, level string) *NSQLogger {
	hl := &NSQLogger{}
	switch ParseLevel(level) {
	case Err:
		hl.fn = l.Err
	case Warn:
		hl.fn = l.Warn
	case Info:
		hl.fn = l.Info
	default:
		hl.fn = l.Debug
	}
	return hl
}

// Output writes the message to our logger
func (l *NSQLogger) Output(calldepth int, s string) error {
	return l.fn(s)
}
