package protocol

import (
	"fmt"
	"io"
)

var (
	onMessage = []byte{'O', 'K'}
)

func NewOkMessage() *Record {
	return NewRecord(TypeOk, onMessage)
}

func ReadDecodeOk(r io.Reader) ([]byte, error) {
	return readDecodeMsg(r, TypeOk)
}

func WriteEncodeOk(w io.Writer) (err error) {
	logger.Debug("writing ok message...")
	if err = EncodeAndWrite(w, NewOkMessage()); err != nil {
		err = fmt.Errorf("error on write ok message, %v", err)
	}
	return
}
