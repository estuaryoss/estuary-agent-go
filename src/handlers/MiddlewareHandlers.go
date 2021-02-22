package handlers

import (
	"encoding/json"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"github.com/estuaryoss/estuary-agent-go/src/logging"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
)

var AddXRequestIdHandler = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestIdRequest := r.Header.Get(constants.X_REQUEST_ID)
		if xRequestIdRequest != "" {
			w.Header().Add(constants.X_REQUEST_ID, xRequestIdRequest)
		} else {
			uniqueId := uuid.New().String()
			w.Header().Add(constants.X_REQUEST_ID, uniqueId)
			r.Header.Add(constants.X_REQUEST_ID, uniqueId)
		}

		next.ServeHTTP(w, r)
	})
}

var LogHttpRequestHandler = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := logging.NewMessageDumper().DumpRequest(r)
		fluentdLogger := logging.GetFluentdInstance()
		result, _ := json.Marshal(fluentdLogger.Emit(constants.NAME+"."+"api",
			request, "DEBUG"))
		log.Println(string(result))

		next.ServeHTTP(w, r)
	})
}

var TokenAuthenticationHandler = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Auth Token
		tokenHeader := r.Header.Get("Token")
		matchedUrl, _ := regexp.MatchString(`.*swaggerui.*`, r.URL.RequestURI())

		//permit swagger ui even though is not auth
		if (tokenHeader == environment.GetInstance().GetConfigEnvVars()[constants.HTTP_AUTH_TOKEN]) || matchedUrl {
			// Delegate request to the given handle
			next.ServeHTTP(w, r)
		} else {
			resp := u.ApiMessage(uint32(constants.UNAUTHORIZED),
				u.GetMessage()[uint32(constants.UNAUTHORIZED)],
				"Invalid Token",
				r.URL.Path)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			u.ApiResponseError(w, resp)
		}

		LogHttpResponse(w)
	})
}

func LogHttpResponse(w http.ResponseWriter) {
	var response interface{}
	if typeof(u.ServerHttpResponse.Response) == "map" {
		response = logging.NewMessageDumper().DumpResponse(w, u.ServerHttpResponse.Response.(map[string]interface{}))
	} else {
		response = logging.NewMessageDumper().DumpResponseString(w, u.ServerHttpResponse.Response.(string))
	}

	fluentdLogger := logging.GetFluentdInstance()
	result, _ := json.Marshal(fluentdLogger.Emit(constants.NAME+"."+"api",
		response, "DEBUG"))
	log.Println(string(result))
}

func typeof(v interface{}) string {
	switch v.(type) {
	case map[string]interface{}:
		return "map"
	case string:
		return "string"
	default:
		return "string"
	}
}
