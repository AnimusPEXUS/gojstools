package arraybuffer

import (
	"io"
	"syscall/js"

	"github.com/AnimusPEXUS/gojswebapi/array"
)

var _ io.Reader = &ArrayBufferReader{}

type ArrayBufferReader struct {
	ab     *ArrayBuffer
	lenght int
	done   int
	EOF    bool
}

func NewArrayBufferReader(ab *ArrayBuffer) (*ArrayBufferReader, error) {

	length, err := ab.Len()
	if err != nil {
		return nil, err
	}

	self := &ArrayBufferReader{
		ab:     ab,
		lenght: length,
		done:   0,
	}

	return self, nil
}

func (self *ArrayBufferReader) Read(p []byte) (n int, err error) {
	len_p := len(p)

	var end_index int

	{
		selfdonelenp := self.done + len_p
		if selfdonelenp > self.lenght {
			end_index = self.done - self.lenght
			self.EOF = true

		} else {
			end_index = selfdonelenp
		}
	}

	slc, err := self.ab.Slice(self.done, &[]int{end_index}[0], nil)
	if err != nil {
		n = 0 // TODO: is this correct?
		return
	}

	arr, err := array.NewArray(array.ArrayTypeUint8, slc.JSValue, nil, nil)
	if err != nil {
		return
	}

	// TODO: probably better error checking needed
	n = js.CopyBytesToGo(p, *arr.JSValue)

	if self.EOF {
		err = io.EOF
	}

	return
}
