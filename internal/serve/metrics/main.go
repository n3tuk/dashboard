package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Attach takes a reference to the Gin engine and attaches all the expected
// endpoints which cam be used by clients through this package.
func Attach(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
