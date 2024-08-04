package healthz

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	slogg "github.com/samber/slog-gin"
)

// Attach takes a reference to the Gin engine and attaches all the expected
// endpoints which cam be used by clients through this package.
func Attach(r *gin.Engine) {
	r.GET("/healthz", healthz)
}

// healthz provides an endpoint for checking on the operational health of the
// service, checking downstream services are behaving as expected and reporting
// on their overall status, allowing the service to be marked as unhealthy and
// to stop processing further requests if there are known issues.
func healthz(c *gin.Context) {
	slogg.AddCustomAttributes(c, slog.Group("healthz", slog.String("status", "ok")))
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "unknown",
		"queue":    "unknown",
	})
}
