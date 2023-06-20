package blob

import (
	"errors"
	"syscall/js"

	"github.com/AnimusPEXUS/gojstools/std/arraybuffer"
	"github.com/AnimusPEXUS/gojstools/std/promise"

	utils_panic "github.com/AnimusPEXUS/utils/panic"
)

func GetGlobalBlobJSValue() js.Value {
	return js.Global().Get("Blob")
}

func IsBlobSupported() bool {
	return !GetGlobalBlobJSValue().IsUndefined()
}

func ValueIsInstanceOfBlob(v js.Value) bool {
	cl := GetGlobalBlobJSValue()
	if cl.IsUndefined() {
		return false
	}
	return v.InstanceOf(cl)
}

type Blob struct {
	JSValue js.Value
}

func NewBlobFromJSValue(jsvalue js.Value) (*Blob, error) {

	if !IsBlobSupported() {
		return nil, errors.New("Blob not supported")
	}

	if !ValueIsInstanceOfBlob(jsvalue) {
		return nil, errors.New("not an instance of Blob")
	}

	self := &Blob{JSValue: jsvalue}
	return self, nil
}

func NewBlobFromArray(array js.Value) (ret *Blob, err error) {

	defer func() {
		if p_err := utils_panic.PanicToError(); p_err != nil {
			err = p_err
		}
	}()

	if !IsBlobSupported() {
		return nil, errors.New("Blob not supported")
	}

	res := GetGlobalBlobJSValue().New(array)

	return &Blob{JSValue: res}, nil
}

func (self *Blob) Size() (int, error) {
	return self.JSValue.Get("size").Int(), nil
}

func (self *Blob) Type() (string, error) {
	return self.JSValue.Get("type").String(), nil
}

func (self *Blob) ArrayBuffer() (*arraybuffer.ArrayBuffer, error) {

	blob_arraybuffer_result := self.JSValue.Call("arrayBuffer")

	pro, err := promise.NewPromiseFromJSValue(blob_arraybuffer_result)
	if err != nil {
		return nil, err
	}

	psucc := make(chan struct{})
	perr := make(chan struct{})
	var array_data js.Value

	pro.Then2(
		js.FuncOf(
			func(
				this js.Value,
				args []js.Value,
			) interface{} {
				if len(args) == 0 {
					perr <- struct{}{}
					return false
				}
				array_data = args[0] // .Get("data")
				psucc <- struct{}{}
				return false
			},
		),

		js.FuncOf(
			func(
				this js.Value,
				args []js.Value,
			) interface{} {
				perr <- struct{}{}
				return false
			},
		),
	)

	select {
	case <-psucc:
		return arraybuffer.NewArrayBufferFromJSValue(array_data)
	case <-perr:
		return nil, errors.New("error getting Blob's ArrayBuffer")
	}

	// return nil, errors.New("invalid behavior")
}

// TODO: maybe int64 is better solution, but I'm not sure
func (self *Blob) Slice(start *int, end *int, contentType *string) (
	*Blob,
	error,
) {
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

	ret_blob := self.JSValue.Call(
		"slice",
		start_p,
		end_p,
		contentType_p,
	)

	return NewBlobFromJSValue(ret_blob)
}

// TODO: maybe later :)
// func (self *Blob) Stream() (*ReadableStream, error)

func (self *Blob) Text() (*promise.Promise, error) {
	blob_text_result := self.JSValue.Call("text")
	pro, err := promise.NewPromiseFromJSValue(blob_text_result)
	if err != nil {
		return nil, err
	}
	return pro, nil
}

func (self *Blob) MakeReader() (*BlobReader, error) {
	return NewBlobReader(self)
}
