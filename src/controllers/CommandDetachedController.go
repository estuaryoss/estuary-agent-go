package controllers

import (
	"fmt"
	"github.com/dinuta/estuary-agent-go/src/constants"
	u "github.com/dinuta/estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
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
	cmdBackground := []string{"./runcmd", "--cid=" + cmd_id, "--args=" + strings.Join(commands, ";;")}
	log.Print(fmt.Sprintf("Starting command '%s' in background", strings.Join(cmdBackground, " ")))
	err := u.StartCommandAndGetError(cmdBackground)
	if err != nil {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.COMMAND_DETACHED_START_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.COMMAND_DETACHED_START_FAILURE)], cmd_id),
			err.Error(),
			r.URL.Path))
		return
	}
	resp := u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		cmd_id,
		r.URL.Path)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	u.ApiResponse(w, resp)
}
