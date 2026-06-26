package purpleair

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestNewDataAPIClientRequiresReadAPIKey(t *testing.T) {
	for _, key := range []string{"", " ", "\t\n"} {
		client, err := NewDataAPIClient(key)

		assert.Nil(t, client)
		assert.EqualError(t, err, "purpleair: API read key is required")
	}
}

func TestNewDataAPIClientUsesFixedDefaults(t *testing.T) {
	client, err := NewDataAPIClient("  api-read-key  ")

	assert.NoError(t, err)
	assert.Equal(t, defaultDataAPIBaseURL, client.baseURL)
	assert.Equal(t, "api-read-key", client.readAPIKey)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, 30*time.Second, client.HTTPClient.Timeout)
	assert.NotNil(t, client.HTTPClient.CheckRedirect)
	assert.Equal(t, http.ErrUseLastResponse, client.HTTPClient.CheckRedirect(nil, nil))
}

func TestDataAPIClientPreservesCallerHTTPClient(t *testing.T) {
	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)

	custom := &http.Client{Timeout: 7 * time.Second}
	client.HTTPClient = custom

	assert.Same(t, custom, client.httpClient())
}

func TestDataAPISensorRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "/v1/sensors/17937", req.URL.Path)
		assert.Equal(t, "name,last_seen,latitude,longitude,pm2.5_atm", req.URL.Query().Get("fields"))
		assert.Empty(t, req.URL.Query().Get("read_key"))
		assert.Equal(t, "api-read-key", req.Header.Get("X-API-Key"))
		assert.NotContains(t, req.URL.String(), "api-read-key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"api_version":"V1.0.15","time_stamp":100,"data_time_stamp":99,"sensor":{"sensor_index":17937,"name":"Test","last_seen":98,"latitude":0,"longitude":0,"pm2.5_atm":0}}`))
	}))
	defer server.Close()

	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.baseURL = server.URL + "/v1"
	client.HTTPClient = server.Client()

	response, err := client.SensorData(context.Background(), 17937, SensorDataOptions{})

	assert.NoError(t, err)
	assert.Equal(t, 17937, response.Sensor.SensorIndex)
	assert.Equal(t, "Test", *response.Sensor.Name)
}

func TestDataAPISensorPropagatesContext(t *testing.T) {
	type contextKey string
	const marker contextKey = "marker"
	ctx := context.WithValue(context.Background(), marker, "request-context")

	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "request-context", req.Context().Value(marker))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: testReadCloser{Reader: strings.NewReader(
				`{"api_version":"V1.0.15","time_stamp":100,"data_time_stamp":99,"sensor":{"sensor_index":17937}}`,
			)},
		}, nil
	})}

	response, err := client.SensorData(ctx, 17937, SensorDataOptions{})

	assert.NoError(t, err)
	assert.Equal(t, 17937, response.Sensor.SensorIndex)
}

func TestDataAPIPrivateSensorRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "private-read-key", req.URL.Query().Get("read_key"))
		assert.NotContains(t, req.URL.RequestURI(), "api-read-key")
		_, _ = w.Write([]byte(`{"api_version":"V1.0.15","time_stamp":100,"data_time_stamp":99,"sensor":{"sensor_index":42}}`))
	}))
	defer server.Close()

	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.baseURL = server.URL
	client.HTTPClient = server.Client()

	response, err := client.SensorData(context.Background(), 42, SensorDataOptions{SensorReadKey: "  private-read-key  "})

	assert.NoError(t, err)
	assert.Equal(t, 42, response.Sensor.SensorIndex)
}

func TestDataAPISensorRejectsInvalidInputsBeforeIO(t *testing.T) {
	requests := 0
	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requests++
		return nil, errors.New("unexpected request")
	})}

	for _, test := range []struct {
		ctx         context.Context
		sensorIndex int
		message     string
	}{
		{nil, 1, "purpleair: context is required"},
		{context.Background(), 0, "purpleair: sensor index must be positive"},
		{context.Background(), -1, "purpleair: sensor index must be positive"},
	} {
		response, requestErr := client.SensorData(test.ctx, test.sensorIndex, SensorDataOptions{})
		assert.Nil(t, response)
		assert.EqualError(t, requestErr, test.message)
	}
	assert.Equal(t, 0, requests)
}

