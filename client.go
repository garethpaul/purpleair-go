package purpleair

import (
	"net/http"
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
