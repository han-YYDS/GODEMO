package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func testMem() {
	v, _ := mem.VirtualMemory()
	fmt.Printf("Total: %v, Available: %v, UsedPercent:%f%%\n", v.Total, v.Available, v.UsedPercent)
	fmt.Println(v)
}

func testCPU() {
	// cpuCount()
	// cpuInfo()
	// cpuRate()
	cpuTime()
}

func cpuCount() {
	// CPU核数统计
	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)
	fmt.Printf("physical count:%d logical count:%d\n", physicalCnt, logicalCnt)
}
func cpuRate() {
	// 收集3s内的CPU总占用率以及各个CPU的使用率
	totalPercent, _ := cpu.Percent(3*time.Second, false)
	perPercents, _ := cpu.Percent(3*time.Second, true)
	fmt.Printf("total percent:%v \n per percents:%v", totalPercent, perPercents)
}

func cpuInfo() {
	// CPU详情
	infos, _ := cpu.Info()
	for _, info := range infos {
		data, _ := json.MarshalIndent(info, "", " ")
		fmt.Print(string(data))
	}
}

func cpuTime() {
	infos, _ := cpu.Times(true)
	for _, info := range infos {
		data, _ := json.MarshalIndent(info, "", " ")
		fmt.Print(string(data))
	}
}

func testStat() {
	mapStat, _ := disk.IOCounters()
	for name, stat := range mapStat {
		fmt.Println(name)
		data, _ := json.MarshalIndent(stat, "", "  ")
		fmt.Println(string(data))
	}
}

func testPartition() {
	infos, _ := disk.Partitions(false)
	for _, info := range infos {
		data, _ := json.MarshalIndent(info, "", "  ")
		fmt.Println(string(data))
	}
}

func testDiskUsage() {
	info, _ := disk.Usage("/home/ubuntu20")
	data, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(data))
}

func testDisk() {
	// testStat()
	// testPartition()
	testDiskUsage()
}

func testBootTime() {
	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)
	fmt.Println(t.Local().Format("2006-01-02 15:04:05"))
}

func testHost() {
	testBootTime()
}

func main() {
	// testMem()
	// testCPU()
	// testDisk()
	testHost()
}
