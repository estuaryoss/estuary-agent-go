package command

import (
	"github.com/estuaryoss/estuary-agent-go/src/models"
	"github.com/estuaryoss/estuary-agent-go/src/utils"
	"os"
	"time"
)

type CommandInMemory struct {
	commandDescription *models.CommandDescription
	commandsMap        map[string]*models.CommandStatus
}

func NewCommandInMemory() *CommandInMemory {
	initAt := time.Now()
	initAtString := utils.GetFormattedTimeAsString(initAt)
	commandInMemory := &CommandInMemory{
		commandDescription: &models.CommandDescription{false, false, initAtString,
			initAtString, 0, 0, "none", map[string]*models.CommandStatus{},
			[]*models.ProcessInfo{}},
		commandsMap: map[string]*models.CommandStatus{},
	}
	return commandInMemory
}

func (cim *CommandInMemory) RunCommands(commands []string) *models.CommandDescription {
	startedAt := time.Now()
	cim.commandDescription.SetPid(os.Getpid())
	cim.commandDescription.SetId("none")
	startedAtString := utils.GetFormattedTimeAsString(startedAt)
	cim.commandDescription.SetStartedAt(startedAtString)

	cim.runCommands(commands)

	finishedAt := time.Now()
	finishedAtString := utils.GetFormattedTimeAsString(finishedAt)
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
	startedAtString := utils.GetFormattedTimeAsString(startedAt)
	commandStatus.SetStartedAt(startedAtString)

	commandDetails := utils.RunCommand(command)

	finishedAt := time.Now()
	finishedAtString := utils.GetFormattedTimeAsString(finishedAt)
	commandStatus.SetFinishedAt(finishedAtString)

	commandStatus.SetDuration(finishedAt.Sub(startedAt).Seconds())
	commandStatus.SetStatus("finished")
	commandStatus.SetCommandDetails(commandDetails)
	cim.commandsMap[command] = commandStatus
}
