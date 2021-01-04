package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

var GetEnvVar = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	env := environment.GetInstance()
	envVar := env.GetEnvAndVirtualEnv()[ps.ByName("name")]

	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		envVar,
		r.URL.Path)
	u.ApiResponse(w, resp)
}

var GetEnvVars = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	env := environment.GetInstance()
	envVars := env.GetEnvAndVirtualEnv()

	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		envVars,
		r.URL.Path)

	u.ApiResponse(w, resp)
}

var SetEnvVars = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	env := environment.GetInstance()
	body, err := ioutil.ReadAll(r.Body)
	attemptedEnvVars := make(map[string]string)
	if err != nil {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.SET_ENV_VAR_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.SET_ENV_VAR_FAILURE)], string(body)),
			err.Error(),
			r.URL.Path))
		return
	}

	err = json.Unmarshal(body, &attemptedEnvVars)
	if err != nil {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.SET_ENV_VAR_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.SET_ENV_VAR_FAILURE)], string(body)),
			err.Error(),
			r.URL.Path))
		return
	}
	addedEnvVars := env.SetEnvVars(attemptedEnvVars)

	u.ApiResponse(w, u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		addedEnvVars,
		r.URL.Path))
}
