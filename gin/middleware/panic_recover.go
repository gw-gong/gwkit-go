package middleware

import (
	"github.com/gw-gong/gwkit-go/gin/res"
	"github.com/gw-gong/gwkit-go/http/code"
	"github.com/gw-gong/gwkit-go/util/common"

	"github.com/gin-gonic/gin"
)

func PanicRecover(c *gin.Context) {
	common.WithRecover(c.Next, func(err interface{}) {
		common.DefaultPanicWithCtx(c.Request.Context(), err)
		res.ResponseError(c, code.ErrInternal)
	})
}
