package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/models"
	"github.com/estuaryoss/estuary-agent-go/src/state"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var CommandDetachedPost = func(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	cmdId := params["cid"]
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

var CommandDetachedGetById = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cmdId := params["cid"]
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

var CommandDetachedGet = func(w http.ResponseWriter, r *http.Request) {
	jsonFileName := state.GetLastCommand()
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
	cmdId := cd.GetId()
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
