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

func NewDOMImplementationFromGlobal() *DOMImplementation {
	self := new(DOMImplementation)
	self.JSValue = GetGlobalDOMImplementationJSValue()
	if self.JSValue.IsUndefined() {
		return nil
	}
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

func (self *DOMImplementation) CreateHTMLDocument(qualifiedNameStr string, publicId string, systemId string) *Document {
	r := self.JSValue.Call("createHTMLDocumentType", qualifiedNameStr, publicId, systemId)
	return NewDocumentFromJsValue(r)
}
