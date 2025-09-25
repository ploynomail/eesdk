package service

import (
	"context"
	"net/http"

	"ee-sdk/internal/core/domain"
)

type HTTPService struct {
	client *http.Client
	config *domain.ServiceConfig
}

func NewHTTPService(config *domain.ServiceConfig) *HTTPService {
	return &HTTPService{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}
}

func (s *HTTPService) Do(ctx context.Context, req *domain.Request) (*domain.Response, error) {
	// Implementation of HTTP request logic
	// ...existing implementation...
	return nil, nil
}

func (s *HTTPService) SetConfig(config *domain.ServiceConfig) {
	s.config = config
	s.client.Timeout = config.Timeout
}

func (s *HTTPService) GetConfig() *domain.ServiceConfig {
	return s.config
}
