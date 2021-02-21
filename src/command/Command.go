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

func (cmd *Command) RunCommands(commands []string) *models.CommandDescription {
	startedAt := time.Now()
	cmd.commandDescription.SetPid(os.Getpid())
	cmd.commandDescription.SetId(cmd.cmdId)
	cmd.commandDescription.SetStarted(true)
	startedAtString := utils.GetFormattedTimeAsString(startedAt)
	cmd.commandDescription.SetStartedAt(startedAtString)
	cmd.setStatusForCommandsAsScheduled(commands)
	utils.WriteFileJson(cmd.outputJsonPath, cmd.commandDescription)

	cmd.runCommands(commands)

	finishedAt := time.Now()
	finishedAtString := utils.GetFormattedTimeAsString(finishedAt)
	cmd.commandDescription.SetFinishedAt(finishedAtString)
	cmd.commandDescription.SetDuration(finishedAt.Sub(startedAt).Seconds())
	cmd.commandDescription.SetStarted(false)
	cmd.commandDescription.SetFinished(true)
	utils.WriteFileJson(cmd.outputJsonPath, cmd.commandDescription)

	return cmd.commandDescription
}

func (cmd *Command) runCommands(commands []string) {
	for _, command := range commands {
		cmd.runCommand(command)
	}

	cmd.commandDescription.SetCommands(cmd.commandsMap)
}

func (cmd *Command) runCommand(command string) {
	startedAt := time.Now()
	startedAtString := utils.GetFormattedTimeAsString(startedAt)
	cmd.commandsMap[command].SetStartedAt(startedAtString)
	cmd.commandsMap[command].SetStatus("in progress")
	cmd.commandDescription.SetCommands(cmd.commandsMap)
	utils.WriteFileJson(cmd.outputJsonPath, cmd.commandDescription)

	var commandDetails *models.CommandDetails
	if cmd.enableStreams == true {
		commandDetails = utils.RunCommandToFile(command, cmd.cmdId)
	} else {
		commandDetails = utils.RunCommandNoFile(command, cmd.cmdId)
	}
	finishedAt := time.Now()
	finishedAtString := utils.GetFormattedTimeAsString(finishedAt)
	cmd.commandsMap[command].SetFinishedAt(finishedAtString)

	cmd.commandsMap[command].SetDuration(finishedAt.Sub(startedAt).Seconds())
	cmd.commandsMap[command].SetStatus("finished")
	cmd.commandsMap[command].SetCommandDetails(commandDetails)
	cmd.commandDescription.SetCommands(cmd.commandsMap)
	utils.WriteFileJson(cmd.outputJsonPath, cmd.commandDescription)
}

func (cmd *Command) setStatusForCommandsAsScheduled(commands []string) {
	for _, command := range commands {
		commandStatus := models.NewCommandStatus()
		commandStatus.SetStatus("scheduled")
		commandStatus.SetCommandDetails(models.NewCommandDetails())
		cmd.commandsMap[command] = commandStatus
	}
}
