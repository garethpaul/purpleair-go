package purpleair

import (
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

	sensor, err := client.SensorWithError("missing")

	assert.Nil(t, sensor)
	if err == nil || !strings.Contains(err.Error(), `no results for sensor "missing"`) {
		t.Fatalf("expected no-results error, got %v", err)
	}
}
