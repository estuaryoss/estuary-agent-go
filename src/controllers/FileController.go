package controllers

import (
	"fmt"
	"github.com/dinuta/estuary-agent-go/src/constants"
	u "github.com/dinuta/estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

var GetFile = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fileName := r.Header.Get("File-Path")
	if fileName == "" {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.HTTP_HEADER_NOT_PROVIDED),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], fileName),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], fileName),
			r.URL.Path))
		return
	}

	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.GET_FILE_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.GET_FILE_FAILURE)], fileName),
			err.Error(),
			r.URL.Path))
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fileName))

	u.ApiResponseByteArray(w, content)
}

var PutFile = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fileName := r.Header.Get("File-Path")
	if fileName == "" {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.HTTP_HEADER_NOT_PROVIDED),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], fileName),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], fileName),
			r.URL.Path))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}

	err := ioutil.WriteFile(fileName, body, 0644)
	if err != nil {
		u.ApiResponse(w, u.ApiMessage(uint32(constants.UPLOAD_FILE_FAILURE),
			u.GetMessage()[uint32(constants.UPLOAD_FILE_FAILURE)],
			err.Error(),
			r.URL.Path))
		return
	}

	u.ApiResponse(w, u.ApiMessage(uint32(constants.SUCCESS),
		u.GetMessage()[uint32(constants.SUCCESS)],
		u.GetMessage()[uint32(constants.SUCCESS)],
		r.URL.Path))
}
