package healthz

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	slogg "github.com/samber/slog-gin"
)

type Health struct {
	Web         bool
	Metrics     bool
	Terminating bool
}

const (
	healthy     = "healthy"
	unhealthy   = "unhealthy"
	terminating = "terminating"
)

var health *Health

func NewHealth() *Health {
	return &Health{
		Web:         false,
		Metrics:     false,
		Terminating: false,
	}
}

// Attach takes a reference to the Gin engine and attaches all the expected
// endpoints which cam be used by clients through this package.
func Attach(r *gin.Engine, h *Health) {
	health = h

	r.GET("/healthz", healthz)
}

// healthz provides an endpoint for checking on the operational health of the
// service, checking downstream services are behaving as expected and reporting
// on their overall status, allowing the service to be marked as unhealthy and
// to stop processing further requests if there are known issues.
func healthz(c *gin.Context) {
	code := http.StatusOK
	status := healthy
	web := healthy
	metrics := healthy

	if !health.Web {
		code = http.StatusServiceUnavailable
		status = unhealthy
		web = unhealthy
	}

	if !health.Metrics {
		code = http.StatusServiceUnavailable
		status = unhealthy
		metrics = unhealthy
	}

	if health.Terminating {
		code = http.StatusGone
		status = terminating
	}

	slogg.AddCustomAttributes(c,
		slog.Group("healthz",
			slog.Int("code", code),
			slog.String("status", status),
			slog.String("web", web),
			slog.String("metrics", metrics),
		),
	)

	c.JSON(code, gin.H{
		"status":  status,
		"web":     web,
		"metrics": metrics,
	})
}
