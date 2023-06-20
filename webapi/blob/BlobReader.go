package blob

import (
	"io"
	"syscall/js"

	"github.com/AnimusPEXUS/gojstools/std/array"
)

var _ io.Reader = &BlobReader{}

type BlobReader struct {
	blob        *Blob
	lenght      int
	start_index int
	EOF         bool
}

func NewBlobReader(blob *Blob) (*BlobReader, error) {

	length, err := blob.Size()
	if err != nil {
		return nil, err
	}

	self := &BlobReader{
		blob:        blob,
		lenght:      length,
		start_index: 0,
	}

	return self, nil
}

func (self *BlobReader) Read(p []byte) (n int, err error) {

	// todo: testing required

	len_p := len(p)

	if self.EOF {
		return 0, io.EOF
	}

	if len_p == 0 {
		// TODO: or should we return error?
		return 0, nil
	}

	var end_index int

	{
		self_start_index_p_lenp := self.start_index + len_p
		if self_start_index_p_lenp > self.lenght {
			end_index = self.lenght
			self.EOF = true
		} else {
			end_index = self_start_index_p_lenp
		}
	}

	if self.start_index == end_index && self.EOF {
		return 0, io.EOF
	}

	defer func() {
		self.start_index = end_index
	}()

	bslc, err := self.blob.Slice(
		&[]int{self.start_index}[0],
		&[]int{end_index}[0],
		nil,
	)
	if err != nil {
		n = 0 // TODO: is this correct?
		return
	}

	ab, err := bslc.ArrayBuffer()
	if err != nil {
		return
	}

	arr, err := array.NewArray(
		array.ArrayTypeUint8,
		ab.JSValue,
		nil,
		nil,
	)
	if err != nil {
		return
	}

	// TODO: probably better error checking needed
	n = js.CopyBytesToGo(p, arr.JSValue)
	return
}
