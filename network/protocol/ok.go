package protocol

import (
	"io"
	"net"
)

var (
	onMessage = []byte{'O', 'K'}
)

type (
	OkMessage []byte
)

func NewOkMessage(privateIP net.IP) *Packet {
	body := Body{
		Type: TypeOk,
		Msg:  HandshakeMessage(append(onMessage, privateIP...)),
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

func (o OkMessage) PrivateIP() net.IP {
	return net.IP(o[2:])
}
