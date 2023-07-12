package dom

import (
	"syscall/js"
)

type DocumentType struct {
	JSValue js.Value
}

func NewDocumentTypeFromJsValue(jsvalue js.Value) *DocumentType {
	// TODO: input check?
	self := &DocumentType{jsvalue}
	return self
}
