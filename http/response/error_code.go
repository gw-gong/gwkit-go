package response

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

var (
	Success = NewErrorCode(0, "success", http.StatusOK)

	// Client errors (4xx)
	ErrParam                 = NewErrorCode(10000, "param error", http.StatusBadRequest)
	ErrInvalidJSON           = NewErrorCode(10001, "invalid json format", http.StatusBadRequest)
	ErrInvalidQueryParam     = NewErrorCode(10002, "invalid query parameter", http.StatusBadRequest)
	ErrMissingRequiredParam  = NewErrorCode(10003, "missing required parameter", http.StatusBadRequest)
	ErrUnauthorized          = NewErrorCode(10010, "unauthorized", http.StatusUnauthorized)
	ErrTokenExpired          = NewErrorCode(10011, "token expired", http.StatusUnauthorized)
	ErrInvalidToken          = NewErrorCode(10012, "invalid token", http.StatusUnauthorized)
	ErrForbidden             = NewErrorCode(10020, "forbidden", http.StatusForbidden)
	ErrPermissionDenied      = NewErrorCode(10021, "permission denied", http.StatusForbidden)
	ErrNotFound              = NewErrorCode(10030, "resource not found", http.StatusNotFound)
	ErrMethodNotAllowed      = NewErrorCode(10040, "method not allowed", http.StatusMethodNotAllowed)
	ErrConflict              = NewErrorCode(10050, "resource conflict", http.StatusConflict)
	ErrTooManyRequests       = NewErrorCode(10060, "too many requests", http.StatusTooManyRequests)
	ErrRequestEntityTooLarge = NewErrorCode(10070, "request entity too large", http.StatusRequestEntityTooLarge)

	// Server Error (5xx)
	ErrInternal           = NewErrorCode(20000, "internal server error", http.StatusInternalServerError)
	ErrDatabase           = NewErrorCode(20001, "database error", http.StatusInternalServerError)
	ErrCacheService       = NewErrorCode(20002, "cache service error", http.StatusInternalServerError)
	ErrThirdPartyService  = NewErrorCode(20003, "third-party service error", http.StatusInternalServerError)
	ErrBadGateway         = NewErrorCode(20010, "bad gateway", http.StatusBadGateway)
	ErrServiceUnavailable = NewErrorCode(20020, "service unavailable", http.StatusServiceUnavailable)
	ErrGatewayTimeout     = NewErrorCode(20030, "gateway timeout", http.StatusGatewayTimeout)
	ErrUnknown            = NewErrorCode(29999, "unknown error", http.StatusInternalServerError)
)
