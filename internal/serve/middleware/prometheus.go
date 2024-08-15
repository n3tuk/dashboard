package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	summary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Subsystem: "http",
		Name:      "response_seconds",
		Help:      "Duration of HTTP requests.",
		//nolint:mnd // ignore
		MaxAge: 15 * time.Second,
		//nolint:mnd // ignore
		Objectives: map[float64]float64{0.25: 0.01, 0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"service"})

	duration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "response_endpoints_seconds",
		Help:      "Duration of HTTP requests.",
		//nolint:mnd // ignore
		Buckets: prometheus.ExponentialBucketsRange(0.00001, 2, 15),
	}, []string{"service", "method", "path", "status"})

	requests = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "http",
		Name:      "request_total",
		Help:      "Count of HTTP requests.",
	}, []string{"service", "method", "path", "status"})

	requestSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "request_size_bytes",
		Help:      "Size of the HTTP requests.",
		//nolint:mnd // ignore
		Buckets: prometheus.ExponentialBuckets(64, 2, 10),
	}, []string{"service", "method", "path", "status"})

	responseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "response_size_bytes",
		Help:      "Size of the HTTP responses.",
		//nolint:mnd // ignore
		Buckets: prometheus.ExponentialBuckets(2, 2, 16),
	}, []string{"service", "method", "path", "status"})

	active = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "http",
		Name:      "request_open",
		Help:      "Number of requests being actively handled.",
	}, []string{"service"})
)

// Prometheus provides instrumentation for the API calls made to a connected
// service, counting both the number of requests being processed, the number
// requested in total, and the time taken to process those requests.
func Prometheus(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := strings.ToUpper(c.Request.Method)

		path := c.FullPath()
		if path == "" {
			path = "404"
		}

		active.WithLabelValues(service).Inc()
		defer active.WithLabelValues(service).Dec()

		timer := time.Now()
		defer func(c *gin.Context, t time.Time) {
			taken := time.Since(t).Seconds()

			status := fmt.Sprintf("%d", c.Writer.Status())
			if status == "0" {
				status = "200"
			}

			responseBytes := float64(c.Writer.Size())

			requestBytes := float64(c.Request.ContentLength)
			if requestBytes < 0 {
				requestBytes = 0
			}

			requests.WithLabelValues(service, method, path, status).Inc()
			duration.WithLabelValues(service, method, path, status).Observe(taken)
			summary.WithLabelValues(service).Observe(taken)
			requestSize.WithLabelValues(service, method, path, status).Observe(requestBytes)
			responseSize.WithLabelValues(service, method, path, status).Observe(responseBytes)
		}(c, timer)

		c.Next()
	}
}
