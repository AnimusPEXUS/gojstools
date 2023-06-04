package events

import (
	"syscall/js"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
	utils_panic "github.com/AnimusPEXUS/utils/panic"
)

type MessageEvent struct {
	Event
}

func NewMessageEventFromJSValue(jsvalue js.Value) (*MessageEvent, error) {
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
		err = utils_panic.PanicToError()
	}()
	ret = gojstoolsutils.JSValueLiteralToPointer(self.Event.JSValue.Get("data"))
	return ret, nil
}

// TODO: testing required
func (self *MessageEvent) GetOrigin() (ret string, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("origin").String()
	return ret, nil
}

func (self *MessageEvent) GetLastEventId() (ret string, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.Event.JSValue.Get("lastEventId").String()
	return ret, nil
}

// TODO: work required
func (self *MessageEvent) GetSource() (ret js.Value, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = gojstoolsutils.JSValueLiteralToPointer(self.Event.JSValue.Get("source"))
	return ret, nil
}

// TODO: work required
func (self *MessageEvent) GetPorts() (ret js.Value, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = gojstoolsutils.JSValueLiteralToPointer(self.Event.JSValue.Get("ports"))
	return ret, nil
}
