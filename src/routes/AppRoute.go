package routes

import (
	"encoding/json"
	"fmt"
	"github.com/estuaryoss/estuary-agent-go/src/controllers"
	"github.com/gorilla/mux"
	"net/http"
	"os"

	"github.com/estuaryoss/estuary-agent-go/src/constants"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
)

var SetupServer = func(appPort string) {
	var router = mux.NewRouter()

	router.Use(TokenAuthentication)

	router.HandleFunc("/ping", controllers.Ping).Methods("GET")
	router.HandleFunc("/env", controllers.GetEnvVars).Methods("GET")
	router.HandleFunc("/env", controllers.SetEnvVars).Methods("POST")
	router.HandleFunc("/env/{name}", controllers.GetEnvVar).Methods("GET")
	router.HandleFunc("/about", controllers.About).Methods("GET")
	router.HandleFunc("/file", controllers.GetFile).Methods("GET")
	router.HandleFunc("/file", controllers.PutFile).Methods("POST")
	router.HandleFunc("/folder", controllers.GetFolder).Methods("GET")
	router.HandleFunc("/command", controllers.CommandPost).Methods("POST")
	router.HandleFunc("/commandparallel", controllers.CommandParallelPost).Methods("POST")
	router.HandleFunc("/commanddetached", controllers.CommandDetachedGet).Methods("GET")
	router.HandleFunc("/commanddetached/{cid}", controllers.CommandDetachedPost).Methods("POST")
	router.HandleFunc("/commanddetached/{cid}", controllers.CommandDetachedGetById).Methods("GET")

	//swagger
	fs := http.FileServer(http.Dir("./swaggerui/"))
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	err := http.ListenAndServe(":"+appPort, router)
	if err != nil {
		fmt.Print(err)
	}
}

var TokenAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Auth Token
		tokenHeader := r.Header.Get("Token")

		if tokenHeader == os.Getenv("HTTP_AUTH_TOKEN") {
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
