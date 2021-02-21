package utils

import (
	"github.com/estuaryoss/estuary-agent-go/src/models"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/estuaryoss/estuary-agent-go/src/state"
	"github.com/mitchellh/go-ps"
)

func KillCmdBackgroundProcess(cmdId string) {
	bgCmdList := state.GetInstance().GetBackgroundCommandList()
	if bgCmdList[cmdId] != nil {
		processesForPid := GetProcessesForPid(bgCmdList[cmdId].Process.Pid)
		log.Printf("Killing processes: %s for pid: %d", processesForPid, bgCmdList[cmdId].Process.Pid)
		KillProcesses(processesForPid) // <- kill processes
		bgCmdList[cmdId].Process.Signal(syscall.SIGTERM)
		delete(bgCmdList, cmdId)
	}
}

func GetAllProcesses() []*models.ProcessInfo {
	var processList = []*models.ProcessInfo{}
	ps, _ := ps.Processes()
	for _, process := range ps {
		pInfo := &models.ProcessInfo{}
		pInfo.Pid = process.Pid()
		pInfo.PPid = process.PPid()
		pInfo.Name = process.Executable()
		processList = append(processList, pInfo)
	}

	return processList
}

func GetAllProcessesByExecName(processName string) []*models.ProcessInfo {
	var processList = []*models.ProcessInfo{}
	ps, _ := ps.Processes()
	for _, process := range ps {
		if strings.Contains(process.Executable(), processName) {
			pInfo := &models.ProcessInfo{}
			pInfo.Pid = process.Pid()
			pInfo.PPid = process.PPid()
			pInfo.Name = process.Executable()
			processList = append(processList, pInfo)
		}
	}

	return processList
}

func GetProcessesForPid(pid int) []*models.ProcessInfo {
	var processList = []*models.ProcessInfo{}
	var currentProcess ps.Process
	pInfoCurrent := &models.ProcessInfo{}
	ps, _ := ps.Processes()

	for _, process := range ps {
		if process.Pid() == pid {
			pInfoCurrent.Pid = process.Pid()
			pInfoCurrent.PPid = process.PPid()
			pInfoCurrent.Name = process.Executable()
			currentProcess = process
			break
		}
	}

	for _, process := range ps {
		if currentProcess != nil {
			if process.PPid() == currentProcess.Pid() {
				pInfo := &models.ProcessInfo{}
				pInfo.Pid = process.Pid()
				pInfo.PPid = process.PPid()
				pInfo.Name = process.Executable()
				processList = append(processList, pInfo)
			}
		}
	}

	log.Printf("Discovered child processes: %s", processList)
	if pInfoCurrent.Pid != 0 {
		processList = append(processList, pInfoCurrent)
	}

	log.Printf("Discovered process list: %s, for pid: %d", processList, pid)
	return processList
}

func GetProcessByPid(pid int) []*models.ProcessInfo {
	var proc []*models.ProcessInfo
	ps, _ := ps.Processes()
	for _, process := range ps {
		if process.Pid() == pid {
			pInfo := &models.ProcessInfo{}
			pInfo.Pid = process.Pid()
			pInfo.PPid = process.PPid()
			pInfo.Name = process.Executable()
			proc = append(proc, pInfo)
		}
	}

	return proc
}

func KillProcesses(processes []*models.ProcessInfo) {
	for _, process := range processes {
		KillProcess(process.GetPid())
	}
}

func KillProcess(pid int) {
	p, err := os.FindProcess(pid)
	if err == nil {
		log.Printf("Killing process %d", pid)
		p.Signal(syscall.SIGTERM)
	}
}

func KillAllCmdBackgroundProcesses() {
	bgCmdList := state.GetInstance().GetBackgroundCommandList()
	for cmdId := range bgCmdList {
		KillCmdBackgroundProcess(cmdId)
		delete(bgCmdList, cmdId)
	}
}
