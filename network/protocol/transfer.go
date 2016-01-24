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

func ReadDecodeTransfer(r io.Reader) (TransferMessage, error) {
	log.Printf("Trying to read Transfer message...")

	transferPack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		log.Printf("Unable to decode package: %s", errDecode)
		return nil, fmt.Errorf("Error on read Transfer package: %v", errDecode)
	}

	if transferPack.Data.Type != TypeTransfer {
		return nil, fmt.Errorf("Got non Transfer message: %+v", transferPack)
	}

	log.Printf("Readed Transfer: %+v", transferPack.Data.Msg)

	return transferPack.Data.Msg.(TransferMessage), nil
}

func WriteEncodeTransfer(w io.Writer, data []byte) (err error) {
	log.Printf("Trying to write Transfer message...")
	if err = EncodeAndWrite(w, NewTransferMessage(data)); err != nil {
		err = fmt.Errorf("Error on write Transfer message: %v", err)
	}
	return
}
