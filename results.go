package purpleair

type Results struct {
    Results []Result `json:"results"`
}

type Result struct {
	ID                           int     `json:"ID"`
	Label                        string  `json:"Label"`
	DEVICELOCATIONTYPE           string  `json:"DEVICE_LOCATIONTYPE,omitempty"`
	THINGSPEAKPRIMARYID          string  `json:"THINGSPEAK_PRIMARY_ID"`
	THINGSPEAKPRIMARYIDREADKEY   string  `json:"THINGSPEAK_PRIMARY_ID_READ_KEY"`
	THINGSPEAKSECONDARYID        string  `json:"THINGSPEAK_SECONDARY_ID"`
	THINGSPEAKSECONDARYIDREADKEY string  `json:"THINGSPEAK_SECONDARY_ID_READ_KEY"`
	Lat                          float64 `json:"Lat"`
	Lon                          float64 `json:"Lon"`
	PM25Value                    string  `json:"PM2_5Value"`
	LastSeen                     int     `json:"LastSeen"`
	Type                         string  `json:"Type,omitempty"`
	Hidden                       string  `json:"Hidden"`
	DEVICEBRIGHTNESS             string  `json:"DEVICE_BRIGHTNESS,omitempty"`
	DEVICEHARDWAREDISCOVERED     string  `json:"DEVICE_HARDWAREDISCOVERED,omitempty"`
	Version                      string  `json:"Version,omitempty"`
	LastUpdateCheck              int     `json:"LastUpdateCheck,omitempty"`
	Created                      int     `json:"Created"`
	Uptime                       string  `json:"Uptime,omitempty"`
	RSSI                         string  `json:"RSSI,omitempty"`
	Adc                          string  `json:"Adc"`
	P03Um                        string  `json:"p_0_3_um"`
	P05Um                        string  `json:"p_0_5_um"`
	P10Um                        string  `json:"p_1_0_um"`
	P25Um                        string  `json:"p_2_5_um"`
	P50Um                        string  `json:"p_5_0_um"`
	P100Um                       string  `json:"p_10_0_um"`
	Pm10Cf1                      string  `json:"pm1_0_cf_1"`
	Pm25Cf1                      string  `json:"pm2_5_cf_1"`
	Pm100Cf1                     string  `json:"pm10_0_cf_1"`
	Pm10Atm                      string  `json:"pm1_0_atm"`
	Pm25Atm                      string  `json:"pm2_5_atm"`
	Pm100Atm                     string  `json:"pm10_0_atm"`
	IsOwner                      int     `json:"isOwner"`
	Humidity                     string  `json:"humidity,omitempty"`
	TempF                        string  `json:"temp_f,omitempty"`
	Pressure                     string  `json:"pressure,omitempty"`
	AGE                          int     `json:"AGE"`
	Stats                        string  `json:"Stats"`
	ParentID                     int     `json:"ParentID,omitempty"`
}

type PurpleAir struct {
	MapVersion       string `json:"mapVersion"`
	BaseVersion      string `json:"baseVersion"`
	MapVersionString string `json:"mapVersionString"`
	Results 		 []Result `json:"results"`
}
