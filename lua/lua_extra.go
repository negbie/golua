package lua

import (
	"fmt"
	"reflect"
)

/* =========================== life cycle ============================= */

func (L *State) AddCloseHandler(callback LuaGoFunction) {
	L.mutex.Lock()
	defer L.mutex.Unlock()
	if L.closeHandlers == nil {
		L.closeHandlers = []LuaGoFunction{callback}
	} else {
		L.closeHandlers = append(L.closeHandlers, callback)
	}
}

func AddCloseHandlerDefault(L *State) int {
	if !L.IsFunction(-1) {
		L.PushString("invalid argument, close callback must be a function")
		return 1
	}
	var callbackRef = L.Ref(LUA_REGISTRYINDEX)
	L.AddCloseHandler(func(L *State) int {
		L.RawGeti(LUA_REGISTRYINDEX, callbackRef)
		_ = L.CallHandle(0, 0, nil)
		return 0
	})
	L.PushNil()
	return 1
}

/* =========================== table ============================= */

func (L *State) GetTableByName(table string, createIfNil bool) (exist bool, err error) {
	L.GetGlobal(table)
	if L.IsNil(-1) {
		L.Pop(1)
		if createIfNil {
			L.NewTable()
			L.SetGlobal(table)
			L.GetGlobal(table)
			return true, nil
		} else {
			return false, nil
		}
	}
	if L.IsTable(-1) {
		return false, fmt.Errorf("is not a table : " + table)
	}
	return true, nil
}

func (L *State) TableSetValue(tableIndex int, key string, val interface{}) (err error) {

	defer func() {
		if err == nil {
			L.SetField(tableIndex, key)
		}
	}()
	if val == nil {
		L.PushNil()
		return
	}

	if gofunc, ok := val.(func(*State) int); ok {
		L.PushGoFunction(gofunc)
		return
	}

	if str, ok := val.(string); ok {
		L.PushString(str)
		return
	}

	if bytes, ok := val.([]byte); ok {
		L.PushBytes(bytes)
		return
	}

	var vref = reflect.ValueOf(val)
	var kind = vref.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var n = AsInt64(val, 0)
		L.PushInteger(n)
	case reflect.Float32, reflect.Float64:
		var n = AsFloat64(val, 0)
		L.PushNumber(n)
	case reflect.Bool:
		var b = val.(bool)
		L.PushBoolean(b)
	case reflect.Func:
		// TODO
		L.PushGoFunction(func(L *State) int {
			vref.Call(nil)
			return 0
		})
	case reflect.Map:
		err = fmt.Errorf("map not support")
	case reflect.Ptr, reflect.Chan, reflect.Array, reflect.Slice:
		var uptr = vref.Pointer()
		L.PushInteger(int64(uptr))
	default:
		L.PushGoStruct(val)
	}
	return nil
}

func (L *State) TableRegister(table string, name string, val interface{}) error {

	if len(table) == 0 {
		table = "_G"
	}

	var _, err = L.GetTableByName(table, true)
	if err != nil {
		return err
	}
	err = L.TableSetValue(-2, name, val)
	L.Pop(1)
	return nil
}

func (L *State) TableRegisters(table string, funcs map[string]interface{}) error {
	if funcs == nil {
		return fmt.Errorf("no function is set")
	}
	var _, err = L.GetTableByName(table, true)
	if err != nil {
		return err
	}
	for key, val := range funcs {
		err = L.TableSetValue(-2, key, val)
		if err != nil {
			break
		}
	}
	L.Pop(1)
	return err
}

func (L *State) TableGetString(tableIndex int, key string, def string, throwNil bool) (string, error) {
	L.GetField(tableIndex, key)
	defer L.Pop(1)

	if L.IsNil(-1) {
		if throwNil {
			return def, fmt.Errorf("nil value")
		} else {
			return def, nil
		}
	}
	if !L.IsString(-1) {
		return def, fmt.Errorf("is not string")
	}
	var result = L.ToString(-1)
	return result, nil
}

func (L *State) TableGetInteger(tableIndex int, key string, def int, throwNil bool) (int, error) {
	L.GetField(tableIndex, key)
	defer L.Pop(1)

	if L.IsNil(-1) {
		if throwNil {
			return def, fmt.Errorf("nil value")
		} else {
			return def, nil
		}
	}

	if !L.IsNumber(-1) {
		return def, fmt.Errorf("is not string")
	}
	var result = L.ToInteger(-1)
	return result, nil
}

func (L *State) TableGetNumber(tableIndex int, key string, def float64, throwNil bool) (float64, error) {
	L.GetField(tableIndex, key)
	defer L.Pop(1)

	if L.IsNil(-1) {
		if throwNil {
			return def, fmt.Errorf("nil value")
		} else {
			return def, nil
		}
	}
	if !L.IsNumber(-1) {
		return def, fmt.Errorf("is not string")
	}
	var result = L.ToNumber(-1)
	return result, nil
}

