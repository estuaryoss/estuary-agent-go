package controllers

import (
	"net/http"

	"github.com/estuaryoss/estuary-agent-go/src/constants"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/gorilla/mux"
)

var GetProcessesByExecName = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	processes := u.GetAllProcessesByExecName(params["name"])

	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		processes,
		r.URL.Path)
	u.ApiResponse(w, resp)
}

var GetProcesses = func(w http.ResponseWriter, r *http.Request) {
	processes := u.GetAllProcesses()

	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		processes,
		r.URL.Path)

	u.ApiResponse(w, resp)
}
