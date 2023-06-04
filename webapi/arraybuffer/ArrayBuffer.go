package arraybuffer

import (
	"errors"
	"syscall/js"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
)

var ERR_ARRAYBUFFER_UNSUPPORTED = errors.New("ArrayBuffer unsupported")

func GetArrayBufferJSValue() (js.Value, error) {
	return gojstoolsutils.JSValueLiteralToPointer(js.Global().Get("ArrayBuffer")), nil
}

func IsArrayBufferSupported() (bool, error) {

	abjv, err := GetArrayBufferJSValue()
	if err != nil {
		return false, err
	}

	undef := abjv.IsUndefined()

	return !undef, nil
}

func IsArrayBuffer(value js.Value) (bool, error) {
	abjv, err := GetArrayBufferJSValue()
	if err != nil {
		return false, err
	}

	return value.InstanceOf(*abjv), nil
}

type ArrayBuffer struct {
	JSValue js.Value
}

func NewArrayBufferFromJSValue(jsvalue js.Value) (*ArrayBuffer, error) {
	self := &ArrayBuffer{JSValue: jsvalue}
	return self, nil
}

func NewArrayBuffer(length int) (*ArrayBuffer, error) {
	jsv_c, err := GetArrayBufferJSValue()
	if err != nil {
		return nil, err
	}

	jsv := gojstoolsutils.JSValueLiteralToPointer(jsv_c.New(length))

	self, err := NewArrayBufferFromJSValue(jsv)
	if err != nil {
		return nil, err
	}

	return self, nil
}

func (self *ArrayBuffer) IsArrayBuffer() (bool, error) {
	return IsArrayBuffer(self.JSValue)
}

func (self *ArrayBuffer) Len() (int, error) {
	return self.JSValue.Get("byteLength").Int(), nil
}

// TODO: maybe int64 is better solution, but I'm not sure
func (self *ArrayBuffer) Slice(begin int, end *int, contentType *string) (*ArrayBuffer, error) {

	begin_p := gojstoolsutils.JSValueLiteralToPointer(js.ValueOf(begin))
	end_p := js.Undefined()
	contentType_p := js.Undefined()

	if end != nil {
		end_p = js.ValueOf(*end)
	}

	if contentType != nil {
		contentType_p = js.ValueOf(*contentType)
	}

	ret_array := self.JSValue.Call("slice", begin_p, end_p, contentType_p)

	return NewArrayBufferFromJSValue(&ret_array)
}

func (self *ArrayBuffer) MakeReader() (*ArrayBufferReader, error) {
	return NewArrayBufferReader(self)
}
