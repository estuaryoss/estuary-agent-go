package utils

import (
	b64 "encoding/base64"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"os"
)

func getBase64HashForTheCommand(command string, cmdId string, suffix string) string {
	return constants.CMD_BACKGROUND_STREAMS_DIR + string(os.PathSeparator) +
		b64.StdEncoding.EncodeToString([]byte(command)) + "_" + cmdId + suffix
}
