package middlewares

import (
	"context"
	"runtime/debug"

	gwkit_res "github.com/gw-gong/gwkit-go/gin/response"
	"github.com/gw-gong/gwkit-go/http/err_code"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/util/common"

	"github.com/gin-gonic/gin"
)

func PanicHandlerWithCtx(ctx context.Context, err interface{}) {
	log.Errorc(ctx, "panic", log.Any("err", err), log.Str("stack", string(debug.Stack())))
}

func PanicRecover(c *gin.Context) {
	gwkit_common.WithRecover(c.Next, func(err interface{}) {
		gwkit_common.DefaultPanicWithCtx(c.Request.Context(), err)
		gwkit_res.ResponseError(c, err_code.ErrInternal)
	})
}
