package widgetcollection

import (
	"errors"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
)

type PageId = uint

type Pager00Page struct {
	PageId  PageId
	Element *elementtreeconstructor.ElementMutator
}

type Pager00Settings struct {
	Pages          []*Pager00Page
	Preselected    PageId
	DisplayElement *elementtreeconstructor.ElementMutator
}

type Pager00 struct {
	etc *elementtreeconstructor.ElementTreeConstructor

	settings *Pager00Settings

	Element *elementtreeconstructor.ElementMutator
}

func NewPager00(
	etc *elementtreeconstructor.ElementTreeConstructor,
	settings *Pager00Settings,
) *Pager00 {

	self := &Pager00{
		settings: settings,
		Element:  settings.DisplayElement,
	}

	self.SwitchPage(self.settings.Preselected)

	return self
}

func (self *Pager00) SwitchPage(id PageId) error {

	var n *elementtreeconstructor.ElementMutator

	for _, i := range self.settings.Pages {
		if i.PageId == id {
			n = i.Element
			break
		}
	}

	if n == nil {
		return errors.New("couldn't find page with this id")
	}

	self.Element.RemoveAllChildren()
	self.Element.AppendChildren(n)

	return nil
}
