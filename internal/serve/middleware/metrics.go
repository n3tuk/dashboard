package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	duration *prometheus.HistogramVec
	count    *prometheus.CounterVec
	size     *prometheus.HistogramVec
	active   *prometheus.GaugeVec
}

// NewMetrics returns a new metrics recorder that implements the recorder
// using Prometheus as the backend.
func NewMetrics(namespace string) *Metrics {
	metrics := &Metrics{
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "The latency of the HTTP requests.",
				//nolint:mnd // these are the building blocks for buckets
				Buckets: prometheus.ExponentialBuckets(0.0005, 2, 12),
			},
			[]string{"service", "handler", "method", "path", "status"},
		),
		count: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "The count of HTTP requests.",
			},
			[]string{"service", "handler", "method", "path", "status"},
		),
		size: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "response_size_bytes",
				Help:      "The size of the HTTP responses.",
				//nolint:mnd // these are the building blocks for buckets
				Buckets: prometheus.ExponentialBuckets(100, 2, 15),
			},
			[]string{"service", "handler", "method", "path", "status"},
		),
		active: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "requests_inflight",
				Help:      "The number of inflight requests being handled at the same time.",
			},
			[]string{"service", "handler"},
		),
	}

	prometheus.MustRegister(metrics.duration, metrics.count, metrics.size, metrics.active)

	return metrics
}

//nolint:lll // acknowledged
func (m Metrics) Record(ctx context.Context, service, handler, method string, status int, duration time.Duration, size int64) {
	m.Count(ctx, service, handler, method, status)
	m.Duration(ctx, service, handler, method, status, duration)
	m.Size(ctx, service, handler, method, status, size)
}

func (m Metrics) Count(_ context.Context, service, handler, method string, status int) {
	m.count.WithLabelValues(service, handler, method, strconv.Itoa(status)).Inc()
}

func (m Metrics) Duration(_ context.Context, service, handler, method string, status int, duration time.Duration) {
	m.duration.WithLabelValues(service, handler, method, strconv.Itoa(status)).Observe(duration.Seconds())
}

func (m Metrics) Size(_ context.Context, service, handler, method string, status int, size int64) {
	m.size.WithLabelValues(service, handler, method, strconv.Itoa(status)).Observe(float64(size))
}

func (m Metrics) Active(_ context.Context, service, handler string, quantity int) {
	m.active.WithLabelValues(service, handler).Add(float64(quantity))
}
