package widgetcollection

import (
	"syscall/js"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
)

type ActiveLabel00 struct {
	Element *elementtreeconstructor.ElementMutator
}

func NewActiveLabel00(
	text string,
	title *string,

	etc *elementtreeconstructor.ElementTreeConstructor,

	onclick func(),
) *ActiveLabel00 {

	self := &ActiveLabel00{
		Element: etc.CreateElement("a").
			ExternalUse(applyAStyle).
			AppendChildren(
				etc.CreateTextNode(text),
			),
	}

	if onclick != nil {
		self.Element.Set(
			"onclick",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					onclick()
					return false
				},
			),
		)
	}

	if title != nil {
		self.Element.Set("title", *title)
	}

	return self
}

func (self *ActiveLabel00) AssignSelf(ref **ActiveLabel00) *ActiveLabel00 {
	*ref = self
	return self
}

func applyAStyle(ed *elementtreeconstructor.ElementMutator) {
	ed.
		SetStyle("color", "blue").
		SetStyle("cursor", "pointer").
		SetStyle("text-decoration", "underline")
}
