package purpleair

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleClient_SensorWithError() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{"results":[{"ID":17937,"Label":"Example Sensor","Lat":0,"Lon":0}]}`)
	}))
	defer server.Close()

	client := NewClient()
	client.baseURL = server.URL + "/json"
	client.HTTPClient = server.Client()

	sensor, err := client.SensorWithError("17937")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(sensor.Results[0].ID)
	fmt.Println(sensor.Results[0].Label)

	// Output:
	// 17937
	// Example Sensor
}

func ExampleClient_SensorWithError_error() {
	client := NewClient()

	sensor, err := client.SensorWithError(" ")

	fmt.Println(sensor == nil)
	fmt.Println(err)

	// Output:
	// true
	// purpleair: sensor id is required
}
