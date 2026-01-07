package internal

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/metric"
)

const (
	// MetricsRequestsCodeTotal is the metric name for requests code total
	MetricsRequestsCodeTotal = "dubbo.requests_code_total"
	// MetricsRequestsSecondsBucket is the metric name for requests seconds bucket
	MetricsRequestsSecondsBucket = "dubbo.requests_seconds_bucket"
)

type Metrics struct {
	InvokedCouner    metric.Int64Counter
	SecondsHistogram metric.Float64Histogram
}

func newMetrics(kind string, meter metric.Meter) (*Metrics, error) {
	var err error

	totalName := fmt.Sprintf("%s/%s", kind, MetricsRequestsCodeTotal)
	requestsCounter, e := newInvokedCounter(meter, totalName)
	if e != nil {
		e = fmt.Errorf("failed to create requests counter: %w", err)
		err = errors.Join(err, e)
	}

	secName := fmt.Sprintf("%s/%s", kind, MetricsRequestsSecondsBucket)
	secondsHistogram, err := newSecondsHistogram(meter, secName)
	if e != nil {
		e = fmt.Errorf("failed to create seconds histogram: %w", err)
		err = errors.Join(err, e)
	}

	return &Metrics{
		InvokedCouner:    requestsCounter,
		SecondsHistogram: secondsHistogram,
	}, nil
}

// suggest histogramName = <client/server>_requests_code_total
func newInvokedCounter(meter metric.Meter, name string) (metric.Int64Counter, error) {
	return meter.Int64Counter(name, metric.WithUnit("{call}"))
}

// newSecondsHistogram
// return metric.Float64Histogram for WithSeconds
// suggest histogramName = <client/server>_requests_seconds_bucket
func newSecondsHistogram(meter metric.Meter, name string) (metric.Float64Histogram, error) {
	return meter.Float64Histogram(
		name,
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1),
	)
}
