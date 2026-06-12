package purpleair

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL     = "https://www.purpleair.com/json"
	defaultHTTPTimeout = 30 * time.Second
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
	if isSupportedBaseURL(baseURL) {
		client.baseURL = baseURL
	}

	return client
}

func isSupportedBaseURL(baseURL string) bool {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	return (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != "" && parsed.User == nil && parsed.Fragment == ""
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
