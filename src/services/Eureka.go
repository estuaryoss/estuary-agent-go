package services

import (
	"fmt"
	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
	appPort, _ := strconv.Atoi(strings.Split(appIpPortArray[1], "/")[0])
	preUrl := strings.Join(strings.Split(appIpPortArray[1], "/")[1:], "/")

	isSsl, _ := strconv.ParseBool(environment.GetInstance().GetConfigEnvVars()[constants.HTTPS_ENABLE])
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "NA"
	}

	instanceInfo := eureka.NewInstanceInfo(hostName+":"+strconv.Itoa(appPort), constants.NAME,
		appIp, appPort, 30, isSsl) //Create a new instanceInfo to register

	instanceInfo.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}

	var protocol = "http"
	if isSsl {
		protocol = "https"
	}

	instanceInfo.Metadata.Map["management.port"] = strconv.Itoa(appPort)
	instanceInfo.InstanceID = hostName + ":" + constants.NAME + ":" + strconv.Itoa(appPort)
	instanceInfo.HomePageUrl = protocol + "://" + appIp + ":" + strconv.Itoa(appPort) + "/"
	instanceInfo.HealthCheckUrl = protocol + "://" + appIp + ":" + strconv.Itoa(appPort) + "/" + preUrl + "/ping"
	instanceInfo.StatusPageUrl = protocol + "://" + appIp + ":" + strconv.Itoa(appPort) + "/" + preUrl + "/ping"

	err = e.client.RegisterInstance(constants.NAME, instanceInfo)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to register to EurekaServer: %s with ip: %s and port: %d",
			fmt.Sprint(e.client.GetCluster()), appIp, appPort))
	}

	go sendHeartbeat(e.client, instanceInfo, 30)
}

func GetEurekaClient() *eureka.Client {
	eurekaServer := environment.GetInstance().GetEnv()[constants.EUREKA_SERVER]
	if eurekaServer != "" {
		client := eureka.NewClient([]string{eurekaServer})

		return client
	}

	return nil
}

func sendHeartbeat(client *eureka.Client, instanceInfo *eureka.InstanceInfo, second int) {
	d := (time.Duration(second) * time.Second)
	tick := time.Tick(d)
	for {
		select {
		case <-tick:
			log.Printf("Sending eureka heartbeat >  [App: '%s', InstanceId: '%s']", instanceInfo.App, instanceInfo.InstanceID)
			client.SendHeartbeat(instanceInfo.App, instanceInfo.InstanceID)
		}
	}
}
