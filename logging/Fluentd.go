package logging

import (
	"fmt"
	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/estuaryoss/estuary-agent-go/src/environment"
	u "github.com/estuaryoss/estuary-agent-go/src/utils"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/shirou/gopsutil/v3/host"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var once sync.Once

type Fluentd struct {
	logger *fluent.Fluent
}

var singleton *Fluentd

func GetInstance() *Fluentd {
	once.Do(
		func() {
			singleton = &Fluentd{
				logger: GetFluentdLogger(),
			}
		})
	return singleton
}

func (f *Fluentd) enrichLog(levelCode string, msg interface{}) map[string]interface{} {
	plat, fam, ver, _ := host.PlatformInformation()
	port := 8080
	if environment.GetInstance().GetEnv()[constants.PORT] != "" {
		port, _ = strconv.Atoi(environment.GetInstance().GetEnv()[constants.PORT])
	}
	constants.About()
	enrichedLog := make(map[string]interface{})
	enrichedLog["name"] = constants.NAME
	enrichedLog["port"] = port
	enrichedLog["version"] = constants.VERSION
	enrichedLog["uname"] = []string{plat, fam, ver}
	enrichedLog["golang"] = runtime.Version()
	enrichedLog["pid"] = os.Getpid()
	enrichedLog["level_code"] = levelCode
	enrichedLog["msg"] = msg
	enrichedLog["timestamp"] = u.GetFormattedTimeAsString(time.Now())

	return enrichedLog
}

func (f *Fluentd) Emit(tag string, msg map[string]interface{}, level string) interface{} {
	consoleMessage := make(map[string]interface{})
	var isMessageSent = true
	enrichedLog := f.enrichLog(level, msg)
	err := f.send(tag, enrichedLog)
	if err != nil {
		isMessageSent = false
	}
	consoleMessage["message"] = enrichedLog
	consoleMessage["emit"] = isMessageSent

	return consoleMessage
}

func (f *Fluentd) send(tag string, msg map[string]interface{}) error {
	if environment.GetInstance().GetEnv()[constants.FLUENTD_IP_PORT] != "" {
		return f.logger.Post(tag, msg)
	}
	return nil
}

func GetFluentdLogger() *fluent.Fluent {
	fluentdIpPort := environment.GetInstance().GetEnv()[constants.FLUENTD_IP_PORT]
	if fluentdIpPort != "" {
		fluentdIpPortArray := strings.Split(fluentdIpPort, ":")
		fluentdIp := fluentdIpPortArray[0]
		fluentdPort, err := strconv.Atoi(fluentdIpPortArray[1])
		if err != nil {
			log.Printf(fmt.Sprintf("Unable to parse port %s to int", fluentdIpPortArray[1]))
		}

		fluent, err := fluent.New(fluent.Config{FluentHost: fluentdIp, FluentPort: fluentdPort})
		if err != nil {
			log.Printf(fmt.Sprintf("Unable to create logger for host: %s and port: %d",
				fluent.Config.FluentHost, fluent.Config.FluentPort))

			return fluent
		}
	}

	return nil
}
