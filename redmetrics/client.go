package redmetrics

import (
	"fmt"

	"github.com/graymeta/gmkit/metrics"
)

// Client is a metrics client that wraps some common RED
// metrics for all services/servers/stores to consume.
type Client struct {
	resource string
}

// NewClient is a client constructor.
func NewClient(resource string) Client {
	return Client{resource: resource}
}

// NewTimer creates a new timer metric.
func (l Client) NewTimer(op string) metrics.StopFunc {
	stat := fmt.Sprintf("%s.%s_duration", l.resource, op)
	if l.resource == "" {
		stat = fmt.Sprintf("%s_duration", op)
	}
	return metrics.NewTimer(stat)
}

// IncReqs increments the number of requests by 1 for the given operation.
func (l Client) IncReqs(op string) {
	stat := fmt.Sprintf("%s.%s_requests", l.resource, op)
	if l.resource == "" {
		stat = fmt.Sprintf("%s_requests", op)
	}
	metrics.Incr(stat, 1)
}

// IncErrs increments the number of errors by 1 for the given operation.
func (l Client) IncErrs(op string) {
	stat := fmt.Sprintf("%s.%s_errs", l.resource, op)
	if l.resource == "" {
		stat = fmt.Sprintf("%s_errs", op)
	}
	metrics.Incr(stat, 1)
}

// Count increments the operation count by the specified count. Example would be
// an op of "insert_item_rows" with a count of 100. Useful in a bulk insert metric.
func (l Client) Count(op string, count int64) {
	stat := fmt.Sprintf("%s.%s", l.resource, op)
	if l.resource == "" {
		stat = op
	}
	metrics.Incr(stat, count)
}

// DecrCount decrements the operation count by the specified count.
func (l Client) DecrCount(op string, count int64) {
	stat := fmt.Sprintf("%s.%s", l.resource, op)
	if l.resource == "" {
		stat = op
	}
	metrics.Incr(stat, count)
}

// REDMetrics is a method for capturing all RED metrics for a given method.
// The RED metrics being request, errors, and duration
func (l Client) REDMetrics(op string) func(error) error {
	record := l.NewTimer(op)
	l.IncReqs(op)
	return func(err error) error {
		record()
		if err != nil {
			l.IncErrs(op)
		}
		return err
	}
}

// REDCountMetrics is a method for capturing all RED metrics for a given method
// in addition to the count metric. Often used alongside a List call or the like.
func (l Client) REDCountMetrics(op string) func(int, error) error {
	record := l.REDMetrics(op)
	return func(count int, err error) error {
		if err == nil {
			l.Count(op, int64(count))
		}
		return record(err)
	}
}
