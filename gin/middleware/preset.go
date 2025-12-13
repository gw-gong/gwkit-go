package middleware

import "github.com/gin-gonic/gin"

func BindBasicMiddlewares(app *gin.Engine, logHttpInfoOptions *LogHttpInfoOptions) {
	app.Use(SetRID)
	app.Use(PanicRecover)
	if logHttpInfoOptions == nil {
		logHttpInfoOptions = &LogHttpInfoOptions{
			LogReqBody:  true,
			LogRespBody: true,
		}
	}
	app.Use(LogHttpReqInfo(*logHttpInfoOptions))
}
