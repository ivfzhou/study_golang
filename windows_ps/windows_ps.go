package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/shirou/gopsutil/v3/winservices"
)

func main() {}

func SampleUse() {
	fmt.Println(host.HostID())

	t, _ := host.BootTime()
	fmt.Println(time.Unix(int64(t), 0))

	info, _ := host.Info()
	fmt.Println(info)

	temperatures, _ := host.SensorsTemperatures()
	fmt.Println(temperatures)

	processes, _ := process.Processes()
	fmt.Println(processes)

	services, _ := winservices.ListServices()
	service := services[2]
	fmt.Println(service.Name, service.Status.State)

	avg, _ := load.Misc()
	fmt.Println(avg)
}
