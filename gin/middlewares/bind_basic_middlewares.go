package middlewares

import "github.com/gin-gonic/gin"

func BindBasicMiddlewares(app *gin.Engine, logRequestBody bool) {
	app.Use(PanicRecover)
	app.Use(InjectLoggerToCtx)
	app.Use(GenerateRID)
	app.Use(LogHttpReqInfo(logRequestBody))
}
