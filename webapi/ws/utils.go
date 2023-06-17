package ws

import (
	"errors"
	"syscall/js"

	"github.com/AnimusPEXUS/gojstools/webapi/array"
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
		obj, err := array.NewArrayFromJSValue(event_data)

		if err != nil {
			goto next1
		}

		ret, err := obj.GetU8Bytes()
		if err != nil {
			return nil, err
		}

		return ret, nil
	}

next1:

	// 	{
	// 		_, err := blob.NewBlobFromJSValue(event_data)

	// 		if err != nil {
	// 			goto next2
	// 		}

	// 		return ret, nil
	// 	}

	// next2:

	return nil, errors.New("couldn't determine event.data type")
}
