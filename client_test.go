package purpleair

import (
	"testing"
	assert "github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	// assert equality
	client := NewClient()

	assert.Equal(t, client.baseURL, "https://www.purpleair.com/json", "error with client")
}