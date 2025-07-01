package http

import (
	"bytes"
	"io"
	"net/http"
)

func ReadAndRestoreReqBody(req *http.Request) ([]byte, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// Close the original Body to avoid resource leaks
	if req.Body != nil {
		_ = req.Body.Close()
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
