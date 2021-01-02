package command

import (
	"fmt"
	"github.com/dinuta/estuary-agent-go/src/models"
	"github.com/dinuta/estuary-agent-go/src/utils"
	"os"
	"time"
)

type Command struct {
	cmdId              string
	outputJsonPath     string
	commandDescription *models.CommandDescription
	commandsMap        map[string]*models.CommandStatus
}

func NewCommand(cmdId string, outputJsonPath string) *Command {
	initAt := time.Now()
	initAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		initAt.Year(), initAt.Month(), initAt.Day(),
		initAt.Hour(), initAt.Minute(), initAt.Second(), initAt.Nanosecond()/1000)

	command := &Command{
		cmdId:          cmdId,
		outputJsonPath: outputJsonPath,
		commandDescription: &models.CommandDescription{false, false, initAtString,
			initAtString, 0, 0, "none", map[string]*models.CommandStatus{}},
		commandsMap: map[string]*models.CommandStatus{},
	}
	return command
}

func (com *Command) RunCommands(commands []string) *models.CommandDescription {
	startedAt := time.Now()
	com.commandDescription.SetPid(os.Getpid())
	com.commandDescription.SetId(com.cmdId)
	com.commandDescription.SetStarted(true)
	startedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		startedAt.Year(), startedAt.Month(), startedAt.Day(),
		startedAt.Hour(), startedAt.Minute(), startedAt.Second(), startedAt.Nanosecond()/1000)
	com.commandDescription.SetStartedAt(startedAtString)
	com.setStatusForCommandsAsScheduled(commands)
	utils.WriteFileJson(com.outputJsonPath, com.commandDescription)

	com.runCommands(commands)

	finishedAt := time.Now()
	finishedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		finishedAt.Year(), finishedAt.Month(), finishedAt.Day(),
		finishedAt.Hour(), finishedAt.Minute(), finishedAt.Second(), finishedAt.Nanosecond()/1000)
	com.commandDescription.SetFinishedAt(finishedAtString)
	com.commandDescription.SetDuration(finishedAt.Sub(startedAt).Seconds())
	com.commandDescription.SetStarted(false)
	com.commandDescription.SetFinished(true)
	utils.WriteFileJson(com.outputJsonPath, com.commandDescription)

	return com.commandDescription
}

func (com *Command) runCommands(commands []string) {
	for _, command := range commands {
		com.runCommand(command)
	}

	com.commandDescription.SetCommands(com.commandsMap)
}

func (com *Command) runCommand(command string) {
	startedAt := time.Now()
	startedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		startedAt.Year(), startedAt.Month(), startedAt.Day(),
		startedAt.Hour(), startedAt.Minute(), startedAt.Second(), startedAt.Nanosecond()/1000)
	com.commandsMap[command].SetStartedAt(startedAtString)
	com.commandsMap[command].SetStatus("in progress")
	com.commandDescription.SetCommands(com.commandsMap)
	utils.WriteFileJson(com.outputJsonPath, com.commandDescription)

	commandDetails := utils.RunCommand(command)

	finishedAt := time.Now()
	finishedAtString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		finishedAt.Year(), finishedAt.Month(), finishedAt.Day(),
		finishedAt.Hour(), finishedAt.Minute(), finishedAt.Second(), finishedAt.Nanosecond()/1000)
	com.commandsMap[command].SetFinishedAt(finishedAtString)

	com.commandsMap[command].SetDuration(finishedAt.Sub(startedAt).Seconds())
	com.commandsMap[command].SetStatus("finished")
	com.commandsMap[command].SetCommandDetails(commandDetails)
	com.commandDescription.SetCommands(com.commandsMap)
	utils.WriteFileJson(com.outputJsonPath, com.commandDescription)
}

func (com *Command) setStatusForCommandsAsScheduled(commands []string) {
	for _, cmd := range commands {
		commandStatus := models.NewCommandStatus()
		commandStatus.SetStatus("scheduled")
		commandStatus.SetCommandDetails(models.NewCommandDetails())
		com.commandsMap[cmd] = commandStatus
	}
}
