package main

import (
	"encoding/json"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"github.com/estuaryoss/estuary-agent-go/src/routes"
	"log"
)

func main() {
	appPort := environment.GetInstance().GetConfigEnvVars()[constants.PORT]

	log.Printf("Running on port: %s\n\n", appPort)
	configEnvVars, _ := json.Marshal(environment.GetInstance().GetConfigEnvVars())
	log.Printf("Config env vars: %s\n\n", configEnvVars)
	systemEnvVars, _ := json.Marshal(environment.GetInstance().GetEnvAndVirtualEnv())
	log.Printf("Environment: %s\n\n", systemEnvVars)
	userEnvVars, _ := json.Marshal(environment.GetInstance().GetVirtualEnv())
	log.Printf("User(virtual) Environment: %s\n\n", userEnvVars)
	routes.SetupServer(appPort)
}
