package collector

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type NetIf struct {
	Iface          string
	InBytes        int64
	InPackages     int64
	InErrors       int64
	InDropped      int64
	InFifoErrs     int64
	InFrameErrs    int64
	InCompressed   int64
	InMulticast    int64
	OutBytes       int64
	OutPackages    int64
	OutErrors      int64
	OutDropped     int64
	OutFifoErrs    int64
	OutCollisions  int64
	OutCarrierErrs int64
	OutCompressed  int64
	TotalBytes     int64
	TotalPackages  int64
	TotalErrors    int64
	TotalDropped   int64
}

/*
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0: 1990350    2838    0    0    0     0          0         0   401351    2218    0    0    0     0       0          0
    lo:   26105     286    0    0    0     0          0         0    26105     286    0    0    0     0       0          0
*/
func NetIfs() ([]*NetIf, error) {
	contents, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	ret := []*NetIf{}

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		raw_line := string(lineBytes)
		idx := strings.Index(raw_line, ":")
		if idx < 0 {
			continue
		}

		netIf := NetIf{}

		line := raw_line[idx+1 : len(raw_line)]
		eth := raw_line[0:idx]
		netIf.Iface = strings.Trim(eth, " ")

		fields := strings.Fields(line)

		netIf.InBytes, _ = strconv.ParseInt(fields[0], 10, 64)
		netIf.InPackages, _ = strconv.ParseInt(fields[1], 10, 64)
		netIf.InErrors, _ = strconv.ParseInt(fields[2], 10, 64)
		netIf.InDropped, _ = strconv.ParseInt(fields[3], 10, 64)
		netIf.InFifoErrs, _ = strconv.ParseInt(fields[4], 10, 64)
		netIf.InFrameErrs, _ = strconv.ParseInt(fields[5], 10, 64)
		netIf.InCompressed, _ = strconv.ParseInt(fields[6], 10, 64)
		netIf.InMulticast, _ = strconv.ParseInt(fields[7], 10, 64)

		netIf.OutBytes, _ = strconv.ParseInt(fields[8], 10, 64)
		netIf.OutPackages, _ = strconv.ParseInt(fields[9], 10, 64)
		netIf.OutErrors, _ = strconv.ParseInt(fields[10], 10, 64)
		netIf.OutDropped, _ = strconv.ParseInt(fields[11], 10, 64)
		netIf.OutFifoErrs, _ = strconv.ParseInt(fields[12], 10, 64)
		netIf.OutCollisions, _ = strconv.ParseInt(fields[13], 10, 64)
		netIf.OutCarrierErrs, _ = strconv.ParseInt(fields[14], 10, 64)
		netIf.OutCompressed, _ = strconv.ParseInt(fields[15], 10, 64)

		netIf.TotalBytes = netIf.InBytes + netIf.OutBytes
		netIf.TotalPackages = netIf.InPackages + netIf.OutPackages
		netIf.TotalErrors = netIf.InErrors + netIf.OutErrors
		netIf.TotalDropped = netIf.InDropped + netIf.OutDropped

		ret = append(ret, &netIf)
	}

	return ret, nil
}
