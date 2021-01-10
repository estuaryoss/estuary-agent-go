package logging

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
)

const (
	HEADERS = "headers"
	BODY    = "body"
)

type MessageDumper struct {
	headers map[string]string
	body    map[string]string
}

func NewMessageDumper() *MessageDumper {
	return &MessageDumper{headers: make(map[string]string),
		body: make(map[string]string)}
}

func (md *MessageDumper) SetHeader(key string, value string) {
	md.headers[key] = value
}

func (md *MessageDumper) GetHeader(name string) string {
	return md.headers[name]
}

func (md *MessageDumper) GetHeaders() map[string]string {
	return md.headers
}

func (md *MessageDumper) DumpRequest(r *http.Request) map[string]interface{} {
	message := make(map[string]interface{})
	headers := md.getRequestHeaders(r)
	body, _ := httputil.DumpRequest(r, true)
	finalMessage := make(map[string]string)
	finalMessage["message"] = string(body)
	message[HEADERS] = headers
	message[BODY] = finalMessage

	return message
}

func (md *MessageDumper) DumpResponse(w http.ResponseWriter, body map[string]interface{}) map[string]interface{} {
	message := make(map[string]interface{})
	headers := md.getResponseHeaders(w)
	responseBody := body
	finalMessage := make(map[string]interface{})
	descriptionString, _ := json.Marshal(responseBody["description"])
	responseBody["description"] = string(descriptionString)
	finalMessage = responseBody

	message[HEADERS] = headers
	message[BODY] = finalMessage

	return message
}

func (md *MessageDumper) DumpResponseString(w http.ResponseWriter, body string) map[string]interface{} {
	message := make(map[string]interface{})
	headers := md.getResponseHeaders(w)
	finalMessage := make(map[string]interface{})

	finalMessage["message"] = body

	message[HEADERS] = headers
	message[BODY] = finalMessage

	return message
}

func (md *MessageDumper) DumpMessage(msg string) map[interface{}]interface{} {
	message := make(map[interface{}]interface{})
	finalMessage := make(map[string]string)
	finalMessage["message"] = msg
	message[HEADERS] = make(map[string]string)
	message[BODY] = finalMessage

	return message
}

func (md *MessageDumper) getRequestHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for name, values := range r.Header {
		for _, value := range values {
			headers[name] = value
		}
	}
	return headers
}

func (md *MessageDumper) getResponseHeaders(w http.ResponseWriter) map[string]string {
	headers := make(map[string]string)
	for name, values := range w.Header() {
		for _, value := range values {
			headers[name] = value
		}
	}
	return headers
}
