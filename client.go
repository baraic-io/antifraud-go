package antifraud

import (
	"crypto/tls"
	"net/http"
	"time"
)

type ClientConfig struct {
	Host   string
	APIKey string
}

type Client struct {
	httpClient *http.Client

	ClientConfig
}

func NewClient(c ClientConfig) (Client, error) {
	client := Client{ClientConfig: c}

	/* TODO: Add conf related TLS InsecureSkipVerify */
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
