package controllers

import (
	"github.com/dinuta/estuary-agent-go/src/command"
	"github.com/dinuta/estuary-agent-go/src/constants"
	u "github.com/dinuta/estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)

var CommandParallelPost = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, _ := ioutil.ReadAll(r.Body)
	commands := u.TrimSpacesAndLineEndings(strings.Split(string(body), "\n"))
	if len(commands) == 0 {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}
	cip := command.NewCommandInParallel()
	commandDescription := cip.RunCommands(commands)
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		commandDescription,
		r.URL.Path)

	u.ApiResponse(w, resp)
}
