package protocol

import "io"

type (
	OkMessage []byte
)

func NewOkMessage() *Packet {
	body := Body{
		Type: TypeOk,
		Msg:  HandshakeMessage([]byte{'O', 'K'}),
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
