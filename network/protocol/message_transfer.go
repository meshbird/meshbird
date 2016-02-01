package protocol

import (
	"fmt"
	"io"
)

func NewTransferMessage(data []byte) *Record {
	rec := NewRecord(TypeTransfer, data)
	rec.Vector = randomBytes(16)
	return rec
}

func WriteEncodeTransfer(w io.Writer, data []byte) (err error) {
	logger.Debug("writing transfer message...")
	if err = EncodeAndWrite(w, NewTransferMessage(data)); err != nil {
		err = fmt.Errorf("error on write transfer message, %v", err)
	}
	return
}
