// Package hostex provides a Go client library for the Hostex API v3.0.0
package hostex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the default Hostex API base URL
	DefaultBaseURL = "https://api.hostex.io/v3"

	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second

	// UserAgent is the user agent string sent with requests
	UserAgent = "hostex-go/1.0.0"
)

// Client is the Hostex API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// Config holds client configuration options
type Config struct {
	// AccessToken is your Hostex API access token (required)
	AccessToken string

	// BaseURL is the Hostex API base URL (optional, defaults to DefaultBaseURL)
	BaseURL string

	// HTTPClient is a custom HTTP client (optional)
	HTTPClient *http.Client

	// Timeout is the HTTP request timeout (optional, defaults to DefaultTimeout)
	Timeout time.Duration
}

// NewClient creates a new Hostex API client
func NewClient(config Config) (*Client, error) {
	if config.AccessToken == "" {
		return nil, fmt.Errorf("access token is required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		timeout := config.Timeout
		if timeout == 0 {
			timeout = DefaultTimeout
		}
		httpClient = &http.Client{
			Timeout: timeout,
		}
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		token:      config.AccessToken,
	}, nil
}

// APIResponse represents a standard Hostex API response
type APIResponse struct {
	RequestID string      `json:"request_id"`
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Data      interface{} `json:"data,omitempty"`
}

// doRequest executes an HTTP request to the Hostex API
func (c *Client) doRequest(ctx context.Context, method, endpoint string, params url.Values, body interface{}) (*APIResponse, error) {
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %w", err)
	}

	// Add query parameters
	if params != nil {
		u.RawQuery = params.Encode()
	}

	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Hostex-Access-Token", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResp.ErrorCode != 200 {
		return &apiResp, fmt.Errorf("API error %d: %s", apiResp.ErrorCode, apiResp.ErrorMsg)
	}

	return &apiResp, nil
}

// Helper function to convert interface{} to specific type
func unmarshalData(data interface{}, v interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, v); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}
