package ws

import (
	"errors"
	"syscall/js"

	// gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
	"github.com/AnimusPEXUS/gojswebapi/events"
	utils_panic "github.com/AnimusPEXUS/utils/panic"
)

type WSReadyState int

const (
	WSReadyState_CONNECTING WSReadyState = 0
	WSReadyState_OPEN                    = 1
	WSReadyState_CLOSING                 = 2
	WSReadyState_CLOSED                  = 3
)

// if both url and js_value are specified, js_value is used
type WSOptions struct {
	URL *string
	// to use existing ws
	JSValue   js.Value
	Protocols []string

	OnClose   func(*events.CloseEvent)   // function(event)
	OnError   func(*events.ErrorEvent)   // function(event)
	OnMessage func(*events.MessageEvent) // function(event)
	OnOpen    func(*events.Event)        // function(event)
}

type WS struct {
	JSValue js.Value
	options *WSOptions
}

func NewWS(options *WSOptions) (*WS, error) {

	if (options.JSValue == nil && options.URL == nil) ||
		(options.JSValue != nil && options.URL != nil) {
		panic("existing socket _or_ url must be supplied")
	}

	self := &WS{
		options: options,
	}

	if options.JSValue != nil {
		self.JSValue = options.JSValue
		options.URL = &([]string{self.JSValue.Get("url").String()}[0])
	} else {
		wsoc_constr := js.Global().Get("WebSocket")
		if wsoc_constr.IsUndefined() {
			return nil, errors.New("WebSocket is undefined")
		}
		url := *options.URL
		wsoc := wsoc_constr.New(url, js.Undefined()) // TODO: options.Protocols
		self.JSValue = &wsoc
		// options.JSValue = self.JSValue
	}

	err := self.SetOnOpen(options.OnOpen)
	if err != nil {
		return nil, err
	}

	err = self.SetOnClose(options.OnClose)
	if err != nil {
		return nil, err
	}

	err = self.SetOnMessage(options.OnMessage)
	if err != nil {
		return nil, err
	}

	err = self.SetOnError(options.OnError)
	if err != nil {
		return nil, err
	}

	return self, nil
}

func (self *WS) SetOnOpen(f func(*events.Event)) (err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()

	if f == nil {
		self.JSValue.Set("onopen", js.Undefined())
		self.options.OnOpen = nil
		return
	}

	self.options.OnOpen = f

	self.JSValue.Set(
		"onopen",
		js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				if self.options.OnOpen != nil {
					ev, err := events.NewEventFromJSValue(&args[0])
					if err != nil {
						return nil
					}
					go self.options.OnOpen(ev)
				} else {
					self.SetOnOpen(nil)
				}
				return nil
			},
		),
	)
	return nil
}

func (self *WS) SetOnClose(f func(*events.CloseEvent)) (err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()

	if f == nil {
		self.JSValue.Set("onclose", js.Undefined())
		self.options.OnClose = nil
		return
	}

	self.options.OnClose = f

	self.JSValue.Set(
		"onclose",
		js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				if self.options.OnClose != nil {
					ev, err := events.NewCloseEventFromJSValue(&args[0])
					if err != nil {
						return nil
					}
					go self.options.OnClose(ev)
				} else {
					self.SetOnClose(nil)
				}
				return nil
			},
		),
	)
	return nil
}

func (self *WS) SetOnMessage(f func(*events.MessageEvent)) (err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()

	if f == nil {
		self.JSValue.Set("onmessage", js.Undefined())
		self.options.OnMessage = nil
		return
	}

	self.options.OnMessage = f

	self.JSValue.Set(
		"onmessage",
		js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				if self.options.OnMessage != nil {
					ev, err := events.NewMessageEventFromJSValue(&args[0])
					if err != nil {
						return nil
					}
					go self.options.OnMessage(ev)
				} else {
					self.SetOnMessage(nil)
				}
				return nil
			},
		),
	)
	return nil
}

func (self *WS) SetOnError(f func(*events.ErrorEvent)) (err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()

	if f == nil {
		self.JSValue.Set("onerror", js.Undefined())
		self.options.OnError = nil
		return
	}

	self.options.OnError = f

	self.JSValue.Set(
		"onerror",
		js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				if self.options.OnError != nil {
					ev, err := events.NewErrorEventFromJSValue(&args[0])
					if err != nil {
						return nil
					}
					go self.options.OnError(ev)
				} else {
					self.SetOnError(nil)
				}
				return nil
			},
		),
	)
	return nil
}

func (self *WS) Close(code *int, reason *string) (err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()

	if reason != nil && code == nil {
		return errors.New("reason can't be specified without code")
	}

	var args []interface{}

	if code != nil {
		args = append(args, *code)
		if reason != nil {
			args = append(args, *reason)
		}
	}

	self.JSValue.Call("close", args...)
	return nil
}

func (self *WS) Send(value js.Value) (err error) {
	// log.Print("WS Send called")
	defer func() {
		err = utils_panic.PanicToError()
	}()

	state, _ := self.ReadyStateGet()
	if state != WSReadyState_OPEN {
		err = errors.New("ws: socket is not open")
		return
	}
	// log.Println("ws state", state)
	// url, _ := self.URLGet()
	// log.Println("ws url", url)

	// v := *value

	// log.Println("value v:", v.Call("toString").String())

	self.JSValue.Call("send", *value)
	return
}

///////////////// properties

func (self *WS) BinaryTypeGet() (ret string, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.JSValue.Get("binaryType").String()
	return
}

func (self *WS) BinaryTypeSet(value string) (err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	self.JSValue.Set("binaryType", value)
	return
}

func (self *WS) BufferedAmountGet() (ret int, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.JSValue.Get("bufferedAmount").Int()
	return
}

func (self *WS) ProtocolGet() (ret string, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.JSValue.Get("protocol").String()
	return
}

func (self *WS) ReadyStateGet() (ret WSReadyState, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = WSReadyState(self.JSValue.Get("readyState").Int())
	return
}

func (self *WS) URLGet() (ret string, err error) {
	defer func() {
		err = utils_panic.PanicToError()
	}()
	ret = self.JSValue.Get("url").String()
	return
}
