package protocol

import (
	"fmt"
	"io"
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
	logger.Debug("reading ok message...")

	okPack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.Error("error on package decode, %v", errDecode)
		return nil, fmt.Errorf("error on read ok package, %v", errDecode)
	}

	if okPack.Data.Type != TypeOk {
		return nil, fmt.Errorf("non ok message received, %+v", okPack)
	}

	logger.Debug("message, %v", okPack.Data.Msg)
	return okPack.Data.Msg.(OkMessage), nil
}

func WriteEncodeOk(w io.Writer) (err error) {
	logger.Debug("writing ok message...")
	if err = EncodeAndWrite(w, NewOkMessage()); err != nil {
		err = fmt.Errorf("error on write ok message, %v", err)
	}
	return
}
