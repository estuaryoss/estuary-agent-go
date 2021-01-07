package state

import (
	"os/exec"
	"sync"
)

type BackgroundCommand struct {
	commands map[string]*exec.Cmd
}

var once sync.Once

var singleton *BackgroundCommand

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
	bgCmd.commands[cmdId] = cmd
}

func (bgCmd *BackgroundCommand) GetBackgroundCommandList() map[string]*exec.Cmd {
	return bgCmd.commands
}
