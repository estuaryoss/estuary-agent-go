package controllers

import (
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"net/http"
)

var Ping = func(w http.ResponseWriter, r *http.Request) {
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		"pong",
		r.URL.Path)
	u.ApiResponse(w, resp)
}
