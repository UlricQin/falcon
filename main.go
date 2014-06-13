package main

import (
	"fmt"
	"github.com/ulricqin/falcon/collector"
)

func main() {
	// printSystemInfo()
	fmt.Println(collector.ListenTcpPorts())
}

func printSystemInfo() {
	fmt.Println(collector.KernelHostname())
	fmt.Println(collector.KernelMaxFiles())
	fmt.Println(collector.KernelMaxProc())
	fmt.Println(collector.LoadAvg())
	fmt.Println(collector.CpuSnapShoot())
	mountPoints, _ := collector.ListMountPoint()
	fmt.Println(mountPoints)

	for idx := range mountPoints {
		fmt.Println(collector.BuildDeviceUsage(mountPoints[idx]))
	}

	netIfs, err := collector.NetIfs()
	fmt.Println("NetIfs.err: ", err)
	for _, netIf := range netIfs {
		fmt.Println(netIf)
	}

	fmt.Println(collector.MemInfo())

	fmt.Println("tcp ports: ", collector.ListenPorts("tcp"))
	fmt.Println("tcp6 ports: ", collector.ListenPorts("tcp6"))
	fmt.Println("udp ports: ", collector.ListenPorts("udp"))
	fmt.Println("udp6 ports: ", collector.ListenPorts("udp6"))

	fmt.Println(collector.SystemUptime())

	diskStats, _ := collector.ListDiskStats()
	for _, ds := range diskStats {
		fmt.Println(ds)
	}
}
