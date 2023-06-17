package events

import (
	"errors"
	"syscall/js"

	utils_panic "github.com/AnimusPEXUS/utils/panic"
)

func ValueIsInstanceOfMessageEvent(v js.Value) bool {
	cl := js.Global().Get("MessageEvent")
	if cl.IsUndefined() {
		return false
	}
	return v.InstanceOf(cl)
}

type MessageEvent struct {
	Event
}

func NewMessageEventFromJSValue(jsvalue js.Value) (*MessageEvent, error) {

	if !ValueIsInstanceOfMessageEvent(jsvalue) {
		return nil, errors.New("not an instance of MessageEvent")
	}

	self := &MessageEvent{}
	r, err := NewEventFromJSValue(jsvalue)
	if err != nil {
		return nil, err
	}
	self.Event = *r
	return self, nil
}

// https://developer.mozilla.org/en-US/docs/Web/API/MessageEvent/data
// says data can be of any type, so, probably, user have to decide what to do
// with it
func (self *MessageEvent) GetData() (ret js.Value, err error) {
	defer func() {
		if p_err := utils_panic.PanicToError(); p_err != nil {
			err = p_err
		}
	}()
	ret = self.Event.JSValue.Get("data")
	return ret, nil
}

// TODO: testing required
func (self *MessageEvent) GetOrigin() (ret string, err error) {
	defer func() {
		if p_err := utils_panic.PanicToError(); p_err != nil {
			err = p_err
		}
	}()
	ret = self.Event.JSValue.Get("origin").String()
	return ret, nil
}

func (self *MessageEvent) GetLastEventId() (ret string, err error) {
	defer func() {
		if p_err := utils_panic.PanicToError(); p_err != nil {
			err = p_err
		}
	}()
	ret = self.Event.JSValue.Get("lastEventId").String()
	return ret, nil
}

// TODO: work required
func (self *MessageEvent) GetSource() (ret js.Value, err error) {
	defer func() {
		if p_err := utils_panic.PanicToError(); p_err != nil {
			err = p_err
		}
	}()
	ret = self.Event.JSValue.Get("source")
	return ret, nil
}

// TODO: work required
func (self *MessageEvent) GetPorts() (ret js.Value, err error) {
	defer func() {
		if p_err := utils_panic.PanicToError(); p_err != nil {
			err = p_err
		}
	}()
	ret = self.Event.JSValue.Get("ports")
	return ret, nil
}
