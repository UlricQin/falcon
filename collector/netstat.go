package collector

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// @param ext e.g. TcpExt or IpExt
func Ext(ext string) (ret map[string]uint64, err error) {
	ret = make(map[string]uint64)
	var contents []byte
	contents, err = ioutil.ReadFile("/proc/net/netstat")
	if err != nil {
		return
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		lineBytes, _, e := reader.ReadLine()
		if e == io.EOF {
			break
		}

		raw_line := string(lineBytes)
		idx := strings.Index(raw_line, ":")
		if idx < 0 {
			continue
		}

		title := strings.TrimSpace(raw_line[:idx])
		if title == ext {
			ths := strings.Fields(strings.TrimSpace(raw_line[idx+1:]))
			// the next line must be values
			lineBytes, _, _ = reader.ReadLine()
			val_line := string(lineBytes)
			tds := strings.Fields(strings.TrimSpace(val_line[idx+1:]))
			for i := 0; i < len(ths); i++ {
				ret[ths[i]], _ = strconv.ParseUint(tds[i], 10, 64)
			}

			return
		}

	}

	return
}
