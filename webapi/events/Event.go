package events

import (
	"errors"
	"syscall/js"
)

func ValueIsInstanceOfEvent(v js.Value) bool {
	cl := js.Global().Get("Event")
	if cl.IsUndefined() {
		return false
	}
	return v.InstanceOf(cl)
}

type Event struct {
	JSValue js.Value
}

func NewEventFromJSValue(jsvalue js.Value) (*Event, error) {

	if !ValueIsInstanceOfEvent(jsvalue) {
		return nil, errors.New("not an instance of Event")
	}

	self := &Event{}
	self.JSValue = jsvalue
	return self, nil
}
