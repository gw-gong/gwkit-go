package response

import (
	"github.com/gw-gong/gwkit-go/gin/middlewares"
	gwkit_res "github.com/gw-gong/gwkit-go/http/response"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_str "github.com/gw-gong/gwkit-go/utils/str"

	"github.com/gin-gonic/gin"
)

type ClientResponse struct {
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
	ErrDetails interface{} `json:"err_details,omitempty"`
}

func responseJson(c *gin.Context, err *gwkit_res.ErrorCode, data interface{}, errDetails interface{}) {
	requestID := getRequestID(c)
	c.Writer.Header().Set(middlewares.HttpHeaderRID, requestID)
	c.JSON(err.HttpStatus, ClientResponse{
		Code:       err.Code,
		Msg:        err.Msg,
		Data:       data,
		ErrDetails: errDetails,
	})
}

// getRequestID retrieves request ID from context or generates a new one
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(middlewares.HttpHeaderRID); exists {
		if rid, ok := requestID.(string); ok {
			return rid
		}
	}
	newRequestID := gwkit_str.GenerateULID()
	log.Warnc(c.Request.Context(), "request id not found, generate a new one", log.Str(middlewares.LoggerKeyRID, newRequestID))
	return newRequestID
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
