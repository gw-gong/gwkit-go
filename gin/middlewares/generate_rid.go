package middlewares

import (
	"github.com/gw-gong/gwkit-go/log"
	gwkit_str "github.com/gw-gong/gwkit-go/utils/str"

	"github.com/gin-gonic/gin"
)

const (
	LoggerKeyRID  = "rid"
	HttpHeaderRID = "X-Request-ID"
)

func GenerateRID(c *gin.Context) {
	var requestID string
	if val := c.GetHeader(HttpHeaderRID); val != "" {
		requestID = val
	} else {
		requestID = gwkit_str.GenerateULID()
	}

	reqCtx := log.WithFields(c.Request.Context(), log.Str(LoggerKeyRID, requestID))
	c.Request = c.Request.WithContext(reqCtx)

	c.Set(HttpHeaderRID, requestID)

	c.Next()
}
