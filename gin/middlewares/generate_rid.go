package middlewares

import (
	"github.com/gw-gong/gwkit-go/log"
	str_utils "github.com/gw-gong/gwkit-go/utils/str"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	LogInfoKeyRID = "rid"
)

func GenerateRID(c *gin.Context) {
	requestId := str_utils.GenerateUUIDName()
	logger := log.GetLoggerFromCtx(c.Request.Context())
	newLogger := logger.With(zap.String(LogInfoKeyRID, requestId))
	reqCtx := log.SetLoggerToCtx(c.Request.Context(), newLogger)
	c.Request = c.Request.WithContext(reqCtx)
	c.Next()
}
