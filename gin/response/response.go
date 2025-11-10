package response

import (
	"net/http"

	"github.com/gw-gong/gwkit-go/global_settings"
	"github.com/gw-gong/gwkit-go/http/err_code"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/utils/trace"

	"github.com/gin-gonic/gin"
)

type ServerResponse struct {
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	RequestID  string      `json:"request_id"`
	Data       interface{} `json:"data,omitempty"`
	ErrDetails interface{} `json:"err_details,omitempty"`
	DebugInfo  interface{} `json:"debug_info,omitempty"`
}

type Option func(*ServerResponse)

func WithDebug(debugInfo interface{}) Option {
	return func(response *ServerResponse) {
		if global_settings.GetEnv() == global_settings.ENV_TEST {
			response.DebugInfo = debugInfo
		}
	}
}

func responseJson(c *gin.Context, err *err_code.ErrorCode, data interface{}, errDetails interface{}, options ...Option) {
	if err == nil {
		err = &err_code.ErrorCode{
			HttpStatus: http.StatusInternalServerError,
			Code:       -1,
			Msg:        "unknown error",
		}
		log.Warnc(c.Request.Context(), "error is nil, set to unknown, this is a bug")
	}
	response := &ServerResponse{
		Code:       err.Code,
		Msg:        err.Msg,
		RequestID:  trace.GetRequestIDFromCtx(c.Request.Context()),
		Data:       data,
		ErrDetails: errDetails,
	}
	for _, option := range options {
		option(response)
	}
	c.JSON(err.HttpStatus, response)
}

// ResponseSuccess sends a success response with data
func ResponseSuccess(c *gin.Context, data interface{}, options ...Option) {
	responseJson(c, err_code.Success, data, nil, options...)
}

// ResponseError sends an error response
func ResponseError(c *gin.Context, err *err_code.ErrorCode, options ...Option) {
	responseJson(c, err, nil, nil, options...)
}

// ResponseErrorWithDetails sends an error response with optional data and detailed error information
func ResponseErrorWithDetails(c *gin.Context, err *err_code.ErrorCode, errDetails interface{}, options ...Option) {
	responseJson(c, err, nil, errDetails, options...)
}
