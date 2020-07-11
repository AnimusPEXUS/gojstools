package elementtreeconstructor

import (
	"github.com/AnimusPEXUS/wasmtools/dom"
)

type ElementTreeConstructor struct {
	document *dom.Document
}

func NewElementTreeConstructor(document *dom.Document) *ElementTreeConstructor {
	self := &ElementTreeConstructor{
		document: document,
	}
	return self
}

func (self *ElementTreeConstructor) CreateTextNode(
	text string,
) *dom.Node {
	return self.document.NewTextNode(text)
}

func (self *ElementTreeConstructor) CreateElement(name string) *ElementMutator {
	return self.CreateElementNS(nil, name)
}

func (self *ElementTreeConstructor) CreateElementNS(namespace *string, name string) *ElementMutator {

	var ret *dom.Element

	if namespace != nil {
		ret = self.document.CreateElementNS(*namespace, name)
	} else {
		ret = self.document.CreateElement(name)
	}

	return NewElementMutatorFromElement(ret)
}

func (self *ElementTreeConstructor) ReplaceChildren(new_children []dom.ToNodeConvertable) {

	n := &dom.Node{self.document.Value}

	for i := n.GetFirstChild(); i != nil; i = n.GetFirstChild() {
		n.RemoveChild(i)
	}

	for _, i := range new_children {
		// log.Println("appending child", i, i.Value)
		n.AppendChild(i.AsNode())
	}

}
