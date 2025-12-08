package trace

import (
	gwkit_str "github.com/gw-gong/gwkit-go/util/str"
)

func GenerateRequestID() string {
	return gwkit_str.GenerateULID()
}

func GenerateTraceID() string {
	return gwkit_str.GenerateULID()
}
