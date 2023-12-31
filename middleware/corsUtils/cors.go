package corsUtils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AllowAllCORS is a middleware function that allows any CORS requests.
func AllowAllCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// If it's OPTIONS method, just return status 204 (No Content)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusAccepted)
			return
		}

		c.Next()
	}
}
