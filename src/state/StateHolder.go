package state

import (
	"fmt"

	"github.com/estuaryoss/estuary-agent-go/src/constants"
)

type StateHolder struct {
	State string
}

func (sh *StateHolder) SetLastCommand(cmdId string) {
	sh.State = fmt.Sprintf(constants.CMD_BACKGROUND_JSON_OUTPUT, cmdId)
}

func (sh *StateHolder) GetLastCommand(cmdId string) string {
	return sh.State
}
