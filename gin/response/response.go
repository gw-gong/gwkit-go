package response

import (
	"net/http"

	"github.com/gw-gong/gwkit-go/http/err_code"
	"github.com/gw-gong/gwkit-go/log"

	"github.com/gin-gonic/gin"
)

type ClientResponse struct {
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data,omitempty"`
	ErrDetails interface{} `json:"err_details,omitempty"`
}

func responseJson(c *gin.Context, err *err_code.ErrorCode, data interface{}, errDetails interface{}) {
	if err == nil {
		err = &err_code.ErrorCode{
			HttpStatus: http.StatusInternalServerError,
			Code:       -1,
			Msg:        "unknown error",
		}
		log.Warnc(c.Request.Context(), "error is nil, set to unknown, this is a bug")
	}
	c.JSON(err.HttpStatus, ClientResponse{
		Code:       err.Code,
		Msg:        err.Msg,
		Data:       data,
		ErrDetails: errDetails,
	})
}

// ResponseSuccess sends a success response with data
func ResponseSuccess(c *gin.Context, data interface{}) {
	responseJson(c, err_code.Success, data, nil)
}

// ResponseError sends an error response
func ResponseError(c *gin.Context, err *err_code.ErrorCode) {
	responseJson(c, err, nil, nil)
}

// ResponseErrorWithDetails sends an error response with optional data and detailed error information
func ResponseErrorWithDetails(c *gin.Context, err *err_code.ErrorCode, errDetails interface{}) {
	responseJson(c, err, nil, errDetails)
}
