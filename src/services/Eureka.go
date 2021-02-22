package services

import (
	"fmt"
	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"log"
	"strconv"
	"strings"
)

type Eureka struct {
	client *eureka.Client
}

func NewEurekaClient() *Eureka {
	return &Eureka{client: GetEurekaClient()}
}

func (e *Eureka) RegisterApp(appIpPort string) {
	if e.client == nil {
		return
	}
	appIpPortArray := strings.Split(appIpPort, ":")
	appIp := appIpPortArray[0]
	appPort, _ := strconv.Atoi(appIpPortArray[1])
	isSsl, _ := strconv.ParseBool(environment.GetInstance().GetConfigEnvVars()[constants.HTTPS_ENABLE])
	instanceInfo := eureka.NewInstanceInfo(appIp, constants.NAME,
		appIp, appPort, 30, isSsl) //Create a new instanceInfo to register

	instanceInfo.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	var protocol = "http"
	if isSsl {
		protocol = "https"
	}
	instanceInfo.HomePageUrl = protocol + "://" + appIp + ":" + strconv.Itoa(appPort) + "/"
	instanceInfo.HealthCheckUrl = protocol + "://" + appIp + ":" + strconv.Itoa(appPort) + "/ping"
	instanceInfo.StatusPageUrl = protocol + "://" + appIp + ":" + strconv.Itoa(appPort) + "/ping"

	err := e.client.RegisterInstance(constants.NAME, instanceInfo)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to register to EurekaServer: %s with ip: %s and port: %d",
			fmt.Sprint(e.client.GetCluster()), appIp, appPort))
	}
}

func GetEurekaClient() *eureka.Client {
	eurekaServer := environment.GetInstance().GetEnv()[constants.EUREKA_SERVER]
	if eurekaServer != "" {
		client := eureka.NewClient([]string{eurekaServer})

		return client
	}

	return nil
}
