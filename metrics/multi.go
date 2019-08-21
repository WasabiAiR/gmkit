package metrics

import (
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/quipo/statsd"
	"github.com/quipo/statsd/event"
)

// MultiClient passes metrics to multiple statsd backends
type MultiClient struct {
	clients []statsd.Statsd
}

// Append adds a new statsd client to the MultiClient
func (c *MultiClient) Append(client statsd.Statsd) {
	c.clients = append(c.clients, client)
}

// CreateSocket calls CreateSocket for each backend
func (c *MultiClient) CreateSocket() error {
	var errs error
	for _, b := range c.clients {
		err := b.CreateSocket()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// CreateTCPSocket calls CreateTCPSocket for each backend
func (c *MultiClient) CreateTCPSocket() error {
	var errs error
	for _, b := range c.clients {
		err := b.CreateTCPSocket()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// Close calls Close for each backend
func (c *MultiClient) Close() error {
	var errs error
	for _, b := range c.clients {
		err := b.Close()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// Incr calls Incr for each backend
func (c *MultiClient) Incr(stat string, count int64) error {
	var errs error
	for _, b := range c.clients {
		err := b.Incr(stat, count)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// Decr calls Decr for each backend
func (c *MultiClient) Decr(stat string, count int64) error {
	var errs error
	for _, b := range c.clients {
		err := b.Decr(stat, count)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// Timing calls Timing for each backend
func (c *MultiClient) Timing(stat string, delta int64) error {
	var errs error
	for _, b := range c.clients {
		err := b.Timing(stat, delta)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// PrecisionTiming calls PrecisionTiming for each backend
func (c *MultiClient) PrecisionTiming(stat string, delta time.Duration) error {
	var errs error
	for _, b := range c.clients {
		err := b.PrecisionTiming(stat, delta)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// Gauge calls Gauge for each backend
func (c *MultiClient) Gauge(stat string, value int64) error {
	var errs error
	for _, b := range c.clients {
		err := b.Gauge(stat, value)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// GaugeDelta calls GaugeDelta for each backend
func (c *MultiClient) GaugeDelta(stat string, value int64) error {
	var errs error
	for _, b := range c.clients {
		err := b.GaugeDelta(stat, value)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// Absolute calls Absolute for each backend
func (c *MultiClient) Absolute(stat string, value int64) error {
	var errs error
	for _, b := range c.clients {
		err := b.Absolute(stat, value)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// Total calls Total for each backend
func (c *MultiClient) Total(stat string, value int64) error {
	var errs error
	for _, b := range c.clients {
		err := b.Total(stat, value)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// FGauge calls FGauge for each backend
func (c *MultiClient) FGauge(stat string, value float64) error {
	var errs error
	for _, b := range c.clients {
		err := b.FGauge(stat, value)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// FGaugeDelta calls FGaugeDelta for each backend
func (c *MultiClient) FGaugeDelta(stat string, value float64) error {
	var errs error
	for _, b := range c.clients {
		err := b.FGaugeDelta(stat, value)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// FAbsolute calls FAbsolute for each backend
func (c *MultiClient) FAbsolute(stat string, value float64) error {
	var errs error
	for _, b := range c.clients {
		err := b.FAbsolute(stat, value)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// SendEvents does nothing.
func (c *MultiClient) SendEvents(events map[string]event.Event) error {
	var errs error
	for _, b := range c.clients {
		err := b.SendEvents(events)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}
