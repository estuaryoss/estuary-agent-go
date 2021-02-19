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
	VERSION = "4.2.0"
)

func About() map[string]interface{} {
	virtualMemory, _ := mem.VirtualMemory()
	plat, fam, ver, _ := host.PlatformInformation()
	hostname, _ := os.Hostname()
	cpuInfo, _ := cpu.Info()
	return map[string]interface{}{
		"system":       runtime.GOOS,
		"platform":     plat,
		"release":      fam,
		"version":      ver,
		"architecture": runtime.GOARCH,
		"machine":      "NA",
		"layer":        getLayer(),
		"hostname":     hostname,
		"cpu":          cpuInfo[0].ModelName,
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
