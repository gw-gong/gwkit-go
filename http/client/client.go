package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

// BaseHTTPClient is a base HTTP client that is used for inheritance
type BaseHTTPClient struct {
	client *http.Client
}

func (c *BaseHTTPClient) Close() {
	if c.client != nil && c.client.Transport != nil {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}

// If you don't know how to pass parameters, you can refer to the following content:
// - timeout: 5
// - maxIdleConnsPerHost: 10
// - maxConnsPerHost: 100
// - tlsConfig: nil
func NewHTTPClient(timeout time.Duration, maxIdleConnsPerHost, maxConnsPerHost int, tlsConfig *tls.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnsPerHost,
			MaxConnsPerHost:     maxConnsPerHost,
			TLSClientConfig:     tlsConfig,
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   3 * time.Second,
			ResponseHeaderTimeout: timeout,
		},
	}
}

type SOCKS5Config struct {
	Address  string
	User     string
	Password string
}

func NewHTTPClientWithSOCKS5(timeout time.Duration, maxIdleConnsPerHost, maxConnsPerHost int, tlsConfig *tls.Config, socks5Config *SOCKS5Config) *http.Client {
	var dialer proxy.Dialer
	if socks5Config == nil || socks5Config.Address == "" {
		dialer = &net.Dialer{Timeout: timeout, KeepAlive: 60 * time.Second}
	} else {
		dialer, _ = proxy.SOCKS5(
			"tcp",
			socks5Config.Address,
			&proxy.Auth{
				User:     socks5Config.User,
				Password: socks5Config.Password,
			},
			&net.Dialer{Timeout: timeout, KeepAlive: 60 * time.Second},
		)
	}

	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   maxIdleConnsPerHost,
			MaxConnsPerHost:       maxConnsPerHost,
			TLSClientConfig:       tlsConfig,
			Dial:                  dialer.Dial,
			TLSHandshakeTimeout:   3 * time.Second,
			ResponseHeaderTimeout: timeout,
		},
	}
}
