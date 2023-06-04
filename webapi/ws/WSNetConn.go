package ws

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"syscall/js"
	"time"

	gojstoolsutils "github.com/AnimusPEXUS/gojstools/utils"
	"github.com/AnimusPEXUS/gojswebapi/array"
	wasmtools_arraybuffer "github.com/AnimusPEXUS/gojswebapi/arraybuffer"
	wasmtools_blob "github.com/AnimusPEXUS/gojswebapi/blob"
	"github.com/AnimusPEXUS/gojswebapi/events"
	"github.com/AnimusPEXUS/utils/worker"
)

type EmptyStruct struct{}

var _ net.Conn = &WSNetConn{}

// TODO: redo. this is fast temporary copypasta
type websocketAddr struct{}

func (a websocketAddr) Network() string {
	return "websocket"
}

func (a websocketAddr) String() string {
	return "websocket/unknown-addr"
}

type WSNetConnOptions struct {
	WS          *WS
	CloseCode   *int
	CloseReason *string
	OnError     func(error)
}

type WSNetConn struct {
	options *WSNetConnOptions

	read_buffer *bytes.Buffer

	inbound_messages       []js.Value
	inbound_messages_mutex sync.Mutex
	inbound_worker         *worker.Worker
	inbound_signal         chan EmptyStruct
	// outbound_messages       []js.Value
	// outbound_messages_mutex sync.Mutex
	// outbound_worker         *worker.Worker
	// outbound_signal         chan EmptyStruct

	central_worker *worker.Worker

	// isopen bool

	WSError error
}

func NewWSNetConn(options *WSNetConnOptions) *WSNetConn {
	self := &WSNetConn{
		options:          options,
		read_buffer:      nil,
		inbound_messages: make([]js.Value, 0),
		// outbound_messages: make([]js.Value, 0),
		inbound_signal: make(chan EmptyStruct),
		// outbound_signal:   make(chan EmptyStruct),
	}

	self.inbound_worker = worker.New(self.inboundWorkerThread)
	// self.outbound_worker = worker.New(self.outboundWorkerThread)
	self.central_worker = worker.New(self.centralWorkerThread)

	return self
}

func (self *WSNetConn) InstallEventHandlersIntoWS() error {
	err := self.options.WS.SetOnMessage(self.onMessage)
	if err != nil {
		return err
	}

	err = self.options.WS.SetOnError(self.onError)
	if err != nil {
		return err
	}

	// err = self.options.WS.SetOnClose(self.onClose)
	// if err != nil {
	// 	return err
	// }

	// err = self.options.WS.SetOnOpen(self.onOpen)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (self *WSNetConn) isopen() bool {
	state, err := self.options.WS.ReadyStateGet()
	return err == nil && state == WSReadyState_OPEN
}

// func (self *WSNetConn) onOpen(event js.Value) {
// 	log.Println("WSNetConn.onOpen")
// 	self.isopen = true
// }

// func (self *WSNetConn) onClose(event js.Value) {
// 	log.Println("WSNetConn.onClose")
// 	self.isopen = false
// }

func (self *WSNetConn) onMessage(event *events.MessageEvent) {
	self.inbound_messages_mutex.Lock()
	defer self.inbound_messages_mutex.Unlock()
	log.Println("WSNetConn.onMessage")

	// TODO: add error handling?

	js_data, err := event.GetData()
	if err != nil {
		log.Println("error:", err.Error())
		return
	}
	self.inbound_messages = append(self.inbound_messages, js_data)

	if len(self.inbound_signal) < 1 {
		self.inbound_signal <- EmptyStruct{}
	}
}

func (self *WSNetConn) onError(event *events.ErrorEvent) {
	log.Println("WSNetConn.onError")
	if self.options.OnError != nil {
		msg, err := event.GetMessage()
		if err != nil {
			log.Println("error:", err.Error())
			return
		}
		self.options.OnError(errors.New(msg))
	}
}

func (self *WSNetConn) GetWorker() worker.WorkerI {
	return self.central_worker
}

func (self *WSNetConn) inboundWorkerThread(
	set_starting func(),
	set_working func(),
	set_stopping func(),
	set_stopped func(),
	is_stop_flag func() bool,
) {
	set_starting()
	defer func() {
		set_stopping()
		go self.central_worker.Stop()
		set_stopped()
	}()

	stop_signal := make(chan EmptyStruct)

	go func() {
		if is_stop_flag() {
			stop_signal <- EmptyStruct{}
			time.Sleep(time.Second)
		}
	}()

	for {
		select {
		case <-stop_signal:
			break
		case <-self.inbound_signal:
			for {
				if is_stop_flag() {
					return
				}
				err := self.processNextInboundMessage()
				if err != nil {
					self.WSError = err
					self.central_worker.Stop()
					return
				}
				if len(self.inbound_messages) == 0 {
					break
				}
			}
		}
	}
}

// func (self *WSNetConn) outboundWorkerThread(
// 	set_starting func(),
// 	set_working func(),
// 	set_stopping func(),
// 	set_stopped func(),
// 	is_stop_flag func() bool,
// ) {
// 	set_starting()
// 	defer func() {
// 		set_stopping()
// 		go self.central_worker.Stop()
// 		set_stopped()
// 	}()

// 	stop_signal := make(chan EmptyStruct)

