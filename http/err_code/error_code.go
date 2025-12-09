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
// It is recommended that the server always returns a 200 status code to avoid interference from other factors in the network request chain,
// such as the gateway returning a 502 status code, which can make error localization difficult.
var (
	Success = NewErrorCode(0, "success", http.StatusOK)
)

var (
	// Client errors
	// ErrParam            = NewErrorCode(100000000, "param error", http.StatusOK)
	// ErrPermissionDenied = NewErrorCode(100000001, "permission denied", http.StatusOK)

	// Server Error
	ErrInternal = NewErrorCode(200000000, "internal server error", http.StatusOK)
	// ErrDatabase = NewErrorCode(200000001, "database error", http.StatusOK)
)
