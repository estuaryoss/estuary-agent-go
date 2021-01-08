package state

import (
	"github.com/mitchellh/go-ps"
	"os/exec"
	"sync"
	"syscall"
)

type BackgroundCommand struct {
	commands map[string]*exec.Cmd
}

var once sync.Once

var singleton *BackgroundCommand

/*
Keeping a list of pointers to the commands started in background
Needed for future process kill
{
	"muCmdId" : (exec.Cmd) pointer
}
*/
func GetInstance() *BackgroundCommand {
	once.Do(
		func() {
			singleton = &BackgroundCommand{
				commands: map[string]*exec.Cmd{},
			}
		})
	return singleton
}

func (bgCmd *BackgroundCommand) AddCmdToCommandList(cmdId string, cmd *exec.Cmd) {
	bgCmd.cleanAlreadyEndedCmdProcesses()
	bgCmd.commands[cmdId] = cmd
}

func (bgCmd *BackgroundCommand) GetBackgroundCommandList() map[string]*exec.Cmd {
	return bgCmd.commands
}

func (bgCmd *BackgroundCommand) cleanAlreadyEndedCmdProcesses() {
	for cmdId, _ := range bgCmd.commands {
		p, _ := ps.FindProcess(bgCmd.commands[cmdId].Process.Pid)
		if p != nil {
			bgCmd.commands[cmdId].Process.Signal(syscall.SIGTERM)
		}
	}
}
