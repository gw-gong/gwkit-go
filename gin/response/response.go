package response

import (
	"fmt"
	"net/http"

	"github.com/gw-gong/gwkit-go/global_settings"
	"github.com/gw-gong/gwkit-go/http/err_code"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_trace "github.com/gw-gong/gwkit-go/util/trace"

	"github.com/gin-gonic/gin"
)

type ServerResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data,omitempty"`
	ErrDetail interface{} `json:"err_detail,omitempty"`
	DebugInfo interface{} `json:"debug_info,omitempty"`
}

type Option func(*ServerResponse)

func WithErrDetail(errDetail interface{}) Option {
	return func(response *ServerResponse) {
		response.ErrDetail = errDetail
	}
}

func WithErrDetailf(format string, a ...interface{}) Option {
	return func(response *ServerResponse) {
		response.ErrDetail = fmt.Sprintf(format, a...)
	}
}

func WithDebug(debugInfo interface{}) Option {
	return func(response *ServerResponse) {
		if global_settings.GetEnv() == global_settings.ENV_TEST {
			response.DebugInfo = debugInfo
		}
	}
}

func WithDebugf(format string, a ...interface{}) Option {
	return func(response *ServerResponse) {
		if global_settings.GetEnv() == global_settings.ENV_TEST {
			response.DebugInfo = fmt.Sprintf(format, a...)
		}
	}
}

func responseJson(c *gin.Context, err *err_code.ErrorCode, data interface{}, options ...Option) {
	if err == nil {
		err = &err_code.ErrorCode{
			HttpStatus: http.StatusInternalServerError,
			Code:       -1,
			Msg:        "unknown error",
		}
		log.Warnc(c.Request.Context(), "error is nil, set to unknown, this is a bug")
	}
	response := &ServerResponse{
		Code:      err.Code,
		Msg:       err.Msg,
		RequestID: gwkit_trace.GetRequestIDFromCtx(c.Request.Context()),
		Data:      data,
	}
	for _, option := range options {
		option(response)
	}
	c.JSON(err.HttpStatus, response)
}

// ResponseSuccess sends a success response with data
func ResponseSuccess(c *gin.Context, data interface{}, options ...Option) {
	responseJson(c, err_code.Success, data, options...)
}

// ResponseError sends an error response
func ResponseError(c *gin.Context, err *err_code.ErrorCode, options ...Option) {
	responseJson(c, err, nil, options...)
}