func (L *State) TableGetBoolean(tableIndex int, key string, def bool, throwNil bool) (bool, error) {
	L.GetField(tableIndex, key)
	defer L.Pop(1)

	if L.IsNil(-1) {
		if throwNil {
			return def, fmt.Errorf("nil value")
		} else {
			return def, nil
		}
	}
	if !L.IsBoolean(-1) {
		return def, fmt.Errorf("is not string")
	}
	var result = L.ToBoolean(-1)
	return result, nil
}

func (L *State) TableGetValue(tableIndex int, key string, def interface{}, throwNil bool) (interface{}, error) {
	L.GetField(tableIndex, key)
	defer L.Pop(1)

	if L.IsNil(-1) {
		if throwNil {
			return def, fmt.Errorf("nil value")
		} else {
			return def, nil
		}
	}

	if L.IsString(-1) {
		return L.ToString(-1), nil
	}

	if L.IsNumber(-1) {
		return L.ToNumber(-1), nil
	}

	if L.IsBoolean(-1) {
		return L.ToBoolean(-1), nil
	}

	return nil, fmt.Errorf("unsupport type : %v", L.Typename(-1))
}

func (L *State) TableGetAndRef(tableIndex int, key string, throwNil bool, judgement func(L *State, tableIndex int, key string) error) (int, error) {
	L.GetField(tableIndex, key)
	defer L.Pop(1)

	if L.IsNil(-1) {
		if throwNil {
			return -1, fmt.Errorf("nil value")
		} else {
			return -1, nil
		}
	}
	if judgement != nil {
		var err = judgement(L, tableIndex, key)
		if err != nil {
			return -1, err
		}
	}
	var result = L.Ref(LUA_REGISTRYINDEX)
	return result, nil
}

/* ================== go life cycle ref keeping ==================== */

func (L *State) KeepGoRef(obj interface{}) int {
	L.gorefmutex.Lock()
	if L.gorefs == nil {
		L.gorefs = make([]interface{}, 16)
	}

	var id = L.gorefptr
	L.gorefptr = L.gorefptr + 1
	if L.gorefptr >= len(L.gorefs) {
		var nova = make([]interface{}, len(L.gorefs)+8)
		copy(nova, L.gorefs)
		L.gorefs = nova
	}
	L.gorefmutex.Unlock()
	return id
}

func (L *State) GetGoRef(ref int) interface{} {
	if L.gorefs == nil {
		return nil
	}
	L.gorefmutex.RLock()
	defer L.gorefmutex.RUnlock()
	if ref >= len(L.gorefs) || ref < 0 {
		return nil
	}
	return L.gorefs[ref]
}

func (L *State) ClearGoRef() {
	if L.gorefs == nil {
		return
	}
	L.gorefmutex.Lock()
	for i := 0; i < len(L.gorefs); i++ {
		L.gorefs[i] = nil
	}
	L.gorefmutex.Unlock()
}

/* ================== extra data ==================== */

func (L *State) GetData(key string) interface{} {
	if L.data == nil {
		return nil
	}
	L.mutex.RLock()
	defer L.mutex.RUnlock()
	return L.data[key]
}

func (L *State) SetData(key string, val interface{}) {
	L.mutex.Lock()
	defer L.mutex.Unlock()
	if L.data == nil {
		L.data = make(map[string]interface{})
	}
	L.data[key] = val
}

func (L *State) String() string {
	return fmt.Sprintf("[%v -> %v]", L.Name, L.Path)
}

/* ==================== info ============================= */

func (L *State) ToInterface(idx int, retptr bool) interface{} {
	var luaType = L.Type(idx)
	switch luaType {
	case LUA_TNIL:
		return nil
	case LUA_TNUMBER:
		return L.ToNumber(idx)
	case LUA_TBOOLEAN:
		return L.ToBoolean(idx)
	case LUA_TSTRING:
		return L.ToString(idx)
	case LUA_TTABLE:
		if retptr {
			return nil
		} else {
			return "[table]"
		}
	case LUA_TFUNCTION:
		if retptr {
			return L.ToGoFunction(idx)
		} else {
			return "[function]"
		}
	case LUA_TUSERDATA:
		if retptr {
			return L.ToUserdata(idx)
		} else {
			return "[userdata]"
		}
	case LUA_TTHREAD:
		if retptr {
			return L.ToThread(idx)
		} else {
			return "[thread]"
		}
	case LUA_TLIGHTUSERDATA:
		if retptr {
			return nil
		} else {
			return "[lightuserdata]"
		}
	case LUA_TNONE:
		return nil
	}
	return nil
}

func (L *State) StackToString() string {
	var info = ""
	var top = L.GetTop()
	for i := top; i >= 1; i-- {
		var luaTypeName = L.TypeString(i)
		var luaValue = L.ToInterface(i, false)
		info = fmt.Sprintf("%v[%d] %v | %v\n", info, i, luaTypeName, luaValue)
	}
	return info
}
