package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Attach takes a reference to the Gin engine and attaches all the expected
// endpoints which cam be used by clients through this package.
func Attach(r *gin.Engine) {
	r.GET("/ping", ping)
}

// ping provides a basic endpoint that just returns a 200 OK response with a
// {"status":"pong"} JSON response, without processing, allowing to test if the
// web service it up and responding to requests, regardless of any other
// downstream service.
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "pong",
	})
}
