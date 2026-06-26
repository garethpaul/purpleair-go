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
	purpleAirUserAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36"
	maxSensorResponseBytes = 1 << 20
	maxSensorResults       = 1024
)

type sensorRequestError struct {
	cause error
}

func (err sensorRequestError) Error() string {
	return "purpleair: request failed"
}

func (err sensorRequestError) Unwrap() error {
	return err.cause
}

// Sensor gets sensor data and returns nil when the compatibility lookup fails.
//
// Deprecated: Use SensorWithError so request, response, and parsing failures
// remain available to the caller.
func (c *Client) Sensor(sensorId string) *PurpleAir {
	pa, err := c.SensorWithError(sensorId)
	if err != nil {
		return nil
	}

	return pa
}

// SensorWithError gets sensor data and returns request, response, and parsing errors.
func (c *Client) SensorWithError(sensorId string) (*PurpleAir, error) {
	return c.SensorWithContext(context.Background(), sensorId)
}

// SensorWithContext gets sensor data with caller-controlled cancellation and deadlines.
func (c *Client) SensorWithContext(ctx context.Context, sensorId string) (sensor *PurpleAir, returnErr error) {
	sensorId = strings.TrimSpace(sensorId)
	if sensorId == "" {
		return nil, fmt.Errorf("purpleair: sensor id is required")
	}
	for _, digit := range sensorId {
		if digit < '0' || digit > '9' {
			return nil, fmt.Errorf("purpleair: sensor id must be a positive integer")
		}
	}
	requestedSensorID, parseErr := strconv.Atoi(sensorId)
	if parseErr != nil || requestedSensorID <= 0 {
		return nil, fmt.Errorf("purpleair: sensor id must be a positive integer")
	}
	if ctx == nil {
		return nil, fmt.Errorf("purpleair: context is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.sensorURL(sensorId), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", purpleAirUserAgent)

	res, getErr := c.httpClient().Do(req)
	if getErr != nil {
		return nil, sensorRequestError{cause: getErr}
	}

	if res == nil {
		return nil, fmt.Errorf("purpleair: request failed with nil response")
	}

	if res.Body == nil {
		return nil, fmt.Errorf("purpleair: response body is empty")
	}

	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil && returnErr == nil {
			sensor = nil
			returnErr = fmt.Errorf("purpleair: close response body: %w", closeErr)
		}
	}()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("purpleair: unexpected status %d", res.StatusCode)
	}

	if res.ContentLength > maxSensorResponseBytes {
		return nil, fmt.Errorf("purpleair: response body exceeds %d bytes", maxSensorResponseBytes)
	}

	body, readErr := ioutil.ReadAll(io.LimitReader(res.Body, maxSensorResponseBytes+1))
	if readErr != nil {
		return nil, fmt.Errorf("purpleair: read response body: %w", readErr)
	}

	if len(body) > maxSensorResponseBytes {
		return nil, fmt.Errorf("purpleair: response body exceeds %d bytes", maxSensorResponseBytes)
	}

	if strings.TrimSpace(string(body)) == "" {
		return nil, fmt.Errorf("purpleair: response body is empty")
	}

	pa, decodeErr := decodeSensorResponse(body)
	if decodeErr != nil {
		return nil, fmt.Errorf("purpleair: decode response body: %w", decodeErr)
	}

	if len(pa.Results) == 0 {
		return nil, fmt.Errorf("purpleair: no results for sensor %q", sensorId)
	}

	matchedRequestedSensor := false
	for index, result := range pa.Results {
		if result.ID <= 0 {
			return nil, fmt.Errorf("purpleair: result %d has invalid sensor id %d", index, result.ID)
		}
		if result.ID == requestedSensorID {
			matchedRequestedSensor = true
		}
	}

	if !matchedRequestedSensor {
		return nil, fmt.Errorf("purpleair: response does not include requested sensor %d", requestedSensorID)
	}

	return pa, nil
}

func decodeSensorResponse(body []byte) (*PurpleAir, error) {
	var envelope struct {
		MapVersion       string          `json:"mapVersion"`
		BaseVersion      string          `json:"baseVersion"`
		MapVersionString string          `json:"mapVersionString"`
		Results          json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	results, err := decodeSensorResults(envelope.Results)
	if err != nil {
		return nil, err
	}

	return &PurpleAir{
		MapVersion:       envelope.MapVersion,
		BaseVersion:      envelope.BaseVersion,
		MapVersionString: envelope.MapVersionString,
		Results:          results,
	}, nil
}

func decodeSensorResults(rawResults json.RawMessage) ([]Result, error) {
	if len(rawResults) == 0 || bytes.Equal(bytes.TrimSpace(rawResults), []byte("null")) {
		return nil, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(rawResults))
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '[' {
		return nil, fmt.Errorf("results must be an array")
	}

	results := make([]Result, 0, 2)
	for decoder.More() {
		if len(results) >= maxSensorResults {
			return nil, fmt.Errorf("too many sensor results (maximum %d)", maxSensorResults)
		}

		var rawResult json.RawMessage
		if err := decoder.Decode(&rawResult); err != nil {
			return nil, err
		}

		var result Result
		if err := json.Unmarshal(rawResult, &result); err != nil {
			return nil, err
		}
		if result.ID > 0 {
			var coordinates struct {
				Lat *float64 `json:"Lat"`
				Lon *float64 `json:"Lon"`
			}
			if err := json.Unmarshal(rawResult, &coordinates); err != nil {
				return nil, err
			}
			if coordinates.Lat == nil || coordinates.Lon == nil {
				return nil, fmt.Errorf("result %d is missing coordinates", len(results))
			}
		}
		if math.IsNaN(result.Lat) || math.IsInf(result.Lat, 0) || math.IsNaN(result.Lon) || math.IsInf(result.Lon, 0) {
			return nil, fmt.Errorf("result %d has non-finite coordinates", len(results))
		}
		if result.Lat < -90 || result.Lat > 90 || result.Lon < -180 || result.Lon > 180 {
			return nil, fmt.Errorf("result %d has out-of-range coordinates (%g, %g)", len(results), result.Lat, result.Lon)
		}
		results = append(results, result)
	}

	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) sensorURL(sensorId string) string {
	baseURL := defaultBaseURL
	if c != nil && c.baseURL != "" {
		baseURL = c.baseURL
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return baseURL + "?show=" + url.QueryEscape(sensorId)
	}

	query := parsed.Query()
	query.Set("show", sensorId)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}
