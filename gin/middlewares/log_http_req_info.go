package middlewares

import (
	// "os"
	"time"

	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/utils/http"

	"github.com/gin-gonic/gin"
)

// If desensitization is needed, set logRequestBody to false
func LogHttpReqInfo(logRequestBody bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read request body (with size limit)
		bodyStr := "hidden"
		if logRequestBody && c.Request.Body != nil {
			body, _ := http.ReadAndRestoreReqBody(c.Request, 1024)
			bodyStr = string(body)
		}

		// Log request information
		log.Infoc(c.Request.Context(), "http request started",
			// Basic request information
			log.String("method", c.Request.Method),      // HTTP method (GET/POST/PUT/DELETE etc.)
			log.String("path", c.Request.URL.Path),      // Request path (e.g. /api/users)
			log.String("query", c.Request.URL.RawQuery), // URL query parameters (e.g. id=123&name=test)

			// Client information
			log.String("ip", c.ClientIP()),                  // Client IP address
			log.String("user_agent", c.Request.UserAgent()), // Client browser/device information
			log.String("referer", c.Request.Referer()),      // Request source page URL

			// Request header information
			log.String("content_type", c.ContentType()),                   // Request body format (e.g. application/json)
			log.String("accept", c.GetHeader("Accept")),                   // Client expected response format
			log.String("x_forwarded_for", c.GetHeader("X-Forwarded-For")), // Original client IP forwarded by proxy server

			// Request body content (limited to 1KB to avoid large files)
			log.String("body", bodyStr),

			// Service information (for distinguishing between multiple instances)
			// zap.String("app_version", "v1.0.0"),              // Application version number
			// zap.String("instance_id", os.Getenv("HOSTNAME")), // Container/hostname
		)

		// Process request
		c.Next()

		// Log response information
		log.Infoc(c.Request.Context(), "http request completed",
			// Response status
			log.Int("status", c.Writer.Status()),      // HTTP status code (e.g. 200/404/500)
			log.Int("response_size", c.Writer.Size()), // Response body size (bytes)

			// Performance metrics
			log.Duration("duration", time.Since(start)), // Request processing duration (from start to end)
		)
	}
}
