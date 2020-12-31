package controllers

import (
	"estuary-agent-go/src/constants"
	u "estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var About = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		constants.About(),
		r.URL.Path)
	u.ApiResponse(w, resp)
}
