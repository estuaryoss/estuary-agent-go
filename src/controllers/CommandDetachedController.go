package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"github.com/estuaryoss/estuary-agent-go/src/models"
	"github.com/estuaryoss/estuary-agent-go/src/state"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var CommandDetachedPost = func(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	cmdId := params["cid"]
	commands := u.TrimSpacesAndLineEndings(
		strings.Split(strings.Trim(string(body), " "), "\n"))
	if len(commands) == 0 {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}
	cmdBackground := []string{"./runcmd",
		"--cid=" + cmdId, "--args=" + strings.Join(commands, ";;"), "--enableStreams=true"}
	log.Print(fmt.Sprintf("Starting command '%s' in background", strings.Join(cmdBackground, " ")))
	cmd := u.GetCommand(cmdBackground)
	err := cmd.Start()
	state.GetInstance().AddCmdToCommandList(cmdId, cmd)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.COMMAND_DETACHED_START_FAILURE),
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

var CommandDetachedPostYaml = func(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	cmdId := params["cid"]
	if len(strings.Trim(string(body), " ")) == 0 {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}

	var yamlConfig = models.NewYamlConfig()
	var configParser = u.NewYamlConfigParser()
	err := yaml.Unmarshal(body, &yamlConfig)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.INVALID_YAML_CONFIG),
			u.GetMessage()[uint32(constants.INVALID_YAML_CONFIG)],
			err.Error(),
			r.URL.Path))
		return
	}
	err = configParser.CheckConfig(yamlConfig)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.INVALID_YAML_CONFIG),
			u.GetMessage()[uint32(constants.INVALID_YAML_CONFIG)],
			err.Error(),
			r.URL.Path))
		return
	}
	envVars := yamlConfig.GetEnv()
	environment.GetInstance().SetEnvVars(envVars)

	cmdBackground := []string{"./runcmd",
		"--cid=" + cmdId, "--args=" + strings.Join(configParser.GetCommandsList(yamlConfig), ";;"), "--enableStreams=true"}
	log.Print(fmt.Sprintf("Starting command '%s' in background", strings.Join(cmdBackground, " ")))
	cmd := u.GetCommand(cmdBackground)
	err = cmd.Start()
	state.GetInstance().AddCmdToCommandList(cmdId, cmd)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.COMMAND_DETACHED_START_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.COMMAND_DETACHED_START_FAILURE)], cmdId),
			err.Error(),
			r.URL.Path))
		return
	}
	commandDescription := cmdId
	models.SetDescription(yamlConfig, commandDescription)
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		models.GetDescription(),
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
	cd := models.NewCommandDescription()

	err := json.Unmarshal(u.ReadFile(jsonFileName), cd)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.GET_COMMAND_DETACHED_INFO_FAILURE),
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
	cd := models.NewCommandDescription()

	if !u.DoesFileExists(jsonFileName) {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.GET_COMMAND_DETACHED_INFO_FAILURE),
			u.GetMessage()[uint32(constants.GET_COMMAND_DETACHED_INFO_FAILURE)],
			fmt.Sprintf("File %s does not exists", jsonFileName),
			r.URL.Path))
		return
	}

	err := json.Unmarshal(u.ReadFile(jsonFileName), cd)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.GET_COMMAND_DETACHED_INFO_FAILURE),
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

var CommandDetachedDelete = func(w http.ResponseWriter, r *http.Request) {
	u.KillAllProcesses()

	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		u.GetMessage()[uint32(constants.SUCCESS)],
		r.URL.Path)

	u.ApiResponse(w, resp)
}

var CommandDetachedDeleteById = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cmdId := params["cid"]
	if state.GetInstance().GetBackgroundCommandList()[cmdId] == nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.COMMAND_DETACHED_STOP_FAILURE),
			u.GetMessage()[uint32(constants.COMMAND_DETACHED_STOP_FAILURE)],
			errors.New(fmt.Sprintf("Exception: There is no active process for command %s", cmdId)).Error(),
			r.URL.Path))
		return
	}
	u.KillProcess(cmdId)

	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		u.GetMessage()[uint32(constants.SUCCESS)],
		r.URL.Path)

	u.ApiResponse(w, resp)
}
