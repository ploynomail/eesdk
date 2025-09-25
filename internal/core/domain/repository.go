package domain

import "context"

// HTTPClient defines the interface for HTTP operations
type HTTPClient interface {
	Do(ctx context.Context, req *Request) (*Response, error)
	SetConfig(config *ServiceConfig)
	GetConfig() *ServiceConfig
}

// ServiceRepository defines the interface for service operations
type ServiceRepository interface {
	Execute(ctx context.Context, req *Request) (*Response, error)
	Health(ctx context.Context) error
}
