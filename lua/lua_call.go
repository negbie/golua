package lua

import "github.com/andrewhare/golua/lua"

func pcall(L *State) int {
	var top = L.GetTop()
	if top == 0 {
		L.PushString("pcall() paramter 1 need to be a function")
		return 1
	}
	if !L.IsFunction(1) {
		L.PushString("pcall() paramter 1 is not a function. current type = " + L.Typename(1))
		return 1
	}

	var callerr = L.Call(top-1, lua.LUA_MULTRET)
	if callerr != nil {
		callerr = LuaErrorTrans(callerr, "\t", "")
		if callerr != nil {
			L.PushString(callerr.Error())
			return 1
		}
		return 0
	}

	L.PushString("true")
	L.Insert(1)
	top = L.GetTop()
	return top
}

func xpcall(L *State) int {
	var top = L.GetTop()
	if top == 0 {
		L.PushString("xpcall() paramter 1 & 2 need to be a function")
		return 1
	}
	if !L.IsFunction(1) {
		L.PushString("xpcall() paramter 1 is not a function. current type = " + L.Typename(1))
		return 1
	}

	if !L.IsFunction(2) {
		L.PushString("xpcall() paramter 2 is not a function. current type = " + L.Typename(1))
		return 1
	}

	var callerr = L.Call(top-1, lua.LUA_MULTRET)
	if callerr != nil {
		callerr = LuaErrorTrans(callerr, "\t", "")
		if callerr != nil {
			L.PushString(callerr.Error())
			return 1
		}
		return 0
	}

	L.PushString("true")
	L.Insert(1)
	top = L.GetTop()
	return top
}
