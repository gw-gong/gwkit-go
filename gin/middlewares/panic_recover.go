package middlewares

import (
	"context"
	"runtime/debug"

	"github.com/gw-gong/gwkit-go/log"
	common_utils "github.com/gw-gong/gwkit-go/utils/common"

	"github.com/gin-gonic/gin"
)

func PanicHandlerWithCtx(ctx context.Context, err interface{}) {
	log.Errorc(ctx, "panic", log.Any("err", err), log.String("stack", string(debug.Stack())))
}

func PanicRecover(c *gin.Context) {
	common_utils.WithRecover(c.Next, common_utils.WithPanicHandler(func(err interface{}) {
		PanicHandlerWithCtx(c.Request.Context(), err)
	}))
}
