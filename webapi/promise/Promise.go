package promise

import (
	"syscall/js"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
)

type Promise struct {
	JSValue js.Value
}

func NewPromiseFromJSValue(jsvalue js.Value) (*Promise, error) {
	self := &Promise{JSValue: jsvalue}
	return self, nil
}

// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/then
func (self *Promise) Then(funcs ...*js.Func) (*Promise, error) {

	funcs2 := make([]interface{}, len(funcs))

	for i := 0; i != len(funcs); i++ {
		funcs2[i] = funcs[i]
	}

	ret, err := NewPromiseFromJSValue(gojstoolsutils.JSValueLiteralToPointer(self.JSValue.Call("then", funcs2...)))
	if err != nil {
		return nil, err
	}

	return ret, nil
}
