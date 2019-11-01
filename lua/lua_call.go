package lua

func pcall(L *State) int {
	var top = L.GetTop()
	if top == 0 {
		L.PushBoolean(false)
		L.PushString("pcall() paramter 1 need to be a function")
		return 2
	}
	if !L.IsFunction(1) {
		L.PushBoolean(false)
		L.PushString("pcall() paramter 1 is not a function. current type = " + L.Typename(1))
		return 2
	}

	var err = L.Call(top-1, LUA_MULTRET)
	if err != nil {
		L.PushBoolean(false)
		var trans = LuaErrorTrans(err, "\t", "")
		if trans == nil {
			L.PushString(err.Error())
		} else {
			L.PushString(trans.Error())
		}
		return 2
	}

	L.PushBoolean(true)
	L.Insert(1)
	top = L.GetTop()
	return top
}

func xpcall(L *State) int {
	var top = L.GetTop()

	if top < 2 {
		L.PushBoolean(false)
		L.PushString("xpcall() paramter 1 & 2 need to be a function")
		return 2
	}

	if !L.IsFunction(1) {
		L.PushBoolean(false)
		L.PushString("xpcall() paramter 1 is not a function. current type = " + L.Typename(1))
		return 2
	}

	if !L.IsFunction(2) {
		L.PushBoolean(false)
		L.PushString("xpcall() paramter 2 is not a function. current type = " + L.Typename(1))
		return 2
	}

	// copy error handler to top
	L.PushValue(2)
	// remove pos 2 error handler
	L.Remove(2)
	// insert error handler to bottom
	L.Insert(1)

	var callerr = L.Call(top-2, LUA_MULTRET)
	if callerr != nil {
		callerr = LuaErrorTrans(callerr, "\t", "")
		if callerr != nil {
			L.PushString(callerr.Error())
			return 1
		}
		return 0
	}

	L.PushBoolean(true)
	L.Insert(1)
	top = L.GetTop()
	return top
}
