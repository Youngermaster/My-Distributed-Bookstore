package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OrderClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewOrderClient(baseURL string, timeout time.Duration) *OrderClient {
	return &OrderClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ProxyRequest forwards a request to the order service
func (c *OrderClient) ProxyRequest(method, path string, body []byte, headers map[string]string) ([]byte, int, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	// Copy headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set Content-Type if not already set and body exists
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("order service unavailable: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	return responseBody, resp.StatusCode, nil
}

// HealthCheck checks if the order service is healthy
func (c *OrderClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("order service is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("order service returned status %d", resp.StatusCode)
	}

	return nil
}
