package purpleair

import (
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL     = "https://www.purpleair.com/json"
	defaultHTTPTimeout = 5 * time.Minute
)

// Client .
type Client struct {
	baseURL    string
	HTTPClient *http.Client
}

// NewClient creates new PurpleApi client
func NewClient() *Client {
	return &Client{
		HTTPClient: defaultHTTPClient(),
		baseURL:    defaultBaseURL,
	}
}

// NewClientWithBaseURL creates a PurpleAir client for a custom endpoint.
func NewClientWithBaseURL(baseURL string) *Client {
	client := NewClient()
	baseURL = strings.TrimSpace(baseURL)
	if baseURL != "" {
		client.baseURL = baseURL
	}

	return client
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultHTTPTimeout,
	}
}

func (c *Client) httpClient() *http.Client {
	if c != nil && c.HTTPClient != nil {
		return c.HTTPClient
	}

	return defaultHTTPClient()
}
