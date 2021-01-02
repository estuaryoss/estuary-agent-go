package controllers

import (
	"fmt"
	"github.com/dinuta/estuary-agent-go/src/constants"
	u "github.com/dinuta/estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strings"
)

var CommandDetachedPost = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, _ := ioutil.ReadAll(r.Body)
	cmd_id := ps.ByName("cid")
	commands := u.TrimSpacesAndLineEndings(strings.Split(string(body), "\n"))
	if len(commands) == 0 {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}

	u.StartCommand(fmt.Sprint("run --cid=%s --args=%s", cmd_id, strings.Join(commands, ";;")))

	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		cmd_id,
		r.URL.Path)
	w.WriteHeader(http.StatusAccepted)
	u.ApiResponse(w, resp)
}
