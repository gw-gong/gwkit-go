package err_code

import "net/http"

type ErrorCode struct {
	HttpStatus int
	Code       int
	Msg        string
}

func NewErrorCode(code int, msg string, httpStatus int) *ErrorCode {
	return &ErrorCode{HttpStatus: httpStatus, Code: code, Msg: msg}
}

// Error implements the error interface
func (e *ErrorCode) Error() string {
	return e.Msg
}

// define [1 client/server][2 service][3 module][3 error]
var (
	Success = NewErrorCode(0, "success", http.StatusOK)

	// Client errors (4xx)
	ErrParam                 = NewErrorCode(100000000, "param error", http.StatusBadRequest)
	ErrInvalidJSON           = NewErrorCode(100000001, "invalid json format", http.StatusBadRequest)
	ErrInvalidQueryParam     = NewErrorCode(100000002, "invalid query parameter", http.StatusBadRequest)
	ErrMissingRequiredParam  = NewErrorCode(100000003, "missing required parameter", http.StatusBadRequest)
	ErrUnauthorized          = NewErrorCode(100000004, "unauthorized", http.StatusUnauthorized)
	ErrTokenExpired          = NewErrorCode(100000005, "token expired", http.StatusUnauthorized)
	ErrInvalidToken          = NewErrorCode(100000006, "invalid token", http.StatusUnauthorized)
	ErrForbidden             = NewErrorCode(100000007, "forbidden", http.StatusForbidden)
	ErrPermissionDenied      = NewErrorCode(100000008, "permission denied", http.StatusForbidden)
	ErrNotFound              = NewErrorCode(100000009, "resource not found", http.StatusNotFound)
	ErrMethodNotAllowed      = NewErrorCode(100000010, "method not allowed", http.StatusMethodNotAllowed)
	ErrConflict              = NewErrorCode(100000011, "resource conflict", http.StatusConflict)
	ErrTooManyRequests       = NewErrorCode(100000012, "too many requests", http.StatusTooManyRequests)
	ErrRequestEntityTooLarge = NewErrorCode(100000013, "request entity too large", http.StatusRequestEntityTooLarge)

	// Server Error (5xx)
	ErrInternal           = NewErrorCode(200000000, "internal server error", http.StatusInternalServerError)
	ErrDatabase           = NewErrorCode(200000001, "database error", http.StatusInternalServerError)
	ErrCacheService       = NewErrorCode(200000002, "cache service error", http.StatusInternalServerError)
	ErrThirdPartyService  = NewErrorCode(200000003, "third-party service error", http.StatusInternalServerError)
	ErrBadGateway         = NewErrorCode(200000004, "bad gateway", http.StatusBadGateway)
	ErrServiceUnavailable = NewErrorCode(200000005, "service unavailable", http.StatusServiceUnavailable)
	ErrGatewayTimeout     = NewErrorCode(200000006, "gateway timeout", http.StatusGatewayTimeout)
)
