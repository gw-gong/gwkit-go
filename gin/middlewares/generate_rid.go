package middlewares

import (
	"github.com/gw-gong/gwkit-go/log"
	str_utils "github.com/gw-gong/gwkit-go/utils/str"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ContextKeyRID = "rid" // request id
)

func GenerateRID(c *gin.Context) {
	requestId := str_utils.GenerateUUIDName()

	logger := log.GetLoggerFromCtx(c.Request.Context())
	newLogger := logger.With(zap.String(ContextKeyRID, requestId))
	reqCtx := log.SetLoggerToCtx(c.Request.Context(), newLogger)
	c.Request = c.Request.WithContext(reqCtx)

	c.Set(ContextKeyRID, requestId)

	c.Next()
}
