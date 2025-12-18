package middleware

import (
	"github.com/gw-gong/gwkit-go/gin/res"
	"github.com/gw-gong/gwkit-go/http/code"
	"github.com/gw-gong/gwkit-go/util"

	"github.com/gin-gonic/gin"
)

func PanicRecover(c *gin.Context) {
	util.WithRecover(c.Next, func(err interface{}) {
		util.DefaultPanicWithCtx(c.Request.Context(), err)
		res.ResponseError(c, code.ErrInternal)
	})
}
