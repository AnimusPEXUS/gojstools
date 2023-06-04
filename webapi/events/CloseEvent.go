package events

import (
	"syscall/js"

	utils_panic "github.com/AnimusPEXUS/utils/panic"
)

type CloseEvent struct {
	Event
}

func NewCloseEventFromJSValue(jsvalue js.Value) (*CloseEvent, error) {
	self := &CloseEvent{}
	r, err := NewEventFromJSValue(jsvalue)
	if err != nil {
		return nil, err
	}
	self.Event = *r
	return self, nil
}

func (self *CloseEvent) GetCode() (ret int, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("code").Int()
	return ret, nil
}

func (self *CloseEvent) GetReason() (ret string, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("reason").String()
	return ret, nil
}

func (self *CloseEvent) GetWasClean() (ret bool, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("wasclean").Bool()
	return ret, nil
}
