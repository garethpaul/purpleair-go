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
