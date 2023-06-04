package window

import (
	"syscall/js"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
	"github.com/AnimusPEXUS/gojswebapi/dom"
)

type Window struct {
	JSValue js.Value
}

func NewWindowFromGlobalThis() (*Window, error) {
	return NewWindowFromJSValue(gojstoolsutils.JSValueLiteralToPointer(js.Global()))
}

func NewWindowFromJSValue(value js.Value) (*Window, error) {
	self := &Window{}
	self.JSValue = value
	return self, nil
}

func (self *Window) Open(url string, windowname string, windowfeatures *string) (*Window, error) {

	windowfeatures_val := js.Undefined()
	if windowfeatures != nil {
		windowfeatures_val = js.ValueOf(*windowfeatures)
	}

	cres := self.JSValue.Call("open", url, windowname, windowfeatures_val)
	ret := &Window{
		JSValue: &cres,
	}
	return ret, nil
}

func (self *Window) GetDocument() *dom.Document {
	doc := self.JSValue.Get("document")
	ret := dom.NewDocumentFromJsValue(&doc)
	return ret
}
