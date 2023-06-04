package blob

import (
	"io"
	"syscall/js"

	"github.com/AnimusPEXUS/gojswebapi/array"
)

var _ io.Reader = &BlobReader{}

type BlobReader struct {
	blob   *Blob
	lenght int
	done   int
	EOF    bool
}

func NewBlobReader(blob *Blob) (*BlobReader, error) {

	length, err := blob.Size()
	if err != nil {
		return nil, err
	}

	self := &BlobReader{
		blob:   blob,
		lenght: length,
		done:   0,
	}

	return self, nil
}

func (self *BlobReader) Read(p []byte) (n int, err error) {
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

	bslc, err := self.blob.Slice(&[]int{self.done}[0], &[]int{end_index}[0], nil)
	if err != nil {
		n = 0 // TODO: is this correct?
		return
	}

	ab, err := bslc.ArrayBuffer()
	if err != nil {
		return
	}

	arr, err := array.NewArray(array.ArrayTypeUint8, ab.JSValue, nil, nil)
	if err != nil {
		return
	}

	// TODO: probably better error checking needed
	n = js.CopyBytesToGo(p, *arr.JSValue)
	return
}
