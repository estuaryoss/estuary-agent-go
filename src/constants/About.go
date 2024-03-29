package constants

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
	"runtime"
)

const (
	NAME    = "estuary-agent"
	VERSION = "4.2.1"
)

func About() map[string]interface{} {
	virtualMemory, _ := mem.VirtualMemory()
	plat, fam, ver, _ := host.PlatformInformation()
	var cpuInfo = "NA"
	var hostNameInfo = "NA"
	hostName, err := os.Hostname()
	if err == nil {
		hostNameInfo = hostName
	}
	cpu_, err := cpu.Info()
	if err == nil {
		cpuInfo = cpu_[0].ModelName
	}

	return map[string]interface{}{
		"system":       runtime.GOOS,
		"platform":     plat,
		"release":      fam,
		"version":      ver,
		"architecture": runtime.GOARCH,
		"machine":      "NA",
		"layer":        getLayer(),
		"hostname":     hostNameInfo,
		"cpu":          cpuInfo,
		"ram":          fmt.Sprint(virtualMemory.Total/(1024*1024*1024), " GB"),
		"golang":       runtime.Version(),
	}
}

func getLayer() string {
	var layer string
	if _, err := os.Stat("/.dockerenv"); err == nil {
		layer = "Docker"
	} else {
		layer = "Machine"
	}
	return layer
}
