package utils

import (
	"github.com/estuaryoss/estuary-agent-go/src/state"
	"syscall"
)

func KillProcess(cmdId string) {
	bgCmdList := state.GetInstance().GetBackgroundCommandList()
	if bgCmdList[cmdId] != nil {
		bgCmdList[cmdId].Process.Signal(syscall.SIGTERM)
		delete(bgCmdList, cmdId)
	}
}

func KillAllProcesses() {
	bgCmdList := state.GetInstance().GetBackgroundCommandList()
	for key, _ := range bgCmdList {
		bgCmdList[key].Process.Signal(syscall.SIGTERM)
		delete(bgCmdList, key)
	}
}
