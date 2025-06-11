package middlewares

import "github.com/gin-gonic/gin"

func BindBasicMiddlewares(app *gin.Engine) {
	app.Use(PanicRecover)
	app.Use(InjectLoggerToCtx)
	app.Use(GenerateRID)
	app.Use(LogHttpReqInfo(true))
}