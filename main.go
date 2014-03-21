package main

import (
	"fmt"
	"github.com/ulricqin/falcon/collector"
)

func main() {
	printSystemInfo()
}

func printSystemInfo() {
	fmt.Println(collector.KernelHostname())
	fmt.Println(collector.KernelMaxFiles())
	fmt.Println(collector.KernelMaxProc())
	fmt.Println(collector.LoadAvg())
	fmt.Println(collector.CpuSnapShoot())
}
