package luar

import (
	"fmt"
	"github.com/camsiabor/golua/lua"
	"path/filepath"
	"runtime"
	"strings"
)

const LUA_PATH = "LUA_PATH"
const LUA_CPATH = "LUA_CPATH"
const LUA_PATH_ABS = "LUA_PATH_ABS"
const LUA_CPATH_ABS = "LUA_CPATH_ABS"

func GetLuaPath(luaPath, luaCPath string) (luaPathAbs, luaCPathAbs, luaPathFull, luaCPathFull string) {

	var err error
	var luaVersion = lua.GetVersionNumber()
	var luaVersionWithoutDot = lua.GetVersionNumberWithoutDot()

	if luaPath, err = filepath.Abs(luaPath); err != nil {
		panic(err)
	}
	if luaCPath, err = filepath.Abs(luaCPath); err != nil {
		panic(err)
	}

	var luaLibSuffix = "so"
	if runtime.GOOS == "windows" {
		luaPath = strings.Replace(luaPath, "\\", "/", -1)
		luaCPath = strings.Replace(luaCPath, "\\", "/", -1)
		luaLibSuffix = "dll"
	}

	if luaPath[:len(luaPath)-1] != "/" {
		luaPath = luaPath + "/"
	}

	if luaCPath[:len(luaCPath)-1] != "/" {
		luaCPath = luaCPath + "/"
	}

	luaPathFull = luaPath + "?.lua;" +
		luaPath + "?init.lua;" +
		luaPath + "?;" +
		luaPath + "lib/?.lua;" +
		luaPath + "lib/?init.lua;" +
		luaPath + "lib/?;"

	luaCPathFull = luaCPath + "lib/?." + luaLibSuffix + ";" +
		luaCPath + "lib/?" + luaVersion + "." + luaLibSuffix + ";" +
		luaCPath + "lib/?" + luaVersionWithoutDot + "." + luaLibSuffix + ";" +
		luaCPath + "lib/load.all." + luaLibSuffix + ";" +
		luaCPath + "lib/?;"

	return luaPath, luaCPath, luaPathFull, luaCPathFull
}

func SetLuaPath(L *lua.State, luaPath, luaCPath string) (err error) {

	var luaPathAbs, luaCPathAbs, luaPathFull, luaCPathFull = GetLuaPath(luaPath, luaCPath)

	L.PushString(luaPathFull)
	L.SetGlobal(LUA_PATH)

	L.PushString(luaCPathFull)
	L.SetGlobal(LUA_CPATH)

	L.PushString(luaPathAbs)
	L.SetGlobal(LUA_PATH_ABS)

	L.PushString(luaCPathAbs)
	L.SetGlobal(LUA_CPATH_ABS)

	L.GetGlobal("package")
	if !L.IsTable(-1) {
		return fmt.Errorf("package is not a table? why? en? @_@?")
	}

	L.PushString(luaPathFull)
	L.SetField(-2, "path")

	L.PushString(luaCPathFull)
	L.SetField(-2, "cpath")

	L.Pop(-1)

	return nil
}
