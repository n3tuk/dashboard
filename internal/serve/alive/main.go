package alive

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Attach takes a reference to the Gin engine and attaches all the expected
// endpoints which cam be used by clients through this package.
func Attach(r *gin.Engine) {
	r.GET("/alive", alive)
	// Set up the default 404 handler here too
	r.NoRoute(notFound)
}

// alive provides a basic endpoint that just returns a 200 OK response with a
// {"status":"alive"} JSON response, without processing, allowing to test if the
// web service it up and responding to requests, regardless of any other
// downstream service.
func alive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// notFound provides the default handler for requests to the API which have no
// route, and therefore cannot be processed, necessitating a 404 (Page Not
// Found) response back to the client.
func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "page-not-found",
		"code":    http.StatusNotFound,
		"message": "Page not found",
	})
}
