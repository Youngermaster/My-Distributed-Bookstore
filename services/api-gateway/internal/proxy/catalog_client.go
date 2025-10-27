package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CatalogClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCatalogClient(baseURL string, timeout time.Duration) *CatalogClient {
	return &CatalogClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ProxyRequest forwards a request to the catalog service
func (c *CatalogClient) ProxyRequest(method, path string, body []byte, headers map[string]string) ([]byte, int, error) {
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
		return nil, http.StatusServiceUnavailable, fmt.Errorf("catalog service unavailable: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	return responseBody, resp.StatusCode, nil
}

// HealthCheck checks if the catalog service is healthy
func (c *CatalogClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("catalog service is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("catalog service returned status %d", resp.StatusCode)
	}

	return nil
}

// Helper methods for specific operations

func (c *CatalogClient) GetBooks(queryParams string) (interface{}, error) {
	path := "/api/v1/books"
	if queryParams != "" {
		path += "?" + queryParams
	}

	body, statusCode, err := c.ProxyRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

func (c *CatalogClient) GetBookByID(id string) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/books/%s", id)

	body, statusCode, err := c.ProxyRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

func (c *CatalogClient) SearchBooks(query string) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/books/search?q=%s", query)

	body, statusCode, err := c.ProxyRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}
