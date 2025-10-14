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
	"github.com/gw-gong/gwkit-go/utils/trace"
)

type BaseHTTPClient struct {
	client *http.Client
}

// If you don't know how to pass parameters, you can refer to the following content:
// - timeout: 5
// - maxIdleConnsPerHost: 10
// - maxConnsPerHost: 100
type BaseHTTPClientCfg struct {
	Timeout             int `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	MaxIdleConnsPerHost int `json:"maxIdleConnsPerHost" yaml:"maxIdleConnsPerHost" mapstructure:"maxIdleConnsPerHost"`
	MaxConnsPerHost     int `json:"maxConnsPerHost" yaml:"maxConnsPerHost" mapstructure:"maxConnsPerHost"`

	// tls config
	// TODO: add tls config

	// socks5 config
	// TODO: add socks5 config
}

func NewHTTPClient(cfg *BaseHTTPClientCfg) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			MaxConnsPerHost:     cfg.MaxConnsPerHost,
			// TLSClientConfig:     tlsConfig,
			Dial: (&net.Dialer{
				Timeout:   time.Duration(cfg.Timeout) * time.Second,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   3 * time.Second,
			ResponseHeaderTimeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// TODO: add tls config
// func NewHTTPClient(timeout time.Duration, maxIdleConnsPerHost, maxConnsPerHost int, tlsConfig *tls.Config) *http.Client {
// 	return &http.Client{
// 		Transport: &http.Transport{
// 			MaxIdleConnsPerHost: maxIdleConnsPerHost,
// 			MaxConnsPerHost:     maxConnsPerHost,
// 			TLSClientConfig:     tlsConfig,
// 			Dial: (&net.Dialer{
// 				Timeout:   timeout,
// 				KeepAlive: 60 * time.Second,
// 			}).Dial,
// 			TLSHandshakeTimeout:   3 * time.Second,
// 			ResponseHeaderTimeout: timeout,
// 		},
// 	}
// }

// TODO: add tls and socks5 config
// func NewHTTPClientWithSOCKS5(timeout time.Duration, maxIdleConnsPerHost, maxConnsPerHost int, tlsConfig *tls.Config, socks5Config *SOCKS5Config) *http.Client {
// 	var dialer proxy.Dialer
// 	if socks5Config == nil || socks5Config.Address == "" {
// 		dialer = &net.Dialer{Timeout: timeout, KeepAlive: 60 * time.Second}
// 	} else {
// 		dialer, _ = proxy.SOCKS5(
// 			"tcp",
// 			socks5Config.Address,
// 			&proxy.Auth{
// 				User:     socks5Config.User,
// 				Password: socks5Config.Password,
// 			},
// 			&net.Dialer{Timeout: timeout, KeepAlive: 60 * time.Second},
// 		)
// 	}

// 	return &http.Client{
// 		Transport: &http.Transport{
// 			MaxIdleConnsPerHost:   maxIdleConnsPerHost,
// 			MaxConnsPerHost:       maxConnsPerHost,
// 			TLSClientConfig:       tlsConfig,
// 			Dial:                  dialer.Dial,
// 			TLSHandshakeTimeout:   3 * time.Second,
// 			ResponseHeaderTimeout: timeout,
// 		},
// 	}
// }

func (c *BaseHTTPClient) Close() {
	if c.client != nil && c.client.Transport != nil {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}

func (c *BaseHTTPClient) DoRequest(ctx context.Context, method, url string, reqBody []byte, headerItems ...HeaderItem) (
	httpStatus int, isTimeout bool, respBody []byte, err error) {
	start := time.Now()
	var resp *http.Response
	defer func() {
		log.Infoc(ctx, "Do request done",
			log.Str("method", method),
			log.Str("url", url),
			log.Str("req_body", string(reqBody)),
			log.Int("http_status", httpStatus),
			log.Bool("is_timeout", isTimeout),
			log.Str("resp_headers", fmt.Sprintf("%v", resp.Header)),
			log.Str("resp_body", string(respBody)),
			log.Duration("duration", time.Since(start)),
		)
	}()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Errorc(ctx, "Do request failed to create request", log.Err(err))
		return HTTPStatusUnknown, false, nil, fmt.Errorf("Do request failed to create request, err: %v", err)
	}

	for _, headerItem := range headerItems {
		req.Header.Add(headerItem.Key, headerItem.Value)
	}

	requestID := trace.GetRequestIDFromCtx(ctx)
	if requestID != "" {
		req.Header.Add(trace.HttpHeaderRequestID, requestID)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err = c.client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Errorc(ctx, "Do request timeout", log.Err(err))
			return HTTPStatusUnknown, true, nil, fmt.Errorf("Do request timeout, err: %v", err)
		}
		return HTTPStatusUnknown, false, nil, err
	}
	defer resp.Body.Close()

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Errorc(ctx, "Do request failed to read response body", log.Int("http_status", resp.StatusCode), log.Err(err))
		return resp.StatusCode, false, nil, fmt.Errorf("Do request failed to read response body, http_status: %d, err: %v", resp.StatusCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("HTTP request failed, status: %d, body: %s", resp.StatusCode, string(respBody))
		log.Errorc(ctx, "Do request failed", log.Str("err_msg", errMsg))
		return resp.StatusCode, false, nil, errors.New(errMsg)
	}

	return resp.StatusCode, false, respBody, nil
}
