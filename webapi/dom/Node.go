package dom

import (
	"syscall/js"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
)

type ToNodeConvertable interface {
	AsNode() *Node
}

type Node struct {
	JSValue js.Value
}

func (self *Node) AppendChild(node *Node) *Node {
	return &Node{gojstoolsutils.JSValueLiteralToPointer(self.JSValue.Call("appendChild", *node.JSValue))}
}

func (self *Node) GetFirstChild() *Node {

	ret := (*Node)(nil)

	r := self.JSValue.Get("firstChild")
	if !r.IsNull() {
		ret = &Node{&r}
	}

	return ret
}

func (self *Node) RemoveChild(c *Node) *Node {
	return &Node{gojstoolsutils.JSValueLiteralToPointer(self.JSValue.Call("removeChild", c.JSValue))}
}

func (self *Node) ParentNode() *Node {
	t := self.JSValue.Get("parentNode")
	if t.IsNull() || t.IsUndefined() {
		return nil
	}
	return &Node{&t}
}

func (self *Node) ParentElement() *Element {
	t := self.JSValue.Get("parentElement")
	if t.IsNull() || t.IsUndefined() {
		return nil
	}
	return &Element{&Node{&t}}
}

func (self *Node) AsNode() *Node {
	return self
}

func (self *Node) RemoveAllChildren() {
	for a := self.GetFirstChild(); a != nil; a = self.GetFirstChild() {
		self.RemoveChild(a)
	}
}