// 	go func() {
// 		if is_stop_flag() {
// 			stop_signal <- EmptyStruct{}
// 			time.Sleep(time.Second)
// 		}
// 	}()

// 	for {
// 		select {
// 		case <-stop_signal:
// 			break
// 		case <-self.outbound_signal:
// 			for {
// 				if is_stop_flag() {
// 					return
// 				}
// 				err := self.processNextOutboundMessage()
// 				if err != nil {
// 					self.WSError = err
// 					self.central_worker.Stop()
// 					return
// 				}
// 				if len(self.outbound_messages) == 0 {
// 					break
// 				}
// 			}
// 		}
// 	}
// }

func (self *WSNetConn) centralWorkerThread(

	set_starting func(),
	set_working func(),
	set_stopping func(),
	set_stopped func(),

	is_stop_flag func() bool,

) {
	set_starting()
	defer func() {
		set_stopping()
		s1 := make(chan worker.WorkerControlChanResult)
		s2 := make(chan worker.WorkerControlChanResult)
		go func() { s1 <- self.inbound_worker.Stop() }()
		// go func() { s2 <- self.outbound_worker.Stop() }()
		<-s1
		<-s2
		set_stopped()
	}()

	{
		s1 := make(chan worker.WorkerControlChanResult)
		s2 := make(chan worker.WorkerControlChanResult)
		go func() { s1 <- self.inbound_worker.Start() }()
		// go func() { s2 <- self.outbound_worker.Start() }()
		<-s1
		<-s2
	}

	set_working()
	for {
		if is_stop_flag() {
			break
		}
		time.Sleep(time.Second)
	}
}

func (self *WSNetConn) processNextInboundMessage() error {
	self.inbound_messages_mutex.Lock()
	defer self.inbound_messages_mutex.Unlock()

	js_data := self.inbound_messages[0]
	// self.inbound_messages = append(self.inbound_messages[0:0], self.inbound_messages[1:]...)
	self.inbound_messages = self.inbound_messages[1:]

	var re io.Reader

	{
		yes, err := wasmtools_blob.IsBlob(js_data)
		if err != nil {
			return err
		}

		if yes {
			res, err := wasmtools_blob.NewBlobFromJSValue(js_data)
			if err != nil {
				return err
			}
			re, err = res.MakeReader()
			if err != nil {
				return err
			}
			goto work_result
		}
	}

	{
		yes, err := wasmtools_arraybuffer.IsArrayBuffer(js_data)
		if err != nil {
			return err
		}

		if yes {
			res, err := wasmtools_arraybuffer.NewArrayBufferFromJSValue(js_data)
			if err != nil {
				return err
			}
			re, err = res.MakeReader()
			if err != nil {
				return err
			}
			goto work_result
		}
	}

	return errors.New("unknown error")

work_result:
	_, err := io.Copy(self.read_buffer, re)
	if err != nil {
		return err
	}

	return nil
}

// func (self *WSNetConn) processNextOutboundMessage() error {

// 	self.outbound_messages_mutex.Lock()
// 	defer self.outbound_messages_mutex.Unlock()

// 	if len(self.outbound_messages) == 0 {
// 		return nil
// 	}

// 	msg := self.outbound_messages[0]

// ok_exit:

// 	self.outbound_messages = self.outbound_messages[1:]

// 	return nil
// }

func (self *WSNetConn) Read(b []byte) (n int, err error) {
	log.Println("WSNetConn Read")
	defer log.Println("WSNetConn Read exit", "n:", n, "err:", err)

	if !self.isopen() {
		return 0, os.ErrClosed
	}

	// make_read:
	if self.read_buffer != nil {
		n, err = self.read_buffer.Read(b)
		if self.read_buffer.Len() == 0 {
			self.read_buffer = nil
		}
		return
	}

	// magic!
	return 666, errors.New("Something strange! Who will you call?")
}

func (self *WSNetConn) Write(b []byte) (n int, err error) {

	initial_b_size := len(b)

	log.Println("WSNetConn Write", b)
	defer log.Println("WSNetConn Write exit", "n:", n, "err:", err)

	if !self.isopen() {
		return 0, os.ErrClosed
	}

	log.Println("got some bytes to write:", b)

	bval, err := array.NewArray(
		array.ArrayTypeUint8,
		gojstoolsutils.JSValueLiteralToPointer(js.ValueOf(len(b))),
		nil,
		nil,
	)
	if err != nil {
		return
	}

	js.CopyBytesToJS(*bval.JSValue, b)

	log.Println("sending...")
	err = self.options.WS.Send(bval.JSValue)
	log.Println("   sending result:", err)
	log.Println("  n:", n, "initial_b_size:", initial_b_size)
	n = initial_b_size
	log.Println("  n:", n, "initial_b_size:", initial_b_size)
	return
}

func (self *WSNetConn) Close() error {
	return self.options.WS.Close(self.options.CloseCode, self.options.CloseReason)
}

func (self *WSNetConn) LocalAddr() net.Addr {
	return websocketAddr{}
}

func (self *WSNetConn) RemoteAddr() net.Addr {
	return websocketAddr{}
}

func (self *WSNetConn) SetDeadline(t time.Time) error {
	return nil
}

func (self *WSNetConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (self *WSNetConn) SetWriteDeadline(t time.Time) error {
	return nil
}
