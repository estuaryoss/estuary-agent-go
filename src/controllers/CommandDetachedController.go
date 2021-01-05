package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/models"
	"github.com/estuaryoss/estuary-agent-go/src/state"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
)

var CommandDetachedPost = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, _ := ioutil.ReadAll(r.Body)
	cmdId := ps.ByName("cid")
	commands := u.TrimSpacesAndLineEndings(strings.Split(string(body), "\n"))
	if len(commands) == 0 {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}
	cmdBackground := []string{"./runcmd", "--cid=" + cmdId, "--args=" + strings.Join(commands, ";;")}
	log.Print(fmt.Sprintf("Starting command '%s' in background", strings.Join(cmdBackground, " ")))
	err := u.StartCommandAndGetError(cmdBackground)
	if err != nil {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.COMMAND_DETACHED_START_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.COMMAND_DETACHED_START_FAILURE)], cmdId),
			err.Error(),
			r.URL.Path))
		return
	}
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		cmdId,
		r.URL.Path)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	state.SetLastCommand(cmdId)
	u.ApiResponse(w, resp)
}

var CommandDetachedGetId = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cmdId := ps.ByName("cid")
	jsonFileName := fmt.Sprintf(constants.CMD_BACKGROUND_JSON_OUTPUT, cmdId)
	var cd *models.CommandDescription

	err := json.Unmarshal(u.ReadFile(jsonFileName), cd)
	if err != nil {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.GET_COMMAND_DETACHED_INFO_FAILURE),
			u.GetMessage()[uint32(constants.GET_COMMAND_DETACHED_INFO_FAILURE)],
			err.Error(),
			r.URL.Path))
		return
	}
	commands := cd.GetCommands()
	for cmd, _ := range commands {
		cmdDetails := commands[cmd].GetCommandDetails()
		cmdDetails.SetOut(string(u.ReadFile(u.GetBase64HashForTheCommand(cmd, cmdId, ".out"))))
		cmdDetails.SetErr(string(u.ReadFile(u.GetBase64HashForTheCommand(cmd, cmdId, ".err"))))
	}

	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		cd,
		r.URL.Path)

	u.ApiResponse(w, resp)
}
