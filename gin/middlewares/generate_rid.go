package middlewares

import (
	"github.com/gw-gong/gwkit-go/log"
	gwkit_trace "github.com/gw-gong/gwkit-go/util/trace"

	"github.com/gin-gonic/gin"
)

func GenerateRID(c *gin.Context) {
	var requestID string
	if val := c.GetHeader(gwkit_trace.HttpHeaderRequestID); val != "" {
		requestID = val
	} else {
		requestID = gwkit_trace.GenerateRequestID()
	}

	// set request id to context
	reqCtx := gwkit_trace.SetRequestIDToCtx(c.Request.Context(), requestID)

	// set request id to log
	reqCtx = log.WithFieldRequestID(reqCtx, requestID)

	c.Request = c.Request.WithContext(reqCtx)

	c.Next()
}
