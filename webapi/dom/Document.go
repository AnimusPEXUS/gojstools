package dom

import (
	"syscall/js"
)

type Document struct {
	JSValue js.Value
}

func NewDocumentFromJsValue(jsvalue js.Value) *Document {
	self := &Document{jsvalue}
	return self
}

func (self *Document) Implementation() *DOMImplementation {
	return NewDOMImplementationFromJSValue(self.JSValue.Get("implementation"))
}

func (self *Document) CreateElementNS(ns string, name string) *Element {
	return &Element{
		&Node{
			self.JSValue.Call(
				"createElementNS",
				ns,
				name,
				js.Undefined(),
			),
		},
	}
}

func (self *Document) CreateElement(name string) *Element {
	return &Element{
		&Node{
			self.JSValue.Call(
				"createElement",
				name,
				js.Undefined(),
			),
		},
	}
}

func (self *Document) NewTextNode(text string) *Node {
	return &Node{
		self.JSValue.Call("createTextNode", text),
	}
}

func (self *Document) GetBody() *Element {
	return &Element{
		&Node{
			self.JSValue.Get("body"),
		},
	}
}
