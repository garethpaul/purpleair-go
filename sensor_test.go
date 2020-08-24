package purpleair

import (
	"testing"
	assert "github.com/stretchr/testify/require"
)

func TestResults(t *testing.T) {
	// assert equality
	client := NewClient()
	s:= client.Sensor("17937")
	ResultsLength:= len(s.Results)
	assert.GreaterOrEqualf(t, ResultsLength, 1, "error message %s", "formatted")
}