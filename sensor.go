package purpleair

import (
	"encoding/json"
    "io/ioutil"
	"log"
	"net/http"
	"time"
)

// Get Sensor Data
func (c *Client) Sensor(sensorId string) (*PurpleAir) {

	url := "https://www.purpleair.com/json?show=" + sensorId

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

    // we initialize our Users array
    var pa PurpleAir

    // we unmarshal our byteArray which contains our
    // jsonFile's content into 'users' which we defined above
    json.Unmarshal(body, &pa)

	return &pa
}