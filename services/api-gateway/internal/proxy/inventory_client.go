package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// InventoryClient proxies HTTP requests to the inventory service.
type InventoryClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewInventoryClient constructs an InventoryClient with the provided base URL and timeout.
func NewInventoryClient(baseURL string, timeout time.Duration) *InventoryClient {
	return &InventoryClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// ProxyRequest forwards an arbitrary request to the inventory service.
func (c *InventoryClient) ProxyRequest(method, path string, body []byte, headers map[string]string) ([]byte, int, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("inventory service unavailable: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	return responseBody, resp.StatusCode, nil
}

// HealthCheck verifies the health endpoint of the inventory service.
func (c *InventoryClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("inventory service is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inventory service returned status %d", resp.StatusCode)
	}

	return nil
}
