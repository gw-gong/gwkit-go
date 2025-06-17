package middlewares

import (
	"github.com/gw-gong/gwkit-go/log"
	gwkit_str "github.com/gw-gong/gwkit-go/utils/str"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyRID = "rid" // request id
)

func GenerateRID(c *gin.Context) {
	requestId := gwkit_str.GenerateUUIDName()

	reqCtx := log.WithFields(c.Request.Context(), log.String(ContextKeyRID, requestId))
	c.Request = c.Request.WithContext(reqCtx)

	c.Set(ContextKeyRID, requestId)

	c.Next()
}
