package routes

import (
	"fmt"
	"net/http"
	"strconv"

	h "github.com/estuaryoss/estuary-agent-go/src/handlers"

	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/controllers"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"github.com/estuaryoss/estuary-agent-go/src/services"
	"github.com/gorilla/mux"
)

var SetupServer = func(appPort string) {
	var router = mux.NewRouter()

	router.Use(h.AddXRequestIdHandler, h.LogHttpRequestHandler, h.TokenAuthenticationHandler)

	router.HandleFunc("/ping", controllers.Ping).Methods("GET")
	router.HandleFunc("/env", controllers.GetEnvVars).Methods("GET")
	router.HandleFunc("/env", controllers.SetEnvVars).Methods("POST")
	router.HandleFunc("/env", controllers.DeleteVirtualEnvVars).Methods("DELETE")
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

	//swagger ui
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
