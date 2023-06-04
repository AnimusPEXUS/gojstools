package events

import (
	"syscall/js"
)

type Event struct {
	JSValue js.Value
}

func NewEventFromJSValue(jsvalue js.Value) (*Event, error) {
	self := &Event{}
	self.JSValue = jsvalue
	return self, nil
}
