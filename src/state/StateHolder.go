package state

import (
	"fmt"

	"github.com/estuaryoss/estuary-agent-go/src/constants"
)

var State = fmt.Sprintf(constants.CMD_BACKGROUND_JSON_OUTPUT, "_")

func SetLastCommand(cmdId string) {
	State = fmt.Sprintf(constants.CMD_BACKGROUND_JSON_OUTPUT, cmdId)
}

func GetLastCommand() string {
	return State
}
