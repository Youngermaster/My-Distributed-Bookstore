package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RecommendationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRecommendationClient(baseURL string, timeout time.Duration) *RecommendationClient {
	return &RecommendationClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ProxyRequest forwards a request to the recommendation service
func (c *RecommendationClient) ProxyRequest(method, path string, body []byte, headers map[string]string) ([]byte, int, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	// Copy headers (including X-User-Id for authentication)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set Content-Type if not already set and body exists
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("recommendation service unavailable: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response: %w", err)
	}

	return responseBody, resp.StatusCode, nil
}

// HealthCheck checks if the recommendation service is healthy
func (c *RecommendationClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/health")
	if err != nil {
		return fmt.Errorf("recommendation service is not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("recommendation service returned status %d", resp.StatusCode)
	}

	return nil
}

// Helper methods for specific operations

func (c *RecommendationClient) GetRecommendations(userID string, limit int) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/recommendations/me?limit=%d", limit)

	headers := map[string]string{
		"X-User-Id": userID,
	}

	body, statusCode, err := c.ProxyRequest("GET", path, nil, headers)
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

func (c *RecommendationClient) GetSimilarBooks(bookID string, limit int) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/recommendations/similar/%s?limit=%d", bookID, limit)

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

func (c *RecommendationClient) GetTrendingBooks(limit, days int) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/recommendations/trending?limit=%d&days=%d", limit, days)

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