func TestDataAPISensorValidatesResponse(t *testing.T) {
	latitude := 37.7
	longitude := -122.4
	pm25 := 5.5
	lastSeen := int64(98)
	name := "Sensor"

	valid := SensorDataResponse{
		APIVersion:    "V1.0.15",
		Timestamp:     100,
		DataTimestamp: 99,
		Sensor: SensorData{
			SensorIndex: 17937,
			Name:        &name,
			LastSeen:    &lastSeen,
			Latitude:    &latitude,
			Longitude:   &longitude,
			PM25ATM:     &pm25,
		},
	}

	tests := []struct {
		name    string
		mutate  func(*SensorDataResponse)
		message string
	}{
		{"sensor identity", func(response *SensorDataResponse) { response.Sensor.SensorIndex = 2 }, "purpleair: response sensor index does not match request"},
		{"API timestamp", func(response *SensorDataResponse) { response.Timestamp = -1 }, "purpleair: response timestamps must not be negative"},
		{"data timestamp", func(response *SensorDataResponse) { response.DataTimestamp = -1 }, "purpleair: response timestamps must not be negative"},
		{"last seen", func(response *SensorDataResponse) { value := int64(-1); response.Sensor.LastSeen = &value }, "purpleair: response timestamps must not be negative"},
		{"partial latitude", func(response *SensorDataResponse) { response.Sensor.Longitude = nil }, "purpleair: response coordinates must both be present or absent"},
		{"partial longitude", func(response *SensorDataResponse) { response.Sensor.Latitude = nil }, "purpleair: response coordinates must both be present or absent"},
		{"latitude range", func(response *SensorDataResponse) { value := 91.0; response.Sensor.Latitude = &value }, "purpleair: response latitude is invalid"},
		{"longitude range", func(response *SensorDataResponse) { value := -181.0; response.Sensor.Longitude = &value }, "purpleair: response longitude is invalid"},
		{"PM2.5 finite", func(response *SensorDataResponse) { value := math.NaN(); response.Sensor.PM25ATM = &value }, "purpleair: response PM2.5 is invalid"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := valid
			test.mutate(&response)
			err := validateDataAPIResponse(&response, 17937)
			assert.EqualError(t, err, test.message)
		})
	}
}

func TestDataAPISensorRejectsMalformedTrailingAndOversizedResponses(t *testing.T) {
	valid := `{"api_version":"V1.0.15","time_stamp":100,"data_time_stamp":99,"sensor":{"sensor_index":17937}}`
	for _, test := range []struct {
		name    string
		body    string
		message string
	}{
		{"empty", "", "purpleair: response body is empty"},
		{"malformed", "{", "purpleair: decode response body"},
		{"trailing", valid + `{}`, "purpleair: response body contains trailing data"},
		{"oversized", strings.Repeat("x", maxDataAPIResponseBytes+1), "purpleair: response body exceeds 1048576 bytes"},
	} {
		t.Run(test.name, func(t *testing.T) {
			client := dataAPIClientWithBody(test.body)
			response, err := client.SensorData(context.Background(), 17937, SensorDataOptions{})
			assert.Nil(t, response)
			assert.Contains(t, err.Error(), test.message)
		})
	}
}

func TestDataAPISensorRejectsOversizedDeclaredBodyBeforeRead(t *testing.T) {
	body := &trackingReadCloser{Reader: strings.NewReader("{}")}
	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:    http.StatusOK,
			ContentLength: maxDataAPIResponseBytes + 1,
			Body:          body,
		}, nil
	})}

	response, requestErr := client.SensorData(context.Background(), 17937, SensorDataOptions{})

	assert.Nil(t, response)
	assert.EqualError(t, requestErr, "purpleair: response body exceeds 1048576 bytes")
	assert.True(t, body.closed)
}

func TestDataAPISensorWrapsResponseReadErrors(t *testing.T) {
	readErr := errors.New("read failed")
	body := &errorReadCloser{Reader: failingReader{err: readErr}}
	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
	})}

	response, requestErr := client.SensorData(context.Background(), 17937, SensorDataOptions{})

	assert.Nil(t, response)
	assert.True(t, errors.Is(requestErr, readErr))
	assert.Contains(t, requestErr.Error(), "purpleair: read response body")
	assert.True(t, body.closed)
}

