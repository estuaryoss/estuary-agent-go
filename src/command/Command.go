package command

import (
	"github.com/estuaryoss/estuary-agent-go/src/models"
	"github.com/estuaryoss/estuary-agent-go/src/utils"
	"os"
	"time"
)

type Command struct {
	cmdId              string
	outputJsonPath     string
	enableStreams      bool
	commandDescription *models.CommandDescription
	commandsMap        map[string]*models.CommandStatus
}

func NewCommand(cmdId string, outputJsonPath string, enableStreams bool) *Command {
	initAt := time.Now()
	initAtString := utils.GetFormattedTimeAsString(initAt)
	command := &Command{
		cmdId:          cmdId,
		outputJsonPath: outputJsonPath,
		enableStreams:  enableStreams,
		commandDescription: &models.CommandDescription{false, false, initAtString,
			initAtString, 0, 0, "none", map[string]*models.CommandStatus{},
			[]*models.ProcessInfo{}},
		commandsMap: map[string]*models.CommandStatus{},
	}
	return command
}

func (com *Command) RunCommands(commands []string) *models.CommandDescription {
	startedAt := time.Now()
	com.commandDescription.SetPid(os.Getpid())
	com.commandDescription.SetId(com.cmdId)
	com.commandDescription.SetStarted(true)
	startedAtString := utils.GetFormattedTimeAsString(startedAt)
	com.commandDescription.SetStartedAt(startedAtString)
	com.setStatusForCommandsAsScheduled(commands)
	utils.WriteFileJson(com.outputJsonPath, com.commandDescription)

	com.runCommands(commands)

	finishedAt := time.Now()
	finishedAtString := utils.GetFormattedTimeAsString(finishedAt)
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
	startedAtString := utils.GetFormattedTimeAsString(startedAt)
	com.commandsMap[command].SetStartedAt(startedAtString)
	com.commandsMap[command].SetStatus("in progress")
	com.commandDescription.SetCommands(com.commandsMap)
	utils.WriteFileJson(com.outputJsonPath, com.commandDescription)

	var commandDetails *models.CommandDetails
	if com.enableStreams == true {
		commandDetails = utils.RunCommandToFile(command, com.cmdId)
	} else {
		commandDetails = utils.RunCommandNoFile(command, com.cmdId)
	}
	finishedAt := time.Now()
	finishedAtString := utils.GetFormattedTimeAsString(finishedAt)
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
