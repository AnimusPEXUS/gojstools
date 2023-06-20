package ws

import (
	"errors"
	"io/ioutil"
	"log"
	"syscall/js"

	"github.com/AnimusPEXUS/gojstools/std/array"
	"github.com/AnimusPEXUS/gojstools/webapi/blob"
	"github.com/AnimusPEXUS/gojstools/webapi/events"
)

func GetByteSliceFromWSMessageEvent(
	event js.Value,
) (
	[]byte,
	error,
) {

	if !events.ValueIsInstanceOfMessageEvent(event) {
		return nil, errors.New("not a MessageEvent instance")
	}

	data := event.Get("data")
	if data.IsUndefined() {
		return nil, errors.New("no 'data' in JS value")
	}

	return GetByteSliceFromWSMessageEventData(data)
}

func GetByteSliceFromWSMessageEventData(
	event_data js.Value,
) (
	[]byte,
	error,
) {

	{
		bl, err := blob.NewBlobFromJSValue(event_data)

		if err != nil {
			goto next1
		}

		bl_rdr, err := bl.MakeReader()
		if err != nil {
			return nil, err
		}

		ret, err := ioutil.ReadAll(bl_rdr)
		if err != nil {
			return nil, err
		}

		log.Println("blob ReadAll result:", ret, string(ret))

		return ret, nil
	}

next1:

	{
		obj, err := array.NewArrayFromJSValue(event_data)

		if err != nil {
			goto next2
		}

		ret, err := obj.GetU8Bytes()
		if err != nil {
			return nil, err
		}

		return ret, nil
	}

next2:

	return nil, errors.New("couldn't determine event.data type")
}
