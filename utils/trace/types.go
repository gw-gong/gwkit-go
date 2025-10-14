package trace

type ContextKeyRequestID struct{}

const (
	LoggerFieldRequestID = "rid"
)

const (
	HttpHeaderRequestID = "X-Request-ID"
)
