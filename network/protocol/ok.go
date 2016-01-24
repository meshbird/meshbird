package protocol

import (
	"fmt"
	"io"
	"log"
)

var (
	onMessage = []byte{'O', 'K'}
)

type (
	OkMessage []byte
)

func NewOkMessage() *Packet {
	body := Body{
		Type: TypeOk,
		Msg:  OkMessage(onMessage),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (o OkMessage) Len() uint16 {
	return uint16(len(o))
}

func (o OkMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(o)
	return int64(n), err
}

func ReadDecodeOk(r io.Reader) (OkMessage, error) {
	log.Printf("Trying to read OK message...")

	okPack, errDecode := ReadAndDecode(r, 6)
	if errDecode != nil {
		log.Printf("Unable to decode package: %s", errDecode)
		return nil, fmt.Errorf("Error on read OK package: %v", errDecode)
	}

	if okPack.Data.Type != TypeOk {
		return nil, fmt.Errorf("Got non OK message: %+v", okPack)
	}

	log.Printf("Readed OK: %+v", okPack.Data.Msg)

	return okPack.Data.Msg.(OkMessage), nil
}

func WriteEncodeOk(w io.Writer) (err error) {
	log.Printf("Trying to write OK message...")
	if err = EncodeAndWrite(w, NewOkMessage()); err != nil {
		err = fmt.Errorf("Error on write OK message: %v", err)
	}
	return
}
