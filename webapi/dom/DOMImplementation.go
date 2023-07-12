package dom

import "syscall/js"

func GetGlobalDOMImplementationJSValue() js.Value {
	return js.Global().Get("DOMImplementation")
}

func IsDOMImplementationSupported() bool {
	return !GetGlobalDOMImplementationJSValue().IsUndefined()
}

type DOMImplementation struct {
	JSValue js.Value
}

func NewDOMImplementationFromJSValue(val js.Value) *DOMImplementation {
	// TODO: input check
	self := new(DOMImplementation)
	self.JSValue = val
	return self
}

// func (self *DOMImplementation) CreateDocument2(namespaceURI, qualifiedNameStr) *XMLDocument {

// }

// func (self *DOMImplementation) CreateDocument3(namespaceURI, qualifiedNameStr, documentType) *XMLDocument {
// }

func (self *DOMImplementation) CreateDocumentType(qualifiedNameStr string, publicId string, systemId string) *DocumentType {
	r := self.JSValue.Call("createDocumentType", qualifiedNameStr, publicId, systemId)
	return NewDocumentTypeFromJsValue(r)
}

func (self *DOMImplementation) CreateHTMLDocument(title string) *Document {
	r := self.JSValue.Call("createHTMLDocument", title)
	return NewDocumentFromJsValue(r)
}
