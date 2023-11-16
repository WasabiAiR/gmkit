package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ernesto-jimenez/httplogger"
)

// HTTPLogger wraps a L and satisfys the interface required by
// https://godoc.org/github.com/ernesto-jimenez/httplogger#HTTPLogger
type HTTPLogger struct {
	fn func(msg string, keyvals ...any)
}

var _ httplogger.HTTPLogger = (*HTTPLogger)(nil)

// NewHTTPLogger creates a new HTTPLogger. level is the log level at which you want
// your HTTP requests logged at.
func NewHTTPLogger(l *L, level string) *HTTPLogger {
	hl := &HTTPLogger{}
	switch ParseLevel(level) {
	case slog.LevelError:
		hl.fn = l.l.Error
	case slog.LevelWarn:
		hl.fn = l.l.Warn
	case slog.LevelInfo:
		hl.fn = l.l.Info
	default:
		hl.fn = l.l.Debug
	}
	return hl
}

// LogRequest doesn't do anything since we'll be logging replies only
func (h *HTTPLogger) LogRequest(*http.Request) {}

// LogResponse logs path, host, status code and duration in milliseconds
func (h *HTTPLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		h.fn("HTTP Request Error",
			"method", req.Method,
			"host", req.Host,
			"path", req.URL.Path,
			"status", "error",
			"durationMs", duration,
			"error", err,
		)
	} else {
		h.fn("HTTP Request",
			"method", req.Method,
			"host", req.Host,
			"path", req.URL.Path,
			"status", res.StatusCode,
			"durationMs", duration,
		)
	}
}
