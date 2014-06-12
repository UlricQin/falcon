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
	MemFree   uint64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
}

var Multi uint64 = 1024

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
				memInfo.Buffers = val * Multi
			case "Cached:":
				memInfo.Cached = val * Multi
			case "MemTotal:":
				memInfo.MemTotal = val * Multi
			case "MemFree:":
				memInfo.MemFree = val * Multi
			case "SwapTotal:":
				memInfo.SwapTotal = val * Multi
			case "SwapFree:":
				memInfo.SwapFree = val * Multi
			}
		}
	}
	memInfo.SwapUsed = memInfo.SwapTotal - memInfo.SwapFree

	return memInfo, nil
}
