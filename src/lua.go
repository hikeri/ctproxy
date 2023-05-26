package src

import (
	"errors"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"log"
	"os"
	"strings"
)

const LuaFile = "script.lua"

var luaScript lua.LState
var luaConfigAvailable bool
var luaValues lua.LValue

var LuaDomainRoutes map[string]string
var LuaAllowedPorts []string
var LuaUserPorts map[string][]string

func LoadLuaConfig() {
	if _, err := os.Stat("./" + LuaFile); errors.Is(err, os.ErrNotExist) {
		return
	}

	luaScript = *lua.NewState()
	defer luaScript.Close()
	if err := luaScript.DoFile(LuaFile); err != nil {
		log.Fatalln("Lua error", err)
	}

	luaValues = luaScript.GetGlobal("conf")
	luaConfigAvailable = luaValues.Type().String() == "table"

	allowedPortsStr := luaScript.GetGlobal("allowedPorts").String()
	LuaAllowedPorts = strings.Split(allowedPortsStr, ",")

	if err := gluamapper.Map(luaScript.GetGlobal("proxy").(*lua.LTable), &LuaDomainRoutes); err != nil {
		log.Println("Cannot import proxy map to Go")
	}

	var ports map[string]string
	if err := gluamapper.Map(luaScript.GetGlobal("userPorts").(*lua.LTable), &ports); err != nil {
		log.Println("Cannot import user ports to Go")
		for user, portsString := range ports {
			LuaUserPorts[user] = strings.Split(portsString, ",")
		}
	}
}

func GetLuaValue(key string) (string, bool) {
	if !luaConfigAvailable {
		return "", false
	}

	exists := luaScript.GetField(luaValues, key).Type().String() != "nil"
	if !exists {
		return "", false
	}

	return luaScript.GetField(luaValues, key).String(), true
}
