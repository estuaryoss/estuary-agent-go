package environment

import (
	"github.com/magiconair/properties"
	"os"
	"strings"
	"sync"
)

var once sync.Once

type Env struct {
	env            map[string]string
	virtualEnv     map[string]string
	virtualEnvSize int
}

var singleton *Env

func GetInstance() *Env {
	once.Do(
		func() {
			singleton = &Env{
				env:            GetEnvAsMap(),
				virtualEnv:     map[string]string{},
				virtualEnvSize: 100,
			}
			singleton.SetEnvVars(GetVirtualEnvAsMapFromFile())
		})
	return singleton
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
	for k, v := range second {
		first[k] = v
	}

	return first
}
