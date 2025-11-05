package grpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ServiceClients holds all service client connections
type ServiceClients struct {
	UserService          *ServiceClient
	CatalogService       *ServiceClient
	OrderService         *ServiceClient
	CartService          *ServiceClient
	InventoryService     *ServiceClient
	RecommendationService *ServiceClient
}

// ServiceClient represents a generic HTTP client for microservices
type ServiceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewServiceClient creates a new service client
func NewServiceClient(baseURL string) *ServiceClient {
	return &ServiceClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get performs a GET request
func (c *ServiceClient) Get(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Post performs a POST request
func (c *ServiceClient) Post(path string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// NewServiceClients initializes all service clients
func NewServiceClients(
	userServiceURL,
	catalogServiceURL,
	orderServiceURL,
	cartServiceURL,
	inventoryServiceURL,
	recommendationServiceURL string,
) *ServiceClients {
	return &ServiceClients{
		UserService:          NewServiceClient(fmt.Sprintf("http://%s", userServiceURL)),
		CatalogService:       NewServiceClient(fmt.Sprintf("http://%s", catalogServiceURL)),
		OrderService:         NewServiceClient(fmt.Sprintf("http://%s", orderServiceURL)),
		CartService:          NewServiceClient(fmt.Sprintf("http://%s", cartServiceURL)),
		InventoryService:     NewServiceClient(fmt.Sprintf("http://%s", inventoryServiceURL)),
		RecommendationService: NewServiceClient(fmt.Sprintf("http://%s", recommendationServiceURL)),
	}
}
