package middlewares

import (
	"github.com/gw-gong/gwkit-go/log"

	"github.com/gin-gonic/gin"
)

func InjectLoggerToCtx(c *gin.Context) {
	reqCtx := c.Request.Context()
	reqCtx = log.SetGlobalLoggerToCtx(reqCtx)
	c.Request = c.Request.WithContext(reqCtx)
	c.Next()
}
