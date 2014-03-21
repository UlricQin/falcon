package main

import (
	"fmt"
	"github.com/ulricqin/falcon/collector"
)

func main() {
	fmt.Println("Hello")
	fmt.Println(collector.KernelHostname())
	fmt.Println(collector.KernelMaxFiles())
	fmt.Println(collector.KernelMaxProc())
	fmt.Println(collector.LoadAvg())
}

