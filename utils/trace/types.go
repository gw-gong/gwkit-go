package trace

type ContextKeyRequestID struct{}
type ContextKeyTraceID struct{}

const (
	LoggerFieldRequestID = "rid"
	LoggerFieldTraceID   = "tid"
)

const (
	HttpHeaderRequestID = "X-Request-ID"
)
