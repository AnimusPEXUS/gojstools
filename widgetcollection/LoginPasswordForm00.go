package widgetcollection

import (
	"syscall/js"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
)

type LoginPasswordForm00 struct {
	etc *elementtreeconstructor.ElementTreeConstructor

	put_submit_button bool

	onedit         func()
	onloginedit    func()
	onpasswordedit func()
	onsubmitclick  func()

	login_input    *elementtreeconstructor.ElementMutator
	password_input *elementtreeconstructor.ElementMutator

	Element *elementtreeconstructor.ElementMutator
}

func NewLoginPasswordForm00(
	etc *elementtreeconstructor.ElementTreeConstructor,
	preset_login string,
	preset_password string,
	put_submit_button bool,
	onedit func(),
	onloginedit func(),
	onpasswordedit func(),
	onsubmitclick func(),
) *LoginPasswordForm00 {

	self := &LoginPasswordForm00{
		etc:               etc,
		put_submit_button: put_submit_button,
		onedit:            onedit,
		onloginedit:       onloginedit,
		onpasswordedit:    onpasswordedit,
		onsubmitclick:     onsubmitclick,
	}

	self.Element = self.etc.CreateElement("div")

	login_input := self.etc.CreateElement("input").
		SetAttribute("type", "text").
		Set(
			"onchange",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					self.loginEdited()
					return false
				},
			),
		)
	self.login_input = login_input

	password_input := self.etc.CreateElement("input").
		SetAttribute("type", "password").
		Set(
			"onchange",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					self.passwordEdited()
					return false
				},
			),
		)
	self.password_input = password_input

	var table *elementtreeconstructor.ElementMutator

	self.Element.AppendChildren(
		self.etc.CreateElement("table").
			AssignSelf(&table).
			AppendChildren(
				self.etc.CreateElement("tr").AppendChildren(
					self.etc.CreateElement("td").AppendChildren(
						self.etc.CreateTextNode("Login"),
					),
					self.etc.CreateElement("td").AppendChildren(
						login_input,
					),
				),
				self.etc.CreateElement("tr").AppendChildren(
					self.etc.CreateElement("td").AppendChildren(
						self.etc.CreateTextNode("Password"),
					),
					self.etc.CreateElement("td").AppendChildren(
						password_input,
					),
				),
			),
	)

	self.SetLogin(preset_login)
	self.SetPassword(preset_password)

	if put_submit_button {
		table.AppendChildren(
			self.etc.CreateElement("tr").AppendChildren(
				self.etc.CreateElement("td").AppendChildren(
				// empty
				),
				self.etc.CreateElement("td").AppendChildren(
					self.etc.CreateElement("button").
						Set(
							"onclick",
							js.FuncOf(
								func(this js.Value, args []js.Value) interface{} {
									self.onsubmitclick()
									return false
								},
							),
						).AppendChildren(
						self.etc.CreateTextNode("Login"),
					),
				),
			),
		)
	}

	return self
}

func (self *LoginPasswordForm00) SetLogin(value string) {
	self.login_input.Set("value", value)
}

func (self *LoginPasswordForm00) SetPassword(value string) {
	self.password_input.Set("value", value)
}

func (self *LoginPasswordForm00) GetLogin() string {
	return self.login_input.GetJsValue("value").String()
}

func (self *LoginPasswordForm00) GetPassword() string {
	return self.password_input.GetJsValue("value").String()
}

func (self *LoginPasswordForm00) loginEdited() {
	if self.onloginedit != nil {
		self.onloginedit()
	}

	if self.onedit != nil {
		self.onedit()
	}
}

func (self *LoginPasswordForm00) passwordEdited() {
	if self.onpasswordedit != nil {
		self.onpasswordedit()
	}

	if self.onedit != nil {
		self.onedit()
	}
}
