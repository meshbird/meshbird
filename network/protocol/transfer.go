package protocol

import (
	"fmt"
	"io"
)

type (
	TransferMessage []byte
)

func NewTransferMessage(data []byte) *Packet {
	body := Body{
		Type:   TypeTransfer,
		Vector: randomBytes(16),
		Msg:    TransferMessage(data),
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

func (m TransferMessage) Bytes() []byte {
	return []byte(m)
}

func WriteEncodeTransfer(w io.Writer, data []byte) (err error) {
	logger.Debug("writing transfer message...")
	if err = EncodeAndWrite(w, NewTransferMessage(data)); err != nil {
		err = fmt.Errorf("error on write transfer message, %v", err)
	}
	return
}
