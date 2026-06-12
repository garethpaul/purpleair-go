package purpleair

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestSensorWithErrorUsesClientConfiguration(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		assert.Equal(t, "/json", r.URL.Path)
		assert.Equal(t, "17937", r.URL.Query().Get("show"))
		assert.Equal(t, purpleAirUserAgent, r.Header.Get("User-Agent"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"ID":17937,"Label":"Test Sensor"}]}`))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL + "/json"
	client.HTTPClient = server.Client()

	sensor, err := client.SensorWithError("17937")

	assert.NoError(t, err)
	assert.Equal(t, 1, requests)
	assert.Len(t, sensor.Results, 1)
	assert.Equal(t, 17937, sensor.Results[0].ID)
	assert.Equal(t, "Test Sensor", sensor.Results[0].Label)
}

func TestSensorWithContextPropagatesCancellation(t *testing.T) {
	type contextKey string
	const requestMarker contextKey = "request-marker"

	ctx := context.WithValue(context.Background(), requestMarker, "sensor-request")
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	client := NewClient()
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "sensor-request", req.Context().Value(requestMarker))
			return nil, req.Context().Err()
		}),
	}

	sensor, err := client.SensorWithContext(ctx, "17937")

	assert.Nil(t, sensor)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled request error, got %v", err)
	}
}

func TestSensorWithContextRejectsNilContext(t *testing.T) {
	requests := 0
	client := NewClient()
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			return nil, fmt.Errorf("nil context must fail before HTTP requests")
		}),
	}

	sensor, err := client.SensorWithContext(nil, "17937")

	assert.Nil(t, sensor)
	assert.EqualError(t, err, "purpleair: context is required")
	assert.Equal(t, 0, requests, "nil context must fail before HTTP requests")

	sensor, err = client.SensorWithContext(nil, "not-a-sensor-id")

	assert.Nil(t, sensor)
	assert.EqualError(t, err, "purpleair: sensor id must be a positive integer")
	assert.Equal(t, 0, requests, "sensor id validation must remain before nil context validation")
}

func TestSensorWithErrorRejectsBlankSensorIDs(t *testing.T) {
	client := NewClient()

	sensor, err := client.SensorWithError(" \t\n")

	assert.Nil(t, sensor)
	if err == nil || !strings.Contains(err.Error(), "sensor id is required") {
		t.Fatalf("expected sensor id error, got %v", err)
	}
}

func TestSensorWithErrorReturnsEmptyBodyErrors(t *testing.T) {
	client := NewClient()
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       nil,
			}, nil
		}),
	}

	sensor, err := client.SensorWithError("17937")

	assert.Nil(t, sensor)
	if err == nil || !strings.Contains(err.Error(), "response body is empty") {
		t.Fatalf("expected empty response body error, got %v", err)
	}
}

func TestSensorWithErrorRejectsOversizedResponseBodies(t *testing.T) {
	client := NewClient()
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(strings.Repeat(" ", maxSensorResponseBytes+1))),
			}, nil
		}),
	}

	sensor, err := client.SensorWithError("17937")

	assert.Nil(t, sensor)
	if err == nil || !strings.Contains(err.Error(), "response body exceeds") {
		t.Fatalf("expected oversized response body error, got %v", err)
	}
}

func TestSensorWithErrorReturnsNilResponseErrors(t *testing.T) {
	client := NewClient()
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, nil
		}),
	}

	sensor, err := client.SensorWithError("17937")

	assert.Nil(t, sensor)
	if err == nil || !strings.Contains(err.Error(), "request failed") {
		t.Fatalf("expected request failure error, got %v", err)
	}
}

func TestSensorWithErrorReturnsStatusErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL + "/json"
	client.HTTPClient = server.Client()

	_, err := client.SensorWithError("17937")

	if err == nil || !strings.Contains(err.Error(), "unexpected status 503") {
		t.Fatalf("expected unexpected status error, got %v", err)
	}
}

func TestSensorWithErrorReturnsMalformedJSONErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[`))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL + "/json"
	client.HTTPClient = server.Client()

	sensor, err := client.SensorWithError("17937")

	assert.Error(t, err)
	assert.Nil(t, sensor)
}

func TestSensorWithErrorReturnsEmptyResultErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL + "/json"
	client.HTTPClient = server.Client()

	sensor, err := client.SensorWithError("17937")

	assert.Nil(t, sensor)
	if err == nil || !strings.Contains(err.Error(), `no results for sensor "17937"`) {
		t.Fatalf("expected no-results error, got %v", err)
	}
}

func TestSensorWithErrorRejectsInvalidResultIDs(t *testing.T) {
	for name, testCase := range map[string]struct {
		payload       string
		expectedIndex int
	}{
		"missing":      {`{"results":[{"Label":"Missing ID"}]}`, 0},
		"zero":         {`{"results":[{"ID":0,"Label":"Zero ID"}]}`, 0},
		"negative":     {`{"results":[{"ID":-1,"Label":"Negative ID"}]}`, 0},
		"later result": {`{"results":[{"ID":17937},{"ID":0}]}`, 1},
	} {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(testCase.payload))
			}))
			defer server.Close()

			client := NewClient()
			client.baseURL = server.URL + "/json"
			client.HTTPClient = server.Client()

			sensor, err := client.SensorWithError("17937")

			assert.Nil(t, sensor)
			expectedError := fmt.Sprintf("result %d has invalid sensor id", testCase.expectedIndex)
			if err == nil || !strings.Contains(err.Error(), expectedError) {
				t.Fatalf("expected invalid result sensor id error, got %v", err)
			}
		})
	}
}

func TestSensorWithErrorAcceptsMultipleValidResultIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"ID":17938},{"ID":17937}]}`))
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL + "/json"
	client.HTTPClient = server.Client()

	sensor, err := client.SensorWithError("17937")

	assert.NoError(t, err)
	assert.Len(t, sensor.Results, 2)
	assert.Equal(t, 17937, sensor.Results[1].ID)
}

func TestSensorWithErrorRejectsInvalidRequestedSensorIDs(t *testing.T) {
	client := NewClient()
	requests := 0
	client.HTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests++
			return nil, errors.New("unexpected request")
		}),
	}

	for _, sensorID := range []string{"0", "-1", "1.5", "sensor"} {
		sensor, err := client.SensorWithError(sensorID)

		assert.Nil(t, sensor)
		assert.EqualError(t, err, "purpleair: sensor id must be a positive integer")
	}

	assert.Equal(t, 0, requests, "invalid sensor IDs must fail before HTTP requests")
}

func TestSensorWithErrorRejectsMismatchedResponseSensorIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"ID":17938},{"ID":17939}]}`))
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.URL + "/json")
	client.HTTPClient = server.Client()

	sensor, err := client.SensorWithError("17937")

	assert.Nil(t, sensor)
	assert.EqualError(t, err, "purpleair: response does not include requested sensor 17937")
}
