package middlewares

import (
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/utils/trace"

	"github.com/gin-gonic/gin"
)

func GenerateRID(c *gin.Context) {
	var requestID string
	if val := c.GetHeader(trace.HttpHeaderRequestID); val != "" {
		requestID = val
	} else {
		requestID = trace.GenerateRequestID()
	}

	reqCtx := trace.SetRequestIDToCtx(c.Request.Context(), requestID)
	reqCtx = log.WithFieldRequestID(reqCtx, requestID)
	c.Request = c.Request.WithContext(reqCtx)

	c.Set(trace.HttpHeaderRequestID, requestID)
	c.Header(trace.HttpHeaderRequestID, requestID) // response header

	c.Next()
}