func TestDataAPISensorReturnsDetailSafeStatusErrors(t *testing.T) {
	for _, status := range []int{401, 403, 404, 429, 500, 503} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			client, err := NewDataAPIClient("api-read-key")
			assert.NoError(t, err)
			client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: status,
					Body:       testReadCloser{Reader: strings.NewReader("provider-secret-body")},
				}, nil
			})}

			response, requestErr := client.SensorData(context.Background(), 17937, SensorDataOptions{SensorReadKey: "private-read-key"})

			assert.Nil(t, response)
			assert.Contains(t, requestErr.Error(), fmt.Sprint(status))
			assert.NotContains(t, requestErr.Error(), "api-read-key")
			assert.NotContains(t, requestErr.Error(), "private-read-key")
			assert.NotContains(t, requestErr.Error(), "provider-secret-body")
		})
	}
}

func TestDataAPISensorClosesStatusResponseBodies(t *testing.T) {
	body := &trackingReadCloser{Reader: strings.NewReader("provider details")}
	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusForbidden, Body: body}, nil
	})}

	response, requestErr := client.SensorData(context.Background(), 17937, SensorDataOptions{})

	assert.Nil(t, response)
	assert.EqualError(t, requestErr, "purpleair: data API forbidden (status 403)")
	assert.True(t, body.closed)
}

func TestDataAPISensorClosesBodiesAndPreservesPrimaryFailure(t *testing.T) {
	body := &errorReadCloser{Reader: strings.NewReader("{"), closeErr: errors.New("close failed")}
	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
	})}

	response, requestErr := client.SensorData(context.Background(), 17937, SensorDataOptions{})

	assert.Nil(t, response)
	assert.True(t, body.closed)
	assert.Contains(t, requestErr.Error(), "decode response body")
	assert.NotContains(t, requestErr.Error(), "close failed")
}

func TestDataAPISensorReturnsCloseFailureAfterSuccessfulDecode(t *testing.T) {
	body := &errorReadCloser{
		Reader:   strings.NewReader(`{"api_version":"V1.0.15","time_stamp":100,"data_time_stamp":99,"sensor":{"sensor_index":17937}}`),
		closeErr: errors.New("close failed"),
	}
	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
	})}

	response, requestErr := client.SensorData(context.Background(), 17937, SensorDataOptions{})

	assert.Nil(t, response)
	assert.EqualError(t, requestErr, "purpleair: close response body: close failed")
	assert.True(t, body.closed)
}

func TestDataAPISensorPreservesRequestCauseWithoutCredentialDetails(t *testing.T) {
	requestErr := context.DeadlineExceeded
	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, requestErr
	})}

	response, err := client.SensorData(context.Background(), 17937, SensorDataOptions{SensorReadKey: "private-read-key"})

	assert.Nil(t, response)
	assert.True(t, errors.Is(err, requestErr))
	assert.Equal(t, "purpleair: data API request failed", err.Error())
	assert.NotContains(t, errors.Unwrap(err).Error(), "api-read-key")
	assert.NotContains(t, errors.Unwrap(err).Error(), "private-read-key")
	var urlErr *url.Error
	assert.False(t, errors.As(err, &urlErr))
}

func TestDataAPISensorRejectsRedirects(t *testing.T) {
	destinationRequests := 0
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		destinationRequests++
	}))
	defer destination.Close()

	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, destination.URL, http.StatusFound)
	}))
	defer source.Close()

	client, err := NewDataAPIClient("api-read-key")
	assert.NoError(t, err)
	client.baseURL = source.URL

	response, requestErr := client.SensorData(context.Background(), 17937, SensorDataOptions{})

	assert.Nil(t, response)
	assert.EqualError(t, requestErr, "purpleair: data API unexpected status 302")
	assert.Equal(t, 0, destinationRequests)
}

type testReadCloser struct {
	io.Reader
}

func (testReadCloser) Close() error {
	return nil
}

func dataAPIClientWithBody(body string) *DataAPIClient {
	client, _ := NewDataAPIClient("api-read-key")
	client.HTTPClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       testReadCloser{Reader: strings.NewReader(body)},
		}, nil
	})}
	return client
}
