package protocol

import (
	"fmt"
	"io"
	"log"
)

type (
	TransferMessage []byte
)

func NewTransferMessage(data []byte) *Packet {
	body := Body{
		Type: TypeTransfer,
		Msg:  TransferMessage(data),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m TransferMessage) Len() uint16 {
	return uint16(len(m))
}

func (m TransferMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}

func WriteEncodeTransfer(w io.Writer, data []byte) (err error) {
	log.Printf("Trying to write Transfer message...")
	if err = EncodeAndWrite(w, NewTransferMessage(data)); err != nil {
		err = fmt.Errorf("Error on write Transfer message: %v", err)
	}
	return
}
