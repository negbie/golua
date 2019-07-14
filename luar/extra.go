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

func GetVal(L *lua.State, idx int) (interface{}, error) {
	if L.IsNoneOrNil(idx) {
		return nil, nil
	}
	var ltype = int(L.Type(idx))
	switch ltype {
	case int(lua.LUA_TNUMBER):
		return L.ToNumber(idx), nil
	case int(lua.LUA_TSTRING):
		return L.ToString(idx), nil
	case int(lua.LUA_TBOOLEAN):
		return L.ToBoolean(idx), nil
	}
	var r interface{}
	var err = LuaToGo(L, idx, &r)
	return r, err
}

func FormatStack(stacks []lua.LuaStackEntry) []lua.LuaStackEntry {
	var count = len(stacks)
	var clones = make([]lua.LuaStackEntry, count)
	for i := 0; i < count; i++ {
		var stack = stacks[i]
		var clone = lua.LuaStackEntry{
			Name: stack.Name,
		}
		var linenum = stack.CurrentLine
		if linenum >= 0 {
			var lines = strings.Split(stack.Source, "\n")
			if linenum < len(lines) {
				clone.ShortSource = lines[linenum-1]
			} else {
				clone.ShortSource = stack.ShortSource
			}
		}
		clone.Source = ""
		clone.CurrentLine = linenum
		clones[i] = clone
	}
	return clones
}

func FormatStackToString(stacks []lua.LuaStackEntry, prefix string, suffix string) string {
	var str = ""
	var count = len(stacks)

	for i := 0; i < count; i++ {
		var stack = stacks[i]
		var source = stack.ShortSource
		var funcname = stack.Name
		var linenum = stack.CurrentLine
		if linenum >= 0 {
			var lines = strings.Split(stack.Source, "\n")
			if linenum < len(lines) {
				source = lines[linenum-1]
			}
		}
		var one = fmt.Sprintf("%s%s %s %d\n%s", prefix, source, funcname, linenum, suffix)
		str = str + one
	}
	return str
}

func FormatStackToMap(stacks []lua.LuaStackEntry) []map[string]interface{} {
	var count = len(stacks)
	var clones = make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		var stack = stacks[i]
		var clone = make(map[string]interface{})
		var linenum = stack.CurrentLine
		if linenum >= 0 {
			var lines = strings.Split(stack.Source, "\n")
			if linenum < len(lines) {
				clone["linesrc"] = lines[linenum-1]
			} else {
				clone["linesrc"] = stack.ShortSource
			}
		}
		clone["line"] = linenum
		clone["func"] = stack.Name
		clones[i] = clone
	}
	return clones
}
