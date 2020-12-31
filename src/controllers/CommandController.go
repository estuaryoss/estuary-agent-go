package controllers

import (
	"estuary-agent-go/src/command"
	"estuary-agent-go/src/constants"
	u "estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)

type Profile struct {
	Name    string
	Hobbies []string
}

var CommandPost = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, _ := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}
	commands := u.TrimSpacesAndLineEndings(strings.Split(string(body), "\n"))
	cim := command.NewCommandInMemory()
	commandDescription := cim.RunCommands(commands)
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		commandDescription,
		r.URL.Path)

	u.ApiResponse(w, resp)
}
