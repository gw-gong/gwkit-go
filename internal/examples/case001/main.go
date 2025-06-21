package main

import (
	"context"

	"github.com/gw-gong/gwkit-go/gin/middlewares"
	gwkit_res "github.com/gw-gong/gwkit-go/gin/response"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"

	"github.com/gin-gonic/gin"
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	gwkit_common.ExitOnErr(context.Background(), err)
	defer syncFn()

	log.Info("应用启动成功")

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	middlewares.BindBasicMiddlewares(app, true)

	app.GET("/", func(c *gin.Context) {
		log.Infoc(c.Request.Context(), "处理请求", log.Str("path", c.Request.URL.Path))
		subProcess(c.Request.Context())
		gwkit_res.ResponseSuccess(c, "success")
	})

	err = app.Run(":8080")
	gwkit_common.ExitOnErr(context.Background(), err)
}

func subProcess(ctx context.Context) {
	log.Debugc(ctx, "子流程执行中", log.Str("sub_process", "data_processing"))
}
