package widgetcollection

import (
	"syscall/js"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
)

type Select00 struct {
	etc *elementtreeconstructor.ElementTreeConstructor

	value_select *elementtreeconstructor.ElementMutator

	Value   string
	Element *elementtreeconstructor.ElementMutator
}

func NewSelect00(
	etc *elementtreeconstructor.ElementTreeConstructor,
	values [][2]string,
	preselected string,
	onchange func(),
) *Select00 {

	self := &Select00{etc: etc}

	self.value_select = self.etc.CreateElement("select")

	for _, i := range values {
		self.AppendOption(i[0], i[1])
	}

	if preselected != "" {
		self.value_select.Set("value", preselected)
	}

	self.value_select.Set(
		"onchange",
		js.FuncOf(
			func(
				this js.Value,
				args []js.Value,
			) interface{} {
				self.Value = self.value_select.GetJsValue("value").String()
				onchange()
				return false
			},
		),
	)

	self.Element = self.value_select

	return self
}

func (self *Select00) AssignSelf(ref **Select00) *Select00 {
	*ref = self
	return self
}

func (self *Select00) AppendOption(key, value string) {

	self.value_select.
		AppendChildren(
			self.etc.
				CreateElement("option").
				Set("value", key).
				AppendChildren(
					self.etc.CreateTextNode(value),
				),
		)
}

func (self *Select00) SetSelected(value string) {
	self.value_select.Set("value", value)
}
