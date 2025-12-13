package middleware

import (
	"context"
	"runtime/debug"

	"github.com/gw-gong/gwkit-go/gin/res"
	"github.com/gw-gong/gwkit-go/http/code"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/util/common"

	"github.com/gin-gonic/gin"
)

func PanicHandlerWithCtx(ctx context.Context, err interface{}) {
	log.Errorc(ctx, "panic", log.Any("err", err), log.Str("stack", string(debug.Stack())))
}

func PanicRecover(c *gin.Context) {
	common.WithRecover(c.Next, func(err interface{}) {
		common.DefaultPanicWithCtx(c.Request.Context(), err)
		res.ResponseError(c, code.ErrInternal)
	})
}
