package purpleair

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	// assert equality
	client := NewClient()

	assert.Equal(t, "https://www.purpleair.com/json", client.baseURL, "error with client")
	assert.Equal(t, 5*time.Minute, client.HTTPClient.Timeout, "error with client timeout")
}

func TestNewClientWithBaseURL(t *testing.T) {
	client := NewClientWithBaseURL(" https://example.test/json?api_key=local ")

	assert.Equal(t, "https://example.test/json?api_key=local", client.baseURL)
	assert.Equal(t, 5*time.Minute, client.HTTPClient.Timeout)
	assert.Equal(t, "https://example.test/json?api_key=local&show=17937", client.sensorURL("17937"))
}

func TestNewClientWithBaseURLFallsBackForBlankValues(t *testing.T) {
	client := NewClientWithBaseURL(" \t\n")

	assert.Equal(t, defaultBaseURL, client.baseURL)
	assert.Equal(t, defaultHTTPTimeout, client.HTTPClient.Timeout)
}

func TestZeroValueClientUsesDefaultTimeout(t *testing.T) {
	client := Client{}

	assert.Equal(t, 5*time.Minute, client.httpClient().Timeout)
	assert.Equal(t, "https://www.purpleair.com/json?show=17937", client.sensorURL("17937"))
}

func TestNilClientUsesDefaultTimeout(t *testing.T) {
	var client *Client

	assert.Equal(t, 5*time.Minute, client.httpClient().Timeout)
	assert.Equal(t, "https://www.purpleair.com/json?show=17937", client.sensorURL("17937"))
}
