package controllers

import (
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var Ping = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		"pong",
		r.URL.Path)
	u.ApiResponse(w, resp)
}
