package services

import (
	"net/http"
	"time"
)

// HTTPClient interface para cliente HTTP
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
	client *http.Client
}

// NewHTTPClient cria uma nova instância do cliente HTTP
func NewHTTPClient() HTTPClient {
	return &httpClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Do executa a requisição HTTP
func (h *httpClient) Do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}
