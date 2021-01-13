package environment

import (
	"os"
	"strings"
	"sync"

	"github.com/estuaryoss/estuary-agent-go/src/constants"
	"github.com/magiconair/properties"
)

var once sync.Once

type Env struct {
	configEnvVars  map[string]string
	env            map[string]string
	virtualEnv     map[string]string
	virtualEnvSize int
}

var singleton *Env

func GetInstance() *Env {
	once.Do(
		func() {
			singleton = &Env{
				configEnvVars:  map[string]string{},
				env:            GetEnvAsMap(),
				virtualEnv:     map[string]string{},
				virtualEnvSize: 100,
			}
			singleton.SetEnvVars(GetVirtualEnvAsMapFromFile())
			singleton.InitConfigEnvVars()
		})
	return singleton
}

func (env *Env) InitConfigEnvVars() map[string]string {
	initEnvVars := make(map[string]string)

	port := "8080"
	if env.GetEnv()[constants.PORT] != "" {
		port = env.GetEnv()[constants.PORT]
	}

	fluentdIpPort := ""
	if env.GetEnv()[constants.FLUENTD_IP_PORT] != "" {
		fluentdIpPort = env.GetEnv()[constants.FLUENTD_IP_PORT]
	}

	eurekaServer := ""
	if env.GetEnv()[constants.EUREKA_SERVER] != "" {
		eurekaServer = env.GetEnv()[constants.EUREKA_SERVER]
	}

	appIpPort := "localhost:8080"
	if env.GetEnv()[constants.APP_IP_PORT] != "" {
		appIpPort = env.GetEnv()[constants.APP_IP_PORT]
	}

	enableHttps := "false"
	if env.GetEnv()[constants.HTTPS_ENABLE] != "" {
		enableHttps = env.GetEnv()[constants.HTTPS_ENABLE]
	}

	httpsCert := "./https/cert.pem"
	if env.GetEnv()[constants.HTTPS_CERT] != "" {
		httpsCert = env.GetEnv()[constants.HTTPS_CERT]
	}

	httpsKey := "./https/key.pem"
	if env.GetEnv()[constants.HTTPS_KEY] != "" {
		enableHttps = env.GetEnv()[constants.HTTPS_KEY]
	}
	initEnvVars[constants.FLUENTD_IP_PORT] = fluentdIpPort
	initEnvVars[constants.PORT] = port
	initEnvVars[constants.EUREKA_SERVER] = eurekaServer
	initEnvVars[constants.APP_IP_PORT] = appIpPort
	initEnvVars[constants.HTTPS_ENABLE] = enableHttps
	initEnvVars[constants.HTTPS_CERT] = httpsCert
	initEnvVars[constants.HTTPS_KEY] = httpsKey

	env.configEnvVars = initEnvVars

	return initEnvVars
}

func (env *Env) GetConfigEnvVars() map[string]string {
	return env.configEnvVars
}

func (env *Env) GetEnv() map[string]string {
	return env.env
}

func (env *Env) GetVirtualEnv() map[string]string {
	return env.virtualEnv
}

func (env *Env) GetEnvAndVirtualEnv() map[string]string {
	return mergeMaps(env.env, env.virtualEnv)
}

func (env *Env) CleanVirtualEnv() {
	for envVarName := range env.virtualEnv {
		delete(env.virtualEnv, envVarName)
	}
}

func (env *Env) SetEnvVar(key string, value string) bool {
	if _, ok := env.env[key]; ok {
		return false
	}
	if _, ok := env.virtualEnv[key]; ok && len(env.virtualEnv) < env.virtualEnvSize && key != "" {
		env.virtualEnv[key] = value
		return true
	}

	if len(env.virtualEnv) < env.virtualEnvSize && key != "" {
		env.virtualEnv[key] = value
		return true
	}

	return false
}

func (env *Env) SetEnvVars(envVars map[string]string) map[string]string {
	addedEnvVars := make(map[string]string)
	for key, value := range envVars {
		if env.SetEnvVar(key, value) {
			addedEnvVars[key] = value
		}
	}

	return addedEnvVars
}

func (env *Env) GetEnvAndVirtualEnvArray() []string {
	environment := setItemInArray(env.GetEnvAndVirtualEnv())

	return environment
}

func GetEnvAsMap() map[string]string {
	environment := setItemInMap(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	return environment
}

func GetVirtualEnvAsMapFromFile() map[string]string {
	props, err := properties.LoadFile("environment.properties", properties.UTF8)

	if err != nil {
		return make(map[string]string)
	}
	return props.Map()
}

var setItemInArray = func(data map[string]string) []string {
	var items []string
	for k, v := range data {
		items = append(items, k+"="+v)
	}
	return items
}

var setItemInMap = func(data []string, getKeyVal func(item string) (key, val string)) map[string]string {
	items := make(map[string]string)
	for _, item := range data {
		key, val := getKeyVal(item)
		items[key] = val
	}
	return items
}

func mergeMaps(first map[string]string, second map[string]string) map[string]string {
	mergedMap := make(map[string]string)
	for k, v := range first {
		if k != "" {
			mergedMap[k] = v
		}
	}
	for k, v := range second {
		if k != "" {
			mergedMap[k] = v
		}
	}
	return mergedMap
}
