package middlewares

import (
	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"
	gwkit_str "github.com/gw-gong/gwkit-go/utils/str"

	"github.com/gin-gonic/gin"
)

const (
	HttpHeaderRID = "X-Request-ID"
)

func GenerateRID(c *gin.Context) {
	var requestID string
	if val := c.GetHeader(HttpHeaderRID); val != "" {
		requestID = val
	} else {
		requestID = gwkit_str.GenerateULID()
	}

	reqCtx := gwkit_common.SetRequestIDToCtx(c.Request.Context(), requestID)
	reqCtx = log.WithFieldRequestID(reqCtx, requestID)
	c.Request = c.Request.WithContext(reqCtx)

	c.Set(HttpHeaderRID, requestID)

	c.Next()
}
