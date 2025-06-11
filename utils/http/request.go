package http

import (
	"bytes"
	"io"
	"net/http"
)

func ReadAndRestoreReqBody(req *http.Request, maxSize int64) ([]byte, error) {
    body, err := io.ReadAll(io.LimitReader(req.Body, maxSize))	// unit is bytes
    if err != nil {
        return nil, err
    }

    req.Body = io.NopCloser(bytes.NewBuffer(body))
    return body, nil
}