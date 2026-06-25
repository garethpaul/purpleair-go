package purpleair

import (
	"net/http"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	// assert equality
	client := NewClient()

	assert.Equal(t, "https://www.purpleair.com/json", client.baseURL, "error with client")
	assert.Equal(t, 30*time.Second, client.HTTPClient.Timeout, "error with client timeout")
}

func TestNewClientWithBaseURL(t *testing.T) {
	client := NewClientWithBaseURL(" https://example.test/json?api_key=local ")

	assert.Equal(t, "https://example.test/json?api_key=local", client.baseURL)
	assert.Equal(t, 30*time.Second, client.HTTPClient.Timeout)
	assert.Equal(t, "https://example.test/json?api_key=local&show=17937", client.sensorURL("17937"))
}

func TestNewClientWithBaseURLFallsBackForBlankValues(t *testing.T) {
	client := NewClientWithBaseURL(" \t\n")

	assert.Equal(t, defaultBaseURL, client.baseURL)
	assert.Equal(t, defaultHTTPTimeout, client.HTTPClient.Timeout)
}

func TestNewClientWithBaseURLFallsBackForInvalidValues(t *testing.T) {
	for _, baseURL := range []string{"://bad-url", "ftp://example.test/json", "https:///missing-host", "https://user:pass@example.test/json", "https://example.test/json#local-token"} {
		t.Run(baseURL, func(t *testing.T) {
			client := NewClientWithBaseURL(baseURL)

			assert.Equal(t, defaultBaseURL, client.baseURL)
			assert.Equal(t, "https://www.purpleair.com/json?show=17937", client.sensorURL("17937"))
		})
	}
}

func TestZeroValueClientUsesDefaultTimeout(t *testing.T) {
	client := Client{}

	assert.Equal(t, 30*time.Second, client.httpClient().Timeout)
	assert.Equal(t, "https://www.purpleair.com/json?show=17937", client.sensorURL("17937"))
}

func TestNilClientUsesDefaultTimeout(t *testing.T) {
	var client *Client

	assert.Equal(t, 30*time.Second, client.httpClient().Timeout)
	assert.Equal(t, "https://www.purpleair.com/json?show=17937", client.sensorURL("17937"))
}

func TestClientPreservesCallerProvidedHTTPTimeout(t *testing.T) {
	redirectPolicy := func(*http.Request, []*http.Request) error {
		return nil
	}
	customHTTPClient := &http.Client{
		Timeout:       2 * time.Minute,
		CheckRedirect: redirectPolicy,
	}
	client := NewClient()
	client.HTTPClient = customHTTPClient

	assert.Same(t, customHTTPClient, client.httpClient())
	assert.Equal(t, 2*time.Minute, client.httpClient().Timeout)
	assert.NotNil(t, client.httpClient().CheckRedirect)
}

func TestDefaultHTTPClientRejectsRedirects(t *testing.T) {
	var nilClient *Client
	clients := map[string]*http.Client{
		"constructor": NewClient().HTTPClient,
		"zero value":  (&Client{}).httpClient(),
		"nil client":  nilClient.httpClient(),
	}

	for name, httpClient := range clients {
		t.Run(name, func(t *testing.T) {
			assert.NotNil(t, httpClient.CheckRedirect)
			assert.Equal(t, http.ErrUseLastResponse, httpClient.CheckRedirect(nil, nil))
		})
	}
}
