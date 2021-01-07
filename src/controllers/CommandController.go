package controllers

import (
	"github.com/estuaryoss/estuary-agent-go/src/command"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"github.com/estuaryoss/estuary-agent-go/src/models"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"strings"
)

var CommandPost = func(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	commands := u.TrimSpacesAndLineEndings(strings.Split(string(body), "\n"))

	if len(commands) == 0 {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
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

var CommandPostYaml = func(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	commands := u.TrimSpacesAndLineEndings(strings.Split(string(body), "\n"))

	if len(commands) == 0 {
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
	envVarsSet := environment.GetInstance().SetEnvVars(envVars)

	cim := command.NewCommandInMemory()
	commandDescription := cim.RunCommands(configParser.GetCommandsList(yamlConfig))
	yamlConfig.SetEnv(envVarsSet)
	models.SetDescription(yamlConfig, commandDescription)
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		models.GetDescription(),
		r.URL.Path)

	u.ApiResponse(w, resp)
}
