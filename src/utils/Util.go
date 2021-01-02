package utils

import (
	"encoding/json"
	"fmt"
	"github.com/dinuta/estuary-agent-go/src/constants"
	"io"
	"net/http"
	"time"
)

func ApiMessage(code uint32, message string, description interface{}, path string) map[string]interface{} {
	t := time.Now()
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000)

	return map[string]interface{}{
		"code":        code,
		"message":     message,
		"description": description,
		"name":        constants.Name,
		"version":     constants.Version,
		"timestamp":   timestamp,
		"path":        path,
	}
}

func ApiResponse(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ApiResponseByteArray(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "application/octet-stream")
	io.Writer.Write(w, data)
}

func ApiResponseZip(w http.ResponseWriter, data []byte) {
	w.Header().Add("Content-Type", "application/zip")
	io.Writer.Write(w, data)
}
