package collector

import (
	"bufio"
	"fmt"
	"github.com/ulricqin/goutils/slicetool"
	"io"
	"os"
	"strconv"
	"strings"
)

var protocol_search = map[string]string{
	"tcp":  " 00000000:0000 0A",
	"tcp6": " 00000000000000000000000000000000:0000 0A",
	"udp":  " 00000000:0000 07",
	"udp6": " 00000000000000000000000000000000:0000 07",
}

func ListenTcpPorts() []int64 {
	ret := []int64{}

	procFile := "/proc/net/tcp"

	f, err := os.Open(procFile)
	if err != nil {
		return ret
	}
	defer f.Close()

	search := protocol_search["tcp"]

	reader := bufio.NewReader(f)

	// ignore the first line
	lineBytes, _, err := reader.ReadLine()
	if err == io.EOF || err != nil {
		return []int64{}
	}

	// handle line:0
	lineBytes, _, err = reader.ReadLine()
	if err == io.EOF || err != nil {
		return []int64{}
	}

	spaceCnt := 0
	for _, b := range lineBytes {
		if b == ' ' {
			spaceCnt++
		} else {
			break
		}
	}

	index := 32 + spaceCnt
	if lineBytes[index] == 'A' {
		rawLine := string(lineBytes)
		fmt.Println(rawLine)
		idx := strings.Index(rawLine, search)
		if idx < 0 {
			fmt.Println("not found")
		} else {
			fmt.Println("found")
			portStr := rawLine[idx-4 : idx]
			i, err := strconv.ParseInt(portStr, 16, 64)
			if err != nil {
				fmt.Println("parse int fail" + err.Error())
			} else {
				ret = append(ret, i)
			}
		}
	}

	for {
		lineBytes, _, err = reader.ReadLine()
		if err == io.EOF || err != nil {
			break
		}

		if lineBytes[index] == 'A' {
			rawLine := string(lineBytes)
			fmt.Println(rawLine)
			idx := strings.Index(rawLine, search)
			if idx < 0 {
				fmt.Println("not found")
			} else {
				fmt.Println("found")
				portStr := rawLine[idx-4 : idx]
				i, err := strconv.ParseInt(portStr, 16, 64)
				if err != nil {
					fmt.Println("parse int fail" + err.Error())
				} else {
					ret = append(ret, i)
				}
			}
		}

	}
	return slicetool.SliceUniqueInt64(ret)
}

// protocol: ['tcp', 'tcp6', 'udp', 'udp6']
func ListenPorts(protocol string) []int64 {
	ret := []int64{}

	procFile := "/proc/net/" + protocol

	f, err := os.Open(procFile)
	if err != nil {
		return ret
	}
	defer f.Close()

	search := protocol_search[protocol]

	reader := bufio.NewReader(f)

	firstLine := true
	for {
		lineBytes, _, err := reader.ReadLine()
		if err == io.EOF || err != nil {
			break
		}

		if firstLine {
			firstLine = false
			continue
		}

		rawLine := string(lineBytes)
		idx := strings.Index(rawLine, search)
		if idx < 0 {
			// if protocol is tcp. 03 maybe in the middle of 0A
			if protocol == "tcp" && strings.Index(rawLine, " 03 ") > 0 {
				continue
			}

			if protocol == "tcp6" {
				continue
			}

			break
		}

		portStr := rawLine[idx-4 : idx]
		i, err := strconv.ParseInt(portStr, 16, 64)
		if err != nil {
			continue
		}

		ret = append(ret, i)

	}
	return slicetool.SliceUniqueInt64(ret)
}

/*
ulric@ubuntu:~$ cat /proc/net/tcp
  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0100007F:0CEA 00000000:0000 0A 00000000:00000000 00:00000000 00000000   115        0 14563 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 13321 1 0000000000000000 100 0 0 10 0
   2: 0100007F:12D6 00000000:0000 0A 00000000:00000000 00:00000000 00000000   999        0 9794 1 0000000000000000 100 0 0 10 0
   3: 0100007F:0277 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 39066 1 0000000000000000 100 0 0 10 0
   4: 8284A8C0:89C3 195EBD5B:0050 08 00000000:00000001 00:00000000 00000000  1000        0 19408 1 0000000000000000 20 4 18 2 -1
   5: 8284A8C0:CEDA 017AA8C0:0C38 01 00000000:00000000 02:00001031 00000000  1000        0 119115 2 0000000000000000 21 4 1 10 -1
   6: 8284A8C0:B151 7D807D4A:1466 01 00000000:00000000 02:0000057E 00000000  1000        0 17091 2 0000000000000000 20 4 28 10 -1
ulric@ubuntu:~$ cat /proc/net/tcp6
  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 0000000000000000FFFF00000100007F:1F45 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000   116        0 12621 1 0000000000000000 100 0 0 10 -1
   1: 00000000000000000000000000000000:1F90 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000   116        0 10657 1 0000000000000000 100 0 0 10 -1
   2: 00000000000000000000000001000000:0277 00000000000000000000000000000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 39065 1 0000000000000000 100 0 0 10 -1
ulric@ubuntu:~$ cat /proc/net/udp
  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops
  944: 00000000:14E9 00000000:0000 07 00000000:00000000 00:00000000 00000000   107        0 8549 2 0000000000000000 0
  962: 00000000:BCFB 00000000:0000 07 00000000:00000000 00:00000000 00000000   107        0 8551 2 0000000000000000 0
 1788: 0100007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 13320 2 0000000000000000 0
 1803: 00000000:0044 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 7961 2 0000000000000000 0
ulric@ubuntu:~$ cat /proc/net/udp6
  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops
  329: 00000000000000000000000000000000:9282 00000000000000000000000000000000:0000 07 00000000:00000000 00:00000000 00000000   107        0 8552 2 0000000000000000 0
  944: 00000000000000000000000000000000:14E9 00000000000000000000000000000000:0000 07 00000000:00000000 00:00000000 00000000   107        0 8550 2 0000000000000000 0
ulric@ubuntu:~$
*/
