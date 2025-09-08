package response

import (
	gwkit_res "github.com/gw-gong/gwkit-go/http/response"
	"github.com/gw-gong/gwkit-go/log"

	"github.com/gin-gonic/gin"
)

type ClientResponse struct {
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data,omitempty"`
	ErrDetails interface{} `json:"err_details,omitempty"`
}

func responseJson(c *gin.Context, err *gwkit_res.ErrorCode, data interface{}, errDetails interface{}) {
	if err == nil {
		err = gwkit_res.ErrUnknown
		log.Warnc(c.Request.Context(), "error is nil, set to unknown, this is a bug", log.Err(err))
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
	responseJson(c, gwkit_res.Success, data, nil)
}

// ResponseError sends an error response
func ResponseError(c *gin.Context, err *gwkit_res.ErrorCode) {
	responseJson(c, err, nil, nil)
}

// ResponseErrorWithDetails sends an error response with optional data and detailed error information
func ResponseErrorWithDetails(c *gin.Context, err *gwkit_res.ErrorCode, errDetails interface{}) {
	responseJson(c, err, nil, errDetails)
}
