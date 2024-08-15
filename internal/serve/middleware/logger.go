package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	slogg "github.com/samber/slog-gin"
)

// Logger provides a structured logging logger which can be used by Gin using
// the new slog package, allowing for easy processing of log data.
func Logger() gin.HandlerFunc {
	return slogg.NewWithConfig(
		slog.Default().WithGroup("gin"),
		slogg.Config{
			WithRequestID: true,
			Filters: []slogg.Filter{
				slogg.IgnoreStatus(
					http.StatusUnauthorized,
					http.StatusNotFound,
				),
			},
		},
	)
}
