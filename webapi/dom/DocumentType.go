package dom

import (
	"syscall/js"
)

type DocumentType struct {
	Node
}

func NewDocumentTypeFromJsValue(jsvalue js.Value) *DocumentType {
	// TODO: input check?
	self := &DocumentType{Node{jsvalue}}
	return self
}
