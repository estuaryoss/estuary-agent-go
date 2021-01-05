package controllers

import (
	"github.com/estuaryoss/estuary-agent-go/src/command"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

var CommandPost = func(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	commands := u.TrimSpacesAndLineEndings(strings.Split(string(body), "\n"))

	if len(commands) == 0 {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}
	cim := command.NewCommandInMemory()
	commandDescription := cim.RunCommands(commands)
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		commandDescription,
		r.URL.Path)

	u.ApiResponse(w, resp)
}
