package purpleair

import (
	"net/http"
	"time"
)

// Client .
type Client struct {
	baseURL    string
	HTTPClient *http.Client
}

// NewClient creates new PurpleApi client
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL: "https://www.purpleair.com/json",
	}
}