package src

import (
	"github.com/sakirsensoy/genv"
	"strconv"
)

var configVars = map[string]string{}

func SetConfig(key string, value string) {
	configVars[key] = value
}

func GetConfig(key string) string {
	if val, ok := configVars[key]; ok {
		return val
	}

	if val, ok := GetLuaValue(key); ok {
		return val
	}

	if val := genv.Key(key).String(); val != "" {
		return val
	}

	return ""
}

func GetConfigBool(key string) bool {
	var val string
	var ok bool

	if val, ok = configVars[key]; !ok {
		if val, ok = GetLuaValue(key); !ok {
			return genv.Key(key).Bool()
		}
	}

	num, _ := strconv.Atoi(val)
	return val == "true" || num > 0
}
