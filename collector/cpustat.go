package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
)

type CpuSnapshoot struct {
	User     uint64 // time spent in user mode
	Nice     uint64 // time spent in user mode with low priority (nice)
	System   uint64 // time spent in system mode
	Idle     uint64 // time spent in the idle task
	Iowait   uint64 // time spent waiting for I/O to complete (since Linux 2.5.41)
	Irq      uint64 // time spent servicing  interrupts  (since  2.6.0-test4)
	SoftIrq  uint64 // time spent servicing softirqs (since 2.6.0-test4)
	Steal    uint64 // time spent in other OSes when running in a virtualized environment
	Guest    uint64 // time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	Total    uint64 // total of all time fields
	Switches uint64 // context switches
}

func CpuNum() int {
	return runtime.NumCPU()
}

func MHz() (mhz string, err error) {
	var contents []byte
	contents, err = ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))

	var line []byte
	for {
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			return
		}

		li := string(line)
		if !strings.Contains(li, "MHz") {
			continue
		}

		arr := strings.Split(li, ":")
		if len(arr) != 2 {
			return "", fmt.Errorf("file content format error")
		}

		return strings.TrimSpace(arr[1]), nil
	}
}

func CpuSnapShoot() (*CpuSnapshoot, error) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	s := &CpuSnapshoot{}

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		fields := strings.Fields(string(line))
		if len(fields) > 0 {
			fieldName := fields[0]
			if fieldName == "cpu" {
				parseCPUFields(fields, s)
			}

			if fieldName == "ctxt" {
				s.Switches, err = strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					continue
				}
			}
		}
	}
	return s, nil
}

func parseCPUFields(fields []string, stat *CpuSnapshoot) {
	numFields := len(fields)
	for i := 1; i < numFields; i++ {
		val, err := strconv.ParseUint(fields[i], 10, 64)
		if err != nil {
			continue
		}

		stat.Total += val
		switch i {
		case 1:
			stat.User = val
		case 2:
			stat.Nice = val
		case 3:
			stat.System = val
		case 4:
			stat.Idle = val
		case 5:
			stat.Iowait = val
		case 6:
			stat.Irq = val
		case 7:
			stat.SoftIrq = val
		case 8:
			stat.Steal = val
		case 9:
			stat.Guest = val
		}
	}
}
