package utils

import (
	"os"
	"syscall"

	"github.com/estuaryoss/estuary-agent-go/src/state"
	"github.com/mitchellh/go-ps"
)

func KillCmdBackgroundProcess(cmdId string) {
	bgCmdList := state.GetInstance().GetBackgroundCommandList()
	if bgCmdList[cmdId] != nil {
		childProcesses := GetChildListForParentProcess(bgCmdList[cmdId].Process.Pid)
		KillProcesses(childProcesses)
		bgCmdList[cmdId].Process.Signal(syscall.SIGTERM)
		delete(bgCmdList, cmdId)
	}
}

func GetChildListForParentProcess(PPid int) []ps.Process {
	var childProcessList = []ps.Process{}
	ps, _ := ps.Processes()
	for _, process := range ps {
		if process.PPid() == PPid {
			childProcessList = append(childProcessList, process)
		}
	}

	return childProcessList
}

func KillProcesses(processes []ps.Process) {
	for _, process := range processes {
		p, err := os.FindProcess(process.Pid())
		if err == nil {
			p.Signal(syscall.SIGTERM)
		}
	}
}

func KillAllCmdBackgroundProcesses() {
	bgCmdList := state.GetInstance().GetBackgroundCommandList()
	for cmdId := range bgCmdList {
		KillCmdBackgroundProcess(cmdId)
		delete(bgCmdList, cmdId)
	}
}
