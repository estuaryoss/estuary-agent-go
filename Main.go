package main

import (
	"log"

	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"github.com/estuaryoss/estuary-agent-go/src/routes"
)

func main() {
	appPort := environment.GetInstance().GetConfigEnvVars()[constants.PORT]

	log.Printf("Running on port %s\n", appPort)
	routes.SetupServer(appPort)
}
