package purpleair

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	// assert equality
	client := NewClient()

	assert.Equal(t, "https://www.purpleair.com/json", client.baseURL, "error with client")
}
