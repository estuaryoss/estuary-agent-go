package command

import (
	"fmt"
	"github.com/dinuta/estuary-agent-go/src/models"
	"github.com/dinuta/estuary-agent-go/src/utils"
	"os"
	"time"
)

type CommandInParallel struct {
	commandDescription models.CommandDescription
	commandsMap        map[string]*models.CommandStatus
}

func NewCommandInParallel() *CommandInParallel {
	commandInParallel := &CommandInParallel{
		commandDescription: models.CommandDescription{false, false, "", "", 0.1,
			10, "none", map[string]*models.CommandStatus{}},
		commandsMap: map[string]*models.CommandStatus{},
	}
	return commandInParallel
}

func (cim *CommandInParallel) RunCommands(commands []string) models.CommandDescription {
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

func (cim *CommandInParallel) runCommands(commands []string) {
	var commandChans []chan *models.CommandStatus
	for _, command := range commands {
		ch := make(chan *models.CommandStatus)
		commandChans = append(commandChans, ch)
		go cim.runCommand(command, ch)
	}

	for _, ch := range commandChans {
		<-ch
	}

	cim.commandDescription.SetCommands(cim.commandsMap)
}

func (cim *CommandInParallel) runCommand(command string, ch chan *models.CommandStatus) {
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

	ch <- commandStatus
}
