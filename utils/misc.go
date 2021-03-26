package misc

import (
	"strings"
	"syscall/js"
)

func StrBoolToBool(value string) bool {
	value = strings.ToLower(value)

	for _, i := range []string{"true", "yes", "1", "enable", "enabled", "on", "positive"} {
		if value == i {
			return true
		}
	}

	return false
}

func JSValueLiteralToPointer(in js.Value) *js.Value {
	return &[]js.Value{in}[0]
}
