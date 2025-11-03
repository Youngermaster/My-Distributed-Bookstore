package proxy

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type AdminClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAdminClient(baseURL string, timeout time.Duration) *AdminClient {
	return &AdminClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *AdminClient) ProxyRequest(method, path string, body []byte, headers map[string]string) ([]byte, int, error) {
	url := c.baseURL + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Copy headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return responseBody, resp.StatusCode, nil
}

func (c *AdminClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
