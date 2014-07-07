package collector

import (
	"github.com/ulricqin/goutils/filetool"
	"strconv"
	"strings"
)

type Loadavg struct {
	Avg1min  float64
	Avg5min  float64
	Avg15min float64
}

func LoadAvg() (*Loadavg, error) {

	loadAvg := Loadavg{}

	data, err := filetool.ReadFileToStringNoLn("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	slice := strings.Split(string(data), " ")
	if loadAvg.Avg1min, err = strconv.ParseFloat(slice[0], 64); err != nil {
		return nil, err
	}
	if loadAvg.Avg5min, err = strconv.ParseFloat(slice[1], 64); err != nil {
		return nil, err
	}
	if loadAvg.Avg15min, err = strconv.ParseFloat(slice[2], 64); err != nil {
		return nil, err
	}

	return &loadAvg, nil
}
