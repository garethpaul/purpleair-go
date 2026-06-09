package purpleair

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const purpleAirUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36"

// Get Sensor Data
func (c *Client) Sensor(sensorId string) *PurpleAir {
	pa, err := c.SensorWithError(sensorId)
	if err != nil {
		log.Fatal(err)
	}

	return pa
}

// SensorWithError gets sensor data and returns request, response, and parsing errors.
func (c *Client) SensorWithError(sensorId string) (*PurpleAir, error) {
	sensorId = strings.TrimSpace(sensorId)
	if sensorId == "" {
		return nil, fmt.Errorf("purpleair: sensor id is required")
	}

	req, err := http.NewRequest(http.MethodGet, c.sensorURL(sensorId), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", purpleAirUserAgent)

	res, getErr := c.httpClient().Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body == nil {
		return nil, fmt.Errorf("purpleair: response body is empty")
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("purpleair: unexpected status %d", res.StatusCode)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	if strings.TrimSpace(string(body)) == "" {
		return nil, fmt.Errorf("purpleair: response body is empty")
	}

	var pa PurpleAir

	if err := json.Unmarshal(body, &pa); err != nil {
		return nil, err
	}

	if len(pa.Results) == 0 {
		return nil, fmt.Errorf("purpleair: no results for sensor %q", sensorId)
	}

	return &pa, nil
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
