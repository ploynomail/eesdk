package domain

import "time"

// Request represents a generic API request
type Request struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    interface{}       `json:"body,omitempty"`
	Timeout time.Duration     `json:"timeout,omitempty"`
}

// Response represents a generic API response
type Response struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
	Error      error             `json:"error,omitempty"`
}

// ServiceConfig holds service configuration
type ServiceConfig struct {
	BaseURL    string        `json:"base_url"`
	APIKey     string        `json:"api_key"`
	Timeout    time.Duration `json:"timeout"`
	RetryCount int           `json:"retry_count"`
}
