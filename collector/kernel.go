package collector

import (
	"github.com/ulricqin/goutils/filetool"
	"os"
)

func KernelMaxFiles() (uint64, error) {
	return filetool.FileToUint64("/proc/sys/fs/file-max")
}

func KernelMaxProc() (uint64, error) {
	return filetool.FileToUint64("/proc/sys/kernel/pid_max")
}

func KernelHostname() (string, error) {
	return os.Hostname()
}
