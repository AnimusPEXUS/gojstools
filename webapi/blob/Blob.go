package blob

import (
	"errors"
	"syscall/js"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
	"github.com/AnimusPEXUS/gojswebapi/arraybuffer"
	"github.com/AnimusPEXUS/gojswebapi/promise"
)

var ERR_BLOB_UNSUPPORTED = errors.New("Blob unsupported")

func GetBlobJSValue() js.Value {
	return gojstoolsutils.JSValueLiteralToPointer(js.Global().Get("Blob"))
}

func IsBlobSupported() bool {
	return !GetBlobJSValue().IsUndefined()
}

func IsBlob(value js.Value) (bool, error) {
	return value.InstanceOf(*GetBlobJSValue()), nil
}

type Blob struct {
	JSValue js.Value
}

func NewBlobFromJSValue(jsvalue js.Value) (*Blob, error) {
	self := &Blob{JSValue: jsvalue}
	return self, nil
}

func (self *Blob) IsBlob() (bool, error) {
	return IsBlob(self.JSValue)
}

func (self *Blob) Size() (int, error) {
	return self.JSValue.Get("size").Int(), nil
}

func (self *Blob) Type() (string, error) {
	return self.JSValue.Get("type").String(), nil
}

func (self *Blob) ArrayBuffer() (*arraybuffer.ArrayBuffer, error) {

	blob_arraybuffer_result := self.JSValue.Call("arrayBuffer")

	pro, err := promise.NewPromiseFromJSValue(&blob_arraybuffer_result)
	if err != nil {
		return nil, err
	}

	psucc := make(chan bool)
	perr := make(chan bool)
	var array_data js.Value

	pro.Then(
		gojstoolsutils.JSFuncLiteralToPointer(
			js.FuncOf(
				func(
					this js.Value,
					args []js.Value,
				) interface{} {
					if len(args) == 0 {
						perr <- true
						return false
					}
					array_data = gojstoolsutils.JSValueLiteralToPointer(args[0].Get("data"))
					psucc <- true
					return false
				},
			)),

		gojstoolsutils.JSFuncLiteralToPointer(
			js.FuncOf(
				func(
					this js.Value,
					args []js.Value,
				) interface{} {
					perr <- true
					return false
				},
			)),
	)

	select {
	case <-psucc:
		return arraybuffer.NewArrayBufferFromJSValue(array_data)
	case <-perr:
		return nil, errors.New("error getting Blob's ArrayBuffer")
	}

	return nil, errors.New("invalid behavior")
}

// TODO: maybe int64 is better solution, but I'm not sure
func (self *Blob) Slice(start *int, end *int, contentType *string) (*Blob, error) {
	start_p := js.Undefined()
	end_p := js.Undefined()
	contentType_p := js.Undefined()

	if start != nil {
		start_p = js.ValueOf(*start)
	}

	if end != nil {
		end_p = js.ValueOf(*end)
	}

	if contentType != nil {
		contentType_p = js.ValueOf(*contentType)
	}

	ret_blob := self.JSValue.Call("slice", start_p, end_p, contentType_p)

	return NewBlobFromJSValue(&ret_blob)
}

// TODO: maybe later :)
// func (self *Blob) Stream() (*ReadableStream, error)

func (self *Blob) Text() (*promise.Promise, error) {
	blob_text_result := self.JSValue.Call("text")
	pro, err := promise.NewPromiseFromJSValue(&blob_text_result)
	if err != nil {
		return nil, err
	}
	return pro, nil
}

func (self *Blob) MakeReader() (*BlobReader, error) {
	return NewBlobReader(self)
}
