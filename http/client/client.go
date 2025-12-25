package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/util/str"
	"github.com/gw-gong/gwkit-go/util/trace"

	jsoniter "github.com/json-iterator/go"
)

const (
	defaultKeepAliveMs           = 60000
	defaultTLSHandshakeTimeoutMs = 3000

	maxLogReqBodyBytes  = 1024 // 1KB
	maxLogRespBodyBytes = 1024 // 1KB
)

type BaseHTTPClientCfg struct {
	TimeoutMs           int `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	MaxIdleConnsPerHost int `json:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost" mapstructure:"maxIdleConnsPerHost"`
	MaxConnsPerHost     int `json:"maxConnsPerHost" yaml:"maxConnsPerHost" mapstructure:"maxConnsPerHost"`
}

var defaultBaseHTTPClientCfg = &BaseHTTPClientCfg{
	TimeoutMs:           5000, // request timeout (milliseconds), including connection, reading response, etc.
	MaxIdleConnsPerHost: 10,   // maximum idle connections per host, for connection reuse
	MaxConnsPerHost:     100,  // maximum connections per host (including idle and active)
}

type BaseHTTPClient interface {
	Close()
	DoRequest(ctx context.Context, method, url string, reqJsonBody interface{}, headerItems ...HeaderItem) (httpStatus int, isTimeout bool, respBody []byte, err error)
}

type baseHTTPClient struct {
	http.Client
}

func NewBaseHTTPClient(cfg *BaseHTTPClientCfg) BaseHTTPClient {
	if cfg == nil {
		cfg = defaultBaseHTTPClientCfg
	}
	if cfg.TimeoutMs <= 0 {
		cfg.TimeoutMs = defaultBaseHTTPClientCfg.TimeoutMs
	}
	if cfg.MaxIdleConnsPerHost <= 0 {
		cfg.MaxIdleConnsPerHost = defaultBaseHTTPClientCfg.MaxIdleConnsPerHost
	}
	if cfg.MaxConnsPerHost <= 0 {
		cfg.MaxConnsPerHost = defaultBaseHTTPClientCfg.MaxConnsPerHost
	}
	return &baseHTTPClient{
		Client: http.Client{
			Timeout: time.Duration(cfg.TimeoutMs) * time.Millisecond,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
				MaxConnsPerHost:     cfg.MaxConnsPerHost,
				// TLSClientConfig:     tlsConfig, // TODO: Support tls config
				DialContext: (&net.Dialer{
					Timeout:   time.Duration(cfg.TimeoutMs) * time.Millisecond,
					KeepAlive: time.Duration(defaultKeepAliveMs) * time.Millisecond,
				}).DialContext,
				TLSHandshakeTimeout: time.Duration(defaultTLSHandshakeTimeoutMs) * time.Millisecond,
			},
		},
	}
}

func (c *baseHTTPClient) Close() {
	if c.Transport != nil {
		if transport, ok := c.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}

func (c *baseHTTPClient) DoRequest(ctx context.Context, method, url string, reqJsonBody interface{}, headerItems ...HeaderItem) (
	httpStatus int, isTimeout bool, respBody []byte, err error) {
	start := time.Now()

	var req *http.Request
	var resp *http.Response
	var reqBody []byte
	defer func() {
		var logFields []log.Field
		logFields = append(logFields, log.Str("method", method))
		logFields = append(logFields, log.Str("url", url))
		logFields = append(logFields, log.Int("http_status", httpStatus))
		logFields = append(logFields, log.Bool("is_timeout", isTimeout))
		logFields = append(logFields, log.Duration("latency", time.Since(start)))
		if req != nil {
			logFields = append(logFields, log.Any("req_headers", req.Header))
			logFields = append(logFields, log.Str("req_body", str.SubStringByByte(string(reqBody), maxLogReqBodyBytes)))
		}
		if resp != nil {
			logFields = append(logFields, log.Any("resp_headers", resp.Header))
			logFields = append(logFields, log.Str("resp_body", str.SubStringByByte(string(respBody), maxLogRespBodyBytes)))
		}
		log.Infoc(ctx, "Do request done", logFields...)
	}()

	if reqJsonBody != nil {
		reqBody, err = jsoniter.Marshal(reqJsonBody)
		if err != nil {
			return HTTPStatusUnknown, false, nil, fmt.Errorf("do request failed to marshal request body, err: %w", err)
		}
	}

	var bodyReader io.Reader
	if len(reqBody) > 0 {
		bodyReader = bytes.NewBuffer(reqBody)
	}
	req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return HTTPStatusUnknown, false, nil, fmt.Errorf("do request failed to create request, err: %w", err)
	}

	for _, headerItem := range headerItems {
		req.Header.Set(headerItem.Key, headerItem.Value)
	}

	c.setHeaderTraceInfo(ctx, req.Header)
	if len(reqBody) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err = c.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return HTTPStatusUnknown, true, nil, fmt.Errorf("do request timeout (context deadline): %w", err)
		}
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return HTTPStatusUnknown, true, nil, fmt.Errorf("do request timeout (network): %w", err)
		}
		return HTTPStatusUnknown, false, nil, err
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, false, nil, fmt.Errorf("do request failed to read response body, http_status: %d, err: %w", resp.StatusCode, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, false, nil, fmt.Errorf("http request failed, status: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return resp.StatusCode, false, respBody, nil
}

func (c *baseHTTPClient) setHeaderTraceInfo(ctx context.Context, header http.Header) {
	if requestID := trace.GetRequestIDFromCtx(ctx); requestID != "" {
		header.Set(trace.HttpHeaderRequestID, requestID)
	}
	if traceID := trace.GetTraceIDFromCtx(ctx); traceID != "" {
		header.Set(trace.HttpHeaderTraceID, traceID)
	}
}
