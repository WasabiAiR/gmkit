package metrics

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/graymeta/env"
	"github.com/quipo/statsd"
)

// DefaultStatsd is the default statsd implementation
var DefaultStatsd statsd.Statsd

func init() {
	service := os.Getenv("gm_service")
	if service == "" {
		service = "unknown"
	}
	DefaultStatsd = Statsd(service)
}

func emitRuntimeStats() {
	// Export number of Goroutines
	Gauge("runtime.num_goroutines", int64(runtime.NumGoroutine()))

	// Export memory stats
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	Gauge("runtime.alloc_bytes", int64(stats.Alloc))
	Gauge("runtime.sys_bytes", int64(stats.Sys))
	Gauge("runtime.malloc_count", int64(stats.Mallocs))
	Gauge("runtime.free_count", int64(stats.Frees))
	Gauge("runtime.heap_objects", int64(stats.HeapObjects))
	Gauge("runtime.total_gc_pause_ns", int64(stats.PauseTotalNs))
	Gauge("runtime.total_gc_runs", int64(stats.NumGC))
}

// Statsd returns a ready to use Statsd interface,
// in case of error connecting to the statsd instance
// returns a no operation interface,
// Connection with Statsd is UDP (fire and forget)
// is recommended to use just one instance, to send all the stats,
// and call Close when the service stops to
func Statsd(service string) statsd.Statsd {
	statsdClient := statsd.NewStatsdClient(
		env.GetenvWithDefault("statsd_host", "127.0.0.1:8125"),
		fmt.Sprintf(
			"%s.%s.",
			env.GetenvWithDefault("gm_env", "prod"),
			service,
		),
	)
	if err := statsdClient.CreateSocket(); err != nil {
		return statsd.NoopClient{}
	}
	return statsdClient
}

// StopFunc stops the timer started with NewTimer and publishes the result
type StopFunc func()

// NewTimer convenience function start and publish Timing stat, typical use is:
//
//	stop := metrics.NewTimer("my_metric")
//	defer stop()
func NewTimer(name string) StopFunc {
	start := MSTime()
	return func() {
		Timing(name, Duration(start))
	}
}

// Constants for computing the size buckets
const (
	Byte  = 1
	KByte = Byte * 1000
	MByte = KByte * 1000
	GByte = MByte * 1000
	TByte = GByte * 1000
)

// Range returns a size range string based on the input size in bytes
func Range(size int64) string {
	var rang string
	switch {
	case size <= 10*MByte:
		rang = "lt_10MB"
	case (size > 10*MByte) && (size <= 50*MByte):
		rang = "10_50MB"
	case (size > 50*MByte) && (size <= 100*MByte):
		rang = "50_100MB"
	case (size > 100*MByte) && (size <= 500*MByte):
		rang = "100_500MB"
	case (size > 500*MByte) && (size <= 1*GByte):
		rang = "500_1000MB"
	case (size > 1*GByte) && (size <= 5*GByte):
		rang = "1_5GB"
	case (size > 5*GByte) && (size <= 10*GByte):
		rang = "5_10GB"
	case (size > 10*GByte) && (size <= 50*GByte):
		rang = "10_50GB"
	case (size > 50*GByte) && (size <= 100*GByte):
		rang = "50_100GB"
	case (size > 100*GByte) && (size <= 500*GByte):
		rang = "100_500GB"
	case (size > 500*GByte) && (size <= 1*TByte):
		rang = "500_1000GB"
	case (size > 1*TByte) && (size <= 2*TByte):
		rang = "1_2TB"
	default:
		rang = "gt_2TB"
	}

	return rang
}

// NewTimerBySize convenience function start and publish Timing stat bucketed by byte size, typical use is:
//
//	stop := metrics.NewTimerBySize("my_metric", size)
//	defer stop()
func NewTimerBySize(name string, size int64) StopFunc {
	return NewTimer(fmt.Sprintf("%s.%s", name, Range(size)))
}

// IncrPanic increment number of panics of the app
func IncrPanic() {
	Incr("panics", 1)
}

// CreateSocket calls CreateSocket on DefaultStatsd
func CreateSocket() error {
	return DefaultStatsd.CreateSocket()
}

// Close calls Close on DefaultStatsd
func Close() error {
	return DefaultStatsd.Close()
}

// Incr calls Incr on DefaultStatsd
func Incr(stat string, count int64) error {
	return DefaultStatsd.Incr(stat, count)
}

// Decr calls Decr on DefaultStatsd
func Decr(stat string, count int64) error {
	return DefaultStatsd.Decr(stat, count)
}

// Timing calls Timing on DefaultStatsd
func Timing(stat string, delta int64) error {
	return DefaultStatsd.Timing(stat, delta)
}

// PrecisionTiming calls PrecisionTiming on DefaultStatsd
func PrecisionTiming(stat string, delta time.Duration) error {
	return DefaultStatsd.PrecisionTiming(stat, delta)
}

// Gauge calls Gauge on DefaultStatsd
func Gauge(stat string, value int64) error {
	return DefaultStatsd.Gauge(stat, value)
}

// GaugeDelta calls GaugeDelta on DefaultStatsd
func GaugeDelta(stat string, value int64) error {
	return DefaultStatsd.GaugeDelta(stat, value)
}

// Absolute calls Absolute on DefaultStatsd
func Absolute(stat string, value int64) error {
	return DefaultStatsd.Absolute(stat, value)
}

// Total calls Total on DefaultStatsd
func Total(stat string, value int64) error {
	return DefaultStatsd.Total(stat, value)
}

// FGauge calls FGauge on DefaultStatsd
func FGauge(stat string, value float64) error {
	return DefaultStatsd.FGauge(stat, value)
}

// FGaugeDelta calls FGaugeDelta on DefaultStatsd
func FGaugeDelta(stat string, value float64) error {
	return DefaultStatsd.FGaugeDelta(stat, value)
}

// FAbsolute calls FAbsolute on DefaultStatsd
func FAbsolute(stat string, value float64) error {
	return DefaultStatsd.FAbsolute(stat, value)
}
