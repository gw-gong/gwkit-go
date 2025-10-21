package middlewares

import "github.com/gin-gonic/gin"

func BindBasicMiddlewares(app *gin.Engine, logHttpInfoOptions *LogHttpInfoOptions) {
	app.Use(PanicRecover)
	app.Use(GenerateRID)
	if logHttpInfoOptions == nil {
		logHttpInfoOptions = &LogHttpInfoOptions{
			LogReqBody:  true,
			LogRespBody: true,
		}
	}
	app.Use(LogHttpReqInfo(*logHttpInfoOptions))
}
