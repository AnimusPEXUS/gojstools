package elementtreeconstructor

import (
	"syscall/js"

	"github.com/AnimusPEXUS/wasmtools/dom"
)

var _ dom.ToNodeConvertable = &ElementMutator{}

type ElementMutator struct {
	Element *dom.Element
}

func NewElementMutatorNS(doc *dom.Document, ns string, name string) *ElementMutator {
	ret := &ElementMutator{doc.CreateElementNS(ns, name)}
	return ret
}

func NewElementMutator(doc *dom.Document, name string) *ElementMutator {
	ret := &ElementMutator{doc.CreateElement(name)}
	return ret
}

func NewElementMutatorFromElement(e *dom.Element) *ElementMutator {
	ret := &ElementMutator{e}
	return ret
}

func NewElementMutatorFromNode(n *dom.Node) *ElementMutator {
	ret := NewElementMutatorFromElement(&dom.Element{Node: *n})
	return ret
}

func (self *ElementMutator) AppendChildren(children ...dom.ToNodeConvertable) *ElementMutator {

	// FIXME: append and remove operations have to be done at Node level

	for _, i := range children {
		self.Element.Append(i.AsNode())
	}
	return self
}

func (self *ElementMutator) RemoveAllChildren() *ElementMutator {
	self.Element.Node.RemoveAllChildren()
	return self
}

func (self *ElementMutator) Remove(children ...dom.ToNodeConvertable) *ElementMutator {
	// FIXME: append and remove operations have to be done at Node level
	for _, i := range children {
		self.Element.RemoveChild(i.AsNode())
	}
	return self
}

func (self *ElementMutator) RemoveFromParent() *ElementMutator {
	self.Parent().Remove(self.AsNode())
	return self
}

func (self *ElementMutator) Parent() *ElementMutator {

	// FIXME: this is Node level operation

	ret := (*ElementMutator)(nil)
	t := self.Element.ParentElement()
	if t != nil {
		ret = NewElementMutatorFromElement(t)
	}
	return ret
}

func (self *ElementMutator) AssignSelf(variable **ElementMutator) *ElementMutator {
	*variable = self
	return self
}

func (self *ElementMutator) AssignSelfDom(variable **dom.Element) *ElementMutator {
	*variable = self.Element
	return self
}

func (self *ElementMutator) ExternalUse(cb func(*ElementMutator)) *ElementMutator {
	cb(self)
	return self
}

func (self *ElementMutator) Call(property string, ret *interface{}, args ...interface{}) *ElementMutator {
	// log.Println("calling", property)
	t := self.Element.Node.Value.Call(property, args...)
	if ret != nil {
		*ret = t
	}
	return self
}

func (self *ElementMutator) SetAttribute(name string, value string) *ElementMutator {
	self.Element.SetAttribute(name, value)
	return self
}

// func (self *ElementMutator) GetAttribute((name string, value string, ret *interface{})) *ElementMutator {
// 	self.Element.SetAttribute(name, value)
// 	if ret != nil {
// 		*ret = t
// 	}
// 	return self
// }
//

func (self *ElementMutator) Set(property string, value interface{}) *ElementMutator {
	self.Element.Node.Value.Set(property, value)
	return self
}

func (self *ElementMutator) Get(property string) interface{} {
	return self.Element.Node.Value.Get(property)
}

func (self *ElementMutator) GetJsValue(property string) js.Value {
	return self.Element.Node.Value.Get(property)
}

func (self *ElementMutator) SelfJsValue() js.Value {
	return self.Element.Node.Value
}

func (self *ElementMutator) GetAssign(property string, ret *interface{}) *ElementMutator {
	t := self.Element.Node.Value.Get(property)
	if ret != nil {
		*ret = t
	}
	return self
}

func (self *ElementMutator) SetStyle(property string, value interface{}) *ElementMutator {
	self.Element.Node.Value.Get("style").Set(property, value)
	return self
}

func (self *ElementMutator) AddListener(
	event string,
	cb func(this js.Value, args []js.Value) interface{},
) *ElementMutator {
	self.Element.Node.Value.Get(event).Call("addListener", js.FuncOf(cb))
	return self
}

func (self *ElementMutator) AddEventListener(
	event string,
	cb func(this js.Value, args []js.Value) interface{},
) *ElementMutator {
	self.Element.Node.Value.Call("addEventListener", event, js.FuncOf(cb))
	return self
}

func (self *ElementMutator) AsElement() *dom.Element {
	return self.Element
}

func (self *ElementMutator) AsNode() *dom.Node {
	return self.AsElement().AsNode()
}
