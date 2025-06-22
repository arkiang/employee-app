package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware allows cross-origin requests from allowed origins.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// You can customize the origin depending on your frontend (e.g. from config)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // replace * with frontend origin if needed
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-USER-ID, X-ROLE")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			// Preflight request
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}