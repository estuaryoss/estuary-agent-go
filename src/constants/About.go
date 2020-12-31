package constants

import (
	"estuary-agent-go/src/io"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
	"runtime"
)

const (
	Name    = "estuary-agent"
	Version = "4.1.0"
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
	layer := "Machine"
	if io.DoesFileExists("/.dockerenv") {
		layer = "Docker"
	}
	return layer
}
