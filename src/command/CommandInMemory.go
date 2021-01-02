package command

import (
	"fmt"
	"github.com/dinuta/estuary-agent-go/src/models"
	"github.com/dinuta/estuary-agent-go/src/utils"
	"os"
	"time"
)

type CommandInMemory struct {
	commandDescription *models.CommandDescription
	commandsMap        map[string]*models.CommandStatus
}

func NewCommandInMemory() *CommandInMemory {
	initAt := time.Now()
	initAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		initAt.Year(), initAt.Month(), initAt.Day(),
		initAt.Hour(), initAt.Minute(), initAt.Second(), initAt.Nanosecond()/1000)
	commandInMemory := &CommandInMemory{
		commandDescription: &models.CommandDescription{false, false, initAtString,
			initAtString, 0, 0, "none", map[string]*models.CommandStatus{}},
		commandsMap: map[string]*models.CommandStatus{},
	}
	return commandInMemory
}

func (cim *CommandInMemory) RunCommands(commands []string) *models.CommandDescription {
	startedAt := time.Now()
	cim.commandDescription.SetPid(os.Getpid())
	cim.commandDescription.SetId("none")
	startedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		startedAt.Year(), startedAt.Month(), startedAt.Day(),
		startedAt.Hour(), startedAt.Minute(), startedAt.Second(), startedAt.Nanosecond()/1000)
	cim.commandDescription.SetStartedAt(startedAtString)

	cim.runCommands(commands)

	finishedAt := time.Now()
	finishedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		finishedAt.Year(), finishedAt.Month(), finishedAt.Day(),
		finishedAt.Hour(), finishedAt.Minute(), finishedAt.Second(), finishedAt.Nanosecond()/1000)
	cim.commandDescription.SetFinishedAt(finishedAtString)
	cim.commandDescription.SetDuration(finishedAt.Sub(startedAt).Seconds())
	cim.commandDescription.SetStarted(false)
	cim.commandDescription.SetFinished(true)

	return cim.commandDescription
}

func (cim *CommandInMemory) runCommands(commands []string) {
	for _, command := range commands {
		cim.runCommand(command)
	}

	cim.commandDescription.SetCommands(cim.commandsMap)
}

func (cim *CommandInMemory) runCommand(command string) {
	commandStatus := models.NewCommandStatus()

	startedAt := time.Now()
	startedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		startedAt.Year(), startedAt.Month(), startedAt.Day(),
		startedAt.Hour(), startedAt.Minute(), startedAt.Second(), startedAt.Nanosecond()/1000)
	commandStatus.SetStartedAt(startedAtString)

	commandDetails := utils.RunCommand(command)

	finishedAt := time.Now()
	finishedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		finishedAt.Year(), finishedAt.Month(), finishedAt.Day(),
		finishedAt.Hour(), finishedAt.Minute(), finishedAt.Second(), finishedAt.Nanosecond()/1000)
	commandStatus.SetFinishedAt(finishedAtString)

	commandStatus.SetDuration(finishedAt.Sub(startedAt).Seconds())
	commandStatus.SetStatus("finished")
	commandStatus.SetCommandDetails(commandDetails)
	cim.commandsMap[command] = commandStatus
}
