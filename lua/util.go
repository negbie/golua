package lua

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func AsInt64(o interface{}, defaultval int64) (r int64) {

	if o == nil {
		return defaultval
	}

	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.Int64:
		return o.(int64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return vref.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(vref.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(vref.Float())
	case reflect.Bool:
		var b = o.(bool)
		if b {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		var s = o.(string)
		var i64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return defaultval
		}
		return i64
	}

	switch o.(type) {
	case time.Time:
		var t = o.(time.Time)
		return int64(t.Unix())
	case *time.Time:
		var t = o.(*time.Time)
		return int64(t.Unix())
	}

	panic(fmt.Errorf("convert not support type %v value %v ", reflect.TypeOf(o), reflect.ValueOf(o)))
}

func AsFloat64(o interface{}, defaultval float64) (r float64) {
	if o == nil {
		return defaultval
	}
	var vref = reflect.ValueOf(o)
	var kind = vref.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(vref.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(vref.Uint())
	case reflect.Float64:
		return o.(float64)
	case reflect.Float32:
		return float64(o.(float32))
	case reflect.Bool:
		var b = o.(bool)
		if b {
			return 1
		} else {
			return 0
		}
	case reflect.String:
		var s = o.(string)
		var f64, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return defaultval
		}
		return f64
	}

	switch o.(type) {
	case time.Time:
		var t = o.(time.Time)
		return float64(t.Unix())
	case *time.Time:
		var t = o.(*time.Time)
		return float64(t.Unix())
	}

	panic(fmt.Errorf("convert not support type %v value %v ", reflect.TypeOf(o), reflect.ValueOf(o)))
}
