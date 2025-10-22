package middlewares

import (
	// "os"
	"bytes"
	"time"

	"github.com/gw-gong/gwkit-go/gin/response"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/utils/http"

	"github.com/gin-gonic/gin"
)

type LogHttpInfoOptions struct {
	LogReqBody  bool `json:"log_request_body" yaml:"log_request_body" mapstructure:"log_request_body"`
	LogRespBody bool `json:"log_response_body" yaml:"log_response_body" mapstructure:"log_response_body"`
}

// If desensitization is needed, set logRequestBody to false
func LogHttpReqInfo(options LogHttpInfoOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read request body
		requestBodyStr := "hidden"
		if options.LogReqBody && c.Request.Body != nil {
			body, err := http.ReadAndRestoreReqBody(c.Request)
			if err != nil {
				requestBodyStr = "failed to read request body"
			} else {
				requestBodyStr = string(body)
			}
		}

		// Log request information
		log.Infoc(c.Request.Context(), "http request started",
			// Basic request information
			log.Str("method", c.Request.Method),      // HTTP method (GET/POST/PUT/DELETE etc.)
			log.Str("path", c.Request.URL.Path),      // Request path (e.g. /api/users)
			log.Str("query", c.Request.URL.RawQuery), // URL query parameters (e.g. id=123&name=test)

			// Client information
			log.Str("ip", c.ClientIP()),                  // Client IP address
			log.Str("user_agent", c.Request.UserAgent()), // Client browser/device information
			log.Str("referer", c.Request.Referer()),      // Request source page URL

			// Request header information
			log.Str("content_type", c.ContentType()),                   // Request body format (e.g. application/json)
			log.Str("accept", c.GetHeader("Accept")),                   // Client expected response format
			log.Str("x_forwarded_for", c.GetHeader("X-Forwarded-For")), // Original client IP forwarded by proxy server

			// request header
			log.Any("headers", c.Request.Header),

			// Request body content
			log.Str("body", requestBodyStr),
		)

		// Replace response writer
		responseBuffer := &bytes.Buffer{}
		responseWriter := response.EnhanceResponseWriter{
			ResponseWriter: c.Writer,
			Buffer:         responseBuffer,
		}
		c.Writer = &responseWriter

		// Process request
		c.Next()

		responseBodyStr := "hidden"
		if options.LogRespBody && c.Writer != nil {
			responseBodyStr = responseBuffer.String()
		}

		// Log response information
		log.Infoc(c.Request.Context(), "http request completed",
			// Response status
			log.Int("status", c.Writer.Status()), // HTTP status code (e.g. 200/404/500)
			log.Int("size", c.Writer.Size()),     // Response body size (bytes)

			// response header
			log.Any("headers", c.Writer.Header()),

			// response body
			log.Str("body", responseBodyStr),

			// Performance metrics
			log.Duration("duration", time.Since(start)), // Request processing duration (from start to end)
		)
	}
}
