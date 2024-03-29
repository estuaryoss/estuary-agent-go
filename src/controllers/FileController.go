package controllers

import (
	"fmt"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

var GetFile = func(w http.ResponseWriter, r *http.Request) {
	fileHeaderName := "File-Path"
	fileName := r.Header.Get(fileHeaderName)
	if fileName == "" {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.HTTP_HEADER_NOT_PROVIDED),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], fileHeaderName),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], fileHeaderName),
			r.URL.Path))
		return
	}

	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.GET_FILE_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.GET_FILE_FAILURE)], fileName),
			err.Error(),
			r.URL.Path))
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fileName))

	u.ApiResponseByteArray(w, content)
}

var PutFile = func(w http.ResponseWriter, r *http.Request) {
	fileHeaderName := "File-Path"
	fileName := r.Header.Get(fileHeaderName)
	if fileName == "" {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.HTTP_HEADER_NOT_PROVIDED),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], fileHeaderName),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], fileHeaderName),
			r.URL.Path))
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.EMPTY_REQUEST_BODY_PROVIDED),
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			u.GetMessage()[uint32(constants.EMPTY_REQUEST_BODY_PROVIDED)],
			r.URL.Path))
		return
	}

	err := ioutil.WriteFile(fileName, body, 0644)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.UPLOAD_FILE_FAILURE),
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
