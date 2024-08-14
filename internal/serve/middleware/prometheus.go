package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PrometheusDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	PrometheusCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of HTTP requests.",
	}, []string{"path", "status"})
)

// PrometheusMiddleware provides instrumentation for the API calls made to a
// connected service, counting both the number of requests being processed, the
// number requested in total, and the time taken to process those requests.
func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var timer *prometheus.Timer

		if c.FullPath() != "" {
			timer = prometheus.NewTimer(
				PrometheusDuration.WithLabelValues(c.FullPath()),
			)
		}

		c.Next()

		if timer != nil {
			timer.ObserveDuration()
			PrometheusCounter.WithLabelValues(c.FullPath(), fmt.Sprintf("%d", c.Writer.Status())).Inc()
		}
	}
}

// Measure abstracts the HTTP handler implementation by only requesting a reporter, this
// reporter will return the required data to be measured.
// it accepts a next function that will be called as the wrapped logic before and after
// measurement actions.
// func (m Middleware) Measure(handlerID string, reporter Reporter, next func()) {
// 	ctx := reporter.Context()
//
// 	// func (r *reporter) Method() string { return r.c.Request.Method }
// 	//
// 	// func (r *reporter) Context() context.Context { return r.c.Request.Context() }
// 	//
// 	// func (r *reporter) URLPath() string { return r.c.FullPath() }
// 	//
// 	// func (r *reporter) StatusCode() int { return r.c.Writer.Status() }
// 	//
// 	// func (r *reporter) BytesWritten() int64 { return int64(r.c.Writer.Size()) }
//
// 	// If there isn't predefined handler ID we
// 	// set that ID as the URL path.
// 	hid := handlerID
// 	if handlerID == "" {
// 		hid = reporter.URLPath()
// 	}
//
// 	// Measure inflights if required.
// 	if !m.cfg.DisableMeasureInflight {
// 		props := metrics.HTTPProperties{
// 			Service: m.cfg.Service,
// 			ID:      hid,
// 		}
// 		m.cfg.Recorder.AddInflightRequests(ctx, props, 1)
// 		defer m.cfg.Recorder.AddInflightRequests(ctx, props, -1)
// 	}
//
// 	// Start the timer and when finishing measure the duration.
// 	start := time.Now()
// 	defer func() {
// 		duration := time.Since(start)
//
// 		// If we need to group the status code, it uses the
// 		// first number of the status code because is the least
// 		// required identification way.
// 		var code string
// 		if m.cfg.GroupedStatus {
// 			code = fmt.Sprintf("%dxx", reporter.StatusCode()/100)
// 		} else {
// 			code = strconv.Itoa(reporter.StatusCode())
// 		}
//
// 		props := metrics.HTTPReqProperties{
// 			Service: m.cfg.Service,
// 			ID:      hid,
// 			Method:  reporter.Method(),
// 			Code:    code,
// 		}
// 		m.cfg.Recorder.ObserveHTTPRequestDuration(ctx, props, duration)
//
// 		// Measure size of response if required.
// 		if !m.cfg.DisableMeasureSize {
// 			m.cfg.Recorder.ObserveHTTPResponseSize(ctx, props, reporter.BytesWritten())
// 		}
// 	}()
//
// 	// Call the wrapped logic.
// 	next()
// }.
