package collector

import (
	"bufio"
	"bytes"
	log "github.com/ulricqin/goutils/logtool"
	"github.com/ulricqin/goutils/slicetool"
	"github.com/ulricqin/goutils/systool"
	"io"
	"strconv"
	"strings"
)

func ListenPorts() []int64 {
	bs, err := systool.CmdOutBytes("ss", "-n", "-l")
	if err != nil {
		log.Error("ss -n -l exec fail: %s", err)
		return []int64{}
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))

	// ignore the first line
	var line []byte
	line, _, err = reader.ReadLine()
	if err == io.EOF || err != nil {
		return []int64{}
	}

	ret := []int64{}

	for {
		line, _, err = reader.ReadLine()
		if err == io.EOF || err != nil {
			break
		}

		arr := strings.Fields(string(line))
		if len(arr) != 4 {
			log.Error("output of [ss -n -l] format error")
			continue
		}

		location := strings.LastIndex(arr[2], ":")
		port := arr[2][location+1:]

		if p, e := strconv.ParseInt(port, 10, 64); e != nil {
			log.Error("parse port to int64 fail: %s", e)
			continue
		} else {
			ret = append(ret, p)
		}

	}

	ret = slicetool.SliceUniqueInt64(ret)

	log.Info("listening ports: %v", ret)
	return ret
}
