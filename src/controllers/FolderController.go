package controllers

import (
	"archive/zip"
	"fmt"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var GetFolder = func(w http.ResponseWriter, r *http.Request) {
	folderName := r.Header.Get("Folder-Path")
	if folderName == "" {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.HTTP_HEADER_NOT_PROVIDED),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], folderName),
			fmt.Sprintf(u.GetMessage()[uint32(constants.HTTP_HEADER_NOT_PROVIDED)], folderName),
			r.URL.Path))
		return
	}

	zipFileName := "response.zip"
	err := zipFolder(folderName, zipFileName)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.FOLDER_ZIP_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.FOLDER_ZIP_FAILURE)], folderName),
			err.Error(),
			r.URL.Path))
		return
	}

	content, err := ioutil.ReadFile(zipFileName)
	if err != nil {
		u.ApiResponseError(w, u.ApiMessage(uint32(constants.GET_FILE_FAILURE),
			fmt.Sprintf(u.GetMessage()[uint32(constants.GET_FILE_FAILURE)], folderName),
			err.Error(),
			r.URL.Path))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(zipFileName))

	u.ApiResponseZip(w, content)
}

func zipFolder(sourceFolder, targetFileName string) error {
	zipfile, err := os.Create(targetFileName)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(sourceFolder)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(sourceFolder)
	}

	filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, sourceFolder))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)

		return err
	})

	return err
}
