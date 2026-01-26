package antifraud

import (
	"crypto/tls"
	"net/http"
	"time"
)

const DefaultValidationCtxDeadlineTimeout = 30

type ClientConfig struct {
	Host   string
	APIKey string

	// Transaction validation context max deadline timeout in seconds.
	// If not provided default=30 value will be setted.
	ValidationCtxDeadlineTimeout int
}

type Client struct {
	httpClient *http.Client

	ClientConfig
}

func NewClient(c ClientConfig) (Client, error) {
	if c.ValidationCtxDeadlineTimeout == 0 {
		c.ValidationCtxDeadlineTimeout = DefaultValidationCtxDeadlineTimeout
	}

	client := Client{ClientConfig: c}

	/* TODO: Add conf related TLS InsecureSkipVerify, move timeouts to conf */
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client.httpClient = httpClient
	return client, nil
}
