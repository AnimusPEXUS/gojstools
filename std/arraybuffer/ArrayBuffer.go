package arraybuffer

import (
	"errors"
	"syscall/js"
)

var ERR_ARRAYBUFFER_UNSUPPORTED = errors.New("ArrayBuffer unsupported")

func GetGlobalArrayBufferJSValue() js.Value {
	return js.Global().Get("ArrayBuffer")
}

func IsArrayBufferSupported() bool {
	return !GetGlobalArrayBufferJSValue().IsUndefined()
}

func ValueIsInstanceOfArrayBuffer(value js.Value) bool {
	abjv := GetGlobalArrayBufferJSValue()
	if abjv.IsUndefined() {
		return false
	}

	return value.InstanceOf(abjv)
}

type ArrayBuffer struct {
	JSValue js.Value
}

func NewArrayBufferFromJSValue(jsvalue js.Value) (*ArrayBuffer, error) {

	if !IsArrayBufferSupported() {
		return nil, errors.New("ArrayBuffer not supported")
	}

	if !ValueIsInstanceOfArrayBuffer(jsvalue) {
		return nil, errors.New("not an instance of ArrayBuffer")
	}

	self := &ArrayBuffer{JSValue: jsvalue}
	return self, nil
}

func NewArrayBuffer(length int) (*ArrayBuffer, error) {

	if !IsArrayBufferSupported() {
		return nil, errors.New("ArrayBuffer not supported")
	}

	jsv_c := GetGlobalArrayBufferJSValue()

	jsv := jsv_c.New(length)

	self, err := NewArrayBufferFromJSValue(jsv)
	if err != nil {
		return nil, err
	}

	return self, nil
}

func (self *ArrayBuffer) Len() (int, error) {
	return self.JSValue.Get("byteLength").Int(), nil
}

// TODO: maybe int64 is better solution, but I'm not sure
func (self *ArrayBuffer) Slice(
	begin int,
	end *int,
	contentType *string,
) (*ArrayBuffer, error) {

	begin_p := js.ValueOf(begin)
	end_p := js.Undefined()
	contentType_p := js.Undefined()

	if end != nil {
		end_p = js.ValueOf(*end)
	}

	if contentType != nil {
		contentType_p = js.ValueOf(*contentType)
	}

	ret_array := self.JSValue.Call(
		"slice",
		begin_p,
		end_p,
		contentType_p,
	)

	return NewArrayBufferFromJSValue(ret_array)
}

func (self *ArrayBuffer) MakeReader() (*ArrayBufferReader, error) {
	return NewArrayBufferReader(self)
}
