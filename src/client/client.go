package client

import (
	"context"
	"ee-sdk/internal/core/domain"
	httpClient "ee-sdk/internal/infrastructure/http"
	"time"
)

// Client provides Go interface for the SDK
type Client struct {
	service domain.ServiceRepository
}

// Config holds client configuration
type Config struct {
	BaseURL    string
	APIKey     string
	Timeout    time.Duration
	RetryCount int
}

// NewClient creates a new SDK client
func NewClient(config *Config) *Client {
	serviceConfig := &domain.ServiceConfig{
		BaseURL:    config.BaseURL,
		APIKey:     config.APIKey,
		Timeout:    config.Timeout,
		RetryCount: config.RetryCount,
	}

	httpSvc := httpClient.NewClient(serviceConfig)

	return &Client{
		service: httpSvc,
	}
}

// ExecuteRequest executes a generic request
func (c *Client) ExecuteRequest(ctx context.Context, method, path string, body interface{}) (*domain.Response, error) {
	req := &domain.Request{
		Method: method,
		Path:   path,
		Body:   body,
	}

	return c.service.Execute(ctx, req)
}

// Health checks service health
func (c *Client) Health(ctx context.Context) error {
	return c.service.Health(ctx)
}
