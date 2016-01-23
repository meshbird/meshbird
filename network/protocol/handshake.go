package protocol

import (
	"io"
)

type (
	HandshakeMessage []byte
)

func NewHandshakePacket(key []byte) *Packet {
	body := Body{
		Type: TypeHandshake,
		Msg:  HandshakeMessage(key),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m HandshakeMessage) Len() uint16 {
	return uint16(len(m))
}

func (m HandshakeMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}
