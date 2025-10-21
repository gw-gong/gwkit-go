package response

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type EnhanceResponseWriter struct {
	gin.ResponseWriter
	Buffer *bytes.Buffer
}

func (w *EnhanceResponseWriter) Write(b []byte) (int, error) {
	if w.Buffer != nil {
		w.Buffer.Write(b)
	}
	return w.ResponseWriter.Write(b)
}
