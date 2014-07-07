package collector

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ulricqin/goutils/filetool"
	log "github.com/ulricqin/goutils/logtool"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type Proc struct {
	Pid     int
	Name    string
	Cmdline string
	State   string
}

func AllProcs() (ps []*Proc, err error) {
	var dirs []string
	dirs, err = filetool.DirsUnder("/proc")
	if err != nil {
		return
	}

	// id dir is a number, it should be a pid. but don't trust it
	dirs_len := len(dirs)
	if dirs_len == 0 {
		return
	}

	var pid int
	var name_state [2]string
	var cmdline string
	for i := 0; i < dirs_len; i++ {
		if pid, err = strconv.Atoi(dirs[i]); err != nil {
			err = nil
			continue
		} else {
			status_file := fmt.Sprintf("/proc/%d/status", pid)
			cmdline_file := fmt.Sprintf("/proc/%d/cmdline", pid)
			if !filetool.IsExist(status_file) || !filetool.IsExist(cmdline_file) {
				continue
			}

			name_state, err = ReadProcStatus(status_file)
			if err != nil {
				log.Error("read %s fail: %s", status_file, err)
				continue
			}

			cmdline, err = filetool.ReadFileToStringNoLn(cmdline_file)
			if err != nil {
				log.Error("read %s fail: %s", cmdline_file, err)
				continue
			}

			p := Proc{Pid: pid, Name: name_state[0], State: name_state[1], Cmdline: cmdline}
			ps = append(ps, &p)
		}
	}

	return
}

func ReadProcStatus(path string) (name_state [2]string, err error) {
	var content []byte
	content, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	reader := bufio.NewReader(bytes.NewBuffer(content))
	name_done := false
	state_done := false
	for {
		bs, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		line := string(bs)

		colonIndex := strings.Index(line, ":")
		if line[0:colonIndex] == "Name" {
			name_state[0] = strings.TrimSpace(line[colonIndex+1:])
			name_done = true
			continue
		}

		if line[0:colonIndex] == "State" {
			name_state[1] = strings.TrimSpace(line[colonIndex+1:])
			state_done = true
			continue
		}

		if name_done && state_done {
			break
		}

	}

	return
}
