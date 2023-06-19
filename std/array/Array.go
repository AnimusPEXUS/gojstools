package array

import (
	"errors"
	"syscall/js"

	utils_panic "github.com/AnimusPEXUS/utils/panic"
)

var ERR_ARRAY_UNSUPPORTED = errors.New("Array unsupported")

func GetGlobalArrayJSValue(type_ ArrayType) (js.Value, error) {
	return js.Global().Get(type_.String()), nil
}

func DetermineArrayType(v js.Value) (ret ArrayType, ok bool) {
	ret = ""
	ok = false
	for _, i := range ArrayTypes {
		global_array_type := js.Global().Get(i.String())
		if global_array_type.IsUndefined() || global_array_type.IsNull() {
			continue
		}
		if v.InstanceOf(global_array_type) {
			ret = i
			ok = true
		}
	}
	return
}

type ArrayType string

func (self ArrayType) String() string {
	return string(self)
}

const (
	ArrayTypeArray        ArrayType = "Array"
	ArrayTypeInt8         ArrayType = "Int8Array"
	ArrayTypeUint8        ArrayType = "Uint8Array"
	ArrayTypeUint8Clamped ArrayType = "Uint8ClampedArray"
	ArrayTypeInt16        ArrayType = "Int16Array"
	ArrayTypeUint16       ArrayType = "Uint16Array"
	ArrayTypeInt32        ArrayType = "Int32Array"
	ArrayTypeUint32       ArrayType = "Uint32Array"
	ArrayTypeFloat32      ArrayType = "Float32Array"
	ArrayTypeFloat64      ArrayType = "Float64Array"
	ArrayTypeBigInt64     ArrayType = "BigInt64Array"
	ArrayTypeBigUint64    ArrayType = "BigUint64Array"
)

var ArrayTypes = []ArrayType{
	ArrayTypeInt8,
	ArrayTypeUint8,
	ArrayTypeUint8Clamped,
	ArrayTypeInt16,
	ArrayTypeUint16,
	ArrayTypeInt32,
	ArrayTypeUint32,
	ArrayTypeFloat32,
	ArrayTypeFloat64,
	ArrayTypeBigInt64,
	ArrayTypeBigUint64,
	ArrayTypeArray,
}

type Array struct {
	JSValue js.Value
}

func NewArray(
	array_type ArrayType,
	length_typedArray_object_or_buffer js.Value,
	byteOffset *js.Value,
	length *js.Value,
) (self *Array, err error) {

	defer func() {
		if p_err := utils_panic.PanicToError(); p_err != nil {
			err = p_err
		}
	}()

	found := false
	for _, i := range ArrayTypes {
		if i == array_type {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("Invalid array type name")
	}

	array_type_s := array_type.String()

	array_type_js := js.Global().Get(array_type_s)
	if array_type_js.IsUndefined() {
		return nil, errors.New(array_type_s + " undefined")
	}

	ud := js.Undefined()

	if byteOffset == nil {
		byteOffset = &ud
	}

	if length == nil {
		length = &ud
	}

	js_array := array_type_js.New(
		length_typedArray_object_or_buffer,
		*byteOffset,
		*length,
	)

	self, err = NewArrayFromJSValue(js_array)
	return self, err
}

func NewArrayFromJSValue(value js.Value) (self *Array, err error) {

	defer func() {
		if p_err := utils_panic.PanicToError(); p_err != nil {
			err = p_err
		}
	}()

	found := false
	for _, i := range ArrayTypes {
		js_type := js.Global().Get(i.String())
		if js_type.IsUndefined() {
			return nil, errors.New(i.String() + " undefined")
		}
		if value.InstanceOf(js_type) {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("unsupported type")
	}

	self = &Array{JSValue: value}
	return self, nil
}

func NewArrayFromByteSlice(data []byte) (self *Array, err error) {
	len_data := len(data)

	self, err = NewArray(
		ArrayTypeUint8,
		js.ValueOf(len_data),
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	res := js.CopyBytesToJS(self.JSValue, data)
	if res != len_data {
		return nil, errors.New("data length doesn't match copied data size")
	}

	return self, nil
}

func (self *Array) Type() (t ArrayType, ok bool) {
	return DetermineArrayType(self.JSValue)
}

func (self *Array) ToString() string {
	return self.JSValue.Call("toString").String()
}

func (self *Array) Length() int {
	return self.JSValue.Get("length").Int()
}

func (self *Array) GetU8Bytes() ([]byte, error) {
	t, ok := self.Type()
	if !ok {
		return nil, errors.New("invalid type")
	}

	if t != ArrayTypeUint8 {
		return nil, errors.New("not an ArrayTypeUint8 array")
	}

	l := self.Length()

	ret := make([]byte, l)
	c_s := js.CopyBytesToGo(ret, self.JSValue)

	if l != c_s {
		return nil, errors.New("copied length misscatch")
	}

	return ret, nil
}
