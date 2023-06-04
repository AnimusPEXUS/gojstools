package dom

import (
	"syscall/js"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
)

type Document struct {
	JSValue js.Value
}

func NewDocumentFromJsValue(jsvalue js.Value) *Document {
	self := &Document{jsvalue}
	return self
}

func (self *Document) CreateElementNS(ns string, name string) *Element {
	return &Element{
		&Node{
			gojstoolsutils.JSValueLiteralToPointer(
				self.JSValue.Call(
					"createElementNS",
					ns,
					name,
					js.Undefined(),
				),
			),
		},
	}
}

func (self *Document) CreateElement(name string) *Element {
	return &Element{
		&Node{
			gojstoolsutils.JSValueLiteralToPointer(
				self.JSValue.Call(
					"createElement",
					name,
					js.Undefined(),
				),
			),
		},
	}
}

func (self *Document) NewTextNode(text string) *Node {
	return &Node{
		gojstoolsutils.JSValueLiteralToPointer(
			self.JSValue.Call("createTextNode", text),
		),
	}
}

func (self *Document) GetBody() *Element {
	return &Element{
		&Node{
			gojstoolsutils.JSValueLiteralToPointer(
				self.JSValue.Get("body"),
			),
		},
	}
}
