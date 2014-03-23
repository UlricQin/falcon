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
	mountPoints, _ := collector.ListMountPoint()
	fmt.Println(mountPoints)

	for idx := range mountPoints{
		fmt.Println(collector.BuildDeviceUsage(mountPoints[idx]))
	}

	netIfs, err := collector.NetIfs()
	fmt.Println("NetIfs.err: ", err)
	for _, netIf := range netIfs{
		fmt.Println(netIf)
	}

	fmt.Println(collector.MemInfo())
}
