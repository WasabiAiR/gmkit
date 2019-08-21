package metrics

import (
	"testing"
	"time"

	"github.com/quipo/statsd/event"
	"github.com/stretchr/testify/require"
)

func TestMultiClient(t *testing.T) {
	count1 := 0
	count2 := 0

	mc := &MultiClient{}
	mc.Append(&mockStatsd{
		IncrFn: func(stat string, count int64) error {
			count1++
			return nil
		},
	})
	mc.Append(&mockStatsd{
		IncrFn: func(stat string, count int64) error {
			count2++
			return nil
		},
	})

	for i := 0; i < 10; i++ {
		mc.Incr("stat", 1)
	}

	require.Equal(t, count1, 10)
	require.Equal(t, count2, 10)
}

type mockStatsd struct {
	IncrFn func(string, int64) error
}

func (c *mockStatsd) CreateSocket() error {
	panic("not implemented")
}

func (c *mockStatsd) CreateTCPSocket() error {
	panic("not implemented")
}

func (c *mockStatsd) Close() error {
	panic("not implemented")
}

func (c *mockStatsd) Incr(stat string, count int64) error {
	if c.IncrFn != nil {
		return c.IncrFn(stat, count)
	}
	panic("not implemented")
}

func (c *mockStatsd) Decr(stat string, count int64) error {
	panic("not implemented")
}

func (c *mockStatsd) Timing(stat string, delta int64) error {
	panic("not implemented")
}

func (c *mockStatsd) PrecisionTiming(stat string, delta time.Duration) error {
	panic("not implemented")
}

func (c *mockStatsd) Gauge(stat string, value int64) error {
	panic("not implemented")
}

func (c *mockStatsd) GaugeDelta(stat string, value int64) error {
	panic("not implemented")
}

func (c *mockStatsd) Absolute(stat string, value int64) error {
	panic("not implemented")
}

func (c *mockStatsd) Total(stat string, value int64) error {
	panic("not implemented")
}

func (c *mockStatsd) FGauge(stat string, value float64) error {
	panic("not implemented")
}

func (c *mockStatsd) FGaugeDelta(stat string, value float64) error {
	panic("not implemented")
}

func (c *mockStatsd) FAbsolute(stat string, value float64) error {
	panic("not implemented")
}

func (c *mockStatsd) SendEvents(events map[string]event.Event) error {
	panic("not implemented")
}
