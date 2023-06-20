package promise

import (
	"errors"
	"syscall/js"
)

func GetGlobalPromiseJSValue() js.Value {
	return js.Global().Get("Promise")
}

func IsPromiseSupported() bool {
	return !GetGlobalPromiseJSValue().IsUndefined()
}

func ValueIsInstanceOfPromise(v js.Value) bool {
	cl := GetGlobalPromiseJSValue()
	if cl.IsUndefined() {
		return false
	}
	return v.InstanceOf(cl)
}

type Promise struct {
	JSValue js.Value
}

func NewPromiseFromJSValue(jsvalue js.Value) (*Promise, error) {

	if !IsPromiseSupported() {
		return nil, errors.New("Blob not supported")
	}

	if !ValueIsInstanceOfPromise(jsvalue) {
		return nil, errors.New("not an instance of Promise")
	}

	self := &Promise{JSValue: jsvalue}
	return self, nil
}

// https://developer.mozilla.org/en-US/docs/Web/
// JavaScript/Reference/Global_Objects/Promise/then
func (self *Promise) Then(cb js.Func) (*Promise, error) {

	ret, err := NewPromiseFromJSValue(
		self.JSValue.Call("then", cb),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (self *Promise) Then2(cb1 js.Func, cb2 js.Func) (*Promise, error) {

	ret, err := NewPromiseFromJSValue(
		self.JSValue.Call("then", cb1, cb2),
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
