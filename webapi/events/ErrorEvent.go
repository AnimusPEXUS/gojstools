package events

import (
	"syscall/js"

	utils_panic "github.com/AnimusPEXUS/utils/panic"
)

type ErrorEvent struct {
	Event
}

func NewErrorEventFromJSValue(jsvalue js.Value) (*ErrorEvent, error) {
	self := &ErrorEvent{}
	r, err := NewEventFromJSValue(jsvalue)
	if err != nil {
		return nil, err
	}
	self.Event = *r
	return self, nil
}

func (self *ErrorEvent) GetMessage() (ret string, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("message").String()
	return ret, nil
}

func (self *ErrorEvent) GetFilename() (ret string, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("filename").String()
	return ret, nil
}

func (self *ErrorEvent) GetLineno() (ret int, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("lineno").Int()
	return ret, nil
}

func (self *ErrorEvent) GetColno() (ret int, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("colno").Int()
	return ret, nil
}
