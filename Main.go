package main

import (
	"github.com/dinuta/estuary-agent-go/src/routes"
	"log"
	"os"
)

func main() {
	var appPort = os.Getenv("PORT")

	if appPort == "" {
		appPort = "8080"
	}

	log.Printf("Running on port %s\n", appPort)
	routes.SetupServer(appPort)
}
