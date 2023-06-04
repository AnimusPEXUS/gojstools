package array

import (
	"errors"
	"syscall/js"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
	utils_panic "github.com/AnimusPEXUS/utils/panic"
)

var ERR_ARRAY_UNSUPPORTED = errors.New("Array unsupported")

func GetArrayJSValue(type_ ArrayType) (js.Value, error) {
	return gojstoolsutils.JSValueLiteralToPointer(js.Global().Get(type_.String())), nil
}

func DetermineArrayType(v js.Value) *ArrayType {
	ret := (*ArrayType)(nil)
	for _, i := range ArrayTypes {
		global_array_type := js.Global().Get(i.String())
		if global_array_type.IsUndefined() || global_array_type.IsNull() {
			continue
		}
		if v.InstanceOf(global_array_type) {
			ret = &([]ArrayType{i}[0])
		}
	}
	return ret
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
	ArrayTypeArray,
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
}

type Array struct {
	JSValue js.Value
}

func NewArray(
	array_type ArrayType,
	length_typedArray_object_or_buffer js.Value,
	byteOffset js.Value,
	length js.Value,
) (self *Array, err error) {

	defer func() {
		err = utils_panic.PanicToError()
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

	array_type_js := gojstoolsutils.JSValueLiteralToPointer(js.Global().Get(array_type_s))
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

	js_array := gojstoolsutils.JSValueLiteralToPointer(
		array_type_js.New(
			length_typedArray_object_or_buffer,
			*byteOffset,
			*length,
		),
	)

	self, err = NewArrayFromJSValue(js_array)
	return self, err
}

func NewArrayFromJSValue(value js.Value) (self *Array, err error) {

	defer func() {
		err = utils_panic.PanicToError()
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

func (self *Array) ToString() string {
	return self.JSValue.Call("toString").String()
}
