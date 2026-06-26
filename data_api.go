package purpleair

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultDataAPIBaseURL   = "https://api.purpleair.com/v1"
	dataAPISensorFields     = "name,last_seen,latitude,longitude,pm2.5_atm"
	maxDataAPIResponseBytes = 1 << 20
)

// DataAPIClient accesses PurpleAir's authenticated Data API.
type DataAPIClient struct {
	HTTPClient *http.Client
	baseURL    string
	readAPIKey string
}

// SensorDataOptions contains optional credentials for one sensor request.
type SensorDataOptions struct {
	SensorReadKey string
}

// SensorDataResponse is the authenticated single-sensor response envelope.
type SensorDataResponse struct {
	APIVersion    string     `json:"api_version"`
	Timestamp     int64      `json:"time_stamp"`
	DataTimestamp int64      `json:"data_time_stamp"`
	Sensor        SensorData `json:"sensor"`
}

// SensorData is the fixed phase-one field set returned by the Data API.
type SensorData struct {
	SensorIndex int      `json:"sensor_index"`
	Name        *string  `json:"name"`
	LastSeen    *int64   `json:"last_seen"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	PM25ATM     *float64 `json:"pm2.5_atm"`
}

type dataAPIRequestError struct {
	cause error
}

func (err dataAPIRequestError) Error() string {
	return "purpleair: data API request failed"
}

func (err dataAPIRequestError) Unwrap() error {
	return err.cause
}

// NewDataAPIClient creates an authenticated PurpleAir Data API client.
func NewDataAPIClient(readAPIKey string) (*DataAPIClient, error) {
	readAPIKey = strings.TrimSpace(readAPIKey)
	if readAPIKey == "" {
		return nil, fmt.Errorf("purpleair: API read key is required")
	}

	return &DataAPIClient{
		HTTPClient: defaultHTTPClient(),
		baseURL:    defaultDataAPIBaseURL,
		readAPIKey: readAPIKey,
	}, nil
}

func (c *DataAPIClient) httpClient() *http.Client {
	if c != nil && c.HTTPClient != nil {
		return c.HTTPClient
	}

	return defaultHTTPClient()
}

// SensorData retrieves authenticated real-time data for one sensor.
func (c *DataAPIClient) SensorData(ctx context.Context, sensorIndex int, options SensorDataOptions) (result *SensorDataResponse, returnErr error) {
	if ctx == nil {
		return nil, fmt.Errorf("purpleair: context is required")
	}
	if c == nil || strings.TrimSpace(c.readAPIKey) == "" {
		return nil, fmt.Errorf("purpleair: API read key is required")
	}
	if sensorIndex <= 0 {
		return nil, fmt.Errorf("purpleair: sensor index must be positive")
	}

	requestURL, err := c.sensorDataURL(sensorIndex, options)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("purpleair: create data API request")
	}
	request.Header.Set("X-API-Key", c.readAPIKey)

	response, err := c.httpClient().Do(request)
	if err != nil {
		return nil, dataAPIRequestError{cause: redactedDataAPIRequestCause(err)}
	}
	if response == nil || response.Body == nil {
		return nil, fmt.Errorf("purpleair: response body is empty")
	}
	defer func() {
		closeErr := response.Body.Close()
		if returnErr == nil && closeErr != nil {
			result = nil
			returnErr = fmt.Errorf("purpleair: close response body: %w", closeErr)
		}
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, dataAPIStatusError(response.StatusCode)
	}
	if response.ContentLength > maxDataAPIResponseBytes {
		return nil, fmt.Errorf("purpleair: response body exceeds %d bytes", maxDataAPIResponseBytes)
	}

	body, err := ioutil.ReadAll(io.LimitReader(response.Body, maxDataAPIResponseBytes+1))
	if err != nil {
		return nil, fmt.Errorf("purpleair: read response body: %w", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("purpleair: response body is empty")
	}
	if len(body) > maxDataAPIResponseBytes {
		return nil, fmt.Errorf("purpleair: response body exceeds %d bytes", maxDataAPIResponseBytes)
	}

	decoded := &SensorDataResponse{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(decoded); err != nil {
		return nil, fmt.Errorf("purpleair: decode response body: %w", err)
	}
	var trailing interface{}
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("purpleair: response body contains trailing data")
		}
		return nil, fmt.Errorf("purpleair: decode response body: %w", err)
	}
	if err := validateDataAPIResponse(decoded, sensorIndex); err != nil {
		return nil, err
	}

	return decoded, nil
}

func redactedDataAPIRequestCause(err error) error {
	for {
		urlErr, ok := err.(*url.Error)
		if !ok || urlErr.Err == nil {
			return err
		}
		err = urlErr.Err
	}
}

func (c *DataAPIClient) sensorDataURL(sensorIndex int, options SensorDataOptions) (string, error) {
	baseURL := strings.TrimRight(c.baseURL, "/")
	parsed, err := url.Parse(baseURL + "/sensors/" + strconv.Itoa(sensorIndex))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return "", fmt.Errorf("purpleair: data API endpoint is invalid")
	}
	query := parsed.Query()
	query.Set("fields", dataAPISensorFields)
	if readKey := strings.TrimSpace(options.SensorReadKey); readKey != "" {
		query.Set("read_key", readKey)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func dataAPIStatusError(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("purpleair: data API unauthorized (status %d)", statusCode)
	case http.StatusForbidden:
		return fmt.Errorf("purpleair: data API forbidden (status %d)", statusCode)
	case http.StatusNotFound:
		return fmt.Errorf("purpleair: data API sensor not found (status %d)", statusCode)
	case http.StatusTooManyRequests:
		return fmt.Errorf("purpleair: data API rate limited (status %d)", statusCode)
	default:
		return fmt.Errorf("purpleair: data API unexpected status %d", statusCode)
	}
}

func validateDataAPIResponse(response *SensorDataResponse, requestedSensorIndex int) error {
	if response.Sensor.SensorIndex != requestedSensorIndex {
		return fmt.Errorf("purpleair: response sensor index does not match request")
	}
	if response.Timestamp < 0 || response.DataTimestamp < 0 ||
		(response.Sensor.LastSeen != nil && *response.Sensor.LastSeen < 0) {
		return fmt.Errorf("purpleair: response timestamps must not be negative")
	}
	if (response.Sensor.Latitude == nil) != (response.Sensor.Longitude == nil) {
		return fmt.Errorf("purpleair: response coordinates must both be present or absent")
	}
	if response.Sensor.Latitude != nil {
		if !isFinite(*response.Sensor.Latitude) || *response.Sensor.Latitude < -90 || *response.Sensor.Latitude > 90 {
			return fmt.Errorf("purpleair: response latitude is invalid")
		}
		if !isFinite(*response.Sensor.Longitude) || *response.Sensor.Longitude < -180 || *response.Sensor.Longitude > 180 {
			return fmt.Errorf("purpleair: response longitude is invalid")
		}
	}
	if response.Sensor.PM25ATM != nil && !isFinite(*response.Sensor.PM25ATM) {
		return fmt.Errorf("purpleair: response PM2.5 is invalid")
	}
	return nil
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
