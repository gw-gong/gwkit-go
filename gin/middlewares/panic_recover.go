package middlewares

import (
	common_utils "github.com/gw-gong/gwkit-go/utils/common"

	"github.com/gin-gonic/gin"
)

func PanicRecover(c *gin.Context) {
    common_utils.WithRecover(c.Next)
}