package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/estuaryoss/estuary-agent-go/logging"
	"github.com/estuaryoss/estuary-agent-go/services"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/controllers"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var SetupServer = func(appPort string) {
	var router = mux.NewRouter()

	router.Use(AddXRequestId, LogHttpRequest, TokenAuthentication)

	router.HandleFunc("/ping", controllers.Ping).Methods("GET")
	router.HandleFunc("/env", controllers.GetEnvVars).Methods("GET")
	router.HandleFunc("/env", controllers.SetEnvVars).Methods("POST")
	router.HandleFunc("/env/{name}", controllers.GetEnvVar).Methods("GET")
	router.HandleFunc("/about", controllers.About).Methods("GET")
	router.HandleFunc("/info", controllers.About).Methods("GET")
	router.HandleFunc("/file", controllers.GetFile).Methods("GET")
	router.HandleFunc("/file", controllers.PutFile).Methods("POST", "PUT")
	router.HandleFunc("/folder", controllers.GetFolder).Methods("GET")
	router.HandleFunc("/command", controllers.CommandPost).Methods("POST")
	router.HandleFunc("/commandyaml", controllers.CommandPostYaml).Methods("POST")
	router.HandleFunc("/commandparallel", controllers.CommandParallelPost).Methods("POST")
	router.HandleFunc("/commanddetached", controllers.CommandDetachedGet).Methods("GET")
	router.HandleFunc("/commanddetached", controllers.CommandDetachedDelete).Methods("DELETE")
	router.HandleFunc("/commanddetached/{cid}", controllers.CommandDetachedPost).Methods("POST")
	router.HandleFunc("/commanddetachedyaml/{cid}", controllers.CommandDetachedPostYaml).Methods("POST")
	router.HandleFunc("/commanddetached/{cid}", controllers.CommandDetachedGetById).Methods("GET")
	router.HandleFunc("/commanddetached/{cid}", controllers.CommandDetachedDeleteById).Methods("DELETE")

	//swagger
	fs := http.FileServer(http.Dir("./swaggerui/"))
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	//eureka registration
	ec := services.NewEurekaClient()
	ec.RegisterApp(environment.GetInstance().GetConfigEnvVars()[constants.APP_IP_PORT])

	var err error
	isHttps, _ := strconv.ParseBool(environment.GetInstance().GetConfigEnvVars()[constants.HTTPS_ENABLE])
	if isHttps == true {
		err = http.ListenAndServeTLS(":"+environment.GetInstance().GetConfigEnvVars()[constants.PORT],
			environment.GetInstance().GetConfigEnvVars()[constants.HTTPS_CERT],
			environment.GetInstance().GetConfigEnvVars()[constants.HTTPS_KEY],
			router)
	} else {
		err = http.ListenAndServe(":"+environment.GetInstance().GetConfigEnvVars()[constants.PORT], router)
	}

	if err != nil {
		fmt.Println(err)
	}
}

var AddXRequestId = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uniqueId := uuid.New().String()
		xRequestIdRequest := r.Header.Get(constants.X_REQUEST_ID)
		if xRequestIdRequest != "" {
			w.Header().Add(constants.X_REQUEST_ID, xRequestIdRequest)
		} else {
			w.Header().Add(constants.X_REQUEST_ID, uniqueId)
			r.Header.Add(constants.X_REQUEST_ID, uniqueId)
		}
		next.ServeHTTP(w, r)
	})
}

var LogHttpRequest = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := logging.NewMessageDumper().DumpRequest(r)
		fluentdLogger := logging.GetInstance()
		result, _ := json.Marshal(fluentdLogger.Emit(constants.NAME+"."+"api",
			request, "DEBUG"))
		log.Println(string(result))

		next.ServeHTTP(w, r)
	})
}

var TokenAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Auth Token
		tokenHeader := r.Header.Get("Token")
		matchedUrl, _ := regexp.MatchString(`.*swaggerui.*`, r.URL.RequestURI())

		//permit swagerui even though is not auth
		if tokenHeader == environment.GetInstance().GetConfigEnvVars()[constants.HTTP_AUTH_TOKEN] || matchedUrl {
			// Delegate request to the given handle
			next.ServeHTTP(w, r)
		} else {
			response, _ := json.Marshal(u.ApiMessage(uint32(constants.UNAUTHORIZED),
				u.GetMessage()[uint32(constants.UNAUTHORIZED)],
				"Invalid Token",
				r.URL.Path))
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, string(response), http.StatusUnauthorized)
		}
	})
}

var LogHttpResponse = func(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
