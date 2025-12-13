package middleware

import (
	// "os"
	"bytes"
	"time"

	"github.com/gw-gong/gwkit-go/gin/res"
	"github.com/gw-gong/gwkit-go/http/util"
	"github.com/gw-gong/gwkit-go/log"

	"github.com/gin-gonic/gin"
)

type LogHttpInfoOptions struct {
	LogReqBody  bool `json:"log_request_body" yaml:"log_request_body" mapstructure:"log_request_body"`
	LogRespBody bool `json:"log_response_body" yaml:"log_response_body" mapstructure:"log_response_body"`
}

func LogHttpReqInfo(options LogHttpInfoOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestBodyStr := "hidden"
		if options.LogReqBody && c.Request.Body != nil {
			body, err := util.ReadAndRestoreReqBody(c.Request)
			if err != nil {
				requestBodyStr = "failed to read request body"
			} else {
				requestBodyStr = string(body)
			}
		}

		log.Infoc(c.Request.Context(), "http request started",
			log.Str("method", c.Request.Method),
			log.Str("path", c.Request.URL.Path),
			log.Str("query", c.Request.URL.RawQuery),
			log.Str("remote_addr", c.Request.RemoteAddr),
			log.Any("headers", c.Request.Header),
			log.Str("body", requestBodyStr),
		)

		// Replace response writer
		responseBuffer := &bytes.Buffer{}
		responseWriter := res.EnhanceResponseWriter{
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

		log.Infoc(c.Request.Context(), "http request completed",
			log.Int("status", c.Writer.Status()),
			log.Int("size", c.Writer.Size()),
			log.Any("headers", c.Writer.Header()),
			log.Str("body", responseBodyStr),
			log.Duration("duration", time.Since(start)),
		)
	}
}
