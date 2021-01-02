package routes

import (
	"encoding/json"
	"fmt"
	"github.com/dinuta/estuary-agent-go/src/constants"
	"github.com/dinuta/estuary-agent-go/src/controllers"
	u "github.com/dinuta/estuary-agent-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
)

var SetupServer = func(appPort string) {
	var router = httprouter.New()
	router.GET("/ping", TokenAuthentication(controllers.Ping))
	router.GET("/env", TokenAuthentication(controllers.GetEnvVars))
	router.POST("/env", TokenAuthentication(controllers.SetEnvVars))
	router.GET("/env/:name", TokenAuthentication(controllers.GetEnvVar))
	router.GET("/about", TokenAuthentication(controllers.About))
	router.GET("/file", TokenAuthentication(controllers.GetFile))
	router.PUT("/file", TokenAuthentication(controllers.PutFile))
	router.GET("/folder", TokenAuthentication(controllers.GetFolder))
	router.POST("/command", TokenAuthentication(controllers.CommandPost))
	router.POST("/commandparallel", TokenAuthentication(controllers.CommandParallelPost))
	router.POST("/commanddetached", TokenAuthentication(controllers.CommandDetachedPost))

	err := http.ListenAndServe(":"+appPort, router)
	if err != nil {
		fmt.Print(err)
	}
}

var TokenAuthentication = func(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Auth Token
		tokenHeader := r.Header.Get("Token")

		if tokenHeader == os.Getenv("HTTP_AUTH_TOKEN") {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			response, _ := json.Marshal(u.ApiMessage(uint32(constants.UNAUTHORIZED),
				u.GetMessage()[uint32(constants.UNAUTHORIZED)],
				"Invalid Token",
				r.URL.Path))
			w.Header().Add("Content-Type", "application/json")
			http.Error(w, string(response), http.StatusUnauthorized)
		}
	}
}
