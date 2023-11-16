package metrics

import (
	"log/slog"
	"time"

	"github.com/graymeta/gmkit/logger"
	"github.com/quipo/statsd"
	"github.com/quipo/statsd/event"
)

// LoggingClient is a client that dumps stats to a Logger
type LoggingClient struct {
	fn func(msg string, keyvals ...interface{})
}

var _ statsd.Statsd = (*LoggingClient)(nil)

// NewLoggingClient creates a new logging client that will log to logger
func NewLoggingClient(l *logger.L, level string) *LoggingClient {
	client := &LoggingClient{}
	switch logger.ParseLevel(level) {
	case slog.LevelError:
		client.fn = l.Err
	case slog.LevelWarn:
		client.fn = l.Warn
	case slog.LevelInfo:
		client.fn = l.Info
	default:
		client.fn = l.Debug
	}
	return client
}

// CreateSocket is a noop
func (c *LoggingClient) CreateSocket() error {
	return nil
}

// CreateTCPSocket is a noop
func (c *LoggingClient) CreateTCPSocket() error {
	return nil
}

// Close is a noop
func (c *LoggingClient) Close() error {
	return nil
}

// Incr logs an Incr operation
func (c *LoggingClient) Incr(stat string, count int64) error {
	c.fn("Incr", stat, count)
	return nil
}

// Decr logs a Decr operation
func (c *LoggingClient) Decr(stat string, count int64) error {
	c.fn("Decr", stat, count)
	return nil
}

// Timing logs a Timing operation
func (c *LoggingClient) Timing(stat string, delta int64) error {
	c.fn("Timing", stat, delta)
	return nil
}

// PrecisionTiming logs a PrecisionTiming operation
func (c *LoggingClient) PrecisionTiming(stat string, delta time.Duration) error {
	c.fn("PrecisionTiming", stat, delta)
	return nil
}

// Gauge logs a Gauge operation
func (c *LoggingClient) Gauge(stat string, value int64) error {
	c.fn("Gauge", stat, value)
	return nil
}

// GaugeDelta logs a GaugeDelta operation
func (c *LoggingClient) GaugeDelta(stat string, value int64) error {
	c.fn("GaugeDelta", stat, value)
	return nil
}

// Absolute logs a Absolute operation
func (c *LoggingClient) Absolute(stat string, value int64) error {
	c.fn("Absolute", stat, value)
	return nil
}

// Total logs a Total operation
func (c *LoggingClient) Total(stat string, value int64) error {
	c.fn("Total", stat, value)
	return nil
}

// FGauge logs a FGauge operation
func (c *LoggingClient) FGauge(stat string, value float64) error {
	c.fn("Fguage", stat, value)
	return nil
}

// FGaugeDelta logs a FGaugeDelta operation
func (c *LoggingClient) FGaugeDelta(stat string, value float64) error {
	c.fn("FGaugeDelta", stat, value)
	return nil
}

// FAbsolute logs a FAbsolute operation
func (c *LoggingClient) FAbsolute(stat string, value float64) error {
	c.fn("FAbsolute", stat, value)
	return nil
}

// SendEvents does nothing.
func (c *LoggingClient) SendEvents(events map[string]event.Event) error {
	return nil
}
