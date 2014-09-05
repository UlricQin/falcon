package collector

import (
	"bufio"
	"bytes"
	log "github.com/ulricqin/goutils/logtool"
	"github.com/ulricqin/goutils/systool"
	"io"
	"strconv"
	"strings"
)

func SocketStatSummary() (m map[string]uint64) {
	m = make(map[string]uint64)
	bs, err := systool.CmdOutBytes("ss", "-s")
	if err != nil {
		log.Error("ss -s exec fail: %s", err)
		return
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))

	// ignore the first line
	line, _, err := reader.ReadLine()
	if err == io.EOF || err != nil {
		return
	}

	for {
		line, _, err = reader.ReadLine()
		if err == io.EOF || err != nil {
			return
		}

		lineStr := string(line)
		if strings.HasPrefix(lineStr, "TCP") {
			left := strings.Index(lineStr, "(")
			right := strings.Index(lineStr, ")")
			if left < 0 || right < 0 {
				continue
			}

			content := lineStr[left+1 : right]
			arr := strings.Split(content, ", ")
			for _, val := range arr {
				fields := strings.Fields(val)
				if fields[0] == "timewait" {
					timewait_arr := strings.Split(fields[1], "/")
					m["timewait"], _ = strconv.ParseUint(timewait_arr[0], 10, 64)
					m["slabinfo.timewait"], _ = strconv.ParseUint(timewait_arr[1], 10, 64)
					continue
				}
				m[fields[0]], _ = strconv.ParseUint(fields[1], 10, 64)
			}
			return
		}
	}

	return
}
