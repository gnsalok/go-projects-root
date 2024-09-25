// middleware/middleware.go
package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs each incoming request and its duration.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log details
		duration := time.Since(startTime)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[%s] %s %d %s", method, path, status, duration)
	}
}

// AuthenticationMiddleware is a placeholder for authentication logic.
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement authentication logic (e.g., JWT verification)
		// If unauthorized:
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		// return

		// For demonstration, we'll allow all requests
		c.Next()
	}
}
