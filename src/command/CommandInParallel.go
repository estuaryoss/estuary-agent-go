package command

import (
	"fmt"
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
	initAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		initAt.Year(), initAt.Month(), initAt.Day(),
		initAt.Hour(), initAt.Minute(), initAt.Second(), initAt.Nanosecond()/1000)
	commandInParallel := &CommandInParallel{
		commandDescription: &models.CommandDescription{false, false, initAtString,
			initAtString, 0, 0, "none", map[string]*models.CommandStatus{}},
		commandsMap: map[string]*models.CommandStatus{},
	}
	return commandInParallel
}

func (cip *CommandInParallel) RunCommands(commands []string) *models.CommandDescription {
	startedAt := time.Now()
	cip.commandDescription.SetPid(os.Getpid())
	cip.commandDescription.SetId("none")
	startedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		startedAt.Year(), startedAt.Month(), startedAt.Day(),
		startedAt.Hour(), startedAt.Minute(), startedAt.Second(), startedAt.Nanosecond()/1000)
	cip.commandDescription.SetStartedAt(startedAtString)

	cip.runCommands(commands)

	finishedAt := time.Now()
	finishedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		finishedAt.Year(), finishedAt.Month(), finishedAt.Day(),
		finishedAt.Hour(), finishedAt.Minute(), finishedAt.Second(), finishedAt.Nanosecond()/1000)
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
	cip.commandsMap[command] = commandStatus

	ch <- commandStatus
}
