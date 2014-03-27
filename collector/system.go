package collector

import (
	"errors"
	"github.com/ulricqin/goutils/filetool"
	"github.com/ulricqin/goutils/systool"
	"strconv"
	"strings"
)

func SystemDate() (string, error) {
	return systool.CmdOutNoLn("/bin/date")
}

func SystemUptime() (arr [3]int64, err error) {
	var content string
	content, err = filetool.ReadFileToStringNoLn("/proc/uptime")
	if err != nil {
		return
	}

	fields := strings.Fields(content)
	if len(fields) < 2 {
		err = errors.New("/proc/uptime parse error")
		return
	}

	secStr := fields[0]
	var secF float64
	secF, err = strconv.ParseFloat(secStr, 64)
	if err != nil {
		return
	}

	minTotal := secF / 60.0
	hourTotal := minTotal / 60.0

	days := int64(hourTotal / 24.0)
	hours := int64(hourTotal) - days*24
	mins := int64(minTotal) - (days * 60 * 24) - (hours * 60)

	return [3]int64{days, hours, mins}, nil
}
