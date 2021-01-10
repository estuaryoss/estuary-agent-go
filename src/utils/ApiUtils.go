package utils

import (
	"encoding/json"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"io"
	"net/http"
	"time"
)

type HttpResponse struct {
	Response interface{}
	Code     int
}

var ServerHttpResponse *HttpResponse = &HttpResponse{}

func ApiMessage(code uint32, message string, description interface{}, path string) map[string]interface{} {
	t := time.Now()
	timestamp := GetFormattedTimeAsString(t)

	return map[string]interface{}{
		"code":        code,
		"message":     message,
		"description": description,
		"name":        constants.NAME,
		"version":     constants.VERSION,
		"timestamp":   timestamp,
		"path":        path,
	}
}

func ApiResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
	ServerHttpResponse = &HttpResponse{Response: data, Code: http.StatusOK}
}

func ApiResponseError(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(data)
	ServerHttpResponse = &HttpResponse{Response: data, Code: http.StatusInternalServerError}
}

func ApiResponseByteArray(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "application/octet-stream")
	io.Writer.Write(w, data)
	ServerHttpResponse = &HttpResponse{Response: string(data), Code: http.StatusOK}
}

func ApiResponseZip(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "application/zip")
	io.Writer.Write(w, data)
	ServerHttpResponse = &HttpResponse{Response: "skipped (zip archive content)", Code: http.StatusOK}
}
