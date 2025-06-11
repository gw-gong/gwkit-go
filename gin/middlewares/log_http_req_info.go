package middlewares

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/gw-gong/gwkit-go/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogHttpReqInfo(c *gin.Context) {
	start := time.Now()
	logger := log.GetLoggerFromCtx(c.Request.Context())

	// Read request body (with size limit)
	var bodyStr string
	if c.Request.Body != nil {
		body, _ := io.ReadAll(io.LimitReader(c.Request.Body, 1024))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore request body
		bodyStr = string(body)
	}

	// Log request information
	logger.Info("http request started",
		// Basic request information
		zap.String("method", c.Request.Method),      // HTTP method (GET/POST/PUT/DELETE etc.)
		zap.String("path", c.Request.URL.Path),      // Request path (e.g. /api/users)
		zap.String("query", c.Request.URL.RawQuery), // URL query parameters (e.g. id=123&name=test)

		// Client information
		zap.String("ip", c.ClientIP()),                  // Client IP address
		zap.String("user_agent", c.Request.UserAgent()), // Client browser/device information
		zap.String("referer", c.Request.Referer()),      // Request source page URL

		// Request header information
		zap.String("content_type", c.ContentType()),                   // Request body format (e.g. application/json)
		zap.String("accept", c.GetHeader("Accept")),                   // Client expected response format
		zap.String("x_forwarded_for", c.GetHeader("X-Forwarded-For")), // Original client IP forwarded by proxy server

		// Request body content (limited to 1KB to avoid large files)
		zap.String("body", bodyStr),

		// Service information (for distinguishing between multiple instances)
		zap.String("app_version", "v1.0.0"),              // Application version number
		zap.String("instance_id", os.Getenv("HOSTNAME")), // Container/hostname
	)

	// Process request
	c.Next()

	// Log response information
	logger.Info("http request completed",
		// Response status
		zap.Int("status", c.Writer.Status()),      // HTTP status code (e.g. 200/404/500)
		zap.Int("response_size", c.Writer.Size()), // Response body size (bytes)

		// Performance metrics
		zap.Duration("duration", time.Since(start)), // Request processing duration (from start to end)
	)
}
