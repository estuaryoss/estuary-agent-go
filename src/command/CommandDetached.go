package command

import (
	"github.com/dinuta/estuary-agent-go/src/models"
)

type CommandDetached struct {
	commandDescription models.CommandDescription
	commandsMap        map[string]*models.CommandStatus
}

func NewCommandDetached() *CommandDetached {
	commandInBackground := &CommandDetached{
		commandDescription: models.CommandDescription{false, false, "", "", 0.1,
			10, "none", map[string]*models.CommandStatus{}},
		commandsMap: map[string]*models.CommandStatus{},
	}
	return commandInBackground
}

func (cim *CommandDetached) RunCommands(commands []string) models.CommandDescription {
	//TBD
	return cim.commandDescription
}
