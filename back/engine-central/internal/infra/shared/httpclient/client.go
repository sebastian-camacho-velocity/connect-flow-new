package httpclient

import (
	"net"
	"net/http"
	"time"
)

type HTTPClientConfig struct {
	Timeout            time.Duration
	MaxIdleConns       int
	IdleConnTimeout    time.Duration
	DisableCompression bool
}

// NewHTTPClient recibe una configuración personalizada por cada cliente externo
func NewHTTPClient(cfg HTTPClientConfig) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:       cfg.MaxIdleConns,
		IdleConnTimeout:    cfg.IdleConnTimeout,
		DisableCompression: cfg.DisableCompression,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}
}
