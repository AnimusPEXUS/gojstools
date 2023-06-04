package elementtreeconstructor

import (
	"github.com/AnimusPEXUS/gojstools/webapi/dom"
)

type ElementTreeConstructor struct {
	Document *dom.Document
}

func NewElementTreeConstructor(document *dom.Document) *ElementTreeConstructor {
	self := &ElementTreeConstructor{
		Document: document,
	}
	return self
}

func (self *ElementTreeConstructor) CreateTextNode(
	text string,
) *dom.Node {
	return self.Document.NewTextNode(text)
}

func (self *ElementTreeConstructor) CreateElement(name string) *ElementMutator {
	return self.CreateElementNS(nil, name)
}

func (self *ElementTreeConstructor) CreateElementNS(namespace *string, name string) *ElementMutator {

	var ret *dom.Element

	if namespace != nil {
		ret = self.Document.CreateElementNS(*namespace, name)
	} else {
		ret = self.Document.CreateElement(name)
	}

	return NewElementMutatorFromElement(ret)
}

func (self *ElementTreeConstructor) ReplaceChildren(new_children []dom.ToNodeConvertable) {

	n := &dom.Node{self.Document.JSValue}

	for i := n.GetFirstChild(); i != nil; i = n.GetFirstChild() {
		n.RemoveChild(i)
	}

	for _, i := range new_children {
		// log.Println("appending child", i, i.Value)
		n.AppendChild(i.AsNode())
	}

}
