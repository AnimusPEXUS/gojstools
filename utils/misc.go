package misc

import (
	"strings"
)

func StrBoolToBool(value string) bool {
	value = strings.ToLower(value)

	for _, i := range []string{
		"+",
		"1",
		"enable",
		"enabled",
		"on",
		"positive",
		"true",
		"y",
		"yes",
	} {
		if value == i {
			return true
		}
	}

	return false
}

// func JSValueLiteralToPointer(in js.Value) *js.Value {
// 	return &[]js.Value{in}[0]
// }

// func JSFuncLiteralToPointer(in js.Func) *js.Func {
// 	return &[]js.Func{in}[0]
// }
