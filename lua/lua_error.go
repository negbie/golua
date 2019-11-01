package lua

import (
	"fmt"
	"strings"
)

type LuaStackEntry struct {
	Name        string
	Source      string
	ShortSource string
	CurrentLine int
}

type LuaError struct {
	code       int
	message    string
	stackTrace []LuaStackEntry
}

func (err *LuaError) Error() string {
	return err.message
}

func (err *LuaError) Code() int {
	return err.code
}

func (err *LuaError) StackTrace() []LuaStackEntry {
	return err.stackTrace
}

func (err *LuaError) StackTraceToString(prefix string, suffix string) string {
	if err.stackTrace == nil {
		return ""
	}

	var str = ""
	var count = len(err.stackTrace)

	for i := 0; i < count; i++ {
		var stack = err.stackTrace[i]
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
