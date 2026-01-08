package middleware

import (
	"github.com/gw-gong/gwkit-go/util/trace"

	"github.com/gin-gonic/gin"
)

func SetRID(c *gin.Context) {
	var requestID string
	if val := c.GetHeader(trace.HttpHeaderRequestID); val != "" {
		requestID = val
	} else {
		requestID = trace.GenerateRequestID()
	}

	// set request id to context
	reqCtx := trace.SetRequestIDToCtx(c.Request.Context(), requestID)

	// set request id to log
	reqCtx = trace.WithLogFieldRequestID(reqCtx, requestID)

	c.Request = c.Request.WithContext(reqCtx)

	c.Next()
}
