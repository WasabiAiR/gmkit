package metrics

import (
	"time"

	"github.com/graymeta/gmkit/logger"

	"github.com/quipo/statsd/event"
)

// LoggingClient is a client that dumps stats to a Logger
type LoggingClient struct {
	fn func(msg any, keyvals ...any) error
}

// NewLoggingClient creates a new logging client that will log to logger
func NewLoggingClient(l *logger.L, level string) *LoggingClient {
	client := &LoggingClient{}
	switch logger.ParseLevel(level) {
	case logger.Err:
		client.fn = l.Err
	case logger.Warn:
		client.fn = l.Warn
	case logger.Info:
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
	return c.fn("Incr", stat, count)
}

// Decr logs a Decr operation
func (c *LoggingClient) Decr(stat string, count int64) error {
	return c.fn("Decr", stat, count)
}

// Timing logs a Timing operation
func (c *LoggingClient) Timing(stat string, delta int64) error {
	return c.fn("Timing", stat, delta)
}

// PrecisionTiming logs a PrecisionTiming operation
func (c *LoggingClient) PrecisionTiming(stat string, delta time.Duration) error {
	return c.fn("PrecisionTiming", stat, delta)
}

// Gauge logs a Gauge operation
func (c *LoggingClient) Gauge(stat string, value int64) error {
	return c.fn("Gauge", stat, value)
}

// GaugeDelta logs a GaugeDelta operation
func (c *LoggingClient) GaugeDelta(stat string, value int64) error {
	return c.fn("GaugeDelta", stat, value)
}

// Absolute logs a Absolute operation
func (c *LoggingClient) Absolute(stat string, value int64) error {
	return c.fn("Absolute", stat, value)
}

// Total logs a Total operation
func (c *LoggingClient) Total(stat string, value int64) error {
	return c.fn("Total", stat, value)
}

// FGauge logs a FGauge operation
func (c *LoggingClient) FGauge(stat string, value float64) error {
	return c.fn("Fguage", stat, value)
}

// FGaugeDelta logs a FGaugeDelta operation
func (c *LoggingClient) FGaugeDelta(stat string, value float64) error {
	return c.fn("FGaugeDelta", stat, value)
}

// FAbsolute logs a FAbsolute operation
func (c *LoggingClient) FAbsolute(stat string, value float64) error {
	return c.fn("FAbsolute", stat, value)
}

// SendEvents does nothing.
func (c *LoggingClient) SendEvents(events map[string]event.Event) error {
	return nil
}
