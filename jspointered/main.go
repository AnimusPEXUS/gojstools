package jspointered

// DESCRIPTION: wraps syscall/js, so all js.Value returned or accepted becoming pointers

// import (
// 	"syscall/js"
// )

// // TODO: continue

// func CopyBytesToGo(dst []byte, src *Value) int {
// 	return js.CopyBytesToGo(dst, *src)
// }

// func CopyBytesToJS(dst *Value, src []byte) int {
// 	return CopyBytesToJS(*dst, src)
// }

// type Error struct {
// 	js.Error
// }

// type Func struct {
// 	js.Func
// }

// func FuncOf(fn func(this *Value, args []*Value) interface{}) *Func {

// 	fn2 := func(this2 js.Value, args2 []js.Value) interface{} {

// 		args3 := make([]*js.Value, len(args))
// 		for ii, i := range args2 {
// 			args3[ii] = &i
// 		}

// 		return fn(&this2, args3)
// 	}

// 	ret := js.FuncOf(fn2)
// 	return &ret
// }

// func (self *Func) Release() {
// 	self.Func.Release()
// }

// type Type struct {
// 	js.Type
// }

// type Value struct {
// 	js.Value
// }

// func Global() (*Value, error) {
// 	ret := &Value{js.Global()}
// 	return ret, nil
// }

// func Null() (*Value, error) {
// 	ret := &Value{js.Null()}
// 	return ret, nil
// }

// func Undefined() (*Value, error) {
// 	ret := &Value{js.Undefined()}
// 	return ret, nil
// }

// func ValueOf(x interface{}) (*Value, error) {
// 	ret := &Value{js.ValueOf(x)}
// 	return ret, nil
// }

// func (self *Value) Bool() bool {
// 	return self.Bool()
// }

// func (self *Value) Call(m string, args ...interface{}) *Value {
// 	self.Call(m, args)
// }
// func (self *Value) Delete(p string)
// func (self *Value) Equal(w Value) bool
// func (self *Value) Float() float64
// func (self *Value) Get(p string) *Value
// func (self *Value) Index(i int) *Value
// func (self *Value) InstanceOf(t *Value) bool
// func (self *Value) Int() int
// func (self *Value) Invoke(args ...interface{}) *Value
// func (self *Value) IsNaN() bool
// func (self *Value) IsNull() bool
// func (self *Value) IsUndefined() bool
// func (self *Value) JSValue() *Value
// func (self *Value) Length() int
// func (self *Value) New(args ...interface{}) *Value
// func (self *Value) Set(p string, x interface{})
// func (self *Value) SetIndex(i int, x interface{})
// func (self *Value) String() string
// func (self *Value) Truthy() bool
// func (self *Value) Type() Type

// type ValueError
//     func (e *ValueError) Error() string
// type Wrapper
