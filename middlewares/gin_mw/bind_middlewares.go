package gin_mw

import "github.com/gin-gonic/gin"

func BindMiddlewares(app *gin.Engine, middlewares ...gin.HandlerFunc) {
	if len(middlewares) > 0 {
		for _, middleware := range middlewares {
			app.Use(middleware)
		}
		return
	}
	app.Use(PanicRecover)
	app.Use(InjectLoggerToCtx)
	app.Use(GenerateRID)
	app.Use(LogHttpReqInfo)
}