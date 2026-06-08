# PurpleAir Golang Parser

<!-- README-OVERVIEW-IMAGE -->
![Project overview](docs/readme-overview.svg)

## Install

`go get -u github.com/garethpaul/purpleair-go`

## Usage

```
package main

import (
	"purpleair"
	"fmt"
)

func main() {
	client := purpleair.NewClient()
	s:= client.Sensor("17937")
	for i := 0; i < len(s.Results); i++ {
        fmt.Println("Air Quality: " + s.Results[i].PM25Value)
    }
}
```
