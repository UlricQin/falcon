package collector

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type Meminfo struct {
	Buffers   uint64
	Cached    uint64
	MemTotal  uint64
	MemUsed   uint64
	MemFree   uint64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
	MemUsedPercent float64
	MemFreePercent float64
	SwapUsedPercent float64
	SwapFreePercent float64
}

func MemInfo() (*Meminfo, error) {
	want := map[string]bool{
		"Buffers:":   true,
		"Cached:":    true,
		"MemTotal:":  true,
		"MemFree:":   true,
		"SwapTotal:": true,
		"SwapFree:":  true}

	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	memInfo := &Meminfo{}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		fields := strings.Fields(string(line))
		fieldName := fields[0]

		_, ok := want[fieldName]
		if ok && len(fields) == 3 {
			val, numerr := strconv.ParseUint(fields[1], 10, 64)
			if numerr != nil {
				continue
			}
			switch fieldName {
			case "Buffers:":
				memInfo.Buffers = val
			case "Cached:":
				memInfo.Cached = val
			case "MemTotal:":
				memInfo.MemTotal = val
			case "MemFree:":
				memInfo.MemFree = val
			case "SwapTotal:":
				memInfo.SwapTotal = val
			case "SwapFree:":
				memInfo.SwapFree = val
			}
		}
	}
	memInfo.MemUsed = memInfo.MemTotal - memInfo.MemFree
	memInfo.SwapUsed = memInfo.SwapTotal - memInfo.SwapFree
	memInfo.MemFreePercent = float64(memInfo.MemFree) * 100.0 / float64(memInfo.MemTotal)
	memInfo.MemUsedPercent = float64(memInfo.MemUsed) * 100.0 / float64(memInfo.MemTotal)
	memInfo.SwapFreePercent = float64(memInfo.SwapFree) * 100.0 / float64(memInfo.SwapTotal)
	memInfo.SwapUsedPercent = float64(memInfo.SwapUsed) * 100.0 / float64(memInfo.SwapTotal)

	return memInfo, nil
}
