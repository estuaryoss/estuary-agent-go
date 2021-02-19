package command

import (
	"github.com/estuaryoss/estuary-agent-go/src/models"
	"github.com/estuaryoss/estuary-agent-go/src/utils"
	"os"
	"time"
)

type CommandInParallel struct {
	commandDescription *models.CommandDescription
	commandsMap        map[string]*models.CommandStatus
}

func NewCommandInParallel() *CommandInParallel {
	initAt := time.Now()
	initAtString := utils.GetFormattedTimeAsString(initAt)
	commandInParallel := &CommandInParallel{
		commandDescription: &models.CommandDescription{false, false, initAtString,
			initAtString, 0, 0, "none", map[string]*models.CommandStatus{},
			[]models.ProcessInfo{}},
		commandsMap: map[string]*models.CommandStatus{},
	}
	return commandInParallel
}

func (cip *CommandInParallel) RunCommands(commands []string) *models.CommandDescription {
	startedAt := time.Now()
	cip.commandDescription.SetPid(os.Getpid())
	cip.commandDescription.SetId("none")
	startedAtString := utils.GetFormattedTimeAsString(startedAt)
	cip.commandDescription.SetStartedAt(startedAtString)

	cip.runCommands(commands)

	finishedAt := time.Now()
	finishedAtString := utils.GetFormattedTimeAsString(finishedAt)
	cip.commandDescription.SetFinishedAt(finishedAtString)
	cip.commandDescription.SetDuration(finishedAt.Sub(startedAt).Seconds())
	cip.commandDescription.SetStarted(false)
	cip.commandDescription.SetFinished(true)

	return cip.commandDescription
}

func (cip *CommandInParallel) runCommands(commands []string) {
	var commandChans []chan *models.CommandStatus
	for _, command := range commands {
		ch := make(chan *models.CommandStatus)
		commandChans = append(commandChans, ch)
		go cip.runCommand(command, ch)
	}

	for _, ch := range commandChans {
		<-ch
	}

	cip.commandDescription.SetCommands(cip.commandsMap)
}

func (cip *CommandInParallel) runCommand(command string, ch chan *models.CommandStatus) {
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
	cip.commandsMap[command] = commandStatus

	ch <- commandStatus
}
